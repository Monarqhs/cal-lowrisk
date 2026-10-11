// Page bootstrap: just the footer year.
// Motion is intentionally minimal — a single CSS-driven hero entrance, no
// scroll-triggered fades (anti-slop: "animate-on-scroll on everything" is a tell).
const yearEl = document.getElementById("year");
if (yearEl) yearEl.textContent = String(new Date().getFullYear());
