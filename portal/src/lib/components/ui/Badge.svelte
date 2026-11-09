<script lang="ts">
  let { label }: { label: string } = $props()
</script>

<span
  class="inline-flex items-center rounded-full bg-[var(--brand)]/10 px-2.5 py-0.5 text-xs font-medium text-[var(--brand)]"
>
  {label}
</span>
