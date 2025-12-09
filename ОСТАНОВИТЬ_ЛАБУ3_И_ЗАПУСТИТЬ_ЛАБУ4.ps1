# Остановка лабы 3 и запуск лабы 4 (PowerShell)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "ОСТАНОВКА ЛАБЫ 3 И ЗАПУСК ЛАБЫ 4" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# ========================================
# ШАГ 1: ОСТАНОВКА ЛАБЫ 3 (backend)
# ========================================
Write-Host "[1/6] Останавливаем лабу 3 (backend)..." -ForegroundColor Yellow
Set-Location C:\spa-frontend\backend
docker-compose down
if ($LASTEXITCODE -ne 0) {
    Write-Host "ОШИБКА: Не удалось остановить контейнеры лабы 3" -ForegroundColor Red
    pause
    exit 1
}

Write-Host "[2/6] Удаляем контейнеры лабы 3..." -ForegroundColor Yellow
docker rm -f deposits-db minio deposits-app adminer 2>$null

Write-Host "[3/6] Очищаем неиспользуемые сети..." -ForegroundColor Yellow
docker network prune -f | Out-Null

Write-Host ""
Write-Host "Лаба 3 остановлена!" -ForegroundColor Green
Write-Host ""

# ========================================
# ШАГ 2: ЗАПУСК ЛАБЫ 4 (SPA_Dep)
# ========================================
Write-Host "[4/6] Переходим в папку лабы 4..." -ForegroundColor Yellow
Set-Location C:\spa-frontend\SPA_Dep
if (-not (Test-Path .)) {
    Write-Host "ОШИБКА: Папка SPA_Dep не найдена!" -ForegroundColor Red
    pause
    exit 1
}

Write-Host "[5/6] Останавливаем старые контейнеры (если есть)..." -ForegroundColor Yellow
docker-compose down 2>&1 | Out-Null
docker rm -f deposits-db minio deposits-app adminer redis 2>$null

Write-Host "[6/6] Запускаем лабу 4..." -ForegroundColor Yellow
docker-compose up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "ОШИБКА: Не удалось запустить контейнеры лабы 4" -ForegroundColor Red
    pause
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "ОЖИДАНИЕ ИНИЦИАЛИЗАЦИИ (15 секунд)..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Start-Sleep -Seconds 15

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "ПРОВЕРКА СТАТУСА КОНТЕЙНЕРОВ" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
docker-compose ps

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "ГОТОВО!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Проверьте доступность:" -ForegroundColor Cyan
Write-Host "- Swagger: http://localhost:3001/swagger/index.html" -ForegroundColor White
Write-Host "- Adminer: http://localhost:8080" -ForegroundColor White
Write-Host "- MinIO: http://localhost:9001" -ForegroundColor White
Write-Host ""
Write-Host "Для просмотра логов выполните:" -ForegroundColor Cyan
Write-Host "docker-compose logs -f deposits-app" -ForegroundColor White
Write-Host ""


