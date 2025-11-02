package pow

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAdaptiveComplexityCalculator(t *testing.T) {
	t.Parallel()

	cfg := Config{ComplexityStep: 100}
	calc := NewAdaptiveComplexityCalculator(cfg)

	assert.NotNil(t, calc)
	assert.Equal(t, cfg, calc.cfg)
}

func TestAdaptiveComplexityCalculator_Calculate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		complexityStep uint
		load           uint64
		expectedResult uint8
	}{
		{
			name:           "zero load",
			complexityStep: 100,
			load:           0,
			expectedResult: 0,
		},
		{
			name:           "load below step",
			complexityStep: 100,
			load:           50,
			expectedResult: 0,
		},
		{
			name:           "load equal to step",
			complexityStep: 100,
			load:           100,
			expectedResult: 1,
		},
		{
			name:           "load multiple of step",
			complexityStep: 100,
			load:           500,
			expectedResult: 5,
		},
		{
			name:           "load not multiple of step",
			complexityStep: 100,
			load:           550,
			expectedResult: 5,
		},
		{
			name:           "max uint8 boundary",
			complexityStep: 1,
			load:           255,
			expectedResult: 255,
		},
		{
			name:           "exceeds uint8 max",
			complexityStep: 1,
			load:           256,
			expectedResult: math.MaxUint8,
		},
		{
			name:           "large load exceeds uint8",
			complexityStep: 100,
			load:           100000,
			expectedResult: math.MaxUint8,
		},
		{
			name:           "very large load",
			complexityStep: 1,
			load:           math.MaxUint64,
			expectedResult: math.MaxUint8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{ComplexityStep: tt.complexityStep}
			calc := NewAdaptiveComplexityCalculator(cfg)

			result := calc.Calculate(tt.load)

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
