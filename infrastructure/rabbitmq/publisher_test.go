package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/rabbitmq/amqp091-go"
)

func TestPublisherQueueRace(t *testing.T) {
	pub := NewPublisher("amqp://guest:guest@localhost:5672")
	defer pub.Close()

	wg := new(sync.WaitGroup)
	for i := range 1000 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := pub.PublishQueue(context.Background(), "orchid-queue", amqp091.Publishing{
				Body: []byte(fmt.Sprintf("index - %d", i)),
			})

			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()

}

func TestPublisherRace(t *testing.T) {
	pub := NewPublisher("amqp://guest:guest@localhost:5672")
	defer pub.Close()

	wg := new(sync.WaitGroup)
	for i := range 1000 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := pub.Publish(context.Background(), "orchid-event", "", Fanout, amqp091.Publishing{
				Body: []byte(fmt.Sprintf("index - %d", i)),
			})

			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()

}
