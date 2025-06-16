# Integration Test Runner Script
# This script starts the Docker Compose stack and runs integration tests

param(
    [switch]$SkipBuild,
    [switch]$KeepRunning
)

Write-Host "🚀 Starting Stock Platform Integration Tests" -ForegroundColor Green

# Check if Docker Compose is available
if (-not (Get-Command docker-compose -ErrorAction SilentlyContinue)) {
    Write-Error "Docker Compose is not installed or not in PATH"
    exit 1
}

# Navigate to project root
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

try {
    # Start Docker Compose stack
    Write-Host "📦 Starting Docker Compose stack..." -ForegroundColor Yellow
    if ($SkipBuild) {
        docker-compose up -d
    } else {
        docker-compose up -d --build
    }
    
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to start Docker Compose stack"
    }

    # Wait for services to be ready
    Write-Host "⏳ Waiting for services to be ready..." -ForegroundColor Yellow
    Start-Sleep -Seconds 30

    # Check service health
    Write-Host "🔍 Checking service health..." -ForegroundColor Yellow
    $services = @(
        @{Name="Product Service"; Port=50053},
        @{Name="Inventory Service"; Port=50054},
        @{Name="Order Service"; Port=50055},
        @{Name="User Service"; Port=50056},
        @{Name="Supplier Service"; Port=50057}
    )

    foreach ($service in $services) {
        $maxRetries = 10
        $retryCount = 0
        $isHealthy = $false
        
        while ($retryCount -lt $maxRetries -and -not $isHealthy) {
            try {
                $result = Test-NetConnection -ComputerName localhost -Port $service.Port -WarningAction SilentlyContinue
                if ($result.TcpTestSucceeded) {
                    Write-Host "✅ $($service.Name) is ready on port $($service.Port)" -ForegroundColor Green
                    $isHealthy = $true
                } else {
                    throw "Port not accessible"
                }
            } catch {
                $retryCount++
                Write-Host "⏳ Waiting for $($service.Name) (attempt $retryCount/$maxRetries)..." -ForegroundColor Yellow
                Start-Sleep -Seconds 5
            }
        }
        
        if (-not $isHealthy) {
            throw "$($service.Name) failed to start after $maxRetries attempts"
        }
    }

    # Run integration tests
    Write-Host "🧪 Running integration tests..." -ForegroundColor Yellow
    Set-Location "tests/integration"
    
    # Initialize Go module if needed
    if (-not (Test-Path "go.sum")) {
        go mod tidy
    }
    
    # Run tests with verbose output
    go test -v ./...
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ All integration tests passed!" -ForegroundColor Green
    } else {
        Write-Host "❌ Some integration tests failed" -ForegroundColor Red
    }

} catch {
    Write-Error "Integration test execution failed: $_"
    exit 1
} finally {
    Set-Location $projectRoot
    
    if (-not $KeepRunning) {
        Write-Host "🛑 Stopping Docker Compose stack..." -ForegroundColor Yellow
        docker-compose down
    } else {
        Write-Host "🔄 Docker Compose stack is still running (use -KeepRunning to stop manually)" -ForegroundColor Cyan
        Write-Host "To stop: docker-compose down" -ForegroundColor Cyan
    }
}

Write-Host "🏁 Integration test run completed" -ForegroundColor Green
