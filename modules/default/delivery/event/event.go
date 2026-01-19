package event

import (
	"context"
	"fmt"
	"orchid-starter/infrastructure/rabbitmq"
	"orchid-starter/internal/bootstrap"

	bunker "github.com/yudhiana/bunker/errors"
	"github.com/yudhiana/logos"
)

// Company event type constants
const (
	EventDefaultName = "default-init-event-name"
)

type eventHandler struct {
	di  *bootstrap.DirectInjection
	log *logos.LogEntry
}

func NewDefaultEventHandler(di *bootstrap.DirectInjection) *eventHandler {
	return &eventHandler{
		di:  di,
		log: di.Log,
	}
}

// Handle processes default init events based on event type
func (eh *eventHandler) Handle(ctx context.Context, event rabbitmq.EventData) error {
	eh.log.Info("Processing default init event", "event_type", event.EventType)
	switch event.EventType {
	case EventDefaultName:
		return eh.DefaultInitEvent(ctx, event)
	default:
		return bunker.New(bunker.StatusUnprocessableEntity).SetMessage(fmt.Sprintf("unknown event type: %s", event.EventType))
	}
}

// GetEventTypes returns the list of event types this handler supports
func (eh *eventHandler) GetEventTypes() []string {
	return []string{
		EventDefaultName,
	}
}

func (eh *eventHandler) DefaultInitEvent(ctx context.Context, event rabbitmq.EventData) error {
	eh.log.Info("event default successfully executed")
	return nil
}
