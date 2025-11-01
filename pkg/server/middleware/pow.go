package middleware

import (
	"context"
	"fmt"

	"github.com/svalasovich/pow-tcp-server/internal/log"
	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

type (
	PoWService interface {
		GenerateData(load uint64) (uint8, []byte)
		Verify(data []byte, nonce []byte) bool
	}

	MetricsService interface {
		Inc()
		RPS() uint64
	}
)

func PoW(powService PoWService, metricsService MetricsService) Middleware {
	logger := log.NewComponentLogger("server.middleware.pow")
	return func(next tcp.ConnectionHandler) tcp.ConnectionHandler {
		handleFunc := func(ctx context.Context, conn tcp.Connection) error {
			metricsService.Inc()
			rps := metricsService.RPS()

			complexity, challenge := powService.GenerateData(rps)
			if err := conn.Write(challenge); err != nil {
				return fmt.Errorf("failed send challenge: %w", err)
			}
			logger.DebugContext(ctx, "sent challenge", "complexity", complexity)

			nonce, err := conn.Read()
			if err != nil {
				return fmt.Errorf("failed read nonce")
			}
			powService.Verify(challenge, nonce)

			return next.Handle(ctx, conn)
		}

		return &middleware{handleFunc}
	}
}
