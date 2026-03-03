package commands

import (
	"net/mail"
	InitHandler "orchid-starter/cmd/cli/handler/init"
	"orchid-starter/internal/bootstrap"

	"github.com/urfave/cli/v3"
)

type BaseCommand struct {
	*cli.Command
	*bootstrap.DirectInjection
}

func NewBaseCommand(di *bootstrap.DirectInjection) *BaseCommand {
	app := &cli.Command{
		Name:    "Command execution for Go API CLI",
		Usage:   "Run task by command CLI for Golang",
		Version: "1.0.0",
		Authors: []any{
			mail.Address{Name: "yudhiana", Address: "yudhiana@orchid-starter.co"},
		},
	}

	return &BaseCommand{
		Command:         app,
		DirectInjection: di,
	}
}

func (base *BaseCommand) GetCommands() *BaseCommand {

	base.Commands = []*cli.Command{
		InitHandler.NewApplication(base.DirectInjection),
		// TODO : add other commands
	}
	return base
}
