@echo off
cd /d %~dp0
taskkill /F /IM node.exe >nul 2>&1
timeout /t 2 /nobreak >nul
echo Запускаем dev server...
npm run dev


