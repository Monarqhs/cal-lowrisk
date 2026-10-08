# cal-lowrisk

> Health tracker — calorie & workout monitoring. Track what goes in, what goes out, and whether your day was balanced.

**cal-lowrisk** is a personal health-tracking project that combines a **calorie tracker** and a **workout tracker**. Users log their meals and workouts; the app computes a daily/weekly/monthly "healthiness" picture (calories in vs. calories out vs. target).

This is a learning project built to practice a modern, enterprise-style stack end to end.

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | **Go** (Gin + GORM) |
| Mobile (user app) | **Flutter** |
| Web Admin | **Next.js** |
| Database | **PostgreSQL** |
| Architecture | **Modular Monolith** (module-first, layered) |

## Architecture at a Glance

A **modular monolith**: one deployable application, internally split into modules
(`user`, `food`, `nutrition`, `exercise`, `workout`, `summary`) with clear boundaries.
Each module keeps the same layering — Controller → Service → Repository → Model + DTO.

> Why modular monolith? Cheap to run (one deploy, free-tier friendly), simple to reason
> about, yet structured so any module can later be extracted into a microservice with
> minimal refactor.

## Documentation

All project documentation lives in [`docs/`](./docs). Start at the
[documentation index](./docs/README.md).

Our working process (we play every role): **BA → SA → UX → Dev**
1. **BRD** — Business Requirements ([`docs/01-business/`](./docs/01-business))
2. **ERD + Architecture** — System design ([`docs/02-system/`](./docs/02-system))
3. **User Flows** — UX ([`docs/03-ux/`](./docs/03-ux))
4. **Data** — Catalog data ([`docs/04-data/`](./docs/04-data))

## Status

🚧 Early stage — currently defining business requirements. No code yet by design;
we document before we build.
