import { state, roomCode, KNOB_MIN, saveScaleConfig } from "./state.svelte.js";

window.addEventListener("beforeunload", () => {
  if (state.username && roomCode) {
    navigator.sendBeacon(
      `${location.origin}/rooms/${roomCode}/leave`,
      new Blob([JSON.stringify({ username: state.username })], { type: "application/json" })
    );
  }
});

let sse = null;
let toastTimer = null;
let prevPhase = -1;

export function showToast(msg) {
  state.toast = msg;
  state.toastVisible = true;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => {
    state.toastVisible = false;
  }, 2600);
}

async function api(method, path, body) {
  const opts = { method, headers: { "Content-Type": "application/json" } };
  if (body) opts.body = JSON.stringify(body);
  const res = await fetch(`${location.origin}/rooms/${roomCode}/${path}`, opts);
  if (!res.ok)
    throw Object.assign(new Error(await res.text()), { status: res.status });
  if (res.status === 204) return null;
  return res.json();
}

function scaleChanged(a, b) {
  return !a || a.length !== b.length || a.some((v, i) => v !== b[i]);
}

function onSnap(s) {
  const phaseChanged = s.phase !== prevPhase;
  prevPhase = s.phase;
  state.snap = s;
  if (phaseChanged && s.phase === 1) {
    state.hasTuned = false;
    state.transmittedFreq = null;
    state.freqIdx = -1;
    state.selectedExtra = -1;
    state.knobAngle = KNOB_MIN;
  }
  if (phaseChanged && s.phase === 2) state.roundNum++;

  // Sync scale from room state.
  if (
    s.scale?.length &&
    (scaleChanged(s.scale, state.scaleItems) ||
      scaleChanged(s.extras, state.extraItems))
  ) {
    state.scaleItems = s.scale;
    state.extraItems = s.extras ?? [];
    state.freqIdx = -1;
    state.selectedExtra = -1;
    state.knobAngle = KNOB_MIN;
    saveScaleConfig(s.scale, s.extras ?? []);
  }
}

export function connectSSE() {
  if (sse) {
    sse.close();
    sse = null;
  }
  sse = new EventSource(`${location.origin}/rooms/${roomCode}/events`);
  sse.onopen = () => {
    state.connected = true;
  };
  sse.onmessage = (e) => {
    try {
      onSnap(JSON.parse(e.data));
    } catch {}
  };
  sse.onerror = () => {
    state.connected = false;
    sse.close();
    setTimeout(() => connectSSE(), 3000);
  };
}

// Try to reconnect using a stored username without registering a new slot.
// Returns true if the player was found in the room.
export async function tryReconnect() {
  const val = state.username;
  if (!val || !roomCode) return false;
  try {
    await api("POST", "reconnect", { username: val });
    connectSSE();
    return true;
  } catch {
    return false;
  }
}

export async function doJoin() {
  const val = state.username.trim().toUpperCase();
  if (!val) {
    showToast("CALLSIGN REQUIRED");
    return;
  }
  state.username = val;
  localStorage.setItem("harmonic_username", val);
  try {
    await api("POST", "join", { username: val });
    connectSSE();
  } catch (err) {
    if (err.status === 409) {
      showToast("NAME TAKEN");
    } else {
      showToast("CONNECTION FAILED");
    }
  }
}

export async function doRename(newName) {
  const oldName = state.username;
  const val = newName.trim().toUpperCase();
  if (!val || val === oldName) return;
  try {
    await api("POST", "rename", { oldUsername: oldName, newUsername: val });
    state.username = val;
    localStorage.setItem("harmonic_username", val);
  } catch (err) {
    if (err.status === 409) {
      showToast("NAME TAKEN");
    } else {
      showToast("RENAME FAILED");
    }
  }
}

export async function doTune() {
  let freq = null;
  if (state.freqIdx >= 0) freq = state.scaleItems[state.freqIdx];
  else if (state.selectedExtra >= 0)
    freq = state.extraItems[state.selectedExtra];

  if (!freq) {
    showToast("SELECT A FREQUENCY");
    return;
  }
  try {
    await api("POST", "tune", { username: state.username, frequency: freq });
    state.hasTuned = true;
    state.transmittedFreq = freq;
  } catch {
    showToast("TRANSMISSION FAILED");
  }
}

export async function doStart() {
  try {
    await api("POST", "start", {});
    state.hasTuned = false;
  } catch {
    showToast("ERROR");
  }
}

export async function doHarmonize() {
  try {
    await api("POST", "harmonize", {});
  } catch {
    showToast("ERROR");
  }
}

export async function doSetScale(scale, extras) {
  try {
    await api("POST", "scale", { scale, extras });
  } catch {
    showToast("UPDATE FAILED");
  }
}
