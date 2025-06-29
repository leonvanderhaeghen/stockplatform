package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/database"
	grpchandlers "github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/interfaces/grpc"
	supplierclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/supplier"
	inventoryclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/inventory"
)

// Server holds the gRPC server and its dependencies
type Server struct {
	grpcServer     *grpc.Server
	healthServer   *health.Server
	config         *config.Config
	database       *database.Database
	logger         *zap.Logger
	supplierClient *supplierclient.Client
	inventoryClient *inventoryclient.Client
}

// New creates a new server instance
func New(cfg *config.Config, db *database.Database, logger *zap.Logger) *Server {
	return &Server{
		config:   cfg,
		database: db,
		logger:   logger,
	}
}

// Initialize sets up the gRPC server with all services
func (s *Server) Initialize() error {
	// Create gRPC server
	s.grpcServer = grpc.NewServer()
	s.healthServer = health.NewServer()

	// Initialize supplier gRPC client
	supplierConfig := supplierclient.Config{
		Address: s.config.SupplierServiceAddr,
	}
	supplierClient, err := supplierclient.New(supplierConfig, s.logger)
	if err != nil {
		s.logger.Error("Failed to initialize supplier client", zap.Error(err))
		return err
	}
	s.supplierClient = supplierClient

	// Initialize inventory gRPC client
	inventoryConfig := inventoryclient.Config{
		Address: s.config.InventoryServiceAddr,
	}
	inventoryClient, err := inventoryclient.New(inventoryConfig, s.logger)
	if err != nil {
		s.logger.Error("Failed to initialize inventory client", zap.Error(err))
		return err
	}
	s.inventoryClient = inventoryClient

	// Initialize application services
	productService := application.NewProductService(s.database.ProductRepo, supplierClient, inventoryClient, s.logger)
	categoryService := application.NewCategoryService(s.database.CategoryRepo, s.logger)

	// Register gRPC services
	productServer := grpchandlers.NewProductServer(productService, categoryService, s.logger)
	productv1.RegisterProductServiceServer(s.grpcServer, productServer)

	// Register health check service
	grpc_health_v1.RegisterHealthServer(s.grpcServer, s.healthServer)
	s.healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable reflection for development
	reflection.Register(s.grpcServer)

	s.logger.Info("gRPC server initialized successfully")
	return nil
}

// Start starts the gRPC server
func (s *Server) Start() error {
	// Create listener
	lis, err := net.Listen("tcp", ":"+s.config.GRPCPort)
	if err != nil {
		s.logger.Error("Failed to create listener", zap.Error(err))
		return err
	}

	// Start server in a goroutine
	go func() {
		s.logger.Info("Starting gRPC server", zap.String("port", s.config.GRPCPort))
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Error("gRPC server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Received shutdown signal, stopping server...")
	return s.Stop()
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() error {
	s.logger.Info("Shutting down gRPC server...")

	// Close supplier client connection
	if s.supplierClient != nil {
		if err := s.supplierClient.Close(); err != nil {
			s.logger.Error("Failed to close supplier client", zap.Error(err))
		}
	}

	// Close inventory client connection
	if s.inventoryClient != nil {
		if err := s.inventoryClient.Close(); err != nil {
			s.logger.Error("Failed to close inventory client", zap.Error(err))
		}
	}

	// Graceful shutdown with timeout
	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
	case <-time.After(10 * time.Second):
		s.logger.Warn("Graceful shutdown timed out, forcing stop")
		s.grpcServer.Stop()
	}

	return nil
}
