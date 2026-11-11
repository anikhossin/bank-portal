<script lang="ts">
  import { createCase, listCases, expireCase, type CaseSummary } from './lib/api'

  // --- auth ---
  let authed = $state(false)
  let loginUser = $state('')
  let loginPass = $state('')
  let loginError = $state('')

  $effect(() => {
    try { authed = sessionStorage.getItem('admin:auth') === '1' } catch {}
  })

  function handleLogin(e: Event) {
    e.preventDefault()
    if (loginUser === 'admin' && loginPass === 'Paternal4-Going5-Smartness1') {
      try { sessionStorage.setItem('admin:auth', '1') } catch {}
      authed = true
    } else {
      loginError = 'Invalid username or password.'
    }
  }

  function logout() {
    try { sessionStorage.removeItem('admin:auth') } catch {}
    authed = false
  }

  // --- cases ---
  let caseId = $state('')
  let iban = $state('')
  let aesReference = $state('')

  let loading = $state(false)
  let error = $state('')

  let cases = $state<CaseSummary[]>([])
  let loadingList = $state(true)
  let copiedId = $state<string | null>(null)
  let expiringId = $state<string | null>(null)

  async function loadList() {
    try { cases = await listCases() } catch {}
    loadingList = false
  }

  $effect(() => { loadList() })

  async function handleSubmit(e: Event) {
    e.preventDefault()
    error = ''
    loading = true
    try {
      await createCase({ case_id: caseId, iban, aes_reference: aesReference })
      caseId = ''
      iban = ''
      aesReference = ''
      await loadList()
    } catch (err) {
      error = err instanceof Error ? err.message : 'Something went wrong'
    } finally {
      loading = false
    }
  }

  async function copyLink(c: CaseSummary) {
    await navigator.clipboard.writeText(c.portal_link)
    copiedId = c.case_id
    setTimeout(() => (copiedId = null), 1500)
  }

  async function handleExpire(c: CaseSummary) {
    expiringId = c.case_id
    try {
      await expireCase(c.case_id)
      await loadList()
    } catch {}
    expiringId = null
  }

  function formatDate(s: string) {
    return new Date(s).toLocaleString('de-DE', { dateStyle: 'short', timeStyle: 'short' })
  }
</script>

<div class="min-h-screen flex flex-col" style="background:var(--bg);color:var(--foreground)">

{#if !authed}
  <!-- Login screen -->
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <p class="text-xl font-semibold tracking-wide">Dagual Bank</p>
        <p class="text-xs mt-1" style="color:var(--muted-foreground)">Admin Portal</p>
      </div>
      <form onsubmit={handleLogin} class="space-y-4 border p-6" style="border-color:var(--border);background:var(--surface)">
        <div class="space-y-1.5">
          <label class="block text-xs font-medium" style="color:var(--muted-foreground)">Username</label>
          <input
            bind:value={loginUser}
            required
            autocomplete="username"
            class="w-full px-3 py-2 text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>
        <div class="space-y-1.5">
          <label class="block text-xs font-medium" style="color:var(--muted-foreground)">Password</label>
          <input
            type="password"
            bind:value={loginPass}
            required
            autocomplete="current-password"
            class="w-full px-3 py-2 text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>
        {#if loginError}
          <p class="text-xs text-red-400">{loginError}</p>
        {/if}
        <button
          type="submit"
          class="w-full py-2.5 text-sm font-semibold text-white"
          style="background:var(--brand)"
        >Sign In</button>
      </form>
    </div>
  </div>
{:else}

  <!-- Header -->
  <header class="border-b px-6 py-4 flex items-center gap-3" style="border-color:var(--border);background:var(--surface)">
    <span class="text-lg font-semibold tracking-wide" style="color:var(--foreground)">Dagual Bank</span>
    <span class="text-xs px-2 py-0.5 font-medium" style="background:var(--brand);color:#fff">ADMIN</span>
    <button onclick={logout} class="ml-auto text-xs" style="color:var(--muted-foreground)">Sign out</button>
  </header>

  <div class="flex flex-1 overflow-hidden">

    <!-- Sidebar -->
    <aside class="w-72 shrink-0 border-r flex flex-col" style="border-color:var(--border);background:var(--surface)">
      <div class="px-5 py-5 border-b" style="border-color:var(--border)">
        <p class="text-xs font-semibold uppercase tracking-widest" style="color:var(--muted-foreground)">Create Case</p>
      </div>

      <form class="flex-1 px-5 py-5 space-y-4 overflow-y-auto" onsubmit={handleSubmit}>

        <div class="space-y-1.5">
          <label class="block text-xs font-medium" style="color:var(--muted-foreground)">Case ID <span class="opacity-50">(optional)</span></label>
          <input
            bind:value={caseId}
            placeholder="CASE-2026-001"
            class="w-full px-3 py-2 text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        <div class="space-y-1.5">
          <label class="block text-xs font-medium" style="color:var(--muted-foreground)">IBAN <span class="text-red-400">*</span></label>
          <input
            bind:value={iban}
            required
            placeholder="DE89 3704 0044 0532 0130 00"
            class="w-full px-3 py-2 text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        <div class="space-y-1.5">
          <label class="block text-xs font-medium" style="color:var(--muted-foreground)">AES Reference <span class="text-red-400">*</span></label>
          <input
            bind:value={aesReference}
            required
            placeholder="AES-REF-00123"
            class="w-full px-3 py-2 text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        {#if error}
          <p class="text-xs text-red-400">{error}</p>
        {/if}

        <button
          type="submit"
          disabled={loading}
          class="w-full py-2.5 text-sm font-semibold text-white transition-colors disabled:opacity-50"
          style="background:var(--brand)"
          onmouseenter={e => (e.currentTarget as HTMLElement).style.background = 'var(--brand-hover)'}
          onmouseleave={e => (e.currentTarget as HTMLElement).style.background = 'var(--brand)'}
        >
          {loading ? 'Creating...' : 'Create Case'}
        </button>

      </form>
    </aside>

    <!-- Main: cases list -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <div class="px-6 py-5 border-b flex items-center justify-between" style="border-color:var(--border)">
        <p class="text-xs font-semibold uppercase tracking-widest" style="color:var(--muted-foreground)">All Cases</p>
        <button
          onclick={loadList}
          class="text-xs px-3 py-1.5 border transition-colors"
          style="border-color:var(--border);color:var(--muted-foreground)"
        >Refresh</button>
      </div>

      <div class="flex-1 overflow-auto">
        {#if loadingList}
          <p class="px-6 py-8 text-sm" style="color:var(--muted-foreground)">Loading...</p>
        {:else if cases.length === 0}
          <p class="px-6 py-8 text-sm" style="color:var(--muted-foreground)">No cases yet.</p>
        {:else}
          <table class="w-full text-sm border-collapse">
            <thead>
              <tr class="text-left border-b" style="border-color:var(--border)">
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Case ID</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">IBAN</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">AES Ref</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Created</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Status</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each cases as c}
                <tr class="border-b transition-colors" style="border-color:var(--border)" onmouseenter={e => (e.currentTarget as HTMLElement).style.background = 'var(--surface)'} onmouseleave={e => (e.currentTarget as HTMLElement).style.background = ''}>
                  <td class="px-4 py-3 font-mono text-xs">{c.case_id}</td>
                  <td class="px-4 py-3 font-mono text-xs" style="color:var(--muted-foreground)">{c.iban}</td>
                  <td class="px-4 py-3 font-mono text-xs" style="color:var(--muted-foreground)">{c.aes_reference}</td>
                  <td class="px-4 py-3 text-xs" style="color:var(--muted-foreground)">{formatDate(c.created_at)}</td>
                  <td class="px-4 py-3">
                    {#if c.expired}
                      <span class="text-xs px-2 py-0.5 font-medium" style="background:#3a1a1a;color:#f87171">Expired</span>
                    {:else}
                      <span class="text-xs px-2 py-0.5 font-medium" style="background:#0f2a1a;color:#4ade80">Active</span>
                    {/if}
                  </td>
                  <td class="px-4 py-3">
                    <div class="flex items-center gap-2">
                      <button
                        onclick={() => copyLink(c)}
                        class="text-xs px-2.5 py-1 border transition-colors"
                        style="border-color:var(--border);color:var(--muted-foreground)"
                      >{copiedId === c.case_id ? 'Copied!' : 'Copy Link'}</button>
                      {#if !c.expired}
                        <button
                          onclick={() => handleExpire(c)}
                          disabled={expiringId === c.case_id}
                          class="text-xs px-2.5 py-1 border transition-colors disabled:opacity-40"
                          style="border-color:#7f1d1d;color:#f87171"
                        >{expiringId === c.case_id ? '...' : 'Expire'}</button>
                      {/if}
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </div>
    </main>

  </div>

{/if}
</div>
