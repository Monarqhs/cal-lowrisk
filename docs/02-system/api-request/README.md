# API Request Collection (Postman) — Cal-Lowrisk

Importable Postman files so QA can exercise the API without any manual setup. These mirror
the contract in [`../api-specs/`](../api-specs/) (the specs are the source of truth).

## Latest version

**v1.0.0** — see [`CHANGELOG.md`](./CHANGELOG.md) for the full history and the current
latest file. The latest collection is always the highest `vX.Y.Z` filename below.

```
api-request/
├── cal-lowrisk.v1.0.0.postman_collection.json   # <- import this (latest)
├── environments/
│   ├── cal-lowrisk-local.postman_environment.json
│   └── cal-lowrisk-uat.postman_environment.json
├── CHANGELOG.md
└── README.md
```

## How to import (Postman)

1. **Import the collection:** Postman → **Import** → drop
   `cal-lowrisk.v1.0.0.postman_collection.json`. It appears as **Cal-Lowrisk API (v1.0.0)**.
2. **Import an environment:** **Import** → drop one file from `environments/`.
3. **Select the environment** (top-right dropdown):
   - `cal-lowrisk-local` → hits a backend you run locally.
   - `cal-lowrisk-uat` → hits the deployed UAT backend (see note below).

## Running the onboarding flow

Open **user-service → 01 Onboarding** and send in order:

1. **Register** — creates a fresh account (the pre-request script generates a unique email
   per run, so you never hit the duplicate-email 409).
2. **Login** — on success it **auto-saves** `{{token}}` and `{{userId}}` into the active
   environment.
3. **Create Profile** — uses `Authorization: Bearer {{token}}` automatically; returns the
   computed daily calorie target.

Or use **Run collection** to execute all three with their assertions in one go.

## Environments & the backend URL

- **Local** (`cal-lowrisk-local`): `baseUrl = http://localhost:8080/api/v1`. Start the
  backend first (`apps/backend/`, see its README), then run requests.
- **UAT** (`cal-lowrisk-uat`): `baseUrl` is a **placeholder** (`REPLACE-ME-AFTER-DEPLOY`).
  The backend is **not hosted yet** — once it is deployed (Render → Neon `cal-lowrisk-uat`,
  see `../deployment.md`), set this environment's `baseUrl` to the deployed URL and QA can
  run everything remotely with no local setup.

## Keeping it current

When the API changes, bump the version per `CHANGELOG.md` and add a new
`cal-lowrisk.vX.Y.Z.postman_collection.json` file (don't silently overwrite), so QA always
knows which import file is the latest.
