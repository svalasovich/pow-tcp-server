package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/svalasovich/pow-tcp-server/internal/config"
	"github.com/svalasovich/pow-tcp-server/pkg/common/pow"
	"github.com/svalasovich/pow-tcp-server/pkg/common/tcp"
)

var cfg = new(Config)

type Config struct {
	PoW pow.Config `mapstructure:"pow"`
	TCP tcp.Config `mapstructure:"tcp"`

	config.Config `mapstructure:",squash"`
}

func InitFlags(cmd *cobra.Command) {
	cmd.Flags().Uint("pow.complexity-step", 8, "step of change in the complexity of PoW with increasing requests")

	cmd.Flags().Duration("tcp.deadline", 10*time.Second, "server connection deadline")
	cmd.Flags().Duration("tcp.keep-alive", 10*time.Second, "server connection keep alive")
	cmd.Flags().String("tcp.address", ":9000", "server port")
}
