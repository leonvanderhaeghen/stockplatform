# Stock Platform API Documentation

## Overview

The Stock Platform is a microservice-based e-commerce and inventory management system built with Go and gRPC. This document provides comprehensive API documentation for all services.

## Architecture

The platform consists of the following services:

- **Gateway Service** (Port 8080) - HTTP REST API Gateway
- **Product Service** (Port 50053) - Product catalog management
- **Inventory Service** (Port 50054) - Stock and inventory tracking
- **Order Service** (Port 50055) - Order processing and management
- **User Service** (Port 50056) - User authentication and management
- **Supplier Service** (Port 50057) - Supplier relationship management

## Authentication

Most endpoints require JWT authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Health Checks

All services implement gRPC health checks:

```bash
# Check service health
grpc_health_probe -addr=localhost:50053
```

## Service APIs

### Gateway Service (HTTP REST)

Base URL: `http://localhost:8080`

#### Authentication Endpoints

```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe"
}
```

#### Product Endpoints

```http
GET /products
GET /products/{id}
POST /products
PUT /products/{id}
DELETE /products/{id}
```

#### Order Endpoints

```http
GET /orders
GET /orders/{id}
POST /orders
PUT /orders/{id}
DELETE /orders/{id}
```

### Product Service (gRPC)

Service: `product.v1.ProductService`
Address: `localhost:50053`

#### Methods

##### CreateProduct
```protobuf
rpc CreateProduct(CreateProductRequest) returns (CreateProductResponse)

message CreateProductRequest {
  string name = 1;
  string description = 2;
  double price = 3;
  string category = 4;
  string sku = 5;
  int32 stock_quantity = 6;
}
```

##### GetProduct
```protobuf
rpc GetProduct(GetProductRequest) returns (GetProductResponse)

message GetProductRequest {
  string id = 1;
}
```

##### ListProducts
```protobuf
rpc ListProducts(ListProductsRequest) returns (ListProductsResponse)

message ListProductsRequest {
  int32 page_size = 1;
  int32 page = 2;
  string category = 3;
  string search = 4;
}
```

##### UpdateProduct
```protobuf
rpc UpdateProduct(UpdateProductRequest) returns (UpdateProductResponse)
```

##### DeleteProduct
```protobuf
rpc DeleteProduct(DeleteProductRequest) returns (DeleteProductResponse)
```

### User Service (gRPC)

Service: `user.v1.UserService`
Address: `localhost:50056`

#### Methods

##### CreateUser
```protobuf
rpc CreateUser(CreateUserRequest) returns (CreateUserResponse)

message CreateUserRequest {
  string email = 1;
  string password = 2;
  string first_name = 3;
  string last_name = 4;
}
```

##### GetUser
```protobuf
rpc GetUser(GetUserRequest) returns (GetUserResponse)

message GetUserRequest {
  string id = 1;
}
```

##### AuthenticateUser
```protobuf
rpc AuthenticateUser(AuthenticateUserRequest) returns (AuthenticateUserResponse)

message AuthenticateUserRequest {
  string email = 1;
  string password = 2;
}
```

##### UpdateUser
```protobuf
rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse)
```

##### DeleteUser
```protobuf
rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse)
```

### Order Service (gRPC)

Service: `order.v1.OrderService`
Address: `localhost:50055`

#### Methods

##### CreateOrder
```protobuf
rpc CreateOrder(CreateOrderRequest) returns (CreateOrderResponse)

message CreateOrderRequest {
  string user_id = 1;
  repeated OrderItem items = 2;
  double total_amount = 3;
  string status = 4;
}

message OrderItem {
  string product_id = 1;
  int32 quantity = 2;
  double price = 3;
}
```

##### GetOrder
```protobuf
rpc GetOrder(GetOrderRequest) returns (GetOrderResponse)

message GetOrderRequest {
  string id = 1;
}
```

##### ListOrders
```protobuf
rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse)

message ListOrdersRequest {
  string user_id = 1;
  int32 page_size = 2;
  int32 page = 3;
  string status = 4;
}
```

##### UpdateOrderStatus
```protobuf
rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse)

message UpdateOrderStatusRequest {
  string id = 1;
  string status = 2;
}
```

### Inventory Service (gRPC)

Service: `inventory.v1.InventoryService`
Address: `localhost:50054`

#### Methods

##### GetInventory
```protobuf
rpc GetInventory(GetInventoryRequest) returns (GetInventoryResponse)

message GetInventoryRequest {
  string product_id = 1;
}
```

##### UpdateInventory
```protobuf
rpc UpdateInventory(UpdateInventoryRequest) returns (UpdateInventoryResponse)

message UpdateInventoryRequest {
  string product_id = 1;
  int32 quantity = 2;
  string operation = 3; // "add" or "subtract"
}
```

##### CheckAvailability
```protobuf
rpc CheckAvailability(CheckAvailabilityRequest) returns (CheckAvailabilityResponse)

message CheckAvailabilityRequest {
  string product_id = 1;
  int32 required_quantity = 2;
}
```

### Supplier Service (gRPC)

Service: `supplier.v1.SupplierService`
Address: `localhost:50057`

#### Methods

##### CreateSupplier
```protobuf
rpc CreateSupplier(CreateSupplierRequest) returns (CreateSupplierResponse)

message CreateSupplierRequest {
  string name = 1;
  string contact_email = 2;
  string contact_phone = 3;
  string address = 4;
}
```

##### GetSupplier
```protobuf
rpc GetSupplier(GetSupplierRequest) returns (GetSupplierResponse)

message GetSupplierRequest {
  string id = 1;
}
```

##### ListSuppliers
```protobuf
rpc ListSuppliers(ListSuppliersRequest) returns (ListSuppliersResponse)

message ListSuppliersRequest {
  int32 page_size = 1;
  int32 page = 2;
}
```

## Error Handling

All services use standard gRPC status codes:

- `OK` (0) - Success
- `INVALID_ARGUMENT` (3) - Invalid request parameters
- `NOT_FOUND` (5) - Resource not found
- `ALREADY_EXISTS` (6) - Resource already exists
- `PERMISSION_DENIED` (7) - Authentication/authorization failed
- `INTERNAL` (13) - Internal server error
- `UNAVAILABLE` (14) - Service temporarily unavailable

## Rate Limiting

The Gateway Service implements rate limiting:
- 100 requests per minute per IP address
- 1000 requests per hour per authenticated user

## Monitoring and Observability

### Health Checks
```bash
# Check all services
curl http://localhost:8080/health

# Check individual gRPC services
grpc_health_probe -addr=localhost:50053
grpc_health_probe -addr=localhost:50054
grpc_health_probe -addr=localhost:50055
grpc_health_probe -addr=localhost:50056
grpc_health_probe -addr=localhost:50057
```

### Metrics
Prometheus metrics are available at:
- Gateway Service: `http://localhost:8080/metrics`
- Prometheus UI: `http://localhost:9090`
- Grafana Dashboard: `http://localhost:3000`

### Tracing
Jaeger tracing is available at:
- Jaeger UI: `http://localhost:16686`

## Development

### Running Locally
```bash
# Start all services with Docker Compose
docker-compose up -d

# Run integration tests
./scripts/run-integration-tests.ps1

# Stop services
docker-compose down
```

### Environment Variables

Key environment variables for configuration:

```bash
# Database
MONGO_URI=mongodb://mongo:27017

# Service Discovery
PRODUCT_SERVICE_ADDR=product-service:50053
INVENTORY_SERVICE_ADDR=inventory-service:50054
ORDER_SERVICE_ADDR=order-service:50055
USER_SERVICE_ADDR=user-service:50056
SUPPLIER_SERVICE_ADDR=supplier-service:50057

# Authentication
JWT_SECRET=your-secret-key-here

# Logging
LOG_LEVEL=info
```

## Examples

### Creating a Product via Gateway
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "category": "Electronics",
    "sku": "LAP-001",
    "stock_quantity": 10
  }'
```

### Creating an Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "user_id": "user123",
    "items": [
      {
        "product_id": "prod123",
        "quantity": 2,
        "price": 999.99
      }
    ],
    "total_amount": 1999.98,
    "status": "pending"
  }'
```

## Support

For questions or issues:
- Check the logs: `docker-compose logs <service-name>`
- Run health checks to verify service status
- Review the integration tests for usage examples
