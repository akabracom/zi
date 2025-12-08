import type { Service, Request } from './types'

// Mock данные для демонстрации без бэкенда
// Ставки соответствуют 2-й лабе (di_web)
export const mockServices: Service[] = [
  {
    id: 1,
    name: 'Январь',
    description: 'Дней: 31 | Ставка: 10.5% | Сезон: Зима | Праздничные дни: 8',
    price: 10.5,
    is_active: true,
    image_url: 'january.png',
    created_at: '2025-12-03T18:27:17.520874Z',
  },
  {
    id: 2,
    name: 'Февраль',
    description: 'Дней: 28 | Ставка: 10.2% | Сезон: Зима | Праздничные дни: 1',
    price: 10.2,
    is_active: true,
    image_url: 'february.png',
    created_at: '2025-12-03T18:27:17.520874Z',
  },
  {
    id: 3,
    name: 'Март',
    description: 'Дней: 31 | Ставка: 10.0% | Сезон: Весна | Праздничные дни: 1',
    price: 10.0,
    is_active: true,
    image_url: 'march.png',
    created_at: '2025-12-03T18:27:17.520874Z',
  },
]

export const mockDraft: Request = {
  id: 1,
  user_id: 1,
  status: 'draft',
  created_at: '2025-12-03T18:27:17.520874Z',
  updated_at: '2025-12-03T18:27:17.520874Z',
  services: [
    {
      request_id: 1,
      service_id: 1,
      quantity: 2,
      service: mockServices[0],
    },
  ],
}

// Начально пустой список - расчеты будут добавляться после отправки
export const mockRequests: Request[] = []

// Функция для фильтрации mock данных
export function filterMockServices(query: string): Service[] {
  if (!query) return mockServices
  const lowerQuery = query.toLowerCase()
  return mockServices.filter(
    (service) =>
      service.name.toLowerCase().includes(lowerQuery) ||
      service.description?.toLowerCase().includes(lowerQuery),
  )
}









