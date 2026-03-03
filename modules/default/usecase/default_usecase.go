package usecase

import (
	"context"

	"orchid-starter/internal/clients"
	"orchid-starter/modules/default/delivery/event/publisher"
	"orchid-starter/modules/default/repository"
	modelUsecase "orchid-starter/modules/default/usecase/models"

	"gorm.io/gorm"
)

type defaultUsecase struct {
	db         *gorm.DB // use for transaction db .. NOTE : don't use for query!
	repository repository.DefaultRepositoryInterface
	client     *clients.Client
	publishing *publisher.EventPublisher
}

func NewDefaultUsecase(db *gorm.DB, r repository.DefaultRepositoryInterface, client *clients.Client, pub *publisher.EventPublisher) DefaultUsecaseInterface {
	return &defaultUsecase{
		db:         db,
		repository: r,
		client:     client,
		publishing: pub,
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

	/*
		Example use for event publisher!
			uc.publishing.PublishDefaultCreated(context.Background(), "orchid-exchange", "orchid.default.created", rabbitmq.Fanout, rabbitmq.Publishing{
				ContentType: "application/json",
				Type:        "orchid-default-created",
				AppId:       os.Getenv("APP_NAME"),
				Headers: map[string]any{
					"request-id": common.GetRequestIDFromContext(ctx),
				},
				MessageId:    uuid.NewString(),
				Timestamp:    time.Now().UTC(),
				DeliveryMode: rabbitmq.Persistent,
				Body:         []byte(`{"message": "Default created event published"}`),
			})

	*/
	return modelUsecase.GetWelcome{
		Message: uc.repository.GetWelcome(ctx).Message,
	}
}
