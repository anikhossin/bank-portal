<script lang="ts">
  import type { Snippet } from 'svelte'
  let { children, class: className = '' }: { children: Snippet; class?: string } = $props()
</script>

<div class="rounded-lg border border-[var(--border)] bg-white shadow-sm {className}">
  {@render children()}
</div>
