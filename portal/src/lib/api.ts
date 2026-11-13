const API_URL = apiBase()

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

export interface CaseData {
  case_id: string
  iban: string
  aes_reference: string
  verwendungszweck: string
  betrag: string
  status: string
  created_at: string
}

export type CaseResult =
  | { kind: 'ok'; data: CaseData }
  | { kind: 'invalid' }
  | { kind: 'expired' }

export async function fetchCase(token: string): Promise<CaseResult> {
  const res = await fetch(`${API_URL}/api/case?token=${encodeURIComponent(token)}`)
  if (res.status === 404) return { kind: 'invalid' }
  if (res.status === 410) return { kind: 'expired' }
  if (!res.ok) return { kind: 'invalid' }
  const data = await res.json()
  return { kind: 'ok', data }
}
