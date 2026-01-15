package bootstrap

import (
	"fmt"

	"orchid-starter/config"
	"orchid-starter/internal/bootstrap/server/applications"

	"github.com/go-chi/chi/v5"
	"github.com/mataharibiz/ward/logging"
)

type Container struct {
	App *chi.Mux
	Cfg *config.LocalConfig
	DI  *DirectInjection
	Log *logging.LogEntry
}

// NewContainer creates a new application container with proper error handling
func NewContainer() (*Container, error) {
	logger := logging.NewLogger()
	logger.Info("Initializing application container...")

	// Load configuration
	cfg := config.GetLocalConfig()
	if cfg == nil {
		return nil, fmt.Errorf("failed to load configuration")
	}
	logger.Info("Configuration loaded successfully")

	// Initialize dependency injection
	di, err := NewDirectInjection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dependencies: %w", err)
	}
	logger.Info("Dependencies initialized successfully")

	// Initialize Iris application
	app := applications.GetChiApplication()
	if app == nil {
		return nil, fmt.Errorf("failed to initialize Iris application")
	}
	logger.Info("Iris application initialized successfully")

	container := &Container{
		App: app,
		Cfg: cfg,
		DI:  di,
		Log: logger,
	}

	logger.Info("Application container initialized successfully")
	return container, nil
}

// Close gracefully shuts down the container and its resources
func (c *Container) Close() error {
	logger := logging.NewLogger()
	logger.Info("Shutting down application container...")

	var errors []error

	// Close dependency injection resources
	if err := c.DI.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to close DI resources: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("container shutdown completed with errors: %v", errors)
	}

	logger.Info("Application container shut down successfully")
	return nil
}

// GetConfig returns the configuration
func (c *Container) GetConfig() *config.LocalConfig {
	return c.Cfg
}

// GetApp returns the Iris application
func (c *Container) GetApp() *chi.Mux {
	return c.App
}

// GetDI returns the dependency injection container
func (c *Container) GetDI() *DirectInjection {
	return c.DI
}

// Legacy support - will be deprecated
func StartContainer() *Container {
	logger := logging.NewLogger()
	logger.Warn("StartContainer() is deprecated, use NewContainer() instead")
	container, err := NewContainer()
	if err != nil {
		logger.Error("Failed to start container", "error", err)
		panic(err) // Maintain the fatal behavior but use panic instead of log.Fatalf
	}
	return container
}
