# Gateway Service (gatewaySvc)

This is the main API gateway for the Stock Platform. It handles all incoming HTTP requests and routes them to the appropriate microservices.

## API Documentation

This service provides a comprehensive RESTful API with the following features:

- User authentication and authorization
- Product catalog management
- Order processing
- Inventory management
- Supplier management

### Viewing Documentation

API documentation is available in two formats:

1. **Interactive Swagger UI**:

   - Start the gateway service
   - Navigate to `http://localhost:8080/swagger/index.html`
   - Explore and test the API endpoints interactively

2. **OpenAPI Specification**:

   - The OpenAPI 3.0 specification is available at `http://localhost:8080/swagger/doc.json`
   - You can import this into tools like Postman or generate client libraries

### Generating Documentation

To update the API documentation:

1. Install the Swag tool:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Run the generation script:

   ```bash
   chmod +x generate_swagger.sh
   ./generate_swagger.sh
   ```

### Authentication

Most endpoints require authentication. Include a JWT token in the `Authorization` header:

```http
Authorization: Bearer <your-jwt-token>
```

### Error Handling

All error responses follow the same format:

```json
{
  "error": "Error message",
  "code": 400
}
```

#### Common HTTP Status Codes

- `200`: Success
- `201`: Resource created
- `400`: Bad request (validation error)
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not found
- `500`: Internal server error

## Service Capabilities

- RESTful API endpoints
- Request routing to appropriate microservices
- Authentication middleware
- Request/response transformation
- OpenAPI documentation via Swagger

## API Specification

The Gateway Service provides comprehensive OpenAPI 3.0 documentation available at `/docs` endpoint.

### Authentication

- JWT Bearer Token (for user authentication)
- API Key (for internal service-to-service authentication)
  - `POST /api/v1/auth/register` - Register a new user
  - `POST /api/v1/auth/login` - Authenticate and get JWT token

### API Endpoints

#### Stock Management

- `GET /stocks` - List all stocks
  - Filters: `symbol`, `exchange`, `sector`
  - Pagination: `limit`, `offset`
- `GET /stocks/{id}` - Get stock details
- `POST /stocks` - Create new stock

#### Order Management

- `GET /orders` - List orders
  - Filters: `status`, `date_from`, `date_to`
- `POST /orders` - Place new order
- `GET /orders/me` - Get current user's orders
- `GET /orders/me/{id}` - Get details of a specific order for current user
- `PUT /orders/{id}/status` - Update order status (admin/staff only)

#### Products

- `GET /products` - List products with filtering and pagination
- `GET /products/{id}` - Get product details
- `POST /products` - Create a new product (admin/staff only)
- `PUT /products/{id}` - Update a product (admin/staff only)
- `DELETE /products/{id}` - Delete a product (admin/staff only)

#### Inventory

- `GET /inventory` - List inventory items (admin/staff only)
- `GET /inventory/{id}` - Get inventory item details (admin/staff only)
- `POST /inventory` - Create a new inventory item (admin/staff only)
- `PUT /inventory/{id}` - Update an inventory item (admin/staff only)
- `POST /inventory/{id}/stock/add` - Add stock to an item (admin/staff only)
- `POST /inventory/{id}/stock/remove` - Remove stock from an item (admin/staff only)

#### Users

- `GET /users/me` - Get current user profile
- `PUT /users/me` - Update current user profile
- `GET /users/me/addresses` - Get user addresses
- `POST /users/me/addresses` - Add a new address
- `GET /users` - List all users (admin only)

#### Suppliers (Admin/Staff only)

- `GET /suppliers` - List all suppliers with pagination and search
- `GET /suppliers/{id}` - Get supplier details
- `POST /suppliers` - Create a new supplier
- `PUT /suppliers/{id}` - Update a supplier
- `DELETE /suppliers/{id}` - Delete a supplier

### Response Formats

- JSON (default)
- XML (via Accept header)

To generate client SDKs from the OpenAPI spec:

```bash
openapi-generator generate -i docs/openapi.yaml -g python -o ./client-sdk
```

## Architecture

The Gateway Service follows a clean architecture pattern:

```text
gatewaySvc/
├── cmd/                 # Command-line entry points
├── internal/
│   ├── rest/            # REST handlers and middleware
│   └── services/        # gRPC client implementations
├── Dockerfile           # Container definition
└── swagger/             # OpenAPI documentation
```

## Configuration

The service can be configured using environment variables:

- `REST_PORT` - Port for REST server (default: 8080)
- `JWT_SECRET` - Secret for JWT token validation
- `PRODUCT_SERVICE_ADDR` - Product service address (default: localhost:50053)
- `INVENTORY_SERVICE_ADDR` - Inventory service address (default: localhost:50054)
- `ORDER_SERVICE_ADDR` - Order service address (default: localhost:50055)
- `USER_SERVICE_ADDR` - User service address (default: localhost:50056)
- `SUPPLIER_SERVICE_ADDR` - Supplier service address (default: localhost:50057)

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
   cd ../supplierSvc && go run cmd/main.go &
   ```

2. Run the gateway service:

   ```bash
   go run cmd/main.go
   ```

3. Access the API at `http://localhost:8080`
4. Access the Swagger documentation at `http://localhost:8080/swagger/`

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
  -e SUPPLIER_SERVICE_ADDR=suppliersvc:50057 \
  stockplatform/gatewaysvc
```

## Authentication and Authorization

The gateway handles authentication using JWT tokens:

1. Users register or login via the auth endpoints
2. The gateway validates credentials with the User Service
3. On successful authentication, a JWT token is issued
4. For protected endpoints, the `Authorization: Bearer <token>` header must be included
5. The gateway validates the token and extracts user information
6. Role-based access control is enforced for admin/staff-only endpoints
