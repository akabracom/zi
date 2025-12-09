@echo off
echo ========================================
echo ПОЛНОЕ ИСПРАВЛЕНИЕ SWAGGER
echo ========================================
echo.
echo Исправлено:
echo 1. Добавлены securityDefinitions в swagger.json
echo 2. Добавлены securityDefinitions в swagger.yaml
echo 3. Исправлен Dockerfile (добавлено копирование docs)
echo.
echo ========================================
echo.

cd /d %~dp0

echo [1/5] Останавливаем контейнеры...
docker-compose stop deposits-app

echo.
echo [2/5] Удаляем старый образ...
docker rmi spa_dep-deposits-app 2>nul

echo.
echo [3/5] Пересобираем контейнер БЕЗ кэша...
docker-compose build --no-cache deposits-app
if %ERRORLEVEL% NEQ 0 (
    echo ОШИБКА при сборке!
    pause
    exit /b 1
)

echo.
echo [4/5] Запускаем контейнер...
docker-compose up -d deposits-app

echo.
echo [5/5] Ожидание запуска (15 секунд)...
timeout /t 15 /nobreak >nul

echo.
echo ========================================
echo ПРОВЕРКА
echo ========================================
echo.
echo Проверяем наличие securityDefinitions...
docker exec -it deposits-app sh -c "cat /root/docs/swagger.json | grep -q securityDefinitions && echo 'OK: securityDefinitions найдены' || echo 'ОШИБКА: securityDefinitions не найдены'"

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo ВАЖНО: Очистите кэш браузера!
echo 1. Откройте в режиме инкогнито: http://localhost:3001/swagger/index.html
echo 2. Или нажмите Ctrl+Shift+R для жесткой перезагрузки
echo.
echo Должна появиться кнопка "Authorize" вверху справа!
echo.
pause


