<#
.SYNOPSIS
    A PowerShell wrapper for make commands on Windows
.DESCRIPTION
    This script provides a way to run make commands on Windows by translating them to the appropriate PowerShell commands.
    It's a workaround for Windows systems that don't have GNU Make installed.
.EXAMPLE
    .\make.ps1 build
    .\make.ps1 test
    .\make.ps1 clean
#>

param(
    [Parameter(Mandatory=$true)]
    [string]$Command,
    
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$Arguments
)

$ErrorActionPreference = "Stop"

# Define the commands and their corresponding PowerShell commands
$commands = @{
    "deps" = {
        Write-Host "Installing dependencies..." -ForegroundColor Cyan
        go get -u google.golang.org/protobuf/cmd/protoc-gen-go `
            google.golang.org/grpc/cmd/protoc-gen-go-grpc `
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway `
            github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 `
            github.com/bufbuild/buf/cmd/buf `
            github.com/fullstorydev/grpcurl/cmd/grpcurl
    }
    
    "generate" = {
        Write-Host "Generating code from proto files..." -ForegroundColor Cyan
        buf generate
    }
    
    "build" = {
        Write-Host "Building all services..." -ForegroundColor Cyan
        Get-ChildItem -Path .\services -Directory | ForEach-Object {
            $serviceName = $_.Name
            Write-Host "Building $serviceName" -ForegroundColor Green
            Set-Location $_.FullName
            go build -o "bin\$serviceName.exe" .
            if ($LASTEXITCODE -ne 0) {
                throw "Failed to build $serviceName"
            }
            Write-Host "Build complete for $serviceName" -ForegroundColor Green
            Set-Location $PSScriptRoot
        }
    }
    
    "test" = {
        Write-Host "Running tests..." -ForegroundColor Cyan
        go test -v ./...
    }
    
    "lint" = {
        Write-Host "Running linters..." -ForegroundColor Cyan
        golangci-lint run ./...
    }
    
    "build-clients" = {
        Write-Host "Building all client applications..." -ForegroundColor Cyan
        $binDir = ".\bin"
        if (-not (Test-Path $binDir)) {
            New-Item -ItemType Directory -Path $binDir | Out-Null
        }
        
        Get-ChildItem -Path .\cmd -Directory | ForEach-Object {
            $clientName = $_.Name
            Write-Host "Building $clientName" -ForegroundColor Green
            $output = "$binDir\$clientName.exe"
            go build -o $output ".\cmd\$clientName"
            if ($LASTEXITCODE -ne 0) {
                throw "Failed to build $clientName"
            }
            Write-Host "Build complete for $clientName" -ForegroundColor Green
        }
    }
    
    "run-grpc-client" = {
        Write-Host "Running gRPC client..." -ForegroundColor Cyan
        go build -o ".\bin\grpc_client.exe" .\cmd\grpc_client
        .\bin\grpc_client.exe
    }
    
    "run-grpc-client-enhanced" = {
        Write-Host "Running enhanced gRPC client..." -ForegroundColor Cyan
        go build -o ".\bin\grpc_client_enhanced.exe" .\cmd\grpc_client_enhanced
        .\bin\grpc_client_enhanced.exe
    }
    
    "clean" = {
        Write-Host "Cleaning build artifacts..." -ForegroundColor Cyan
        go clean
        if (Test-Path ".\bin") {
            Remove-Item -Recurse -Force ".\bin"
        }
        
        Get-ChildItem -Path .\services -Directory | ForEach-Object {
            $binDir = "$($_.FullName)\bin"
            if (Test-Path $binDir) {
                Write-Host "Cleaning $($_.Name)" -ForegroundColor Yellow
                Remove-Item -Recurse -Force $binDir
            }
        }
    }
    
    "dev" = {
        Write-Host "Starting development environment..." -ForegroundColor Cyan
        docker-compose up --build -d
    }
    
    "dev-stop" = {
        Write-Host "Stopping development environment..." -ForegroundColor Cyan
        docker-compose down
    }
    
    "dev-logs" = {
        docker-compose logs -f
    }
    
    "help" = {
        Write-Host @"
Available commands:
  .\make.ps1 deps       - Install dependencies
  .\make.ps1 generate   - Generate code from proto files
  .\make.ps1 build      - Build all services
  .\make.ps1 test       - Run tests
  .\make.ps1 lint       - Run linters
  .\make.ps1 clean      - Clean build artifacts
  .\make.ps1 build-clients - Build all client applications
  .\make.ps1 run-grpc-client - Run the gRPC client
  .\make.ps1 run-grpc-client-enhanced - Run the enhanced gRPC client
  .\make.ps1 dev        - Start the development environment
  .\make.ps1 dev-stop   - Stop the development environment
  .\make.ps1 dev-logs   - View logs from the development environment
"@
    }
}

# Main execution
try {
    if ($commands.ContainsKey($Command)) {
        & $commands[$Command] @Arguments
    } else {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host "Use '.\make.ps1 help' to see available commands"
        exit 1
    }
} catch {
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}
