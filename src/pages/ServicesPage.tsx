import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { addServiceToDraft, getDraftRequest, getServices } from '../api/http'
import type { Service } from '../api/types'
import { useAuth, useRole } from '../context/AuthContext'

// Флаг для проверки mock режима
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true' || false

export function ServicesPage() {
  const [services, setServices] = useState<Service[]>([])
  const [query, setQuery] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [info, setInfo] = useState<string | null>(null)
  const [cartCount, setCartCount] = useState(0)
  const [imageError, setImageError] = useState<Set<number>>(new Set())
  const [cartIconBroken, setCartIconBroken] = useState(false)
  const role = useRole()
  const { user } = useAuth()

  function handleImageError(serviceId: number) {
    setImageError((prev) => {
      const next = new Set(prev)
      next.add(serviceId)
      return next
    })
  }

  const renderPlaceholder = (name: string) => (
    <div
      className="service-image-placeholder"
      style={{
        width: '100%',
        height: '150px',
        backgroundColor: '#e0e0e0',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        color: '#666',
        fontSize: '16px',
        fontWeight: '500',
        borderRadius: '4px',
        marginBottom: '10px',
        border: '2px dashed #999',
      }}
    >
      {name}
    </div>
  )

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const data = await getServices(query)
      setServices(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки месяцев')
    } finally {
      setLoading(false)
    }
  }

  async function loadCartCount() {
    if (!user) {
      setCartCount(0)
      return
    }
    try {
      const draft = await getDraftRequest()
      const count = draft.services.reduce((sum, item) => sum + item.quantity, 0)
      setCartCount(count)
    } catch {
      setCartCount(0)
    }
  }

  useEffect(() => {
    void load()
    void loadCartCount()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  async function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    await load()
  }

  async function handleAddToDraft(serviceId: number) {
    setInfo(null)
    setError(null)
    try {
      await addServiceToDraft(serviceId)
      setInfo('Месяц добавлен в расчет вклада')
      await loadCartCount()
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Не удалось добавить месяц в расчет',
      )
    }
  }

  return (
    <div className="page">
      <h2 className="catalog-title">Банковские месяцы</h2>
      <p className="catalog-subtitle">Рассчитайте проценты по вкладу для каждого месяца с учетом количества дней и ставок</p>
      <div className="header-bar">
        <form onSubmit={handleSearch} className="toolbar">
          <input
            type="search"
            placeholder="Поиск по названию месяца..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <button
            type="submit"
            disabled={loading}
          >
            Найти
          </button>
        </form>
        <Link to="/draft" className="btn btn-cart" aria-label="Корзина">
          {USE_MOCK ? (
            <span className="cart-icon" style={{
              fontSize: '24px',
              display: 'inline-block',
              width: '24px',
              height: '24px',
              lineHeight: '24px',
              textAlign: 'center'
            }}>🛒</span>
          ) : (
            cartIconBroken ? (
              <span className="cart-icon" style={{
                fontSize: '24px',
                display: 'inline-block',
                width: '24px',
                height: '24px',
                lineHeight: '24px',
                textAlign: 'center'
              }}>🛒</span>
            ) : (
              <img
                src="/images/basket.png"
                alt="Cart"
                className="cart-icon"
                onError={() => setCartIconBroken(true)}
              />
            )
          )}
          {cartCount > 0 && (
            <span className="cart-count">{cartCount}</span>
          )}
        </Link>
      </div>

      {error && <div className="form-error">{error}</div>}
      {info && <div className="form-success">{info}</div>}

      {loading && <p>Загрузка...</p>}

      {!loading && !error && (
        <div className="cards-grid">
          {services.map((s) => (
            <article key={s.id} className="card service-card">
              {/* Изображение или заглушка */}
              {USE_MOCK ? (
                renderPlaceholder(s.name)
              ) : s.image_url && !imageError.has(s.id) ? (
                // В реальном режиме - изображение из бэкенда через прокси
                <img
                  src={
                    s.image_url.startsWith('http')
                      ? s.image_url
                      : `/images/${s.image_url.replace(/^(img\/|images\/|deposits\/img\/)/, '')}`
                  }
                  alt={s.name}
                  className="service-image"
                  onError={() => handleImageError(s.id)}
                />
              ) : (
                // Если image_url пустой или не загрузилось, показываем заглушку
                renderPlaceholder(s.name)
              )}
              <h3>{s.name}</h3>
              {s.description && <p className="muted">{s.description}</p>}
              <p className="price">Ставка: {s.price.toFixed(1)}%</p>
            {role === 'guest' ? (
              <p className="muted">
                Чтобы рассчитать вклад, <strong>войдите</strong> или{' '}
                <strong>зарегистрируйтесь</strong>.
              </p>
            ) : (
              <button
                type="button"
                className="app-button app-button-primary btn-add-to-cart"
                onClick={() => handleAddToDraft(s.id)}
              >
                Добавить в расчет
              </button>
            )}
            </article>
          ))}
        </div>
      )}

      {!loading && services.length === 0 && !error && (
        <p className="muted" style={{ textAlign: 'center', padding: '40px' }}>
          Месяцев не найдено.
        </p>
      )}
    </div>
  )
}


