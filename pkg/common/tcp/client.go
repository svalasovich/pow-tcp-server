package tcp

import (
	"context"
	"fmt"
	"net"
	"time"
)

type (
	Client struct {
		cfg Config
	}
)

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Connect(ctx context.Context) (Connection, error) {
	dialer := &net.Dialer{
		KeepAlive: c.cfg.KeepAlive,
	}

	conn, err := dialer.DialContext(ctx, "tcp", c.cfg.Address)
	if err != nil {
		return Connection{}, fmt.Errorf("failed to connect to TCP server: %w", err)
	}

	if err := conn.SetDeadline(time.Now().Add(c.cfg.Deadline)); err != nil {
		return Connection{}, fmt.Errorf("failed set deadline: %w", err)
	}

	return Connection{conn}, nil
}
