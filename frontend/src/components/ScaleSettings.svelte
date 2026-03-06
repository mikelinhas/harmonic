<script>
  import { onMount } from "svelte";
  import {
    state as gs,
    SCALE_PRESETS,
    saveScaleConfig,
    KNOB_MIN,
  } from "../state.svelte.js";
  import { doSetScale } from "../api.js";

  let { onclose } = $props();

  onMount(() => {
    const handleKeydown = (e) => {
      if (e.key === "Escape") {
        onclose();
        return;
      }
      if (
        (e.key === "a" || e.key === "A") &&
        e.target.tagName !== "INPUT" &&
        e.target.tagName !== "TEXTAREA"
      )
        save();
    };
    window.addEventListener("keydown", handleKeydown);
    return () => window.removeEventListener("keydown", handleKeydown);
  });

  let scaleInput = $state(gs.scaleItems.join(", "));
  let extrasInput = $state(gs.extraItems.join(", "));
  let error = $state("");

  let rangeMin = $state(0);
  let rangeMax = $state(100);
  let rangeStep = $state(5);

  function applyRange() {
    const min = Number(rangeMin);
    const max = Number(rangeMax);
    const step = Math.max(1, Number(rangeStep));
    if (isNaN(min) || isNaN(max) || min >= max) {
      error = "INVALID RANGE";
      return;
    }
    const values = [];
    for (let v = min; v <= max; v += step) values.push(String(v));
    scaleInput = values.join(", ");
    error = "";
  }

  function parseItems(str) {
    return str
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);
  }

  function applyPreset(preset) {
    scaleInput = preset.scale.join(", ");
    extrasInput = preset.extras.join(", ");
    error = "";
  }

  function save() {
    const scale = parseItems(scaleInput);
    const extras = parseItems(extrasInput);
    if (scale.length < 1) {
      error = "SCALE MUST HAVE AT LEAST ONE VALUE";
      return;
    }
    gs.scaleItems = scale;
    gs.extraItems = extras;
    gs.freqIdx = -1;
    gs.selectedExtra = -1;
    gs.knobAngle = KNOB_MIN;
    saveScaleConfig(scale, extras);
    doSetScale(scale, extras);
    onclose();
  }

  function onBackdropClick(e) {
    if (e.target === e.currentTarget) onclose();
  }
</script>

<div
  class="backdrop"
  role="presentation"
  onclick={onBackdropClick}
  onkeydown={null}
>
  <div class="modal">
    <div class="modal-header">
      <div class="text-label tracking-title text-green">CONFIGURATION</div>
      <button class="close-btn" onclick={onclose}>✕</button>
    </div>

    <div class="modal-body">
      <!-- Presets -->
      <div class="field">
        <div class="field-label">PRESETS</div>
        <div class="preset-row">
          {#each SCALE_PRESETS as preset (preset.id)}
            <button class="preset-btn" onclick={() => applyPreset(preset)}>
              {preset.label}
            </button>
          {/each}
        </div>
      </div>

      <div class="field mb-4">
        <div class="range-row">
          <label class="range-field">
            <span class="range-label">MIN</span>
            <input type="number" class="range-input" bind:value={rangeMin} />
          </label>
          <label class="range-field">
            <span class="range-label">MAX</span>
            <input type="number" class="range-input" bind:value={rangeMax} />
          </label>
          <label class="range-field">
            <span class="range-label">STEP</span>
            <input
              type="number"
              class="range-input"
              min="1"
              bind:value={rangeStep}
            />
          </label>
          <button class="preset-btn" onclick={applyRange}>GENERATE</button>
        </div>
      </div>

      <!-- Scale items -->
      <div class="field mb-4">
        <label class="field-label" for="scale-input">SCALE VALUES</label>
        <input
          id="scale-input"
          class="text-input"
          bind:value={scaleInput}
          placeholder="1, 2, 3, 5, 8, 13, 21"
          oninput={() => (error = "")}
        />
        <div class="preview-tags">
          {#each parseItems(scaleInput) as item, i (i)}
            <span class="tag scale-tag">{item}</span>
          {/each}
        </div>
      </div>

      <!-- Extra items -->
      <div class="field">
        <label class="field-label" for="extras-input">EXTRA BUTTONS</label>
        <input
          id="extras-input"
          class="text-input"
          bind:value={extrasInput}
          placeholder="?, ☕"
          oninput={() => (error = "")}
        />
        <div class="preview-tags">
          {#each parseItems(extrasInput) as item, i (i)}
            <span class="tag extra-tag">{item}</span>
          {/each}
        </div>
      </div>

      {#if error}
        <div class="error-msg">{error}</div>
      {/if}
    </div>

    <div class="modal-footer">
      <button class="cancel-btn" onclick={onclose}>CANCEL [ESC]</button>
      <button class="save-btn" onclick={save}>APPLY [A]</button>
    </div>
  </div>
</div>

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.75);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
    padding: 16px;
  }

  .modal {
    background: var(--color-surface);
    border: 1px solid var(--color-edge-2);
    border-radius: 8px;
    width: 100%;
    max-width: 460px;
    box-shadow: 0 24px 80px rgba(0, 0, 0, 0.9);
    display: flex;
    flex-direction: column;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--color-edge);
  }

  .close-btn {
    background: none;
    border: none;
    color: var(--color-muted);
    cursor: pointer;
    font-size: 0.75rem;
    padding: 4px 8px;
    border-radius: 3px;
    transition: color 0.15s;
  }
  .close-btn:hover {
    color: var(--color-bright);
  }

  .modal-body {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 20px;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .field-label {
    font-family: var(--font-retro);
    font-size: 0.75rem;
    letter-spacing: 0.2em;
    color: var(--color-muted);
  }

  .text-input {
    background: var(--color-surface-2);
    border: 1px solid var(--color-edge-2);
    border-radius: 4px;
    color: var(--color-bright);
    font-family: var(--font-retro);
    font-size: 0.75rem;
    letter-spacing: 0.1em;
    padding: 8px 12px;
    outline: none;
    transition: border-color 0.15s;
  }
  .text-input:focus {
    border-color: var(--color-green-dim);
  }
  .text-input::placeholder {
    color: var(--color-muted);
  }

  .preset-row {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .range-row {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: flex-end;
  }

  .range-field {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .range-label {
    font-family: var(--font-retro);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
    color: var(--color-muted);
  }

  .range-input {
    background: var(--color-surface-2);
    border: 1px solid var(--color-edge-2);
    border-radius: 4px;
    color: var(--color-bright);
    font-family: var(--font-retro);
    font-size: 0.75rem;
    letter-spacing: 0.05em;
    padding: 5px 8px;
    outline: none;
    width: 64px;
    transition: border-color 0.15s;
  }
  .range-input:focus {
    border-color: var(--color-green-dim);
  }

  .preset-btn {
    font-family: var(--font-retro);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
    padding: 5px 10px;
    border: 1px solid var(--color-edge-2);
    border-radius: 3px;
    background: transparent;
    color: var(--color-mid);
    cursor: pointer;
    transition: all 0.15s;
  }
  .preset-btn:hover {
    border-color: var(--color-green-dim);
    color: var(--color-green);
  }

  .preview-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 5px;
    min-height: 22px;
  }

  .tag {
    font-family: var(--font-retro);
    font-size: 0.625rem;
    letter-spacing: 0.1em;
    padding: 2px 8px;
    border-radius: 3px;
  }

  .scale-tag {
    background: rgba(170, 255, 0, 0.07);
    border: 1px solid var(--color-green-dim);
    color: var(--color-green);
  }

  .extra-tag {
    background: rgba(255, 149, 0, 0.07);
    border: 1px solid rgba(255, 149, 0, 0.4);
    color: var(--color-orange);
  }

  .error-msg {
    font-family: var(--font-retro);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
    color: var(--color-danger);
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    padding: 14px 20px;
    border-top: 1px solid var(--color-edge);
  }

  .cancel-btn {
    font-family: var(--font-retro);
    font-size: 0.625rem;
    letter-spacing: 0.2em;
    padding: 8px 16px;
    border: 1px solid var(--color-edge-2);
    border-radius: 4px;
    background: transparent;
    color: var(--color-mid);
    cursor: pointer;
    transition: all 0.15s;
  }
  .cancel-btn:hover {
    color: var(--color-bright);
    border-color: var(--color-mid);
  }

  .save-btn {
    font-family: var(--font-retro);
    font-size: 0.625rem;
    letter-spacing: 0.2em;
    padding: 8px 20px;
    border: none;
    border-radius: 4px;
    background: var(--color-green);
    color: #000;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.15s;
  }
  .save-btn:hover {
    background: #ccff33;
    box-shadow: 0 0 16px rgba(170, 255, 0, 0.2);
  }
</style>
