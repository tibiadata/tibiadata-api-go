package main

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestParseTibiaDataLogLevel(t *testing.T) {
	assert.Equal(t, slog.LevelDebug, parseTibiaDataLogLevel("debug"))
	assert.Equal(t, slog.LevelWarn, parseTibiaDataLogLevel("warn"))
	assert.Equal(t, slog.LevelError, parseTibiaDataLogLevel("error"))
	assert.Equal(t, slog.LevelInfo, parseTibiaDataLogLevel("info"))
	assert.Equal(t, slog.LevelInfo, parseTibiaDataLogLevel("unknown"))
}

func TestNewTibiaDataLogHandlerJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(newTibiaDataLogHandler(&buf, "json", slog.LevelInfo))
	logger.Info("hello", "key", "value")
	assert.Contains(t, buf.String(), `"msg":"hello"`)
	assert.Contains(t, buf.String(), `"key":"value"`)
}

func TestGinAccessLogMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(newTibiaDataLogHandler(&buf, "text", slog.LevelInfo)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	router := gin.New()
	router.Use(ginAccessLogMiddleware())
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, buf.String(), "http request")
	assert.Contains(t, buf.String(), "/items")
}

func TestGinAccessLogMiddlewareIncludesCloudflareRayID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureCloudflareRayLoggingForTest(t, true)

	router := gin.New()
	router.Use(ginAccessLogMiddleware())
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set(cloudflareRayHeader, "abc123-FRA")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.Contains(t, buf.String(), "request_id=abc123-FRA")
}

func TestGinAccessLogMiddlewareOmitsMissingCloudflareRayID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureCloudflareRayLoggingForTest(t, true)

	router := gin.New()
	router.Use(ginAccessLogMiddleware())
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.NotContains(t, buf.String(), "request_id=")
}

func TestGinAccessLogMiddlewareOmitsCloudflareRayIDWhenDisabled(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureCloudflareRayLoggingForTest(t, false)

	router := gin.New()
	router.Use(ginAccessLogMiddleware())
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set(cloudflareRayHeader, "abc123-FRA")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.NotContains(t, buf.String(), "request_id=")
}

func TestGinRecoveryMiddlewareIncludesCloudflareRayID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureCloudflareRayLoggingForTest(t, true)

	router := gin.New()
	router.Use(ginRecoveryMiddleware())
	router.GET("/boom", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	req.Header.Set(cloudflareRayHeader, "abc123-FRA")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.Contains(t, buf.String(), "request_id=abc123-FRA")
}

func TestGinRecoveryMiddlewareOmitsMissingCloudflareRayID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureCloudflareRayLoggingForTest(t, true)

	router := gin.New()
	router.Use(ginRecoveryMiddleware())
	router.GET("/boom", func(c *gin.Context) {
		panic("boom")
	})

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/boom", nil))

	assert.NotContains(t, buf.String(), "request_id=")
}

func configureCloudflareRayLoggingForTest(t *testing.T, enabled bool) {
	t.Helper()
	t.Cleanup(configureCloudflareRayLogging)
	t.Setenv(cloudflareRayLoggingEnv, strconv.FormatBool(enabled))
	configureCloudflareRayLogging()
}

func TestTraceContextHandlerAddsTraceAttributes(t *testing.T) {
	var buf bytes.Buffer
	traceID, _ := trace.TraceIDFromHex("0123456789abcdef0123456789abcdef")
	spanID, _ := trace.SpanIDFromHex("0123456789abcdef")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
	}))

	handler := traceContextHandler{handler: newTibiaDataLogHandler(&buf, "text", slog.LevelInfo)}
	require.NoError(t, handler.Handle(ctx, slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0)))

	assert.Contains(t, buf.String(), "trace_id=0123456789abcdef0123456789abcdef")
	assert.Contains(t, buf.String(), "span_id=0123456789abcdef")
}

func captureLogs(t *testing.T, level slog.Level) *bytes.Buffer {
	var buf bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(traceContextHandler{handler: newTibiaDataLogHandler(&buf, "text", level)}))
	t.Cleanup(func() { slog.SetDefault(previous) })
	return &buf
}

func TestGinAccessLogLevels(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)

	router := gin.New()
	router.Use(ginAccessLogMiddleware(), ginRecoveryMiddleware())
	router.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/boom", func(c *gin.Context) { panic("boom") })

	for _, path := range []string{"/ping", "/missing", "/boom"} {
		router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, path, nil))
	}

	out := buf.String()
	assert.NotContains(t, out, "path=/ping")
	assert.Contains(t, out, "level=WARN msg=\"http request\"")
	assert.Contains(t, out, "status=404")
	assert.Contains(t, out, "level=ERROR msg=\"panic recovered\"")
	assert.Contains(t, out, "level=ERROR msg=\"http request\"")
	assert.Contains(t, out, "status=500")
}
