@echo off
echo ========================================
echo БЫСТРЫЙ ПЕРЕЗАПУСК DEV SERVER
echo ========================================
echo.

cd /d %~dp0

echo Останавливаем процессы Node.js...
taskkill /F /IM node.exe >nul 2>&1
timeout /t 1 /nobreak >nul

echo Запускаем dev server...
echo.
echo Откройте: http://localhost:5173
echo.

npm run dev


