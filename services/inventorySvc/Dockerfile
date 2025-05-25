FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inventorySvc ./cmd/main.go

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set environment variables
ENV TZ=UTC \
    APP_ENV=production

# Copy the binary from the builder stage
COPY --from=builder /app/inventorySvc .

# Expose the application port
EXPOSE 50054

# Command to run the application
CMD ["./inventorySvc"]
