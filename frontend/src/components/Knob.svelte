<script>
  import { onMount } from "svelte";
  import { state, KNOB_MIN, KNOB_MAX } from "../state.svelte.js";
  import SigBar from "./SigBar.svelte";
  import Button from "./Button.svelte";

  // Ticks live inside the physical rotation range, leaving dead zones at each end
  const TICK_MIN = -110;
  const TICK_MAX = 110;

  const knobStep = $derived(
    (TICK_MAX - TICK_MIN) / Math.max(1, state.scaleItems.length - 1),
  );

  const selectedFreq = $derived(
    state.freqIdx >= 0
      ? state.scaleItems[state.freqIdx]
      : state.selectedExtra >= 0
        ? state.extraItems[state.selectedExtra]
        : "—",
  );

  const sigHeight = $derived.by(() => {
    if (state.selectedExtra >= 0) return 0;
    if (state.freqIdx < 0) return 0;
    if (state.scaleItems.length === 1) return 100;
    return (state.freqIdx / (state.scaleItems.length - 1)) * 90 + 10;
  });

  const knobTicks = $derived.by(() => {
    const len = state.scaleItems.length;
    const dense = len > 20;
    const stride = dense ? Math.ceil(len / 25) : 1;
    return state.scaleItems.map((_, i) => {
      const angle = TICK_MIN + i * knobStep;
      const rad = ((angle - 90) * Math.PI) / 180;
      const x = 90 + 82 * Math.cos(rad);
      const y = 90 + 82 * Math.sin(rad);
      const visible = !dense || i % stride === 0 || i === len - 1;
      const on = i <= state.freqIdx;
      return {
        style: `left:${x - 1}px;top:${y - 6}px;transform:rotate(${angle}deg)`,
        on,
        visible,
      };
    });
  });

  function selectExtra(i) {
    state.selectedExtra = i;
    state.freqIdx = -1;
    state.knobAngle = KNOB_MIN;
  }

  onMount(() => {
    const setFreqIdx = (idx) => {
      const len = state.scaleItems.length;
      idx = Math.max(0, Math.min(len - 1, idx));
      if (idx === state.freqIdx) return;
      state.freqIdx = idx;
      state.selectedExtra = -1;
      state.knobAngle =
        TICK_MIN + idx * ((TICK_MAX - TICK_MIN) / Math.max(1, len - 1));
      const fc = document.getElementById("freq-card");
      if (fc) {
        fc.classList.remove("flash");
        void fc.offsetWidth;
        fc.classList.add("flash");
      }
    };

    let dragging = false,
      startY = 0,
      startAngle = 0;

    const angleToIdx = (a) => {
      const step =
        (TICK_MAX - TICK_MIN) / Math.max(1, state.scaleItems.length - 1);
      return Math.round((a - TICK_MIN) / step);
    };

    const knobBody = document.getElementById("knob-body");
    const knobOuter = document.getElementById("knob-outer");

    knobBody.addEventListener("mousedown", (e) => {
      dragging = true;
      startY = e.clientY;
      startAngle = state.knobAngle;
      e.preventDefault();
    });
    knobBody.addEventListener(
      "touchstart",
      (e) => {
        dragging = true;
        startY = e.touches[0].clientY;
        startAngle = state.knobAngle;
        e.preventDefault();
      },
      { passive: false },
    );

    window.addEventListener("mousemove", (e) => {
      if (!dragging) return;
      const a = Math.max(
        KNOB_MIN,
        Math.min(KNOB_MAX, startAngle + (startY - e.clientY) * 2.2),
      );
      state.knobAngle = a;
      setFreqIdx(angleToIdx(a));
    });
    window.addEventListener(
      "touchmove",
      (e) => {
        if (!dragging) return;
        const a = Math.max(
          KNOB_MIN,
          Math.min(
            KNOB_MAX,
            startAngle + (startY - e.touches[0].clientY) * 2.2,
          ),
        );
        state.knobAngle = a;
        setFreqIdx(angleToIdx(a));
      },
      { passive: true },
    );

    window.addEventListener("mouseup", () => {
      dragging = false;
    });
    window.addEventListener("touchend", () => {
      dragging = false;
    });

    const handleKeydown = (e) => {
      if (e.target.tagName === "INPUT" || e.target.tagName === "TEXTAREA") return;
      if (state.showSettings) return;
      if (e.key === "ArrowUp") {
        e.preventDefault();
        setFreqIdx(state.freqIdx < 0 ? 0 : state.freqIdx + 1);
      } else if (e.key === "ArrowDown") {
        e.preventDefault();
        setFreqIdx(state.freqIdx < 0 ? 0 : state.freqIdx - 1);
      } else if (e.key === "ArrowRight") {
        e.preventDefault();
        if (state.extraItems.length === 0) return;
        selectExtra(state.selectedExtra < 0 ? 0 : (state.selectedExtra + 1) % state.extraItems.length);
      } else if (e.key === "ArrowLeft") {
        e.preventDefault();
        if (state.extraItems.length === 0) return;
        const prev = state.selectedExtra <= 0 ? state.extraItems.length - 1 : state.selectedExtra - 1;
        selectExtra(prev);
      }
    };
    window.addEventListener("keydown", handleKeydown);

    knobOuter.addEventListener(
      "wheel",
      (e) => {
        e.preventDefault();
        if (state.freqIdx < 0 && state.selectedExtra >= 0) {
          // re-engage scale from the nearest end when scrolling from an extra
          setFreqIdx(e.deltaY > 0 ? 0 : state.scaleItems.length - 1);
        } else {
          setFreqIdx(state.freqIdx + (e.deltaY > 0 ? -1 : 1));
        }
      },
      { passive: false },
    );

    return () => window.removeEventListener("keydown", handleKeydown);
  });
</script>

<!-- Center column -->
<div
  class="max-[600px]:border-edge flex flex-col items-center gap-4 px-5 py-6 max-[600px]:col-start-1 max-[600px]:col-end-2 max-[600px]:row-start-2 max-[600px]:border-r"
>
  <div
    class="freq-display relative flex w-full max-w-xs flex-col items-center justify-center gap-2 overflow-hidden rounded border border-[#1e1e1e] bg-black max-[600px]:max-w-full"
    style="aspect-ratio:4/2.5"
  >
    <div
      class="freq-card font-display text-readout text-green z-10 min-w-16 text-center leading-none font-black"
      id="freq-card"
      class:dim={state.freqIdx < 0 && state.selectedExtra < 0}
      class:untransmitted={(state.freqIdx >= 0 || state.selectedExtra >= 0) && selectedFreq !== state.transmittedFreq}
    >
      {selectedFreq}
    </div>
  </div>

  <div class="flex flex-col items-center gap-2 select-none">
    <div class="knob-outer relative h-45 w-45" id="knob-outer">
      {#each knobTicks as tick, i (i)}
        {#if tick.visible}
          <div
            class="tick absolute h-1.5 w-0.5 rounded-sm bg-[#2a2a2a]"
            class:on={tick.on}
            style={tick.style}
          ></div>
        {/if}
      {/each}
      <div
        class="knob-body absolute inset-4 cursor-grab rounded-full active:cursor-grabbing"
        id="knob-body"
      >
        <div
          class="knob-rim pointer-events-none absolute inset-4 rounded-full border border-[#2c2c2c]"
        ></div>
        <div
          class="knob-ind pointer-events-none absolute bottom-2.5"
          style={`transform: rotate(${state.knobAngle - 180}deg)`}
        ></div>
      </div>
    </div>
    {#if state.extraItems.length > 0}
      <div class="flex flex-wrap justify-center gap-1.5">
        {#each state.extraItems as extra, i (extra)}
          <Button
            variant="sm"
            active={state.selectedExtra === i}
            onclick={() => selectExtra(i)}
          >{extra}</Button>
        {/each}
      </div>
    {/if}
  </div>
</div>

<!-- Right column — signal bar -->
<div
  class="border-edge flex flex-col items-center gap-2 border-l py-5 max-[600px]:col-start-2 max-[600px]:col-end-3 max-[600px]:row-start-2"
>
  <SigBar height={sigHeight} />
</div>

<style>
  /* Frequency display scanline + vignette */
  .freq-display::before {
    content: "";
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 1;
    background: repeating-linear-gradient(
      0deg,
      rgba(170, 255, 0, 0.025) 0px,
      rgba(170, 255, 0, 0.025) 1px,
      transparent 1px,
      transparent 3px
    );
  }
  .freq-display::after {
    content: "";
    position: absolute;
    inset: 0;
    pointer-events: none;
    z-index: 2;
    background: radial-gradient(
      ellipse at center,
      transparent 60%,
      rgba(0, 0, 0, 0.6)
    );
  }
  .freq-card {
    text-shadow:
      0 0 30px var(--color-green),
      0 0 60px rgba(170, 255, 0, 0.4);
    transition: all 0.15s;
  }
  @keyframes cardFlash {
    0% {
      text-shadow:
        0 0 60px var(--color-green),
        0 0 120px var(--color-green);
    }
    100% {
      text-shadow:
        0 0 30px var(--color-green),
        0 0 60px rgba(170, 255, 0, 0.4);
    }
  }
  .freq-card.dim {
    color: #222;
    text-shadow: none;
  }
  .freq-card.untransmitted {
    color: #4a5c28;
    text-shadow: 0 0 12px rgba(100, 130, 40, 0.25);
  }

  /* Knob body */
  .knob-body {
    background: radial-gradient(circle at 32% 28%, #3c3c3c, #141414 70%);
    box-shadow:
      0 8px 24px rgba(0, 0, 0, 0.9),
      0 2px 6px rgba(0, 0, 0, 0.9),
      inset 0 1px 0 rgba(255, 255, 255, 0.07),
      inset 0 -3px 6px rgba(0, 0, 0, 0.6);
  }
  .knob-ind {
    position: absolute;
    bottom: 10px;
    left: calc(50% - 2px);
    width: 4px;
    height: 14px;
    background: var(--color-green);
    border-radius: 2px;
    box-shadow: 0 0 8px var(--color-green);
    transform-origin: 2px -51px;
    transition: transform 0.12s cubic-bezier(0.23, 1, 0.32, 1);
  }
  .tick.on {
    background: var(--color-green);
    box-shadow: 0 0 4px var(--color-green);
  }


</style>
