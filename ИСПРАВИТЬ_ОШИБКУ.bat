@echo off
echo ========================================
echo ИСПРАВЛЕНИЕ ОШИБКИ TAURI
echo ========================================
echo.

cd /d %~dp0

echo [1/3] Останавливаем все процессы...
taskkill /F /IM node.exe >nul 2>&1
taskkill /F /IM cargo.exe >nul 2>&1

echo [2/3] Очищаем кэш Tauri...
if exist "src-tauri\target" (
    echo Удаляем папку src-tauri\target...
    rmdir /s /q "src-tauri\target"
    echo Кэш очищен!
) else (
    echo Папка target не найдена.
)

echo [3/3] Очищаем кэш Cargo...
cd src-tauri
if exist "Cargo.lock" (
    echo Обновляем Cargo.lock...
    cargo clean
)
cd ..

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Теперь запустите: npm run tauri:dev
echo.
pause

