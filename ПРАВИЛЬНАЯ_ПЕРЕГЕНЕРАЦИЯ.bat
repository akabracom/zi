@echo off
echo ========================================
echo ПРАВИЛЬНАЯ ПЕРЕГЕНЕРАЦИЯ SWAGGER
echo ========================================
echo.

cd /d %~dp0

echo [1/4] Проверяем наличие swag...
where swag >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo swag не найден. Устанавливаем...
    go install github.com/swaggo/swag/cmd/swag@latest
    if %ERRORLEVEL% NEQ 0 (
        echo ОШИБКА: Не удалось установить swag
        echo Установите вручную: go install github.com/swaggo/swag/cmd/swag@latest
        pause
        exit /b 1
    )
    echo swag установлен!
) else (
    echo swag найден!
)

echo.
echo [2/4] Перегенерируем Swagger документацию...
swag init -g cmd/server/main.go
if %ERRORLEVEL% NEQ 0 (
    echo ОШИБКА при генерации Swagger!
    pause
    exit /b 1
)

echo.
echo [3/4] Проверяем наличие securityDefinitions в swagger.json...
findstr /C:"securityDefinitions" docs\swagger.json >nul
if %ERRORLEVEL% EQU 0 (
    echo === УСПЕХ! securityDefinitions найдены в swagger.json ===
) else (
    echo === ОШИБКА: securityDefinitions НЕ найдены ===
    echo Проверьте комментарии в cmd/server/main.go
    pause
    exit /b 1
)

echo.
echo [4/4] Пересобираем контейнер...
docker-compose stop deposits-app
docker-compose build --no-cache deposits-app
docker-compose up -d deposits-app

echo.
echo ========================================
echo ОЖИДАНИЕ ЗАПУСКА (15 секунд)...
echo ========================================
timeout /t 15 /nobreak >nul

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Проверьте:
echo 1. http://localhost:3001/swagger/doc.json - должен содержать securityDefinitions
echo 2. http://localhost:3001/swagger/index.html - должна быть кнопка "Authorize"
echo.
echo ВАЖНО: Очистите кэш браузера (Ctrl+Shift+R или режим инкогнито)!
echo.
pause


