package handler

import (
	"orchid-starter/internal/bootstrap"

	api "orchid-starter/modules/default/delivery/api/rest"

	"github.com/go-chi/chi/v5"
)

// SetupDefaultRoutes configures the default routes with proper dependency injection
func SetupDefaultRoutes(app *chi.Mux, container *bootstrap.Container) {
	container.Log.Info("Initialize default handler...")

	// Get DI from container instead of creating new instance
	di := container.GetDI()

	app.Group(func(defaultParty chi.Router) {
		api.NewDefaultHandler(defaultParty, di)
	})
}
