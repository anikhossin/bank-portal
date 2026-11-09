<script lang="ts">
  import type { Snippet } from 'svelte'

  let {
    variant = 'default',
    type = 'button',
    disabled = false,
    onclick,
    children,
  }: {
    variant?: 'default' | 'outline'
    type?: 'button' | 'submit'
    disabled?: boolean
    onclick?: () => void
    children: Snippet
  } = $props()

  const base =
    'inline-flex items-center justify-center gap-2 rounded-md text-sm font-medium px-4 py-2 transition-colors disabled:opacity-50 disabled:pointer-events-none'
  const variants = {
    default: 'bg-[#1e3a5f] text-white hover:bg-[#162c47]',
    outline: 'border border-[var(--border)] bg-white hover:bg-[var(--muted)] text-[var(--foreground)]',
  }
</script>

<button class="{base} {variants[variant]}" {type} {disabled} {onclick}>
  {@render children()}
</button>
