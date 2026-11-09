<script lang="ts">
  import Button from './lib/components/ui/Button.svelte'
  import Card from './lib/components/ui/Card.svelte'
  import Input from './lib/components/ui/Input.svelte'
  import Select from './lib/components/ui/Select.svelte'
  import Textarea from './lib/components/ui/Textarea.svelte'
  import { createCase } from './lib/api'

  let caseId = $state('')
  let iban = $state('')
  let aesReference = $state('')
  let issueType = $state<'failed_transfer' | 'dispute'>('failed_transfer')
  let issueDescription = $state('')

  let loading = $state(false)
  let error = $state('')
  let result = $state<{ case_id: string; portal_link: string } | null>(null)
  let copied = $state(false)

  async function handleSubmit(e: Event) {
    e.preventDefault()
    error = ''
    loading = true
    try {
      result = await createCase({
        case_id: caseId,
        iban,
        aes_reference: aesReference,
        issue_type: issueType,
        issue_description: issueDescription,
      })
    } catch (err) {
      error = err instanceof Error ? err.message : 'Something went wrong'
    } finally {
      loading = false
    }
  }

  async function copyLink() {
    if (!result) return
    await navigator.clipboard.writeText(result.portal_link)
    copied = true
    setTimeout(() => (copied = false), 1500)
  }

  function startNewCase() {
    result = null
    caseId = ''
    iban = ''
    aesReference = ''
    issueType = 'failed_transfer'
    issueDescription = ''
  }
</script>

<div class="min-h-screen bg-white">
  <header class="border-b border-[var(--border)] px-6 py-4">
    <h1 class="text-base font-semibold text-[var(--foreground)]">Case Portal — Admin</h1>
  </header>

  <main class="mx-auto max-w-xl px-6 py-10">
    {#if result}
      <Card class="p-6">
        <h2 class="text-sm font-medium text-[var(--muted-foreground)]">Case created</h2>
        <p class="mt-1 text-lg font-semibold">{result.case_id}</p>

        <div class="mt-4 flex items-center gap-2 rounded-md border border-[var(--border)] bg-[var(--muted)] px-3 py-2">
          <span class="flex-1 truncate text-sm text-[var(--foreground)]">{result.portal_link}</span>
          <Button variant="outline" onclick={copyLink}>{copied ? 'Copied' : 'Copy'}</Button>
        </div>

        <div class="mt-6">
          <Button variant="outline" onclick={startNewCase}>Create another case</Button>
        </div>
      </Card>
    {:else}
      <Card class="p-6">
        <form class="space-y-4" onsubmit={handleSubmit}>
          <div>
            <label for="case_id" class="mb-1 block text-sm font-medium">Case ID</label>
            <Input id="case_id" bind:value={caseId} placeholder="CASE-2026-001 (optional)" />
          </div>

          <div>
            <label for="iban" class="mb-1 block text-sm font-medium">IBAN</label>
            <Input id="iban" bind:value={iban} placeholder="GB29 NWBK 6016 1331 9268 19" required />
          </div>

          <div>
            <label for="aes" class="mb-1 block text-sm font-medium">AES Reference</label>
            <Input id="aes" bind:value={aesReference} placeholder="AES-REF-00123" required />
          </div>

          <div>
            <label for="issue_type" class="mb-1 block text-sm font-medium">Issue Type</label>
            <Select
              id="issue_type"
              bind:value={issueType}
              options={[
                { label: 'Failed Transfer', value: 'failed_transfer' },
                { label: 'Dispute', value: 'dispute' },
              ]}
            />
          </div>

          <div>
            <label for="desc" class="mb-1 block text-sm font-medium">Issue Description</label>
            <Textarea id="desc" bind:value={issueDescription} placeholder="Describe the issue..." required />
          </div>

          {#if error}
            <p class="text-sm text-red-600">{error}</p>
          {/if}

          <Button type="submit" disabled={loading}>
            {loading ? 'Creating...' : 'Create Case'}
          </Button>
        </form>
      </Card>
    {/if}
  </main>
</div>
