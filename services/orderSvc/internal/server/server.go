package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	orderv1 "github.com/leonvanderhaeghen/stockplatform/services/orderSvc/api/gen/go/proto/order/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/database"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
	grpcintf "github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/interfaces/grpc"
)

// Server holds the gRPC server and its dependencies
type Server struct {
	grpcServer *grpc.Server
	config     *config.Config
	database   *database.Database
	logger     *zap.Logger
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

	// Initialize event publisher (could be Kafka or in-memory; nil for now)
	var publisher domain.EventPublisher
	// Create event service
	eventService := application.NewEventService(publisher, s.logger)

	// Initialize order service
	orderService := application.NewOrderService(s.database.OrderRepo, eventService, s.logger)

	// Create service config for POS transactions
	serviceConfig := &domain.ServiceConfig{
		InventoryServiceAddr: s.config.InventoryServiceAddr,
	}

	// Initialize order inventory service
	orderInventoryService, err := application.NewOrderInventoryService(
		orderService,
		s.config.InventoryServiceAddr,
		s.logger,
	)
	if err != nil {
		return err
	}
	defer orderInventoryService.Close()

	// Initialize POS transaction service
	posTransactionService := application.NewPOSTransactionService(orderService, serviceConfig)

	// Initialize gRPC handlers
	orderServer := grpcintf.NewOrderServer(orderService, posTransactionService, s.logger)

	// Register gRPC services
	orderv1.RegisterOrderServiceServer(s.grpcServer, orderServer)

	// Register health check service
	healthServer := grpcintf.NewHealthServer(s.logger)
	grpc_health_v1.RegisterHealthServer(s.grpcServer, healthServer)

	// Enable reflection for development
	reflection.Register(s.grpcServer)

	s.logger.Info("gRPC server initialized successfully")
	return nil
}

// Start starts the gRPC server with graceful shutdown
func (s *Server) Start() error {
	// Create listener
	lis, err := net.Listen("tcp", ":"+s.config.GRPCPort)
	if err != nil {
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

	s.logger.Info("Shutting down gRPC server...")

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

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("Graceful shutdown timed out, forcing stop")
		s.grpcServer.Stop()
		return ctx.Err()
	}
}
