// main.go for CLI app
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"orchid-starter/cmd/cli/commands"
	"orchid-starter/config"
	"orchid-starter/internal/bootstrap/container"
	"orchid-starter/observability/sentry"

	_ "github.com/joho/godotenv/autoload"
	"github.com/yudhiana/logos"
)

func main() {
	di, err := container.NewApplicationContainer(config.GetLocalConfig())
	if err != nil {
		panic("Failed to initialize dependencies: " + err.Error())
	}

	sentry.InitSentry()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := commands.NewBaseCommand(di).GetCommands()
	done := make(chan error, 1)

	go func() {
		done <- app.Run(ctx, os.Args)
	}()

	select {
	case err := <-done:
		if err != nil {
			logos.NewLogger().Error("Application error", "err", err)
		}
	case <-ctx.Done():
		logos.NewLogger().Info("🔴Signal received, stopping CLI...")
	}

	logos.NewLogger().Info("🔴Shutting down CLI Application...")

	if err := di.Close(); err != nil {
		logos.NewLogger().Error("Shutdown error", "err", err)
	}

	logos.NewLogger().Info("🔴Cli Application exited properly")
}
