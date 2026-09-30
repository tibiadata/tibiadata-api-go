package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func initTibiaDataLogging() {
	level := tibiaDataLogLevel()
	format := strings.ToLower(getEnv("TIBIADATA_LOG_FORMAT", "text"))
	handler := newTibiaDataLogHandler(os.Stdout, format, level)
	slog.SetDefault(slog.New(handler))
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

func traceLogAttrs(ctx context.Context) []any {
	_ = ctx
	return nil // extended when OpenTelemetry tracing (#749) is wired
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
		attrs = append(attrs, traceLogAttrs(c.Request.Context())...)

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		slog.Info("http request", attrs...)
	}
}

func tibiaDataLogFatal(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
	} else {
		slog.Error(msg)
	}
	os.Exit(1)
}
