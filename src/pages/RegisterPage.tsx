import { type FormEvent, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { register } from '../api/http'
import { useAuth } from '../context/AuthContext'

export function RegisterPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const resp = await register({ name, email, password })
      login(resp)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="card">
      <h2>Регистрация</h2>
      <form onSubmit={handleSubmit} className="form">
        <label className="form-field">
          <span>Имя</span>
          <input
            type="text"
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </label>
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
          {loading ? 'Регистрируем...' : 'Зарегистрироваться'}
        </button>
      </form>
      <p className="muted">
        Уже есть аккаунт? <Link to="/login">Войти</Link>
      </p>
    </div>
  )
}












