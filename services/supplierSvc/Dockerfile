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
WORKDIR /app/services/supplierSvc

# Set GOPROXY with multiple mirrors and retry logic for go mod download
ENV GOPROXY=https://proxy.golang.org,direct
RUN \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for i in 1 2 3; do \
        go mod download && break || \
        (echo "Attempt $i failed, retrying..." && sleep 5); \
    done

# Create necessary directories for generated code
RUN mkdir -p /app/pkg/gen/go/supplier/v1

# Generate protobuf code
RUN protoc -I=/app/proto \
    --go_out=/app/pkg/gen/go/supplier/v1 \
    --go_opt=paths=source_relative \
    --go-grpc_out=/app/pkg/gen/go/supplier/v1 \
    --go-grpc_opt=paths=source_relative \
    /app/proto/supplier/v1/supplier.proto

# Verify the generated files
RUN ls -la /app/pkg/gen/go/supplier/v1/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /supplierSvc ./cmd/main.go

# Final stage
FROM alpine:latest

# Install necessary packages for MongoDB tools
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from builder to the final image
COPY --from=builder /supplierSvc .

# Command to run the executable
CMD ["/root/supplierSvc"]
