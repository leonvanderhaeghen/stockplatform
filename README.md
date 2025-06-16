# Stock Platform

A scalable stock management and e-commerce platform built with Go microservices and Vue.js. The system provides a comprehensive solution for product management, inventory tracking, order processing, and user management with a modern microservices architecture.

## Table of Contents

- [Architecture](#architecture)
- [Services](#services)
- [gRPC Architecture](#grpc-architecture)
- [API Documentation](#api-documentation)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)

## Client Applications

The project includes several client applications for interacting with the Stock Platform:

- **gRPC Client**: Basic gRPC client for the Stock Platform
- **Enhanced gRPC Client**: More feature-rich client with better logging
- **TCP Client**: Simple TCP client for connectivity testing
- **Test Clients**: Various test clients for development and debugging

For more information, see the [Client Applications Documentation](./cmd/README.md).

## Architecture

The Stock Platform is built using a microservices architecture with the following key components:

### Key Components

- **Microservices**: Each service is implemented in Go with gRPC for efficient internal communication
- **Service-Owned APIs**: Each service owns its protobuf definitions and generated code
- **Client Abstractions**: Shared client libraries for clean inter-service communication
- **API Gateway**: REST/HTTP API with OpenAPI documentation for client applications
- **Frontend**: Vue.js 3 with TypeScript and Vite
- **Backend**: Go with gRPC microservices
- **Database**: MongoDB for product catalog, PostgreSQL for orders and users
- **Message Broker**: NATS for event-driven communication
- **API Gateway**: gRPC-Gateway for HTTP/JSON to gRPC translation
- **Authentication**: JWT with refresh tokens
- **Containerization**: Docker and Docker Compose
- **Orchestration**: Kubernetes (optional)
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus and Grafana for metrics collection and visualization
- **Distributed Tracing**: Jaeger for end-to-end request tracing

### System Architecture Diagram

```mermaid
graph LR
    A[Frontend\nVue.js SPA] <--> B[API Gateway\nREST/HTTP]
    B --> C[gRPC Services]
    C --> D[User Service]
    C --> E[Product Service]
    C --> F[Order Service]
    C --> G[Inventory Service]
    C --> H[Supplier Service]
    
    style A fill:#f9f,stroke:#333,stroke-width:2px
    style B fill:#bbf,stroke:#333,stroke-width:2px
    style C fill:#9cf,stroke:#333,stroke-width:2px
    style D,E,F,G,H fill:#9f9,stroke:#333,stroke-width:2px
```

### Microservices Communication

```mermaid
graph TB
    subgraph "Client Abstractions (/pkg/clients)"
        PC[Product Client]
        IC[Inventory Client]
        OC[Order Client]
        UC[User Client]
        SC[Supplier Client]
    end
    
    subgraph "Services"
        PS[Product Service]
        IS[Inventory Service]
        OS[Order Service]
        US[User Service]
        SS[Supplier Service]
    end
    
    PC --> PS
    IC --> IS
    OC --> OS
    UC --> US
    SC --> SS
    
    OS -.-> PC
    OS -.-> IC
    OS -.-> UC
    
    style PC,IC,OC,UC,SC fill:#bbf,stroke:#333,stroke-width:2px
    style PS,IS,OS,US,SS fill:#9f9,stroke:#333,stroke-width:2px
```

## gRPC Architecture

The platform uses a **service-owned protobuf architecture** that ensures clean service boundaries and loose coupling:

### Architecture Principles

- **Service Autonomy**: Each service owns its API contract (protobuf definitions)
- **Loose Coupling**: Services communicate only through client abstractions
- **Clean Boundaries**: No direct imports of generated protobuf code between services
- **Maintainability**: API changes don't break other services directly

### Directory Structure

```text
services/
├── productSvc/
│   ├── api/
│   │   ├── proto/product/v1/product.proto    # Service-owned proto definition
│   │   └── gen/go/proto/product/v1/          # Service-generated code
│   ├── cmd/
│   └── internal/
├── inventorySvc/
│   ├── api/
│   │   ├── proto/inventory/v1/inventory.proto
│   │   └── gen/go/proto/inventory/v1/
│   └── ...
└── ...

pkg/clients/                                   # Shared client abstractions
├── product/client.go                          # Product service client
├── inventory/client.go                        # Inventory service client
├── order/client.go                           # Order service client
├── user/client.go                            # User service client
└── supplier/client.go                        # Supplier service client
```

### Client Abstractions

Each service has a corresponding client abstraction in `/pkg/clients/` that:

- Encapsulates gRPC connection management
- Provides high-level methods for service operations
- Handles logging and error wrapping
- Maintains stable interfaces for inter-service communication

**Example Usage:**

```go
// Services use client abstractions, never direct protobuf imports
import "github.com/leonvanderhaeghen/stockplatform/pkg/clients/product"

// Create client
productClient, err := product.NewClient("product-service:50053", logger)
if err != nil {
    return err
}
defer productClient.Close()

// Use client methods
resp, err := productClient.GetProduct(ctx, &productv1.GetProductRequest{
    Id: "product-123",
})
```

### Code Generation

Each service generates its own protobuf code using `buf generate`:

```bash
# Generate code for a specific service
cd services/productSvc
buf generate

# Or generate for all services
make generate-proto
```

## Services

### 1. Product Service (productSvc)

Manages the product catalog, including product creation, updates, and retrieval.

**Key Features:**

- CRUD operations for products
- Product categorization and attributes
- Product search and filtering
- Price management

### 2. Inventory Service (inventorySvc)

Tracks stock levels and manages inventory items.

**Key Features:**

- Real-time stock level tracking
- Low stock alerts
- Stock adjustments (add/remove)
- Multiple locations support
- Reorder point management

### 3. Order Service (orderSvc)

Handles order processing, from creation to fulfillment.

**Key Features:**

- Order creation and management
- Order status tracking
- Order history
- Returns and refunds processing
- Invoice generation
- Shipping and fulfillment tracking
- Order history and reporting

### 4. User Service (userSvc)

Manages user authentication, profiles, and access control.

**Key Features:**

- User registration and authentication
- Role-based access control (RBAC)
- Profile management
- Address book
- Password reset and recovery
- User profile management
- Address management

### 5. Gateway Service (gatewaySvc)

Provides a unified REST API for client applications, translating between REST and gRPC.

**Key Features:**

- RESTful API endpoints
- Request validation
- Rate limiting
- CORS support
- Request/Response logging
- Authentication middleware
- Request/response transformation
- OpenAPI documentation

## API Documentation

The API Gateway provides RESTful endpoints for client applications. The full OpenAPI documentation is available at `/swagger/` when running the gateway service.

### Key API Endpoints

#### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Authenticate and get JWT token

#### Products

- `GET /api/v1/products` - List products with filtering and pagination
- `GET /api/v1/products/{id}` - Get product details
- `POST /api/v1/products` - Create a new product (admin/staff only)
- `PUT /api/v1/products/{id}` - Update a product (admin/staff only)
- `DELETE /api/v1/products/{id}` - Delete a product (admin/staff only)

#### Inventory

- `GET /api/v1/inventory` - List inventory items (admin/staff only)
- `GET /api/v1/inventory/{id}` - Get inventory item details (admin/staff only)
- `POST /api/v1/inventory` - Create a new inventory item (admin/staff only)
- `PUT /api/v1/inventory/{id}` - Update an inventory item (admin/staff only)
- `POST /api/v1/inventory/{id}/stock/add` - Add stock to an item (admin/staff only)
- `POST /api/v1/inventory/{id}/stock/remove` - Remove stock from an item (admin/staff only)

#### Orders

- `GET /api/v1/orders/me` - Get current user's orders
- `GET /api/v1/orders/me/{id}` - Get details of a specific order for current user
- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - List all orders (admin/staff only)
- `PUT /api/v1/orders/{id}/status` - Update order status (admin/staff only)

#### Users

- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile
- `GET /api/v1/users/me/addresses` - Get user addresses
- `POST /api/v1/users/me/addresses` - Add a new address
- `GET /api/v1/admin/users` - List all users (admin only)

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)

### Quick Start with Docker Compose

1. **Clone the repository**

   ```bash
   git clone https://github.com/leonvanderhaeghen/stockplatform.git
   cd stockplatform
   ```

2. **Copy environment configuration**

   ```bash
   cp .env.example .env
   ```

3. **Start all services with Docker Compose**

   ```bash
   # Enable BuildKit for faster builds
   export DOCKER_BUILDKIT=1
   
   # Start all services
   docker-compose up -d
   ```

4. **Verify services are running**

   ```bash
   docker-compose ps
   ```

5. **Access the application**

   - **API Gateway**: [http://localhost:8080](http://localhost:8080)
   - **Frontend**: [http://localhost:3001](http://localhost:3001)
   - **MongoDB Express**: [http://localhost:8083](http://localhost:8083) (admin/admin123)
   - **Grafana**: [http://localhost:3000](http://localhost:3000) (admin/admin)
   - **Prometheus**: [http://localhost:9090](http://localhost:9090)
   - **Jaeger**: [http://localhost:16686](http://localhost:16686)

### Health Checks

Check service health:

```bash
# Gateway health check
curl http://localhost:8080/api/v1/health

# Individual service logs
docker-compose logs product-service
docker-compose logs inventory-service
```

### Development Setup (Alternative)

For local development without Docker:

1. **Install dependencies**

   ```bash
   make deps
   ```

2. **Generate protobuf code**

   ```bash
   make generate
   ```

3. **Start MongoDB**

   ```bash
   # Using Docker
   docker run -d --name mongodb -p 27017:27017 mongo:7.0
   
   # Or using local installation
   mongod --dbpath /path/to/your/db
   ```

4. **Start services individually**

   ```bash
   # Terminal 1 - Product Service
   cd services/productSvc && go run cmd/main.go
   
   # Terminal 2 - Inventory Service  
   cd services/inventorySvc && go run cmd/main.go
   
   # Terminal 3 - Order Service
   cd services/orderSvc && go run cmd/main.go
   
   # Terminal 4 - User Service
   cd services/userSvc && go run cmd/main.go
   
   # Terminal 5 - Supplier Service
   cd services/supplierSvc && go run cmd/main.go
   
   # Terminal 6 - Gateway Service
   cd services/gatewaySvc && go run cmd/main.go
   ```

### Environment Variables

The application uses Docker Compose environment variables. Key configurations:

#### Docker Compose Environment

- `MONGO_URI=mongodb://admin:admin123@mongodb:27017` - MongoDB connection (Docker internal)
- `JWT_SECRET=your-secret-key-here-change-in-production` - JWT signing secret
- Service addresses use Docker service names (e.g., `product-service:50053`)

#### Local Development Environment

- `MONGO_URI=mongodb://localhost:27017` - MongoDB connection (local)
- Service addresses use localhost (e.g., `localhost:50053`)

## Development Workflow

### Code Organization

Each service follows a clean architecture pattern with service-owned protobuf definitions:

```text
services/
  └── [serviceName]/
      ├── api/                 # API definitions and generated code
      │   ├── proto/           # Protobuf definitions (.proto files)
      │   │   └── [service]/v1/[service].proto
      │   └── gen/go/          # Generated Go code (gitignored)
      │       └── proto/[service]/v1/
      ├── cmd/                 # Command-line entry points
      ├── internal/
      │   ├── config/          # Configuration management
      │   ├── database/        # Database initialization and repositories
      │   ├── server/          # Server setup and dependency injection
      │   ├── domain/          # Domain models and interfaces
      │   ├── application/     # Application services (business logic)
      │   ├── infrastructure/  # External dependencies (database, etc.)
      │   └── interfaces/      # Interface adapters (gRPC handlers, etc.)
      ├── buf.gen.yaml         # Protocol buffer generation config
      ├── buf.work.yaml        # Protocol buffer workspace config
      └── Dockerfile           # Container definition

pkg/clients/                  # Shared client abstractions
  ├── [service]/client.go     # Client abstraction for each service
  └── ...
```

### Protobuf Development Guidelines

#### 1. Service-Owned Proto Files

- Each service owns its `.proto` files in `services/[serviceName]/api/proto/[service]/v1/`
- Never import proto files from other services
- Use semantic versioning for proto packages (v1, v2, etc.)

#### 2. Code Generation

Generate protobuf code for a specific service:

```bash
cd services/productSvc
buf generate
```

Generate for all services:

```bash
make generate-proto
```

#### 3. Client Abstractions

- All inter-service communication must use client abstractions from `/pkg/clients/`
- Never import generated protobuf code directly from other services
- Client abstractions provide stable interfaces and handle connection management

**Example: Adding a new service method**

1. Update the service's `.proto` file:
   ```protobuf
   service ProductService {
     rpc GetProduct(GetProductRequest) returns (GetProductResponse);
     rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse);
     // Add new method
     rpc UpdateProductPrice(UpdateProductPriceRequest) returns (UpdateProductPriceResponse);
   }
   ```

2. Regenerate protobuf code:
   ```bash
   cd services/productSvc
   buf generate
   ```

3. Implement the method in the service handler

4. Add the method to the client abstraction in `/pkg/clients/product/client.go`:
   ```go
   func (c *Client) UpdateProductPrice(ctx context.Context, req *productv1.UpdateProductPriceRequest) (*productv1.UpdateProductPriceResponse, error) {
       // Implementation
   }
   ```

#### 4. Breaking Changes

- Use semantic versioning for major API changes (v1 → v2)
- Maintain backward compatibility within the same version
- Coordinate breaking changes across dependent services

### Making Changes

1. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the architecture guidelines

3. If modifying protobuf definitions:
   ```bash
   # Regenerate code
   cd services/[serviceName]
   buf generate
   
   # Update client abstractions if needed
   # Test compilation
   go build ./pkg/clients/...
   ```

4. Run tests:

   ```bash
   go test ./...
   ```

5. Commit your changes with a descriptive message following Conventional Commits:

   ```bash
   git commit -m "feat: add new product price update functionality"
   ```

6. Push changes and create a pull request

## Testing

### Unit Tests

Run unit tests for all services:

```bash
go test ./...
```

Or for a specific service:

```bash
cd services/userSvc
go test ./...
```

### Integration Tests

Integration tests require a running MongoDB instance:

```bash
go test ./... -tags=integration
```

### End-to-End Tests

End-to-end tests require all services to be running:

```bash
cd tests
go test -v ./e2e
```

## Deployment

### Docker Deployment

Build and run all services using Docker Compose:

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment

Kubernetes manifests are available in the `k8s/` directory:

```bash
kubectl apply -f k8s/
```

## Contributing

### Development Guidelines

- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages
- Write tests for all new features
- Update documentation when making changes
- Keep the main branch stable
- Use linters to ensure code quality

### Code Review Process

1. Create a pull request with a clear description
2. Ensure CI passes on your branch
3. Request review from at least one team member
4. Address review comments
5. Merge once approved

## License

This project is licensed under the MIT License - see the LICENSE file for details.
