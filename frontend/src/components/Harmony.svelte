<script>
  import { onMount } from "svelte";
  import { state, roomCode } from "../state.svelte.js";
  import { doReset } from "../api.js";
  import Button from "./Button.svelte";
  import SigBar from "./SigBar.svelte";

  const players = $derived(state.snap?.players ?? []);
  const numericFreqs = $derived(
    players.map((p) => parseFloat(p.frequency)).filter((v) => !isNaN(v)),
  );
  const avgFreq = $derived(
    numericFreqs.length
      ? (numericFreqs.reduce((a, b) => a + b, 0) / numericFreqs.length).toFixed(
          1,
        )
      : "—",
  );
  const harmonyLabel = $derived.by(() => {
    const n = numericFreqs;
    if (n.length < 2) return "—";
    const mean = n.reduce((a, b) => a + b, 0) / n.length;
    const stddev = Math.sqrt(
      n.reduce((a, b) => a + (b - mean) ** 2, 0) / n.length,
    );
    return stddev < 1.5 ? "HIGH" : stddev < 4 ? "MED" : "LOW";
  });
  const harmonyClass = $derived(
    { HIGH: "text-green", MED: "text-orange", LOW: "text-danger" }[
      harmonyLabel
    ] ?? "text-green",
  );

  onMount(() => {
    const handleKeydown = (e) => {
      if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA")
        return;
      if (state.showSettings) return;
      if (e.key === "n" || e.key === "N") {
        e.preventDefault();
        doReset();
      }
    };
    window.addEventListener("keydown", handleKeydown);
    return () => window.removeEventListener("keydown", handleKeydown);
  });
  const vuBarData = $derived.by(() => {
    return players.map((p) => {
      if (!p.frequency) {
        return { name: p.name, frequency: "—", hPct: 0, isMe: p.name === state.username };
      }
      const idx = state.scaleItems.indexOf(p.frequency);
      let hPct;
      if (idx >= 0) {
        hPct = state.scaleItems.length === 1 ? 100 : (idx / (state.scaleItems.length - 1)) * 90 + 10;
      } else {
        hPct = 0; // extra item (e.g. "?", "☕")
      }
      return { name: p.name, frequency: p.frequency, hPct, isMe: p.name === state.username };
    });
  });
</script>

<div class="flex flex-col">
  <!-- VU bars -->
  <div class="border-edge border-b px-7 pt-12 pb-5">
    <div class="flex flex-wrap items-end justify-center gap-4">
      {#each vuBarData as bar (bar.name)}
        <div
          class="flex w-14 flex-col items-center gap-1.5"
        >
          <div
            class="text-green font-display min-w-10 px-1.5 py-0.5 text-center text-xs font-bold"
          >
            {bar.frequency}
          </div>
          <SigBar height={bar.hPct} barClass="h-36" />
          <div
            class="text-micro tracking-code text-mid min-h-12 pt-2 text-center"
            class:me={bar.isMe}
          >
            {bar.name.toUpperCase()}
          </div>
        </div>
      {/each}
    </div>
  </div>

  <!-- Actions -->
  <div class="flex flex-wrap justify-center gap-3 px-7 py-5">
    <Button class="px-9 text-sm" onclick={doReset}>NEW SESSION [N]</Button>
  </div>
</div>

<style>
  .me {
    color: var(--color-green);
  }
</style>
