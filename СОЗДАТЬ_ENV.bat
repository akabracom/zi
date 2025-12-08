@echo off
echo ========================================
echo СОЗДАНИЕ .env ФАЙЛА
echo ========================================
echo.

cd /d %~dp0

echo VITE_USE_MOCK=true > .env

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Создан файл .env с VITE_USE_MOCK=true
echo.
pause

