---
name: add-module
description: Scaffold a new module in the cal-lowrisk modular monolith (Go) and its SQL migrations. Use when adding a new bounded context/module (e.g. user, food, nutrition, exercise, workout, summary), creating its layered code, or writing/ordering golang-migrate migrations. Covers the module-first layout, Controller→Service→Repository→Model+DTO layering, per-module migration folders, dependency ordering, and seed migrations.
metadata:
  project: cal-lowrisk
  version: "0.1"
---

# Add a Module — cal-lowrisk

How to add a new module to the modular monolith and wire its migrations. Follow the
project conventions (see steering `product.md`): module-first layout, layered code,
GORM for runtime queries only, **golang-migrate** for schema with **per-module folders**.

## 1. Module code layout (module-first)

Each module is a bounded context under `internal/modules/<module>/`, keeping the same
layers. Dependencies flow inward (Controller → Service → Repository). Use interfaces for
dependency injection so layers are testable.

```text
internal/modules/<module>/
├── controller.go      # HTTP handlers (Gin). Parse/validate request, call service, format response. No business logic.
├── service.go         # Business logic. Interface + impl. Returns DTOs, not models. No HTTP, no direct DB.
├── repository.go      # Data access via GORM. Interface + impl. Only DB ops.
├── model.go           # GORM structs = DB entities. No business logic, no HTTP.
├── dto/
│   ├── request.go     # Request DTOs with validation tags.
│   └── response.go    # Response DTOs. No DB tags.
└── module.go          # Wiring: constructs repo→service→controller, registers routes.
```

Rules (match the workplace Go standard):
- Controller: HTTP only; use the shared standard response (`success`/`error`).
- Service: all business rules; depends on repository **interfaces**; returns DTOs.
- Repository: GORM queries only; exposes an interface.
- Model: struct + GORM tags only. Prefer a shared base (id, created_at, updated_at,
  deleted_at for soft delete) consistent across modules.
- Cross-module calls go through the other module's **service interface**, never by
  reaching into its repository or tables directly. This keeps boundaries clean and
  preserves the microservice "extraction path".

## 2. Register the module

In `cmd/api/main.go` (single entry point for the monolith), construct the module and
register its routes. Order construction so dependencies (e.g. `user`) are available to
modules that need them.

## 3. Migrations (golang-migrate, per-module folders)

Schema is owned by SQL migrations — **never** GORM AutoMigrate. One folder per module.

```text
migrations/
└── <module>/
    ├── 000001_create_<table>.up.sql
    ├── 000001_create_<table>.down.sql
    └── ...
```

- Write both `.up.sql` and `.down.sql` (every migration must be reversible).
- Each module uses its **own tracking table** (`-x-migrations-table schema_migrations_<module>`),
  so modules stay isolated.
- **Dependency order matters** because of foreign keys. The canonical order is:
  `user` (roles → users) → `food` → `exercise` → `nutrition` → `workout` → `summary`.
  A module's migrations must run **after** any module it has FKs into.

### Makefile execution (explicit order)
Add the new module to `migrate-up`/`migrate-down` in the correct position:

```makefile
migrate-up:
	migrate -path migrations/user      -database "$(DB_URL)" -x-migrations-table schema_migrations_user up
	migrate -path migrations/food      -database "$(DB_URL)" -x-migrations-table schema_migrations_food up
	migrate -path migrations/exercise  -database "$(DB_URL)" -x-migrations-table schema_migrations_exercise up
	migrate -path migrations/nutrition -database "$(DB_URL)" -x-migrations-table schema_migrations_nutrition up
	migrate -path migrations/workout   -database "$(DB_URL)" -x-migrations-table schema_migrations_workout up
# migrate-down should reverse the order (workout → ... → user), stepping down 1 per folder.
```

> When adding a module with FKs, insert it **after** its dependencies in `migrate-up`
> and **before** them in `migrate-down`.

## 4. Seed migrations (master/reference data)

Master data (roles, food/exercise catalogs) is seeded via dedicated migrations in the
owning module's folder, numbered after its structural migrations, e.g.
`migrations/user/000010_seed_roles.up.sql` inserting the `user` and `admin` roles.
Keep seed data idempotent where practical (e.g. `ON CONFLICT DO NOTHING`).

## 5. Checklist when adding a module

- [ ] Create `internal/modules/<module>/` with controller, service, repository, model, dto, module.go
- [ ] Service depends on repository **interface**; cross-module access via service interfaces only
- [ ] Register construction + routes in `cmd/api/main.go`
- [ ] Create `migrations/<module>/` with paired up/down SQL
- [ ] Add the module to `migrate-up` (after deps) and `migrate-down` (before deps) in the Makefile
- [ ] Add seed migrations if the module owns master/reference data
- [ ] Respect the privacy boundary: admin features must not expose users' private logs
- [ ] Document the module's tables in `docs/02-system/erd.md`
