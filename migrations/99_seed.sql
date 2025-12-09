INSERT INTO users (email, name, password, role) VALUES
  ('student@example.com', 'Student User', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user')
ON CONFLICT (email) DO NOTHING;

-- Удаляем старые услуги перед вставкой новых месяцев
DELETE FROM services WHERE name IN ('Консультация', 'Установка', 'Обслуживание', 'Услуга 1');

-- Вставляем банковские месяцы с изображениями
-- Ставки соответствуют 2-й лабе (di_web)
INSERT INTO services (name, description, price, image_url) VALUES
  ('Январь', 'Дней: 31 | Ставка: 10.5% | Сезон: Зима | Праздничные дни: 8', 10.50, 'january.png'),
  ('Февраль', 'Дней: 28 | Ставка: 10.2% | Сезон: Зима | Праздничные дни: 1', 10.20, 'february.png'),
  ('Март', 'Дней: 31 | Ставка: 10.0% | Сезон: Весна | Праздничные дни: 1', 10.00, 'march.png')
ON CONFLICT (name) DO UPDATE SET
  description = EXCLUDED.description,
  price = EXCLUDED.price,
  image_url = EXCLUDED.image_url;
