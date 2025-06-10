package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/spf13/viper"
	supplierv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/supplier/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/bootstrap"
	mongorepo "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/infrastructure/mongodb"
	grpchandler "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/interfaces/grpc"
)

type Config struct {
	Database struct {
		URI  string `mapstructure:"uri"`
		Name string `mapstructure:"name"`
	} `mapstructure:"database"`
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg, err := LoadConfig("./config/config.yaml")
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Disconnect(context.Background())

	// Initialize repository
	supplierRepo := mongorepo.NewSupplierRepository(mongoClient.Database(cfg.Database.Name), "suppliers")

	// Initialize services
	supplierService := application.NewSupplierService(supplierRepo)

	// Register supplier adapters
	bootstrap.RegisterAdapters(supplierService)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	supplierServer := grpchandler.NewSupplierServer(supplierService)
	supplierv1.RegisterSupplierServiceServer(grpcServer, supplierServer)

	// Start gRPC server
	addr := ":" + cfg.Server.Port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// Graceful shutdown
	go func() {
		logger.Info("Starting gRPC server", zap.String("addr", addr))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Gracefully stop the server
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		logger.Warn("Server forced to shutdown")
		grpcServer.Stop()
	case <-stopped:
		logger.Info("Server stopped gracefully")
	}
}
