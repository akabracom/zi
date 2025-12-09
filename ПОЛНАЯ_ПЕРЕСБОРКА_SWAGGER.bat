@echo off
echo ========================================
echo ПОЛНАЯ ПЕРЕСБОРКА С ИСПРАВЛЕННЫМ SWAGGER
echo ========================================
echo.

cd /d %~dp0

echo [1/5] Останавливаем контейнеры...
docker-compose stop deposits-app

echo.
echo [2/5] Удаляем старый образ (если есть)...
docker rmi spa_dep-deposits-app 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Образ не найден или уже удален
)

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
echo ПРОВЕРКА ФАЙЛОВ В КОНТЕЙНЕРЕ
echo ========================================
echo.
echo Проверяем наличие securityDefinitions в swagger.json...
docker exec -it deposits-app cat /app/docs/swagger.json | findstr "securityDefinitions"
if %ERRORLEVEL% EQU 0 (
    echo.
    echo === УСПЕХ! securityDefinitions найдены ===
) else (
    echo.
    echo === ПРЕДУПРЕЖДЕНИЕ: securityDefinitions не найдены ===
    echo Проверьте файлы вручную
)

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo ВАЖНО: Очистите кэш браузера!
echo 1. Откройте Swagger в режиме инкогнито (Ctrl+Shift+N)
echo 2. Или нажмите Ctrl+Shift+R для жесткой перезагрузки
echo.
echo URL: http://localhost:3001/swagger/index.html
echo.
echo Должна появиться кнопка "Authorize" вверху справа!
echo.
pause


