package repository

import (
	"context"
	modelDomain "orchid-starter/modules/default/domain/models"
	openTelemetri "orchid-starter/observability/open-telemetri"

	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type defaultRepository struct {
	esClient *elasticsearch.Client
	db       *gorm.DB
	otel     *openTelemetri.OTel
}

func NewDefaultRepository(db *gorm.DB, es *elasticsearch.Client, otel *openTelemetri.OTel) DefaultRepositoryInterface {
	return &defaultRepository{
		esClient: es,
		db:       db,
		otel:     otel,
	}
}

func (repo *defaultRepository) GetWelcome(ctx context.Context) modelDomain.Welcome {
	ctx, span := repo.otel.StartSpan(ctx, "repository", openTelemetri.GetFuncName())
	span.SetAttributes(openTelemetri.MakeTags(map[string]any{
		"repo":   "default repo",
		"method": "GetWelcome",
	})...)
	defer span.End()

	return modelDomain.Welcome{
		Message: "Welcome to orchid-starter...",
	}
}

func (repo *defaultRepository) WithTx(tx *gorm.DB) DefaultRepositoryInterface {
	return &defaultRepository{
		esClient: repo.esClient,
		db:       tx,
	}
}
