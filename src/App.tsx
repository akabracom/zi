import { BrowserRouter, Route, Routes } from 'react-router-dom'
import './App.css'
import { AuthProvider } from './context/AuthContext'
import { Layout } from './components/Layout'
import { ProtectedRoute } from './components/ProtectedRoute'
import { ServicesPage } from './pages/ServicesPage'
import { DraftPage } from './pages/DraftPage'
import { RequestsPage } from './pages/RequestsPage'
import { LoginPage } from './pages/LoginPage'
import { RegisterPage } from './pages/RegisterPage'

function App() {
  return (
    <BrowserRouter basename={import.meta.env.BASE_URL || '/zi/'}>
      <AuthProvider>
        <Routes>
          <Route element={<Layout />}>
            <Route index element={<ServicesPage />} />
            <Route element={<ProtectedRoute />}>
              <Route path="draft" element={<DraftPage />} />
              <Route path="requests" element={<RequestsPage />} />
            </Route>
            <Route path="login" element={<LoginPage />} />
            <Route path="register" element={<RegisterPage />} />
            <Route path="*" element={<ServicesPage />} />
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
