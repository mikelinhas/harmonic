<script>
  import { onMount } from "svelte";
  import { state, roomCode } from "../state.svelte.js";
  import { doJoin } from "../api.js";
  import Button from "./Button.svelte";

  onMount(() => {
    const handleKeydown = (e) => {
      if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA") return;
      if (e.key === "s" || e.key === "S") {
        e.preventDefault();
        doJoin();
      }
    };
    window.addEventListener("keydown", handleKeydown);
    return () => window.removeEventListener("keydown", handleKeydown);
  });
</script>

<div
  class="flex flex-1 flex-col items-center justify-center gap-7 px-5 py-16 text-center"
>
  <div class="text-tiny tracking-caps text-muted">
    ROOM // <span class="text-green">{roomCode}</span>
  </div>
  <div
    class="font-display tracking-title text-bright font-bold"
    style="font-size:clamp(18px,4vw,26px)"
  >
    CONNECT TO SESSION
  </div>
  <div class="flex w-full max-w-xs flex-col gap-1.5">
    <label
      class="text-label tracking-caps text-muted text-left"
      for="username-input"
    >
      OPERATOR CALLSIGN
    </label>
    <input
      class="field-input text-bright font-retro tracking-code placeholder:text-surface-3 w-full border-0 border-b bg-transparent py-2 text-sm uppercase transition-colors duration-200 outline-none"
      id="username-input"
      type="text"
      placeholder="ENTER NAME"
      autocomplete="off"
      spellcheck="false"
      maxlength="20"
      bind:value={state.username}
      onkeydown={(e) => e.key === "Enter" && doJoin()}
    />
  </div>
  <Button onclick={doJoin}>⬡ START [S]</Button>
</div>

<style>
  .field-input {
    border-color: var(--color-green-dim);
  }
  .field-input:focus {
    border-color: var(--color-green);
  }
</style>
