package config

import (
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Log        Log
		Monitoring Monitoring
	}

	Log struct {
		Level      string `mapstructure:"level" validate:"oneof=debug info warn error"`
		ShowSource bool   `mapstructure:"show-source"`
		JSONFormat bool   `mapstructure:"json-format"`
	}

	Monitoring struct {
		Server HTTPServer
	}

	HTTPServer struct {
		Port uint `mapstructure:"port" validate:"max=65535"`
	}

	GRPCServer struct {
		Port uint `mapstructure:"port" validate:"max=65535"`
	}
)

func Init(cfg any, cmd *cobra.Command) error {
	readEnvVariables()

	if err := readCmdVariables(cmd); err != nil {
		return fmt.Errorf("failed to read command line variables: %w", err)
	}

	if err := readConfigFile(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func readEnvVariables() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "__"))
	viper.AutomaticEnv()
}

func readCmdVariables(cmd *cobra.Command) error {
	var err error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err = viper.BindPFlag(f.Name, f); err != nil {
			return
		}
	})
	if err != nil {
		return err
	}

	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if err = viper.BindPFlag(f.Name, f); err != nil {
			return
		}
	})

	return err
}

func readConfigFile() error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(viper.GetString(configFlag))

	if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return err
	}

	return nil
}

func unmarshal(cfg any) error {
	hooks := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		bytesHexHook,
		mapAsStringHook,
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	))

	return viper.Unmarshal(cfg, hooks)
}

func bytesHexHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() != reflect.String {
		return data, nil
	}

	if to != reflect.SliceOf(reflect.TypeOf(byte(0))) {
		return data, nil
	}

	dataHex := data.(string)
	if dataHex == "" {
		return ([]byte)(nil), nil
	}

	dataHex = strings.TrimPrefix(dataHex, "0x")

	return hex.DecodeString(dataHex)
}

func mapAsStringHook(from reflect.Kind, to reflect.Kind, data interface{}) (interface{}, error) {
	if from != reflect.String || to != reflect.Map {
		return data, nil
	}

	dataString := data.(string)
	dataString = strings.Trim(dataString, "[]")
	if dataString == "" {
		return map[string]any{}, nil
	}

	reader := csv.NewReader(strings.NewReader(dataString))
	record, err := reader.Read()
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(record))
	for _, pair := range record {
		key, value, found := strings.Cut(pair, "=")
		if !found {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}
