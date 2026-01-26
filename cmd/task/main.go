package main

import (
	"os"
	"sort"

	"orchid-starter/cmd/task/init"
	"orchid-starter/config"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/observability/sentry"

	"github.com/urfave/cli"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()
	app.Name = "Command execution for Go API CLI"
	app.Usage = "Run task by command CLI for Golang"
	app.Author = "backend 2025"
	app.Version = "1.0.0"

	sentry.InitSentry()

	di, err := bootstrap.NewDirectInjection(config.GetLocalConfig())
	if err != nil {
		panic("Failed to initialize dependencies: " + err.Error())
	}
	defer di.Close() // Ensure cleanup even if app panics

	app.Commands = []cli.Command{
		initTaskApplication.NewInitTask(di),
		// TODO : add other commands
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	err = app.Run(os.Args)
	if err != nil {
		panic(err.Error())
	}
}
