@echo off
echo ========================================
echo ПЕРЕСБОРКА С ИСПРАВЛЕНИЕМ ИЗОБРАЖЕНИЙ
echo ========================================
echo.

cd /d %~dp0

echo Останавливаем контейнер...
docker-compose stop deposits-app

echo Пересобираем контейнер...
docker-compose build --no-cache deposits-app

echo Запускаем контейнер...
docker-compose up -d deposits-app

echo Ждем 15 секунд для инициализации...
timeout /t 15 /nobreak >nul

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Бэкенд пересобран и запущен
echo Проверьте логи: docker-compose logs -f deposits-app
echo.

pause


