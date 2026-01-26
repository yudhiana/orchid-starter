package gqlServer

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

func NewGQLServer(container *bootstrap.Container) *Server {
	srv := &Server{
		cfg: container.GetConfig(),
		// app: container.GetApp(),
	}

	// Setup Global middlewares before server initialization
	// srv.setupMiddlewares(middleware.Debug)

	// Setup routes after server initialization
	srv.setupRoutes(container)
	return srv
}

func (s *Server) Run() error {
	// Log version and server info
	cfg := s.cfg
	logos.NewLogger().Info("Version", "version", cfg.AppVersion, "host", cfg.AppHost, "port", cfg.AppPort)

	// Start server
	address := cfg.AppHost + ":" + cfg.AppPort
	return http.ListenAndServe(address, s.app)
}

func (s *Server) setupRoutes(container *bootstrap.Container) {
	// Use centralized route management

	// Define all route setup functions
	routeSetups := []handler.RouteSetup{
		handler.SetupDefaultRoutes, // Root routes
		handler.GQLRoutes,
		// Add more route setups here as your application grows
	}
	handler.SetupAllRoutes(s.app, container, routeSetups...)
}

// func (s *Server) setupMiddlewares(middlewares ...func(iris.Context)) {
// 	for _, middleware := range middlewares {
// 		s.app.Use(middleware)
// 	}
// }
