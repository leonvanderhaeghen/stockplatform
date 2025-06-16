package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// HealthServer implements the gRPC health check service
type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	logger *zap.Logger
}

// NewHealthServer creates a new health check server
func NewHealthServer(logger *zap.Logger) *HealthServer {
	return &HealthServer{
		logger: logger.Named("health_server"),
	}
}

// Check performs a health check
func (h *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.logger.Debug("Health check requested", zap.String("service", req.GetService()))
	
	// For now, always return serving
	// In production, you'd check database connectivity, dependencies, etc.
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch performs a streaming health check
func (h *HealthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	// Send initial status
	if err := stream.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}); err != nil {
		h.logger.Error("Failed to send health status", zap.Error(err))
		return status.Error(codes.Internal, "failed to send health status")
	}
	
	// Keep the stream open
	<-stream.Context().Done()
	return stream.Context().Err()
}
