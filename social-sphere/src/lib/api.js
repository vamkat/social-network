const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'

// ==========================================
// CLIENT-SIDE API REQUEST
// ==========================================
export async function apiRequest(endpoint, options = {}) {
  const isFormData = options.body instanceof FormData

  const headers = {
    ...options.headers,
  }

  // Only set Content-Type for non-FormData requests
  if (!isFormData) {
    headers['Content-Type'] = 'application/json'
  }

  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    credentials: 'include',
    headers,
  })

  if (!res.ok) {
    const error = await res.json().catch(() => ({}));
    throw new Error(error.error || `API error: ${res.status}`);
  }

  return res.json()
}