-- Обновление ставок в БД для соответствия 2-й лабе
-- Выполните в Adminer: http://localhost:8080

-- Обновляем ставки для месяцев
UPDATE services 
SET price = 10.50,
    description = 'Дней: 31 | Ставка: 10.5% | Сезон: Зима | Праздничные дни: 8'
WHERE name = 'Январь';

UPDATE services 
SET price = 10.20,
    description = 'Дней: 28 | Ставка: 10.2% | Сезон: Зима | Праздничные дни: 1'
WHERE name = 'Февраль';

UPDATE services 
SET price = 10.00,
    description = 'Дней: 31 | Ставка: 10.0% | Сезон: Весна | Праздничные дни: 1'
WHERE name = 'Март';

-- Проверка результата
SELECT id, name, price, description FROM services WHERE name IN ('Январь', 'Февраль', 'Март');


