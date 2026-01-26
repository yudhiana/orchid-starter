package InitHandler

import (
	"orchid-starter/cmd/cli/handler"
	"orchid-starter/internal/bootstrap"
	defaultEventHandler "orchid-starter/modules/default/delivery/event"

	"github.com/urfave/cli"
)

// NewApplication creates a CLI application for company event handling
func NewApplication(di *bootstrap.DirectInjection) cli.Command {
	config := handler.EventHandlerConfig{
		Name:         "cli-init-handler",
		Alias:        "cih",
		Usage:        "run cli-init-handler",
		Description:  "cli-init-handler",
		QueueName:    "cli-init-event-queue",
		ExchangeName: "event",
		LoggerPrefix: "init-handler",
	}

	return handler.CreateEventHandlerApplication(config, registerHandlers)(di)
}

// registerHandlers registers company-specific event handlers
func registerHandlers(baseHandler *handler.BaseEventHandler) {
	// Initialize event handler
	defaultHandler := defaultEventHandler.NewDefaultEventHandler(baseHandler.GetDI())
	baseHandler.RegisterHandler(defaultHandler)
}
