package log

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type Logger struct {
	*slog.Logger
}

func Init(cfg Config) {
	options := &slog.HandlerOptions{AddSource: cfg.ShowSource, Level: parseLogLevel(cfg.Level)}

	var handler slog.Handler
	if cfg.JSONFormat {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	slog.SetDefault(slog.New(handler))
}

func NewComponentLogger(name string) *Logger {
	return &Logger{Logger: slog.With("component", name)}
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func parseLogLevel(level string) slog.Level {
	level = strings.ToUpper(level)
	switch {
	case level == "DEBUG":
		return slog.LevelDebug
	case level == "INFO":
		return slog.LevelInfo
	case level == "WARN":
		return slog.LevelWarn
	default:
		return slog.LevelError
	}
}
