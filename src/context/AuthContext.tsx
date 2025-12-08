import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import type { AuthResponse, User, UserRole } from '../api/types'

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
}

interface AuthContextValue extends AuthState {
  login: (data: AuthResponse) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

const STORAGE_KEY = 'spa_deposits_auth'

function loadInitialState(): AuthState {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return { user: null, token: null, loading: false }
    }
    const parsed = JSON.parse(raw) as { user: User; token: string }
    localStorage.setItem('auth_token', parsed.token)
    return { user: parsed.user, token: parsed.token, loading: false }
  } catch {
    return { user: null, token: null, loading: false }
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>(() => loadInitialState())

  useEffect(() => {
    if (state.token && state.user) {
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({ token: state.token, user: state.user }),
      )
      localStorage.setItem('auth_token', state.token)
    } else {
      localStorage.removeItem(STORAGE_KEY)
      localStorage.removeItem('auth_token')
    }
  }, [state.token, state.user])

  const value = useMemo<AuthContextValue>(
    () => ({
      ...state,
      login: (data: AuthResponse) => {
        setState({
          user: data.user,
          token: data.token,
          loading: false,
        })
      },
      logout: () => {
        setState({ user: null, token: null, loading: false })
      },
    }),
    [state],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}

export function useRole(): UserRole {
  const { user } = useAuth()
  return user?.role ?? 'guest'
}











