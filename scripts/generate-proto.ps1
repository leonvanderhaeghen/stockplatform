# Ensure the pkg/gen directory exists
$genDir = "pkg/gen"
if (-not (Test-Path -Path $genDir)) {
    New-Item -ItemType Directory -Path $genDir -Force | Out-Null
}

# Function to generate protobuf code
function Generate-Protobuf {
    param (
        [string]$protoFile,
        [string]$goPackage
    )
    
    Write-Host "Generating code for $protoFile..."
    
    # Ensure the output directory exists
    $outputDir = "pkg/gen/$goPackage"
    if (-not (Test-Path -Path $outputDir)) {
        New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
    }
    
    # Generate the code
    protoc -I=./proto `
        --go_out=./pkg/gen `
        --go_opt=paths=source_relative `
        --go-grpc_out=./pkg/gen `
        --go-grpc_opt=paths=source_relative `
        $protoFile
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to generate code for $protoFile"
        exit 1
    }
}

# Generate code for each service
Generate-Protobuf -protoFile "product/v1/product.proto" -goPackage "product/v1"
Generate-Protobuf -protoFile "user/v1/user.proto" -goPackage "user/v1"
Generate-Protobuf -protoFile "order/v1/order.proto" -goPackage "order/v1"
Generate-Protobuf -protoFile "inventory/v1/inventory.proto" -goPackage "inventory/v1"

Write-Host "All protobuf files generated successfully!" -ForegroundColor Green
