package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(map[string]any) error

// EventData struct for event data
type EventData struct {
	EventType string     `json:"event_type,omitempty"`
	Data      any        `json:"data,omitempty"`
	Date      *time.Time `json:"date,omitempty"`
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	tag     string
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	conn, errConn := amqp.Dial(amqpURI)
	if errConn != nil {
		return nil, fmt.Errorf("dial connection error: %w", errConn)
	}

	channel, errGetChan := conn.Channel()
	if errGetChan != nil {
		return nil, fmt.Errorf("getting amqp channel error: %w", errGetChan)
	}

	if errExchange := channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); errExchange != nil {
		return nil, fmt.Errorf("declare exchange error: %w", errExchange)
	}

	queue, errQueue := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if errQueue != nil {
		return nil, fmt.Errorf("declare queueu error: %w", errQueue)
	}

	if errBind := channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); errBind != nil {
		return nil, fmt.Errorf("bind queue error: %w", errBind)
	}

	return &Consumer{
		conn:    conn,
		channel: channel,
	}, nil
}

func (c *Consumer) Consume(autoAck bool, handler Handler) error {
	deliveries, errDeliv := c.channel.Consume(
		c.queue.Name, // name
		c.tag,        // consumerTag,
		autoAck,      // autoAck
		false,        // exclusive
		false,        // noLocal
		false,        // noWait
		nil,          // arguments

	)

	if errDeliv != nil {
		return fmt.Errorf("consume queue error: %w", errDeliv)
	}

	var mapData map[string]any
	for d := range deliveries {
		if errUnMarshal := json.Unmarshal(d.Body, &mapData); errUnMarshal != nil {
			return fmt.Errorf("unmarshal queue body error: %w", errUnMarshal)
		}

		if errHandler := handler(mapData); errHandler != nil {
			return fmt.Errorf("handler error: %w", errHandler)
		}
	}

	return nil
}
