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
	tibiaDataTelemetryStop    = func(context.Context) error { return nil }
)

// initTibiaDataTelemetry configures the global OpenTelemetry tracer and meter
// providers when TIBIADATA_OTEL_ENABLED is true. Exporters and SDK settings are
// read from the standard OTEL_* environment variables.
func initTibiaDataTelemetry(ctx context.Context) (func(context.Context) error, error) {
	noopShutdown := func(context.Context) error { return nil }

	tibiaDataTelemetryEnabled = getEnvAsBool(telemetryEnabledEnv, false)
	if !tibiaDataTelemetryEnabled {
		slog.Info("TibiaData API opentelemetry", "enabled", false)
		return noopShutdown, nil
	}

	res, err := tibiaDataTelemetryResource(ctx)
	if err != nil {
		return noopShutdown, err
	}

	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return noopShutdown, err
	}
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(spanExporter),
		sdktrace.WithResource(res),
	)

	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return noopShutdown, errors.Join(err, tracerProvider.Shutdown(ctx))
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricReader),
		sdkmetric.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Warn("TibiaData API opentelemetry error", "error", err)
	}))

	slog.Info("TibiaData API opentelemetry",
		"enabled", true,
		"traces_exporter", envOrDefault("OTEL_TRACES_EXPORTER", "otlp"),
		"metrics_exporter", envOrDefault("OTEL_METRICS_EXPORTER", "otlp"),
	)

	return func(ctx context.Context) error {
		return errors.Join(tracerProvider.Shutdown(ctx), meterProvider.Shutdown(ctx))
	}, nil
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

// ginTelemetryMiddleware instruments requests with spans and HTTP server
// metrics. Probe endpoints are recorded in metrics but not traced.
func ginTelemetryMiddleware() gin.HandlerFunc {
	otelMiddleware := otelgin.Middleware(
		telemetryService,
		otelgin.WithTracerProvider(probeFilteringTracerProvider{otel.GetTracerProvider()}),
	)
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
