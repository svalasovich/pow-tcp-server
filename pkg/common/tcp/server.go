package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/svalasovich/pow-tcp-server/internal/log"
)

type (
	Server struct {
		cfg     Config
		handler ConnectionHandler
		logger  *log.Logger
	}
)

func NewServer(cfg Config) *Server {
	return &Server{
		cfg:    cfg,
		logger: log.NewComponentLogger("tcp.server"),
	}
}

func (s *Server) AddHandler(handler ConnectionHandler) {
	s.handler = handler
}

func (s *Server) Start(ctx context.Context) error {
	if s.handler == nil {
		return errors.New("no handler")
	}

	s.logger.Info("start TCP server", "address", s.cfg.Address)

	cfg := net.ListenConfig{KeepAlive: s.cfg.KeepAlive}
	listener, err := cfg.Listen(ctx, "tcp", s.cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}

	go func() {
		s.listenConnections(ctx, listener.(*net.TCPListener))
	}()

	s.logger.Info("TCP server started")

	return nil
}

func (s *Server) listenConnections(ctx context.Context, listener *net.TCPListener) {
	var wg sync.WaitGroup
	defer wg.Done()

	listenCtx, listenCancel := context.WithCancel(ctx)
	defer listenCancel()

	for {
		ctx, cancel := context.WithTimeout(listenCtx, s.cfg.Deadline)
		if err := listener.SetDeadline(time.Now().Add(s.cfg.Deadline)); err != nil {
			continue
		}

		conn, err := listener.Accept()
		if err != nil {
			var opErr net.Error
			if errors.As(err, &opErr) && opErr.Timeout() {
				continue
			}
		}

		wg.Go(func() {
			connection := Connection{conn}
			defer connection.Close()
			defer cancel()

			err := s.handler.Handle(ctx, connection)
			if err != nil {
				s.logger.ErrorContext(ctx, "failed handle request", "error", err)
			}
		})

		select {
		case <-ctx.Done():
			cancel()
			break
		default:
		}
	}
}
