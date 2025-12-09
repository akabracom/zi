@echo off
echo ========================================
echo ДИАГНОСТИКА ПРОБЛЕМ С АУТЕНТИФИКАЦИЕЙ
echo ========================================
echo.

cd /d %~dp0

echo [1/5] Проверка логов приложения...
echo.
docker-compose logs deposits-app --tail=30
echo.

echo [2/5] Проверка статуса контейнеров...
docker-compose ps
echo.

echo [3/5] Проверка Redis...
docker exec -it redis redis-cli PING
echo.

echo [4/5] Проверка пользователей в БД...
docker exec -it deposits-db psql -U user -d deposits -c "SELECT id, email, name, role FROM users;"
echo.

echo [5/5] Проверка подключения к БД из приложения...
docker exec -it deposits-app sh -c "echo 'SELECT COUNT(*) FROM users;' | psql -h deposits-db -U user -d deposits" 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Предупреждение: Не удалось подключиться к БД из контейнера
) else (
    echo Подключение к БД работает!
)
echo.

echo ========================================
echo ДИАГНОСТИКА ЗАВЕРШЕНА
echo ========================================
echo.
echo Если видите ошибки выше, отправьте их для анализа.
echo.
pause


