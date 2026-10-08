# Environments & Deployment — cal-lowrisk

| Field | Value |
|---|---|
| **Project** | cal-lowrisk |
| **Document** | Environments, configuration & deployment strategy |
| **Version** | 0.1 (Draft) |
| **Status** | Draft — pending review |
| **Author** | Product Owner + Kiro (SA) |
| **Phase** | System Analysis (companion to `architecture.md`) |
| **Last updated** | 2026-10-08 |

> **Process note:** Companion to [`architecture.md`](./architecture.md). Records *where*
> and *how* cal-lowrisk runs: the environments, how one codebase targets each, how config
> and secrets differ, and how migrations are promoted. Everything here stays within
> **free-tier** limits (BRD §7). It is a plan — nothing is provisioned yet.

---

## 1. Core principle — one codebase, many targets

**We do NOT duplicate code per environment.** There is **one** backend codebase, **one**
Next.js codebase, **one** Flutter codebase. What differs between environments is
**configuration (env vars / secrets) and the deploy target** — never the source.

```mermaid
flowchart LR
    subgraph Source["Source code (single)"]
        BE["Go backend"]
        WEB["Next.js admin"]
        APP["Flutter app"]
    end

    subgraph UAT["UAT environment"]
        BEU["backend @ UAT config"]
        DBU[("Neon project<br/>cal-lowrisk-uat")]
    end

    subgraph PROD["PROD environment"]
        BEP["backend @ PROD config"]
        DBP[("Neon project<br/>cal-lowrisk-prod")]
    end

    BE -->|".env.uat / UAT secrets"| BEU --> DBU
    BE -->|".env.prod / PROD secrets"| BEP --> DBP
    WEB -->|"API_URL=uat"| BEU
    WEB -->|"API_URL=prod"| BEP
    APP -->|"flavor uat"| BEU
    APP -->|"flavor prod"| BEP
```

---

## 2. Environments

Two long-lived environments:

| Env | Purpose | Audience |
|---|---|---|
| **UAT** | Integration/testing; validate migrations & features before release. | Us (test accounts, seed data). |
| **PROD** | The "live" environment. | Real usage. |

> Migrations and features always land in **UAT first**, are verified, then promoted to
> **PROD** (BRD-style BA→SA→UX→Dev discipline extended to releases).

---

## 3. Database — Neon (two separate projects)

**Decision:** use **two separate Neon projects**, not two branches of one project.

```text
Neon account
├── project: cal-lowrisk-uat   → database for UAT
└── project: cal-lowrisk-prod  → database for PROD
```

### Why two projects (not two branches)
Neon's **compute quota is scoped per project** (100 CU-hours/project/month on the Free
plan). Two projects therefore give us **two independent 100 CU-hour budgets** and real
"prod vs non-prod" isolation. With a single project + two branches, both environments
would **share one 100 CU-hour budget**, and — per Neon's Free-plan behavior — once the
quota is hit, the **non-primary branch (UAT) is suspended until the next month** while the
primary (PROD) keeps running. Separate projects avoid that coupling entirely.

> Branches still have a place later: within each project we can spin up a short-lived
> branch for a **per-PR preview**, then delete it. That's a feature use case, not a
> long-lived environment.

### Free-tier facts that shape our usage
Source: Neon Free plan docs (see §9). Content rephrased for licensing compliance.

| Fact | Value | Our handling |
|---|---|---|
| Compute | 100 CU-hours / project / month | Two projects = 2× budget; fine for low traffic. |
| Branches | up to 10 / project | Used for PR previews only. |
| Storage | ~0.5–1 GB / project | Ample for MVP. |
| **Auto-suspend** | after **5 min idle** (fixed on Free) | First request after idle = cold start (~hundreds of ms). Clients/health checks must tolerate it; set sane DB connect timeouts. |
| **Connection pooling** | PgBouncer; hostname contains **`-pooler`** | **We MUST use the pooled connection string** and keep the Go pool small. |

### Connection pooling (hard requirement)
Go + GORM opens a pool of real Postgres connections; Neon Free is connection-limited. So:
- Use the **pooled** (`-pooler`) connection string in `DATABASE_URL`.
- Keep the app pool small: e.g. `SetMaxOpenConns(5–10)`, `SetMaxIdleConns(2)`,
  `SetConnMaxLifetime(~5m)`. Tune later.
- For **migrations**, prefer the **direct (non-pooled)** connection string — some
  migration operations don't play well with transaction-mode pooling.

---

## 4. Configuration & secrets (per environment)

Same variables, different values per env. **Secrets are never committed** (BRD §7). Local
dev uses git-ignored `.env` files; CI/CD uses the platform's secret store.

| Variable | UAT | PROD | Notes |
|---|---|---|---|
| `DATABASE_URL` | Neon `cal-lowrisk-uat` **pooled** URL | Neon `cal-lowrisk-prod` **pooled** URL | App runtime (GORM). |
| `DATABASE_URL_DIRECT` | UAT **direct** URL | PROD **direct** URL | Migrations (golang-migrate). |
| `JWT_SECRET` | distinct value | distinct value | Never reuse across envs. |
| `APP_ENV` | `uat` | `prod` | For logging/behavior flags. |
| `REDIS_URL` | *(deferred)* | *(deferred)* | When the summary cache is enabled. |

**Client config:**
- **Next.js admin:** `API_URL` (or `NEXT_PUBLIC_API_URL`) points at the matching backend;
  Vercel maps *preview* → UAT backend, *production* → PROD backend.
- **Flutter:** build **flavors** `uat` / `prod`, each compiled with its API base URL.

**Local files (git-ignored):**
```text
.env            # local dev (your own throwaway Neon project or UAT)
.env.uat        # UAT values (NOT committed)
.env.prod       # PROD values (NOT committed)
environment-variables.example   # committed template with blank values
```

---

## 5. Migrations across environments (golang-migrate)

Same migration files (`migrations/<module>/*.{up,down}.sql`) run against each env — only
the target `DB_URL` changes. Schema is owned by **golang-migrate**, per-module folders,
Makefile-ordered (see `architecture.md` §5.1 and the `add-module` skill).

**Promotion flow (UAT → PROD):**
```mermaid
flowchart LR
    A["Write migration<br/>(up + down SQL)"] --> B["make migrate-up DB_URL=UAT_DIRECT"]
    B --> C["Verify in UAT<br/>(info / validate / smoke test)"]
    C --> D{"OK?"}
    D -->|"no"| A
    D -->|"yes"| E["make migrate-up DB_URL=PROD_DIRECT"]
    E --> F["Verify in PROD"]
```

Rules (carried from the `add-module` skill + migration best practice):
- Always test in **UAT first**; promote to PROD only after it passes.
- Every migration has a paired **`.up.sql` / `.down.sql`**.
- **Never edit an applied migration** — add a new one (checksum/version integrity).
- Use the **direct** (non-pooled) connection string for migrations.
- Seed/reference data via idempotent seed migrations (`ON CONFLICT DO NOTHING`);
  seeded rows use fixed hardcoded UUIDs (see `erd.md` §2).

> Makefile targets accept `DB_URL` so the same `migrate-up`/`migrate-down` run against
> either environment. Exact targets are defined during the Dev phase.

---

## 6. Deploy targets (free-tier)

| Component | Host | Env mapping |
|---|---|---|
| Go backend | **Render** (free) | **Two services** → UAT / PROD (see §8.1) |
| Next.js admin | **Vercel** (free) | preview → UAT, production → PROD |
| Flutter app | Built artifacts (store/side-load) | flavor uat / prod |
| Database | **Neon** | `cal-lowrisk-uat` / `cal-lowrisk-prod` |

> The design stays host-agnostic — anything that runs a Go binary with env vars works — so
> a future host swap (e.g. to Fly.io for Singapore-region latency) is config, not code.

---

## 7. Release / branch strategy (repo)

Kept lightweight for a solo team (contrast the heavier branch-per-env model in some
enterprise setups):

- **`main`** is the source of truth; no direct commits (branch + PR, per steering).
- Suggested: merging to `main` deploys to **UAT** automatically (CI/CD); **PROD** deploy
  is a **manual promotion/approval** step (a tag or a manually-triggered workflow).
- CI/CD injects per-env secrets from the platform's secret store (GitHub Actions secrets),
  mirroring §4. Never echo secrets into logs.

> This gives us the enterprise *promote* discipline (UAT → approve → PROD) without
> maintaining parallel long-lived code branches.

---

## 8. Decisions (resolved)

These were open items; now decided for the free-tier flow. Rationale is kept so the Dev
phase understands *why*.

### 8.1 Backend host — **Render, two services (UAT + PROD)**
One service per environment, mirroring the two-Neon-project isolation. Render chosen for
the simplest DX (connect repo → deploy), so effort goes into learning the code, not ops.
- **Consequence (cold start):** Render Free spins a service down after ~15 min idle, and
  Neon suspends after 5 min — so the **first request after idle is slow** (potentially
  ~20–40 s combined). Accepted for a learning MVP (see 8.3).
- **Config:** each Render service gets its own env vars/secrets per §4 (UAT values vs PROD
  values); same image/code, different config.
- **Revisit later:** if latency for Indonesian users or cold start becomes a real problem,
  re-evaluate a region-closer host (e.g. Fly.io, Singapore) or a paid tier — a host swap,
  not a code change.

### 8.2 PROD promotion trigger — **manual now, git tag later**
Merging to `main` auto-deploys **UAT**. **PROD** is promoted by a **manually-triggered
workflow** (GitHub Actions `workflow_dispatch`) — full control, zero overhead for a solo
team. **Later**, move to **git tags** (e.g. `v1.0.0`) so PROD releases are versioned and
line up with migration version numbers. Environment-approval gates are skipped for now
(little value when the approver is the same person).

### 8.3 Cold start — **accept it; NO keep-warm**
We deliberately do **not** add a keep-warm ping. On Neon's Free plan a keep-warm would
keep the DB from suspending and **burn the 100 CU-hour budget**, which could get PROD
suspended at month-end — counterproductive. Scale-to-zero is a *feature* at this tier.
- **Mitigation is client-side, not compute-side:** set sensible DB connect-timeout +
  retry so a cold start reads as a brief "loading", not an error; clients show a loading
  state on first request.
- **Revisit:** only if PROD gets real users and cold start becomes a genuine issue —
  then consider a paid tier, not a keep-warm hack.

### 8.4 Local dev DB — **local PostgreSQL (Docker)**
Day-to-day development runs against a **local Postgres in Docker**: free, fast, offline,
and it never touches the Neon budget. Our migrations (golang-migrate, plain Postgres) run
identically locally and on Neon. **UAT is not used as a dev DB** — that would blur the
dev/test boundary and make UAT unreliable for pre-PROD validation.

### 8.5 Seed data — **two layers**
1. **Structural/master seed via migrations** (idempotent, `ON CONFLICT DO NOTHING`):
   roles + a base Indonesian-oriented `food`/`exercise` catalog. This is **shared by UAT
   and PROD** — the catalog is reference data PROD needs too. Seeded rows use fixed
   hardcoded UUIDs (see `erd.md` §2).
2. **UAT-only test seed via a separate script** (e.g. `scripts/seed-uat.*`), NOT a
   migration: test user + admin accounts and a few sample meal/workout logs, so it can
   **never leak into PROD**. Scope kept small — enough to validate the daily/weekly/monthly
   summaries meaningfully.

> **Never clone PROD data into UAT** — it would violate the privacy boundary (§4.1 / ADM-6;
> admins must not see users' private logs), and PROD has no data to clone anyway.

---

## 9. References

- Neon Free plan limits & quotas — https://neon.com/faqs/free-plan-limits-and-quotas
- Neon Free tier compute/branch suspension behavior — https://neon.tech/docs/reference/technical-preview-free-tier/
- Neon connection pooling (PgBouncer, `-pooler` host) — https://neon.com/faqs/postgres-hosting-options-auto-pause-database
- Neon auto-suspend (5-min, fixed on Free) — https://neon.com/faqs/databases-automatically-scale-serverless-environments

> Figures reflect Neon's Free plan as of 2026-10 and may change; verify against the live
> pages above before provisioning.
