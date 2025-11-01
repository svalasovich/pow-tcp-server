package log

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/svalasovich/golang-template/internal/config"
)

type Logger struct {
	*slog.Logger
}

func Init(cfg config.Log) {
	// TODO integrate opentelemetry logging and use it
	options := &slog.HandlerOptions{AddSource: cfg.ShowSource, Level: parseLogLevel(cfg.Level)}

	var handler slog.Handler
	if cfg.JSONFormat {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	slog.SetDefault(slog.New(&requestIDHandler{Handler: handler}))
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
