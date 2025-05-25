# Client Applications

This directory contains various client applications for interacting with the Stock Platform.

## Available Clients

### 1. gRPC Client
A basic gRPC client for the Stock Platform.

**Build and Run:**
```bash
make run-grpc-client
```

### 2. Enhanced gRPC Client
A more feature-rich gRPC client with better logging and error handling.

**Build and Run:**
```bash
make run-grpc-client-enhanced
```

### 3. TCP Client
A simple TCP client that can connect to the server (useful for basic connectivity testing).

**Build and Run:**
```bash
make run-tcp-client
```

### 4. Test gRPC Client
A test client with JSON output formatting.

**Build and Run:**
```bash
make run-test-grpc-client
```

### 5. Minimal gRPC Client
A minimal gRPC client for basic testing.

**Build and Run:**
```bash
make run-test-grpc-minimal
```

## Building All Clients

To build all client applications at once:

```bash
make build-clients
```

This will place the compiled binaries in the `bin/` directory.

## Environment Variables

Some clients support configuration via environment variables:

- `GRPC_SERVER_ADDR`: The address of the gRPC server (default: `localhost:50053`)

Example:
```bash
export GRPC_SERVER_ADDR=localhost:50053
make run-grpc-client
```
