<script>
  import { onMount } from "svelte";
  import { state as appState, roomCode } from "./state.svelte.js";
  import { tryReconnect } from "./api.js";
  import Join from "./components/Join.svelte";
  import Tuning from "./components/Tuning.svelte";
  import Harmony from "./components/Harmony.svelte";
  import Toast from "./components/Toast.svelte";
  import ScaleSettings from "./components/ScaleSettings.svelte";
  import Button from "./components/Button.svelte";

  const phase = $derived(appState.snap?.phase ?? -1);

  // Room picker state (used only when at root /)
  let roomInput = $state("");

  function goToRoom() {
    const code = roomInput.trim().toUpperCase();
    if (!code) return;
    location.href = `/rooms/${encodeURIComponent(code)}`;
  }

  onMount(async () => {
    // Try to silently reconnect if we have a stored username for this room.
    if (roomCode && appState.username) {
      await tryReconnect();
    }

    const handleKeydown = (e) => {
      if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA")
        return;
      if ((e.key === "x" || e.key === "X") && !appState.showSettings) {
        e.preventDefault();
        appState.showSettings = true;
      }
    };
    window.addEventListener("keydown", handleKeydown);
    return () => window.removeEventListener("keydown", handleKeydown);
  });
</script>

<header class="mb-2 flex w-full items-center justify-between">
  <div class="flex flex-col gap-0.5">
    <span class="text-label tracking-caps text-muted uppercase"
      >INPUT TERMINAL</span
    >
    <span
      class="hdr-title font-display tracking-caps text-green font-bold"
      style="font-size:clamp(16px,3vw,22px)"
    >
     VOTING MACHINE 
    </span>
  </div>
  {#if roomCode}
    <Button variant="sm" onclick={() => (appState.showSettings = true)}
      >⚙ SETTINGS [X]</Button
    >
  {/if}
</header>

<main
  class="bg-surface border-edge-2 flex w-full flex-col overflow-hidden rounded-lg border shadow-[0_24px_80px_rgba(0,0,0,0.9)]"
>
  {#if !roomCode}
    <!-- Room picker -->
    <div
      class="flex flex-1 flex-col items-center justify-center gap-7 px-5 py-16 text-center"
    >
      <div
        class="font-display tracking-title text-bright font-bold"
        style="font-size:clamp(18px,4vw,26px)"
      >
        ENTER ROOM
      </div>
      <div class="flex w-full max-w-xs flex-col gap-1.5">
        <label
          class="text-label tracking-caps text-muted text-left"
          for="room-input"
        >
          ROOM CODE
        </label>
        <input
          class="text-bright font-retro tracking-code placeholder:text-muted border-green-dim focus:border-green w-full border-0 border-b bg-transparent py-2 text-sm uppercase transition-colors duration-200 outline-none"
          id="room-input"
          type="text"
          placeholder="WALRUS"
          autocomplete="off"
          spellcheck="false"
          maxlength="20"
          bind:value={roomInput}
          onkeydown={(e) => e.key === "Enter" && goToRoom()}
        />
      </div>
      <Button onclick={goToRoom}>⬡ JOIN ROOM</Button>
    </div>
  {:else if !appState.snap}
    <Join />
  {:else if phase === 1}
    <Tuning />
  {:else if phase === 2}
    <Harmony />
  {/if}

  <div
    class="border-edge flex flex-wrap items-center justify-between gap-2 border-t bg-[#111] px-7 py-2"
  >
    <div class="text-label tracking-code text-muted flex items-center gap-1.5">
      {#if appState.connected}
        <div class="bg-green size-1.5 rounded-full shadow-[0_0_4px_#aaff00]"></div>
        POWER: ON
      {:else}
        <div class="size-1.5 rounded-full border border-current opacity-40"></div>
        POWER: OFF
      {/if}
    </div>
    {#if roomCode}
      <div
        class="text-label tracking-code text-muted flex items-center gap-1.5"
      >
        ROOM: <span class="text-green">{roomCode}</span>
      </div>
    {/if}
  </div>
</main>

<Toast />

{#if appState.showSettings}
  <ScaleSettings onclose={() => (appState.showSettings = false)} />
{/if}

<style>
  .hdr-title {
    text-shadow: 0 0 24px rgba(170, 255, 0, 0.5);
  }
</style>
