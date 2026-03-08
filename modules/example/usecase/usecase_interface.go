package usecase

import (
	"context"
	modelUsecase "orchid-starter/modules/example/usecase/models"
)

type ExampleUsecaseInterface interface {
	GetWelcome(ctx context.Context) modelUsecase.GetWelcome
}
