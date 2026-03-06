<script>
  import { state } from "../state.svelte.js";
</script>

<div
  class="toast bg-surface-2 border-edge-2 text-bright text-tiny tracking-code pointer-events-none fixed bottom-7 inset-x-0 mx-auto w-fit z-50 rounded border px-4 py-2 transition-transform duration-300"
  class:up={state.toastVisible}
>
  {state.toast}
</div>

<style>
  .toast {
    transform: translateY(80px);
    transition-timing-function: cubic-bezier(0.23, 1, 0.32, 1);
  }
  .toast.up {
    transform: translateY(0);
  }
</style>
