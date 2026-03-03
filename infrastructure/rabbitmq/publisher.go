package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yudhiana/logos"
)

var (
	publisherInstance *Publisher
	publisherMu       sync.Mutex
)

type Publishing amqp.Publishing

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

	uri    string
	closed bool
}

type PublisherInterface interface {
	Publish(ctx context.Context, exchange, routingKey string, kind Kind, msg Publishing) error
	PublishQueue(ctx context.Context, queue string, msg Publishing) error
	Close() error
}

// NewPublisher opens a connection and channel, declares the exchange and
// returns a publisher instance. exchangeType is typically "direct"
// / "topic" / "fanout" etc.
func NewPublisher(amqpURI string) PublisherInterface {
	publisherMu.Lock()
	defer publisherMu.Unlock()

	if publisherInstance != nil {
		return publisherInstance
	}

	pub := &Publisher{
		uri: amqpURI,
	}

	if err := pub.connect(); err != nil {
		panic(err)
	}

	publisherInstance = pub
	return publisherInstance
}

func (p *Publisher) connect() error {
	conn, err := amqp.Dial(p.uri)
	logos.NewLogger().Info("dialing", "URI", p.uri)
	if err != nil {
		return fmt.Errorf("dial connection error: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("getting amqp channel error: %w", err)
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

// Publish sends it with the given routing key.
func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, kind Kind, msg Publishing) error {
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

	if err := p.channel.ExchangeDeclare(
		exchange,     // name
		string(kind), // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		p.channel.Close()
		p.conn.Close()
		return fmt.Errorf("declare exchange error: %w", err)
	}

	if err := p.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing(msg),
	); err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	logos.NewLogger().Info("publishing", "message", msg)
	return nil
}

// Publish sends it with the given routing key.
func (p *Publisher) PublishQueue(ctx context.Context, queue string, msg Publishing) error {
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

	_, err := p.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		p.channel.Close()
		p.conn.Close()
		return fmt.Errorf("declare exchange error: %w", err)
	}

	if err := p.channel.PublishWithContext(
		ctx,
		"", // exchange
		queue,
		false, // mandatory
		false, // immediate
		amqp.Publishing(msg),
	); err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	logos.NewLogger().Info("publishing", "message", msg)
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
