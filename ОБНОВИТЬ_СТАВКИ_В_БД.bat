@echo off
echo ========================================
echo ОБНОВЛЕНИЕ СТАВОК В БД
echo ========================================
echo.

cd /d %~dp0

echo Выполняем SQL запросы для обновления ставок...
echo.

docker exec -i deposits-db psql -U user -d deposits -c "UPDATE services SET price = 10.50, description = 'Дней: 31 | Ставка: 10.5%% | Сезон: Зима | Праздничные дни: 8' WHERE name = 'Январь';"

docker exec -i deposits-db psql -U user -d deposits -c "UPDATE services SET price = 10.20, description = 'Дней: 28 | Ставка: 10.2%% | Сезон: Зима | Праздничные дни: 1' WHERE name = 'Февраль';"

docker exec -i deposits-db psql -U user -d deposits -c "UPDATE services SET price = 10.00, description = 'Дней: 31 | Ставка: 10.0%% | Сезон: Весна | Праздничные дни: 1' WHERE name = 'Март';"

echo.
echo Проверяем результат...
docker exec -i deposits-db psql -U user -d deposits -c "SELECT id, name, price, description FROM services WHERE name IN ('Январь', 'Февраль', 'Март');"

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Ставки обновлены в БД
echo.

pause


