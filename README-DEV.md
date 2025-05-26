# Stock Platform - Development Guide

This document contains detailed instructions for setting up the development environment and working with the codebase.

## Prerequisites

- Go 1.22 or later
- Node.js 18+ (for frontend development)
- Docker and Docker Compose
- Buf CLI (`brew install bufbuild/buf/buf` or see [installation](https://docs.buf.build/installation))
- protoc (Protocol Buffers compiler)
- grpcurl (for testing gRPC services)

## Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/leonvanderhaeghen/stockplatform.git
   cd stockplatform
   ```

2. **Install dependencies**
   ```bash
   # Install Go tools
   make deps
   
   # Install development tools
   make tools
   ```

3. **Start the development environment**
   ```bash
   # Start all services (MongoDB, Redis, Jaeger, etc.)
   docker-compose up -d
   
   # Or use the convenience command
   make dev
   ```

4. **Generate code from protobuf definitions**
   ```bash
   make generate
   ```

5. **Build and run services**
   ```bash
   # Build all services
   make build
   
   # Run a specific service (example: product service)
   cd services/productSvc
   go run main.go
   ```

## Development Workflow

### Working with Protocol Buffers

1. Add or modify `.proto` files in the `proto/` directory
2. Generate Go code:
   ```bash
   make generate
   ```

### Running Tests

```bash
# Run all tests
make test

# Run tests for a specific package
cd services/productSvc
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting and Code Quality

```bash
# Run linters
make lint

# Format code
gofmt -w .
```

### API Documentation

After starting the gateway service, API documentation will be available at:
- Swagger UI: `http://localhost:8080/swagger/`
- gRPC reflection: `grpcurl -plaintext localhost:50051 list`

## Service Architecture

### Product Service
- Manages product catalog
- Handles CRUD operations for products and categories
- Exposes gRPC and REST APIs

### Inventory Service
- Tracks stock levels
- Handles inventory updates and reservations
- Uses MongoDB for persistence

### Order Service
- Processes orders
- Manages order lifecycle
- Integrates with payment providers

### User Service
- Handles authentication and authorization
- Manages user accounts and permissions
- Uses JWT for authentication

### API Gateway
- Single entry point for all external requests
- Routes requests to appropriate services
- Handles authentication and rate limiting

## Monitoring and Observability

- **Jaeger**: Distributed tracing at `http://localhost:16686`
- **Prometheus**: Metrics at `http://localhost:9090`
- **Grafana**: Dashboards at `http://localhost:3000` (admin/admin)
- **Mongo Express**: Web UI for MongoDB at `http://localhost:8081`

## Deployment

### Building for Production

```bash
# Build all services for Linux
GOOS=linux GOARCH=amd64 make build

# Build Docker images
docker-compose -f docker-compose.prod.yml build
```

### Kubernetes

See the `deploy/kubernetes` directory for Kubernetes manifests and Helm charts.

## Troubleshooting

### Common Issues

1. **Protobuf Generation Fails**
   - Ensure all dependencies are installed (`make deps`)
   - Check that `protoc` is in your PATH
   - Verify `.proto` file syntax

2. **Connection Issues**
   - Check that all services are running (`docker ps`)
   - Verify network configuration in `docker-compose.yml`
   - Check service logs (`docker-compose logs <service>`)

3. **Database Issues**
   - Ensure MongoDB is running and accessible
   - Check connection strings in service configurations
   - Verify database indexes are created

## Contributing

1. Create a feature branch (`git checkout -b feature/amazing-feature`)
2. Commit your changes (`git commit -m 'Add some amazing feature'`)
3. Push to the branch (`git push origin feature/amazing-feature`)
4. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
