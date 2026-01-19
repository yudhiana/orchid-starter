package handler

import (
	"context"
	"orchid-starter/infrastructure/rabbitmq"
)

// EventHandlerInterface defines the contract for event handlers
type EventHandlerInterface interface {
	Handle(ctx context.Context, event rabbitmq.EventData) error
	GetEventTypes() []string
}

// EventHandlerConfig holds configuration for event handlers
type EventHandlerConfig struct {
	Name         string
	Alias        string
	Usage        string
	Description  string
	QueueName    string
	ExchangeName string
	LoggerPrefix string
}
