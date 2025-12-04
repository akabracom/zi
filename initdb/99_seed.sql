-- initdb/99_seed.sql
INSERT INTO deposit_offers (
  name, days, interest_rate_pct, start_date, end_date, season, category, holidays, image_key, description, is_active
) VALUES
  ('Январь 2026', 31, 10.50, '2026-01-01', '2026-01-31', 'Зима', 'Сезонный', 'Новогодние', 'january.png', 'Высокая ставка в праздники', TRUE),
  ('Февраль 2026', 28, 10.20, '2026-02-01', '2026-02-28', 'Зима', 'Сезонный', '23 февраля', 'february.png', 'Подарок защитникам', TRUE),
  ('Март 2026', 31, 10.00, '2026-03-01', '2026-03-31', 'Весна', 'Сезонный', '8 марта', 'march.png', 'Женский день', TRUE)
ON CONFLICT (name) DO NOTHING;