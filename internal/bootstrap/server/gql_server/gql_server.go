package gqlServer

import (
	"net/http"
	"orchid-starter/config"
	"orchid-starter/internal/bootstrap/container"

	"github.com/go-chi/chi/v5"
	"github.com/yudhiana/logos"
)

type Server struct {
	cfg *config.LocalConfig
	app *chi.Mux
}

func NewGQLServer(container *container.Container) *Server {
	srv := &Server{
		cfg: container.GetConfig(),
		// app: container.GetApp(),
	}

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

func (s *Server) setupRoutes(container *container.Container) {
	GQLRoutes(s.app, container)
}
