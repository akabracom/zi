@echo off
echo ========================================
echo НАСТРОЙКА ДЛЯ РАБОТЫ С БЭКЕНДОМ
echo ========================================
echo.

cd /d %~dp0

echo [1/4] Удаляем .env (выключаем mock режим)...
if exist .env (
    del .env
    echo Файл .env удален
) else (
    echo Файл .env не найден (уже удален)
)
echo.

echo [2/4] Проверяем бэкенд...
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

echo [3/4] Возвращаемся в spa-deposits...
cd ..\spa-deposits

echo [4/4] Останавливаем старые процессы Node.js...
taskkill /F /IM node.exe >nul 2>&1
timeout /t 2 /nobreak >nul

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Запускаем dev server...
echo Откройте: http://localhost:5173
echo Бэкенд: http://localhost:3001
echo.
echo Данные будут браться с бэкенда
echo.

npm run dev


