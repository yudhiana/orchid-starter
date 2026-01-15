package server

import (
	"net/http"
	"orchid-starter/config"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/internal/bootstrap/server/applications/handler"

	"github.com/go-chi/chi/v5"
	"github.com/yudhiana/logos"
)

type Server struct {
	cfg *config.LocalConfig
	app *chi.Mux
}

func NewServer(container *bootstrap.Container) *Server {
	srv := &Server{
		cfg: container.GetConfig(),
		app: container.GetApp(),
	}

	// Setup Global middlewares before server initialization
	// srv.setupMiddlewares(middleware.SetAPIVersion, middleware.Debug, middleware.Prometheus)

	// Setup routes after server initialization
	srv.setupRoutes(container)

	return srv
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes(container *bootstrap.Container) {
	// Use centralized route management

	// Define all route setup functions
	routeSetups := []handler.RouteSetup{
		handler.SetupDefaultRoutes, // Root routes
		// Add more route setups here as your application grows
	}

	if s.cfg.AppEnv == "dev" {
		routeSetups = append(routeSetups, handler.GQLRoutes)
	}

	handler.SetupAllRoutes(s.app, container, routeSetups...)
}

// func (s *Server) setupMiddlewares(middlewares ...func(iris.Context)) {
// 	for _, middleware := range middlewares {
// 		s.app.Use(middleware)
// 	}
// }

func (s *Server) Run() error {
	// Log version and server info
	cfg := s.cfg
	logos.NewLogger().Info("Version", "version", cfg.AppVersion, "host", cfg.AppHost, "port", cfg.AppPort)

	// Start server
	address := cfg.AppHost + ":" + cfg.AppPort
	return http.ListenAndServe(address, s.app)
}
