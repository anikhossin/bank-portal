<script lang="ts">
  let { label, value }: { label: string; value: string } = $props()
  let copied = $state(false)

  async function copy() {
    await navigator.clipboard.writeText(value)
    copied = true
    setTimeout(() => (copied = false), 1500)
  }
</script>

<div class="space-y-1.5">
  <p class="text-sm font-bold text-[var(--foreground)]">{label}</p>
  <div class="relative">
    <input
      type="text"
      readonly
      {value}
      class="w-full rounded-md border border-[var(--border)] bg-white px-4 py-3.5 pr-20 text-base tracking-wide text-[var(--foreground)] focus:outline-none cursor-default select-all"
    />
    <button
      type="button"
      onclick={copy}
      class="absolute right-3 top-1/2 -translate-y-1/2 text-sm font-semibold text-[#1f2978] hover:text-[#182068] transition-colors"
    >
      {copied ? 'Kopiert!' : 'Kopieren'}
    </button>
  </div>
</div>
