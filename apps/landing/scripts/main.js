// Page bootstrap: dynamic year, and lazy-init the Three.js hero scene with graceful
// fallback (no WebGL / load failure / reduced motion -> the page still works fully).

// Current year in the footer.
const yearEl = document.getElementById("year");
if (yearEl) yearEl.textContent = String(new Date().getFullYear());

// Detect WebGL support before loading Three.js (avoids shipping the lib if unusable).
function hasWebGL() {
  try {
    const c = document.createElement("canvas");
    return !!(window.WebGLRenderingContext && (c.getContext("webgl") || c.getContext("experimental-webgl")));
  } catch (_) {
    return false;
  }
}

const canvas = document.getElementById("hero-canvas");
const note = document.getElementById("webgl-note");

function showFallbackNote() {
  if (note) note.hidden = false;
  if (canvas) canvas.style.display = "none";
}

if (canvas && hasWebGL()) {
  // Only load the 3D scene when the hero is near the viewport (it is, at top — so this
  // mostly guards against wasted work if the markup ever changes).
  const start = () => {
    import("./hero-scene.js")
      .then((m) => m.initHeroScene(canvas))
      .catch((err) => {
        console.warn("Hero 3D scene failed to load:", err);
        showFallbackNote();
      });
  };

  if ("IntersectionObserver" in window) {
    const io = new IntersectionObserver((entries, obs) => {
      if (entries.some((e) => e.isIntersecting)) {
        obs.disconnect();
        start();
      }
    });
    io.observe(canvas);
  } else {
    start();
  }
} else {
  showFallbackNote();
}
