export type UserRole = 'guest' | 'user' | 'moderator'

export interface User {
  id: number
  email: string
  name: string
  role: UserRole
  created_at?: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
}

export interface Service {
  id: number
  name: string
  description?: string
  price: number
  image_url?: string
  is_active: boolean
  created_at: string
}

export interface RequestService {
  request_id: number
  service_id: number
  quantity: number
  service: Service
}

export type RequestStatus = 'draft' | 'formed' | 'completed' | string

export interface Request {
  id: number
  user_id: number
  status: RequestStatus
  created_at: string
  updated_at: string
  services: RequestService[]
}












