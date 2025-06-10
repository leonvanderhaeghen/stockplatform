package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	productv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/product/v1"
	supplierv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/supplier/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/infrastructure/mongodb"
	handlergrpc "github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/interfaces/grpc"
	handlerhttp "github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/interfaces/http"
)

// Config holds the application configuration
type Config struct {
	GRPCPort string `env:"GRPC_PORT,default=50053"`
	HTTPPort string `env:"HTTP_PORT,default=3001"`
	MongoURI string `env:"MONGO_URI,default=mongodb://localhost:27017"`
}

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting product service...")

	// Load configuration from environment variables
	config := Config{
		GRPCPort: os.Getenv("GRPC_PORT"),
		HTTPPort: os.Getenv("HTTP_PORT"),
		MongoURI: os.Getenv("MONGO_URI"),
	}

	// Set default values if environment variables are not set
	if config.GRPCPort == "" {
		config.GRPCPort = "50053"
	}
	if config.HTTPPort == "" {
		config.HTTPPort = "3001"
	}
	if config.MongoURI == "" {
		config.MongoURI = "mongodb://localhost:27017"
	}
	logger.Info("Configuration loaded",
		zap.String("grpc_port", config.GRPCPort),
		zap.String("http_port", config.HTTPPort),
		zap.String("mongo_uri", config.MongoURI),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	logger.Info("Connecting to MongoDB...",
		zap.String("uri", config.MongoURI),
	)

	// Set up client options with timeouts and retry settings
	clientOptions := options.Client()
	clientOptions.ApplyURI(config.MongoURI)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetSocketTimeout(30 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)
	clientOptions.SetMaxPoolSize(10)

	logger.Info("MongoDB client options",
		zap.String("uri", config.MongoURI),
		zap.Duration("connect_timeout", 10*time.Second),
		zap.Duration("socket_timeout", 30*time.Second),
	)

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal("Failed to create MongoDB client",
			zap.Error(err),
		)
	}

	// Log MongoDB server status
	serverStatus, err := mongoClient.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to list MongoDB databases",
			zap.Error(err),
		)
	} else {
		logger.Info("Successfully connected to MongoDB server",
			zap.Strings("databases", serverStatus),
		)
	}

	// Set up a goroutine to monitor connection state
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Reduced frequency to reduce log spam
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				start := time.Now()
				err := mongoClient.Ping(ctx, nil)
				duration := time.Since(start)

				if err != nil {
					logger.Error("MongoDB ping failed",
						zap.Error(err),
						zap.Duration("duration", duration),
					)
				} else {
					logger.Debug("MongoDB ping successful",
						zap.Duration("duration", duration),
					)
				}
			case <-ctx.Done():
				logger.Info("Stopping MongoDB connection monitor")
				return
			}
		}
	}()

	// Verify initial connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		logger.Fatal("Failed to ping MongoDB",
			zap.Error(err),
		)
	}

	logger.Info("Successfully connected to MongoDB")

	defer func() {
		logger.Info("Disconnecting from MongoDB...")
		if err := mongoClient.Disconnect(ctx); err != nil {
			logger.Error("Failed to disconnect from MongoDB",
				zap.Error(err),
			)
		} else {
			logger.Info("Successfully disconnected from MongoDB")
		}
	}()

	// This block was removed as we moved the ping check right after connection

	// Initialize repositories
	dbName := "stockplatform"
	logger.Info("Initializing MongoDB repositories",
		zap.String("database", dbName),
	)
	db := mongoClient.Database(dbName)

	productRepo := mongodb.NewProductRepository(db, logger)
	categoryRepo := mongodb.NewCategoryRepository(db, logger)
	supplierRepo := mongodb.NewSupplierRepository(db, logger)

	// Initialize services
	supplierSvc := application.NewSupplierService(supplierRepo, logger)
	productSvc := application.NewProductService(productRepo, supplierSvc, logger)
	categorySvc := application.NewCategoryService(categoryRepo, logger)

	// Initialize gRPC server with interceptors
	grpcServer := grpcserver.NewServer(
		grpcserver.ChainUnaryInterceptor(
			loggingInterceptor(logger),
		),
	)

	// Register services with the gRPC server
	productServer := handlergrpc.NewProductServer(productSvc, categorySvc, logger)
	productv1.RegisterProductServiceServer(grpcServer, productServer)

	// Register supplier service
	supplierServer := handlergrpc.NewSupplierServer(supplierSvc, logger)
	supplierv1.RegisterSupplierServiceServer(grpcServer, supplierServer)

	// Enable reflection for gRPC CLI tools
	reflection.Register(grpcServer)

	// Start gRPC server
	addr := "0.0.0.0:" + config.GRPCPort
	logger.Info("Starting gRPC server",
		zap.String("address", addr),
		zap.String("port", config.GRPCPort),
	)

	// Enable gRPC logging
	grpcserver.EnableTracing = true
	grpcLogger := logger.Named("grpc")
	grpcLogLevel := zap.DebugLevel
	grpcLogger = grpcLogger.WithOptions(zap.IncreaseLevel(grpcLogLevel))
	grpcLogger.Info("gRPC logging enabled", zap.Stringer("level", grpcLogLevel))

	// Log network interfaces for debugging
	logNetworkInterfaces(logger)

	// Try to create a TCP listener with detailed error handling
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		handleListenError(logger, err, addr, config.GRPCPort)
	}

	// Log listener details
	logger.Info("TCP listener created successfully",
		zap.String("local_address", listener.Addr().String()),
		zap.String("network", listener.Addr().Network()),
	)

	// Log when the server starts accepting connections
	logger.Info("gRPC server is ready to accept connections")

	// Log all available network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn("Failed to get network interfaces", zap.Error(err))
	} else {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				logger.Warn("Failed to get addresses for interface",
					zap.String("interface", i.Name),
					zap.Error(err),
				)
				continue
			}
			for _, addr := range addrs {
				logger.Debug("Network interface",
					zap.String("interface", i.Name),
					zap.String("address", addr.String()),
				)
			}
		}
	}

	// Create HTTP server
	httpServer := handlerhttp.NewServer(
		config.HTTPPort,
		productSvc,
		categorySvc,
		logger,
	)

	// Start gRPC server in a goroutine
	grpcErrCh := make(chan error, 1)
	go func() {
		logger.Info("Starting gRPC server", zap.String("address", addr))
		if err := grpcServer.Serve(listener); err != nil {
			grpcErrCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	// Start HTTP server in a goroutine
	httpErrCh := make(chan error, 1)
	go func() {
		logger.Info("Starting HTTP server", zap.String("port", config.HTTPPort))
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			httpErrCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-grpcErrCh:
		logger.Error("gRPC server failed", zap.Error(err))
	case err := <-httpErrCh:
		logger.Error("HTTP server failed", zap.Error(err))
	case sig := <-quit:
		logger.Info("Shutting down server...", zap.String("signal", sig.String()))
	}

	// Create a context with a timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown HTTP server", zap.Error(err))
	}

	// Shutdown gRPC server
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	// Wait for gRPC server to stop or timeout
	select {
	case <-stopped:
		logger.Info("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		grpcServer.Stop()
		logger.Warn("gRPC server forcefully stopped after timeout")
	}

	logger.Info("Server stopped")
}

// handleListenError provides detailed error information for TCP listener errors
func handleListenError(logger *zap.Logger, err error, addr, port string) {
	// Try to get more detailed error information
	if opErr, ok := err.(*net.OpError); ok {
		switch opErr.Op {
		case "listen":
			logger.Fatal("Failed to start TCP listener",
				zap.String("address", addr),
				zap.String("port", port),
				zap.String("operation", opErr.Op),
				zap.String("network", opErr.Net),
				zap.String("error_type", fmt.Sprintf("%T", opErr.Err)),
				zap.Error(opErr.Err),
			)
		case "dial":
			logger.Fatal("Failed to dial address",
				zap.String("address", addr),
				zap.String("operation", opErr.Op),
				zap.String("network", opErr.Net),
				zap.Error(opErr.Err),
			)
		default:
			logger.Fatal("Network operation failed",
				zap.String("operation", opErr.Op),
				zap.String("network", opErr.Net),
				zap.Error(opErr.Err),
			)
		}
	}

	// For non-OpError errors
	logger.Fatal("Failed to create listener",
		zap.String("address", addr),
		zap.String("port", port),
		zap.String("error_type", fmt.Sprintf("%T", err)),
		zap.Error(err),
	)
}

// logNetworkInterfaces logs all available network interfaces and their addresses
func logNetworkInterfaces(logger *zap.Logger) {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.Warn("Failed to get network interfaces", zap.Error(err))
		return
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			logger.Warn("Failed to get addresses for interface",
				zap.String("interface", i.Name),
				zap.Error(err),
			)
			continue
		}

		for _, addr := range addrs {
			logger.Debug("Network interface",
				zap.String("interface", i.Name),
				zap.String("address", addr.String()),
			)
		}
	}
}

// loggingInterceptor logs gRPC requests with method and duration
func loggingInterceptor(logger *zap.Logger) grpcserver.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpcserver.UnaryServerInfo, handler grpcserver.UnaryHandler) (interface{}, error) {
		start := time.Now()
		logger.Info("Processing gRPC request",
			zap.String("method", info.FullMethod),
		)

		// Call the handler
		resp, err := handler(ctx, req)

		// Log the request
		logger.Info("Request processed",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}
