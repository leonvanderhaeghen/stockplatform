# Debug Test for Specific Failing Endpoints
$ErrorActionPreference = "Stop"

$base = "http://localhost:8080/api/v1"
$adminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluM0BleGFtcGxlLmNvbSIsImV4cCI6MTc1MDYwMDE5NSwibmFtZSI6IlBsYXRmb3JtIEFkbWluIiwicm9sZSI6IkFETUlOIiwic3ViIjoiYzUyNDFmMzgtMTg4Ni00NzNiLTkxNmEtYWE2OTdlY2MzYjcwIn0.deEW39o5s7U6XkrCPRfVVo8B2f6kD0JxUIBp8IHdotw"

$adminHeaders = @{
    Authorization = "Bearer $adminToken"
    "Content-Type" = "application/json"
}

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Uri,
        [hashtable]$Headers = @{},
        [string]$Body = $null
    )
    
    Write-Host "`nTesting: $Name" -ForegroundColor Cyan
    Write-Host "  URL: $Uri" -ForegroundColor Gray
    
    try {
        $params = @{
            Method = $Method
            Uri = $Uri
            Headers = $Headers
        }
        
        if ($Body) {
            $params.Body = $Body
            Write-Host "  Body: $Body" -ForegroundColor Gray
        }
        
        $response = Invoke-RestMethod @params
        Write-Host "  SUCCESS" -ForegroundColor Green
        Write-Host "  Response: $($response | ConvertTo-Json -Depth 2)" -ForegroundColor White
        return $response
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "  FAILED: Status $statusCode" -ForegroundColor Red
        Write-Host "  Error: $($_.Exception.Message)" -ForegroundColor Red
        
        # Try to get more detailed error info
        try {
            $errorResponse = $_.Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($errorResponse)
            $errorBody = $reader.ReadToEnd()
            if ($errorBody) {
                Write-Host "  Error Body: $errorBody" -ForegroundColor Red
            }
        } catch {
            Write-Host "  Could not read error response body" -ForegroundColor Yellow
        }
        return $null
    }
}

Write-Host "=== Debugging Specific Failing Endpoints ===" -ForegroundColor Yellow

# Test 1: List All Users (Admin) - 500 error
Test-Endpoint "List All Users (Admin)" "GET" "$base/admin/users" $adminHeaders

# Test 2: Create Supplier - 400 error (with corrected body)
$createSupplierBody = @{
    name = "Test Supplier Co"
    contactPerson = "John Smith"
    email = "contact@testsupplier.com"
    phone = "+1234567890"
    address = "123 Supplier Street"
    city = "Supplier City"
    state = "SC"
    postalCode = "12345"
    country = "USA"
    currency = "USD"
    leadTimeDays = 7
    paymentTerms = "NET30"
} | ConvertTo-Json

$supplierResponse = Test-Endpoint "Create Supplier" "POST" "$base/suppliers" $adminHeaders $createSupplierBody

# Test 3: List Supplier Adapters - 500 error
Test-Endpoint "List Supplier Adapters" "GET" "$base/suppliers/adapters" $adminHeaders

# Test 4: Get Adapter Capabilities - 404 error
Test-Endpoint "Get Adapter Capabilities" "GET" "$base/suppliers/adapters/generic/capabilities" $adminHeaders

# Test 5: Test Adapter Connection (if supplier was created)
if ($supplierResponse -and $supplierResponse.id) {
    Test-Endpoint "Test Adapter Connection" "POST" "$base/suppliers/$($supplierResponse.id)/test-connection" $adminHeaders
}

Write-Host "`n=== Debug Test Completed ===" -ForegroundColor Yellow
