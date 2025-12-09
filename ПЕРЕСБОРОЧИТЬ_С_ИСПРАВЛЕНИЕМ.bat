@echo off
echo ========================================
echo ПЕРЕСБОРКА С ИСПРАВЛЕНИЕМ GO MOD
echo ========================================
echo.
echo Исправлен Dockerfile:
echo - Используется GOPROXY=direct (обход proxy)
echo - Отключена проверка GOSUMDB
echo.

cd /d %~dp0

echo [1/3] Останавливаем контейнер...
docker-compose stop deposits-app

echo.
echo [2/3] Пересобираем контейнер (без кэша)...
echo Это может занять несколько минут...
echo.
docker-compose build --no-cache deposits-app

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ========================================
    echo ОШИБКА ПРИ СБОРКЕ!
    echo ========================================
    echo.
    echo Попробуйте:
    echo 1. Проверить интернет-соединение
    echo 2. Использовать VPN (если есть ограничения)
    echo 3. Подождать и попробовать снова
    echo.
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
echo Проверьте логи:
echo docker-compose logs deposits-app
echo.
pause
