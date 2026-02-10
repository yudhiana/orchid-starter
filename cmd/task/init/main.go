package initTaskApplication

import (
	"context"
	"orchid-starter/cmd/task/init/handler/task"
	modelInit "orchid-starter/cmd/task/init/model"
	"orchid-starter/internal/bootstrap"

	"github.com/urfave/cli/v3"
)

func NewInitTask(di *bootstrap.DirectInjection) *cli.Command {
	return &cli.Command{
		Name:    "init-task",
		Aliases: []string{"init-task"},
		Usage:   "Run init-task",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:  "id",
				Value: 1,
			},
			&cli.BoolFlag{
				Name: "count",
			},
			&cli.StringFlag{
				Name:     "name",
				Required: false,
			},
		},
		Action: func(c context.Context, cmd *cli.Command) error {
			return task.NewTask(di, modelInit.Init{
				ID: cmd.Uint64("id"),
			}).Start()
		},
	}
}
