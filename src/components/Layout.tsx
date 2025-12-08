import { Link, NavLink, Outlet } from 'react-router-dom'
import { useAuth, useRole } from '../context/AuthContext'

export function Layout() {
  const { user, logout } = useAuth()
  const role = useRole()

  return (
    <div className="app-root">
      <header className="app-header">
        <Link to="/" className="app-header-title">
          <h1>КАЛЬКУЛЯТОР ПРОЦЕНТОВ ПО ВКЛАДАМ</h1>
        </Link>
        <div className="app-header-right">
          {user ? (
            <>
              <span className="app-user">
                {user.name || user.email}
              </span>
              <button
                type="button"
                className="app-button app-button-outline"
                onClick={logout}
              >
                Выйти
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="app-button app-button-outline">
                Войти
              </Link>
              <Link to="/register" className="app-button app-button-primary">
                Регистрация
              </Link>
            </>
          )}
        </div>
      </header>

      <main className="app-main">
        {role !== 'guest' && (
          <nav className="page-nav">
            <NavLink to="/draft" className="page-nav-link">
              Расчет
            </NavLink>
            <NavLink to="/requests" className="page-nav-link">
              Расчеты
            </NavLink>
          </nav>
        )}
        <Outlet />
      </main>
    </div>
  )
}


