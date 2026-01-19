package handler

import (
	"orchid-starter/infrastructure/rabbitmq"
	"orchid-starter/internal/bootstrap"

	"github.com/urfave/cli"
)

// HandlerRegistrationFunc defines a function type for registering specific handlers
type HandlerRegistrationFunc func(baseHandler *BaseEventHandler)

// CreateEventHandlerApplication creates a generic CLI application for event handling
func CreateEventHandlerApplication(
	config EventHandlerConfig,
	registerHandlers HandlerRegistrationFunc,
) func(di *bootstrap.DirectInjection) cli.Command {
	return func(di *bootstrap.DirectInjection) cli.Command {
		return cli.Command{
			Name:        config.Name,
			Aliases:     []string{config.Alias},
			Usage:       config.Usage,
			Description: config.Description,
			Action: func(ctx *cli.Context) (err error) {
				consumer, errConn := rabbitmq.NewConsumer("localhost:5672", "events_exchange", "fanout", "events_queue", "#", "event_consumer")
				if errConn != nil {
					return errConn
				}

				// Create base handler and register specific handlers
				baseHandler := NewBaseEventHandler(di, config)
				registerHandlers(baseHandler)

				consumer.Consume(false, baseHandler.EventHandler)
				return
			},
		}
	}
}
