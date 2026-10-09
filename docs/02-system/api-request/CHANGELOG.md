# Postman Collection — Changelog

The importable collection is versioned with [Semantic Versioning](https://semver.org/).
The **latest** file is always the highest `vX.Y.Z` in this folder (see the table).

| Version | File | Date | Summary |
|---|---|---|---|
| **1.0.0** | `cal-lowrisk.v1.0.0.postman_collection.json` | 2026-10-10 | Initial release. `user-service / 01 Onboarding`: Register, Login, Create Profile. Auto-saves `{{token}}`/`{{userId}}` on login; per-request assertions mirror the spec. |

## Versioning rules
- **MAJOR** (`2.0.0`) — breaking change to existing requests (renamed/removed endpoint, changed required field, incompatible body shape).
- **MINOR** (`1.1.0`) — backward-compatible additions (new module folder, new request, new optional field).
- **PATCH** (`1.0.1`) — fixes that don't change the contract (example values, descriptions, test-script tweaks).

## Conventions
- One collection file per version; the filename carries the version so QA can see the latest at a glance.
- `info.version` inside the JSON matches the filename.
- Keep in sync with `docs/02-system/api-specs/` — the specs are the source of truth; the collection follows them.
