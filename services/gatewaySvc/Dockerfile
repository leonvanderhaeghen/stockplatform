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
ENV GOPROXY=https://proxy.golang.org,direct
RUN \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for i in 1 2 3; do \
        go mod download && break || \
        (echo "Attempt $i failed, retrying..." && sleep 5); \
    done

# Generate protobuf code for all services
# Product Service
RUN mkdir -p ../../pkg/gen/go/product/v1 && \
    protoc -I=../../proto \
      --go_out=../../pkg/gen/go/product/v1 \
      --go_opt=paths=source_relative \
      --go-grpc_out=../../pkg/gen/go/product/v1 \
      --go-grpc_opt=paths=source_relative \
      ../../proto/product/v1/product.proto

# User Service
RUN mkdir -p ../../pkg/gen/go/user/v1 && \
    protoc -I=../../proto \
      --go_out=../../pkg/gen/go/user/v1 \
      --go_opt=paths=source_relative \
      --go-grpc_out=../../pkg/gen/go/user/v1 \
      --go-grpc_opt=paths=source_relative \
      ../../proto/user/v1/user.proto

# Order Service
RUN mkdir -p ../../pkg/gen/go/order/v1 && \
    protoc -I=../../proto \
      --go_out=../../pkg/gen/go/order/v1 \
      --go_opt=paths=source_relative \
      --go-grpc_out=../../pkg/gen/go/order/v1 \
      --go-grpc_opt=paths=source_relative \
      ../../proto/order/v1/order.proto

# Inventory Service
RUN mkdir -p ../../pkg/gen/go/inventory/v1 && \
    protoc -I=../../proto \
      --go_out=../../pkg/gen/go/inventory/v1 \
      --go_opt=paths=source_relative \
      --go-grpc_out=../../pkg/gen/go/inventory/v1 \
      --go-grpc_opt=paths=source_relative \
      ../../proto/inventory/v1/inventory.proto

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /gatewaySvc ./cmd/main.go

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
