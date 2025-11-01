package monitoring

import (
	"fmt"

	"github.com/hellofresh/health-go/v5"
)

type HealthProvider struct {
	*health.Health
}

func NewHealthProvider(name Name, version Version) (*HealthProvider, error) {
	appHealth, err := health.New(
		health.WithComponent(health.Component{
			Name:    string(name),
			Version: string(version),
		}),
		health.WithSystemInfo(),
		health.WithMaxConcurrent(1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create health check: %w", err)
	}

	return &HealthProvider{Health: appHealth}, nil
}
