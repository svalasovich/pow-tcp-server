package log

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

const RequestIDContextKey = middleware.RequestIDKey

type requestIDHandler struct {
	slog.Handler
}

func (r *requestIDHandler) Handle(ctx context.Context, record slog.Record) error {
	requestID := ctx.Value(RequestIDContextKey)
	if requestID != nil {
		record.Add("requestID", requestID)
	}

	return r.Handler.Handle(ctx, record)
}

func (r *requestIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &requestIDHandler{Handler: r.Handler.WithAttrs(attrs)}
}

func (r *requestIDHandler) WithGroup(name string) slog.Handler {
	return &requestIDHandler{Handler: r.Handler.WithGroup(name)}
}
