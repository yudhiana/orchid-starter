package publisher

import (
	"context"
	"orchid-starter/infrastructure/rabbitmq"
)

type EventPublisher struct {
	publishing rabbitmq.PublisherInterface
}

func NewEventPublisher(pub rabbitmq.PublisherInterface) *EventPublisher {
	return &EventPublisher{
		publishing: pub,
	}
}

// TODO: put here to add others event publisher
func (p *EventPublisher) PublishExampleCreated(ctx context.Context, exchange, routingKey string, kind rabbitmq.Kind, msg rabbitmq.Publishing) error {
	return p.publishing.Publish(ctx, exchange, routingKey, kind, msg)
}

func (p *EventPublisher) PublishQueueExampleCreated(ctx context.Context, queue string, msg rabbitmq.Publishing) error {
	return p.publishing.PublishQueue(ctx, queue, msg)
}
