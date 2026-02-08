package usecase

import (
	"context"

	"orchid-starter/internal/clients"
	"orchid-starter/modules/default/repository"
	modelUsecase "orchid-starter/modules/default/usecase/models"

	"gorm.io/gorm"
)

type defaultUsecase struct {
	db         *gorm.DB // use for transaction db .. NOTE : don't use for query!
	repository repository.DefaultRepositoryInterface
	client     *clients.Client
}

func NewDefaultUsecase(db *gorm.DB, r repository.DefaultRepositoryInterface, client *clients.Client) DefaultUsecaseInterface {
	return &defaultUsecase{
		db:         db,
		repository: r,
		client:     client,
	}
}

func (uc *defaultUsecase) GetWelcome(ctx context.Context) (result modelUsecase.GetWelcome) {
	/*
		    NOTE: This usecase is responsible for encapsulating business logic and does not handle data creation or persistence.
			Example use for transaction!

			Manual way
			---------------------------
			tx := uc.db.Begin()
			repo := uc.repository.WithTx(tx)

			repo.Welcome()
			tx.Commit()

			OR

			Traditional way
			---------------------------
			uc.db.Transaction(func(tx *gorm.DB) error {
				repo := uc.repository.WithTx(tx)
				result = repo.Welcome()
				return nil
			})
	*/

	return modelUsecase.GetWelcome{
		Message: uc.repository.GetWelcome(ctx).Message,
	}
}
