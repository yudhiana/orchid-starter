package api

import (
	"orchid-starter/http"
	"orchid-starter/internal/bootstrap"
	"orchid-starter/modules/default/repository"
	"orchid-starter/modules/default/usecase"

	v2 "orchid-starter/modules/default/delivery/api/rest/v2"

	"github.com/go-chi/chi/v5"
)

func NewDefaultHandler(app chi.Router, di *bootstrap.DirectInjection) {

	defaultRepository := repository.NewDefaultRepository(di.GetMySQL(), di.GetElasticsearch())

	// Get the comprehensive client for all API operations
	client := di.GetClient()

	// Initialize usecase with client access
	defaultUseCase := usecase.NewDefaultUsecase(di.GetMySQL(), defaultRepository, client)
	defaultV2 := v2.NewDefaultHandler(defaultUseCase)

	/*
		// this example is for versioning
		v2 := versioning.NewGroup(app, ">=2.0.0")
		v2.Get("/", defaultV2.Welcome)
	*/

	// app.Get("/metrics", iris.FromStd(promHttp.Handler()))
	app.Get("/", defaultV2.Welcome)
	app.Get("/health-check", http.HealthCheckHandler)
	// app.OnErrorCode(iris.StatusNotFound, http.NotFoundHandler)
}
