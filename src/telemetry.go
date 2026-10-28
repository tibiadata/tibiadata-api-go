package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

const (
	telemetryEnabledEnv = "TIBIADATA_OTEL_ENABLED"
	telemetryService    = "tibiadata-api-go"
)

var (
	tibiaDataTelemetryEnabled bool
	tibiaDataTracingEnabled   bool
	tibiaDataMetricsEnabled   bool
	tibiaDataLogsEnabled      bool
	tibiaDataTelemetryStop    = func(context.Context) error { return nil }
)

// initTibiaDataTelemetry configures global OpenTelemetry providers when
// TIBIADATA_OTEL_ENABLED is true. Exporters and SDK settings are read from the
// standard OTEL_* environment variables.
func initTibiaDataTelemetry(ctx context.Context) (func(context.Context) error, error) {
	noopShutdown := func(context.Context) error { return nil }

	tibiaDataTelemetryEnabled = getEnvAsBool(telemetryEnabledEnv, false)
	tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = false, false, false
	if !tibiaDataTelemetryEnabled {
		slog.Info("TibiaData API opentelemetry", "enabled", false)
		return noopShutdown, nil
	}

	res, err := tibiaDataTelemetryResource(ctx)
	if err != nil {
		return noopShutdown, err
	}

	var shutdowns []func(context.Context) error
	previousTracerProvider := otel.GetTracerProvider()
	previousMeterProvider := otel.GetMeterProvider()
	previousLoggerProvider := otel.GetLoggerProvider()
	var tracerProviderConfigured, meterProviderConfigured, loggerProviderConfigured bool
	shutdown := func(ctx context.Context) error {
		var errs []error
		for _, s := range shutdowns {
			errs = append(errs, s(ctx))
		}
		return errors.Join(errs...)
	}
	fail := func(err error) (func(context.Context) error, error) {
		err = errors.Join(err, shutdown(ctx))
		if tracerProviderConfigured {
			otel.SetTracerProvider(previousTracerProvider)
		}
		if meterProviderConfigured {
			otel.SetMeterProvider(previousMeterProvider)
		}
		if loggerProviderConfigured {
			otel.SetLoggerProvider(previousLoggerProvider)
		}
		tibiaDataTelemetryEnabled = false
		tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = false, false, false
		return noopShutdown, err
	}

	// Signals whose exporter is "none" keep the global no-op provider, so no
	// spans or measurements are created only to be discarded.
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return fail(err)
	}
	if !autoexport.IsNoneSpanExporter(spanExporter) {
		tracerProvider := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(spanExporter),
			sdktrace.WithResource(res),
		)
		otel.SetTracerProvider(tracerProvider)
		tracerProviderConfigured = true
		tibiaDataTracingEnabled = true
		shutdowns = append(shutdowns, tracerProvider.Shutdown)
	}

	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return fail(err)
	}
	if !autoexport.IsNoneMetricReader(metricReader) {
		meterProvider := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(metricReader),
			sdkmetric.WithResource(res),
		)
		otel.SetMeterProvider(meterProvider)
		meterProviderConfigured = true
		tibiaDataMetricsEnabled = true
		shutdowns = append(shutdowns, meterProvider.Shutdown)
	}

	logExporter, err := autoexport.NewLogExporter(ctx)
	if err != nil {
		return fail(err)
	}
	if !autoexport.IsNoneLogExporter(logExporter) {
		loggerProvider := sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
			sdklog.WithResource(res),
		)
		otel.SetLoggerProvider(loggerProvider)
		loggerProviderConfigured = true
		configureTibiaDataOTelLogging(loggerProvider)
		tibiaDataLogsEnabled = true
		shutdowns = append(shutdowns, loggerProvider.Shutdown)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		// Avoid sending exporter errors back through the same OTel log pipeline.
		logOpenTelemetryError(err)
	}))

	slog.Info("TibiaData API opentelemetry",
		"enabled", true,
		"traces_exporter", envOrDefault("OTEL_TRACES_EXPORTER", "otlp"),
		"metrics_exporter", envOrDefault("OTEL_METRICS_EXPORTER", "otlp"),
		"logs_exporter", envOrDefault("OTEL_LOGS_EXPORTER", "otlp"),
	)

	return shutdown, nil
}

func tibiaDataTelemetryResource(ctx context.Context) (*resource.Resource, error) {
	// Attributes from OTEL_SERVICE_NAME and OTEL_RESOURCE_ATTRIBUTES take precedence.
	defaults := resource.NewSchemaless(
		attribute.String("service.name", telemetryService),
		attribute.String("service.version", TibiaDataBuildRelease),
	)
	detected, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
	)
	if err != nil && detected == nil {
		return nil, err
	}
	return resource.Merge(defaults, detected)
}

func envOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

type skipTraceKey struct{}

// tibiaDataTelemetryActive reports whether request tracing or metrics should
// be installed. Log export does not require Gin instrumentation.
func tibiaDataTelemetryActive() bool {
	return tibiaDataTelemetryEnabled && (tibiaDataTracingEnabled || tibiaDataMetricsEnabled)
}

// ginTelemetryMiddleware instruments requests with spans and HTTP server
// metrics. Probe endpoints are recorded in metrics but not traced.
func ginTelemetryMiddleware() gin.HandlerFunc {
	otelMiddleware := otelgin.Middleware(
		telemetryService,
		otelgin.WithTracerProvider(probeFilteringTracerProvider{otel.GetTracerProvider()}),
	)
	if !tibiaDataTracingEnabled {
		return otelMiddleware
	}
	return func(c *gin.Context) {
		if isProbePath(c.Request.URL.Path) {
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), skipTraceKey{}, true))
		}
		otelMiddleware(c)
	}
}

// probeFilteringTracerProvider returns non-recording spans for requests whose
// context is marked with skipTraceKey, so otelgin still records their metrics.
type probeFilteringTracerProvider struct {
	trace.TracerProvider
}

func (p probeFilteringTracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return probeFilteringTracer{Tracer: p.TracerProvider.Tracer(name, opts...)}
}

type probeFilteringTracer struct {
	trace.Tracer
}

var noopTracer = noop.NewTracerProvider().Tracer("")

func (t probeFilteringTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if skip, _ := ctx.Value(skipTraceKey{}).(bool); skip {
		return noopTracer.Start(ctx, name, opts...)
	}
	return t.Tracer.Start(ctx, name, opts...)
}

func setupTibiaDataTelemetry() {
	shutdown, err := initTibiaDataTelemetry(context.Background())
	if err != nil {
		tibiaDataTelemetryEnabled = false
		slog.Error("TibiaData API opentelemetry setup failed, continuing without telemetry", "error", err)
		return
	}
	tibiaDataTelemetryStop = shutdown
}

// shutdownTibiaDataTelemetry flushes and stops the telemetry providers.
func shutdownTibiaDataTelemetry() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tibiaDataTelemetryStop(ctx); err != nil {
		slog.Warn("TibiaData API opentelemetry shutdown failed", "error", err)
	}
}
