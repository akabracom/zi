@echo off
echo ========================================
echo ПЕРЕСБОРКА С ПОДДЕРЖКОЙ COOKIE
echo ========================================
echo.
echo Добавлено:
echo 1. Сохранение JWT токена в cookie при Login/Register
echo 2. Поддержка чтения токена из cookie в middleware
echo.

cd /d %~dp0

echo [1/3] Останавливаем контейнер...
docker-compose stop deposits-app

echo.
echo [2/3] Пересобираем контейнер (без кэша)...
docker-compose build --no-cache deposits-app

if %ERRORLEVEL% NEQ 0 (
    echo ОШИБКА при сборке!
    pause
    exit /b 1
)

echo.
echo [3/3] Запускаем контейнер...
docker-compose up -d deposits-app

echo.
echo ========================================
echo ОЖИДАНИЕ ЗАПУСКА (15 секунд)...
echo ========================================
timeout /t 15 /nobreak >nul

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Теперь:
echo 1. Выполните POST /auth/login в Swagger
echo 2. Откройте DevTools → Application → Cookies
echo 3. Должна появиться cookie "session_token" с JWT токеном!
echo.
pause


