package handler

import (
	"context"
	"orchid-starter/config"
	"orchid-starter/infrastructure/rabbitmq"
	"orchid-starter/internal/bootstrap/container"

	"github.com/urfave/cli/v3"
)

// HandlerRegistrationFunc defines a function type for registering specific handlers
type HandlerRegistrationFunc func(baseHandler *BaseEventHandler)

// CreateEventHandlerApplication creates a generic CLI application for event handling
func CreateEventHandlerApplication(
	handlerConfig EventHandlerConfig,
	appConfig *config.LocalConfig,
	registerHandlers HandlerRegistrationFunc,
) func(di *container.DirectInjection) *cli.Command {
	return func(di *container.DirectInjection) *cli.Command {
		return &cli.Command{
			Name:        handlerConfig.Name,
			Aliases:     []string{handlerConfig.Alias},
			Usage:       handlerConfig.Usage,
			Description: handlerConfig.Description,
			Action: func(ctx context.Context, cmd *cli.Command) (err error) {
				rabbitConfig := appConfig.RabbitMQConfig
				consumer, errConn := rabbitmq.NewConsumer(
					rabbitConfig.AmqpURI(),
					handlerConfig.ExchangeName,
					string(rabbitmq.Fanout),
					handlerConfig.QueueName,
					"",
					"event_consumer",
				)
				if errConn != nil {
					return errConn
				}

				// Create base handler and register specific handlers
				baseHandler := NewBaseEventHandler(di, handlerConfig)
				registerHandlers(baseHandler)

				consumer.Consume(false, baseHandler.EventHandler)
				return
			},
		}
	}
}
