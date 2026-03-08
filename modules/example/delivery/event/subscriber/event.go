package subscriber

import (
	"context"
	"fmt"
	"orchid-starter/infrastructure/rabbitmq"
	"orchid-starter/internal/bootstrap/container"

	bunker "github.com/yudhiana/bunker/errors"
	"github.com/yudhiana/logos"
)

// Company event type constants
const (
	EventexampleName = "example-init-event-name"
)

type eventHandler struct {
	di  *container.DirectInjection
	log *logos.LogEntry
}

func NewExampleEventHandler(di *container.DirectInjection) *eventHandler {
	return &eventHandler{
		di:  di,
		log: logos.NewLogger(),
	}
}

// Handle processes example init events based on event type
func (eh *eventHandler) Handle(ctx context.Context, event rabbitmq.Publishing) error {
	eh.log.Info("Processing example init event", "event_type", event.Type)
	switch event.Type {
	case EventexampleName:
		return eh.exampleInitEvent(ctx, event)
	default:
		return bunker.New(bunker.StatusUnprocessableEntity).SetMessage(fmt.Sprintf("unknown event type: %s", event.Type))
	}
}

// GetEventTypes returns the list of event types this handler supports
func (eh *eventHandler) GetEventTypes() []string {
	return []string{
		EventexampleName,
	}
}

func (eh *eventHandler) exampleInitEvent(ctx context.Context, event rabbitmq.Publishing) error {
	eh.log.Info("event example successfully executed")
	return nil
}
