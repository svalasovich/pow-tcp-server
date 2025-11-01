package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

func TestNewProvider(t *testing.T) {
	// given
	health := &HealthProvider{}
	ready := &ReadyProvider{}
	metric := &MetricProvider{}

	// when
	result := NewProvider(health, ready, metric)

	// then
	assert.Equal(t, health, result.Health)
	assert.Equal(t, metric, result.Metric)
	assert.Equal(t, metric, otel.GetMeterProvider())
}
