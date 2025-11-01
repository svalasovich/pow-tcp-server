package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/svalasovich/pow-tcp-server/internal/config"
	"github.com/svalasovich/pow-tcp-server/internal/log"
)

var root = &cobra.Command{
	Use:     "server",
	Version: "0.0.1",
	Short:   "Golang TCP server with PoW",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if err := config.Init(cfg, cmd); err != nil {
			return fmt.Errorf("failed to init config: %w", err)
		}

		log.Init(cfg.Log)

		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return bootstrap(cmd.Context(), cfg)
	},
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	config.InitConfigFileFlags(root)
	config.InitLogFlags(root)
	InitFlags(root)

	cobra.CheckErr(root.ExecuteContext(ctx))
}
