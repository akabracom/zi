Write-Host "========================================" -ForegroundColor Cyan
Write-Host "ПОЛНЫЙ ПЕРЕЗАПУСК DEV SERVER" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Set-Location $PSScriptRoot

Write-Host "[1/4] Останавливаем все процессы Node.js..." -ForegroundColor Yellow
Get-Process -Name node -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Seconds 2

Write-Host "[2/4] Очищаем кэш npm..." -ForegroundColor Yellow
npm cache clean --force 2>$null

Write-Host "[3/4] Проверяем .env файл..." -ForegroundColor Yellow
if (Test-Path ".env") {
    Write-Host "Файл .env найден" -ForegroundColor Green
    Get-Content .env
} else {
    Write-Host "Файл .env не найден" -ForegroundColor Yellow
}

Write-Host "[4/4] Запускаем dev server..." -ForegroundColor Yellow
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "ГОТОВО! Dev server запускается..." -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Откройте в браузере: http://localhost:5173" -ForegroundColor Green
Write-Host ""
Write-Host "Если изменения не видны:" -ForegroundColor Yellow
Write-Host "1. Нажмите Ctrl+Shift+R (жесткая перезагрузка)" -ForegroundColor Yellow
Write-Host "2. Или Ctrl+F5" -ForegroundColor Yellow
Write-Host ""

npm run dev


