# Orchid Starter

A starter template for building modern, scalable backend services using **Go** with a clean and modular architecture.

---

## 🧠 Overview

**Orchid Starter** is a backend starter project written in **Go**, designed to help developers bootstrap new services quickly while following best practices in project structure, configuration management, and extensibility.

This project is suitable for:
- Backend services
- GraphQL-based APIs
- Microservices
- Production-ready Go applications

---

## 🚀 Key Features

- 🧩 **Modular Architecture**  
  Clear separation of concerns using domain-based modules.
- 🌐 **GraphQL Ready**  
  Includes GraphQL schema and `gqlgen` configuration.
- 🐳 **Docker Support**  
  Ready for containerized development and deployment.
- ⚙️ **Environment Configuration**  
  Uses environment variables with `.env.example` as reference.
- 🔐 **Security & Middleware Layer**  
  Structured place for authentication, authorization, and middleware.
- 📊 **Observability Ready**  
  Dedicated folder for logging, metrics, and tracing.

---

## 🛠️ Getting Started

### 📥 Clone the Repository

```bash
git clone https://github.com/yudhiana/orchid-starter.git
cd orchid-starter
```


## ⚙️ Environment Variables

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Adjust the values based on your local setup.
```bash
APP_ENV=local
APP_PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/dbname
```

## 📦 Install Dependencies
```bash
go mod download
```


## 🧠 Generate Code

If using gqlgen, generate GraphQL boilerplate code:

```bash
./scripts/generate_gql.sh
```


If wants to generate module boilerplate code:
```bash
./scripts/generate_module.sh $module_name
```


## 📁 Project Structure
```text
.
├── cmd/                # Application entrypoint
├── config/             # Application configuration
├── constants/          # Global constants
├── docker/             # Docker configuration
├── gql/                # GraphQL schema definitions
├── http/               # HTTP server and handlers
├── infrastructure/     # External integrations (DB, cache, etc)
├── internal/           # Core application logic
├── middleware/         # HTTP / auth middleware
├── modules/            # Feature-based domain modules
├── observability/      # Logging, metrics, tracing
├── pkg/                # Shared / reusable packages
├── scripts/            # Utility scripts
├── security/           # Security-related logic
├── .env.example        # Environment variable template
├── gqlgen.yml          # GraphQL generator config
└── go.mod              # Go module definition
```