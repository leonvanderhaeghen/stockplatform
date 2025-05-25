# Product Service (productSvc)

This service manages the product catalog for the Stock Platform, providing CRUD operations for products.

## Features

- CRUD operations for products
- Product categorization and attributes
- Product search and filtering
- Price management
- Product image handling

## Architecture

The Product Service follows a clean architecture pattern:

```
productSvc/
├── cmd/                 # Command-line entry points
├── internal/
│   ├── domain/          # Domain models and interfaces
│   ├── application/     # Business logic
│   ├── infrastructure/  # External dependencies (MongoDB)
│   └── interfaces/      # Interface adapters (gRPC)
├── Dockerfile           # Container definition
├── buf.gen.yaml         # Protocol buffer generation config
└── buf.work.yaml        # Protocol buffer workspace config
```

## API

The service exposes a gRPC API defined in `proto/product/v1/product.proto`:

### Key Endpoints

- `CreateProduct` - Create a new product
- `GetProduct` - Get product details by ID
- `UpdateProduct` - Update an existing product
- `DeleteProduct` - Delete a product
- `ListProducts` - List products with filtering options
- `SearchProducts` - Search products by name, description, or other attributes
- `GetProductsByCategory` - Get products in a specific category

## Domain Model

The core domain entity is:

- **Product**: Represents a product with properties like name, description, SKU, price, categories, and custom attributes

## Configuration

The service can be configured using environment variables:

- `GRPC_PORT` - Port for gRPC server (default: 50053)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)

## Development

### Prerequisites

- Go 1.22+
- MongoDB 7.0+
- Protocol buffer compiler (protoc)
- Buf CLI

### Running Locally

1. Start MongoDB:

```bash
docker-compose up -d mongodb
```

2. Run the service:

```bash
go run cmd/main.go
```

### Testing

Run unit tests:

```bash
go test ./...
```

Run integration tests (requires MongoDB):

```bash
go test -tags=integration ./...
```

## Docker

Build the Docker image:

```bash
docker build -t stockplatform/productsvc .
```

Run the container:

```bash
docker run -p 50053:50053 -e MONGO_URI=mongodb://mongodb:27017 stockplatform/productsvc
```

## Interacting with the Service

You can use gRPCurl to interact with the service:

```bash
# List products
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' localhost:50053 product.v1.ProductService/ListProducts

# Get product by ID
grpcurl -plaintext -d '{"id": "product-id-here"}' localhost:50053 product.v1.ProductService/GetProduct

# Create a new product
grpcurl -plaintext -d '{
  "product": {
    "name": "Test Product",
    "description": "This is a test product",
    "sku": "TEST-001",
    "price": 19.99,
    "cost": 9.99,
    "categories": ["test", "sample"],
    "active": true
  }
}' localhost:50053 product.v1.ProductService/CreateProduct
```

## Data Storage

Products are stored in MongoDB with the following structure:

```json
{
  "_id": "unique-product-id",
  "name": "Product Name",
  "description": "Product description",
  "sku": "PRODUCT-SKU-001",
  "price": 19.99,
  "cost": 9.99,
  "categories": ["category1", "category2"],
  "active": true,
  "images": ["image_url_1", "image_url_2"],
  "attributes": {
    "color": "red",
    "size": "medium",
    "weight": "500g"
  },
  "created_at": "2025-05-25T15:30:00Z",
  "updated_at": "2025-05-25T15:30:00Z"
}
```
