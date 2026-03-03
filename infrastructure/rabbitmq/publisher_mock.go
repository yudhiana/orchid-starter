package rabbitmq

import "context"

type MockPublisher struct {
	Published []Publishing
	Closed    bool
}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		Published: []Publishing{},
	}
}

func (p *MockPublisher) Publish(ctx context.Context, exchange, routingKey string, kind Kind, msg Publishing) error {
	p.Published = append(p.Published, msg)
	return nil
}

func (p *MockPublisher) PublishQueue(ctx context.Context, queue string, msg Publishing) error {
	p.Published = append(p.Published, msg)
	return nil
}

func (p *MockPublisher) Close() error {
	p.Closed = true
	return nil
}
