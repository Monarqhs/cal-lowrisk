# cal-lowrisk — Backend

Go **modular monolith** (Gin + GORM) — one deployable binary, split internally into
modules. This is the backend app of the cal-lowrisk monorepo (`apps/backend/`).

> Architecture, data model, and environments are documented in
> [`/docs/02-system/`](../../docs/02-system/). This README is the practical how-to.

## Layout

```
apps/backend/
├── cmd/api/main.go              # single entry point; wires modules + routes
├── internal/
│   ├── modules/                 # bounded contexts (user, food, exercise, nutrition, workout, summary)
│   │   └── <module>/            # controller → service → repository → model + dto
│   └── shared/                  # config, database, response envelope, base model
├── migrations/<module>/         # golang-migrate SQL (per-module, dependency-ordered)
├── Makefile                     # run / build / test / migrate targets
└── environment-variables.example
```

Rules (see `/docs/02-system/architecture.md` + the `add-module` skill):
- Dependencies flow inward: Controller → Service → Repository, via interfaces.
- Cross-module access only through another module's **service interface**.
- Schema is owned by **golang-migrate**, NOT GORM AutoMigrate. GORM = runtime queries only.
- Primary keys are native `uuid`, app-generated **UUIDv7** (`internal/shared/model.Base`).

## Prerequisites

- Go 1.25+
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI (for migrations)
- A PostgreSQL database (local via Docker for dev; **Neon** for UAT/PROD — deployment.md)

## Setup

```bash
cp environment-variables.example .env   # then fill in DATABASE_URL(S) + JWT_SECRET
make tidy
```

## Run

```bash
make run        # starts the API on :$PORT (default 8080)
# health check:
curl localhost:8080/healthz
```

## Migrations

Use the **direct** (non-pooled) connection string for migrations.

```bash
make migrate-up    DB_URL="$DATABASE_URL_DIRECT"   # apply all modules in FK order
make migrate-down  DB_URL="$DATABASE_URL_DIRECT"   # roll back one step per module
make migrate-status DB_URL="$DATABASE_URL_DIRECT"
make migrate-create M=user NAME=create_users       # scaffold paired up/down SQL
```

Module order (FKs): `user → food → exercise → nutrition → workout`.

## Build / test

```bash
make build      # -> bin/api
make test
```
