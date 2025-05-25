# Inventory Service (inventorySvc)

This service manages inventory items and stock levels for the Stock Platform.

## Features

- Real-time stock level tracking
- Low stock alerts
- Stock adjustments (add/remove)
- Multiple locations support
- Reorder point management

## Architecture

The Inventory Service follows a clean architecture pattern:

```
inventorySvc/
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

The service exposes a gRPC API defined in `proto/inventory/v1/inventory.proto`:

### Key Endpoints

- `CreateInventoryItem` - Create a new inventory item
- `GetInventoryItem` - Get inventory item details by ID
- `GetInventoryByProduct` - Get inventory items for a specific product
- `GetInventoryBySku` - Get inventory item by SKU
- `UpdateInventoryItem` - Update an inventory item
- `DeleteInventoryItem` - Delete an inventory item
- `ListInventory` - List inventory items with filtering options
- `AddStock` - Add stock to an inventory item
- `RemoveStock` - Remove stock from an inventory item
- `CheckLowStock` - Check for items with low stock levels

## Configuration

The service can be configured using environment variables:

- `GRPC_PORT` - Port for gRPC server (default: 50054)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `PRODUCT_SERVICE_ADDR` - Product service address (default: localhost:50053)

## Development

### Prerequisites

- Go 1.22+
- MongoDB 7.0+
- Protocol buffer compiler (protoc)
- Buf CLI

### Running Locally

1. Start MongoDB and the Product Service:

```bash
docker-compose up -d mongodb
cd ../productSvc && go run cmd/main.go
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

Run integration tests (requires MongoDB and Product Service):

```bash
go test -tags=integration ./...
```

## Docker

Build the Docker image:

```bash
docker build -t stockplatform/inventorysvc .
```

Run the container:

```bash
docker run -p 50054:50054 \
  -e MONGO_URI=mongodb://mongodb:27017 \
  -e PRODUCT_SERVICE_ADDR=productsvc:50053 \
  stockplatform/inventorysvc
```

## Interacting with the Service

You can use gRPCurl to interact with the service:

```bash
# List inventory items
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' localhost:50054 inventory.v1.InventoryService/ListInventory

# Get inventory item by ID
grpcurl -plaintext -d '{"id": "inventory-id-here"}' localhost:50054 inventory.v1.InventoryService/GetInventoryItem

# Add stock to an inventory item
grpcurl -plaintext -d '{"id": "inventory-id-here", "quantity": 10, "reason": "Restocking", "reference": "PO-12345"}' localhost:50054 inventory.v1.InventoryService/AddStock
```
