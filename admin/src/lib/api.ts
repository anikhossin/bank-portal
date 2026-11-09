const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export interface CreateCasePayload {
  case_id: string
  iban: string
  aes_reference: string
  issue_type: 'failed_transfer' | 'dispute'
  issue_description: string
}

export interface CreateCaseResult {
  case_id: string
  portal_link: string
}

// All validation and business logic lives server-side in the Go API;
// this just forwards the form and surfaces whatever the server decides.
export async function createCase(payload: CreateCasePayload): Promise<CreateCaseResult> {
  const res = await fetch(`${API_URL}/api/cases`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  const data = await res.json()
  if (!res.ok) {
    throw new Error(data.error || 'Failed to create case')
  }
  return data
}
