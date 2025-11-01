package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/svalasovich/golang-template/internal/config"
	"github.com/svalasovich/golang-template/internal/log"
)

var Version = "N/A"

var root = &cobra.Command{
	Use:     "golang-template",
	Version: Version,
	Short:   "Golang GitHub template",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if err := config.Init(cfg, cmd); err != nil {
			return fmt.Errorf("failed to init config: %w", err)
		}

		log.Init(cfg.Log)

		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		logger := func() fxevent.Logger {
			return &fxevent.SlogLogger{Logger: slog.New(log.NewComponentLogger("dependency-injection").Handler())}
		}
		fx.New(BootstrapCommon(cmd), BootstrapApp(), fx.WithLogger(logger)).Run()

		return nil
	},
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	config.InitConfigFileFlags(root)
	config.InitLogFlags(root)
	config.InitMonitoringFlags(root)

	cobra.CheckErr(root.ExecuteContext(ctx))
}
