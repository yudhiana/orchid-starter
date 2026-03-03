// main.go for CLI app
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"orchid-starter/cmd/cli/commands"
	"orchid-starter/config"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/observability/sentry"

	_ "github.com/joho/godotenv/autoload"
	"github.com/yudhiana/logos"
)

func main() {
	di, err := bootstrap.NewDirectInjection(config.GetLocalConfig())
	if err != nil {
		panic("Failed to initialize dependencies: " + err.Error())
	}
	defer di.Close() // Ensure cleanup even if app panics

	sentry.InitSentry()

	go func() {
		// Wait for interrupt signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logos.NewLogger().Info("🔴Shutting Cli Application...")

		// Graceful shutdown timeout
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// perform close in goroutine so we can respect timeout
		done := make(chan error, 1)
		go func() {
			done <- di.Close()
		}()

		select {
		case err := <-done:
			if err != nil {
				logos.NewLogger().Error("🔴Forced shutdown", "err", err)
			}
		case <-ctx.Done():
			logos.NewLogger().Error("🔴Shutdown timed out", "err", ctx.Err())
		}

		logos.NewLogger().Info("🔴Cli Application exited properly")
	}()

	app := commands.NewBaseCommand(di).GetCommands()
	err = app.Run(context.Background(), os.Args)
	if err != nil {
		panic(err.Error())
	}
}
