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
}

func NewDefaultRepository(db *gorm.DB, es *elasticsearch.Client) DefaultRepositoryInterface {
	return &defaultRepository{
		esClient: es,
		db:       db,
	}
}

func (repo *defaultRepository) GetWelcome(ctx context.Context) modelDomain.Welcome {
	tp := openTelemetri.GetTraceProvider(ctx)

	ctx, span := tp.Tracer.Start(ctx, "defaultRepository.GetWelcome")
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
