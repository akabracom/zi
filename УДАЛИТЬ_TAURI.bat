@echo off
echo ========================================
echo УДАЛЕНИЕ TAURI ИЗ 5 ЛАБЫ
echo ========================================
echo.

cd /d %~dp0

if exist "src-tauri" (
    echo Удаляем папку src-tauri...
    rmdir /s /q "src-tauri"
    echo Папка src-tauri удалена!
) else (
    echo Папка src-tauri не найдена. Возможно, уже удалена.
)

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo 5 лаба теперь без Tauri.
echo.
pause

