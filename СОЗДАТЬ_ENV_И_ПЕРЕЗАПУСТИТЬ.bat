@echo off
echo ========================================
echo СОЗДАНИЕ .env И ПЕРЕЗАПУСК
echo ========================================
echo.

cd /d %~dp0

echo [1/3] Создаем файл .env...
echo VITE_USE_MOCK=true > .env
echo.
echo Проверяем содержимое:
type .env
echo.

echo [2/3] Останавливаем процессы Node.js...
taskkill /F /IM node.exe >nul 2>&1
timeout /t 2 /nobreak >nul

echo [3/3] Запускаем dev server...
echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Откройте: http://localhost:5173
echo Нажмите Ctrl+Shift+R для жесткой перезагрузки
echo.

npm run dev


