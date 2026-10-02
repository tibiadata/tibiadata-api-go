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
)

func initTibiaDataLogging() {
	level := tibiaDataLogLevel()
	format := strings.ToLower(getEnv("TIBIADATA_LOG_FORMAT", "text"))
	handler := newTibiaDataLogHandler(os.Stdout, format, level)
	slog.SetDefault(slog.New(handler))

	if v := strings.TrimSpace(os.Getenv("TIBIADATA_LOG_LEVEL")); v != "" && !isKnownLogLevel(v) {
		slog.Warn("unknown TIBIADATA_LOG_LEVEL, falling back to info", "value", v)
	}
	if format != "text" && format != "json" {
		slog.Warn("unknown TIBIADATA_LOG_FORMAT, falling back to text", "value", format)
	}
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
		slog.Error("panic recovered",
			"panic", fmt.Sprint(recovered),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"stack", string(debug.Stack()),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

func tibiaDataLogFatal(msg string, err error) {
	if err != nil {
		slog.Error(msg, "error", err)
	} else {
		slog.Error(msg)
	}
	os.Exit(1)
}
