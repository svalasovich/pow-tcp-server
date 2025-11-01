package middleware

import (
	"context"

	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

type (
	Middleware func(tcp.ConnectionHandler) tcp.ConnectionHandler

	middleware struct {
		handleFunc func(ctx context.Context, conn tcp.Connection) error
	}
)

func Chain(handler tcp.ConnectionHandler, middleware ...Middleware) tcp.ConnectionHandler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

func (m *middleware) Handle(ctx context.Context, conn tcp.Connection) error {
	return m.handleFunc(ctx, conn)
}
