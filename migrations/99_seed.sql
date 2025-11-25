INSERT INTO users (email, name, password, role) VALUES
  ('student@example.com', 'Student User', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user')
ON CONFLICT (email) DO NOTHING;

INSERT INTO services (name, description, price) VALUES
  ('Консультация', 'Консультация специалиста', 1000.00),
  ('Установка', 'Установка оборудования', 5000.00),
  ('Обслуживание', 'Техническое обслуживание', 3000.00)
ON CONFLICT (name) DO NOTHING;
