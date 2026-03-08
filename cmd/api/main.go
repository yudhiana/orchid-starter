package main

import (
	"log"

	"orchid-starter/internal/bootstrap/container"
	restfulServer "orchid-starter/internal/bootstrap/server/restful-server"
	"orchid-starter/observability/prometheus"
	"orchid-starter/observability/sentry"

	_ "github.com/joho/godotenv/autoload"
)

func init() {
	prometheus.InitPrometheus()
}

func main() {
	// Initialize bootstrap container
	container, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application container: %v", err)
	}
	defer container.Close()

	sentry.InitSentry()

	// Initialize and start server
	srv := restfulServer.NewServer(container)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
