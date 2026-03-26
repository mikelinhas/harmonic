// ─── Constants ────────────────────────────────────────────────────────────────
export const KNOB_MIN = -135;
export const KNOB_MAX = 135;
export const TOTAL_SEGS = 20;

export const SCALE_PRESETS = [
  { id: "fibonacci", label: "FIBONACCI", scale: ["1", "2", "3", "5", "8", "13", "21"], extras: ["?", "☕"] },
  { id: "tshirt",    label: "T-SHIRT",   scale: ["XS", "S", "M", "L", "XL"],            extras: ["?"] },
  { id: "powers",    label: "POWERS ×2", scale: ["1", "2", "4", "8", "16", "32"],        extras: ["?", "∞"] },
  { id: "range10",  label: "1–10",     scale: Array.from({ length: 10 }, (_, i) => String(i+1)), extras: ["?"] },
];

const _parts = location.pathname.split("/").filter(Boolean);
// roomCode is null when at the root picker page (no /rooms/{code} in path)
export const roomCode =
  _parts.length >= 2 && _parts[0] === "rooms"
    ? decodeURIComponent(_parts[1]).toUpperCase()
    : null;

// ─── Scale config persistence ─────────────────────────────────────────────────
function loadScaleConfig() {
  try {
    const s = localStorage.getItem("harmonic_scale_cfg");
    if (s) return JSON.parse(s);
  } catch {}
  return { scale: ["1", "2", "3", "5", "8", "13", "21"], extras: ["?", "☕"] };
}

export function saveScaleConfig(scale, extras) {
  localStorage.setItem("harmonic_scale_cfg", JSON.stringify({ scale, extras }));
}

const cfg = loadScaleConfig();

// ─── Reactive state (single source of truth) ─────────────────────────────────
export const state = $state({
  snap: null,
  username: localStorage.getItem("harmonic_username") || "",
  freqIdx: -1,
  selectedExtra: -1,
  knobAngle: KNOB_MIN,
  hasTuned: false,
  transmittedFreq: null,
  roundNum: 1,
  connected: false,
  inRoom: false,
  toast: "",
  toastVisible: false,
  scaleItems: cfg.scale,
  extraItems: cfg.extras,
  showSettings: false,
});
