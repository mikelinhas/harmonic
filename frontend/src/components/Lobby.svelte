<script>
  import { onMount } from "svelte";
  import { state, roomCode } from "../state.svelte.js";
  import { doStart } from "../api.js";
  import Button from "./Button.svelte";

  const players = $derived(state.snap?.players ?? []);
  onMount(() => {
    const handleKeydown = (e) => {
      if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA") return;
      if (e.key === "n" || e.key === "N") {
        e.preventDefault();
        doStart();
      }
    };
    window.addEventListener("keydown", handleKeydown);
    return () => window.removeEventListener("keydown", handleKeydown);
  });
</script>

<div class="flex flex-1 flex-col items-center gap-6 px-5 py-10">
  <div class="flex flex-col items-center gap-5 py-5">
    <div class="text-micro tracking-caps text-muted">
      ROOM // <span class="text-green">{roomCode}</span>
    </div>
    <div
      class="font-display text-bright tracking-caps font-bold"
      style="font-size:0.75ren"
    >
      STANDBY MODE
    </div>
    <div class="flex w-full max-w-xs flex-col gap-1.5">
      {#each players as p (p.name)}
        <div
          class="player-row bg-surface-2 border-edge tracking-code text-mid flex items-center justify-between rounded border px-3 py-2 text-xs"
          class:me={p.name === state.username}
        >
          <span>{p.name.toUpperCase()}</span>
          {#if p.name === state.username}
            <span class="text-muted text-micro">(YOU)</span>
          {/if}
        </div>
      {/each}
    </div>
    <Button onclick={doStart}>NEW SESSION [N]</Button>
  </div>
</div>

<style>
  .player-row.me {
    color: var(--color-green);
    border-color: var(--color-green-dim);
  }
</style>
