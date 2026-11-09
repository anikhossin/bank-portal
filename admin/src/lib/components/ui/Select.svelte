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
  class="w-full rounded-md border border-[var(--border)] bg-white px-3 py-2 text-sm text-[var(--foreground)] focus:outline-none focus:ring-2 focus:ring-[#1e3a5f]/30 focus:border-[#1e3a5f]"
>
  {#each options as opt (opt.value)}
    <option value={opt.value}>{opt.label}</option>
  {/each}
</select>
