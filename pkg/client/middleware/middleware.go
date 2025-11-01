package middleware

import (
	"context"

	"github.com/svalasovich/pow-tcp-server/pkg/client"
	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

type (
	Middleware func(next client.Client) client.Client

	middleware struct {
		connectFunc func(ctx context.Context) (tcp.Connection, error)
	}
)

func Chain(f client.Client, middlewares ...Middleware) client.Client {
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}

	return f
}

func (m *middleware) Connect(ctx context.Context) (tcp.Connection, error) {
	return m.connectFunc(ctx)
}
