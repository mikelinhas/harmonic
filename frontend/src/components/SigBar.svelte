<script>
  let { height = 0, barClass = "flex-1" } = $props();
</script>

<div
  class="border-edge relative flex w-7 flex-col-reverse overflow-hidden rounded border bg-[#0a0a0a] {barClass}"
>
  <div
    class="fill w-full"
    style={`height: ${height}%; background-size: 100% ${height > 0 ? (100 / height * 100).toFixed(1) : '10000'}%`}
  ></div>
  <div class="grid-overlay pointer-events-none absolute inset-0"></div>
</div>
<style>
  .fill {
    background-image: linear-gradient(
      to top,
      #00dd44 0%,
      #aaff00 30%,
      #ffdd00 60%,
      #ff8800 80%,
      #ff2200 100%
    );
    background-position: bottom center;
    background-repeat: no-repeat;
    transition: height 0.45s cubic-bezier(0.23, 1, 0.32, 1);
    box-shadow: 0 0 10px rgba(170, 255, 0, 0.3);
  }
  .grid-overlay {
    background: repeating-linear-gradient(
      180deg,
      transparent,
      transparent 5px,
      rgba(0, 0, 0, 0.35) 5px,
      rgba(0, 0, 0, 0.35) 6px
    );
  }
</style>
