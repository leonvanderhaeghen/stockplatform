# PowerShell script to generate protobuf and gRPC code

# Set error action preference
$ErrorActionPreference = "Stop"

# Get the root directory of the project
$rootDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$protoDir = Join-Path $rootDir "proto"
$outputDir = Join-Path $rootDir "pkg\gen\go"

# Clean up existing generated files
Write-Host "Cleaning up generated protobuf files..."
if (Test-Path $outputDir) {
    Remove-Item -Recurse -Force $outputDir
}

# Create output directories
New-Item -ItemType Directory -Force -Path $outputDir | Out-Null

# Generate code for each service
$services = @("inventory", "order", "product", "user")

foreach ($service in $services) {
    $serviceProtoDir = Join-Path -Path $protoDir -ChildPath $service
    $serviceVersionDir = Join-Path -Path $serviceProtoDir -ChildPath "v1"
    $protoFile = Join-Path -Path $serviceVersionDir -ChildPath "$service.proto"
    
    if (Test-Path $protoFile) {
        Write-Host "Generating code for $service..."
        
        try {
            # Calculate relative path from proto file to root
            $relativePath = "..\..\.."
            
            # Generate Go code
            & protoc `
                -I="$protoDir" `
                --go_out="paths=source_relative:$outputDir" `
                --go-grpc_out="paths=source_relative:$outputDir" `
                "$protoFile"
                
            if ($LASTEXITCODE -ne 0) {
                throw "protoc failed with exit code $LASTEXITCODE"
            }
            
            Write-Host "Successfully generated code for $service" -ForegroundColor Green
        }
        catch {
            Write-Error "Failed to generate code for $service : $_"
            exit 1
        }
    }
    else {
        Write-Warning "Proto file not found: $protoFile"
    }
}

Write-Host "`nCode generation completed successfully!" -ForegroundColor Green
