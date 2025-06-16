package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/rest"
	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/services"
)

// ServiceClients holds all service client instances
type ServiceClients struct {
	ProductSvc   services.ProductService
	InventorySvc services.InventoryService
	OrderSvc     services.OrderService
	UserSvc      services.UserService
	SupplierSvc  services.SupplierService
}

// Server holds the REST server and its dependencies
type Server struct {
	restServer *rest.Server
	config     *config.Config
	logger     *zap.Logger
}

// New creates a new server instance
func New(cfg *config.Config, logger *zap.Logger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
	}
}

// Initialize sets up the server with all services
func (s *Server) Initialize() error {
	// Initialize service clients
	serviceClients, err := s.initServices()
	if err != nil {
		return err
	}

	// Initialize REST server
	s.restServer = rest.NewServer(
		serviceClients.ProductSvc,
		serviceClients.InventorySvc,
		serviceClients.OrderSvc,
		serviceClients.UserSvc,
		serviceClients.SupplierSvc,
		s.config.JWT.Secret,
		s.config.Server.Port,
		s.logger,
	)

	s.restServer.SetupRoutes()

	s.logger.Info("Server initialized successfully")
	return nil
}

// Start starts the REST server with graceful shutdown
func (s *Server) Start() error {
	// Start server in a goroutine
	go func() {
		s.logger.Info("Starting REST server")
		if err := s.restServer.Start(); err != nil {
			s.logger.Fatal("Failed to start REST server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown
	return s.restServer.Shutdown(ctx)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.restServer.Shutdown(ctx)
}

// initServices initializes all service clients
func (s *Server) initServices() (*ServiceClients, error) {
	productSvc, err := services.NewProductService(s.config.Services.ProductAddr, s.logger)
	if err != nil {
		return nil, err
	}

	inventorySvc, err := services.NewInventoryService(s.config.Services.InventoryAddr, s.logger)
	if err != nil {
		return nil, err
	}

	orderSvc, err := services.NewOrderService(s.config.Services.OrderAddr, s.logger)
	if err != nil {
		return nil, err
	}

	userSvc, err := services.NewUserService(s.config.Services.UserAddr, s.logger)
	if err != nil {
		return nil, err
	}

	supplierSvc, err := services.NewSupplierService(s.config.Services.SupplierAddr, s.logger)
	if err != nil {
		return nil, err
	}

	return &ServiceClients{
		ProductSvc:   productSvc,
		InventorySvc: inventorySvc,
		OrderSvc:     orderSvc,
		UserSvc:      userSvc,
		SupplierSvc:  supplierSvc,
	}, nil
}
