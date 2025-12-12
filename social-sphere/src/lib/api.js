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

// // ==========================================
// // SERVER-SIDE API REQUEST
// // ==========================================
// export async function serverApiRequest(endpoint, options = {}) {
//     const isFormData = options.body instanceof FormData;
    
//     const headers = { ...options.headers };
    
//     if (!isFormData) {
//         headers['Content-Type'] = 'application/json';
//     }
    
//     try {
//         const { cookies } = await import('next/headers');
//         const cookieStore = cookies();
//         const jwtCookie = cookieStore.get("jwt");
        
//         if (jwtCookie) {
//             headers['Cookie'] = `jwt=${jwtCookie.value}`;
//         } else {
//             console.warn('No JWT cookie found in server request');
//             // Don't throw - let the API handle unauthorized requests
//         }
//     } catch (error) {
//         console.error('Failed to access cookies:', error);
//     }
    
//     const res = await fetch(`${API_BASE}${endpoint}`, {
//         ...options,
//         headers,
//         cache: 'no-store',
//     });
    
//     if (!res.ok) {
//         const error = await res.json().catch(() => ({}));
//         throw new Error(error.error || `API error: ${res.status}`);
//     }
    
//     return res.json();
// }