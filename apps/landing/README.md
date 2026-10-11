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
- **Hero visual:** a static illustration (`assets/hero-illustration.png`). Its own
  background is `#f4f1ec`, which is the page's `--paper` token — so the image blends into
  the page with no visible frame. If you swap in an image with a different background,
  either match `--paper` to it or re-add a frame/blend in `.hero-frame` / `.hero-illustration`.
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
│   ├── hero-illustration.png   # your hero image (hero-illustration.svg = fallback)
│   └── icons/        # line-art SVGs (balance, macros, target, bowl, log)
├── _headers          # Cloudflare Pages security + cache headers
└── README.md
```

## Swapping the hero illustration

The hero shows `assets/hero-illustration.png`. Replace that file with your own image
(keep the filename, or update the `src` in `index.html`). A line-art placeholder
`hero-illustration.svg` loads automatically if the PNG is missing (via `onerror`).

The page's `--paper` token is `#f4f1ec` — the current illustration's own background — so
the image blends in with no visible frame. If your replacement has a different background,
either update `--paper` to match it, or re-introduce a frame/blend on `.hero-frame` /
`.hero-illustration` in `main.css`.

> Use an image you have the rights to (your own render, or a free/licensed asset).

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

## Design discipline

This page was passed against the [nutlope/hallmark](https://github.com/nutlope/hallmark)
anti-AI-slop rule-set. Named tells deliberately avoided: the AI nav (wordmark-left /
links-center / CTA-right), the 3–4-col icon-tile feature grid, per-section eyebrows, the
4-column link footer, fade-up-on-every-section, pure #000/#fff, gradient/aurora heroes,
italic headers. Colours are OKLCH tokens (tinted neutrals, one warm accent); motion is a
single CSS hero entrance; punctuation uses em-dashes and curly quotes. The CSS top carries
the Hallmark macrostructure stamp + a pre-emit self-critique.

## Assets (optimized)

All raster assets are compressed WebP/PNG, sized for how they're actually displayed:

| File | Use | Size |
|---|---|---|
| `hero-illustration.webp` | hero (shown ≤720px, 1440px for retina) | ~46 KB |
| `hero-illustration.svg` | line-art fallback if the webp fails | ~1 KB |
| `Cal-Logo-2.webp` | navbar wordmark (16:9, CSS-cropped) | ~6 KB |
| `Cal-Logo.webp` | square 1:1 logo (spare, for future use) | ~6 KB |
| `favicon-32.png` / `favicon-180.png` | tab icon / apple-touch icon | ~1 / 28 KB |
| `og-image.svg` | social share card | ~1 KB |

The heavy source files (multi-MB JPEGs the logos/hero were exported as) were replaced by
these. To regenerate from new sources, resize + re-encode with Pillow/`cwebp` to similar
dimensions.

## TODO
- Swap Anton → **Harse** once the licensed font file is available (self-host via `@font-face`,
  update `--font-display`).
- Replace the SVG `og-image` with a rendered PNG once branding is finalized.
