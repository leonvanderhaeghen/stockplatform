.PHONY: deps generate build test lint clean build-clients run-grpc-client run-grpc-client-enhanced run-tcp-client run-test-grpc-client run-test-grpc-minimal

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BUF=buf

# Directories
PROTO_DIR=./proto
SERVICES_DIR=./services

# Install dependencies
deps:
	$(GOGET) -u google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		github.com/bufbuild/buf/cmd/buf \
		github.com/fullstorydev/grpcurl/cmd/grpcurl

# Generate protobuf and gRPC code
generate:
	$(BUF) generate

# Build all services
build:
	@echo "Building all services..."
	@for dir in $(wildcard $(SERVICES_DIR)/*/); do \
		echo "Building $$(basename $$dir)"; \
		cd $$dir && $(GOBUILD) -o bin/$$(basename $$dir) . || exit 1; \
		echo "Build complete for $$(basename $$dir)"; \
	done

# Run tests
test:
	$(GOTEST) -v ./...

# Run linters
lint:
	golangci-lint run ./...

# Build all clients
build-clients:
	@echo "Building all client applications..."
	@for dir in $(wildcard ./cmd/*/); do \
		echo "Building $$(basename $$dir)"; \
		$(GOBUILD) -o bin/$$(basename $$dir) ./cmd/$$(basename $$dir) || exit 1; \
		echo "Build complete for $$(basename $$dir)"; \
	done

# Run gRPC client
run-grpc-client:
	$(GOBUILD) -o bin/grpc_client ./cmd/grpc_client
	./bin/grpc_client

# Run enhanced gRPC client
run-grpc-client-enhanced:
	$(GOBUILD) -o bin/grpc_client_enhanced ./cmd/grpc_client_enhanced
	./bin/grpc_client_enhanced

# Run TCP client
run-tcp-client:
	$(GOBUILD) -o bin/tcp_client ./cmd/tcp_client
	./bin/tcp_client

# Run test gRPC client
run-test-grpc-client:
	$(GOBUILD) -o bin/test_grpc_client ./cmd/test_grpc_client
	./bin/test_grpc_client

# Run test gRPC minimal client
run-test-grpc-minimal:
	$(GOBUILD) -o bin/test_grpc_minimal ./cmd/test_grpc_minimal
	./bin/test_grpc_minimal

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/
	@for dir in $(wildcard $(SERVICES_DIR)/*/); do \
		echo "Cleaning $$(basename $$dir)"; \
		rm -f $$dir/bin/*; \
	done

# Install tools
tools:
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Run the application in development mode
dev: generate build
	@echo "Starting development environment..."
	@docker-compose up --build -d

# Stop the development environment
dev-stop:
	@echo "Stopping development environment..."
	@docker-compose down

# View logs
dev-logs:
	@docker-compose logs -f

# Help
dev-help:
	@echo "Available commands:"
	@echo "  make dev       - Start the development environment"
	@echo "  make dev-stop  - Stop the development environment"
	@echo "  make dev-logs  - View logs from the development environment"
	@echo "  make generate  - Generate code from proto files"
	@echo "  make build     - Build all services"
	@echo "  make test      - Run tests"
	@echo "  make lint      - Run linters"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make deps      - Install dependencies"
	@echo "  make tools     - Install development tools"
