import { useState } from 'react'
import './AuthForm.css'
import { ApiError, login, signup } from './api'

type Mode = 'login' | 'signup'

type Props = {
  onAuthenticated: () => void
}

function AuthForm({ onAuthenticated }: Props) {
  const [mode, setMode] = useState<Mode>('login')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [info, setInfo] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setInfo(null)
    setSubmitting(true)

    try {
      if (mode === 'signup') {
        await signup({ email, password })
        // Pas de connexion automatique : l'utilisateur doit prouver
        // qu'il connaît son mot de passe en se connectant lui-même,
        // comme pour n'importe quel compte existant.
        setMode('login')
        setPassword('')
        setInfo('Compte créé ! Connecte-toi avec tes identifiants.')
        return
      }

      await login({ email, password })
      onAuthenticated()
    } catch (err) {
      // Seul ApiError porte un message destiné à l'utilisateur ; une
      // autre erreur (réseau coupé, etc.) reste générique exprès.
      setError(err instanceof ApiError ? err.message : 'Impossible de contacter le serveur')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form className="auth-form" onSubmit={handleSubmit}>
      <h2>{mode === 'login' ? 'Connexion' : 'Inscription'}</h2>

      <div className="auth-field">
        <label htmlFor="email">Email</label>
        <input
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
      </div>

      <div className="auth-field">
        <label htmlFor="password">Mot de passe</label>
        <input
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          minLength={8}
          required
        />
      </div>

      {info && <p className="auth-info">{info}</p>}

      {error && (
        <p className="auth-error" role="alert">
          {error}
        </p>
      )}

      <button className="auth-submit" type="submit" disabled={submitting}>
        {mode === 'login' ? 'Se connecter' : "S'inscrire"}
      </button>

      <button
        className="auth-switch"
        type="button"
        onClick={() => {
          setMode(mode === 'login' ? 'signup' : 'login')
          setError(null)
          setInfo(null)
        }}
      >
        {mode === 'login' ? "Pas de compte ? S'inscrire" : 'Déjà un compte ? Se connecter'}
      </button>
    </form>
  )
}

export default AuthForm
