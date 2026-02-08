package repository

import (
	"context"
	modelDomain "orchid-starter/modules/default/domain/models"

	"gorm.io/gorm"
)

type DefaultRepositoryInterface interface {
	GetWelcome(ctx context.Context) modelDomain.Welcome
	WithTx(tx *gorm.DB) DefaultRepositoryInterface
}
