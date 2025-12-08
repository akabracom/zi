import { type FormEvent, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { login } from '../api/http'
import { useAuth } from '../context/AuthContext'

export function LoginPage() {
  const { login: saveAuth } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const resp = await login({ email, password })
      saveAuth(resp)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="card">
      <h2>Вход</h2>
      <form onSubmit={handleSubmit} className="form">
        <label className="form-field">
          <span>Email</span>
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </label>
        <label className="form-field">
          <span>Пароль</span>
          <input
            type="password"
            required
            minLength={6}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>
        {error && <div className="form-error">{error}</div>}
        <button
          type="submit"
          className="app-button app-button-primary"
          disabled={loading}
        >
          {loading ? 'Входим...' : 'Войти'}
        </button>
      </form>
      <p className="muted">
        Нет аккаунта? <Link to="/register">Зарегистрироваться</Link>
      </p>
    </div>
  )
}












