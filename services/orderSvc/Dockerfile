# Build stage for the application
FROM golang:1.23.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev protoc protobuf-dev

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Set the working directory to the service
WORKDIR /app/services/orderSvc

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

# Download dependencies with retries and verbose output
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

# Create necessary directories for generated code
RUN mkdir -p api/gen/go/order/v1

# Generate protobuf code
RUN protoc -I=api/proto \
    --go_out=api/gen/go/order/v1 \
    --go_opt=paths=source_relative \
    --go-grpc_out=api/gen/go/order/v1 \
    --go-grpc_opt=paths=source_relative \
    api/proto/order/v1/order.proto

# Verify the generated files
RUN ls -la api/gen/go/order/v1/

# Build the application with retries and verbose output
RUN \
    echo "Building application..." && \
    for i in 1 2 3; do \
        echo "Build attempt $i/3..." && \
        CGO_ENABLED=0 GOOS=linux go build -v -o /orderSvc ./cmd/main.go && break || \
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
COPY --from=builder /orderSvc .

# Copy the generated protobuf code
COPY --from=builder /app/services/orderSvc/api/gen /app/api/gen/

# Expose the application port
EXPOSE 50055

# Command to run the application
CMD ["./orderSvc"]
