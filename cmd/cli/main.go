// main.go for CLI app
package main

import (
	"context"
	"net/mail"
	"os"

	InitHandler "orchid-starter/cmd/cli/handler/init"
	"orchid-starter/config"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/observability/sentry"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "Command execution for Go API CLI",
		Usage:   "Run task by command CLI for Golang",
		Version: "1.0.0",
		Authors: []any{
			mail.Address{Name: "yudhiana", Address: "yudhiana@orchid-starter.co"},
		},
	}

	di, err := bootstrap.NewDirectInjection(config.GetLocalConfig())
	if err != nil {
		panic("Failed to initialize dependencies: " + err.Error())
	}
	defer di.Close() // Ensure cleanup even if app panics

	sentry.InitSentry()

	app.Commands = []*cli.Command{
		InitHandler.NewApplication(di),
		// TODO : add other commands
	}

	err = app.Run(context.Background(), os.Args)
	if err != nil {
		panic(err.Error())
	}
}
