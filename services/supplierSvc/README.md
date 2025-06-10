# Supplier Service

This is the Supplier Service for the Stock Platform. It handles all supplier-related operations.

## Features

- Create, read, update, and delete suppliers
- List suppliers with pagination and search
- Validate supplier data
- gRPC API for inter-service communication

## Prerequisites

- Go 1.18+
- MongoDB 4.4+
- Protocol Buffers compiler (protoc)

## Configuration

Copy `config/config.example.yaml` to `config/config.yaml` and update the configuration as needed.

## Running the Service

1. Start MongoDB
2. Run the service:
   ```bash
   go run cmd/main.go
   ```

## gRPC API

The service exposes the following gRPC endpoints:

- `CreateSupplier` - Create a new supplier
- `GetSupplier` - Get a supplier by ID
- `UpdateSupplier` - Update an existing supplier
- `DeleteSupplier` - Delete a supplier by ID
- `ListSuppliers` - List suppliers with pagination and search

## Development

### Generating Protobuf Code

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/gen/go/supplier/v1/supplier.proto
```

### Testing

Run tests:

```bash
go test ./...
```

## License

MIT
