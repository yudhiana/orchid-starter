package usecase

import (
	"context"
	modelUsecase "orchid-starter/modules/default/usecase/models"
)

type DefaultUsecaseInterface interface {
	GetWelcome(ctx context.Context) modelUsecase.GetWelcome
}
