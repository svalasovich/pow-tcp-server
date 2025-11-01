package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/svalasovich/pow-tcp-server/internal/log"
)

type (
	Config struct {
		Log log.Config
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
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	))

	return viper.Unmarshal(cfg, hooks)
}
