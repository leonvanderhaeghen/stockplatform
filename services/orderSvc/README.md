# Order Service (orderSvc)

This service handles order processing, from creation to fulfillment, for the Stock Platform.

## Features

- Order creation and management
- Order status tracking
- Payment processing integration
- Shipping and fulfillment tracking
- Order history and reporting

## Architecture

The Order Service follows a clean architecture pattern:

```
orderSvc/
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

The service exposes a gRPC API defined in `proto/order/v1/order.proto`:

### Key Endpoints

- `CreateOrder` - Create a new order
- `GetOrder` - Get order details by ID
- `GetUserOrder` - Get a specific order for a user
- `GetUserOrders` - Get all orders for a user
- `UpdateOrderStatus` - Update the status of an order
- `ListOrders` - List orders with filtering options
- `AddPayment` - Add payment information to an order
- `AddTracking` - Add tracking information to an order
- `CancelOrder` - Cancel an order

## Domain Model

The core domain entities include:

- **Order**: Represents a customer order with items, shipping info, and payment details
- **OrderItem**: Represents an individual item in an order
- **OrderStatus**: Enum representing the possible states of an order (PENDING, PAID, PROCESSING, SHIPPED, DELIVERED, CANCELLED)
- **Payment**: Information about payments associated with an order
- **Tracking**: Shipping and tracking information for an order

## Configuration

The service can be configured using environment variables:

- `GRPC_PORT` - Port for gRPC server (default: 50055)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `PRODUCT_SERVICE_ADDR` - Product service address (default: localhost:50053)
- `INVENTORY_SERVICE_ADDR` - Inventory service address (default: localhost:50054)
- `USER_SERVICE_ADDR` - User service address (default: localhost:50056)

## Development

### Prerequisites

- Go 1.22+
- MongoDB 7.0+
- Protocol buffer compiler (protoc)
- Buf CLI

### Running Locally

1. Start MongoDB and dependent services:

```bash
docker-compose up -d mongodb
cd ../productSvc && go run cmd/main.go &
cd ../inventorySvc && go run cmd/main.go &
cd ../userSvc && go run cmd/main.go &
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

Run integration tests (requires MongoDB and dependent services):

```bash
go test -tags=integration ./...
```

## Docker

Build the Docker image:

```bash
docker build -t stockplatform/ordersvc .
```

Run the container:

```bash
docker run -p 50055:50055 \
  -e MONGO_URI=mongodb://mongodb:27017 \
  -e PRODUCT_SERVICE_ADDR=productsvc:50053 \
  -e INVENTORY_SERVICE_ADDR=inventorysvc:50054 \
  -e USER_SERVICE_ADDR=usersvc:50056 \
  stockplatform/ordersvc
```

## Order Flow

1. **Order Creation**:
   - Validate product availability and user information
   - Reserve inventory
   - Create the order with status PENDING

2. **Payment Processing**:
   - Update payment information
   - Update order status to PAID

3. **Order Fulfillment**:
   - Update order status to PROCESSING
   - Reduce inventory quantity
   - Add tracking information
   - Update order status to SHIPPED

4. **Delivery**:
   - Update order status to DELIVERED

5. **Cancellation**:
   - If the order is cancelled, release reserved inventory
   - Update order status to CANCELLED
