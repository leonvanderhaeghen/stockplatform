package grpcclient

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	productv1 "github.com/leonvanderhaeghen/stockplatform/gen/go/product/v1"
)

// NewProductClient creates a new gRPC client for the product service
func NewProductClient(addr string, timeout time.Duration) (*grpc.ClientConn, productv1.ProductServiceClient, error) {
	// Configure logger with better output
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return nil, nil, err
	}
	defer logger.Sync()

	logger.Info("Connecting to gRPC server", 
		zap.String("address", addr),
		zap.Duration("timeout", timeout),
	)

	// Set up a connection to the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Error("Failed to connect to gRPC server", 
			zap.String("address", addr),
			zap.Error(err),
		)
		return nil, nil, err
	}

	logger.Info("Successfully connected to gRPC server")
	client := productv1.NewProductServiceClient(conn)
	
	return conn, client, nil
}

// TestGetProduct tests the GetProduct endpoint with the provided client
func TestGetProduct(client productv1.ProductServiceClient, id string) (*productv1.GetProductResponse, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	
	logger.Info("Testing GetProduct", zap.String("id", id))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetProduct(ctx, &productv1.GetProductRequest{Id: id})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			logger.Error("gRPC error",
				zap.String("code", st.Code().String()),
				zap.String("message", st.Message()),
				zap.Error(err),
			)
		} else {
			logger.Error("Non-gRPC error", zap.Error(err))
		}
		return nil, err
	}

	logger.Info("GetProduct response",
		zap.Any("product", resp.GetProduct()),
	)
	
	return resp, nil
}
