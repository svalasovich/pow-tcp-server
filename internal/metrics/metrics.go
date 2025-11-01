package metrics

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/svalasovich/pow-tcp-server/internal/log"
)

const rpsCalculateTick = 5 * time.Second

type Service struct {
	totalRequests atomic.Uint64
	rps           atomic.Uint64
	logger        *log.Logger
}

func NewService() *Service {
	return &Service{
		logger: log.NewComponentLogger("metrics.service"),
	}
}

func (s *Service) Inc() {
	s.totalRequests.Add(1)
}

func (s *Service) Start(ctx context.Context) {
	s.logger.Info("start metrics service", "tick", rpsCalculateTick)
	go func() {
		ticker := time.NewTicker(rpsCalculateTick)
		defer ticker.Stop()

		previousTotalRequests := uint64(0)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				currentTotalRequests := s.totalRequests.Load()
				rps := currentTotalRequests - previousTotalRequests
				previousTotalRequests = currentTotalRequests

				s.rps.Store(rps)

				s.logger.Debug("RPS refreshed", "rps", rps)
			}
		}
	}()
}

func (s *Service) RPS() uint64 {
	return s.rps.Load()
}
