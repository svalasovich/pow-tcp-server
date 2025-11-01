package tcp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type (
	ConnectionHandler interface {
		Handle(ctx context.Context, conn Connection) error
	}

	Connection struct {
		net.Conn
	}
)

func (c Connection) Write(data []byte) error {
	err := binary.Write(c.Conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return fmt.Errorf("failed send size message: %w", err)
	}

	_, err = c.Conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed send message: %w", err)
	}

	return nil
}

func (c Connection) Read() ([]byte, error) {
	var length uint64
	if err := binary.Read(c.Conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed receive size message: %w", err)
	}

	result := make([]byte, length)
	if _, err := io.ReadFull(c.Conn, result); err != nil {
		return nil, fmt.Errorf("failed receive message: %w", err)
	}

	return result, nil
}
