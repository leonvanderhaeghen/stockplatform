# User Service (userSvc)

This service manages user authentication, profiles, and access control for the Stock Platform.

## Features

- User registration and authentication
- JWT-based authorization
- Role-based access control (customer, staff, admin)
- User profile management
- Address management

## Architecture

The User Service follows a clean architecture pattern:

```
userSvc/
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

The service exposes a gRPC API defined in `proto/user/v1/user.proto`:

### Key Endpoints

- `Register` - Register a new user
- `Login` - Authenticate a user and return a JWT token
- `GetUser` - Get user details by ID
- `GetUserByEmail` - Get user details by email
- `UpdateProfile` - Update user profile information
- `ChangePassword` - Change user password
- `ActivateUser` - Activate a user account
- `DeactivateUser` - Deactivate a user account
- `ListUsers` - List users with filtering options
- `CreateAddress` - Create a new address for a user
- `GetAddresses` - Get all addresses for a user
- `GetDefaultAddress` - Get a user's default address
- `UpdateAddress` - Update a user address
- `DeleteAddress` - Delete a user address
- `SetDefaultAddress` - Set an address as the default for a user

## Configuration

The service can be configured using environment variables:

- `GRPC_PORT` - Port for gRPC server (default: 50056)
- `MONGO_URI` - MongoDB connection string (default: mongodb://localhost:27017)
- `JWT_SECRET` - Secret for JWT token generation

## Development

### Prerequisites

- Go 1.22+
- MongoDB 7.0+
- Protocol buffer compiler (protoc)
- Buf CLI

### Running Locally

1. Start MongoDB:

```bash
docker-compose up -d mongodb
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

Run integration tests (requires MongoDB):

```bash
go test -tags=integration ./...
```

## Docker

Build the Docker image:

```bash
docker build -t stockplatform/usersvc .
```

Run the container:

```bash
docker run -p 50056:50056 -e MONGO_URI=mongodb://mongodb:27017 stockplatform/usersvc
```
