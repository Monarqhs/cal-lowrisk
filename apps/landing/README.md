# Callium — Landing Page

Marketing/landing site for **Callium** (the team/company) and its flagship product in
development, **cal-lowrisk**. Built to support the Claude for Startups application
alongside a `@callium` email domain.

> cal-lowrisk is in active **R&D** — the page states this clearly. It's a product preview,
> not a live app.

## Look & feel

- **Notion-style line-art**, monochrome base (near-black ink on warm off-white) with brand
  color as accent only.
- **Fonts:** **Anton** for titles/display (a free placeholder for **Harse** — swap later via
  the `--font-title` CSS variable), **Poppins** for headings + body.
- Brand accents from `docs/03-ux/user-flows.md` §11: orange `#FF6B3D`; macros protein=blue,
  carbs=green, fat=amber.
- Logo is a **wordmark** ("Callium" in the title typeface) — placeholder until an official
  logo exists.

## Stack

- Plain **HTML + CSS + ES modules**, no build step.
- **Animations:** the official LottieFiles player **`dotlottie-wc`** (web component) via CDN
  — renders both `.lottie` and `.json`. See `<head>` in `index.html`.
- **Line-art icons:** local single-color SVGs in `assets/icons/`, tinted via CSS masks
  (so hover can recolor them to brand orange).
- `scripts/main.js` adds the footer year + a reduced-motion-aware reveal-on-scroll.

## Structure

```
apps/landing/
├── index.html
├── styles/main.css
├── scripts/main.js
├── assets/
│   ├── favicon.svg, og-image.svg
│   ├── icons/        # line-art SVGs (balance, macros, target, bowl, log)
│   └── lottie/       # (optional) self-hosted .lottie/.json animations
├── _headers          # Cloudflare Pages security + cache headers
└── README.md
```

## Swapping the hero animation

The hero uses a `<dotlottie-wc>` element with a **placeholder** `src`. To use the chosen
"Healthy Diet" animation (free under the Lottie Simple License):

- **Option A — hotlink:** on LottieFiles, open the animation → **Asset & Embed** → copy the
  `.lottie` (or `.json`) URL → paste it into the `src` of `<dotlottie-wc class="hero-lottie">`
  in `index.html`.
- **Option B — self-host (preferred for clean licensing):** download the `.lottie`/`.json`
  into `assets/lottie/` and set `src="/assets/lottie/<file>"`. Add attribution to the author
  here if the license requires it.

Both `.lottie` and `.json` work with the same player.

## Run locally

ES modules + the web component need `http://` (not `file://`):

```bash
cd apps/landing
python3 -m http.server 5173       # or: npx serve -l 5173
# open http://localhost:5173
```

## Deploy to Cloudflare Pages

1. Push the repo (`uat`/`main`).
2. Cloudflare → **Workers & Pages → Create → Pages → Connect to Git** → `Monarqhs/cal-lowrisk`.
3. Build settings: **Framework preset:** None · **Build command:** *(empty)* ·
   **Build output directory:** `apps/landing` · **Root directory:** *(repo root)*.
4. Deploy → a `*.pages.dev` URL. HTTPS is automatic.

### Custom domain (callium.app / callium.com)
Pages project → **Custom domains → Set up a custom domain** → add `callium.app` (+ `www`).
If DNS is on Cloudflare, records are added automatically; HTTPS is automatic.

### Email (`@callium`)
Keep DNS on Cloudflare and use **Email Routing** (free) to forward
`call.help@callium.app` to a personal inbox. Full mailbox (Google Workspace / Zoho) later.

## TODO
- Swap Anton → **Harse** once the licensed font file is available (self-host via `@font-face`,
  update `--font-title`).
- Swap the placeholder hero `.lottie` for the chosen "Healthy Diet" animation.
- Replace the SVG `og-image` with a rendered PNG once branding is finalized.
