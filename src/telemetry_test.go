package main

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelLog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

type recordingLogExporter struct {
	mu      sync.Mutex
	records []sdklog.Record
}

func (e *recordingLogExporter) Export(_ context.Context, records []sdklog.Record) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, record := range records {
		e.records = append(e.records, record.Clone())
	}
	return nil
}

func (e *recordingLogExporter) Shutdown(context.Context) error   { return nil }
func (e *recordingLogExporter) ForceFlush(context.Context) error { return nil }

func (e *recordingLogExporter) exportedRecords() []sdklog.Record {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]sdklog.Record(nil), e.records...)
}

// setupTestTelemetry installs in-memory tracer and meter providers as the
// global providers and restores the previous state when the test ends.
func setupTestTelemetry(t *testing.T) (*tracetest.SpanRecorder, *sdkmetric.ManualReader) {
	t.Helper()

	recorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	prevTP, prevMP, prevProp := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetTextMapPropagator()
	prevEnabled, prevTracing, prevMetrics, prevLogs := tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = true, true, true, false

	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetMeterProvider(prevMP)
		otel.SetTextMapPropagator(prevProp)
		tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = prevEnabled, prevTracing, prevMetrics, prevLogs
		_ = tracerProvider.Shutdown(context.Background())
		_ = meterProvider.Shutdown(context.Background())
	})

	return recorder, reader
}

func requestDurationCount(t *testing.T, reader *sdkmetric.ManualReader, route string) uint64 {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	var count uint64
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			if m.Name != "http.server.request.duration" {
				continue
			}
			hist, ok := m.Data.(metricdata.Histogram[float64])
			require.True(t, ok)
			for _, dp := range hist.DataPoints {
				if v, ok := dp.Attributes.Value("http.route"); ok && v.AsString() == route {
					count += dp.Count
				}
			}
		}
	}
	return count
}

func TestInitTibiaDataTelemetryDisabledByDefault(t *testing.T) {
	t.Setenv(telemetryEnabledEnv, "")
	prevTP := otel.GetTracerProvider()

	shutdown, err := initTibiaDataTelemetry(context.Background())

	require.NoError(t, err)
	assert.False(t, tibiaDataTelemetryEnabled)
	assert.Equal(t, prevTP, otel.GetTracerProvider())
	assert.NoError(t, shutdown(context.Background()))
}

func TestInitTibiaDataTelemetryEnabled(t *testing.T) {
	for _, test := range []struct {
		traces, metrics, logs             string
		wantTraces, wantMetrics, wantLogs bool
	}{
		{"console", "none", "none", true, false, false},
		{"none", "console", "none", false, true, false},
	} {
		t.Run(test.traces+"/"+test.metrics+"/"+test.logs, func(t *testing.T) {
			t.Setenv(telemetryEnabledEnv, "true")
			t.Setenv("OTEL_TRACES_EXPORTER", test.traces)
			t.Setenv("OTEL_METRICS_EXPORTER", test.metrics)
			t.Setenv("OTEL_LOGS_EXPORTER", test.logs)
			t.Setenv("OTEL_METRIC_EXPORT_INTERVAL", "3600000")
			prevTP, prevMP, prevLP, prevProp := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetLoggerProvider(), otel.GetTextMapPropagator()
			prevLogger := slog.Default()
			t.Cleanup(func() {
				otel.SetTracerProvider(prevTP)
				otel.SetMeterProvider(prevMP)
				otel.SetLoggerProvider(prevLP)
				otel.SetTextMapPropagator(prevProp)
				slog.SetDefault(prevLogger)
				tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = false, false, false, false
			})

			shutdown, err := initTibiaDataTelemetry(context.Background())

			require.NoError(t, err)
			assert.True(t, tibiaDataTelemetryEnabled)
			assert.Equal(t, test.wantTraces, tibiaDataTracingEnabled)
			assert.Equal(t, test.wantMetrics, tibiaDataMetricsEnabled)
			assert.Equal(t, test.wantLogs, tibiaDataLogsEnabled)
			_, tracesSDK := otel.GetTracerProvider().(*sdktrace.TracerProvider)
			assert.Equal(t, test.wantTraces, tracesSDK)
			_, metricsSDK := otel.GetMeterProvider().(*sdkmetric.MeterProvider)
			assert.Equal(t, test.wantMetrics, metricsSDK)
			_, logsSDK := otel.GetLoggerProvider().(*sdklog.LoggerProvider)
			assert.Equal(t, test.wantLogs, logsSDK)
			assert.NoError(t, shutdown(context.Background()))
		})
	}
}

func TestInitTibiaDataTelemetryEnablesLogs(t *testing.T) {
	t.Setenv(telemetryEnabledEnv, "true")
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	t.Setenv("OTEL_METRICS_EXPORTER", "none")
	t.Setenv("OTEL_LOGS_EXPORTER", "console")
	prevTP, prevMP, prevLP, prevProp := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetLoggerProvider(), otel.GetTextMapPropagator()
	prevLogger := slog.Default()
	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetMeterProvider(prevMP)
		otel.SetLoggerProvider(prevLP)
		otel.SetTextMapPropagator(prevProp)
		slog.SetDefault(prevLogger)
		tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = false, false, false, false
	})

	shutdown, err := initTibiaDataTelemetry(context.Background())
	require.NoError(t, err)
	assert.True(t, tibiaDataLogsEnabled)
	_, logsSDK := otel.GetLoggerProvider().(*sdklog.LoggerProvider)
	assert.True(t, logsSDK)
	require.NoError(t, shutdown(context.Background()))
}

func TestInitTibiaDataTelemetryRestoresProvidersWhenLogExporterIsInvalid(t *testing.T) {
	t.Setenv(telemetryEnabledEnv, "true")
	t.Setenv("OTEL_TRACES_EXPORTER", "console")
	t.Setenv("OTEL_METRICS_EXPORTER", "none")
	t.Setenv("OTEL_LOGS_EXPORTER", "invalid")
	prevTP, prevMP, prevLP := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetLoggerProvider()
	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetMeterProvider(prevMP)
		otel.SetLoggerProvider(prevLP)
		tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = false, false, false, false
	})

	_, err := initTibiaDataTelemetry(context.Background())
	require.Error(t, err)
	assert.Equal(t, prevTP, otel.GetTracerProvider())
	assert.Equal(t, prevMP, otel.GetMeterProvider())
	assert.Equal(t, prevLP, otel.GetLoggerProvider())
	assert.False(t, tibiaDataTelemetryEnabled)
}

func TestSlogBridgeExportsContextAndHonorsLevel(t *testing.T) {
	var stdout bytes.Buffer
	previousLogger := slog.Default()
	previousStdoutHandler := tibiaDataLogStdoutHandler
	previousLevel := tibiaDataLogConfiguredLevel
	t.Cleanup(func() {
		slog.SetDefault(previousLogger)
		tibiaDataLogStdoutHandler = previousStdoutHandler
		tibiaDataLogConfiguredLevel = previousLevel
	})

	tibiaDataLogConfiguredLevel = slog.LevelInfo
	tibiaDataLogStdoutHandler = newTibiaDataLogHandler(&stdout, "text", slog.LevelInfo)
	exporter := &recordingLogExporter{}
	provider := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exporter)))
	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })
	configureTibiaDataOTelLogging(provider)

	traceID, _ := trace.TraceIDFromHex("0123456789abcdef0123456789abcdef")
	spanID, _ := trace.SpanIDFromHex("0123456789abcdef")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	}))
	slog.InfoContext(ctx, "bridged message", "key", "value")
	slog.DebugContext(ctx, "filtered message")

	records := exporter.exportedRecords()
	require.Len(t, records, 1)
	assert.Equal(t, "bridged message", records[0].Body().AsString())
	assert.Equal(t, otelLog.SeverityInfo, records[0].Severity())
	assert.Equal(t, traceID, records[0].TraceID())
	assert.Equal(t, spanID, records[0].SpanID())
	var foundValue string
	records[0].WalkAttributes(func(kv attribute.KeyValue) bool {
		if kv.Key == "key" {
			foundValue = kv.Value.AsString()
		}
		return true
	})
	assert.Equal(t, "value", foundValue)
	assert.Contains(t, stdout.String(), "trace_id="+traceID.String())
	assert.Contains(t, stdout.String(), "span_id="+spanID.String())
}

func TestTibiaDataTelemetryResourceDefaults(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "")
	res, err := tibiaDataTelemetryResource(context.Background())
	require.NoError(t, err)

	name, ok := res.Set().Value("service.name")
	assert.True(t, ok)
	assert.Equal(t, telemetryService, name.AsString())

	t.Setenv("OTEL_SERVICE_NAME", "custom")
	res, err = tibiaDataTelemetryResource(context.Background())
	require.NoError(t, err)
	name, _ = res.Set().Value("service.name")
	assert.Equal(t, "custom", name.AsString())
}

func TestGinTelemetryMiddlewareTracesAndCorrelatesLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder, reader := setupTestTelemetry(t)
	logs := captureLogs(t, slog.LevelInfo)

	router := gin.New()
	router.Use(ginTelemetryMiddleware())
	router.Use(ginAccessLogMiddleware())
	router.GET("/v4/worlds/:world", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/v4/worlds/Antica", nil)
	req.Header.Set("traceparent", "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
	router.ServeHTTP(httptest.NewRecorder(), req)

	spans := recorder.Ended()
	require.Len(t, spans, 1)
	span := spans[0]
	assert.Equal(t, trace.SpanKindServer, span.SpanKind())
	assert.Equal(t, "0af7651916cd43dd8448eb211c80319c", span.SpanContext().TraceID().String())
	assert.Equal(t, "b7ad6b7169203331", span.Parent().SpanID().String())

	assert.Contains(t, logs.String(), "trace_id="+span.SpanContext().TraceID().String())
	assert.Contains(t, logs.String(), "span_id="+span.SpanContext().SpanID().String())

	assert.Equal(t, uint64(1), requestDurationCount(t, reader, "/v4/worlds/:world"))
}

func TestGinTelemetryMiddlewareSkipsProbeTracesButKeepsMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder, reader := setupTestTelemetry(t)

	router := gin.New()
	router.Use(ginTelemetryMiddleware())
	var probeCtx context.Context
	router.GET("/healthz", func(c *gin.Context) {
		probeCtx = c.Request.Context()
		c.Status(http.StatusOK)
	})

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/healthz", nil))

	assert.Empty(t, recorder.Ended())
	assert.False(t, trace.SpanFromContext(probeCtx).IsRecording())
	assert.Equal(t, uint64(1), requestDurationCount(t, reader, "/healthz"))
}

func TestTibiaDataHTMLDataCollectorPropagatesTraceContext(t *testing.T) {
	recorder, _ := setupTestTelemetry(t)

	prevClient := tibiaDataClient
	initTibiaDataClient()
	t.Cleanup(func() { tibiaDataClient = prevClient })

	var traceparent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceparent = r.Header.Get("traceparent")
		_, _ = w.Write([]byte("<html>ok</html>"))
	}))
	defer server.Close()

	ctx, parent := otel.Tracer("test").Start(context.Background(), "parent")
	_, err := TibiaDataHTMLDataCollector(ctx, TibiaDataRequestStruct{
		Method:  http.MethodGet,
		URL:     server.URL,
		RawBody: true,
	})
	parent.End()
	require.NoError(t, err)

	var client sdktrace.ReadOnlySpan
	for _, s := range recorder.Ended() {
		if s.SpanKind() == trace.SpanKindClient {
			client = s
		}
	}
	require.NotNil(t, client)
	assert.Equal(t, parent.SpanContext().SpanID(), client.Parent().SpanID())
	assert.Contains(t, traceparent, parent.SpanContext().TraceID().String())
	assert.Contains(t, traceparent, client.SpanContext().SpanID().String())
}

func TestTibiaDataTelemetryActive(t *testing.T) {
	prevEnabled, prevTracing, prevMetrics, prevLogs := tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled
	t.Cleanup(func() {
		tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = prevEnabled, prevTracing, prevMetrics, prevLogs
	})

	for _, tc := range []struct {
		enabled, tracing, metrics, logs, want bool
	}{
		{false, true, true, true, false},
		{true, false, false, false, false},
		{true, false, false, true, false},
		{true, false, true, false, true},
		{true, true, false, false, true},
	} {
		tibiaDataTelemetryEnabled, tibiaDataTracingEnabled, tibiaDataMetricsEnabled, tibiaDataLogsEnabled = tc.enabled, tc.tracing, tc.metrics, tc.logs
		assert.Equal(t, tc.want, tibiaDataTelemetryActive(), "%+v", tc)
	}
}
