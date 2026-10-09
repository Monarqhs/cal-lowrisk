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

.PHONY: run tidy build test migrate-up migrate-down migrate-status migrate-create

run:
	go run ./cmd/api

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

## migrate-up: apply all module migrations in dependency order
migrate-up:
	@for m in $(MODULES_UP); do \
		echo ">> migrate up: $$m"; \
		$(MIGRATE) -path migrations/$$m -database "$(DB_URL)" -x-migrations-table schema_migrations_$$m up || exit 1; \
	done

## migrate-down: roll back one step per module, reverse order
migrate-down:
	@for m in $(MODULES_DOWN); do \
		echo ">> migrate down 1: $$m"; \
		$(MIGRATE) -path migrations/$$m -database "$(DB_URL)" -x-migrations-table schema_migrations_$$m down 1 || exit 1; \
	done

## migrate-status: show version per module (usage: make migrate-status)
migrate-status:
	@for m in $(MODULES_UP); do \
		echo ">> $$m:"; \
		$(MIGRATE) -path migrations/$$m -database "$(DB_URL)" -x-migrations-table schema_migrations_$$m version || true; \
	done

## migrate-create: scaffold a paired up/down file (usage: make migrate-create M=user NAME=create_users)
migrate-create:
	@test -n "$(M)" || (echo "set M=<module>"; exit 1)
	@test -n "$(NAME)" || (echo "set NAME=<description>"; exit 1)
	$(MIGRATE) create -ext sql -dir migrations/$(M) -seq $(NAME)
