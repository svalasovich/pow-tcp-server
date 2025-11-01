package config

import (
	"github.com/spf13/cobra"
)

const configFlag = "config"

func InitConfigFileFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(configFlag, "c", ".", "set path to config file. Example: etc/.app")
}

func InitLogFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("log.level", "l", "info", "logging level")
	cmd.Flags().Bool("log.show-source", false, "show logs in JSON format")
	cmd.Flags().Bool("log.json-format", false, "show the source code position of the log statement")
}
