package monitoring

import (
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"

	"github.com/svalasovich/golang-template/internal/config"
	"github.com/svalasovich/golang-template/internal/http"
	"github.com/svalasovich/golang-template/internal/log"
)

type (
	Name string

	Version string

	Provider struct {
		Health *HealthProvider
		Ready  *ReadyProvider
		Metric *MetricProvider
	}
)

func NewProvider(health *HealthProvider, ready *ReadyProvider, metric *MetricProvider) *Provider {
	otel.SetLogger(logr.FromSlogHandler(log.NewComponentLogger("monitoring").Handler()))
	otel.SetMeterProvider(metric)

	return &Provider{
		Health: health,
		Ready:  ready,
		Metric: metric,
	}
}

func NewServer(cfg config.Monitoring, provider *Provider) *http.Server {
	server := http.NewServer(cfg.Server)

	server.Router().Get("/metrics", provider.Metric.Handler().ServeHTTP)
	server.Router().Get("/health", provider.Health.Handler().ServeHTTP)
	server.Router().Get("/ready", provider.Ready.Handler().ServeHTTP)

	return server
}
