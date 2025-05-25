# Test gRPC service using grpcurl

# Read the request JSON
$json = Get-Content -Raw .\testdata\create_product_request.json

# Call the gRPC service
& grpcurl -plaintext -d $json 127.0.0.1:50053 product.v1.ProductService/CreateProduct
