@echo off
echo ========================================
echo ИСПРАВЛЕНИЕ ПАРОЛЯ ПОЛЬЗОВАТЕЛЯ
echo ========================================
echo.

cd /d %~dp0

echo Вариант 1: Удалить старого пользователя и создать нового через API
echo.
echo Вариант 2: Обновить пароль в БД правильным хешем
echo.

echo Выполняем SQL для обновления пароля...
echo Пароль: password123
echo.

docker exec -it deposits-db psql -U user -d deposits -c "UPDATE users SET password = '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy' WHERE email = 'user@example.com';"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo === ПАРОЛЬ ОБНОВЛЕН! ===
    echo Теперь попробуйте войти:
    echo Email: user@example.com
    echo Password: password123
) else (
    echo.
    echo ОШИБКА при обновлении пароля
    echo Попробуйте через Adminer вручную
)

echo.
pause


