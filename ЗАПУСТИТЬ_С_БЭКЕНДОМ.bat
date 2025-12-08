@echo off
echo ========================================
echo ЗАПУСК С БЭКЕНДОМ
echo ========================================
echo.

cd /d %~dp0

echo [1/4] Проверяем бэкенд...
cd ..\SPA_Dep
docker-compose ps deposits-app | findstr "Up" >nul
if %ERRORLEVEL% NEQ 0 (
    echo Бэкенд не запущен. Запускаем...
    docker-compose up -d deposits-app
    echo Ждем 15 секунд для инициализации...
    timeout /t 15 /nobreak >nul
) else (
    echo Бэкенд уже запущен!
)

echo.
echo [2/4] Возвращаемся в spa-deposits...
cd ..\spa-deposits

echo [3/4] Создаем .env с VITE_USE_MOCK=false...
echo VITE_USE_MOCK=false > .env

echo [4/4] Запускаем dev server...
echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Откройте в браузере: http://localhost:5173
echo Бэкенд: http://localhost:3001
echo.
echo Данные будут браться с бэкенда через proxy
echo.
npm run dev


