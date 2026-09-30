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
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, buf.String(), "http request")
	assert.Contains(t, buf.String(), "/ping")
}

func TestTraceLogAttrsEmptyWithoutTracing(t *testing.T) {
	assert.Empty(t, traceLogAttrs(context.Background()))
}
