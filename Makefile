# Stock Platform Makefile

.PHONY: help deps generate-proto build test clean docker-build docker-up docker-down lint format

# Default target
help:
	@echo "Stock Platform Development Commands"
	@echo ""
	@echo "Setup:"
	@echo "  deps              Install dependencies"
	@echo "  generate-proto    Generate protobuf code for all services"
	@echo ""
	@echo "Development:"
	@echo "  build             Build all services"
	@echo "  test              Run all tests"
	@echo "  lint              Run linters"
	@echo "  format            Format code"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build      Build Docker images"
	@echo "  docker-up         Start services with Docker Compose"
	@echo "  docker-down       Stop Docker Compose services"
	@echo ""
	@echo "Utilities:"
	@echo "  clean             Clean build artifacts"
	@echo "  health-check      Check all service health endpoints"

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "Installing buf..."
	go install github.com/bufbuild/buf/cmd/buf@latest
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate protobuf code for all services
generate-proto:
	@echo "Generating protobuf code for all services..."
	@for service in productSvc inventorySvc orderSvc userSvc supplierSvc storeSvc; do \
		echo "Generating code for $$service..."; \
		cd services/$$service ; buf generate ; cd ../..; \
	done
	@echo "Protobuf code generation complete!"

# Build all services
build:
	@echo "Building all services..."
	@for service in productSvc inventorySvc orderSvc userSvc supplierSvc storeSvc gatewaySvc; do \
		echo "Building $$service..."; \
		cd services/$$service ; go build -o ../../bin/$$service ./cmd/main.go ; cd ../..; \
	done
	@echo "Building client abstractions..."
	go build ./pkg/clients/...
	@echo "Build complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run integration tests (requires running services)
test-integration:
	@echo "Running integration tests..."
	go test ./... -tags=integration

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run ./...
	@echo "Running buf lint on proto files..."
	@for service in productSvc inventorySvc orderSvc userSvc supplierSvc storeSvc; do \
		echo "Linting $$service proto files..."; \
		cd services/$$service ; buf lint ; cd ../..; \
	done

# Format code
format:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting proto files..."
	@for service in productSvc inventorySvc orderSvc userSvc supplierSvc storeSvc; do \
		echo "Formatting $$service proto files..."; \
		cd services/$$service ; buf format -w ; cd ../..; \
	done

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@for service in productSvc inventorySvc orderSvc userSvc supplierSvc storeSvc; do \
		echo "Cleaning generated code for $$service..."; \
		rm -rf services/$$service/api/gen/; \
	done
	@echo "Clean complete!"

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d
	@echo "Services started! Check health with: make health-check"

docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-logs:
	@echo "Showing Docker Compose logs..."
	docker-compose logs -f

# Health check all services
health-check:
	@echo "Checking service health..."
	@echo "Gateway Service (HTTP):"
	@curl -s http://localhost:8080/health || echo "  ❌ Gateway service not responding"
	@echo ""
	@echo "gRPC Services:"
	@for port in 50053 50054 50055 50056 50057 50058; do \
		service_name=$$(case $$port in \
			50053) echo "Product Service" ;; \
			50054) echo "Inventory Service" ;; \
			50055) echo "Order Service" ;; \
			50056) echo "User Service" ;; \
			50057) echo "Supplier Service" ;; \
			50058) echo "Store Service" ;; \
		esac); \
		echo "  $$service_name ($$port):"; \
		curl -s http://localhost:$$port/health >/dev/null ; echo "    ✅ Healthy" || echo "    ❌ Not responding"; \
	done

# Development shortcuts
dev-product:
	@echo "Starting Product Service for development..."
	cd services/productSvc ; go run cmd/main.go

dev-inventory:
	@echo "Starting Inventory Service for development..."
	cd services/inventorySvc ; go run cmd/main.go

dev-order:
	@echo "Starting Order Service for development..."
	cd services/orderSvc ; go run cmd/main.go

dev-user:
	@echo "Starting User Service for development..."
	cd services/userSvc ; go run cmd/main.go

dev-supplier:
	@echo "Starting Supplier Service for development..."
	cd services/supplierSvc ; go run cmd/main.go

dev-store:
	@echo "Starting Store Service for development..."
	cd services/storeSvc ; go run cmd/main.go

dev-gateway:
	@echo "Starting Gateway Service for development..."
	cd services/gatewaySvc ; go run cmd/main.go

# Quick setup for new developers
setup: deps generate-proto
	@echo "Setting up development environment..."
	cp .env.example .env
	@echo "✅ Development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Edit .env file with your configuration"
	@echo "2. Run 'make docker-up' to start all services"
	@echo "3. Run 'make health-check' to verify services are running"
	@echo "4. See 'make help' for more commands"

# Verify protobuf architecture
verify-proto:
	@echo "Verifying protobuf architecture..."
	@echo "Checking that client abstractions compile..."
	go build ./pkg/clients/...
	@echo "✅ Client abstractions compile successfully"
	@echo ""
	@echo "Checking service-owned proto structure..."
	@for service in productSvc inventorySvc orderSvc userSvc supplierSvc storeSvc; do \
		if [ -f "services/$$service/api/proto/*/v1/*.proto" ]; then \
			echo "  ✅ $$service has service-owned proto files"; \
		else \
			echo "  ❌ $$service missing proto files"; \
		fi; \
	done
	@echo ""
	@echo "Protobuf architecture verification complete!"
