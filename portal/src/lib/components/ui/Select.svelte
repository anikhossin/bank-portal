<script lang="ts">
  let {
    value = $bindable(''),
    id = '',
    options,
  }: { value?: string; id?: string; options: { label: string; value: string }[] } = $props()
</script>

<select
  {id}
  bind:value
  class="w-full rounded-md border border-[var(--border)] bg-white px-3 py-2 text-sm text-[var(--foreground)] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]/30 focus:border-[var(--brand)]"
>
  {#each options as opt (opt.value)}
    <option value={opt.value}>{opt.label}</option>
  {/each}
</select>
