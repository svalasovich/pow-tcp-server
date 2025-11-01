package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/svalasovich/golang-template/internal/config"
	"github.com/svalasovich/golang-template/internal/log"
)

type Server struct {
	cfg    config.GRPCServer
	logger *log.Logger
	*grpc.Server
}

func NewServer(cfg config.GRPCServer, opt ...grpc.ServerOption) *Server {
	// TODO Add tests
	return &Server{
		cfg:    cfg,
		logger: log.NewComponentLogger("grpc-server"),
		Server: grpc.NewServer(opt...),
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen network: %w", err)
	}

	go func() {
		if err := s.Server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("failed to serve GRPC server: %w", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	signal := make(chan struct{})

	go func() {
		defer close(signal)
		s.GracefulStop()
	}()

	select {
	case <-ctx.Done():
		s.Stop()
		return fmt.Errorf("failed to graceful stop GRPC server: %w", ctx.Err())
	case <-signal:
		return nil
	}
}
