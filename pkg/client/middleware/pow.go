package middleware

import (
	"context"
	"fmt"

	"github.com/svalasovich/pow-tcp-server/pkg/client"
	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

type PoWEngine interface {
	Solve(ctx context.Context, data []byte) ([]byte, error)
}

func PoW(engine PoWEngine) Middleware {
	return func(client client.Client) client.Client {
		connectFunc := func(ctx context.Context) (tcp.Connection, error) {
			conn, err := client.Connect(ctx)
			if err != nil {
				return tcp.Connection{}, err
			}

			data, err := conn.Read()
			if err != nil {
				return tcp.Connection{}, fmt.Errorf("failed to read challenge: %w", err)
			}

			nonce, err := engine.Solve(ctx, data)
			if err != nil {
				return tcp.Connection{}, fmt.Errorf("failed to solve challenge: %w", err)
			}

			if err := conn.Write(nonce); err != nil {
				return tcp.Connection{}, fmt.Errorf("failed to send nonce: %w", err)
			}

			return conn, nil
		}

		return &middleware{connectFunc}
	}
}
