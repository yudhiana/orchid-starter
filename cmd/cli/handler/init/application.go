package InitHandler

import (
	"orchid-starter/cmd/cli/handler"
	"orchid-starter/config"
	"orchid-starter/internal/bootstrap/container"
	"orchid-starter/modules/example/delivery/event/subscriber"

	"github.com/urfave/cli/v3"
)

// NewApplication creates a CLI application for company event handling
func NewApplication(di *container.DirectInjection) *cli.Command {
	handlerConfig := handler.EventHandlerConfig{
		Name:         "cli-init-handler",
		Alias:        "cih",
		Usage:        "run cli-init-handler",
		Description:  "cli-init-handler",
		QueueName:    "orchid-queue",
		ExchangeName: "orchid-event",
		LoggerPrefix: "init-handler",
	}

	appConfig := config.GetLocalConfig()
	return handler.CreateEventHandlerApplication(handlerConfig, appConfig, registerHandlers)(di)
}

// registerHandlers registers company-specific event handlers
func registerHandlers(baseHandler *handler.BaseEventHandler) {
	// Initialize event handler
	exampleHandler := subscriber.NewExampleEventHandler(baseHandler.GetDI())
	baseHandler.RegisterHandler(exampleHandler)
}
