package http

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/svalasovich/golang-template/internal/config"
	"github.com/svalasovich/golang-template/internal/log"
)

type Server struct {
	logger *log.Logger
	router *chi.Mux
	*http.Server
}

func NewServer(cfg config.HTTPServer) *Server {
	router := newRouter()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	return &Server{
		router: router,
		Server: server,
		logger: log.NewComponentLogger("http-server"),
	}
}

func newRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Recoverer, middleware.CleanPath)

	return router
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to listen network: %w", err)
	}

	go func() {
		if err := s.Server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("failed to serve HTTP server: %w", err)
		}
	}()

	return nil
}
