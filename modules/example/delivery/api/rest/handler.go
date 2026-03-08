package api

import (
	"orchid-starter/internal/bootstrap/container"
	httpUtil "orchid-starter/internal/bootstrap/server/restful_server/http_util"
	"orchid-starter/modules/example/repository"
	"orchid-starter/modules/example/usecase"

	v1 "orchid-starter/modules/example/delivery/api/rest/v1"
	"orchid-starter/modules/example/delivery/event/publisher"

	"github.com/go-chi/chi/v5"
)

func NewExampleHandler(app chi.Router, di *container.DirectInjection) {
	tracer := di.GetTracer()

	exampleRepository := repository.NewExampleRepository(di.GetMySQL(), di.GetElasticsearch(), tracer)

	// Get the comprehensive client for all API operations
	client := di.GetClient()

	// define event publisher
	pub := publisher.NewEventPublisher(di.GetPublisher())

	// Initialize usecase with client access
	exampleUseCase := usecase.NewExampleUsecase(di.GetMySQL(), exampleRepository, client, pub, tracer)
	exampleV1 := v1.NewExampleHandler(exampleUseCase, tracer)

	app.Get("/", exampleV1.Welcome)
	app.Get("/error-check", exampleV1.ErrorResponse)
	app.Get("/health-check", httpUtil.HealthCheckHandler)
}
