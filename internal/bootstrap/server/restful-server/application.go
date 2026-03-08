package restfulServer

import (
	"orchid-starter/internal/bootstrap/container"
	api "orchid-starter/modules/example/delivery/api/rest"

	"github.com/go-chi/chi/v5"
)

// SetupExampleRoutes configures the example routes with proper dependency injection
func SetupExampleRoutes(app *chi.Mux, container *container.Container) {
	container.Log.Info("Initialize example handler...")

	// Get DI from container instead of creating new instance
	di := container.GetDI()

	app.Group(func(exampleParty chi.Router) {
		api.NewExampleHandler(exampleParty, di)
	})
}
