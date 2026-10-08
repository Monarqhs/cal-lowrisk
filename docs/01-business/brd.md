# Business Requirements Document (BRD) — cal-lowrisk

| Field | Value |
|---|---|
| **Project** | cal-lowrisk |
| **Document** | Business Requirements Document (BRD) |
| **Version** | 0.1 (Draft) |
| **Status** | Draft — pending review |
| **Author** | Product Owner + Kiro (BA) |
| **Last updated** | 2026-10-04 |

> **Process note:** This BRD defines *what* we build and *why*. The data model (ERD)
> and system design follow in `docs/02-system/`, and are derived from this document —
> not the other way around. Requirements here are the source of truth for scope.

---

## 1. Executive Summary

**cal-lowrisk** is a personal health-tracking application that combines a **calorie
tracker** and a **workout tracker** into a single daily picture of balance.

Users log the food they eat and the workouts they do. The app then shows whether their
day — and their week and month — is **balanced against their personal goal** (lose,
maintain, or gain weight), by comparing calories consumed, calories burned, and a
personalized daily target, with a macronutrient (protein / carbs / fat) breakdown.

The product targets a general audience, with food data oriented toward **Indonesian
meals**. It is built as a learning project across a modern stack (Go, Flutter, Next.js,
PostgreSQL) using a **modular monolith** architecture.

> ⚠️ **Disclaimer:** cal-lowrisk is a self-tracking and educational tool. Calorie and
> macronutrient figures are estimates. It does **not** provide medical, dietary, or
> professional health advice, and is not a medical device.

---

## 2. Problem Statement & Goals

### 2.1 Problem Statement

> **"People find it hard to know whether their food intake and physical activity in a
> day are balanced against their goal."**

Existing tools are often too complex, or their food databases do not reflect local
(Indonesian) meals, making consistent tracking difficult.

### 2.2 Product Goals

| # | Goal |
|---|---|
| G1 | Let users easily log meals (with portions) and workouts in a day. |
| G2 | Show a clear daily balance: calories **in** vs **out** vs **target**, plus status (surplus / deficit / balanced). |
| G3 | Aggregate that balance over **weekly** and **monthly** views to show progress over time. |
| G4 | Personalize the daily target based on the user's self-selected goal and body profile. |
| G5 | Include a **macronutrient** (protein / carbs / fat) breakdown, not just total calories. |
| G6 | Provide admins the tools to curate the food and exercise catalogs that make the app useful. |

### 2.3 Success Indicators (qualitative, for a learning MVP)

- A user can complete the core loop: **log a meal → log a workout → see today's balance vs their goal.**
- Weekly and monthly summaries render correctly from logged data.
- An admin can populate and maintain the food and exercise catalogs.

---

## 3. Target Users & Personas

### 3.1 Audience

General audience interested in tracking health. Users **self-select their goal** on
sign-up, so the product serves multiple intents rather than one niche. Food content is
oriented toward **Indonesian cuisine**.

### 3.2 Personas

**Persona A — "The Goal-Setter" (primary user)**
- Wants to lose / maintain / gain weight and needs to know if today supports that goal.
- May be a beginner; wants clarity over raw data dumps.
- Eats local food that generic apps often lack.

**Persona B — "The Admin / Curator"**
- Maintains the quality and coverage of the food and exercise catalogs.
- Manages user accounts (activate / deactivate) and views aggregate statistics.
- Does **not** inspect individual users' private meal/workout data.

---

## 4. User Roles

For the MVP there are **two roles**. The system models roles as a first-class concept so
new roles (e.g. Nutritionist/Coach) can be added later without rework.

| Role | Client | Responsibilities |
|---|---|---|
| **User** | Flutter (mobile) | Log meals & workouts; view personal summaries (daily/weekly/monthly) vs goal. |
| **Admin** | Next.js (web) | Manage food & exercise catalogs; view aggregate statistics; manage user accounts. |

> **Design note:** roles are modeled as a separate entity (not a boolean `is_admin`),
> so scaling to more roles later is additive.

### 4.1 Admin / Privacy Boundary (important)

- Admin manages **accounts** (activate / deactivate via soft delete) and sees **basic
  account info** and **aggregate** statistics only.
- Admin does **not** access an individual user's private meal or workout logs.
- This boundary is a hard requirement for the MVP and must be enforced in the backend.

---

## 5. Scope

### 5.1 In Scope (MVP)

- User registration, authentication, and profile with **body data** (e.g. weight,
  height, age, sex, activity level) and a **self-selected goal** (lose / maintain / gain).
- **Nutrition logging:** pick food from the catalog **or** create a **custom food**;
  specify a **portion**; calories & macros computed accordingly.
- **Workout logging:** pick an exercise from the catalog and record the amount
  (reps / duration / steps depending on exercise type); estimated calories burned.
- **Macronutrients** (protein / carbs / fat) tracked for catalog foods and in summaries.
- **Summaries:** daily, weekly, and monthly views of calories in / out / net vs target,
  with surplus / deficit / balanced status and a macro breakdown.
- **Personalized daily target** derived from the user's profile and goal (BMR/TDEE-style
  estimation).
- **Admin:** CRUD food catalog; CRUD exercise catalog; read-only aggregate dashboard;
  user-account management (list, activate/deactivate via soft delete).

### 5.2 Out of Scope (MVP) — deferred to roadmap

| Deferred feature | Rationale |
|---|---|
| Barcode scan / food photo recognition | Requires external API / ML; high complexity. |
| Social features / sharing / leaderboards | Needs inter-user relationships & privacy design. |
| Push notifications / reminders | Needs FCM + scheduling; fits a later phase. |
| Meal plans / recommendations | Needs complex logic / ML. |
| Smartwatch / Google Fit integration | Device integration complexity. |
| Water, sleep, mood tracking | Widens focus beyond calories + workout. |
| Custom **exercise** creation | Catalog is sufficient; custom calorie estimates are unreliable. |
| Health **score (0–100)** and **streaks** | Scoring definition needs iteration; a Phase 2 engagement feature. |
| Admin **audit trail** | Infrastructure prepared conceptually; enabled in a later phase. |

---

## 6. Functional Requirements

Requirements are grouped by module to align with the modular-monolith architecture.
Each requirement uses **MoSCoW** priority: **M**ust / **S**hould / **C**ould.

### 6.1 Module: `user` (Profile, Goal)

| ID | Requirement | Priority |
|---|---|---|
| USR-1 | A user can register with email + password. | M |
| USR-2 | A user can log in and receive an auth token (JWT). | M |
| USR-3 | A user has a profile with body data: weight, height, age, sex, activity level. | M |
| USR-4 | A user selects a **goal**: lose / maintain / gain weight. | M |
| USR-5 | A user can update their profile and goal; targets recalculate accordingly. | M |
| USR-6 | Each user is assigned a **role** (default: User). | M |

### 6.2 Module: `food` (Food Catalog — admin-curated)

| ID | Requirement | Priority |
|---|---|---|
| FOOD-1 | The catalog stores foods with name, **calories**, **protein**, **carbs**, **fat**. | M |
| FOOD-2 | Each food defines a **portion basis** (e.g. per serving or per 100 g) used for calculation. | M |
| FOOD-3 | Users can search/browse the catalog when logging a meal. | M |
| FOOD-4 | Admins can create, read, update, and delete catalog foods. | M |
| FOOD-5 | Catalog is oriented toward Indonesian foods (seed data). | S |

### 6.3 Module: `nutrition` (Meal Logging — user data)

| ID | Requirement | Priority |
|---|---|---|
| NUT-1 | A user can log a meal entry for a given date and meal type (breakfast / lunch / dinner / snack). | M |
| NUT-2 | A meal entry references either a **catalog food** or a **custom food**, with a **portion**. | M |
| NUT-3 | Calories & macros for an entry are computed from the food data × portion. | M |
| NUT-4 | A user can create a **custom food** (name + calories; macros **optional**), **visible only to that user**. | M |
| NUT-5 | Each food reference is tagged with its **source** (`catalog` or `custom`). | M |
| NUT-6 | A user can view, edit, and delete their own meal entries. | M |
| NUT-7 | Custom foods created by a user are reusable in future entries. | S |

### 6.4 Module: `exercise` (Exercise Catalog — admin-curated)

| ID | Requirement | Priority |
|---|---|---|
| EXE-1 | The catalog stores exercises with name and a **unit type** (reps / duration / distance-steps). | M |
| EXE-2 | Each exercise stores the data needed to **estimate calories burned** for a given amount. | M |
| EXE-3 | Users can search/browse the exercise catalog when logging a workout. | M |
| EXE-4 | Admins can create, read, update, and delete catalog exercises. | M |
| EXE-5 | **No custom exercises** in MVP (catalog only). | M |

### 6.5 Module: `workout` (Workout Logging — user data)

| ID | Requirement | Priority |
|---|---|---|
| WRK-1 | A user can log a workout entry for a given date, referencing a catalog exercise. | M |
| WRK-2 | A user records the **amount** in the exercise's unit (reps / minutes / steps). | M |
| WRK-3 | Estimated calories burned are computed from the exercise data × amount (and body data where relevant). | M |
| WRK-4 | A user can view, edit, and delete their own workout entries. | M |

### 6.6 Module: `summary` (Healthiness / Balance)

| ID | Requirement | Priority |
|---|---|---|
| SUM-1 | Compute a **daily target** (calories, and macro targets) from the user's profile + goal (BMR/TDEE-style). | M |
| SUM-2 | **Daily summary:** calories in (meals), calories out (workouts), net, target, and status (surplus / deficit / balanced). | M |
| SUM-3 | Daily summary includes a **macro breakdown** (protein / carbs / fat vs targets). | M |
| SUM-4 | **Weekly summary:** aggregates/averages across the week with on-track indication. | M |
| SUM-5 | **Monthly summary:** aggregates/averages across the month with trend indication. | M |
| SUM-6 | Summary computation is designed to allow a **cache layer** (e.g. Redis) to be added later without redesign. | S |

### 6.7 Module: `admin` (Web Admin — Next.js)

| ID | Requirement | Priority |
|---|---|---|
| ADM-1 | Admin authenticates and is authorized via the Admin role. | M |
| ADM-2 | Admin can manage the **food catalog** (CRUD). | M |
| ADM-3 | Admin can manage the **exercise catalog** (CRUD). | M |
| ADM-4 | Admin can view a **read-only aggregate dashboard** (e.g. total users, popular foods). | S |
| ADM-5 | Admin can **list user accounts** and **activate/deactivate** them (soft delete). | M |
| ADM-6 | Admin **cannot** view an individual user's private meal/workout data. | M |

---

## 7. Non-Functional Requirements

These are intentionally lightweight for a learning MVP; full detail will live in
`docs/02-system/`.

| Area | Requirement |
|---|---|
| **Architecture** | Modular monolith (module-first), layered: Controller → Service → Repository → Model + DTO. |
| **Cost** | Must run comfortably on **free-tier** services (single deployable backend, managed Postgres). |
| **Database** | PostgreSQL as the single source of truth. |
| **Security** | Passwords hashed; JWT auth; role-based authorization; secrets never committed. |
| **Privacy** | Admin/user data boundary enforced server-side (see §4.1). |
| **Maintainability** | Clear module boundaries; interfaces for dependency injection; testable layers. |
| **Extensibility** | Design "paths" for deferred features (roles, cache, score/streak, audit trail) without building them now. |
| **Localization** | Food catalog oriented to Indonesian meals (UI language TBD). |

---

## 8. Assumptions & Constraints

### 8.1 Assumptions
- Food calorie/macro values are **estimates**; the catalog is manually curated for the MVP.
- Workout calorie burn is an **estimate** and inherently approximate.
- Daily targets use standard BMR/TDEE-style formulas (defined in the system design phase).
- Primary food content is Indonesian cuisine.

### 8.2 Constraints
- Built as a learning project; must stay within free-tier infrastructure.
- Backend: Go (Gin + GORM). Mobile: Flutter. Web admin: Next.js. DB: PostgreSQL.
- Team is two "people" (Product Owner + Kiro) covering all roles sequentially.

---

## 9. Roadmap (Post-MVP)

| Phase | Candidate features |
|---|---|
| **Phase 2** | Health **score (0–100)** and **streaks** (gamification); macro **targets** tuned per goal; admin **audit trail**; push notifications/reminders (FCM). |
| **Phase 3** | Custom exercises; meal plans / recommendations; social features; smartwatch / Google Fit integration; barcode / photo food input; additional roles (e.g. Nutritionist/Coach). |

> **Design philosophy:** "Prepare the path, don't pave it." MVP data models should
> accommodate these without large migrations — e.g. clean timestamps for future streaks,
> a role entity for future roles, a cache-friendly summary design for Redis.

---

## 10. Open Questions (to resolve in System Design)

These are deferred to the SA/ERD phase intentionally:

1. **Portion basis for foods** — per serving, per 100 g, or support both? (FOOD-2)
2. **Custom food persistence** — reusable personal foods vs one-off entries. (NUT-7; current intent: reusable.)
3. **Summary computation** — on-the-fly vs precomputed/stored, and where the cache boundary sits. (SUM-6)
4. **Calorie-burn estimation model** — formula per exercise unit type (reps / duration / steps), and whether body weight factors in. (WRK-3)
5. **Daily target formula** — exact BMR/TDEE variant and activity multipliers. (SUM-1)

---

## 11. Disclaimer

cal-lowrisk is a personal tracking and educational tool. All calorie and macronutrient
values are estimates and may not be accurate. The application does not provide medical,
nutritional, or professional health advice and must not be used as a substitute for
consultation with a qualified professional. It is not a medical device.
