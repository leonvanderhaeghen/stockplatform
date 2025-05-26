package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
	productv1 "github.com/leonvanderhaeghen/stockplatform/gen/go/product/v1"
)

func main() {
	// Configure enhanced logger
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a custom options slice
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024 * 1024 * 10), // 10MB
			grpc.MaxCallSendMsgSize(1024 * 1024 * 10), // 10MB
		),
	}

	// Set up a connection to the server
	addr := "localhost:50053"
	if envAddr := os.Getenv("GRPC_SERVER_ADDR"); envAddr != "" {
		addr = envAddr
	}

	timeout := 5 * time.Second
	
	logger.Info("Connecting to gRPC server", 
		zap.String("address", addr),
		zap.Duration("timeout", timeout),
	)

	// Connect with context and custom options
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, options...)
	if err != nil {
		logger.Fatal("Failed to connect to gRPC server", 
			zap.String("address", addr),
			zap.Error(err),
		)
	}
	defer conn.Close()

	logger.Info("Successfully connected to gRPC server")
	
	// Create client and run tests
	client := productv1.NewProductServiceClient(conn)
	
	// Test GetProduct with different IDs
	testIDs := []string{"test-id-1", "test-id-2", "non-existent-id"}
	for _, id := range testIDs {
		resp, err := grpcclient.TestGetProduct(client, id)
		if err != nil {
			logger.Error("GetProduct failed",
				zap.String("id", id),
				zap.Error(err),
			)
			continue
		}
		
		logger.Info("GetProduct succeeded",
			zap.String("id", id),
			zap.Any("product", resp.GetProduct()),
		)
	}
}
