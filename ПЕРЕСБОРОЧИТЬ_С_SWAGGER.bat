@echo off
echo ========================================
echo ПЕРЕСБОРКА С ИСПРАВЛЕННЫМ SWAGGER
echo ========================================
echo.

cd /d %~dp0

echo [1/3] Останавливаем контейнер...
docker-compose stop deposits-app

echo.
echo [2/3] Пересобираем контейнер (без кэша)...
docker-compose build --no-cache deposits-app

echo.
echo [3/3] Запускаем контейнер...
docker-compose up -d deposits-app

echo.
echo ========================================
echo ОЖИДАНИЕ ЗАПУСКА (10 секунд)...
echo ========================================
timeout /t 10 /nobreak >nul

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Проверьте Swagger: http://localhost:3001/swagger/index.html
echo Должна появиться кнопка "Authorize" вверху справа!
echo.
pause


