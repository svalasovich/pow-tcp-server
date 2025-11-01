package main

import (
	"context"
	"time"

	"github.com/svalasovich/pow-tcp-server/internal/log"
	"github.com/svalasovich/pow-tcp-server/pkg/client"
	"github.com/svalasovich/pow-tcp-server/pkg/client/middleware"
	"github.com/svalasovich/pow-tcp-server/pkg/common/pow"
	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

func bootstrap(ctx context.Context, cfg *Config) error {
	powEngine := pow.NewEngine()

	powMiddleware := middleware.PoW(powEngine)

	tcpClient := tcp.NewClient(cfg.TCP)
	client := client.NewQuoteClient(middleware.Chain(tcpClient, powMiddleware))

	logger := log.NewComponentLogger("bootstrap")
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			now := time.Now()
			quote, err := client.FetchRandomQuote(ctx)
			if err != nil {
				logger.Error("failed fetch quote", "error", err)
				continue
			}

			logger.Info(quote, "duration", time.Since(now))
		}
	}
}
