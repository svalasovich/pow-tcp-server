package client

import (
	"context"
	"fmt"

	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

type (
	Client interface {
		Connect(ctx context.Context) (tcp.Connection, error)
	}

	Quote struct {
		client Client
	}
)

func NewQuoteClient(client Client) *Quote {
	return &Quote{client: client}
}

func (q *Quote) FetchRandomQuote(ctx context.Context) (string, error) {
	conn, err := q.client.Connect(ctx)
	if err != nil {
		return "", fmt.Errorf("failed create connection: %w", err)
	}
	defer conn.Close()

	rawQuote, err := conn.Read()
	if err != nil {
		return "", fmt.Errorf("failed read quote: %w", err)
	}

	return string(rawQuote), nil
}
