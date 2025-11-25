-- migrations/01_users.sql
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT,
  password TEXT NOT NULL,
  role VARCHAR(20) NOT NULL DEFAULT 'user',
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создаем тестовых пользователей
-- Пароль для всех: "password123"
-- Хеш сгенерирован с помощью bcrypt (cost=10)
INSERT INTO users (email, name, password, role) VALUES
  ('user@example.com', 'Test User', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user'),
  ('moderator@example.com', 'Test Moderator', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'moderator');
