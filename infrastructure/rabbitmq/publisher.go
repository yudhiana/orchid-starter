package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	publisherInstance *Publisher
	publisherMu       sync.Mutex
)

type Kind string

const (
	Fanout  Kind = "fanout"
	Topic   Kind = "topic"
	Direct  Kind = "direct"
	Headers Kind = "headers"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	mu sync.Mutex
	wg sync.WaitGroup

	exchange string
	uri      string

	kind   Kind
	closed bool
}

// NewPublisher opens a connection and channel, declares the exchange and
// returns a publisher instance. exchangeType is typically "direct"
// / "topic" / "fanout" etc.
func NewPublisher(amqpURI, exchange string, exchangeType Kind) (*Publisher, error) {
	publisherMu.Lock()
	defer publisherMu.Unlock()

	if publisherInstance != nil {
		return publisherInstance, nil
	}

	pub := &Publisher{
		uri:      amqpURI,
		exchange: exchange,
		kind:     exchangeType,
	}

	if err := pub.connect(); err != nil {
		return nil, err
	}

	publisherInstance = pub
	return publisherInstance, nil
}

func (p *Publisher) connect() error {
	conn, err := amqp.Dial(p.uri)
	if err != nil {
		return fmt.Errorf("dial connection error: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("getting amqp channel error: %w", err)
	}

	if err = ch.ExchangeDeclare(
		p.exchange,     // name
		string(p.kind), // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("declare exchange error: %w", err)
	}

	p.conn = conn
	p.channel = ch

	return nil
}

func (p *Publisher) ensureConnection() error {
	if p.conn != nil && !p.conn.IsClosed() {
		return nil
	}

	return p.connect()
}

// Publish marshals msg to JSON and sends it with the given routing key.
func (p *Publisher) Publish(ctx context.Context, routingKey string, msg any) error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return errors.New("publisher already closed")
	}

	p.wg.Add(1)
	p.mu.Unlock()
	defer p.wg.Done()

	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.ensureConnection(); err != nil {
		return err
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	if err = p.channel.PublishWithContext(
		ctx,
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now().UTC(),
		},
	); err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	return nil
}

// Close releases the channel and connection.
func (p *Publisher) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}

	p.closed = true
	p.mu.Unlock()

	// Wait for in-flight publishes without holding the lock
	p.wg.Wait()

	// Re-acquire lock for cleanup
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel != nil {
		if errCh := p.channel.Close(); errCh != nil {
			return errCh
		}
	}

	if p.conn != nil {
		if errConn := p.conn.Close(); errConn != nil {
			return errConn
		}
	}
	return nil
}
