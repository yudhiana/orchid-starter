package main

import (
	"log"
	"orchid-starter/internal/bootstrap/container"
	gqlServer "orchid-starter/internal/bootstrap/server/gql_server"
	"orchid-starter/observability/sentry"
)

func main() {
	// Initialize bootstrap container
	container, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application container: %v", err)
	}
	defer container.Close()

	sentry.InitSentry()

	// Initialize and start server
	srv := gqlServer.NewGQLServer(container)
	log.Println("Starting server...")
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
