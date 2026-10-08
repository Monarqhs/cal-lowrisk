# User Flows & Screen Inventory — cal-lowrisk

| Field | Value |
|---|---|
| **Project** | cal-lowrisk |
| **Document** | UX — user flows, sitemap, screen inventory, low-fi wireframes |
| **Version** | 0.1 (Draft) |
| **Status** | Draft — pending review |
| **Author** | Product Owner + Kiro (UX) |
| **Phase** | UX (follows the SA phase: `erd.md`, `architecture.md`, `deployment.md`) |
| **Last updated** | 2026-10-08 |

> **Process note:** The UX phase defines *how the product is used* — journeys, screens, and
> interactions — derived from the BRD and the data model. It is **docs-as-code**: flows as
> Mermaid, wireframes as text/ASCII (low-fidelity, intentionally). High-fidelity visual
> design happens next **in Figma** (see §9), using the **screen inventory** here as the
> design brief. No pixel mockups live in this repo; link Figma from `docs/03-ux/` later.

---

## 1. UX principles (from the BRD)

| # | Principle | Source |
|---|---|---|
| P1 | **Clarity over data dumps.** The core question is "am I balanced vs my goal today?" — answer it first, details after. | BRD §2.2 G2, Persona A |
| P2 | **Fast logging.** Logging a meal/workout is the most frequent action; keep it few-taps. | G1 |
| P3 | **Local-first content.** Indonesian foods and household serving sizes feel native. | §3.1, FOOD-5 |
| P4 | **Personalized, not generic.** Everything frames against the user's goal and target. | G4 |
| P5 | **Privacy is visible.** User data is theirs; admins never see individual logs. | §4.1 |
| P6 | **Two distinct clients.** User = Flutter (mobile, log + see balance). Admin = Next.js (web, curate + manage). Don't blur them. | §4 |

---

## 2. Sitemap

### 2.1 User app (Flutter)
```mermaid
flowchart TD
    Splash --> AuthCheck{Logged in?}
    AuthCheck -->|no| Login
    AuthCheck -->|yes| Home
    Login --> Register
    Register --> Onboarding["Onboarding: profile + goal"]
    Onboarding --> Home

    Home["Home / Today's Balance"] --> LogMeal["Log Meal"]
    Home --> LogWorkout["Log Workout"]
    Home --> Summary["Summaries (week/month)"]
    Home --> Profile["Profile & Goal"]

    LogMeal --> FoodSearch["Food search (catalog + my custom)"]
    FoodSearch --> CustomFood["Create custom food"]
    FoodSearch --> PortionPick["Pick portion (serving/grams)"]
    LogWorkout --> ExerciseSearch["Exercise search (catalog)"]
    ExerciseSearch --> AmountEntry["Enter amount"]
```

### 2.2 Admin web (Next.js)
```mermaid
flowchart TD
    ALogin["Admin Login"] --> Dashboard["Dashboard (aggregate)"]
    Dashboard --> FoodCat["Food Catalog (CRUD)"]
    Dashboard --> ExerciseCat["Exercise Catalog (CRUD)"]
    Dashboard --> Users["User Accounts (list, activate/deactivate)"]

    FoodCat --> FoodForm["Food create/edit + servings"]
    ExerciseCat --> ExerciseForm["Exercise create/edit"]
```

> **Privacy boundary (P5):** the admin sitemap has **no path** to any individual user's
> meal/workout logs. "User Accounts" shows account info + status only.

---

## 3. Core user journey — the loop (BRD success indicator)

The MVP's heartbeat: **log a meal → log a workout → see today's balance vs goal.**

```mermaid
flowchart LR
    Start([Open app]) --> Home["Home: today's net vs target"]
    Home --> LM["Log a meal<br/>(pick food + portion)"]
    LM --> Home2["Home updates:<br/>calories in ↑"]
    Home2 --> LW["Log a workout<br/>(pick exercise + amount)"]
    LW --> Home3["Home updates:<br/>calories out ↑"]
    Home3 --> See["See status:<br/>deficit / balanced / surplus<br/>+ macro breakdown"]
    See --> Done([Done])
```

---

## 4. Detailed flows (User app)

### 4.1 Registration → Onboarding (goal + profile)
First-run sets up everything the target calculation needs (`erd.md` §4.3; USR-1..4).

```mermaid
flowchart TD
    A[Register: email + password] --> B[Pick goal: lose / maintain / gain]
    B --> C["Body profile: weight, height, age, sex, activity level"]
    C --> D{Valid?}
    D -->|no| C
    D -->|yes| E["Compute target (Mifflin-St Jeor → TDEE → goal)"]
    E --> F[Home with target set]
```
- **Why goal first, then body data:** goal frames the whole app (P4); body data is the
  numeric input to the target.
- Target is computed on-the-fly (never stored) — see `erd.md` §5 / `architecture.md` §5.4.

### 4.2 Log a meal (most frequent — keep it fast, P2)
```mermaid
flowchart TD
    A["Tap Log Meal"] --> B["Choose meal type: breakfast / lunch / dinner / snack"]
    B --> C["Search food"]
    C --> D{Found?}
    D -->|"yes, catalog"| E["Select food"]
    D -->|"yes, my custom"| E
    D -->|no| F["Create custom food: name + calories, macros optional"]
    F --> E
    E --> G{Has servings?}
    G -->|yes| H["Pick serving + quantity"]
    G -->|no| I["Enter grams"]
    H --> J["Confirm, entry saved, grams resolved"]
    I --> J
    J --> K["Home updates: calories in + macros"]
```
- Custom foods are **reusable** and appear in future searches (NUT-7; `erd.md` §4.4).
- Portion always resolves to **grams** under the hood (NUT-3), even if picked by serving.

### 4.3 Log a workout
```mermaid
flowchart TD
    A["Tap Log Workout"] --> B["Search exercise catalog"]
    B --> C["Select exercise"]
    C --> D{Unit type?}
    D -->|duration| E["Enter minutes"]
    D -->|reps| F["Enter reps"]
    D -->|"distance / steps"| G["Enter steps or distance"]
    E --> H["Estimate burn: MET x weight x time"]
    F --> H
    G --> H
    H --> I["Confirm, entry saved"]
    I --> J["Home updates: calories out"]
```
- No custom exercises in MVP (EXE-5). Burn uses MET + body weight (`erd.md` §4.8).

### 4.4 View summaries (daily / weekly / monthly)
```mermaid
flowchart LR
    Home["Daily (default on Home)"] --> Week["Weekly: avg net, on-track?"]
    Week --> Month["Monthly: trend"]
    Home -.-> Macro["Macro breakdown: protein/carbs/fat vs target"]
```
- Daily is the default view (P1); weekly/monthly are tabs/segments.
- All computed on-the-fly from logs (SUM-2..5).

---

## 5. Detailed flows (Admin web)

### 5.1 Catalog curation (food / exercise)
```mermaid
flowchart TD
    A[Admin login] --> B[Dashboard]
    B --> C[Food Catalog list]
    C --> D{Action}
    D -->|create| E[New food: name, per-100g calories+macros]
    D -->|edit| F[Edit food]
    D -->|delete| G[Soft-delete food]
    E --> H[Optionally add servings: label → grams]
    F --> H
    H --> I[Saved → appears in user search]
```
- Exercise catalog mirrors this (name, unit_type, MET). Deletes are soft (`erd.md`).

### 5.2 User account management (within the privacy boundary)
```mermaid
flowchart TD
    A[Users list: email, status, role] --> B{Action}
    B -->|activate| C[Clear deleted_at]
    B -->|deactivate| D[Set deleted_at = soft delete]
    A -.->|NOT ALLOWED| X["View a user's meal/workout logs ❌"]
```
- Deactivate = soft delete (ADM-5). **No screen or API exposes a user's private logs**
  (ADM-6 / P5) — this is a hard architectural invariant, not just a hidden button.

---

## 6. Screen inventory (the Figma design brief)

The concrete list of screens to design. Each note what it must show/do. **This is what you
design in Figma (§9).**

### 6.1 User app (Flutter) — mobile
| # | Screen | Must contain | Key actions |
|---|---|---|---|
| U1 | **Splash / Auth check** | logo | route to Login or Home |
| U2 | **Login** | email, password, link to Register | log in |
| U3 | **Register** | email, password, confirm | continue to onboarding |
| U4 | **Onboarding — Goal** | 3 goal options (lose/maintain/gain), short explainer | select goal |
| U5 | **Onboarding — Profile** | weight, height, age, sex, activity level | compute target, go Home |
| U6 | **Home / Today** | **big balance** (in vs out vs target), status chip (deficit/balanced/surplus), macro mini-bars, quick-add buttons (meal, workout), today's entries list | tap to log; open summaries/profile |
| U7 | **Log Meal — meal type** | breakfast/lunch/dinner/snack picker | choose, go to search |
| U8 | **Food Search** | search box, results (catalog + "my foods" tagged), "Create custom food" CTA | select or create |
| U9 | **Create Custom Food** | name, calories (required), macros (optional) | save (private, reusable) |
| U10 | **Portion Picker** | serving options (if any) + quantity, OR grams input; live calorie/macro preview | confirm |
| U11 | **Log Workout — search** | search box, catalog results | select exercise |
| U12 | **Amount Entry** | input adapts to unit_type (minutes/reps/steps), live burn estimate | confirm |
| U13 | **Summaries** | segmented Daily/Weekly/Monthly; net, target, status, macro breakdown, simple trend | switch range |
| U14 | **Profile & Goal** | current body data + goal, edit; shows current target | edit (recomputes target) |

### 6.2 Admin web (Next.js) — desktop
| # | Screen | Must contain | Key actions |
|---|---|---|---|
| A1 | **Admin Login** | email, password | log in (Admin role only) |
| A2 | **Dashboard** | aggregate stats (total users, popular foods) — **no individual data** | navigate |
| A3 | **Food Catalog** | table (name, calories/100g, macros), search, create button | CRUD |
| A4 | **Food Form** | fields + **servings sub-list** (label → grams) | save |
| A5 | **Exercise Catalog** | table (name, unit_type, MET), create button | CRUD |
| A6 | **Exercise Form** | name, unit_type, MET | save |
| A7 | **User Accounts** | table (email, status, role) — **no logs** | activate/deactivate |

---

## 7. Low-fi wireframes (text) — the two most important screens

Intentionally rough — just layout intent for the Figma brief.

### 7.1 Home / Today (U6) — the screen that answers "am I balanced?" (P1)
```text
┌─────────────────────────────┐
│  cal-lowrisk        [profile]│
│                             │
│      TODAY — Oct 8          │
│   ┌───────────────────┐     │
│   │   NET: -320 kcal   │     │  ← big, the headline
│   │  [ DEFICIT ✓ ]     │     │  ← status chip (color-coded)
│   └───────────────────┘     │
│   in 1850  out 420  tgt 2150│  ← supporting numbers
│                             │
│  Macros                     │
│  P ▓▓▓▓░░  C ▓▓▓░░░  F ▓▓░░░ │  ← vs target
│                             │
│  [ + Log Meal ] [ + Workout]│  ← fast actions (P2)
│                             │
│  Today's entries            │
│  • Nasi goreng   520 kcal   │
│  • Lari 20 min  -180 kcal   │
│  ...                        │
│  [ Daily | Weekly | Month ] │  ← to summaries (U13)
└─────────────────────────────┘
```

### 7.2 Portion Picker (U10) — where logging feels local (P3)
```text
┌─────────────────────────────┐
│ ‹ Back      Nasi Goreng      │
│                             │
│ Portion                     │
│  ( ) 1 piring  (≈ 250 g)    │  ← household servings (P3)
│  ( ) 1 centong (≈ 100 g)    │
│  (•) Grams: [ 180 ] g       │  ← fallback
│                             │
│  Quantity: [ − ] 1 [ + ]    │
│                             │
│  Preview                    │
│   520 kcal  P 12  C 78  F 16│  ← live calc
│                             │
│         [  Add to log  ]    │
└─────────────────────────────┘
```

---

## 8. Cross-cutting interaction notes

- **Cold start (free-tier):** first request after idle may be slow (`deployment.md` §8.3).
  Every data screen needs a **loading state** and a **retry-friendly error state** —
  never a dead white screen. Design these states in Figma, not just the happy path.
- **Empty states:** Home with no entries ("Log your first meal"), empty search, no custom
  foods yet — design them; they're the true first-run experience.
- **Validation inline:** profile numbers, custom-food calories, workout amount — show
  errors next to the field (P1 clarity).
- **Status color system:** deficit / balanced / surplus need a consistent, accessible color
  language (define as tokens in Figma → become theme constants in code).
- **Role separation:** user flows never show admin UI and vice-versa; they are different
  apps (P6).

---

## 9. Handoff to Figma — "ready-for-design" checklist

The screen inventory (§6) is your design brief. To make the later **design → code** step
(via the Figma Power, in the Dev phase) as clean as possible:

- [ ] One Figma file (or two: **Mobile** = U1–U14, **Admin** = A1–A7).
- [ ] **Name each frame** to match the inventory (e.g. `Home / Today`, `Portion Picker`,
      `Food Catalog`) — frame names become screen/component names in code.
- [ ] Use **Auto Layout** — translates far more accurately to Flutter/CSS layout.
- [ ] Define **styles/variables** for the status colors, spacing, and typography — these
      become design tokens / theme constants in code.
- [ ] Build repeated elements as **components** (status chip, entry row, macro bar, table row).
- [ ] Design the **non-happy states** called out in §8 (loading, empty, error).
- [ ] Ensure the file is **accessible to the Figma Power** when we reach the Dev phase.

> You can design manually or use Figma's **Make** to draft — either way the integration
> reads your final design. The cleaner the structure, the less "AI slop" and the higher
> the design-code parity.

---

## 10. Open items for review

1. **Navigation pattern (mobile):** bottom tab bar (Home / Summaries / Profile) vs
   Home-centric with push navigation. (Leaning bottom tabs for a tracker.)
2. **Weekly/monthly visualization:** simple numbers vs a lightweight chart — how much
   visualization for MVP without heavy chart deps.
3. **Meal-type requirement:** is `meal_type` mandatory per entry, or default to a
   time-of-day guess the user can change?
4. **Custom food editing:** can a user edit/delete their custom foods from a dedicated
   "My Foods" screen (beyond inline during logging)?
5. **Admin dashboard metrics (A2):** exact aggregate widgets for MVP (ties to deployment
   §8 and ADM-4), staying within the privacy boundary.
