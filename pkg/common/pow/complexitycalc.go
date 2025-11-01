package pow

import "math"

type AdaptiveComplexityCalculator struct {
	cfg Config
}

func NewAdaptiveComplexityCalculator(cfg Config) *AdaptiveComplexityCalculator {
	return &AdaptiveComplexityCalculator{
		cfg: cfg,
	}
}

func (a *AdaptiveComplexityCalculator) Calculate(load uint64) uint8 {
	calculatedComplexity := load / uint64(a.cfg.ComplexityStep)
	if calculatedComplexity > math.MaxUint8 {
		return math.MaxUint8
	}

	return uint8(calculatedComplexity)
}
