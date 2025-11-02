package pow

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	t.Parallel()

	engine := NewEngine()

	require.NotNil(t, engine)
	assert.NotNil(t, engine.logger)
}

func TestGenerateData(t *testing.T) {
	t.Parallel()

	engine := NewEngine()

	tests := []struct {
		name       string
		complexity uint8
	}{
		{"complexity 0", 0},
		{"complexity 10", 10},
		{"complexity 255", 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := engine.GenerateData(tt.complexity)

			assert.Len(t, data, uint8Size+randomDataSize)
			assert.Equal(t, tt.complexity, data[0])
		})
	}
}

func TestGenerateDataRandomness(t *testing.T) {
	t.Parallel()

	engine := NewEngine()
	complexity := uint8(10)

	data1 := engine.GenerateData(complexity)
	data2 := engine.GenerateData(complexity)

	assert.NotEqual(t, data1, data2, "generated data should be random and different")
	assert.Equal(t, complexity, data1[0])
	assert.Equal(t, complexity, data2[0])
}

func TestVerifyComplexity(t *testing.T) {
	t.Parallel()

	engine := NewEngine()

	tests := []struct {
		name       string
		complexity uint8
		nonce      []byte
		shouldPass bool
	}{
		{
			name:       "complexity 0 always passes",
			complexity: 0,
			nonce:      []byte{0, 0, 0, 0, 0, 0, 0, 1},
			shouldPass: true,
		},
		{
			name:       "invalid nonce fails verification",
			complexity: 10,
			nonce:      []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := engine.GenerateData(tt.complexity)
			result := engine.Verify(data, tt.nonce)

			if tt.shouldPass {
				assert.True(t, result)
			} else {
				assert.False(t, result)
			}
		})
	}
}

func TestSolveAndVerify(t *testing.T) {
	t.Parallel()

	engine := NewEngine()
	complexity := uint8(0) // Use low complexity for faster tests

	data := engine.GenerateData(complexity)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	nonce, err := engine.Solve(ctx, data)
	require.NoError(t, err)
	assert.NotEmpty(t, nonce)
	assert.True(t, engine.Verify(data, nonce))
}
