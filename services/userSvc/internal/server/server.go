package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/database"
	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/handlers"
	grpchandlers "github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/interfaces/grpc"
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

	// Initialize services
	userService := application.NewUserService(
		s.database.UserRepo,
		s.database.AddressRepo,
		s.config.JWTSecret,
		s.logger,
	)
	_ = application.NewPermissionService(s.database.PermissionRepo) // Initialize but don't use directly

	// Initialize AuthService
	authConfig := &application.AuthConfig{
		JWTSecret:       []byte(s.config.JWTSecret),
		TokenDuration:   24 * time.Hour,
		RefreshDuration: 7 * 24 * time.Hour,
	}
	authService, err := application.NewAuthService(s.database.UserRepo, authConfig, s.logger)
	if err != nil {
		return fmt.Errorf("failed to initialize auth service: %w", err)
	}

	// Initialize user auth service
	userAuthService, err := application.NewUserAuthService(
		userService,
		s.config.OrderSvcURL,
		s.logger,
	)
	if err != nil {
		return err
	}
	defer userAuthService.Close()

	// Initialize gRPC handlers
	userServer := grpchandlers.NewUserServer(userService, s.logger)
	authHandler := handlers.NewAuthHandler(authService, s.logger)

	// Register gRPC services
	userv1.RegisterUserServiceServer(s.grpcServer, userServer)
	userv1.RegisterAuthServiceServer(s.grpcServer, authHandler)

	// Register health check service
	healthServer := grpchandlers.NewHealthServer(s.logger)
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
