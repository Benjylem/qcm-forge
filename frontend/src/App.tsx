import { useEffect, useState } from 'react'
import AuthForm from './AuthForm'
import { logout, me } from './api'

type AuthState = 'loading' | 'anonymous' | { userId: number }

function App() {
  const [auth, setAuth] = useState<AuthState>('loading')

  function checkSession() {
    me()
      .then((data) => setAuth({ userId: data.user_id }))
      .catch(() => setAuth('anonymous'))
  }

  useEffect(checkSession, [])

  async function handleLogout() {
    await logout().catch(() => {}) // déconnexion best-effort côté UI
    setAuth('anonymous')
  }

  return (
    <main>
      <h1>qcm-forge</h1>

      {auth === 'loading' && <p>Chargement...</p>}

      {auth === 'anonymous' && <AuthForm onAuthenticated={checkSession} />}

      {typeof auth === 'object' && (
        <div>
          <p>Connecté (utilisateur #{auth.userId})</p>
          <button onClick={handleLogout}>Se déconnecter</button>
        </div>
      )}
    </main>
  )
}

export default App
