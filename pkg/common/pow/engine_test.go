package pow

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()

	require.NotNil(t, engine)
	assert.NotNil(t, engine.logger)
}

func TestGenerateData(t *testing.T) {
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
	engine := NewEngine()
	complexity := uint8(10)

	data1 := engine.GenerateData(complexity)
	data2 := engine.GenerateData(complexity)

	assert.NotEqual(t, data1, data2, "generated data should be random and different")
	assert.Equal(t, complexity, data1[0])
	assert.Equal(t, complexity, data2[0])
}

func TestSerializeDeserializeNonce(t *testing.T) {
	tests := []struct {
		name   string
		salt   []byte
		proofs []uint32
	}{
		{
			name:   "basic case",
			salt:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
			proofs: []uint32{100, 200, 300},
		},
		{
			name:   "empty proofs",
			salt:   []byte{0, 0, 0, 0, 0, 0, 0, 0},
			proofs: []uint32{},
		},
		{
			name:   "large values",
			salt:   []byte{255, 255, 255, 255, 255, 255, 255, 255},
			proofs: []uint32{4294967295, 4294967294, 4294967293},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce := serializeNonce(tt.salt, tt.proofs)
			salt, proofs := deserializeNonce(nonce)

			assert.Equal(t, tt.salt, salt)
			require.Len(t, proofs, len(tt.proofs))
			assert.Equal(t, tt.proofs, proofs)
		})
	}
}

func TestSolveAndVerify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping solve test in short mode")
	}

	engine := NewEngine()
	complexity := uint8(0) // Use low complexity for faster tests

	data := engine.GenerateData(complexity)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nonce, err := engine.Solve(ctx, data)
	require.NoError(t, err)
	assert.NotEmpty(t, nonce)
	assert.True(t, engine.Verify(data, nonce))
}
