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
	"${APP_NAME}/internal/clients"
	"${APP_NAME}/modules/${MODULE_SLUG}/repository"

	"github.com/mataharibiz/ward/logging"
	"gorm.io/gorm"
)

type ${lowerCamelModule}Usecase struct {
	db         *gorm.DB // use for transaction db .. NOTE : don't use for query!
	repository repository.${CamelModule}RepositoryInterface
	client     *clients.Client
	log        *logging.LogEntry
}

func New${CamelModule}Usecase(db *gorm.DB, r repository.${CamelModule}RepositoryInterface, client *clients.Client) ${CamelModule}UsecaseInterface {
	return &${lowerCamelModule}Usecase{
		db:         db,
		repository: r,
		client:     client,
		log:        logging.NewLogger(),
	}
}
EOF
      ;;

    "repository/${FILE_NAME}_repository.go")
      cat <<EOF
import (
    "github.com/mataharibiz/ward/logging"
    "gorm.io/gorm"
)

type ${lowerCamelModule}Repository struct {
	db  *gorm.DB
	log *logging.LogEntry
}

func New${CamelModule}Repository(db *gorm.DB) ${CamelModule}RepositoryInterface {
	return &${lowerCamelModule}Repository{
		db:  db,
		log: logging.NewLogger(),
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
  esac
}


# Function to determine package name from file path
get_package_name() {
  local file_path="$1"
  local dir_name=$(dirname "$file_path")

  case "$dir_name" in
    "delivery/api/rest/v2") echo "v2" ;;
    "delivery/api/rest") echo "rest" ;;
    "delivery/event") echo "event" ;;
    "repository") echo "repository" ;;
    "usecase") echo "usecase" ;;
    "model/db") echo "modelDB" ;;
    "model/request") echo "modelRequest" ;;
    "model/response") echo "modelResponse" ;;
    "model/repository") echo "modelRepository" ;;
    "model/usecase") echo "modelUsecase" ;;
    *) echo "main" ;;
  esac
}

# Struktur paths relatif dari MODULES_ROOT
paths=(
  "delivery/api/rest/handler.go"
  "delivery/api/rest/v2/handler_v2.go"
  "delivery/event/event.go"
  "repository/${FILE_NAME}_repository.go"
  "repository/repository_interface.go"
  "usecase/${FILE_NAME}_usecase.go"
  "usecase/usecase_interface.go"
  "model/db/${FILE_NAME}.go"
  "model/request/request_${FILE_NAME}.go"
  "model/response/response_${FILE_NAME}.go"
  "model/repository/${FILE_NAME}_repo_input.go"
  "model/usecase/${FILE_NAME}_usecase_input.go"
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



