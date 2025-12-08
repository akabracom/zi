@echo off
echo ========================================
echo БЫСТРЫЙ ПЕРЕЗАПУСК С БЭКЕНДОМ
echo ========================================
echo.

cd /d %~dp0

echo Останавливаем старые процессы...
taskkill /F /IM node.exe >nul 2>&1
timeout /t 2 /nobreak >nul

echo Запускаем dev server...
echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Откройте: http://localhost:5173
echo.
echo ВАЖНО: Нажмите Ctrl+Shift+R в браузере!
echo.
echo Запускаем...
echo.

npm run dev


