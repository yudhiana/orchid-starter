package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"orchid-starter/internal/bootstrap"
	"orchid-starter/internal/common"
	"orchid-starter/observability/sentry"

	"github.com/mataharibiz/sange/v2"
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

// SearchEngineEventHandler processes search engine events using the registry system
func (h *BaseEventHandler) SearchEngineEventHandler(body map[string]any) {
	startTime := time.Now()

	var event sange.EventData
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
			h.retryEvent(event, processingError)

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
			h.retryEvent(event, fmt.Errorf("panic occurred during event processing error: %v", r))
		}
	}()

	// Parse event data
	if processingError = h.parseEventData(body, &event); processingError != nil {
		return
	}

	// Validate parsed event
	if processingError = h.validateEvent(&event); processingError != nil {
		return
	}

	ctx := context.Background()

	// Route event to registered handler
	processingError = h.routeEvent(ctx, event)
}

// parseEventData safely parses the event body into EventData struct
func (h *BaseEventHandler) parseEventData(body map[string]any, event *sange.EventData) error {
	// Marshal the body to JSON
	dataMarshal, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal event body: %w", err)
	}

	// Unmarshal into EventData struct
	if err := json.Unmarshal(dataMarshal, event); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return nil
}

// validateEvent validates the incoming event data
func (h *BaseEventHandler) validateEvent(event *sange.EventData) error {
	if event == nil {
		return fmt.Errorf("event data is nil")
	}

	if event.EventType == "" {
		return fmt.Errorf("event type is empty")
	}

	return nil
}

// routeEvent routes the event to the registered handler
func (h *BaseEventHandler) routeEvent(ctx context.Context, event sange.EventData) error {
	handler, exists := h.handlers[event.EventType]

	if !exists {
		return sange.NewError(sange.NotFound, ErrNoHandlerRegistered, "no handler registered event type", "event_type", event.EventType)
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

func (h *BaseEventHandler) retryEvent(event sange.EventData, err error) {
	// retry event to global queue
	event.RmqRetryConfig.ErrorMessage = err.Error()
	retryEvent := sange.EventData{
		EventType: "retry_event",
		Data:      event,
	}

	if h.isIgnoreEventRetry(event) {
		return
	}

	delayMs := common.GetInt64Env("RETRY_DELAY_MS", 3000)
	retryEvent.PublishToDelayExchange("dmp_retry", delayMs)
}

func (h *BaseEventHandler) isIgnoreEventRetry(event sange.EventData) (ok bool) {
	ignoreEvents := map[string]bool{
		"search-keyword-store-cron":          true,
		"search-keyword-remove-cron":         true,
		"company-updated-institution":        true,
		"product-status-pkp-updated":         true,
		"product-category-updated":           true,
		"hotfix-product-status-pkp-rollback": true,
	}
	return ignoreEvents[event.EventType]
}
