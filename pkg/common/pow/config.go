package pow

type Config struct {
	ComplexityStep uint `mapstructure:"complexity-step" validate:"required"`
}
