const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export interface CaseData {
  case_id: string
  iban: string
  aes_reference: string
  issue_type: 'failed_transfer' | 'dispute'
  issue_description: string
  status: string
  created_at: string
}

export type CaseResult =
  | { kind: 'ok'; data: CaseData }
  | { kind: 'invalid' }
  | { kind: 'used' }

// Interprets the server's 404/410/200 responses; no business logic here,
// the API alone decides whether a token is valid, used, or fresh.
export async function fetchCase(token: string): Promise<CaseResult> {
  const res = await fetch(`${API_URL}/api/case?token=${encodeURIComponent(token)}`)
  if (res.status === 404) return { kind: 'invalid' }
  if (res.status === 410) return { kind: 'used' }
  if (!res.ok) return { kind: 'invalid' }
  const data = await res.json()
  return { kind: 'ok', data }
}
