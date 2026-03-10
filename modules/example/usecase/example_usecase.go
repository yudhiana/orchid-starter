package usecase

import (
	"context"

	"orchid-starter/clients"
	"orchid-starter/modules/example/delivery/event/publisher"
	"orchid-starter/modules/example/repository"
	modelUsecase "orchid-starter/modules/example/usecase/models"
	openTelemetry "orchid-starter/observability/open-telemetry"

	"gorm.io/gorm"
)

type exampleUsecase struct {
	db         *gorm.DB // use for transaction db .. NOTE : don't use for query!
	repository repository.ExampleRepositoryInterface
	client     *clients.Client
	publishing *publisher.EventPublisher
	otel       *openTelemetry.OTel
}

func NewExampleUsecase(db *gorm.DB, r repository.ExampleRepositoryInterface, client *clients.Client, pub *publisher.EventPublisher, otel *openTelemetry.OTel) ExampleUsecaseInterface {
	return &exampleUsecase{
		db:         db,
		repository: r,
		client:     client,
		publishing: pub,
		otel:       otel,
	}
}

func (uc *exampleUsecase) GetWelcome(ctx context.Context) (result modelUsecase.GetWelcome) {
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


			// Example use for event publisher!
			uc.publishing.PublishExampleCreated(context.Background(), "orchid-event", "orchid.example.created", rabbitmq.Fanout, rabbitmq.Publishing{
				ContentType: "application/json",
				Type:        "example-init-event-name",
				AppId:       os.Getenv("APP_NAME"),
				Headers: map[string]any{
					"request-id": common.GetRequestIDFromContext(ctx),
				},
				MessageId:    uuid.NewString(),
				Timestamp:    time.Now().UTC(),
				DeliveryMode: rabbitmq.Persistent,
				Body:         []byte(`{"message": "example created event published"}`),
			})

	*/

	ctx, span := uc.otel.StartSpan(ctx, "usecase", openTelemetry.GetFuncName())
	defer span.End()
	return modelUsecase.GetWelcome{
		Message: uc.repository.GetWelcome(ctx).Message,
	}
}
