<script lang="ts">
  let {
    value = $bindable(''),
    placeholder = '',
    required = false,
    id = '',
  }: { value?: string; placeholder?: string; required?: boolean; id?: string } = $props()
</script>

<input
  {id}
  {placeholder}
  {required}
  bind:value
  class="w-full rounded-md border border-[var(--border)] bg-white px-3 py-2 text-sm text-[var(--foreground)] placeholder:text-[var(--muted-foreground)] focus:outline-none focus:ring-2 focus:ring-[var(--brand)]/30 focus:border-[var(--brand)]"
/>
