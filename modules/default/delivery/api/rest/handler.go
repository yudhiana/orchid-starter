package api

import (
	"orchid-starter/http"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/modules/default/repository"
	"orchid-starter/modules/default/usecase"

	v2 "orchid-starter/modules/default/delivery/api/rest/v2"
	"orchid-starter/modules/default/delivery/event/publisher"

	"github.com/go-chi/chi/v5"
)

func NewDefaultHandler(app chi.Router, di *bootstrap.DirectInjection) {
	tracer := di.GetTracer()

	defaultRepository := repository.NewDefaultRepository(di.GetMySQL(), di.GetElasticsearch(), tracer)

	// Get the comprehensive client for all API operations
	client := di.GetClient()

	// define event publisher
	pub := publisher.NewEventPublisher(di.GetPublisher())

	// Initialize usecase with client access
	defaultUseCase := usecase.NewDefaultUsecase(di.GetMySQL(), defaultRepository, client, pub, tracer)
	defaultV2 := v2.NewDefaultHandler(defaultUseCase, tracer)

	app.Get("/", defaultV2.Welcome)
	app.Get("/error-check", defaultV2.ErrorResponse)
	app.Get("/health-check", http.HealthCheckHandler)
}
