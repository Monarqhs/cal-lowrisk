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
- **Schema-per-module:** each module owns a dedicated PostgreSQL **schema named after the
  module** (`"user".users`, `food.food`, `nutrition.nutrition_log`, …), not a single
  `public`. Mirrors the modular-monolith boundaries + eases the microservice extraction
  path. `user` is a reserved word → the `"user"` schema is **always double-quoted**. Each
  module's migration tracking table lives **inside its own schema**
  (`"<module>"."schema_migrations"`). Cross-schema FKs are allowed, but cross-module data
  access goes through the other module's **service interface** (never a direct cross-schema
  JOIN). Models set schema via GORM `TableName()` + `model.Qualify(schema, table)`. See
  erd.md §2 + the `add-module` skill.
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
- 🔧 **Dev (in progress):** branch model in use: `feat/*` → PR → **`uat`** (tests on Neon
  `cal-lowrisk-uat`) → PR → `main`. `main` IS production now (landing is deployed from it).
  - ✅ `apps/backend/` scaffolded: Gin entry, shared layer (config, GORM database, response
    envelope, Base model with UUIDv7 + `Qualify`, **`internal/app.NewRouter`**, **`shared/auth`**
    = bcrypt + HS256 JWT + Gin middleware), Makefile, migrations. Builds + vets clean (Go 1.25).
  - ✅ **Neon UAT connected.** `apps/backend/.env` (git-ignored): pooled `DATABASE_URL` +
    direct `DATABASE_URL_DIRECT` + JWT. Project `cal-lowrisk-uat` = project-id
    `rough-field-65178844`, db `neondb`, PG 18, AWS Singapore. `.env` URLs must be
    double-quoted (the `&` in the query string breaks `source .env` otherwise).
  - ✅ **Schema-per-module** applied; `user` module lives in the `"user"` schema; migrations
    `000001..000004` run & verified on UAT (tracking `"user"."schema_migrations"` v4).
  - ✅ **`user` module BUILT, TESTED, MERGED** (PR #12→uat, #13 tests→uat, #17 uat→main).
    Onboarding vertical slice done: `POST /api/v1/auth/register`, `/auth/login` (JWT),
    `/me/profile` (auth+role; computes Mifflin-St Jeor → TDEE → goal target). Layered
    controller→service→repository→model+dto per the `add-module` skill. Smoke-tested live
    on UAT via Postman.
  - ✅ **Contract-first API specs** at `docs/02-system/api-specs/` (README conventions +
    `user-service/01. Onboarding User.md` with a numbered Test cases list TC-ON-*).
  - ✅ **Postman collection** (versioned) at `docs/02-system/api-request/`
    (`cal-lowrisk.v1.0.0.postman_collection.json` + local/uat environments + CHANGELOG).
  - ✅ **Test suite** at `apps/backend/tests/` (unit + E2E). E2E auto-skips unless
    `TEST_DATABASE_URL` is set; run it against a disposable **Neon branch** (create from
    uat via Neon Power, delete after). All TC-ON-* green. `make test-unit` / `test-e2e`.
  - ⚠️ `golang-migrate` CLI + `golangci-lint` fail to install in-sandbox (Go toolchain
    clash). Run migrations via Neon Power `run_sql` (one stmt at a time), keep
    `"<module>"."schema_migrations"` in sync; rely on `go vet` + `gofmt` for linting.
  - ⚠️ Neon `cal-lowrisk-prod` NOT created yet (deferred; create at first real prod backend
    deploy — Render, per deployment.md). UAT branch is the Neon prod-equivalent for now.
  - ⚠️ Do NOT run Neon's generic TS deploy flow — golang-migrate owns the schema.
  - ⏭️ **Next backend:** build the `food` module (catalog) end-to-end — spec → service →
    tests, same pattern. Then `exercise` → `nutrition` → `workout` → `summary`.

- ✅ **Landing page + domain + email (Callium marketing):** done to support a Claude for
  Startups application.
  - `apps/landing/` — static site (HTML/CSS/ESM, no build), **Notion-style line-art**,
    Anton (display, placeholder for **Harse**) + Poppins. Passed an anti-slop pass against
    the **nutlope/hallmark** rule-set. Assets optimized to WebP (~90 KB total). Hero image
    + logo are the user's own assets; `--paper` is `#f4f1ec` to match the hero bg.
  - **Deployed on Cloudflare Pages** (project `cal-lowrisk`, auto-deploy from `main`,
    output dir `apps/landing`). Live at `cal-lowrisk.pages.dev`.
  - **Domain `callium.app`** bought on Cloudflare (~$8.20/yr) → custom domain attached to
    Pages (CNAME `@` → `cal-lowrisk.pages.dev`); was "Initializing" at handoff (SSL
    provisioning, up to ~15 min).
  - ✅ **Email Routing DONE** on `callium.app` (DNS records Locked). Destination
    `ghalihyulio12@gmail.com` added + **Verified**. Two routing rules **Active**, both
    forwarding to that Gmail: `call.help@callium.app` (public/support, used in the landing
    footer) and `yulio.ananda@callium.app` (founder identity, used for Claude Console +
    program/official comms). Forwarding **tested end-to-end** (inbound email arrives in the
    Gmail inbox). Catch-all rule left Disabled (Drop) by choice.
    - Optional next (not done): enable Gmail "Send mail as" so the user can reply from an
      `@callium.app` address (needs SMTP relay — Zoho free / Resend), and/or a full mailbox.
  - 🔁 **Claude for Startups — APPLIED TWICE; 1st REJECTED, 2nd pending.** Created a
    **Claude Console Organization** `Callium` (entity: Small/medium business; Indonesia)
    under founder email `yulio.ananda@callium.app`. Website `https://callium.app`, founded
    Oct 2026, positioned as an **educational self-tracking tool, not medical advice**.
    Skipped buying API credits (org credits $0.00). Answered honestly throughout:
    early-stage/**bootstrapped, no outside funding** ("Not yet raised"; the "raised outside
    funding" box is left UNCHECKED and must stay that way — self/parent money without a real
    documented equity round is NOT "outside funding"; a false claim is unverifiable and
    backfires).
    - **1st application → REJECTED** ("we couldn't approve Callium's application this time").
      Email said the program is for *verified* startups founded in the last 5y or funded in
      the last 2y; rejections are usually either out-of-window OR **couldn't be verified from
      the details**. Reapply tips given: **(1) use a work email matching the company domain**,
      (2) add company website (already done). Likely root cause = **verification**, not
      eligibility — most probably because the domain was **bought the same day** (too new)
      and `yulio.ananda@callium.app` is **forwarding-only (Cloudflare Email Routing), not a
      real mailbox**.
    - **Contact channels explored and are DEAD ENDS:** the rejection came from a marketing
      address on `mail.anthropic.com`, which Anthropic's help docs confirm are
      **outbound-only / unmonitored**; the Console **"Get help" Fin AI agent** has **no
      knowledge of the Startups program** and can't help. There is no easy human-reconsider
      path — the intended flow is simply **fix the issues and reapply**.
    - **2nd application → SUBMITTED, pending review** (same honest data; industry set to
      "Machine Learning" this time, which is slightly off — Callium is a *consumer*
      health/fitness app that *uses* AI, not an ML company; "Consumer"/"Health & Fitness"
      would be more accurate, but industry is optional and not a core verification criterion).
      ⏳ Decision usually arrives in minutes–3 business days via email to `yulio.ananda@`
      (forwards to Gmail — watch inbox + spam).
    - **If rejected again, the real fix (do these before a 3rd try):** (a) let the domain
      **age** ~1–2 weeks; (b) set up a **real mailbox** on `@callium.app` (e.g. Zoho free or
      Google Workspace) instead of forwarding-only — NOTE this needs the MX records switched
      away from Cloudflare Email Routing, which would break the current `call.help@`/
      `yulio.ananda@`→Gmail forwarding, so decide the trade-off deliberately. Don't rush
      another reapply while the domain is brand new.
    - ⚠️ Context: Anthropic **paused** the program's two headline perks ($1,000 API credits
      + 1 free year of Claude Team) on ~9 Oct 2026 after a demand surge (re-reviewing). So
      the application is likely queued/under re-review; benefit grant is NOT guaranteed and
      credits often require institutional funding (which Callium doesn't have).

### Dev plan (agreed order)
Per module, build end-to-end before the next: **API spec (contract-first, with a Test
cases list) → migrations (Neon) → service (controller→service→repository→model+dto) →
Postman collection entry → unit + E2E tests (E2E on a disposable Neon branch) → PR to
`uat` → promote `uat`→`main`.** Order: `user` ✅ → `food` → `exercise` → `nutrition` →
`workout` → `summary`. Then mock the mobile app in Flutter (Figma → Flutter).

### Team / communication preferences (saved as global learnings too)
- Converse with the user in **Bahasa Indonesia**; keep everything written to the repo
  (code, docs, commits, PR titles/bodies, branch names, comments) in **English**.
- DB: native `uuid` PKs, app-generated **UUIDv7** (see erd.md §2).
- Testing: contract-first; specs are source of truth; tests live centralized under
  `apps/backend/tests/` (unit + E2E); E2E uses a disposable Neon branch.

## Pointers
- Full requirements: `docs/01-business/brd.md`
- Data model: `docs/02-system/erd.md` · Architecture: `docs/02-system/architecture.md` ·
  Environments/deploy: `docs/02-system/deployment.md`
- **API specs (contract-first, source of truth):** `docs/02-system/api-specs/`
- **Postman collection (versioned):** `docs/02-system/api-request/`
- UX flows + screen inventory + visual tokens: `docs/03-ux/user-flows.md`
- Documentation index: `docs/README.md`
- How to add a module (code + migrations): skill `add-module`
- Backend code + how-to: `apps/backend/` (+ its README) · tests: `apps/backend/tests/`
- **Landing site:** `apps/landing/` (static; Cloudflare Pages → `cal-lowrisk.pages.dev` →
  `callium.app`) (+ its README)
- Figma (mobile): file "cal-lowrisk — Mobile" in team "callium project" (via Figma Power)
- Powers connected: Figma (design-to-code), Neon (DB), Postman (API testing)
- Hosting: Cloudflare (Pages for landing + Registrar for `callium.app` + Email Routing)

## For a new session (handover)
This file is auto-loaded. State as of this handoff:
- **Backend:** `user` module (onboarding) is built, tested, and merged to `main`. Neon
  `cal-lowrisk-uat` (project-id `rough-field-65178844`) connected; `apps/backend/.env`
  holds the (double-quoted) pooled+direct URLs + JWT. **Next:** the `food` module,
  end-to-end, following the Dev plan above and the `add-module` skill.
- **Landing/marketing:** `apps/landing/` is LIVE on Cloudflare Pages at
  `cal-lowrisk.pages.dev`, auto-deploying from `main`. Domain `callium.app` bought on
  Cloudflare and attached as a custom domain (may still be finishing SSL).
- **✅ Email Routing DONE:** Cloudflare Email Routing on `callium.app` fully set up —
  destination `ghalihyulio12@gmail.com` Verified; two Active rules forwarding to it:
  `call.help@callium.app` (public/support) and `yulio.ananda@callium.app` (founder); inbound
  forwarding tested end-to-end. Catch-all left Disabled (Drop). Optional next (not done):
  Gmail "Send mail as" to reply from an `@callium.app` address, and/or a full mailbox.
- **🔁 Claude for Startups — 1st REJECTED, 2nd pending:** Claude Console Org `Callium` (SMB,
  Indonesia) under `yulio.ananda@callium.app`, website `callium.app`, answered honestly as
  bootstrapped/no funding. 1st app rejected on a **verification** basis (likely: domain bought
  same day = too new, and the email is **forwarding-only, not a real mailbox**). Contact paths
  are dead ends (marketing addresses are outbound-only; Console "Get help" Fin doesn't know the
  program). 2nd app submitted, awaiting review (email to `yulio.ananda@`→Gmail; watch spam).
  **Before a 3rd try:** let the domain age ~1–2 weeks AND set up a real `@callium.app` mailbox
  (Zoho/Workspace — breaks the current Cloudflare forwarding; decide the trade-off). Headline
  perks are paused by Anthropic (~9 Oct 2026), so benefits aren't guaranteed regardless.
- **Fonts:** landing title uses **Anton** as a placeholder for **Harse** (not free on
  Google Fonts). When the user supplies a licensed Harse file, self-host via `@font-face`
  and update `--font-display` in `apps/landing/styles/main.css`.
- If Neon Power tools stop enumerating, starting a fresh session fixed it before.
  `cal-lowrisk-prod` (Neon) still not created — do that at first real backend deploy (Render).

#[[file:docs/01-business/brd.md]]
