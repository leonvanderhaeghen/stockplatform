# ----- Comprehensive API Smoke Test -----
$ErrorActionPreference = "Stop"

# Configuration
$base = "http://localhost:8080/api/v1"
$testResults = @{}
$totalTests = 0
$passedTests = 0

# Helper function to make API calls and track results
function Test-ApiCall {
    param(
        [string]$TestName,
        [string]$Method,
        [string]$Uri,
        [hashtable]$Headers = @{},
        [string]$Body = $null,
        [int[]]$ExpectedStatusCodes = @(200, 201)
    )
    
    $global:totalTests++
    Write-Host "Testing: $TestName" -ForegroundColor Cyan
    
    try {
        $params = @{
            Method = $Method
            Uri = $Uri
            Headers = $Headers
        }
        
        if ($Body) {
            $params.Body = $Body
        }
        
        $response = Invoke-RestMethod @params
        
        if ($response) {
            Write-Host "  PASS: $TestName" -ForegroundColor Green
            $global:passedTests++
            $global:testResults[$TestName] = @{ Status = "PASS"; Response = $response }
            return $response
        }
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -in $ExpectedStatusCodes) {
            Write-Host "  PASS: $TestName (Expected status: $statusCode)" -ForegroundColor Green
            $global:passedTests++
            $global:testResults[$TestName] = @{ Status = "PASS"; StatusCode = $statusCode }
            return $null
        } else {
            Write-Host "  FAIL: $TestName - $($_.Exception.Message)" -ForegroundColor Red
            $global:testResults[$TestName] = @{ Status = "FAIL"; Error = $_.Exception.Message }
            return $null
        }
    }
}

Write-Host "=== Starting Comprehensive API Smoke Test ===" -ForegroundColor Yellow
Write-Host "Base URL: $base`n"

# ============= HEALTH CHECK =============
Write-Host "=== Health Check ===" -ForegroundColor Magenta
Test-ApiCall "Health Check" "GET" "$base/health"

# ============= AUTHENTICATION TESTS =============
Write-Host "`n=== Authentication Tests ===" -ForegroundColor Magenta

# Test user registration (correct format: firstName, lastName instead of name)
$registerBody = @{
    firstName = "Test"
    lastName = "User"
    email = "testuser@example.com"
    password = "TestPassword123!"
    phone = "+1234567890"
} | ConvertTo-Json

$registerResponse = Test-ApiCall "User Registration" "POST" "$base/auth/register" @{"Content-Type" = "application/json"} $registerBody

# Test user login
$loginBody = @{
    email = "testuser@example.com"
    password = "TestPassword123!"
} | ConvertTo-Json

$loginResponse = Test-ApiCall "User Login" "POST" "$base/auth/login" @{"Content-Type" = "application/json"} $loginBody

# Extract JWT token for authenticated requests
$userToken = $null
if ($loginResponse -and $loginResponse.token) {
    $userToken = $loginResponse.token
    Write-Host "  User JWT Token obtained" -ForegroundColor Green
}

# Use admin token for admin operations
$adminToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluM0BleGFtcGxlLmNvbSIsImV4cCI6MTc1MDYwMDE5NSwibmFtZSI6IlBsYXRmb3JtIEFkbWluIiwicm9sZSI6IkFETUlOIiwic3ViIjoiYzUyNDFmMzgtMTg4Ni00NzNiLTkxNmEtYWE2OTdlY2MzYjcwIn0.deEW39o5s7U6XkrCPRfVVo8B2f6kD0JxUIBp8IHdotw"

$userHeaders = @{ 
    Authorization = "Bearer $userToken"
    "Content-Type" = "application/json" 
}

$adminHeaders = @{ 
    Authorization = "Bearer $adminToken"
    "Content-Type" = "application/json" 
}

# ============= USER MANAGEMENT TESTS =============
Write-Host "`n=== User Management Tests ===" -ForegroundColor Magenta

if ($userToken) {
    # Get current user profile
    $userProfile = Test-ApiCall "Get User Profile" "GET" "$base/users/me" $userHeaders
    $userId = $null
    if ($userProfile -and $userProfile.id) {
        $userId = $userProfile.id
    }
    
    # Update user profile (correct format: firstName, lastName)
    $updateProfileBody = @{
        firstName = "Updated Test"
        lastName = "User"
        phone = "+1987654321"
    } | ConvertTo-Json
    
    Test-ApiCall "Update User Profile" "PUT" "$base/users/me" $userHeaders $updateProfileBody
    
    # Change password (correct endpoint: /users/me/password)
    $changePasswordBody = @{
        currentPassword = "TestPassword123!"
        newPassword = "NewPassword123!"
    } | ConvertTo-Json
    
    Test-ApiCall "Change Password" "PUT" "$base/users/me/password" $userHeaders $changePasswordBody
    
    # Create user address
    $createAddressBody = @{
        street = "123 Test Street"
        city = "Test City"
        state = "TS"
        zipCode = "12345"
        country = "USA"
        isDefault = $true
    } | ConvertTo-Json
    
    $addressResponse = Test-ApiCall "Create User Address" "POST" "$base/users/me/addresses" $userHeaders $createAddressBody
    
    if ($addressResponse -and $addressResponse.id) {
        $addressId = $addressResponse.id
        
        # Get user addresses
        Test-ApiCall "Get User Addresses" "GET" "$base/users/me/addresses" $userHeaders
        
        # Update address
        $updateAddressBody = @{
            street = "456 Updated Street"
            city = "Updated City"
        } | ConvertTo-Json
        
        Test-ApiCall "Update User Address" "PUT" "$base/users/me/addresses/$addressId" $userHeaders $updateAddressBody
        
        # Set default address
        Test-ApiCall "Set Default Address" "PUT" "$base/users/me/addresses/$addressId/default" $userHeaders
        
        # Delete address (create another one first)
        $deleteAddressBody = @{
            street = "789 Delete Street"
            city = "Delete City"
            state = "DS"
            zipCode = "54321"
            country = "USA"
        } | ConvertTo-Json
        
        $deleteAddressResponse = Test-ApiCall "Create Address to Delete" "POST" "$base/users/me/addresses" $userHeaders $deleteAddressBody
        if ($deleteAddressResponse -and $deleteAddressResponse.id) {
            Test-ApiCall "Delete User Address" "DELETE" "$base/users/me/addresses/$($deleteAddressResponse.id)" $userHeaders
        }
    }
}

# Admin user management tests (correct endpoint: /admin/users)
Write-Host "`n=== Admin User Management Tests ===" -ForegroundColor Magenta

# List all users (admin only) - correct endpoint
$usersResponse = Test-ApiCall "List All Users (Admin)" "GET" "$base/admin/users" $adminHeaders

if ($userId) {
    # Deactivate user (correct endpoint: PUT /admin/users/{id}/deactivate)
    Test-ApiCall "Deactivate User (Admin)" "PUT" "$base/admin/users/$userId/deactivate" $adminHeaders
    
    # Activate user (correct endpoint: PUT /admin/users/{id}/activate)
    Test-ApiCall "Activate User (Admin)" "PUT" "$base/admin/users/$userId/activate" $adminHeaders
}

# ============= PRODUCT MANAGEMENT TESTS =============
Write-Host "`n=== Product Management Tests ===" -ForegroundColor Magenta

# List categories (public endpoint)
$categoriesResponse = Test-ApiCall "List Product Categories" "GET" "$base/products/categories"

# Create category (requires admin/staff auth)
$createCategoryBody = @{
    name = "Test Electronics"
    description = "Test category for electronics"
} | ConvertTo-Json

$categoryResponse = Test-ApiCall "Create Product Category" "POST" "$base/products/categories" $adminHeaders $createCategoryBody
$categoryId = $null
if ($categoryResponse -and $categoryResponse.id) {
    $categoryId = $categoryResponse.id
}

# List products (public endpoint)
$productsResponse = Test-ApiCall "List Products" "GET" "$base/products"

# Create product (requires admin/staff auth)
if ($categoryId) {
    $createProductBody = @{
        name = "Test Wireless Headphones"
        description = "High-quality test headphones"
        price = 199.99
        categoryId = $categoryId
        sku = "TWH-001"
        brand = "TestBrand"
    } | ConvertTo-Json

    $productResponse = Test-ApiCall "Create Product" "POST" "$base/products" $adminHeaders $createProductBody
    $productId = $null
    if ($productResponse -and $productResponse.id) {
        $productId = $productResponse.id
        
        # Get product by ID
        Test-ApiCall "Get Product by ID" "GET" "$base/products/$productId"
        
        # Update product
        $updateProductBody = @{
            name = "Updated Test Headphones"
            price = 249.99
        } | ConvertTo-Json

        Test-ApiCall "Update Product" "PUT" "$base/products/$productId" $adminHeaders $updateProductBody
    }
}

# ============= SUPPLIER MANAGEMENT TESTS =============
Write-Host "`n=== Supplier Management Tests ===" -ForegroundColor Magenta

# Create supplier
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

$supplierResponse = Test-ApiCall "Create Supplier" "POST" "$base/suppliers" $adminHeaders $createSupplierBody
$supplierId = $null
if ($supplierResponse -and $supplierResponse.id) {
    $supplierId = $supplierResponse.id
}

# List suppliers
Test-ApiCall "List Suppliers" "GET" "$base/suppliers" $adminHeaders

if ($supplierId) {
    # Get supplier by ID
    Test-ApiCall "Get Supplier by ID" "GET" "$base/suppliers/$supplierId" $adminHeaders
}

# List adapters
Test-ApiCall "List Supplier Adapters" "GET" "$base/suppliers/adapters" $adminHeaders

# Get adapter capabilities
Test-ApiCall "Get Adapter Capabilities" "GET" "$base/suppliers/adapters/generic/capabilities" $adminHeaders

if ($supplierId) {
    # Test adapter connection
    Test-ApiCall "Test Adapter Connection" "POST" "$base/suppliers/$supplierId/test-connection" $adminHeaders
}

# ============= INVENTORY MANAGEMENT TESTS =============
Write-Host "`n=== Inventory Management Tests ===" -ForegroundColor Magenta

# List inventory
Test-ApiCall "List Inventory" "GET" "$base/inventory" $adminHeaders

if ($productId) {
    # Get inventory by product ID (correct endpoint: /inventory/product/{productId})
    Test-ApiCall "Get Inventory by Product" "GET" "$base/inventory/product/$productId" $adminHeaders
}

# ============= ORDER MANAGEMENT TESTS =============
Write-Host "`n=== Order Management Tests ===" -ForegroundColor Magenta

if ($userToken -and $productId) {
    # Create order
    $createOrderBody = @{
        items = @(
            @{
                productId = $productId
                quantity = 2
                price = 199.99
            }
        )
        shippingAddress = @{
            street = "123 Test Street"
            city = "Test City"
            state = "TS"
            zipCode = "12345"
            country = "USA"
        }
        notes = "Test order"
    } | ConvertTo-Json -Depth 3
    
    $orderResponse = Test-ApiCall "Create Order" "POST" "$base/orders" $userHeaders $createOrderBody
    $orderId = $null
    if ($orderResponse -and $orderResponse.id) {
        $orderId = $orderResponse.id
        
        # List user orders
        Test-ApiCall "List User Orders" "GET" "$base/orders" $userHeaders
        
        # Get order by ID
        Test-ApiCall "Get Order by ID" "GET" "$base/orders/$orderId" $userHeaders
        
        # Update order
        $updateOrderBody = @{
            notes = "Updated test order"
        } | ConvertTo-Json
        
        Test-ApiCall "Update Order" "PUT" "$base/orders/$orderId" $userHeaders $updateOrderBody
        
        # Cancel order
        Test-ApiCall "Cancel Order" "PUT" "$base/orders/$orderId/cancel" $userHeaders
    }
}

# ============= FINAL SUMMARY =============
Write-Host "`n=== Test Summary ===" -ForegroundColor Yellow

# Get final counts
$finalCounts = @{}
try {
    $categoriesResp = Invoke-RestMethod -Uri "$base/products/categories"
    if ($categoriesResp.data -and $categoriesResp.data.total) {
        $finalCounts.Categories = $categoriesResp.data.total
    } else {
        $finalCounts.Categories = $categoriesResp.total
    }
} catch { $finalCounts.Categories = "Error" }

try {
    $productsResp = Invoke-RestMethod -Uri "$base/products"
    if ($productsResp.data -and $productsResp.data.total) {
        $finalCounts.Products = $productsResp.data.total
    } else {
        $finalCounts.Products = $productsResp.total
    }
} catch { $finalCounts.Products = "Error" }

try {
    $suppliersResp = Invoke-RestMethod -Uri "$base/suppliers" -Headers $adminHeaders
    if ($suppliersResp.data -and $suppliersResp.data.total) {
        $finalCounts.Suppliers = $suppliersResp.data.total
    } else {
        $finalCounts.Suppliers = $suppliersResp.total
    }
} catch { $finalCounts.Suppliers = "Error" }

try {
    if ($userToken) {
        $ordersResp = Invoke-RestMethod -Uri "$base/orders" -Headers $userHeaders
        if ($ordersResp.data -and $ordersResp.data.total) {
            $finalCounts.Orders = $ordersResp.data.total
        } else {
            $finalCounts.Orders = $ordersResp.total
        }
    } else {
        $finalCounts.Orders = "No user token"
    }
} catch { $finalCounts.Orders = "Error" }

Write-Host "`nFinal Entity Counts:" -ForegroundColor Cyan
$finalCounts.GetEnumerator() | ForEach-Object { 
    Write-Host "  $($_.Key): $($_.Value)" -ForegroundColor White
}

Write-Host "`nTest Results Summary:" -ForegroundColor Cyan
Write-Host "  Total Tests: $totalTests" -ForegroundColor White
Write-Host "  Passed: $passedTests" -ForegroundColor Green
Write-Host "  Failed: $($totalTests - $passedTests)" -ForegroundColor Red
Write-Host "  Success Rate: $([math]::Round(($passedTests / $totalTests) * 100, 2))%" -ForegroundColor Yellow

if ($passedTests -eq $totalTests) {
    Write-Host "`n ALL TESTS PASSED! API is functioning correctly." -ForegroundColor Green
} else {
    Write-Host "`n  Some tests failed. Check the output above for details." -ForegroundColor Yellow
    
    Write-Host "`nFailed Tests:" -ForegroundColor Red
    $testResults.GetEnumerator() | Where-Object { $_.Value.Status -eq "FAIL" } | ForEach-Object {
        Write-Host "  FAIL: $($_.Key): $($_.Value.Error)" -ForegroundColor Red
    }
}

Write-Host "`n=== Comprehensive API Smoke Test Completed ===" -ForegroundColor Yellow
# ----- End of Comprehensive API Smoke Test -----