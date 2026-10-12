type Credentials = {
  email: string
  password: string
}

class ApiError extends Error {}

// Toutes les routes passent par ici. credentials: 'include' force le
// navigateur à envoyer le cookie de session même si fetch, par
// défaut, serait plus restrictif selon le contexte.
async function request(path: string, options: RequestInit = {}) {
  const res = await fetch(path, {
    ...options,
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...options.headers },
  })

  // Lire en texte d'abord : certaines réponses n'ont aucun corps
  // (signup 201, login 200, logout 204), et res.json() sur un corps
  // vide lève une erreur de parsing qui n'a rien à voir avec une
  // vraie erreur réseau. res.text() ne lève jamais, même sur "".
  const text = await res.text()
  const data = text ? JSON.parse(text) : null

  if (!res.ok) {
    throw new ApiError(data?.error ?? res.statusText)
  }

  return data
}

export function signup(creds: Credentials) {
  return request('/signup', { method: 'POST', body: JSON.stringify(creds) })
}

export function login(creds: Credentials) {
  return request('/login', { method: 'POST', body: JSON.stringify(creds) })
}

export function logout() {
  return request('/logout', { method: 'POST' })
}

export function me(): Promise<{ user_id: number }> {
  return request('/me')
}

export { ApiError }
