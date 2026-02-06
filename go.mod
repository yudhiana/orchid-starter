module orchid-starter

go 1.24.5

// ==============================
// Core Web Framework
// ==============================
require github.com/go-chi/chi/v5 v5.2.4

// ==============================
// GraphQL
// ==============================
require (
	github.com/99designs/gqlgen v0.17.85
	github.com/vektah/gqlparser/v2 v2.5.31
)

// ==============================
// Database & ORM
// ==============================
require (
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.30.1
)

// ==============================
// Messaging
// ==============================
require github.com/rabbitmq/amqp091-go v1.10.0

// ==============================
// Elasticsearch
// ==============================
require (
	github.com/elastic/elastic-transport-go/v8 v8.7.0
	github.com/elastic/go-elasticsearch/v9 v9.1.0
)

// ==============================
// Observability & Monitoring
// ==============================
require (
	github.com/getsentry/sentry-go v0.35.0
	github.com/prometheus/client_golang v1.23.2
)

// ==============================
// Configuration
// ==============================
require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/joho/godotenv v1.5.1
)

// ==============================
// HTTP Client
// ==============================
require github.com/go-resty/resty/v2 v2.16.5

// ==============================
// CLI
// ==============================
require github.com/urfave/cli v1.22.17

// ==============================
// Template Engine
// ==============================
require github.com/tyler-sommer/stick v1.0.6

// ==============================
// Internal Modules
// ==============================
require (
	github.com/yudhiana/bunker v1.0.0
	github.com/yudhiana/logos v1.1.0
)

// ==============================
// mail service
// ==============================
require gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df

// ==============================
// Testing
// ==============================
require github.com/DATA-DOG/go-sqlmock v1.5.2

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-yaml v1.19.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/urfave/cli/v3 v3.6.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)

tool github.com/99designs/gqlgen
