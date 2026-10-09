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
  deleted_at for soft delete) consistent across modules. Each model declares its
  schema-qualified table via `TableName()` using the shared helper, e.g.
  `func (User) TableName() string { return model.Qualify("user", "users") }` — this is how
  the per-module schema (§3) is enforced at the ORM layer (no reliance on `search_path`).
- Cross-module calls go through the other module's **service interface**, never by
  reaching into its repository or tables directly. This keeps boundaries clean and
  preserves the microservice "extraction path".

## 2. Register the module

In `cmd/api/main.go` (single entry point for the monolith), construct the module and
register its routes. Order construction so dependencies (e.g. `user`) are available to
modules that need them.

## 3. Migrations (golang-migrate, per-module folders + per-module schema)

Schema is owned by SQL migrations — **never** GORM AutoMigrate. One folder per module.
Each module also owns a **dedicated PostgreSQL schema named after the module** (DB-level
bounded context): tables live in `"<module>"."<table>"`, not in `public`.

```text
migrations/
└── <module>/
    ├── 000001_create_<table>.up.sql   # starts with CREATE SCHEMA IF NOT EXISTS "<module>";
    ├── 000001_create_<table>.down.sql # ends by DROP-ing the (now empty) schema
    └── ...
```

- Write both `.up.sql` and `.down.sql` (every migration must be reversible).
- **Schema per module:** the module name IS the schema name. Qualify every table as
  `"<module>".<table>`. The module name `user` is a SQL **reserved word** → the `"user"`
  schema must **always be double-quoted**; other module schemas don't need quoting but
  stay consistent. Migration `000001` begins with `CREATE SCHEMA IF NOT EXISTS "<module>";`
  and its `.down.sql` ends with `DROP SCHEMA IF EXISTS "<module>" RESTRICT;` after dropping
  the table.
- Each module keeps its **own tracking table inside its own schema**:
  `-x-migrations-table '"<module>"."schema_migrations"' -x-migrations-table-quoted=1`,
  so modules stay fully isolated. (golang-migrate won't `CREATE SCHEMA` for its tracking
  table, so the Makefile `CREATE SCHEMA IF NOT EXISTS` first — see below.)
- **Cross-schema FKs are allowed** (single DB today), but cross-module *data access* must
  go through the other module's **service interface** (see §1), never a direct cross-schema
  JOIN — this preserves the microservice extraction path.
- **Dependency order matters** because of foreign keys. The canonical order is:
  `user` (roles → users) → `food` → `exercise` → `nutrition` → `workout` → `summary`.
  A module's migrations must run **after** any module it has FKs into.

### Makefile execution (explicit order)
The Makefile loops modules in dependency order. For each module it (1) ensures the schema
exists, then (2) runs migrate against the schema-qualified, quoted tracking table:

```makefile
# MODULES_UP   = user food exercise nutrition workout
# For each $$m:
psql "$(DB_URL)" -v ON_ERROR_STOP=1 -c "CREATE SCHEMA IF NOT EXISTS \"$$m\";"
migrate -path migrations/$$m -database "$(DB_URL)" \
        -x-migrations-table "\"$$m\".\"schema_migrations\"" -x-migrations-table-quoted=1 up
# migrate-down reverses the order (workout → ... → user), stepping down 1 per folder.
```

> When adding a module with FKs, insert it **after** its dependencies in `MODULES_UP`
> and **before** them in `MODULES_DOWN`.
>
> **Sandbox note:** the golang-migrate CLI may fail to install (Go toolchain clash). In
> the Kiro sandbox, run the SQL via the Neon Power (`run_sql`, one statement at a time)
> and keep `"<module>"."schema_migrations"` in sync manually (set `version`, `dirty=false`).

## 4. Seed migrations (master/reference data)

Master data (roles, food/exercise catalogs) is seeded via dedicated migrations in the
owning module's folder, numbered after its structural migrations, e.g.
`migrations/user/000004_seed_roles.up.sql` inserting the `user` and `admin` roles into
`"user".role`. Seed into the module's **schema-qualified** table, and keep seed data
idempotent where practical (e.g. `ON CONFLICT DO NOTHING`).

## 5. Checklist when adding a module

- [ ] Create `internal/modules/<module>/` with controller, service, repository, model, dto, module.go
- [ ] Service depends on repository **interface**; cross-module access via service interfaces only
- [ ] Register construction + routes in `cmd/api/main.go`
- [ ] Create `migrations/<module>/` with paired up/down SQL; `000001` does
      `CREATE SCHEMA IF NOT EXISTS "<module>"`, all tables are `"<module>".<table>`, and
      the down of `000001` drops the schema
- [ ] Each model's `TableName()` returns `model.Qualify("<module>", "<table>")`
- [ ] Add the module to `MODULES_UP` (after deps) and `MODULES_DOWN` (before deps) in the Makefile
- [ ] Add seed migrations if the module owns master/reference data
- [ ] Respect the privacy boundary: admin features must not expose users' private logs
- [ ] Document the module's tables in `docs/02-system/erd.md`
