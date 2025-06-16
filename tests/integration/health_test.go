package integration

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// TestServiceHealth tests health endpoints of all services
func TestServiceHealth(t *testing.T) {
	services := map[string]string{
		"product-service":   "localhost:50053",
		"inventory-service": "localhost:50054",
		"order-service":     "localhost:50055",
		"user-service":      "localhost:50056",
		"supplier-service":  "localhost:50057",
	}

	for serviceName, address := range services {
		t.Run(serviceName, func(t *testing.T) {
			testServiceHealth(t, serviceName, address)
		})
	}
}

func testServiceHealth(t *testing.T, serviceName, address string) {
	// Create gRPC connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to connect to %s at %s: %v", serviceName, address, err)
	}
	defer conn.Close()

	// Create health client
	healthClient := grpc_health_v1.NewHealthClient(conn)

	// Check health
	resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Health check failed for %s: %v", serviceName, err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("Service %s is not serving, status: %v", serviceName, resp.Status)
	}

	t.Logf("✅ %s is healthy", serviceName)
}
