package monitoring

import (
	"fmt"

	"github.com/hellofresh/health-go/v5"
)

type ReadyProvider struct {
	*health.Health
}

func NewReadyProvider(name Name, version Version) (*ReadyProvider, error) {
	appReady, err := health.New(
		health.WithComponent(health.Component{
			Name:    string(name),
			Version: string(version),
		}),
		health.WithMaxConcurrent(1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ready check: %w", err)
	}

	return &ReadyProvider{Health: appReady}, nil
}
