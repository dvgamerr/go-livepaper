@echo off
setlocal

REM Get version from VERSION file
set /p VERSION=<VERSION
echo Building livepaper version %VERSION%

REM Build with version embedded
go build -ldflags "-X main.VERSION=%VERSION%" -o livepaper.exe

echo Build complete
endlocal