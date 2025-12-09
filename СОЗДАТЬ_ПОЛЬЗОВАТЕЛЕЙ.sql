-- Создание тестовых пользователей для лабы 4
-- Выполнить в Adminer или через docker exec

-- Проверка существующих пользователей
SELECT id, email, name, role, created_at FROM users;

-- Создание пользователей
-- Пароль для всех: "password123"
-- Хеш сгенерирован с помощью bcrypt (cost=10)
INSERT INTO users (email, name, password, role) VALUES
  ('user@example.com', 'Test User', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user'),
  ('moderator@example.com', 'Test Moderator', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'moderator')
ON CONFLICT (email) DO NOTHING;

-- Проверка результата
SELECT id, email, name, role FROM users;


