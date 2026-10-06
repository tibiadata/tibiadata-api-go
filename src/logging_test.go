package main

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func TestGinAccessLogMiddlewareIncludesConfiguredRequestIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureRequestLogHeadersForTest(t, "X-Request-ID", "X-Trace-ID")

	router := gin.New()
	router.Use(ginAccessLogMiddleware())
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set(requestIDHeader, "abc123-FRA")
	req.Header.Set(correlationIDHeader, "kong-request-456")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.Contains(t, buf.String(), "request_id="+obfuscatedRequestLogValue("abc123-FRA"))
	assert.Contains(t, buf.String(), "correlation_id="+obfuscatedRequestLogValue("kong-request-456"))
	assert.NotContains(t, buf.String(), "abc123-FRA")
	assert.NotContains(t, buf.String(), "kong-request-456")
}

func TestGinAccessLogMiddlewareOmitsUnconfiguredRequestIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureRequestLogHeadersForTest(t, "", "")

	router := gin.New()
	router.Use(ginAccessLogMiddleware())
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("CF-Ray", "abc123-FRA")
	req.Header.Set("X-Correlation-ID", "kong-request-456")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.NotContains(t, buf.String(), "request_id=")
	assert.NotContains(t, buf.String(), "correlation_id=")
}

func TestGinAccessLogMiddlewareOmitsSensitiveConfiguredHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureRequestLogHeadersForTest(t, "Authorization", "Cookie")

	router := gin.New()
	router.Use(ginAccessLogMiddleware())
	router.GET("/items", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("Authorization", "test-auth-header-value")
	req.Header.Set("Cookie", "session=secret-cookie")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.NotContains(t, buf.String(), "test-auth-header-value")
	assert.NotContains(t, buf.String(), "secret-cookie")
	assert.NotContains(t, buf.String(), "request_id=")
	assert.NotContains(t, buf.String(), "correlation_id=")
}

func TestGinRecoveryMiddlewareIncludesConfiguredRequestIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureRequestLogHeadersForTest(t, "X-Request-ID", "X-Trace-ID")

	router := gin.New()
	router.Use(ginRecoveryMiddleware())
	router.GET("/boom", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	req.Header.Set(requestIDHeader, "abc123-FRA")
	req.Header.Set(correlationIDHeader, "kong-request-456")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.Contains(t, buf.String(), "request_id="+obfuscatedRequestLogValue("abc123-FRA"))
	assert.Contains(t, buf.String(), "correlation_id="+obfuscatedRequestLogValue("kong-request-456"))
	assert.NotContains(t, buf.String(), "abc123-FRA")
	assert.NotContains(t, buf.String(), "kong-request-456")
}

func TestGinRecoveryMiddlewareOmitsSensitiveConfiguredHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	buf := captureLogs(t, slog.LevelInfo)
	configureRequestLogHeadersForTest(t, "Authorization", "Cookie")

	router := gin.New()
	router.Use(ginRecoveryMiddleware())
	router.GET("/boom", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	req.Header.Set("Authorization", "test-auth-header-value")
	req.Header.Set("Cookie", "session=secret-cookie")
	router.ServeHTTP(httptest.NewRecorder(), req)

	assert.NotContains(t, buf.String(), "test-auth-header-value")
	assert.NotContains(t, buf.String(), "secret-cookie")
	assert.NotContains(t, buf.String(), "request_id=")
	assert.NotContains(t, buf.String(), "correlation_id=")
}

func TestTraceLogAttrsEmptyWithoutTracing(t *testing.T) {
	assert.Empty(t, traceLogAttrs(context.Background()))
}

func configureRequestLogHeadersForTest(t *testing.T, requestHeader, correlationHeader string) {
	t.Helper()
	t.Cleanup(configureRequestLogHeaders)
	t.Setenv(requestIDHeaderEnv, requestHeader)
	t.Setenv(correlationIDHeaderEnv, correlationHeader)
	configureRequestLogHeaders()
}

func captureLogs(t *testing.T, level slog.Level) *bytes.Buffer {
	var buf bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(newTibiaDataLogHandler(&buf, "text", level)))
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
