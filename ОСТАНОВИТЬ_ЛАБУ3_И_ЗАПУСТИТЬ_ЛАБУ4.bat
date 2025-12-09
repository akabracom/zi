@echo off
echo ========================================
echo ОСТАНОВКА ЛАБЫ 3 И ЗАПУСК ЛАБЫ 4
echo ========================================
echo.

REM ========================================
REM ШАГ 1: ОСТАНОВКА ЛАБЫ 3 (backend)
REM ========================================
echo [1/6] Останавливаем лабу 3 (backend)...
cd /d C:\spa-frontend\backend
docker-compose down
if %ERRORLEVEL% NEQ 0 (
    echo ОШИБКА: Не удалось остановить контейнеры лабы 3
    pause
    exit /b 1
)

echo [2/6] Удаляем контейнеры лабы 3...
docker rm -f deposits-db minio deposits-app adminer 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Предупреждение: Некоторые контейнеры уже удалены
)

echo [3/6] Очищаем неиспользуемые сети...
docker network prune -f >nul

echo.
echo Лаба 3 остановлена!
echo.

REM ========================================
REM ШАГ 2: ЗАПУСК ЛАБЫ 4 (SPA_Dep)
REM ========================================
echo [4/6] Переходим в папку лабы 4...
cd /d C:\spa-frontend\SPA_Dep
if %ERRORLEVEL% NEQ 0 (
    echo ОШИБКА: Папка SPA_Dep не найдена!
    pause
    exit /b 1
)

echo [5/6] Останавливаем старые контейнеры (если есть)...
docker-compose down >nul 2>&1
docker rm -f deposits-db minio deposits-app adminer redis 2>nul

echo [6/6] Запускаем лабу 4...
docker-compose up -d
if %ERRORLEVEL% NEQ 0 (
    echo ОШИБКА: Не удалось запустить контейнеры лабы 4
    pause
    exit /b 1
)

echo.
echo ========================================
echo ОЖИДАНИЕ ИНИЦИАЛИЗАЦИИ (15 секунд)...
echo ========================================
timeout /t 15 /nobreak >nul

echo.
echo ========================================
echo ПРОВЕРКА СТАТУСА КОНТЕЙНЕРОВ
echo ========================================
docker-compose ps

echo.
echo ========================================
echo ГОТОВО!
echo ========================================
echo.
echo Проверьте доступность:
echo - Swagger: http://localhost:3001/swagger/index.html
echo - Adminer: http://localhost:8080
echo - MinIO: http://localhost:9001
echo.
echo Для просмотра логов выполните:
echo docker-compose logs -f deposits-app
echo.
pause


