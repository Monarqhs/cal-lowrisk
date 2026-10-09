# Backend Tests — cal-lowrisk

Centralized test suite for the Go backend. Tests are **contract-first**: they verify the
"Test cases" sections of the API specs in `docs/02-system/api-specs/`, so tests can't drift
from the contract. Case IDs (e.g. `TC-ON-13`) map directly to the spec.

## Layout

```
tests/
├── unit/        # fast, no HTTP, no DB — pure logic (BMR/TDEE/target, password hashing)
├── e2e/         # drive the real HTTP router against a real PostgreSQL database
└── helpers/     # test harness: boots the router + DB, request/JSON helpers, migrations
```

## Running

### Unit tests (always, no setup)
```bash
go test ./tests/unit/...
```

### Everything (unit + E2E)
E2E tests connect to a **disposable PostgreSQL database** via `TEST_DATABASE_URL`. If that
variable is **unset**, E2E tests **skip automatically** (so `go test ./...` stays green on a
machine without a DB).

```bash
# unit only (E2E auto-skip)
go test ./...

# unit + E2E against a disposable DB
export TEST_DATABASE_URL="postgresql://.../neondb?sslmode=require"
go test ./...
```

## E2E database: use a disposable Neon branch

Never point E2E at the shared UAT data. Create a throwaway branch from `uat`, run against
it, then delete it. In this project we use the Neon Power:

1. `create_branch` (name e.g. `e2e-user-tests`, parent = the `uat`/production branch) — it
   copies the parent at HEAD, so the `"user"` schema is already present.
2. `get_connection_string` for that branch → export as `TEST_DATABASE_URL`.
3. `go test ./...`
4. `delete_branch` to clean up.

The harness also applies `migrations/user/*.up.sql` on startup (idempotent), so E2E works
even against an empty branch. Each test calls `TruncateAll` to isolate state while keeping
the seeded roles.

> **Sandbox note:** run the server/tests inside one shell invocation with the env exported;
> background processes don't persist across tool calls here. `golangci-lint` in the sandbox
> targets an older Go than the module, so rely on `go vet` + `gofmt`.
