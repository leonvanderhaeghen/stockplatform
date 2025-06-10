# Build stage for the application
FROM golang:1.23.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev protoc protobuf-dev

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Set the working directory to the service
WORKDIR /app/services/gatewaySvc

# Set GOPROXY with multiple mirrors and retry logic for go mod download
ENV GOPROXY=https://proxy.golang.org,https://goproxy.cn,direct \
    GOPRIVATE= \
    GOSUMDB=off \
    GOPROXY_FALLBACK=1 \
    GONOPROXY= \
    GONOSUMDB=* \
    GOTOOLCHAIN=local \
    GOCACHE=/tmp/go-build \
    GOMODCACHE=/go/pkg/mod \
    GOMAXPROCS=4
RUN \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    echo "Setting up Go module cache..." && \
    for i in 1 2 3 4 5; do \
        echo "Downloading dependencies (Attempt $i/5)..." && \
        go mod download -x && break || \
        (echo "Attempt $i failed, retrying in 5 seconds..." && sleep 5); \
        [ $i -eq 5 ] && echo "All download attempts failed" && exit 1; \
    done && \
    echo "Dependencies downloaded successfully"

# Generate protobuf code for all services
RUN echo "Generating protobuf code..." && \
    # Create directories
    mkdir -p ../../pkg/gen/go/product/v1 && \
    mkdir -p ../../pkg/gen/go/user/v1 && \
    mkdir -p ../../pkg/gen/go/order/v1 && \
    mkdir -p ../../pkg/gen/go/inventory/v1 && \
    # Generate code for each service
    protoc -I=../../proto \
      --go_out=../../pkg/gen \
      --go_opt=paths=source_relative \
      --go-grpc_out=../../pkg/gen \
      --go-grpc_opt=paths=source_relative \
      product/v1/product.proto \
      user/v1/user.proto \
      order/v1/order.proto \
      inventory/v1/inventory.proto && \
    echo "Protobuf code generated successfully"

# Build the application with retries
RUN \
    echo "Building application..." && \
    for i in 1 2 3; do \
        echo "Build attempt $i/3..." && \
        CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o /gatewaySvc ./cmd/main.go && break || \
        (echo "Build attempt $i failed, retrying in 5 seconds..." && sleep 5); \
        [ $i -eq 3 ] && echo "All build attempts failed" && exit 1; \
    done && \
    echo "Application built successfully"

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set environment variables
ENV TZ=UTC \
    APP_ENV=production

# Copy the binary from builder to the final image
COPY --from=builder /gatewaySvc .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./gatewaySvc"]
