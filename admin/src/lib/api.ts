const API_URL = apiBase()

// Empty means same-origin `/api`, which Vite proxies to the Go server.
// A localhost URL baked in at build time is ignored on a public host, because
// the visitor's browser would call their own machine and fail with NetworkError.
function apiBase() {
  const configured = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '')
  if (!configured) return ''
  if (typeof window === 'undefined') return configured
  const pageHost = window.location.hostname
  const pointsAtLocalMachine = /^https?:\/\/(localhost|127\.0\.0\.1)(:|\/|$)/.test(configured)
  const pageIsLocal = pageHost === 'localhost' || pageHost === '127.0.0.1'
  if (pointsAtLocalMachine && !pageIsLocal) return ''
  return configured
}

export interface CreateCasePayload {
  case_id: string
  iban: string
  aes_reference: string
  verwendungszweck: string
  betrag: string
}

export interface CreateCaseResult {
  case_id: string
  portal_link: string
}

export interface CaseSummary {
  case_id: string
  iban: string
  aes_reference: string
  verwendungszweck: string
  betrag: string
  status: string
  created_at: string
  portal_link: string
  expired: boolean
}

export async function createCase(payload: CreateCasePayload): Promise<CreateCaseResult> {
  const res = await fetch(`${API_URL}/api/cases`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `HTTP ${res.status}`)
  }
  return res.json()
}

export async function listCases(): Promise<CaseSummary[]> {
  const res = await fetch(`${API_URL}/api/cases`)
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  return res.json()
}

export async function deleteCase(caseId: string): Promise<void> {
  const res = await fetch(`${API_URL}/api/cases/delete`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ case_id: caseId }),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `HTTP ${res.status}`)
  }
}

export async function expireCase(caseId: string): Promise<void> {
  const res = await fetch(`${API_URL}/api/cases/expire`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ case_id: caseId }),
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `HTTP ${res.status}`)
  }
}
