# Callium — Landing Page

Marketing/landing site for **Callium** (the team/company) and its flagship product in
development, **cal-lowrisk**. Built as a lightweight static site with a **Three.js** hero
"balance" scene. Intended to support the Claude for Startups application alongside a
`@callium` email domain.

> cal-lowrisk is in active **R&D** — the page states this clearly. It is a product preview,
> not a live app.

## Stack

- Plain **HTML + CSS + ES modules** — no build step required.
- **Three.js** loaded via an ESM CDN import map (see `index.html`). The hero scene is
  lazy-loaded and degrades gracefully (no WebGL / load failure / `prefers-reduced-motion`
  → the page works fully without 3D).
- Visual tokens from `docs/03-ux/user-flows.md` §11: Inter, primary orange `#FF6B3D`,
  macro colors (protein=blue, carbs=green, fat=amber), warm/clean/calm.

## Structure

```
apps/landing/
├── index.html
├── styles/main.css
├── scripts/
│   ├── main.js          # year, WebGL detection, lazy-init of the hero scene
│   └── hero-scene.js    # Three.js balance scene (ES module)
├── assets/              # favicon.svg, og-image.svg
├── _headers             # Cloudflare Pages security + cache headers
└── README.md
```

## Run locally

No build needed — serve the folder over HTTP (ES modules + import maps need http://, not file://):

```bash
cd apps/landing
python3 -m http.server 5173
# open http://localhost:5173
```

(Any static server works: `npx serve`, `php -S`, etc.)

## Deploy to Cloudflare Pages

1. Push this repo (branch `uat`/`main`).
2. Cloudflare dashboard → **Workers & Pages → Create → Pages → Connect to Git** → select
   `Monarqhs/cal-lowrisk`.
3. Build settings:
   - **Framework preset:** None
   - **Build command:** *(leave empty — static site)*
   - **Build output directory:** `apps/landing`
   - **Root directory:** *(repo root)*
4. Deploy. Cloudflare gives a `*.pages.dev` URL immediately.

### Connect the custom domain (callium.app / callium.com)

1. Register the domain (Cloudflare Registrar is cheapest at cost; `.app` requires HTTPS,
   which Cloudflare provides automatically).
2. Pages project → **Custom domains → Set up a custom domain** → enter `callium.app`
   (and `www`). If the domain is on Cloudflare DNS, records are added automatically.
3. HTTPS is automatic. Done.

> `.app` vs `.com`: `.app` reads modern/techy and forces HTTPS; `.com` is the most
> universally trusted. The page is domain-agnostic — the email `call.help@callium.app`
> in the footer is the only hard-coded reference; update it if you choose `.com`.

## Email domain (`@callium`)

Buying the domain enables `call.help@callium.app`. Cheapest path: keep DNS on Cloudflare
and use **Cloudflare Email Routing** (free) to forward `call.help@callium.app` to a
personal inbox — enough to send/receive for the application. A full mailbox (Google
Workspace / Zoho) can come later.

## TODO / revisit
- Replace the SVG `og-image` with a rendered PNG once the brand is finalized.
- Revise the Three.js scene per feedback (first pass implements the agreed "balance" idea).
- Consider migrating to the planned Next.js app if the landing grows beyond a single page.
