# API Specifications — cal-lowrisk

> **Contract-first.** These documents are the **source of truth** for the HTTP API. The
> backend services implement them, and the test suite (`apps/backend/tests/`) verifies
> them. When behaviour and spec disagree, the spec wins — change the spec first (via PR),
> then the code and tests. This keeps unit/E2E tests from drifting (see `apps/backend`
> `tests/` and the discussion in the project history).

## Layout

```
api-specs/
└── <module>-service/            # one folder per backend module (bounded context)
    ├── 01. <Use Case>.md        # one file per use case / flow (may span several endpoints)
    ├── 02. <Use Case>.md
    └── ...
```

- A file documents a **use case / flow**, not a single endpoint — e.g. "Onboarding User"
  covers register → login → create profile.
- Numbered (`01.`, `02.`, …) in the order a user/consumer naturally meets them.

## Shared conventions (apply to every spec unless a spec overrides)

### Base URL & versioning
- All endpoints are under `/api/v1`. Paths in specs are written relative to that prefix
  (a spec path `POST /auth/register` means `POST /api/v1/auth/register`).

### Response envelope
Every response uses one envelope (`internal/shared/response`):

```jsonc
// success
{ "success": true, "data": { /* payload */ } }

// error
{ "success": false, "error": { "code": "MACHINE_CODE", "message": "Human readable." } }
```

- `data` is present only on success; `error` only on failure.
- Clients branch on `success` and, for errors, on `error.code` (stable machine string) —
  never on `error.message` (human text, may change / localise later).

### Auth
- Protected endpoints require `Authorization: Bearer <JWT>`.
- The JWT is issued by the login endpoint. Its claims include at least the user `id`
  (subject) and `role` name. TTL per `JWT_TTL_HOURS` (default 24h).
- **Privacy boundary (hard rule):** a `user`-role token may only act on **its own** data.
  Admin endpoints require the `admin` role. Server enforces this; never trust the client.

### Standard error codes
| `error.code` | HTTP | Meaning |
|---|---|---|
| `VALIDATION_ERROR` | 400 | Request body/params failed validation. `message` describes the first/aggregated problem. |
| `UNAUTHENTICATED` | 401 | Missing/invalid/expired token, or bad credentials on login. |
| `FORBIDDEN` | 403 | Authenticated but not allowed (role / ownership). |
| `NOT_FOUND` | 404 | Resource does not exist (or is not visible to this caller). |
| `CONFLICT` | 409 | Violates a uniqueness/state rule (e.g. email already registered). |
| `INTERNAL` | 500 | Unexpected server error. |

Specs may add module-specific codes; list them in the spec's own error table.

### Conventions for field shapes
- Timestamps: ISO-8601 UTC strings (`2026-10-09T20:16:22Z`).
- IDs: UUID strings (UUIDv7, see `erd.md` §2).
- Money/none; nutrition & body numbers are JSON numbers (backed by `NUMERIC` — see `erd.md`).

## Modules
| Spec folder | Backend module | Status |
|---|---|---|
| `user-service/` | `user` | 🔧 in progress (Onboarding) |
| `food-service/` | `food` | ⏳ planned |
| `exercise-service/` | `exercise` | ⏳ planned |
| `nutrition-service/` | `nutrition` | ⏳ planned |
| `workout-service/` | `workout` | ⏳ planned |
| `summary-service/` | `summary` | ⏳ planned |

> "service" in the folder name reflects the **bounded context** (and the microservice
> extraction path), even though the backend is a single modular monolith today.
