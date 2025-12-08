import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  Request,
  Service,
} from './types'
import { filterMockServices, mockDraft } from './mock'

// Флаг для переключения между mock и реальным API
// Установите USE_MOCK=true в .env или измените здесь для демонстрации без бэкенда
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true' || false

const JSON_HEADERS: HeadersInit = {
  'Content-Type': 'application/json',
}

function getToken(): string | null {
  return localStorage.getItem('auth_token')
}

async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
  opts?: { auth?: boolean },
): Promise<T> {
  const headers: HeadersInit = {
    ...(options.headers || {}),
  }

  if (options.body && !(options.body instanceof FormData)) {
    Object.assign(headers, JSON_HEADERS)
  }

  const needAuth = opts?.auth ?? false
  const token = getToken()

  if (needAuth && token) {
    ;(headers as Record<string, string>).Authorization = `Bearer ${token}`
  }

  const response = await fetch(path, {
    ...options,
    headers,
  })

  if (!response.ok) {
    const text = await response.text()
    let message = `HTTP ${response.status}`
    try {
      const parsed = JSON.parse(text) as { error?: string }
      if (parsed.error) message = parsed.error
    } catch {
      // ignore
    }
    throw new Error(message)
  }

  if (response.status === 204) {
    // @ts-expect-error – no body
    return null
  }

  return (await response.json()) as T
}

// ===== AUTH =====

export function login(payload: LoginRequest): Promise<AuthResponse> {
  return apiFetch<AuthResponse>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

export function register(payload: RegisterRequest): Promise<AuthResponse> {
  return apiFetch<AuthResponse>('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
}

// ===== SERVICES =====

export function getServices(query?: string): Promise<Service[]> {
  if (USE_MOCK) {
    // Имитация асинхронного запроса с задержкой
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(filterMockServices(query || ''))
      }, 300)
    })
  }
  const q = query ? `?query=${encodeURIComponent(query)}` : ''
  // GET /services допускает гостей, поэтому auth не обязателен
  return apiFetch<Service[]>(`/api/services${q}`, {}, { auth: !!getToken() })
}

export function createService(data: {
  name: string
  price: number
}): Promise<Service> {
  return apiFetch<Service>(
    '/api/services',
    {
      method: 'POST',
      body: JSON.stringify(data),
    },
    { auth: true },
  )
}

export function updateService(
  id: number,
  data: { name?: string; price?: number },
): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(
    `/api/services/${id}`,
    {
      method: 'PUT',
      body: JSON.stringify(data),
    },
    { auth: true },
  )
}

export function deleteService(id: number): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(
    `/api/services/${id}`,
    { method: 'DELETE' },
    { auth: true },
  )
}

export async function uploadServiceImage(
  id: number,
  file: File,
): Promise<{ image_url: string }> {
  const form = new FormData()
  form.append('image', file)

  const token = getToken()
  const headers: HeadersInit = {}
  if (token) {
    ;(headers as Record<string, string>).Authorization = `Bearer ${token}`
  }

  const response = await fetch(`/api/services/${id}/image`, {
    method: 'POST',
    body: form,
    headers,
  })

  if (!response.ok) {
    throw new Error('Не удалось загрузить изображение')
  }

  return (await response.json()) as { image_url: string }
}

// ===== REQUESTS =====

export function getDraftRequest(): Promise<Request> {
  if (USE_MOCK) {
    return new Promise((resolve) => {
      setTimeout(() => resolve(mockDraft), 300)
    })
  }
  return apiFetch<Request>('/api/requests/draft', {}, { auth: true })
}

export function addServiceToDraft(serviceId: number): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(
    '/api/requests/add-service',
    {
      method: 'POST',
      body: JSON.stringify({ service_id: serviceId }),
    },
    { auth: true },
  )
}

export function removeServiceFromRequest(
  requestId: number,
  serviceId: number,
): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(
    '/api/requests/remove-service',
    {
      method: 'DELETE',
      body: JSON.stringify({ request_id: requestId, service_id: serviceId }),
    },
    { auth: true },
  )
}

export function updateServiceQuantity(
  requestId: number,
  serviceId: number,
  quantity: number,
): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(
    `/api/requests/${requestId}/services/${serviceId}/quantity`,
    {
      method: 'PUT',
      body: JSON.stringify({ quantity }),
    },
    { auth: true },
  )
}

export function submitRequest(
  requestId: number,
): Promise<{ status: string }> {
  if (USE_MOCK) {
    // В mock режиме добавляем заявку в список и возвращаем успех
    return new Promise((resolve) => {
      setTimeout(() => {
        // Получаем черновик и добавляем его в список с обновленным статусом
        const draft = { ...mockDraft, id: requestId, status: 'formed' as const }
        // Проверяем, нет ли уже такой заявки
        const exists = mockRequestsStorage.find(r => r.id === requestId)
        if (!exists) {
          mockRequestsStorage.push(draft)
        } else {
          // Обновляем существующую
          const index = mockRequestsStorage.findIndex(r => r.id === requestId)
          if (index >= 0) {
            mockRequestsStorage[index] = { ...mockRequestsStorage[index], status: 'formed' }
          }
        }
        resolve({ status: 'submitted' })
      }, 300)
    })
  }
  return apiFetch<{ status: string }>(
    `/api/requests/${requestId}/submit`,
    { method: 'PUT' },
    { auth: true },
  )
}

export function completeRequest(
  requestId: number,
): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(
    `/api/requests/${requestId}/complete`,
    { method: 'PUT' },
    { auth: true },
  )
}

export function deleteRequestLogical(
  requestId: number,
): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(
    `/api/requests/${requestId}/delete`,
    { method: 'POST' },
    { auth: true },
  )
}

// Глобальное хранилище для mock заявок (чтобы они сохранялись после отправки)
let mockRequestsStorage: Request[] = []

export function listRequests(params?: {
  date_from?: string
  date_to?: string
}): Promise<Request[]> {
  if (USE_MOCK) {
    return new Promise((resolve) => {
      setTimeout(() => resolve(mockRequestsStorage), 300)
    })
  }
  const search = new URLSearchParams()
  if (params?.date_from) search.set('date_from', params.date_from)
  if (params?.date_to) search.set('date_to', params.date_to)
  const qs = search.toString()
  const url = qs ? `/api/requests?${qs}` : '/api/requests'
  return apiFetch<Request[]>(url, {}, { auth: true })
}



