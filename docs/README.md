# cal-lowrisk — Documentation

Welcome to the documentation home for **cal-lowrisk**. This is the single source of
truth for the project. Everything here is **docs-as-code**: plain Markdown (and Mermaid
for diagrams), versioned alongside the code and reviewed through pull requests.

## How this is organized

We work through the project the way a real delivery team does, except the "team" is
just the two of us (you + Kiro) wearing every hat in sequence:

| Phase | Role | Folder | Produces |
|---|---|---|---|
| 1️⃣ | Business Analyst (BA) | [`01-business/`](./01-business) | BRD — *what & why* |
| 2️⃣ | System Analyst (SA) | [`02-system/`](./02-system) | ERD + architecture — *how data & system are modeled* |
| 3️⃣ | UX Designer | [`03-ux/`](./03-ux) | User flows & wireframes — *how it's used* |
| 4️⃣ | Developer | (code) | The implementation |

Supporting reference data (food/exercise catalogs, etc.) lives in [`04-data/`](./04-data).

## Index

### 01 — Business
- [`brd.md`](./01-business/brd.md) — Business Requirements Document ✅

### 02 — System
- [`erd.md`](./02-system/erd.md) — Entity Relationship Diagram (Mermaid) ✅
- `architecture.md` — Architecture decisions (modular monolith, layering) *(next)*

### 03 — UX
- `user-flows.md` — User journeys & flows *(later)*

### 04 — Data
- Catalog data (food, exercise) — CSV or links to spreadsheets *(later)*

## Conventions

- **Format:** Markdown for text, [Mermaid](https://mermaid.js.org/) for diagrams
  (GitHub renders Mermaid natively — diagrams stay versioned and diff-able).
- **Source of truth:** this repo. Heavy binary assets (visual wireframes, large
  spreadsheets) may live in Google Drive / Figma and be *linked* from here.
- **Process:** we document before we build. BRD is finalized before the ERD, because
  the data model should follow business needs — not the other way around.
