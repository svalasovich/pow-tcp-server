package log

type Config struct {
	Level      string `mapstructure:"level" validate:"oneof=debug info warn error"`
	ShowSource bool   `mapstructure:"show-source"`
	JSONFormat bool   `mapstructure:"json-format"`
}
