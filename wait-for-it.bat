@echo off
setlocal

set "host=%1"
set "port=%2"
set "timeout=30"

:loop
timeout /t 1 >nul
echo Checking %host%:%port%...
nslookup %host% >nul 2>&1 && (
    for /f "tokens=*" %%i in ('powershell -command "(New-Object System.Net.Sockets.TcpClient).Connect('%host%', %port%)"') do (
        echo Connected to %host%:%port%
        exit /b 0
    )
)

set /a timeout-=1
if %timeout% gtr 0 goto loop

echo Timeout waiting for %host%:%port%
exit /b 1