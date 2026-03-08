package repository

import (
	"context"
	modelDomain "orchid-starter/modules/example/domain/models"

	"gorm.io/gorm"
)

type ExampleRepositoryInterface interface {
	GetWelcome(ctx context.Context) modelDomain.Welcome
	WithTx(tx *gorm.DB) ExampleRepositoryInterface
}
