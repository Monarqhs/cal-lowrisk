// Hero "balance" scene — a Three.js gimmick that mirrors the product idea:
// two glowing orbs (calories IN vs OUT) gently orbit a central target point, with
// macro-colored particles (protein=blue, carbs=green, fat=amber) drifting around.
// Lightweight, pointer-reactive, and degrades gracefully (caller handles no-WebGL).

import * as THREE from "three";

const BRAND = {
  orange: 0xff6b3d,
  protein: 0x3b82f6,
  carbs: 0x22c55e,
  fat: 0xf59e0b,
};

export function initHeroScene(canvas) {
  const prefersReduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const renderer = new THREE.WebGLRenderer({
    canvas,
    antialias: true,
    alpha: true,
    powerPreference: "high-performance",
  });
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

  const scene = new THREE.Scene();
  const camera = new THREE.PerspectiveCamera(50, 1, 0.1, 100);
  camera.position.set(0, 0, 9);

  // Lights
  scene.add(new THREE.AmbientLight(0xffffff, 0.65));
  const key = new THREE.DirectionalLight(0xffffff, 1.1);
  key.position.set(4, 6, 8);
  scene.add(key);
  const warm = new THREE.PointLight(BRAND.orange, 1.4, 40);
  warm.position.set(-5, -2, 6);
  scene.add(warm);

  // Central "target" — a soft wireframe sphere representing the goal
  const target = new THREE.Mesh(
    new THREE.IcosahedronGeometry(1.15, 2),
    new THREE.MeshStandardMaterial({
      color: BRAND.orange,
      roughness: 0.35,
      metalness: 0.1,
      wireframe: true,
      transparent: true,
      opacity: 0.55,
    })
  );
  scene.add(target);

  // Two orbiting orbs: calories IN (green) and OUT (amber)
  const orbMat = (c) =>
    new THREE.MeshStandardMaterial({ color: c, roughness: 0.25, metalness: 0.2, emissive: c, emissiveIntensity: 0.25 });
  const orbIn = new THREE.Mesh(new THREE.SphereGeometry(0.52, 32, 32), orbMat(BRAND.carbs));
  const orbOut = new THREE.Mesh(new THREE.SphereGeometry(0.52, 32, 32), orbMat(BRAND.fat));
  scene.add(orbIn, orbOut);

  // Macro particles
  const COUNT = 420;
  const positions = new Float32Array(COUNT * 3);
  const colors = new Float32Array(COUNT * 3);
  const macroColors = [BRAND.protein, BRAND.carbs, BRAND.fat].map((c) => new THREE.Color(c));
  for (let i = 0; i < COUNT; i++) {
    const r = 2.6 + Math.random() * 3.4;
    const theta = Math.random() * Math.PI * 2;
    const phi = Math.acos(2 * Math.random() - 1);
    positions[i * 3] = r * Math.sin(phi) * Math.cos(theta);
    positions[i * 3 + 1] = r * Math.sin(phi) * Math.sin(theta) * 0.6;
    positions[i * 3 + 2] = r * Math.cos(phi);
    const c = macroColors[i % 3];
    colors[i * 3] = c.r;
    colors[i * 3 + 1] = c.g;
    colors[i * 3 + 2] = c.b;
  }
  const pGeo = new THREE.BufferGeometry();
  pGeo.setAttribute("position", new THREE.BufferAttribute(positions, 3));
  pGeo.setAttribute("color", new THREE.BufferAttribute(colors, 3));
  const particles = new THREE.Points(
    pGeo,
    new THREE.PointsMaterial({ size: 0.07, vertexColors: true, transparent: true, opacity: 0.9, depthWrite: false })
  );
  scene.add(particles);

  // Group we tilt toward the pointer
  const world = new THREE.Group();
  world.add(target, orbIn, orbOut, particles);
  scene.add(world);

  // Pointer parallax
  const pointer = { x: 0, y: 0, tx: 0, ty: 0 };
  function onPointerMove(e) {
    const nx = (e.clientX / window.innerWidth) * 2 - 1;
    const ny = (e.clientY / window.innerHeight) * 2 - 1;
    pointer.tx = nx * 0.35;
    pointer.ty = -ny * 0.25;
  }
  window.addEventListener("pointermove", onPointerMove, { passive: true });

  // Resize to canvas box
  function resize() {
    const w = canvas.clientWidth || window.innerWidth;
    const h = canvas.clientHeight || window.innerHeight;
    renderer.setSize(w, h, false);
    camera.aspect = w / h;
    camera.updateProjectionMatrix();
  }
  const ro = new ResizeObserver(resize);
  ro.observe(canvas);
  resize();

  // Animation
  const clock = new THREE.Clock();
  let running = true;

  function frame() {
    if (!running) return;
    const t = clock.getElapsedTime();
    const speed = prefersReduced ? 0 : 1;

    // orbit the two orbs on opposite phases (in vs out, seeking balance)
    const R = 2.3;
    orbIn.position.set(Math.cos(t * 0.6 * speed) * R, Math.sin(t * 0.6 * speed) * R * 0.5, Math.sin(t * 0.4 * speed) * 0.8);
    orbOut.position.set(Math.cos(t * 0.6 * speed + Math.PI) * R, Math.sin(t * 0.6 * speed + Math.PI) * R * 0.5, Math.cos(t * 0.4 * speed) * 0.8);

    target.rotation.y = t * 0.15 * speed;
    target.rotation.x = t * 0.08 * speed;
    particles.rotation.y = t * 0.03 * speed;

    // ease the whole scene toward the pointer
    pointer.x += (pointer.tx - pointer.x) * 0.05;
    pointer.y += (pointer.ty - pointer.y) * 0.05;
    world.rotation.y = pointer.x;
    world.rotation.x = pointer.y;

    renderer.render(scene, camera);
    requestAnimationFrame(frame);
  }
  requestAnimationFrame(frame);

  // Pause when tab/section not visible (saves battery/CPU)
  document.addEventListener("visibilitychange", () => {
    running = !document.hidden;
    if (running) { clock.start(); requestAnimationFrame(frame); }
  });

  // Cleanup handle (not strictly needed for a single-page landing)
  return () => {
    running = false;
    ro.disconnect();
    window.removeEventListener("pointermove", onPointerMove);
    renderer.dispose();
    pGeo.dispose();
  };
}
