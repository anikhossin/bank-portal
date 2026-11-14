<script lang="ts">
  import { createCase, listCases, expireCase, deleteCase, type CaseSummary } from './lib/api'

  // --- auth ---
  let authed = $state(false)
  try { authed = sessionStorage.getItem('admin:auth') === '1' } catch {}
  let loginUser = $state('')
  let loginPass = $state('')
  let loginError = $state('')

  function handleLogin(e: Event) {
    e.preventDefault()
    if (loginUser === 'admin' && loginPass === 'Paternal4-Going5-Smartness1') {
      try { sessionStorage.setItem('admin:auth', '1') } catch {}
      authed = true
    } else {
      loginError = 'Ungültiger Benutzername oder Passwort.'
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
  let verwendungszweck = $state('')
  let betrag = $state('')

  let loading = $state(false)
  let error = $state('')

  let cases = $state<CaseSummary[]>([])
  let loadingList = $state(true)
  let copiedId = $state<string | null>(null)
  let expiringId = $state<string | null>(null)
  let deletingId = $state<string | null>(null)
  let panel = $state<'create' | 'list'>('list')

  async function loadList() {
    try { cases = await listCases() } catch {}
    loadingList = false
  }

  loadList()

  async function handleSubmit(e: Event) {
    e.preventDefault()
    error = ''
    loading = true
    try {
      await createCase({ case_id: caseId, iban, aes_reference: aesReference, verwendungszweck, betrag })
      caseId = ''
      iban = ''
      aesReference = ''
      verwendungszweck = ''
      betrag = ''
      await loadList()
      panel = 'list'
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

  async function handleDelete(c: CaseSummary) {
    deletingId = c.case_id
    try {
      await deleteCase(c.case_id)
      cases = cases.filter(x => x.case_id !== c.case_id)
    } catch {}
    deletingId = null
  }

  function formatDate(s: string) {
    return new Date(s).toLocaleString('de-DE', { dateStyle: 'short', timeStyle: 'short' })
  }

  function formatBetrag() {
    const raw = betrag.trim().replace(/\./g, '').replace(',', '.')
    const num = parseFloat(raw)
    if (!isNaN(num)) betrag = num.toLocaleString('de-DE', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  }
</script>

<div class="min-h-screen lg:h-dvh flex flex-col overflow-x-hidden" style="background:var(--bg);color:var(--foreground)">

{#if !authed}
  <!-- Login screen -->
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <p class="text-xl font-semibold tracking-wide">Deutsche Bank</p>
        <p class="text-xs mt-1" style="color:var(--muted-foreground)">Admin Portal</p>
      </div>
      <form onsubmit={handleLogin} class="space-y-4 border p-6" style="border-color:var(--border);background:var(--surface)">
        <div class="space-y-1.5">
          <label for="login-user" class="block text-xs font-medium" style="color:var(--muted-foreground)">Benutzername</label>
          <input
            id="login-user"
            bind:value={loginUser}
            required
            autocomplete="username"
            class="w-full px-3 py-3 lg:py-2 text-base lg:text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>
        <div class="space-y-1.5">
          <label for="login-pass" class="block text-xs font-medium" style="color:var(--muted-foreground)">Passwort</label>
          <input
            id="login-pass"
            type="password"
            bind:value={loginPass}
            required
            autocomplete="current-password"
            class="w-full px-3 py-3 lg:py-2 text-base lg:text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>
        {#if loginError}
          <p class="text-xs text-red-400">{loginError}</p>
        {/if}
        <button
          type="submit"
          class="w-full min-h-11 py-2.5 text-sm font-semibold text-white"
          style="background:var(--brand)"
        >Anmelden</button>
      </form>
    </div>
  </div>
{:else}

  <!-- Header -->
  <header class="border-b px-4 py-3 lg:px-6 lg:py-4 flex items-center gap-3" style="border-color:var(--border);background:var(--surface)">
    <span class="text-lg font-semibold tracking-wide" style="color:var(--foreground)">Deutsche Bank</span>
    <span class="text-xs px-2 py-0.5 font-medium" style="background:var(--brand);color:#fff">ADMIN</span>
    <button onclick={logout} class="ml-auto min-h-11 px-2 text-sm" style="color:var(--muted-foreground)">Abmelden</button>
  </header>

  <div class="lg:hidden grid grid-cols-2 border-b" style="border-color:var(--border);background:var(--surface)">
    <button
      type="button"
      class="min-h-11 text-sm font-medium"
      style="color:{panel === 'create' ? 'var(--foreground)' : 'var(--muted-foreground)'};border-bottom:2px solid {panel === 'create' ? 'var(--brand)' : 'transparent'}"
      aria-pressed={panel === 'create'}
      onclick={() => (panel = 'create')}
    >Fall erstellen</button>
    <button
      type="button"
      class="min-h-11 text-sm font-medium"
      style="color:{panel === 'list' ? 'var(--foreground)' : 'var(--muted-foreground)'};border-bottom:2px solid {panel === 'list' ? 'var(--brand)' : 'transparent'}"
      aria-pressed={panel === 'list'}
      onclick={() => (panel = 'list')}
    >Alle Fälle</button>
  </div>

  <div class="flex flex-col lg:flex-row lg:flex-1 lg:min-h-0 lg:overflow-hidden">

    <!-- Sidebar -->
    <aside class="w-full lg:w-72 lg:shrink-0 border-b lg:border-b-0 lg:border-r flex flex-col {panel === 'list' ? 'max-lg:hidden' : ''}" style="border-color:var(--border);background:var(--surface)">
      <div class="px-5 py-5 border-b" style="border-color:var(--border)">
        <p class="text-xs font-semibold uppercase tracking-widest" style="color:var(--muted-foreground)">Fall erstellen</p>
      </div>

      <form class="flex-1 px-4 py-4 lg:px-5 lg:py-5 space-y-4 lg:overflow-y-auto" onsubmit={handleSubmit}>

        <div class="space-y-1.5">
          <label for="case-id" class="block text-xs font-medium" style="color:var(--muted-foreground)">Fallnummer <span class="text-red-400">*</span></label>
          <input
            id="case-id"
            bind:value={caseId}
            required
            placeholder="CASE-2026-001"
            class="w-full px-3 py-3 lg:py-2 text-base lg:text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        <div class="space-y-1.5">
          <label for="iban" class="block text-xs font-medium" style="color:var(--muted-foreground)">IBAN <span class="text-red-400">*</span></label>
          <input
            id="iban"
            bind:value={iban}
            required
            placeholder="DE89 3704 0044 0532 0130 00"
            class="w-full px-3 py-3 lg:py-2 text-base lg:text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        <div class="space-y-1.5">
          <label for="empfaenger" class="block text-xs font-medium" style="color:var(--muted-foreground)">Empfänger <span class="text-red-400">*</span></label>
          <input
            id="empfaenger"
            bind:value={aesReference}
            required
            placeholder="Max Mustermann"
            class="w-full px-3 py-3 lg:py-2 text-base lg:text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        <div class="space-y-1.5">
          <label for="verwendungszweck" class="block text-xs font-medium" style="color:var(--muted-foreground)">Verwendungszweck <span class="text-red-400">*</span></label>
          <input
            id="verwendungszweck"
            bind:value={verwendungszweck}
            required
            placeholder="Rückbuchung Ref. 2026-001"
            class="w-full px-3 py-3 lg:py-2 text-base lg:text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        <div class="space-y-1.5">
          <label for="betrag" class="block text-xs font-medium" style="color:var(--muted-foreground)">Betrag (EUR) <span class="text-red-400">*</span></label>
          <input
            id="betrag"
            bind:value={betrag}
            required
            placeholder="1.250,00"
            onblur={formatBetrag}
            class="w-full px-3 py-3 lg:py-2 text-base lg:text-sm border outline-none focus:border-[var(--brand)]"
            style="background:var(--surface2);border-color:var(--border);color:var(--foreground)"
          />
        </div>

        {#if error}
          <p class="text-xs text-red-400">{error}</p>
        {/if}

        <button
          type="submit"
          disabled={loading}
          class="w-full min-h-11 py-2.5 text-sm font-semibold text-white transition-colors disabled:opacity-50"
          style="background:var(--brand)"
          onmouseenter={e => (e.currentTarget as HTMLElement).style.background = 'var(--brand-hover)'}
          onmouseleave={e => (e.currentTarget as HTMLElement).style.background = 'var(--brand)'}
        >
          {loading ? 'Wird erstellt...' : 'Fall erstellen'}
        </button>

      </form>
    </aside>

    <!-- Main: cases list -->
    <main class="flex-1 min-w-0 flex flex-col lg:overflow-hidden {panel === 'create' ? 'max-lg:hidden' : ''}">
      <div class="px-4 py-4 lg:px-6 lg:py-5 border-b flex items-center justify-between gap-3" style="border-color:var(--border)">
        <p class="text-xs font-semibold uppercase tracking-widest" style="color:var(--muted-foreground)">Alle Fälle</p>
        <button
          onclick={loadList}
          class="min-h-11 lg:min-h-0 text-sm lg:text-xs px-3 border transition-colors"
          style="border-color:var(--border);color:var(--muted-foreground)"
        >Aktualisieren</button>
      </div>

      <div class="flex-1 lg:overflow-auto">
        {#if loadingList}
          <p class="px-4 lg:px-6 py-8 text-sm" style="color:var(--muted-foreground)">Laden...</p>
        {:else if cases.length === 0}
          <p class="px-4 lg:px-6 py-8 text-sm" style="color:var(--muted-foreground)">Noch keine Fälle.</p>
        {:else}
          <ul class="lg:hidden">
            {#each cases as c (c.case_id)}
              <li class="px-4 py-4 border-b" style="border-color:var(--border)">
                <div class="flex items-start justify-between gap-3">
                  <div class="min-w-0">
                    <p class="font-mono text-sm break-all">{c.case_id}</p>
                    <p class="mt-1 text-sm break-all" style="color:var(--muted-foreground)">{c.aes_reference} · {c.iban}</p>
                    <p class="mt-1 text-xs" style="color:var(--muted-foreground)">{formatDate(c.created_at)}</p>
                  </div>
                  <div class="shrink-0 text-right">
                    <p class="text-sm tabular-nums">{c.betrag ? c.betrag + ' €' : '—'}</p>
                    {#if c.expired}
                      <span class="mt-1 inline-block text-xs px-2 py-0.5 font-medium" style="background:#3a1a1a;color:#f87171">Abgelaufen</span>
                    {:else}
                      <span class="mt-1 inline-block text-xs px-2 py-0.5 font-medium" style="background:#0f2a1a;color:#4ade80">Aktiv</span>
                    {/if}
                  </div>
                </div>
                <div class="mt-3 flex flex-wrap gap-2">
                  <button
                    type="button"
                    onclick={() => copyLink(c)}
                    class="min-h-11 px-3 text-sm border"
                    style="border-color:var(--border);color:var(--muted-foreground)"
                  >{copiedId === c.case_id ? 'Kopiert!' : 'Link kopieren'}</button>
                  {#if !c.expired}
                    <button
                      type="button"
                      onclick={() => handleExpire(c)}
                      disabled={expiringId === c.case_id}
                      class="min-h-11 px-3 text-sm border disabled:opacity-40"
                      style="border-color:#7f1d1d;color:#f87171"
                    >{expiringId === c.case_id ? '...' : 'Deaktivieren'}</button>
                  {/if}
                  <button
                    type="button"
                    onclick={() => handleDelete(c)}
                    disabled={deletingId === c.case_id}
                    class="min-h-11 px-3 text-sm border disabled:opacity-40"
                    style="border-color:#4b1a1a;color:#ef4444"
                  >{deletingId === c.case_id ? '...' : 'Löschen'}</button>
                </div>
              </li>
            {/each}
          </ul>
          <table class="hidden lg:table w-full text-sm border-collapse">
            <thead>
              <tr class="text-left border-b" style="border-color:var(--border)">
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Fallnummer</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">IBAN</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Empfänger</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Betrag</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Erstellt</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Status</th>
                <th class="px-4 py-3 text-xs font-semibold uppercase tracking-wide" style="color:var(--muted-foreground)">Aktionen</th>
              </tr>
            </thead>
            <tbody>
              {#each cases as c (c.case_id)}
                <tr class="border-b transition-colors" style="border-color:var(--border)" onmouseenter={e => (e.currentTarget as HTMLElement).style.background = 'var(--surface)'} onmouseleave={e => (e.currentTarget as HTMLElement).style.background = ''}>
                  <td class="px-4 py-3 font-mono text-xs">{c.case_id}</td>
                  <td class="px-4 py-3 font-mono text-xs" style="color:var(--muted-foreground)">{c.iban}</td>
                  <td class="px-4 py-3 text-xs" style="color:var(--muted-foreground)">{c.aes_reference}</td>
                  <td class="px-4 py-3 text-xs" style="color:var(--muted-foreground)">{c.betrag ? c.betrag + ' €' : '—'}</td>
                  <td class="px-4 py-3 text-xs" style="color:var(--muted-foreground)">{formatDate(c.created_at)}</td>
                  <td class="px-4 py-3">
                    {#if c.expired}
                      <span class="text-xs px-2 py-0.5 font-medium" style="background:#3a1a1a;color:#f87171">Abgelaufen</span>
                    {:else}
                      <span class="text-xs px-2 py-0.5 font-medium" style="background:#0f2a1a;color:#4ade80">Aktiv</span>
                    {/if}
                  </td>
                  <td class="px-4 py-3">
                    <div class="flex items-center gap-2">
                      <button
                        onclick={() => copyLink(c)}
                        class="text-xs px-2.5 py-1 border transition-colors"
                        style="border-color:var(--border);color:var(--muted-foreground)"
                      >{copiedId === c.case_id ? 'Kopiert!' : 'Link kopieren'}</button>
                      {#if !c.expired}
                        <button
                          onclick={() => handleExpire(c)}
                          disabled={expiringId === c.case_id}
                          class="text-xs px-2.5 py-1 border transition-colors disabled:opacity-40"
                          style="border-color:#7f1d1d;color:#f87171"
                        >{expiringId === c.case_id ? '...' : 'Deaktivieren'}</button>
                      {/if}
                      <button
                        onclick={() => handleDelete(c)}
                        disabled={deletingId === c.case_id}
                        class="text-xs px-2.5 py-1 border transition-colors disabled:opacity-40"
                        style="border-color:#4b1a1a;color:#ef4444"
                      >{deletingId === c.case_id ? '...' : 'Löschen'}</button>
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
