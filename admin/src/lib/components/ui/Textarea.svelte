<script lang="ts">
  let {
    value = $bindable(''),
    placeholder = '',
    required = false,
    id = '',
    rows = 4,
  }: { value?: string; placeholder?: string; required?: boolean; id?: string; rows?: number } = $props()
</script>

<textarea
  {id}
  {placeholder}
  {required}
  {rows}
  bind:value
  class="w-full rounded-md border border-[var(--border)] bg-white px-3 py-2 text-sm text-[var(--foreground)] placeholder:text-[var(--muted-foreground)] focus:outline-none focus:ring-2 focus:ring-[#1e3a5f]/30 focus:border-[#1e3a5f] resize-none"
></textarea>
