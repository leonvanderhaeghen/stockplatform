# ---------- Build Stage ----------
FROM golang:1.23.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Install buf for protobuf generation
RUN wget -O /usr/local/bin/buf https://github.com/bufbuild/buf/releases/download/v1.28.1/buf-Linux-x86_64 \
    && chmod +x /usr/local/bin/buf

# Enable Go module proxy & caching
ENV GOPROXY=https://proxy.golang.org,direct

# Set working directory at repository root
WORKDIR /app

# Copy entire project so that inter-service imports resolve
COPY . .

# Switch to store service directory
WORKDIR /app/services/storeSvc

# Cache go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for i in 1 2 3; do go mod download && break || (echo "go mod retry $i" && sleep 5); done

# Re-generate protobuf code to ensure compatibility with current gRPC version
RUN rm -rf api/gen/go && \
    mv ../../buf.work.yaml ../../buf.work.yaml.bak 2>/dev/null || true && \
    buf generate && \
    mv ../../buf.work.yaml.bak ../../buf.work.yaml 2>/dev/null || true && \
    go mod tidy

# Build the service binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /storeSvc ./cmd/main.go

# ---------- Runtime Stage ----------
FROM alpine:latest

# Install runtime deps
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

ENV TZ=UTC \
    APP_ENV=production

# Copy binary & generated code
COPY --from=builder /storeSvc .
COPY --from=builder /app/services/storeSvc/api/gen /app/api/gen/

# Expose gRPC port
EXPOSE 50058

CMD ["./storeSvc"]
