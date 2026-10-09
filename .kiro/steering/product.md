---
inclusion: always
---

# cal-lowrisk — Project Context

> This file is auto-loaded into every Kiro session in this workspace. It is the
> "memory" of the project so any session (or any Kiro agent) immediately knows what
> we're building, why, and the decisions made so far. Keep it current as decisions evolve.

## What we're building

**cal-lowrisk** is a personal **health tracker** that combines a **calorie tracker**
and a **workout tracker**. Users log meals and workouts; the app shows whether their
day/week/month is **balanced against their personal goal** (lose / maintain / gain
weight), including a **macronutrient** (protein / carbs / fat) breakdown.

- **Problem statement:** "People find it hard to know whether their food intake and
  physical activity in a day are balanced against their goal."
- **Audience:** general; users self-select their goal on sign-up. Food content is
  oriented toward **Indonesian meals**.
- **Nature:** a learning project to practice a modern enterprise-style stack end to end.
- ⚠️ It is an educational/self-tracking tool, **not** medical advice or a medical device.

## Tech stack (decided)

| Layer | Choice |
|---|---|
| Backend | **Go** (Gin + GORM) |
| Mobile (user app) | **Flutter** |
| Web admin | **Next.js** |
| Database | **PostgreSQL** — **Neon** (free tier, AWS Singapore region) |
| Architecture | **Modular monolith**, module-first, in a **monorepo** (`apps/`) |
| Migrations | **golang-migrate** (SQL versioned, per-module) — see skill `add-module` |
| Mobile UI source | **Figma** (via Figma Power) — design-to-code into Flutter |

### Key tech decisions & rationale
- **Modular monolith, not microservices** — cheap (one deploy, free-tier friendly),
  simple, yet structured so a module can later be extracted into a service with minimal
  refactor. Layering mirrors the user's workplace Go standard.
- **Module-first layout** (`internal/modules/<module>/...`), each module keeping the same
  layers: **Controller → Service → Repository → Model + DTO**. Dependencies flow inward;
  use interfaces for DI.
- **GORM for runtime queries, NOT for schema.** Schema (DDL) is owned by SQL migrations
  via **golang-migrate**. Do **not** use GORM AutoMigrate for structure — we want explicit,
  reviewable, reversible SQL (user wants to "see the detail" and match enterprise practice).
- **Migrations are split per module** (Option 2 / "Cara A"): one folder per module under
  `migrations/<module>/`, each with its own tracking table, executed in dependency order
  via the Makefile. See the `add-module` skill for the exact procedure.
- **Firebase:** not used for now (would hide the backend we want to learn). Only FCM is a
  possible later add for push notifications.
- **Redis:** deferred. The `summary` module is designed with a `Cache` interface so Redis
  (e.g. Upstash free tier) can be added later without redesign.
- **Monorepo layout:** one repo, apps separated — `apps/backend/` (Go modular monolith,
  one `go.mod`/binary/deploy), `apps/mobile/` (Flutter, later), `apps/admin/` (Next.js,
  later). Separating apps is NOT microservices; the backend stays a single monolith with
  modules in `internal/modules/`.
- **Primary keys:** native PostgreSQL `uuid` on every table, app-generated **UUIDv7**
  (Go `uuid.NewV7()` via `google/uuid` v1.6+). Not bigint, not varchar. Seeded rows use
  fixed hardcoded UUIDs. (Also saved as a global learning.) See `docs/02-system/erd.md` §2.
- **Enum-like values** (sex, activity_level, goal, meal_type, food.source,
  exercise.unit_type) use **CHECK constraints**; only `role` is a seeded table.
- **Environments:** two **separate Neon projects** — `cal-lowrisk-uat` and
  `cal-lowrisk-prod` (compute quota is per-project → isolated budgets). Backend hosted on
  **Render** (two services). One codebase, config differs per env. Migrations use the
  **direct** (non-pooled) conn; app runtime uses the **pooled** (`-pooler`) conn with a
  small pool. PROD promotion is manual (`workflow_dispatch`) now, git tags later. No
  keep-warm (would burn Neon's free CU-hours). Local dev DB: Postgres via Docker — but in
  the Kiro sandbox we use Neon directly. See `docs/02-system/deployment.md`.
- **Visual design direction:** warm, clean, calm; hero number on Home; color-coded macros
  (protein=blue, carbs=green, fat=amber). Tokens: font **Inter** (NOT SF Pro), primary
  orange **#FF6B3D**, 4px spacing scale, soft radii, **Lucide** icons (MIT). Mobile nav =
  bottom tab bar + center FAB. See `docs/03-ux/user-flows.md` §11.

## Modules (bounded contexts)

| Module | Responsibility | Data kind |
|---|---|---|
| `user` | Auth, profile (weight/height/age/sex/activity), goal, role | master (per user) |
| `food` | Food catalog + calories + macros (admin-curated) | master/reference |
| `nutrition` | Meal logs (per user); references catalog food OR custom food + portion | transactional |
| `exercise` | Exercise catalog + unit type + calorie-burn data (admin-curated) | master/reference |
| `workout` | Workout logs (per user); references catalog exercise + amount | transactional |
| `summary` | Daily/weekly/monthly balance: calories in/out/target + macros | aggregation |

**Master vs transactional data distinction is deliberate** and reflected in migrations:
all structure via migrations; master/reference data (roles, catalogs) gets dedicated
**seed migrations**.

## Scope (MVP)

**In scope:** auth + profile + goal; nutrition logging (catalog + portion + **custom
food**, private per user, source-tagged); workout logging (**catalog only**, no custom);
macros (protein/carbs/fat); daily/weekly/monthly summaries vs a personalized BMR/TDEE
target; admin = catalog CRUD (food + exercise) + aggregate dashboard + user-account
management (list, activate/deactivate via soft delete).

**Out of scope (roadmap):** barcode/photo input, social/sharing, push notifications,
meal plans, smartwatch/Google Fit, water/sleep/mood, custom exercises, health
score/streaks, admin audit trail.

### Roles & privacy (hard rule)
Two roles: **User** (Flutter) and **Admin** (Next.js). Roles are modeled as a
first-class entity (not a boolean) to allow future roles. **Privacy boundary:** admin
manages *accounts* (activate/deactivate) and sees *aggregate* stats only — admin must
**never** access an individual user's private meal/workout logs. Enforce this server-side.

## Working process (we play every role)

**BA → SA → UX → Dev**, documented before built (docs-as-code). We use branch + PR for
every change (no direct commits to `main`). Current repo: `Monarqhs/cal-lowrisk`.

### Progress so far
- ✅ **BA:** repo + docs structure + root README (PR #1); **BRD** `docs/01-business/brd.md` (PR #2).
- ✅ **SA:** `docs/02-system/` complete — `erd.md` (data model, UUIDv7, 5 BRD open
  questions resolved), `architecture.md` (modular monolith, layering, privacy boundary),
  `deployment.md` (two Neon projects, Render, migration promotion). (PRs #4–#8.)
  - The 5 BRD §10 questions are RESOLVED (see erd.md §1): per-100g + optional serving;
    reusable custom food (one `food` table, `source`+`owner_user_id`); summary computed
    on-the-fly (no table) behind a Cache interface; MET + body weight for burn;
    Mifflin-St Jeor → TDEE → goal for target, computed on-the-fly.
- ✅ **UX:** `docs/03-ux/user-flows.md` complete (PR #9) + 5 UX open items resolved +
  visual direction (PR #10). Flutter screens drafted in **Figma** (file "cal-lowrisk —
  Mobile", team "callium project"): all 13 user screens (U1–U15; Home has a ring variant
  and the chosen **Balance Bar** variant on the "Explore" page). Admin (A1–A7) not yet
  designed.
- 🔧 **Dev (in progress):** branch `feat/backend-scaffold`.
  - ✅ `apps/backend/` scaffolded: Gin entry point, shared layer (config, GORM database,
    response envelope, Base model with UUIDv7), Makefile (golang-migrate targets),
    `migrations/{user,food,exercise,nutrition,workout}/`. Builds + vets clean (Go 1.25).
  - ✅ `user` module migrations WRITTEN (not yet run): `migrations/user/000001..000004`
    = role, users, user_profile, seed_roles.
  - ⏳ **Neon:** project `cal-lowrisk-uat` created (project-id `rough-field-65178844`,
    AWS Singapore). Neon Power connected but tool enumeration was flaky in-session.
  - ⏳ **Blocked on:** getting the Neon connection strings into `apps/backend/.env`
    (git-ignored) — pooled (`DATABASE_URL`) + direct (`DATABASE_URL_DIRECT`). Then: run
    `user` migrations, verify tables + seed, build the `user` module API (register/login/
    profile with JWT), test endpoints (Postman Power), then mock mobile in Flutter.
  - ⚠️ `golang-migrate` CLI failed to install in-sandbox (Go toolchain version clash) —
    will run migrations via Neon Power `execute_sql` or a small Go runner instead.
  - ⚠️ Do NOT run Neon's generic TS deploy flow (`neon config init`/`neon.ts`/`neon deploy`)
    — it conflicts with golang-migrate owning the schema. We only need the connection string.

### Dev plan (agreed order)
master-table migrations → backend API (per module, start with `user` vertical slice) →
test endpoints → mock mobile in Flutter. Build one module end-to-end before the next.

## Pointers
- Full requirements: `docs/01-business/brd.md`
- Data model: `docs/02-system/erd.md` · Architecture: `docs/02-system/architecture.md` ·
  Environments/deploy: `docs/02-system/deployment.md`
- UX flows + screen inventory + visual tokens: `docs/03-ux/user-flows.md`
- Documentation index: `docs/README.md`
- How to add a module (code + migrations): skill `add-module`
- Backend code + how-to: `apps/backend/` (+ its README)
- Figma (mobile): file "cal-lowrisk — Mobile" in team "callium project" (via Figma Power)
- Powers connected: Figma (design-to-code), Neon (DB), Postman (API testing)

## For a new session (handover)
This file is auto-loaded. To resume Dev: the backend is on branch `feat/backend-scaffold`.
Next concrete step is getting Neon connection strings into `apps/backend/.env` (pooled +
direct), then running the `user` migrations and verifying `role`/`users`/`user_profile` +
seeded roles. Prefer the Neon Power `get_connection_string` (project-id
`rough-field-65178844`); if its tools don't enumerate, ask the user to paste the pooled
and direct strings from the Neon dashboard Connect page.

#[[file:docs/01-business/brd.md]]
