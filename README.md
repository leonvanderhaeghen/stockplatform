# Stock Platform

A scalable stock management and e-commerce platform built with Go microservices and Vue.js. The system provides a comprehensive solution for product management, inventory tracking, order processing, and user management with a modern microservices architecture.

## Table of Contents

- [Architecture](#architecture)
- [Services](#services)
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
    
    style A fill:#f9f,stroke:#333,stroke-width:2px
    style B fill:#bbf,stroke:#333,stroke-width:2px
    style C fill:#9cf,stroke:#333,stroke-width:2px
    style D,E,F,G fill:#9f9,stroke:#333,stroke-width:2px
                                 | gRPC
                                 v
 +------------+  +-------------+  +-----------+  +-----------+
 | Product    |  | Inventory   |  | Order     |  | User      |
 | Service    |  | Service     |  | Service   |  | Service   |
 +------------+  +-------------+  +-----------+  +-----------+
       |               |               |              |
       +---------------+---------------+--------------+
                                |
                                v
                         +---------------+
                         |   MongoDB     |
                         +---------------+
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

- Go 1.22+

### Development Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/stockplatform.git
   cd stockplatform
   ```

2. **Install dependencies**

   ```bash
   make deps
   ```

3. **Generate protobuf code**

   ```bash
   make generate
   ```

4. **Build all services**

   ```bash
   make build
   ```

5. **Start the services** (in separate terminal windows):

   ```bash
   # Start User Service
   cd services/userSvc
   go run cmd/main.go
   ```

   ```bash
   # Start Product Service
   cd services/productSvc
   go run cmd/main.go
   ```

   ```bash
   # Start Inventory Service
   cd services/inventorySvc
   go run cmd/main.go
   ```

   ```bash
   # Start Order Service
   cd services/orderSvc
   go run cmd/main.go
   ```

   ```bash
   # Start Gateway Service
   cd services/gatewaySvc
   go run cmd/main.go
   ```

6. Access the following URLs:

   - API Gateway: [http://localhost:8080](http://localhost:8080)
   - Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Environment Variables

Each service can be configured using environment variables:

#### User Service

- `GRPC_PORT` - Port for gRPC server (default: 50056)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `JWT_SECRET` - Secret for JWT token generation

#### Product Service

- `GRPC_PORT` - Port for gRPC server (default: 50053)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)

#### Inventory Service

- `GRPC_PORT` - Port for gRPC server (default: 50054)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `PRODUCT_SERVICE_ADDR` - Product service address (default: localhost:50053)

#### Order Service

- `GRPC_PORT` - Port for gRPC server (default: 50055)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `PRODUCT_SERVICE_ADDR` - Product service address (default: localhost:50053)
- `INVENTORY_SERVICE_ADDR` - Inventory service address (default: localhost:50054)
- `USER_SERVICE_ADDR` - User service address (default: localhost:50056)

#### Gateway Service

- `REST_PORT` - Port for REST server (default: 8080)
- `JWT_SECRET` - Secret for JWT token validation
- `PRODUCT_SERVICE_ADDR` - Product service address (default: localhost:50053)
- `INVENTORY_SERVICE_ADDR` - Inventory service address (default: localhost:50054)
- `ORDER_SERVICE_ADDR` - Order service address (default: localhost:50055)
- `USER_SERVICE_ADDR` - User service address (default: localhost:50056)

## Development Workflow

### Code Organization

Each service follows a clean architecture pattern:

```text
services/
  └── [serviceName]/
      ├── cmd/                 # Command-line entry points
      ├── internal/
      │   ├── domain/          # Domain models and interfaces
      │   ├── application/     # Application services (business logic)
      │   ├── infrastructure/  # External dependencies (database, etc.)
      │   └── interfaces/      # Interface adapters (gRPC, etc.)
      ├── Dockerfile           # Container definition
      ├── buf.gen.yaml         # Protocol buffer generation config
      └── buf.work.yaml        # Protocol buffer workspace config
```

### Making Changes

1. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes

3. Run tests:

   ```bash
   go test ./...
   ```

4. Commit your changes with a descriptive message following Conventional Commits:

   ```bash
   git commit -m "feat: add new product attribute functionality"
   ````

5. Push changes and create a pull request

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
