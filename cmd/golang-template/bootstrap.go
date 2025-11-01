package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/svalasovich/golang-template/internal/config"
	"github.com/svalasovich/golang-template/internal/http"
	"github.com/svalasovich/golang-template/internal/monitoring"
)

var cfg = new(config.Config)

func BootstrapApp() fx.Option {
	return fx.Provide(
	// Add here app components
	)
}

func BootstrapCommon(cmd *cobra.Command) fx.Option {
	return fx.Module("common",
		fx.Provide(
			func() config.Monitoring { return cfg.Monitoring },
			func() monitoring.Version { return monitoring.Version(cmd.Version) },
			func() monitoring.Name { return monitoring.Name(cmd.Name()) },
			monitoring.NewHealthProvider,
			monitoring.NewMetricProvider,
			monitoring.NewProvider,
			monitoring.NewServer,
		),
		fx.Invoke(invokeServer),
	)
}

func invokeServer(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if err := server.ListenAndServe(); err != nil {
				return fmt.Errorf("failed to start monitoring server: %w", err)
			}

			return nil
		},
		OnStop: server.Shutdown,
	})
}
