package main

import (
	"context"
	"fmt"

	"github.com/svalasovich/pow-tcp-server/internal/metrics"
	"github.com/svalasovich/pow-tcp-server/pkg/common/pow"
	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
	"github.com/svalasovich/pow-tcp-server/pkg/repository"
	"github.com/svalasovich/pow-tcp-server/pkg/server"
	"github.com/svalasovich/pow-tcp-server/pkg/server/middleware"
	"github.com/svalasovich/pow-tcp-server/pkg/service"
)

func bootstrap(ctx context.Context, cfg *Config) error {
	powService := pow.NewService(pow.NewEngine(), pow.NewAdaptiveComplexityCalculator(cfg.PoW))

	metricsService := metrics.NewService()
	metricsService.Start(ctx)

	quoteService := service.NewQuote(repository.NewQuote())
	quoteHanlder := server.NewQuoteHandler(quoteService)
	powMiddleware := middleware.PoW(powService, metricsService)

	tcpServer := tcp.NewServer(cfg.TCP)
	tcpServer.AddHandler(middleware.Chain(quoteHanlder, powMiddleware))
	err := tcpServer.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	<-ctx.Done()

	return nil
}
