package pow

type (
	ComplexityCalculator interface {
		Calculate(load uint64) uint8
	}

	EngineInterface interface {
		Verify(data []byte, nonce []byte) bool
		GenerateData(complexity uint8) []byte
	}

	Service struct {
		complexityCalculator ComplexityCalculator
		engine               EngineInterface
	}
)

func NewService(engine EngineInterface, complexityCalculator ComplexityCalculator) *Service {
	return &Service{
		engine:               engine,
		complexityCalculator: complexityCalculator,
	}
}

func (s *Service) GenerateData(load uint64) (uint8, []byte) {
	complexity := s.complexityCalculator.Calculate(load)

	return complexity, s.engine.GenerateData(complexity)
}

func (s *Service) Verify(data []byte, nonce []byte) bool {
	return s.engine.Verify(data, nonce)
}
