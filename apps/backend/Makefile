# cal-lowrisk backend — Makefile
# Schema is owned by golang-migrate, per-module folders, run in dependency order.
# See docs/02-system/architecture.md §5.1 and the add-module skill.
#
# Migrations use the DIRECT (non-pooled) connection string (deployment.md §3):
#   make migrate-up   DB_URL="postgres://...direct...?sslmode=require"
#   make migrate-down DB_URL="postgres://...direct...?sslmode=require"

MIGRATE ?= migrate
DB_URL  ?= $(DATABASE_URL_DIRECT)

# Module order matters because of FKs. Add new modules in the correct position.
MODULES_UP   = user food exercise nutrition workout
MODULES_DOWN = workout nutrition exercise food user

# Each module owns a PostgreSQL schema (= module name) and keeps its own migration
# tracking table INSIDE that schema: "<module>"."schema_migrations". The schema name is
# always double-quoted ("user" is a reserved word). golang-migrate won't CREATE SCHEMA
# for its tracking table, so we ensure the schema exists first (psql), then point migrate
# at the schema-qualified, quoted tracking table. See product.md + add-module §3.
PSQL ?= psql

.PHONY: run tidy build test test-unit test-e2e migrate-up migrate-down migrate-status migrate-create

run:
	go run ./cmd/api

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

## test: run everything (E2E auto-skips unless TEST_DATABASE_URL is set)
test:
	go test ./...

## test-unit: fast unit tests only (no DB)
test-unit:
	go test ./tests/unit/...

## test-e2e: E2E against a disposable DB (set TEST_DATABASE_URL to a throwaway Neon branch)
test-e2e:
	@test -n "$(TEST_DATABASE_URL)" || (echo "set TEST_DATABASE_URL to a disposable DB (see tests/README.md)"; exit 1)
	go test ./tests/e2e/... -v

## migrate-up: apply all module migrations in dependency order
migrate-up:
	@for m in $(MODULES_UP); do \
		echo ">> migrate up: $$m"; \
		$(PSQL) "$(DB_URL)" -v ON_ERROR_STOP=1 -c "CREATE SCHEMA IF NOT EXISTS \"$$m\";" || exit 1; \
		$(MIGRATE) -path migrations/$$m -database "$(DB_URL)" -x-migrations-table "\"$$m\".\"schema_migrations\"" -x-migrations-table-quoted=1 up || exit 1; \
	done

## migrate-down: roll back one step per module, reverse order
migrate-down:
	@for m in $(MODULES_DOWN); do \
		echo ">> migrate down 1: $$m"; \
		$(MIGRATE) -path migrations/$$m -database "$(DB_URL)" -x-migrations-table "\"$$m\".\"schema_migrations\"" -x-migrations-table-quoted=1 down 1 || exit 1; \
	done

## migrate-status: show version per module (usage: make migrate-status)
migrate-status:
	@for m in $(MODULES_UP); do \
		echo ">> $$m:"; \
		$(MIGRATE) -path migrations/$$m -database "$(DB_URL)" -x-migrations-table "\"$$m\".\"schema_migrations\"" -x-migrations-table-quoted=1 version || true; \
	done

## migrate-create: scaffold a paired up/down file (usage: make migrate-create M=user NAME=create_users)
migrate-create:
	@test -n "$(M)" || (echo "set M=<module>"; exit 1)
	@test -n "$(NAME)" || (echo "set NAME=<description>"; exit 1)
	$(MIGRATE) create -ext sql -dir migrations/$(M) -seq $(NAME)
