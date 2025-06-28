package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	storev1 "github.com/leonvanderhaeghen/stockplatform/services/storeSvc/api/gen/go/proto/store/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/database"
	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/service"
)

// Server represents the gRPC server
type Server struct {
	config   *config.Config
	database *database.Database
	grpcSrv  *grpc.Server
}

// New creates a new server instance
func New(cfg *config.Config, db *database.Database) *Server {
	return &Server{
		config:   cfg,
		database: db,
	}
}

// Start starts the gRPC server
func (s *Server) Start() error {
	// Create gRPC server
	s.grpcSrv = grpc.NewServer()

	// Register store service
	storeService := service.NewStoreService(s.database)
	storev1.RegisterStoreServiceServer(s.grpcSrv, storeService)

	// Start listening
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	log.Printf("Store service gRPC server starting on %s", addr)

	// Start serving
	if err := s.grpcSrv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() {
	if s.grpcSrv != nil {
		log.Println("Stopping gRPC server...")
		s.grpcSrv.GracefulStop()
	}
}
