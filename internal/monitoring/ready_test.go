package monitoring

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/hellofresh/health-go/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReadyProvider(t *testing.T) {
	// given
	name := gofakeit.Word()
	appVersion := gofakeit.Word()

	// when
	provider, err := NewReadyProvider(Name(name), Version(appVersion))

	// then
	require.NoError(t, err)

	check := provider.Health.Measure(context.Background())
	assert.Equal(t, name, check.Component.Name)
	assert.Equal(t, appVersion, check.Component.Version)
	assert.Equal(t, health.StatusOK, check.Status)
}
