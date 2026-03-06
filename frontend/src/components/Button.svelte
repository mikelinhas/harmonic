<script>
  let {
    variant = "primary",
    onclick,
    disabled,
    active = false,
    class: className = "",
    children,
  } = $props();
</script>

<button
  class="font-retro inline-flex cursor-pointer items-center justify-center gap-2 rounded whitespace-nowrap uppercase transition-all duration-150
    {variant === 'primary'
    ? 'tracking-caps bg-green disabled:bg-surface-3 disabled:text-muted border-0 px-6 py-3 text-xs font-bold text-black hover:bg-[#ccff33] hover:shadow-[0_0_20px_rgba(170,255,0,0.18)] active:scale-[.97] disabled:cursor-not-allowed disabled:shadow-none'
    : variant === 'sm'
    ? `tracking-[0.2em] border bg-transparent px-3 py-1.5 text-[9px] ${active ? 'border-[var(--color-green)] text-[var(--color-green)] bg-[rgba(170,255,0,0.06)] shadow-[0_0_8px_rgba(170,255,0,0.15)]' : 'border-[var(--color-edge-2)] text-[var(--color-muted)] hover:border-[var(--color-green-dim)] hover:text-[var(--color-green)]'}`
    : 'tracking-caps text-mid border-edge-2 hover:border-green-dim hover:text-green border bg-transparent px-6 py-3 text-xs'}
    {className}"
  {onclick}
  {disabled}
>
  {@render children?.()}
</button>
