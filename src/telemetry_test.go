package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// setupTestTelemetry installs in-memory tracer and meter providers as the
// global providers and restores the previous state when the test ends.
func setupTestTelemetry(t *testing.T) (*tracetest.SpanRecorder, *sdkmetric.ManualReader) {
	t.Helper()

	recorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	reader := sdkmetric.NewManualReader()
	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	prevTP, prevMP, prevProp := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetTextMapPropagator()
	prevEnabled := tibiaDataTelemetryEnabled

	otel.SetTracerProvider(tracerProvider)
	otel.SetMeterProvider(meterProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tibiaDataTelemetryEnabled = true

	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetMeterProvider(prevMP)
		otel.SetTextMapPropagator(prevProp)
		tibiaDataTelemetryEnabled = prevEnabled
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
	t.Setenv(telemetryEnabledEnv, "true")
	t.Setenv("OTEL_TRACES_EXPORTER", "none")
	t.Setenv("OTEL_METRICS_EXPORTER", "none")
	prevTP, prevMP, prevProp := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetTextMapPropagator()
	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetMeterProvider(prevMP)
		otel.SetTextMapPropagator(prevProp)
		tibiaDataTelemetryEnabled = false
	})

	shutdown, err := initTibiaDataTelemetry(context.Background())

	require.NoError(t, err)
	assert.True(t, tibiaDataTelemetryEnabled)
	_, isSDK := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	assert.True(t, isSDK)
	assert.NoError(t, shutdown(context.Background()))
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
	_, err := TibiaDataHTMLDataCollector(TibiaDataRequestStruct{
		Method:  http.MethodGet,
		URL:     server.URL,
		RawBody: true,
		Ctx:     ctx,
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
