package tcp

import "time"

type Config struct {
	Deadline  time.Duration `mapstructure:"deadline" validate:"min=1s,max=120s"`
	KeepAlive time.Duration `mapstructure:"keep-alive" validate:"min=1s,max=120s"`
	Address   string        `mapstructure:"address" validate:"required"`
}
