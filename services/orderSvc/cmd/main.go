package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	orderv1 "stockplatform/pkg/gen/order/v1"
	"stockplatform/services/orderSvc/internal/application"
	grpchandler "stockplatform/services/orderSvc/internal/interfaces/grpc"
	"stockplatform/services/orderSvc/internal/infrastructure/mongodb"
)

// Config holds the application configuration
type Config struct {
	GRPCPort              string `env:"GRPC_PORT,default=50055"`
	MongoURI              string `env:"MONGO_URI,default=mongodb://localhost:27017"`
	ProductServiceAddr    string `env:"PRODUCT_SERVICE_ADDR,default=localhost:50053"`
	InventoryServiceAddr  string `env:"INVENTORY_SERVICE_ADDR,default=localhost:50054"`
}

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting order service...")

	// Load configuration
	config := Config{
		GRPCPort:             "50055",
		MongoURI:             "mongodb://localhost:27017",
		ProductServiceAddr:   "localhost:50053",
		InventoryServiceAddr: "localhost:50054",
	}
	
	// Check for environment variables
	if port := os.Getenv("GRPC_PORT"); port != "" {
		config.GRPCPort = port
	}
	
	if mongoURI := os.Getenv("MONGO_URI"); mongoURI != "" {
		config.MongoURI = mongoURI
	}
	
	if productSvcAddr := os.Getenv("PRODUCT_SERVICE_ADDR"); productSvcAddr != "" {
		config.ProductServiceAddr = productSvcAddr
	}
	
	if inventorySvcAddr := os.Getenv("INVENTORY_SERVICE_ADDR"); inventorySvcAddr != "" {
		config.InventoryServiceAddr = inventorySvcAddr
	}
	
	logger.Info("Configuration loaded", 
		zap.String("grpc_port", config.GRPCPort),
		zap.String("mongo_uri", config.MongoURI),
		zap.String("product_service_addr", config.ProductServiceAddr),
		zap.String("inventory_service_addr", config.InventoryServiceAddr),
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

	// Initialize repositories
	dbName := "stockplatform"
	collectionName := "orders"
	logger.Info("Initializing MongoDB repository", 
		zap.String("database", dbName),
		zap.String("collection", collectionName),
	)
	db := mongoClient.Database(dbName)
	orderRepo := mongodb.NewOrderRepository(db, collectionName, logger)

	// Initialize application services
	orderSvc := application.NewOrderService(orderRepo, logger)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC services
	orderv1.RegisterOrderServiceServer(
		grpcServer,
		grpchandler.NewOrderServer(orderSvc, logger),
	)

	// Enable reflection for gRPC CLI tools
	reflection.Register(grpcServer)

	// Start gRPC server
	addr := "0.0.0.0:" + config.GRPCPort
	logger.Info("Starting gRPC server", 
		zap.String("address", addr),
	)

	// Create a TCP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", 
			zap.String("address", addr),
			zap.Error(err),
		)
	}

	// Log listener details
	logger.Info("TCP listener created successfully",
		zap.String("local_address", listener.Addr().String()),
		zap.String("network", listener.Addr().Network()),
	)

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Fatal("Failed to serve gRPC server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Gracefully stop the gRPC server
	grpcServer.GracefulStop()
	logger.Info("Server stopped gracefully")
}
