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
| Database | **PostgreSQL** (free tier: Neon or Supabase) |
| Architecture | **Modular monolith**, module-first |
| Migrations | **golang-migrate** (SQL versioned, Flyway-style) — see skill `add-module` |

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
- ✅ Repo created; docs structure + root README + `.gitignore` (PR #1, merged).
- ✅ **BRD** written: `docs/01-business/brd.md` (PR #2).
- 🔜 **Next: ERD + architecture** in `docs/02-system/` (SA phase). The BRD §10 lists 5
  open questions to resolve during ERD:
  1. Food portion basis (per serving vs per 100 g).
  2. Custom food persistence (intent: reusable personal foods).
  3. Summary computation: on-the-fly vs precomputed; where the cache boundary sits.
  4. Calorie-burn estimation per exercise unit type (reps/duration/steps); body weight?
  5. Daily target formula (BMR/TDEE variant + activity multipliers).

## Pointers
- Full requirements: `docs/01-business/brd.md`
- Documentation index: `docs/README.md`
- How to add a module (code + migrations): skill `add-module`

#[[file:docs/01-business/brd.md]]
