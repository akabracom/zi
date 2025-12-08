@echo off
echo ========================================
echo ОЧИСТКА КЭША TAURI
echo ========================================
echo.

cd /d %~dp0

if exist "src-tauri\target" (
    echo Удаляем папку src-tauri\target...
    rmdir /s /q "src-tauri\target"
    echo Кэш очищен!
) else (
    echo Папка target не найдена.
)

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Теперь запустите: npm run tauri:dev
echo.
pause

