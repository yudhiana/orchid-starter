package gqlServer

import (
	"orchid-starter/internal/bootstrap/container"
	gqlHandler "orchid-starter/modules/example/delivery/api/gql"

	"github.com/go-chi/chi/v5"
)

func GQLRoutes(app *chi.Mux, container *container.Container) {
	container.Log.Info("Initialize default handler...")

	// Get DI from container instead of creating new instance
	di := container.GetDI()

	app.Route("/gql", func(graphHandler chi.Router) {
		graphHandler.Post("/query", gqlHandler.NewGraphHandler(di).GQLHandler())
		graphHandler.Get("/playground", gqlHandler.PlaygroundHandler())
	})
}
