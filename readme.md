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
├── cmd/                # Application entrypoints (API, CLI, GraphQL, background tasks)
├── config/             # Application configuration and config models
├── constants/          # Global constants (API, context keys, timezone, etc.)
├── docker/             # Docker build configuration
├── gql/                # GraphQL schemas, resolvers, and generated code
├── http/               # HTTP handlers (health check, common responses)
├── infrastructure/     # External service integrations (DB, message broker, search)
├── internal/           # Application bootstrap, server setup, clients, and shared utilities
├── middleware/         # HTTP / GraphQL middleware
├── modules/            # Feature-based domain modules (usecase, repository, delivery)
├── observability/      # Monitoring, metrics, and error tracking
├── pkg/                # Reusable packages shared across modules
├── scripts/            # Helper scripts for code generation and tooling
└── security/           # Security-related utilities (hashing, etc.)

```


## 🏗 High-Level Architecture
flowchart TD
    subgraph Entrypoints["Entrypoints (cmd/)"]
        API["API Server<br/>cmd/api"]
        GQL["GraphQL Server<br/>cmd/gql"]
        CLI["CLI<br/>cmd/cli"]
        TASK["Background Tasks<br/>cmd/task"]
    end

    subgraph Bootstrap["Application Bootstrap (internal/bootstrap)"]
        APP["App Initialization"]
        DI["Dependency Injection"]
        SERVER["HTTP / GQL Server Setup"]
    end

    subgraph Middleware["Middleware"]
        MW["Auth / Context / Logging"]
    end

    subgraph Modules["Business Logic (modules/)"]
        USECASE["Usecase<br/>Business Rules"]
        DELIVERY["Delivery<br/>HTTP / GQL Handlers"]
        REPO["Repository<br/>Interfaces"]
    end

    subgraph Infrastructure["Infrastructure"]
        DB["MySQL"]
        MQ["RabbitMQ"]
        ES["Elasticsearch"]
    end

    subgraph Shared["Shared & Cross-cutting"]
        COMMON["internal/common"]
        CLIENTS["internal/clients"]
        PKG["pkg"]
        OBS["observability"]
        SEC["security"]
    end

    API --> Bootstrap
    GQL --> Bootstrap
    CLI --> Bootstrap
    TASK --> Bootstrap

    Bootstrap --> MW
    MW --> DELIVERY
    DELIVERY --> USECASE
    USECASE --> REPO

    REPO --> DB
    REPO --> MQ
    REPO --> ES

    USECASE --> COMMON
    DELIVERY --> COMMON
    USECASE --> CLIENTS

    Bootstrap --> OBS
    Bootstrap --> SEC
