#!/bin/bash

# Create necessary directories
mkdir -p pkg/gen/go/product/v1

# Generate protobuf files
echo "Generating protobuf files..."
protoc -I=./proto \
  --go_out=./pkg/gen/go/product/v1 \
  --go_opt=paths=source_relative \
  --go-grpc_out=./pkg/gen/go/product/v1 \
  --go-grpc_opt=paths=source_relative \
  ./proto/product/v1/product.proto

# Fix import paths in generated files
echo "Fixing import paths..."
find ./pkg/gen -type f -name "*.go" -exec sed -i '' 's/"stockplatform\//"github.com\/leonvanderhaeghen\/stockplatform\//g' {} \;

echo "Setup complete!"
