@echo off
REM Build script for Windows

echo Building SerdeVal for Windows...

REM Build the binary
go build -o serdeval.exe ./cmd/serdeval

if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    exit /b 1
)

echo Build successful! Binary created: serdeval.exe
echo.
echo Usage examples:
echo   serdeval.exe validate config.json
echo   serdeval.exe web
echo   type data.json ^| serdeval.exe validate