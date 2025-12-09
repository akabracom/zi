-- Исправление паролей для пользователей
-- Выполнить в Adminer

-- Обновить пароль для user@example.com
-- Пароль: password123
-- Хеш: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
UPDATE users 
SET password = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy' 
WHERE email = 'user@example.com';

-- Обновить пароль для moderator@example.com
UPDATE users 
SET password = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy' 
WHERE email = 'moderator@example.com';

-- Проверка результата
SELECT id, email, name, role, 
       LEFT(password, 20) as password_preview,
       LENGTH(password) as password_length
FROM users;


