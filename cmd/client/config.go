package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/svalasovich/pow-tcp-server/internal/config"
	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

var cfg = new(Config)

type Config struct {
	Interval time.Duration `mapstructure:"interval"`
	TCP      tcp.Config    `mapstructure:"tcp"`

	config.Config `mapstructure:",squash"`
}

func InitFlags(cmd *cobra.Command) {
	cmd.Flags().Duration("interval", time.Millisecond, "client fetch quote interval")

	cmd.Flags().Duration("tcp.deadline", 10*time.Second, "client connection deadline")
	cmd.Flags().Duration("tcp.keep-alive", 10*time.Second, "client connection keep alive")
	cmd.Flags().String("tcp.address", "localhost:9000", "server address")
}
