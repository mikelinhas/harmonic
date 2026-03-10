<script>
  import { onMount } from "svelte";
  import { state as appState } from "../state.svelte.js";
  import { doTune, doHarmonize, doRename } from "../api.js";
  import Knob from "./Knob.svelte";
  import Button from "./Button.svelte";

  let nameInput = $state(localStorage.getItem("harmonic_username") || "");
  let renameTimer = null;

  function onNameInput() {
    clearTimeout(renameTimer);
    renameTimer = setTimeout(() => doRename(nameInput), 600);
  }

  const players = $derived(appState.snap?.players ?? []);
const selectedFreq = $derived(
    appState.freqIdx >= 0
      ? appState.scaleItems[appState.freqIdx]
      : appState.selectedExtra >= 0
        ? appState.extraItems[appState.selectedExtra]
        : null,
  );
  const canTransmit = $derived(selectedFreq !== null);
  const shouldPulse = $derived(canTransmit && selectedFreq !== transmittedValue);

  let transmittedValue = $state(null);
  let harmonizeConfirm = $state(false);
  let prevSnapPhase = $state(-1);
  $effect(() => {
    const phase = appState.snap?.phase ?? -1;
    if (phase === 1 && prevSnapPhase !== 1) transmittedValue = null;
    prevSnapPhase = phase;
  });

  function handleHarmonize() {
    const allTuned = players.every(p => p.hasTuned);
    if (!allTuned && !harmonizeConfirm) {
      harmonizeConfirm = true;
      return;
    }
    harmonizeConfirm = false;
    doHarmonize();
  }


  async function handleTransmit() {
    if (!canTransmit) return;
    transmittedValue = selectedFreq;
    await doTune();
  }

  onMount(() => {
    const handleKeydown = (e) => {
      if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA") return;
      if (appState.showSettings) return;
      if (e.key === "t" || e.key === "T") {
        e.preventDefault();
        if (canTransmit) handleTransmit();
      } else if (e.key === "r" || e.key === "R") {
        e.preventDefault();
        handleHarmonize();
      } else if ((e.key === "y" || e.key === "Y") && harmonizeConfirm) {
        e.preventDefault();
        handleHarmonize();
      } else if (e.key === "Escape" && harmonizeConfirm) {
        harmonizeConfirm = false;
      }
    };
    window.addEventListener("keydown", handleKeydown);
    return () => window.removeEventListener("keydown", handleKeydown);
  });
</script>

<div class="flex flex-col">
  <div
    class="grid min-h-115 grid-cols-[190px_1fr_70px] grid-rows-[1fr_auto] max-[600px]:grid-cols-[1fr_56px] max-[600px]:grid-rows-[auto_auto_auto]"
  >
    <!-- Left column -->
    <div
      class="border-edge max-[600px]:border-edge flex flex-col gap-5 border-r px-5 py-6 max-[600px]:col-start-1 max-[600px]:col-end-3 max-[600px]:border-r-0 max-[600px]:border-b"
    >
      <div>
        <label class="text-label tracking-code text-muted mb-1 block" for="operator-name">OPERATOR</label>
        <input
          id="operator-name"
          class="tracking-code text-bright border-edge-2 w-full border-b bg-transparent pb-1 text-sm uppercase outline-none transition-colors duration-200 focus:border-[var(--color-green)]"
          type="text"
          autocomplete="off"
          spellcheck="false"
          maxlength="20"
          bind:value={nameInput}
          oninput={onNameInput}
        />
      </div>
      <div>
        <div class="text-label tracking-code text-muted mb-2">
          LINKED OPERATORS
        </div>
        <div class="flex flex-col gap-1">
          {#each players as p (p.name)}
            <div
              class="plist-item text-tiny tracking-ui flex items-center justify-between"
              class:me={p.name === appState.username}
            >
              <span>{p.name.toUpperCase()}</span>
              <div class="tuned-dot" class:lit={p.hasTuned}></div>
            </div>
          {/each}
        </div>
      </div>
    </div>

    <!-- Center + Right columns rendered by Knob -->
    <Knob />

    <!-- Bottom bar -->
    <div
      class="border-edge col-start-1 col-end-4 flex flex-wrap items-center justify-end gap-3 border-t px-6 py-3.5 max-[600px]:col-end-3"
    >
      <!-- Main actions -->
      <div class="flex flex-wrap items-center gap-2.5">
        <Button variant="ghost" onclick={handleHarmonize}>REVEAL [R]</Button>
        <div class="flex items-stretch gap-2">
          <Button onclick={handleTransmit} disabled={!canTransmit} class={shouldPulse ? 'pulse' : ''}>TRANSMIT [T]</Button>
          <div class="transmitted-box" class:has-value={transmittedValue !== null}>
            {transmittedValue ?? ''}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

{#if harmonizeConfirm}
  <div class="backdrop" role="presentation" onclick={() => harmonizeConfirm = false} onkeydown={null}>
    <div class="modal" role="dialog" onclick={(e) => e.stopPropagation()} onkeydown={null}>
      <div class="modal-header">
        <div class="modal-title">REVEAL</div>
        <button class="close-btn" onclick={() => harmonizeConfirm = false}>✕</button>
      </div>
      <div class="modal-body">
        <p class="modal-message">Not all operators have transmitted a frequency. Reveal anyway?</p>
        <div class="untimed-list">
          {#each players.filter(p => !p.hasTuned) as p (p.name)}
            <div class="untimed-item">
              <div class="tuned-dot"></div>
              <span>{p.name.toUpperCase()}</span>
            </div>
          {/each}
        </div>
      </div>
      <div class="modal-footer">
        <button class="cancel-btn" onclick={() => harmonizeConfirm = false}>CANCEL [ESC]</button>
        <button class="confirm-btn" onclick={handleHarmonize}>YES [Y]</button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* Player list */
  .plist-item {
    color: var(--color-mid);
  }
  .plist-item.me {
    color: var(--color-green);
  }
  .tuned-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: var(--color-surface-3);
    border: 1px solid var(--color-edge-2);
    transition: all 0.3s;
  }
  .tuned-dot.lit {
    background: var(--color-green);
    border-color: var(--color-green);
    box-shadow: 0 0 5px var(--color-green);
  }

  @keyframes transmit-pulse {
    0% { box-shadow: 0 0 0 0 rgba(170, 255, 0, 0.5); }
    25% { box-shadow: 0 0 0 8px rgba(170, 255, 0, 0); }
    100% { box-shadow: 0 0 0 8px rgba(170, 255, 0, 0); }
  }
  :global(button.pulse) {
    animation: transmit-pulse 2.5s ease-out 1s infinite;
  }

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
    max-width: 360px;
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

  .modal-title {
    font-family: var(--font-retro);
    font-size: 0.75rem;
    letter-spacing: 0.2em;
    color: var(--color-green);
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

  .modal-message {
    font-family: var(--font-retro);
    font-size: 0.75rem;
    letter-spacing: 0.08em;
    color: var(--color-mid);
    line-height: 1.6;
    margin: 0;
  }

  .untimed-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .untimed-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-family: var(--font-retro);
    font-size: 0.625rem;
    letter-spacing: 0.15em;
    color: var(--color-muted);
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

  .confirm-btn {
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
  .confirm-btn:hover {
    background: #ccff33;
  }

  .transmitted-box {
    font-family: var(--font-display);
    font-size: 0.75rem;
    font-weight: 900;
    min-width: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #000;
    border: 1px solid var(--color-edge-2);
    border-radius: 4px;
    color: var(--color-green);
    letter-spacing: 0.05em;
    padding: 0 0.5rem;
  }
  .transmitted-box.has-value {
    text-shadow: 0 0 10px rgba(170, 255, 0, 0.5);
  }
</style>
