// pkg/logger/logger.go
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type contextKey string

const loggerKey contextKey = "logger"

// Setup initialises the global slog logger.
// level: "debug", "info", "warn", "error"
// format: "json" or "text"
func Setup(level string, format string, output io.Writer) *slog.Logger {
	if output == nil {
		output = os.Stdout
	}

	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: level == "debug",
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// WithContext returns a new context carrying the logger.
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts the logger from context, falling back to the default logger.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
