package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	otelLog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

const (
	cloudflareRayHeader     = "Cf-Ray"
	cloudflareRayLoggingEnv = "TIBIADATA_LOG_CF_RAY"
)

var cloudflareRayLoggingEnabled bool
var tibiaDataLogStdoutHandler slog.Handler
var tibiaDataLogConfiguredLevel slog.Leveler

func initTibiaDataLogging() {
	level := tibiaDataLogLevel()
	format := strings.ToLower(getEnv("TIBIADATA_LOG_FORMAT", "text"))
	tibiaDataLogConfiguredLevel = level
	tibiaDataLogStdoutHandler = newTibiaDataLogHandler(os.Stdout, format, level)
	slog.SetDefault(slog.New(traceContextHandler{handler: tibiaDataLogStdoutHandler}))
	configureCloudflareRayLogging()

	if v := strings.TrimSpace(os.Getenv("TIBIADATA_LOG_LEVEL")); v != "" && !isKnownLogLevel(v) {
		slog.Warn("unknown TIBIADATA_LOG_LEVEL, falling back to info", "value", v)
	}
	if format != "text" && format != "json" {
		slog.Warn("unknown TIBIADATA_LOG_FORMAT, falling back to text", "value", format)
	}
}

func configureTibiaDataOTelLogging(provider otelLog.LoggerProvider) {
	stdoutHandler := traceContextHandler{handler: tibiaDataLogStdoutHandler}
	otelHandler := levelHandler{
		handler: otelslog.NewHandler(telemetryService, otelslog.WithLoggerProvider(provider)),
		level:   tibiaDataLogConfiguredLevel,
	}
	slog.SetDefault(slog.New(slog.NewMultiHandler(stdoutHandler, otelHandler)))
}

func logOpenTelemetryError(err error) {
	slog.New(traceContextHandler{handler: tibiaDataLogStdoutHandler}).Warn(
		"TibiaData API opentelemetry error",
		"error", err,
	)
}

func configureCloudflareRayLogging() {
	cloudflareRayLoggingEnabled = getEnvAsBool(cloudflareRayLoggingEnv, false)
}

func tibiaDataLogLevel() slog.Leveler {
	if levelStr := strings.TrimSpace(os.Getenv("TIBIADATA_LOG_LEVEL")); levelStr != "" {
		return parseTibiaDataLogLevel(levelStr)
	}
	if getEnvAsBool("DEBUG_MODE", false) {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func isKnownLogLevel(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug", "info", "warn", "warning", "error":
		return true
	}
	return false
}

func parseTibiaDataLogLevel(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func newTibiaDataLogHandler(w io.Writer, format string, level slog.Leveler) slog.Handler {
	opts := &slog.HandlerOptions{Level: level}
	if strings.ToLower(format) == "json" {
		return slog.NewJSONHandler(w, opts)
	}
	return slog.NewTextHandler(w, opts)
}

type levelHandler struct {
	handler slog.Handler
	level   slog.Leveler
}

func (h levelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level() && h.handler.Enabled(ctx, level)
}

func (h levelHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler.Handle(ctx, record)
}

func (h levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return levelHandler{handler: h.handler.WithAttrs(attrs), level: h.level}
}

func (h levelHandler) WithGroup(name string) slog.Handler {
	return levelHandler{handler: h.handler.WithGroup(name), level: h.level}
}

type traceContextHandler struct {
	handler slog.Handler
}

func (h traceContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h traceContextHandler) Handle(ctx context.Context, record slog.Record) error {
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		record = record.Clone()
		record.AddAttrs(
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
		)
	}
	return h.handler.Handle(ctx, record)
}

func (h traceContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return traceContextHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h traceContextHandler) WithGroup(name string) slog.Handler {
	return traceContextHandler{handler: h.handler.WithGroup(name)}
}

func requestLogAttrs(r *http.Request) []any {
	if !cloudflareRayLoggingEnabled {
		return nil
	}
	if requestID := r.Header.Get(cloudflareRayHeader); requestID != "" {
		return []any{"request_id", requestID}
	}
	return nil
}

func ginRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

func ginAccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
			"client_ip", c.ClientIP(),
		}
		attrs = append(attrs, requestLogAttrs(c.Request)...)

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		status := c.Writer.Status()
		level := slog.LevelInfo
		switch {
		case isProbePath(c.Request.URL.Path):
			level = slog.LevelDebug
		case status >= http.StatusInternalServerError:
			level = slog.LevelError
		case status >= http.StatusBadRequest:
			level = slog.LevelWarn
		}

		slog.Log(c.Request.Context(), level, "http request", attrs...)
	}
}

func isProbePath(path string) bool {
	switch path {
	case "/ping", "/health", "/healthz", "/readyz":
		return true
	}
	return false
}

func ginRecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		attrs := []any{
			"panic", fmt.Sprint(recovered),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"stack", string(debug.Stack()),
		}
		attrs = append(attrs, requestLogAttrs(c.Request)...)
		slog.ErrorContext(c.Request.Context(), "panic recovered", attrs...)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

func tibiaDataLogFatal(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
	} else {
		slog.Error(msg)
	}
	shutdownTibiaDataTelemetry()
	os.Exit(1)
}
