package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"orchid-starter/infrastructure/rabbitmq"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/observability/sentry"

	bunker "github.com/yudhiana/bunker/errors"
	"github.com/yudhiana/logos"
)

var ErrNoHandlerRegistered = errors.New("no handler registered event type")

// BaseEventHandler provides common functionality for all event handlers
type BaseEventHandler struct {
	di       *bootstrap.DirectInjection
	log      *logos.LogEntry
	handlers map[string]EventHandlerInterface
	config   EventHandlerConfig
}

// NewBaseEventHandler creates a new base event handler
func NewBaseEventHandler(di *bootstrap.DirectInjection, config EventHandlerConfig) *BaseEventHandler {
	di.Log.Info(fmt.Sprintf("Initialize %s search engine event handler", config.LoggerPrefix))

	handler := &BaseEventHandler{
		di:       di,
		log:      di.Log,
		handlers: make(map[string]EventHandlerInterface),
		config:   config,
	}

	return handler
}

// RegisterHandler registers an event handler for its supported event types
func (h *BaseEventHandler) RegisterHandler(handler EventHandlerInterface) {
	eventTypes := handler.GetEventTypes()
	for _, eventType := range eventTypes {
		if _, exists := h.handlers[eventType]; exists {
			h.log.Warn("Event handler already registered", "event_type", eventType)
			panic(fmt.Errorf("event handler already registered. event_type: %s", eventType))
		}

		h.handlers[eventType] = handler
		h.log.Info("Registered event handler", "event_type", eventType)
	}

	h.log.Info("Event handlers registered successfully", "total_handlers", len(h.handlers))
}

// EventHandler processes search engine events using the registry system
func (h *BaseEventHandler) EventHandler(body map[string]any) (err error) {
	startTime := time.Now()

	var event rabbitmq.EventData
	var processingError error

	// Defer logging and error handling
	defer func() {
		processingTime := time.Since(startTime)

		if processingError != nil {
			if errors.Is(processingError, ErrNoHandlerRegistered) {
				return
			}

			h.log.Error("Event processing failed",
				"event_type", event.EventType,
				"action", "event-processing",
				"error", processingError,
				"processing_time_ms", processingTime.Milliseconds())

			sentry.SentryLogger(processingError, body)

		} else {
			h.log.Info("Event processed successfully",
				"event_type", event.EventType,
				"action", "event-processing",
				"processing_time_ms", processingTime.Milliseconds())
		}

		// Recovery from panics
		if r := recover(); r != nil {
			h.log.Error("Panic occurred during event processing",
				"panic", r,
				"action", "event-processing",
				"event_type", event.EventType,
				"recovery_time", time.Since(startTime))

			sentry.SentryLogger(fmt.Errorf("panic occurred during event processing error: %v", r), body)
		}
	}()

	ctx := context.Background()

	// Route event to registered handler
	processingError = h.routeEvent(ctx, event)
	return
}

// routeEvent routes the event to the registered handler
func (h *BaseEventHandler) routeEvent(ctx context.Context, event rabbitmq.EventData) error {
	handler, exists := h.handlers[event.EventType]

	if !exists {
		return bunker.New(bunker.StatusUnprocessableEntity).SetMessage("no handler registered event type")
	}

	h.log.Info("Processing event", "event_type", event.EventType)
	return handler.Handle(ctx, event)
}

// GetRegisteredEventTypes returns all registered event types
func (h *BaseEventHandler) GetRegisteredEventTypes() []string {
	eventTypes := make([]string, 0, len(h.handlers))
	for eventType := range h.handlers {
		eventTypes = append(eventTypes, eventType)
	}

	return eventTypes
}

// HealthCheck returns the health status of the event handler
func (h *BaseEventHandler) HealthCheck() map[string]any {
	return map[string]any{
		"status":              "healthy",
		"registered_handlers": len(h.handlers),
		"supported_events":    h.GetRegisteredEventTypes(),
	}
}

// GetDI returns the dependency injection container
func (h *BaseEventHandler) GetDI() *bootstrap.DirectInjection {
	return h.di
}
