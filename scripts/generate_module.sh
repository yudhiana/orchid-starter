#!/bin/bash
set -a
source .env
set +a


to_snake_case() {
  echo "$1" | tr '-' '_'
}

to_camel_case() {
  local s
  s=$(to_snake_case "$1")
  echo "$s" | sed -r 's/(^|_)([a-z])/\U\2/g'
}

to_lower_camel_case() {
  local camel
  camel=$(to_camel_case "$1")
  echo "${camel,}"
}

MODULE_SLUG="${1:-default}"    
MODULE_KEY="$(to_snake_case "$MODULE_SLUG")"
MODULES_ROOT="./modules/${MODULE_SLUG}"
FILE_NAME="$(to_snake_case "$MODULE_SLUG")"

generate_file_content() {
  local path="$1"
  local package_name="$2"


  local CamelModule
  local lowerCamelModule

  CamelModule=$(to_camel_case "$MODULE_KEY")
  lowerCamelModule=$(to_lower_camel_case "$MODULE_KEY")

  echo "package $package_name"
  echo ""

  case "$path" in
    "repository/repository_interface.go")
      cat <<EOF
import (
	"gorm.io/gorm"
)

type ${CamelModule}RepositoryInterface interface {
	WithTx(tx *gorm.DB) ${CamelModule}RepositoryInterface
}
EOF
      ;;

    "usecase/usecase_interface.go")
      cat <<EOF
type ${CamelModule}UsecaseInterface interface {

}
EOF
      ;;

    "usecase/${FILE_NAME}_usecase.go")
      cat <<EOF
import (
	"${APP_NAME}/clients"
	"${APP_NAME}/modules/${MODULE_SLUG}/repository"

	"github.com/yudhiana/logos"
	"gorm.io/gorm"
)

type ${lowerCamelModule}Usecase struct {
	db         *gorm.DB // use for transaction db only!
	repository repository.${CamelModule}RepositoryInterface
	client     *clients.Client
	log        *logos.LogEntry
}

func New${CamelModule}Usecase(db *gorm.DB, r repository.${CamelModule}RepositoryInterface, client *clients.Client) ${CamelModule}UsecaseInterface {
	return &${lowerCamelModule}Usecase{
		db:         db,
		repository: r,
		client:     client,
		log:        logos.NewLogger(),
	}
}
EOF
      ;;

    "repository/${FILE_NAME}_repository.go")
      cat <<EOF
import (
    "github.com/yudhiana/logos"
    "gorm.io/gorm"
)

type ${lowerCamelModule}Repository struct {
	db  *gorm.DB
	log *logos.LogEntry
}

func New${CamelModule}Repository(db *gorm.DB) ${CamelModule}RepositoryInterface {
	return &${lowerCamelModule}Repository{
		db:  db,
		log: logos.NewLogger(),
	}
}

func (r *${lowerCamelModule}Repository) WithTx(tx *gorm.DB) ${CamelModule}RepositoryInterface {
	return &${lowerCamelModule}Repository{
		db:  tx,
		log: r.log,
	}
}
EOF
      ;;

    "delivery/event/publisher/publisher.go")
      cat <<EOF
import (
	"context"
	"${APP_NAME}/infrastructure/rabbitmq"
)

type EventPublisher struct {
	publishing rabbitmq.PublisherInterface
}

func NewEventPublisher(pub rabbitmq.PublisherInterface) *EventPublisher {
	return &EventPublisher{
		publishing: pub,
	}
}

// TODO: put here to add others event publisher
func (p *EventPublisher) Publish${CamelModule}Created(ctx context.Context, exchange, routingKey string, kind rabbitmq.Kind, msg rabbitmq.Publishing) error {
	return p.publishing.Publish(ctx, exchange, routingKey, kind, msg)
}

func (p *EventPublisher) PublishQueue${CamelModule}Created(ctx context.Context, queue string, msg rabbitmq.Publishing) error {
	return p.publishing.PublishQueue(ctx, queue, msg)
}

EOF
      ;;

    "delivery/event/subscriber/subscriber.go")
      cat <<EOF
import (
	"context"
	"fmt"
	"${APP_NAME}/infrastructure/rabbitmq"
	"${APP_NAME}/internal/bootstrap/container"

	bunker "github.com/yudhiana/bunker/errors"
	"github.com/yudhiana/logos"
)

// example event type constants
const (
	Event${CamelModule}Name = "example-${lowerCamelModule}-event-name"
)

type eventHandler struct {
	di  *container.DirectInjection
	log *logos.LogEntry
}

func New${CamelModule}EventHandler(di *container.DirectInjection) *eventHandler {
	return &eventHandler{
		di:  di,
		log: logos.NewLogger(),
	}
}

// Handle processes example init events based on event type
func (eh *eventHandler) Handle(ctx context.Context, event rabbitmq.Publishing) error {
	eh.log.Info("Processing example init event", "event_type", event.Type)
	switch event.Type {
	case Event${CamelModule}Name:
		return eh.${lowerCamelModule}Event(ctx, event)
	default:
		return bunker.New(bunker.StatusUnprocessableEntity).SetMessage(fmt.Sprintf("unknown event type: %s", event.Type))
	}
}

// GetEventTypes returns the list of event types this handler supports
func (eh *eventHandler) GetEventTypes() []string {
	return []string{
		Event${CamelModule}Name,
	}
}

func (eh *eventHandler) ${lowerCamelModule}Event(ctx context.Context, event rabbitmq.Publishing) error {
	eh.log.Info("event ${lowerCamelModule} successfully executed")
	return nil
}


EOF
      ;;
  esac
}


# Function to determine package name from file path
get_package_name() {
  local file_path="$1"
  local dir_name=$(dirname "$file_path")

  case "$dir_name" in
    "delivery/api/rest/v1") echo "v1" ;;
    "delivery/api/rest") echo "rest" ;;
    "delivery/api/gql") echo "gqlHandler" ;;
    "delivery/event/publisher") echo "publisher" ;;
    "delivery/event/subscriber") echo "subscriber" ;;
    "repository") echo "repository" ;;
    "usecase") echo "usecase" ;;
    "domain/models") echo "modelDomain" ;;
    "domain/errors") echo "errors" ;;
    "domain/events") echo "events" ;;
    "delivery/models/request") echo "modelRequest" ;;
    "delivery/models/response") echo "modelResponse" ;;
    "repository/models") echo "modelRepository" ;;
    "usecase/models") echo "modelUsecase" ;;
    *) echo "main" ;;
  esac
}

# Struktur paths relatif dari MODULES_ROOT
paths=(
  "domain/models/${FILE_NAME}.go"
  "domain/errors/${FILE_NAME}_errors.go"
  "domain/events/${FILE_NAME}_events.go"
  "delivery/models/request/request_${FILE_NAME}.go"
  "delivery/models/response/response_${FILE_NAME}.go"
  "delivery/api/rest/handler.go"
  "delivery/api/gql/gql_handler.go"
  "delivery/api/rest/v1/handler_v1.go"
  "delivery/event/publisher/publisher.go"
  "delivery/event/subscriber/subscriber.go"
  "repository/${FILE_NAME}_repository.go"
  "repository/repository_interface.go"
  "usecase/${FILE_NAME}_usecase.go"
  "usecase/usecase_interface.go"
  "repository/models/${FILE_NAME}_repo_input.go"
  "usecase/models/${FILE_NAME}_usecase_input.go"
)

# Loop buat direktori dan file
for path in "${paths[@]}"; do
  dir_path="${MODULES_ROOT}/$(dirname "$path")"
  file_path="${MODULES_ROOT}/${path}"
  package_name=$(get_package_name "$path")

  mkdir -p "$dir_path"
  if [ ! -f "$file_path" ]; then
    generate_file_content "$path" "$package_name" "$MODULE_SLUG" > "$file_path"
    echo "Created: $file_path (package: $package_name)"
  else
    echo "Already exists: $file_path"
  fi
done

echo "✅ Module '${MODULE_SLUG}' structure created at: $MODULES_ROOT"



