@echo off
echo ========================================
echo ПОЛНЫЙ ПЕРЕЗАПУСК С БЭКЕНДОМ
echo ========================================
echo.

cd /d %~dp0

echo [1/5] Удаляем .env (выключаем mock)...
if exist .env (
    del .env
    echo Файл .env удален
) else (
    echo Файл .env не найден
)
echo.

echo [2/5] Останавливаем старые процессы Node.js...
taskkill /F /IM node.exe >nul 2>&1
timeout /t 2 /nobreak >nul
echo.

echo [3/5] Проверяем бэкенд...
cd ..\SPA_Dep
docker-compose ps deposits-app | findstr "Up" >nul
if %ERRORLEVEL% EQU 0 (
    echo Бэкенд уже запущен!
) else (
    echo Бэкенд не запущен. Запускаем...
    docker-compose up -d deposits-app
    echo Ждем 15 секунд для инициализации...
    timeout /t 15 /nobreak >nul
)
echo.

echo [4/5] Обновляем ставки в БД...
docker exec -i deposits-db psql -U user -d deposits -c "UPDATE services SET price = 10.50, description = 'Дней: 31 | Ставка: 10.5%% | Сезон: Зима | Праздничные дни: 8' WHERE name = 'Январь';" >nul 2>&1
docker exec -i deposits-db psql -U user -d deposits -c "UPDATE services SET price = 10.20, description = 'Дней: 28 | Ставка: 10.2%% | Сезон: Зима | Праздничные дни: 1' WHERE name = 'Февраль';" >nul 2>&1
docker exec -i deposits-db psql -U user -d deposits -c "UPDATE services SET price = 10.00, description = 'Дней: 31 | Ставка: 10.0%% | Сезон: Весна | Праздничные дни: 1' WHERE name = 'Март';" >nul 2>&1
echo Ставки обновлены в БД
echo.

echo [5/5] Возвращаемся в spa-deposits и запускаем dev server...
cd ..\spa-deposits

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Откройте: http://localhost:5173
echo Бэкенд: http://localhost:3001
echo.
echo ВАЖНО: Нажмите Ctrl+Shift+R в браузере для жесткой перезагрузки!
echo.
echo Запускаем dev server...
echo.

npm run dev


