package handler

import (
	"orchid-starter/internal/bootstrap"

	"github.com/go-chi/chi/v5"
)

func GQLRoutes(app *chi.Mux, container *bootstrap.Container) {
	container.Log.Info("Initialize default handler...")

	// Get DI from container instead of creating new instance
	// di := container.GetDI()

	// app.PartyFunc("/", func(defaultParty iris.Party) {
	// 	gqlHandler.NewDefaultGQLHandler(app, di)
	// })

	// app.PartyFunc("/gql", func(graphHandler iris.Party) {
	// 	graphHandler.Post("/query", gqlHandler.NewGraphHandler(di).GQLHandler())
	// 	graphHandler.Get("/playground", gqlHandler.PlaygroundHandler())
	// })
}
