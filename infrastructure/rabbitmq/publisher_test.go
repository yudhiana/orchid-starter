package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPublisherQueueRace(t *testing.T) {
	pub := NewPublisher("amqp://guest:guest@localhost:5672")
	defer pub.Close()

	wg := new(sync.WaitGroup)
	for i := range 1000 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := pub.PublishQueue(context.Background(), "orchid-queue", Publishing{
				Body:      []byte(fmt.Sprintf("index - %d", i)),
				Timestamp: time.Now().UTC(),
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
			err := pub.Publish(context.Background(), "orchid-event", "", Fanout, Publishing{
				Body:      []byte(fmt.Sprintf("index - %d", i)),
				Timestamp: time.Now().UTC(),
			})

			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()

}
