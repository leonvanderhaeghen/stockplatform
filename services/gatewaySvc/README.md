# Gateway Service (gatewaySvc)

This service acts as an API Gateway, providing a unified REST API for client applications and translating between REST and gRPC.

## Features

- RESTful API endpoints
- Request routing to appropriate microservices
- Authentication middleware
- Request/response transformation
- OpenAPI documentation via Swagger

## Architecture

The Gateway Service follows a clean architecture pattern:

```
gatewaySvc/
├── cmd/                 # Command-line entry points
├── internal/
│   ├── rest/            # REST handlers and middleware
│   └── services/        # gRPC client implementations
├── Dockerfile           # Container definition
└── swagger/             # OpenAPI documentation
```

## API

The service exposes a REST API with the following key endpoints:

### Authentication

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Authenticate and get JWT token

### Products

- `GET /api/v1/products` - List products with filtering and pagination
- `GET /api/v1/products/{id}` - Get product details
- `POST /api/v1/products` - Create a new product (admin/staff only)
- `PUT /api/v1/products/{id}` - Update a product (admin/staff only)
- `DELETE /api/v1/products/{id}` - Delete a product (admin/staff only)

### Inventory

- `GET /api/v1/inventory` - List inventory items (admin/staff only)
- `GET /api/v1/inventory/{id}` - Get inventory item details (admin/staff only)
- `POST /api/v1/inventory` - Create a new inventory item (admin/staff only)
- `PUT /api/v1/inventory/{id}` - Update an inventory item (admin/staff only)
- `POST /api/v1/inventory/{id}/stock/add` - Add stock to an item (admin/staff only)
- `POST /api/v1/inventory/{id}/stock/remove` - Remove stock from an item (admin/staff only)

### Orders

- `GET /api/v1/orders/me` - Get current user's orders
- `GET /api/v1/orders/me/{id}` - Get details of a specific order for current user
- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - List all orders (admin/staff only)
- `PUT /api/v1/orders/{id}/status` - Update order status (admin/staff only)

### Users

- `GET /api/v1/users/me` - Get current user profile
- `PUT /api/v1/users/me` - Update current user profile
- `GET /api/v1/users/me/addresses` - Get user addresses
- `POST /api/v1/users/me/addresses` - Add a new address
- `GET /api/v1/admin/users` - List all users (admin only)

## Configuration

The service can be configured using environment variables:

- `REST_PORT` - Port for REST server (default: 8080)
- `JWT_SECRET` - Secret for JWT token validation
- `PRODUCT_SERVICE_ADDR` - Product service address (default: localhost:50053)
- `INVENTORY_SERVICE_ADDR` - Inventory service address (default: localhost:50054)
- `ORDER_SERVICE_ADDR` - Order service address (default: localhost:50055)
- `USER_SERVICE_ADDR` - User service address (default: localhost:50056)

## Development

### Prerequisites

- Go 1.22+
- All backend microservices running

### Running Locally

1. Start dependent services:

```bash
cd ../productSvc && go run cmd/main.go &
cd ../inventorySvc && go run cmd/main.go &
cd ../orderSvc && go run cmd/main.go &
cd ../userSvc && go run cmd/main.go &
```

2. Run the gateway service:

```bash
go run cmd/main.go
```

3. Access the API at http://localhost:8080
4. Access the Swagger documentation at http://localhost:8080/swagger/

### Testing

Run unit tests:

```bash
go test ./...
```

Run integration tests (requires all microservices):

```bash
go test -tags=integration ./...
```

## Docker

Build the Docker image:

```bash
docker build -t stockplatform/gatewaysvc .
```

Run the container:

```bash
docker run -p 8080:8080 \
  -e JWT_SECRET=your-secret-key \
  -e PRODUCT_SERVICE_ADDR=productsvc:50053 \
  -e INVENTORY_SERVICE_ADDR=inventorysvc:50054 \
  -e ORDER_SERVICE_ADDR=ordersvc:50055 \
  -e USER_SERVICE_ADDR=usersvc:50056 \
  stockplatform/gatewaysvc
```

## Authentication

The gateway handles authentication using JWT tokens:

1. Users register or login via the auth endpoints
2. The gateway validates credentials with the User Service
3. On successful authentication, a JWT token is issued
4. For protected endpoints, the `Authorization: Bearer <token>` header must be included
5. The gateway validates the token and extracts user information
6. Role-based access control is enforced for admin/staff-only endpoints
