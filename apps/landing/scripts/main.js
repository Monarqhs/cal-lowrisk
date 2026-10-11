// Page bootstrap: dynamic year + a lightweight reveal-on-scroll for engagement.
// The hero animation is handled by the <dotlottie-wc> web component (loaded in <head>).

// Footer year
const yearEl = document.getElementById("year");
if (yearEl) yearEl.textContent = String(new Date().getFullYear());

// Reveal-on-scroll: fade/slide elements in as they enter the viewport.
// Respects reduced-motion (we simply show everything, no animation).
const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
const revealTargets = document.querySelectorAll(".card, .step, .section-title, .section-lead, .team-points li");

if (reduceMotion || !("IntersectionObserver" in window)) {
  revealTargets.forEach((el) => el.classList.add("is-in"));
} else {
  revealTargets.forEach((el) => el.classList.add("reveal"));
  const io = new IntersectionObserver(
    (entries, obs) => {
      entries.forEach((e) => {
        if (e.isIntersecting) {
          e.target.classList.add("is-in");
          obs.unobserve(e.target);
        }
      });
    },
    { threshold: 0.12, rootMargin: "0px 0px -40px 0px" }
  );
  revealTargets.forEach((el) => io.observe(el));
}
