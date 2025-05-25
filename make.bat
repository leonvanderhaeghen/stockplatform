@echo off
REM A batch file wrapper for make commands on Windows

if "%~1"=="" (
    echo Usage: make ^<command^>
    echo.
    echo Available commands:
    echo   make deps       - Install dependencies
    echo   make generate   - Generate code from proto files
    echo   make build      - Build all services
    echo   make test       - Run tests
    echo   make lint       - Run linters
    echo   make clean      - Clean build artifacts
    echo   make build-clients - Build all client applications
    echo   make run-grpc-client - Run the gRPC client
    echo   make run-grpc-client-enhanced - Run the enhanced gRPC client
    echo   make dev        - Start the development environment
    echo   make dev-stop   - Stop the development environment
    echo   make dev-logs   - View logs from the development environment
    exit /b 1
)

setlocal enabledelayedexpansion

set COMMAND=%~1
shift

:: Execute the command using the PowerShell script
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0make.ps1" %COMMAND% %*

if !ERRORLEVEL! NEQ 0 (
    echo Error: Command failed with error code !ERRORLEVEL!
    exit /b !ERRORLEVEL!
)
