@echo off
echo ========================================
echo БЫСТРЫЙ СТАРТ ДЛЯ ЛАБЫ 5
echo ========================================
echo.

cd /d %~dp0

echo Выберите режим:
echo 1. Mock режим (без бэкенда)
echo 2. Реальный режим (с бэкендом)
echo.
set /p mode="Введите номер (1 или 2): "

if "%mode%"=="1" (
    echo.
    echo [1/2] Создаем .env с VITE_USE_MOCK=true...
    echo VITE_USE_MOCK=true > .env
    
    echo [2/2] Запускаем dev server...
    echo.
    echo Откройте: http://localhost:5173
    echo.
    npm run dev
) else if "%mode%"=="2" (
    echo.
    echo [1/3] Проверяем бэкенд...
    cd ..\SPA_Dep
    docker-compose ps deposits-app | findstr "Up" >nul
    if %ERRORLEVEL% NEQ 0 (
        echo Бэкенд не запущен. Запускаем...
        docker-compose up -d deposits-app
        echo Ждем 10 секунд...
        timeout /t 10 /nobreak >nul
    ) else (
        echo Бэкенд уже запущен!
    )
    
    echo.
    echo [2/3] Создаем .env с VITE_USE_MOCK=false...
    cd ..\spa-deposits
    echo VITE_USE_MOCK=false > .env
    
    echo [3/3] Запускаем dev server...
    echo.
    echo Откройте: http://localhost:5173
    echo Бэкенд: http://localhost:3001
    echo.
    npm run dev
) else (
    echo Неверный выбор!
    pause
    exit /b 1
)


