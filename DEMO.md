# Инструкция по демонстрации 5-й лабораторной работы

## Порядок показа

### 1. Показать три страницы фронтенда с mock без запущенного сервиса

1. **Включить mock режим:**
   - Создайте файл `.env` в корне проекта `spa-deposits`
   - Добавьте строку: `VITE_USE_MOCK=true`
   - Перезапустите dev server: `npm run dev`

2. **Показать три страницы:**
   - **Главная страница** (`/`) - список банковских месяцев с поиском
   - **Страница черновика** (`/draft`) - расчет вклада (требует авторизации)
   - **Страница заявок** (`/requests`) - список всех расчетов (требует авторизации)

3. **Показать работу поиска:**
   - Введите "Январь" в поисковую строку
   - Показать, что фильтрация работает на клиенте

### 2. Показать страницы фронтенда с бэкендом

1. **Выключить mock режим:**
   - В файле `.env` установите: `VITE_USE_MOCK=false`
   - Или удалите файл `.env`
   - Перезапустите dev server

2. **Убедиться, что бэкенд запущен:**
   ```bash
   cd SPA_Dep
   docker-compose up -d
   ```

3. **Показать в Network (F12 → Network):**
   - Открыть страницу со списком услуг
   - Показать запрос: `GET http://localhost:5173/api/services`
   - Объяснить, что запрос идет через Vite proxy на `http://localhost:3001/api/services`
   - Показать, что в Headers видно проксирование

### 3. Показать названия и адреса из MinIO для изображений

1. **В Network показать запросы к изображениям:**
   - `GET http://localhost:5173/images/january.png`
   - Объяснить, что запрос проксируется на `http://localhost:3001/images/january.png`
   - Показать, что бэкенд получает изображение из MinIO по пути `img/january.png`

2. **Показать в коде (`ServicesPage.tsx`):**
   ```tsx
   src={s.image_url.startsWith('http') 
     ? s.image_url 
     : `http://localhost:3001/images/${s.image_url.replace(/^(img\/|images\/|deposits\/img\/)/, '')}`}
   ```
   - Объяснить, что `image_url` из БД содержит только имя файла (например, `january.png`)
   - Код формирует полный URL для загрузки через бэкенд

### 4. Внести изменения в БД, показать их во фронтенде

1. **Добавить новый месяц через API или напрямую в БД:**
   ```sql
   INSERT INTO services (name, description, price, image_url) 
   VALUES ('Апрель', 'Дней: 30 | Ставка: 6.0% | Сезон: Весна', 6.0, 'april.png');
   ```

2. **Обновить страницу во фронтенде:**
   - Показать, что новый месяц появился в списке
   - Показать в Network новый запрос к API

### 5. Объяснить код компонентов для фильтрации

#### Компонент: `ServicesPage.tsx`

**Хуки:**
- `useState<Service[]>([])` - состояние для списка услуг
- `useState('')` - состояние для поискового запроса
- `useState(false)` - состояние загрузки
- `useEffect(() => { load() }, [])` - загрузка данных при монтировании компонента

**Props (передаваемые данные):**
- Компонент не принимает props, но использует:
  - `useAuth()` - получает данные пользователя из Context
  - `useRole()` - получает роль пользователя из Context

**Фильтрация:**
```tsx
async function load() {
  setLoading(true)
  try {
    const data = await getServices(query)  // Передаем query в API
    setServices(data)  // Обновляем состояние
  } catch (err) {
    setError(...)
  } finally {
    setLoading(false)
  }
}

async function handleSearch(e: React.FormEvent) {
  e.preventDefault()
  await load()  // Вызываем load с текущим query
}
```

**Вызовы fetch:**
- `getServices(query)` - вызывает `fetch('/api/services?query=...')`
- Запрос проксируется через Vite на `http://localhost:3001/api/services?query=...`
- Бэкенд фильтрует данные и возвращает отфильтрованный список

**Жизненный цикл:**
1. Компонент монтируется → `useEffect` вызывает `load()`
2. Пользователь вводит текст → `onChange` обновляет `query`
3. Пользователь нажимает "Найти" → `handleSearch` вызывает `load()` с новым `query`
4. `load()` делает fetch запрос → обновляет `services` → компонент перерисовывается

## Структура проекта

```
spa-deposits/
├── src/
│   ├── api/
│   │   ├── http.ts      # Функции для работы с API
│   │   ├── mock.ts      # Mock данные для демонстрации
│   │   └── types.ts     # TypeScript типы
│   ├── components/
│   │   ├── Layout.tsx   # Основной layout с header
│   │   └── ProtectedRoute.tsx  # Защита роутов
│   ├── context/
│   │   └── AuthContext.tsx  # Context для авторизации
│   ├── pages/
│   │   ├── ServicesPage.tsx  # Список услуг с поиском
│   │   ├── DraftPage.tsx     # Черновик расчета
│   │   ├── RequestsPage.tsx  # Список заявок
│   │   ├── LoginPage.tsx      # Вход
│   │   └── RegisterPage.tsx   # Регистрация
│   └── App.tsx          # Роутинг
├── vite.config.ts       # Конфигурация Vite с proxy
└── package.json
```

## Proxy настройка

В `vite.config.ts`:
```typescript
proxy: {
  '/api': {
    target: 'http://localhost:3001',
    changeOrigin: true,
    secure: false,
  },
  '/images': {
    target: 'http://localhost:3001',
    changeOrigin: true,
    secure: false,
  },
}
```

Это позволяет:
- Запросы к `/api/*` проксируются на `http://localhost:3001/api/*`
- Запросы к `/images/*` проксируются на `http://localhost:3001/images/*`
- Избежать проблем с CORS










