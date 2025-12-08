@echo off
echo ========================================
echo ПОЛНЫЙ ПЕРЕЗАПУСК DEV SERVER
echo ========================================
echo.

cd /d %~dp0

echo [1/4] Останавливаем все процессы Node.js...
taskkill /F /IM node.exe >nul 2>&1
timeout /t 2 /nobreak >nul

echo [2/4] Очищаем кэш npm...
call npm cache clean --force >nul 2>&1

echo [3/4] Удаляем node_modules и package-lock.json (опционально)...
echo Пропускаем удаление node_modules для экономии времени
REM rmdir /s /q node_modules >nul 2>&1
REM del package-lock.json >nul 2>&1

echo [4/4] Запускаем dev server...
echo.
echo ========================================
echo ГОТОВО! Dev server запускается...
echo ========================================
echo.
echo Откройте в браузере: http://localhost:5173
echo.
echo Если изменения не видны:
echo 1. Нажмите Ctrl+Shift+R (жесткая перезагрузка)
echo 2. Или Ctrl+F5
echo.

npm run dev


