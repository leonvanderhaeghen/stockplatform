# Developer Guide - Stock Platform

This guide provides comprehensive information for developers working on the Stock Platform microservices architecture.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Protobuf Architecture](#protobuf-architecture)
- [Development Setup](#development-setup)
- [Working with Services](#working-with-services)
- [Inter-Service Communication](#inter-service-communication)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Architecture Overview

The Stock Platform follows a microservices architecture with the following key principles:

- **Service Autonomy**: Each service owns its API contract and data
- **Loose Coupling**: Services communicate through well-defined interfaces
- **Clean Architecture**: Clear separation of concerns within each service
- **Domain-Driven Design**: Services are organized around business domains

### Services

| Service | Port | Purpose | Database |
|---------|------|---------|----------|
| gatewaySvc | 8080 | HTTP API Gateway | None |
| productSvc | 50053 | Product catalog management | MongoDB |
| inventorySvc | 50054 | Stock level tracking | MongoDB |
| orderSvc | 50055 | Order processing | MongoDB |
| userSvc | 50056 | User management & auth | MongoDB |
| supplierSvc | 50057 | Supplier management | MongoDB |

## Protobuf Architecture

### Service-Owned Proto Files

Each service owns its protobuf definitions and generated code:

```
services/
├── productSvc/
│   ├── api/
│   │   ├── proto/product/v1/product.proto    # Service API definition
│   │   └── gen/go/proto/product/v1/          # Generated Go code
│   ├── buf.gen.yaml                          # Code generation config
│   └── buf.work.yaml                         # Workspace config
└── ...
```

### Client Abstractions

All inter-service communication uses client abstractions in `/pkg/clients/`:

```
pkg/clients/
├── product/client.go      # Product service client
├── inventory/client.go    # Inventory service client
├── order/client.go        # Order service client
├── user/client.go         # User service client
└── supplier/client.go     # Supplier service client
```

### Key Rules

1. **Never import generated protobuf code directly from other services**
2. **Always use client abstractions for inter-service communication**
3. **Each service generates its own protobuf code**
4. **Proto files are versioned (v1, v2, etc.)**

## Development Setup

### Prerequisites

- Go 1.22+
- Docker and Docker Compose
- buf CLI for protobuf generation
- MongoDB (for local development)

### Quick Start

1. **Clone and setup environment:**
   ```bash
   git clone https://github.com/leonvanderhaeghen/stockplatform.git
   cd stockplatform
   cp .env.example .env
   ```

2. **Start with Docker Compose (Recommended):**
   ```bash
   docker-compose up -d
   ```

3. **Or setup for local development:**
   ```bash
   # Install dependencies
   go mod download
   
   # Generate protobuf code for all services
   make generate-proto
   
   # Start MongoDB
   docker run -d --name mongodb -p 27017:27017 mongo:7.0
   
   # Start services individually (in separate terminals)
   cd services/productSvc && go run cmd/main.go
   cd services/inventorySvc && go run cmd/main.go
   # ... etc
   ```

### Verification

Check that all services are running:

```bash
# Health checks
curl http://localhost:8080/health
curl http://localhost:50053/health
curl http://localhost:50054/health
curl http://localhost:50055/health
curl http://localhost:50056/health
curl http://localhost:50057/health
```

## Working with Services

### Service Structure

Each service follows clean architecture principles:

```
services/[serviceName]/
├── api/                    # API definitions
│   ├── proto/             # Protobuf definitions
│   └── gen/go/            # Generated code (gitignored)
├── cmd/
│   └── main.go            # Entry point (minimal)
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database setup & repositories
│   ├── server/            # Server setup & dependency injection
│   ├── domain/            # Business entities
│   ├── application/       # Business logic
│   └── interfaces/        # gRPC handlers
├── buf.gen.yaml           # Protobuf generation config
└── Dockerfile
```

### Adding a New Service

1. **Create service directory structure:**
   ```bash
   mkdir -p services/newSvc/{api/proto/new/v1,cmd,internal/{config,database,server,domain,application,interfaces}}
   ```

2. **Create protobuf definition:**
   ```protobuf
   // services/newSvc/api/proto/new/v1/new.proto
   syntax = "proto3";
   package new.v1;
   option go_package = "github.com/leonvanderhaeghen/stockplatform/services/newSvc/api/gen/go/proto/new/v1;newv1";
   
   service NewService {
     rpc GetItem(GetItemRequest) returns (GetItemResponse);
   }
   
   message GetItemRequest {
     string id = 1;
   }
   
   message GetItemResponse {
     Item item = 1;
   }
   
   message Item {
     string id = 1;
     string name = 2;
   }
   ```

3. **Create buf configuration:**
   ```yaml
   # services/newSvc/buf.gen.yaml
   version: v1
   plugins:
     - plugin: buf.build/protocolbuffers/go
       out: api/gen/go
       opt: paths=source_relative
     - plugin: buf.build/grpc/go
       out: api/gen/go
       opt: paths=source_relative
   ```

4. **Generate protobuf code:**
   ```bash
   cd services/newSvc
   buf generate
   ```

5. **Create client abstraction:**
   ```go
   // pkg/clients/new/client.go
   package new
   
   import (
       "context"
       "fmt"
       
       newv1 "github.com/leonvanderhaeghen/stockplatform/services/newSvc/api/gen/go/proto/new/v1"
       "go.uber.org/zap"
       "google.golang.org/grpc"
   )
   
   type Client struct {
       client newv1.NewServiceClient
       conn   *grpc.ClientConn
       logger *zap.Logger
   }
   
   func NewClient(address string, logger *zap.Logger) (*Client, error) {
       conn, err := grpc.Dial(address, grpc.WithInsecure())
       if err != nil {
           return nil, fmt.Errorf("failed to connect to new service: %w", err)
       }
       
       return &Client{
           client: newv1.NewNewServiceClient(conn),
           conn:   conn,
           logger: logger,
       }, nil
   }
   
   func (c *Client) Close() error {
       return c.conn.Close()
   }
   
   func (c *Client) GetItem(ctx context.Context, req *newv1.GetItemRequest) (*newv1.GetItemResponse, error) {
       c.logger.Debug("Getting item", zap.String("id", req.Id))
       
       resp, err := c.client.GetItem(ctx, req)
       if err != nil {
           c.logger.Error("Failed to get item", zap.Error(err))
           return nil, fmt.Errorf("failed to get item: %w", err)
       }
       
       return resp, nil
   }
   ```

6. **Implement service logic following the existing patterns**

## Inter-Service Communication

### Using Client Abstractions

Always use client abstractions for inter-service communication:

```go
// ✅ Correct: Use client abstraction
import "github.com/leonvanderhaeghen/stockplatform/pkg/clients/product"

func (s *OrderService) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
    // Use product client to validate products
    productClient, err := product.NewClient("product-service:50053", s.logger)
    if err != nil {
        return nil, err
    }
    defer productClient.Close()
    
    productResp, err := productClient.GetProduct(ctx, &productv1.GetProductRequest{
        Id: req.ProductId,
    })
    if err != nil {
        return nil, err
    }
    
    // Continue with order creation...
}
```

```go
// ❌ Incorrect: Direct protobuf import from another service
import productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
```

### Error Handling

Client abstractions provide consistent error handling:

```go
resp, err := productClient.GetProduct(ctx, req)
if err != nil {
    // Error is already wrapped by the client
    s.logger.Error("Failed to get product", zap.Error(err))
    return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
}
```

## Best Practices

### Protobuf Development

1. **Use semantic versioning for packages:**
   ```protobuf
   package product.v1;  // Use v1, v2, etc.
   ```

2. **Add field numbers sequentially:**
   ```protobuf
   message Product {
     string id = 1;
     string name = 2;
     string description = 3;
     // Never reuse field numbers
   }
   ```

3. **Use optional for nullable fields:**
   ```protobuf
   message Product {
     string id = 1;
     optional string description = 2;  // Can be null
   }
   ```

4. **Document your services:**
   ```protobuf
   // ProductService provides operations for managing products
   service ProductService {
     // GetProduct retrieves a product by ID
     rpc GetProduct(GetProductRequest) returns (GetProductResponse);
   }
   ```

### Service Development

1. **Keep main.go minimal:**
   ```go
   func main() {
       cfg, err := config.Load()
       if err != nil {
           log.Fatal("Failed to load config:", err)
       }
       
       if err := server.Run(cfg); err != nil {
           log.Fatal("Server failed:", err)
       }
   }
   ```

2. **Use dependency injection:**
   ```go
   // In server package
   func Run(cfg *config.Config) error {
       db, err := database.Connect(cfg.DatabaseURL)
       if err != nil {
           return err
       }
       
       productRepo := repository.NewProductRepository(db)
       productSvc := service.NewProductService(productRepo)
       productHandler := handler.NewProductHandler(productSvc)
       
       // Setup gRPC server with handlers
   }
   ```

3. **Implement proper logging:**
   ```go
   logger.Info("Processing request",
       zap.String("method", "GetProduct"),
       zap.String("product_id", req.Id),
   )
   ```

### Testing

1. **Unit tests for business logic:**
   ```go
   func TestProductService_GetProduct(t *testing.T) {
       // Test business logic in isolation
   }
   ```

2. **Integration tests for handlers:**
   ```go
   func TestProductHandler_GetProduct(t *testing.T) {
       // Test gRPC handler with real database
   }
   ```

3. **Client abstraction tests:**
   ```go
   func TestProductClient_GetProduct(t *testing.T) {
       // Test client abstraction behavior
   }
   ```

## Troubleshooting

### Common Issues

1. **Protobuf generation fails:**
   ```bash
   # Ensure buf is installed
   go install github.com/bufbuild/buf/cmd/buf@latest
   
   # Check buf.gen.yaml syntax
   cd services/productSvc
   buf lint
   ```

2. **Import path errors:**
   ```bash
   # Regenerate protobuf code
   cd services/productSvc
   buf generate
   
   # Update go.mod if needed
   go mod tidy
   ```

3. **gRPC connection errors:**
   ```bash
   # Check service is running
   curl http://localhost:50053/health
   
   # Check Docker network connectivity
   docker-compose logs product-service
   ```

4. **Database connection issues:**
   ```bash
   # Check MongoDB is running
   docker-compose logs mongodb
   
   # Verify connection string in .env
   echo $MONGO_URI
   ```

### Debugging

1. **Enable debug logging:**
   ```bash
   export LOG_LEVEL=debug
   ```

2. **Use gRPC reflection for testing:**
   ```bash
   grpcurl -plaintext localhost:50053 list
   grpcurl -plaintext localhost:50053 product.v1.ProductService/GetProduct
   ```

3. **Monitor service health:**
   ```bash
   # Check all service health endpoints
   for port in 50053 50054 50055 50056 50057; do
     echo "Checking port $port:"
     curl -s http://localhost:$port/health || echo "Failed"
   done
   ```

### Performance

1. **Connection pooling:**
   ```go
   // Reuse gRPC connections
   var productClient *product.Client
   
   func init() {
       var err error
       productClient, err = product.NewClient("product-service:50053", logger)
       if err != nil {
           log.Fatal(err)
       }
   }
   ```

2. **Context timeouts:**
   ```go
   ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
   defer cancel()
   
   resp, err := productClient.GetProduct(ctx, req)
   ```

3. **Graceful shutdown:**
   ```go
   // Implemented in all services
   c := make(chan os.Signal, 1)
   signal.Notify(c, os.Interrupt, syscall.SIGTERM)
   <-c
   
   // Cleanup connections
   productClient.Close()
   ```

## Additional Resources

- [Protocol Buffers Language Guide](https://developers.google.com/protocol-buffers/docs/proto3)
- [gRPC Go Documentation](https://grpc.io/docs/languages/go/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
