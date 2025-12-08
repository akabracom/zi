import { Navigate, Outlet } from 'react-router-dom'
import { useAuth, useRole } from '../context/AuthContext'

interface Props {
  allowedRoles?: ('user' | 'moderator')[]
}

export function ProtectedRoute({ allowedRoles }: Props) {
  const { user } = useAuth()
  const role = useRole()

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (allowedRoles && !allowedRoles.includes(role as 'user' | 'moderator')) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}











