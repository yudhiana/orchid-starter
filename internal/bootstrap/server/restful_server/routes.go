package restfulServer

import (
	"orchid-starter/internal/bootstrap/container"

	"github.com/go-chi/chi/v5"
)

// RouteSetup represents a function that sets up routes
type RouteSetup func(app *chi.Mux, container *container.Container)

// SetupAllRoutes configures all application routes in an organized manner
func SetupAllRoutes(app *chi.Mux, container *container.Container, routeSetups ...RouteSetup) {
	container.Log.Info("Setting up all application routes...")

	// Execute all route setups
	for _, setupFunc := range routeSetups {
		setupFunc(app, container)
	}

	container.Log.Info("All application routes configured successfully")
}
