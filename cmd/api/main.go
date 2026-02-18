package main

import (
	"log"

	"orchid-starter/internal/bootstrap"
	"orchid-starter/internal/bootstrap/server"
	"orchid-starter/observability/prometheus"
	"orchid-starter/observability/sentry"

	_ "github.com/joho/godotenv/autoload"
)

func init() {
	prometheus.InitPrometheus()
}

func main() {
	// Initialize bootstrap container
	container, err := bootstrap.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application container: %v", err)
	}
	defer container.Close()

	sentry.InitSentry()

	// Initialize and start server
	srv := server.NewServer(container)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
