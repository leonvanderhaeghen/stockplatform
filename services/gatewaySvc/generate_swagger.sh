#!/bin/bash

# Install swag if not already installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate swagger docs
echo "Generating Swagger documentation..."
swag init -g ./cmd/main.go -o ./docs --outputTypes "go,json"

echo "Swagger documentation generated successfully!"
echo "You can view the documentation at http://localhost:8080/swagger/index.html"
