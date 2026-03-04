package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(Publishing) error

// EventData struct for event data

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

	for delivery := range deliveries {
		errHandler := handler(Publishing{
			AppId:           delivery.AppId,
			UserId:          delivery.UserId,
			MessageId:       delivery.MessageId,
			CorrelationId:   delivery.CorrelationId,
			Headers:         delivery.Headers,
			ReplyTo:         delivery.ReplyTo,
			Expiration:      delivery.Expiration,
			Type:            delivery.Type,
			Body:            delivery.Body,
			ContentType:     delivery.ContentType,
			ContentEncoding: delivery.ContentEncoding,
			DeliveryMode:    delivery.DeliveryMode,
			Timestamp:       delivery.Timestamp,
			Priority:        delivery.Priority,
		})

		if !autoAck {
			if errHandler != nil {
				if errNack := delivery.Nack(false, true); errNack != nil {
					return fmt.Errorf("nack error: %w", errNack)
				}
				continue
			}

			if errAck := delivery.Ack(false); errAck != nil {
				return fmt.Errorf("ack error: %w", errAck)
			}
		}
	}

	return nil
}
