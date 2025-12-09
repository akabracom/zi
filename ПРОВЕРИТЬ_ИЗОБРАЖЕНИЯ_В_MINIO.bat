@echo off
echo ========================================
echo ПРОВЕРКА ИЗОБРАЖЕНИЙ В MINIO
echo ========================================
echo.

cd /d %~dp0

echo Проверяем, запущен ли MinIO...
docker ps | findstr "minio" >nul
if %ERRORLEVEL% NEQ 0 (
    echo MinIO не запущен. Запускаем...
    docker-compose up -d minio
    echo Ждем 5 секунд для инициализации...
    timeout /t 5 /nobreak >nul
) else (
    echo MinIO запущен!
)
echo.

echo Проверяем содержимое MinIO...
echo.
docker exec -it minio ls -la /data/deposits/img/ 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Папка img/ не найдена или пуста
    echo.
    echo Проверяем корень bucket deposits...
    docker exec -it minio ls -la /data/deposits/ 2>nul
    if %ERRORLEVEL% NEQ 0 (
        echo.
        echo Bucket deposits не найден или пуст
        echo.
        echo Изображения нужно загрузить через API или скопировать вручную
    )
) else (
    echo.
    echo Изображения найдены в MinIO!
)

echo.
echo ========================================
echo ИНСТРУКЦИЯ
echo ========================================
echo.
echo Если изображений нет:
echo 1. Загрузите через API: POST /api/services/{id}/image
echo 2. Или используйте заглушки (они уже работают)
echo.
echo Если изображения есть, но не отображаются:
echo 1. Проверьте логи бэкенда: docker-compose logs deposits-app
echo 2. Проверьте путь в БД: SELECT name, image_url FROM services;
echo 3. Проверьте URL в браузере: http://localhost:3001/images/january.png
echo.

pause


