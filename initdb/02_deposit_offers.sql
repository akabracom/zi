-- initdb/02_deposit_offers.sql
-- Создаём таблицу deposit_offers
CREATE TABLE IF NOT EXISTS deposit_offers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    days INTEGER NOT NULL,
    interest_rate_pct DOUBLE PRECISION NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    season TEXT,
    category TEXT,
    holidays TEXT,
    image_key TEXT,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE
);

-- Добавляем уникальный индекс ПОСЛЕ создания таблицы
CREATE UNIQUE INDEX IF NOT EXISTS idx_deposit_offers_name ON deposit_offers(name);