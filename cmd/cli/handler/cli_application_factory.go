package handler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"orchid-starter/internal/bootstrap"
	"orchid-starter/internal/common"

	"github.com/mataharibiz/sange/v2"
	"github.com/yudhiana/logos"
	"github.com/urfave/cli"
)

// HandlerRegistrationFunc defines a function type for registering specific handlers
type HandlerRegistrationFunc func(baseHandler *BaseEventHandler)

// CreateEventHandlerApplication creates a generic CLI application for event handling
func CreateEventHandlerApplication(
	config EventHandlerConfig,
	registerHandlers HandlerRegistrationFunc,
) func(di *bootstrap.DirectInjection) cli.Command {
	return func(di *bootstrap.DirectInjection) cli.Command {
		return cli.Command{
			Name:        config.Name,
			Aliases:     []string{config.Alias},
			Usage:       config.Usage,
			Description: config.Description,
			Action: func(ctx *cli.Context) (err error) {
				// Setup exchange configuration
				exchange := sange.ExchangeConfigurations{
					Name:    sange.GetEnv("DMP_EXCHANGE", config.ExchangeName),
					Kind:    sange.Fanout,
					Durable: true,
				}

				// Setup queue configuration
				queue := sange.QueueConfigurations{
					Name:       config.QueueName,
					Durable:    true,
					AutoDelete: true,
				}

				preCount, _ := strconv.Atoi(os.Getenv("RMQ_PREFETCH_COUNT"))

				// Setup context and signal
				goCtx, cancel := context.WithCancel(context.Background())
				defer cancel()

				sigs := make(chan os.Signal, 1)
				signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

				// Create base handler and register specific handlers
				baseHandler := NewBaseEventHandler(di, config)
				registerHandlers(baseHandler)

				wg := new(sync.WaitGroup)

				// Setup consumer
				consumer := sange.NewConsumer().SetPrefetchCount(preCount)
				consumer.SubscribeWithGf(
					goCtx,
					wg,
					exchange,
					queue,
					baseHandler.SearchEngineEventHandler,
				)

				// Wait for shutdown signal
				<-sigs
				logos.NewLogger().Info("Received shutdown signal. Completing pending tasks...")

				// Step 1: Stop receiving new messages
				cancel()
				consumer.Client.Close()

				// Step 2: Wait for consumers
				done := make(chan struct{})
				go func() {
					wg.Wait()               // Wait for consumers
					sange.FlushPublishers() // Wait for publishers
					close(done)
				}()

				// Timeout handling
				timeout := common.GetIntEnv("GRACEFUL_SHUTDOWN_TIMEOUT", 10)
				select {
				case <-done:
					logos.NewLogger().Info("All tasks completed successfully")
				case <-time.After(time.Duration(timeout) * time.Second):
					logos.NewLogger().Info(fmt.Sprintf("Force shutdown timeout of %v elapsed, forcing exit", timeout))
				}

				// Step 3: cleanup
				errClose := di.Close()
				if errClose != nil {
					logos.NewLogger().Error("Failed to close DI", "error", errClose)
				}
				return
			},
		}
	}
}
