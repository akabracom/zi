@echo off
echo ========================================
echo ЗАПУСК В MOCK РЕЖИМЕ (БЕЗ БЭКЕНДА)
echo ========================================
echo.

cd /d %~dp0

echo [1/2] Создаем .env с VITE_USE_MOCK=true...
echo VITE_USE_MOCK=true > .env

echo [2/2] Запускаем dev server...
echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Откройте в браузере: http://localhost:5173
echo.
echo Данные будут браться из mock.ts (без бэкенда)
echo.
npm run dev


