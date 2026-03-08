package repository

import (
	"context"
	modelDomain "orchid-starter/modules/example/domain/models"
	openTelemetry "orchid-starter/observability/open-telemetry"

	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type exampleRepository struct {
	esClient *elasticsearch.Client
	db       *gorm.DB
	otel     *openTelemetry.OTel
}

func NewExampleRepository(db *gorm.DB, es *elasticsearch.Client, otel *openTelemetry.OTel) ExampleRepositoryInterface {
	return &exampleRepository{
		esClient: es,
		db:       db,
		otel:     otel,
	}
}

func (repo *exampleRepository) GetWelcome(ctx context.Context) modelDomain.Welcome {
	ctx, span := repo.otel.StartSpan(ctx, "repository", openTelemetry.GetFuncName())
	span.SetAttributes(openTelemetry.MakeTags(map[string]any{
		"repo":   "example repo",
		"method": "GetWelcome",
	})...)
	defer span.End()

	return modelDomain.Welcome{
		Message: "Welcome to orchid-starter...",
	}
}

func (repo *exampleRepository) WithTx(tx *gorm.DB) ExampleRepositoryInterface {
	return &exampleRepository{
		esClient: repo.esClient,
		db:       tx,
	}
}
