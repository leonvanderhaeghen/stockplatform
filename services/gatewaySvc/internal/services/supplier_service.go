package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
)

// SupplierServiceImpl implements the SupplierService interface
type SupplierServiceImpl struct {
	client *grpcclient.SupplierClient
	logger *zap.Logger
}

// NewSupplierService creates a new instance of SupplierServiceImpl
func NewSupplierService(supplierServiceAddr string, logger *zap.Logger) (SupplierService, error) {
	// Create a gRPC client
	client, err := grpcclient.NewSupplierClient(supplierServiceAddr, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier client: %w", err)
	}

	return &SupplierServiceImpl{
		client: client,
		logger: logger.Named("supplier_service"),
	}, nil
}

// CreateSupplier creates a new supplier
func (s *SupplierServiceImpl) CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.Supplier, error) {
	s.logger.Debug("CreateSupplier", 
		zap.String("name", req.GetName()),
	)
	
	supplier, err := s.client.CreateSupplier(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create supplier",
			zap.String("name", req.GetName()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}
	return supplier, nil
}

// GetSupplier gets a supplier by ID
func (s *SupplierServiceImpl) GetSupplier(ctx context.Context, id string) (*supplierv1.Supplier, error) {
	s.logger.Debug("GetSupplier",
		zap.String("id", id),
	)
	
	supplier, err := s.client.GetSupplier(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get supplier",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}
	return supplier, nil
}

// UpdateSupplier updates a supplier
func (s *SupplierServiceImpl) UpdateSupplier(ctx context.Context, req *supplierv1.UpdateSupplierRequest) (*supplierv1.Supplier, error) {
	s.logger.Debug("UpdateSupplier",
		zap.String("id", req.GetId()),
		zap.String("name", req.GetName()),
	)
	
	supplier, err := s.client.UpdateSupplier(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update supplier",
			zap.String("id", req.GetId()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}
	return supplier, nil
}

// DeleteSupplier deletes a supplier by ID
func (s *SupplierServiceImpl) DeleteSupplier(ctx context.Context, id string) error {
	s.logger.Debug("DeleteSupplier",
		zap.String("id", id),
	)
	
	success, err := s.client.DeleteSupplier(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete supplier",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete supplier: %w", err)
	}
	
	if !success {
		return fmt.Errorf("supplier deletion failed")
	}
	return nil
}

// ListSuppliers lists suppliers with pagination and search
func (s *SupplierServiceImpl) ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*supplierv1.Supplier, int32, error) {
	s.logger.Debug("ListSuppliers",
		zap.Int32("page", page),
		zap.Int32("pageSize", pageSize),
		zap.String("search", search),
	)
	
	suppliers, total, err := s.client.ListSuppliers(ctx, page, pageSize, search)
	if err != nil {
		s.logger.Error("Failed to list suppliers",
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to list suppliers: %w", err)
	}
	return suppliers, total, nil
}

// Close closes the client connection
func (s *SupplierServiceImpl) Close() error {
	s.logger.Debug("Closing supplier service connection")
	return s.client.Close()
}

// ListAdapters lists all available supplier adapters
func (s *SupplierServiceImpl) ListAdapters(ctx context.Context) ([]*supplierv1.SupplierAdapter, error) {
	s.logger.Debug("ListAdapters")
	
	adapters, err := s.client.ListAdapters(ctx)
	if err != nil {
		s.logger.Error("Failed to list adapters",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list adapters: %w", err)
	}
	return adapters, nil
}

// GetAdapterCapabilities gets capabilities of a supplier adapter
func (s *SupplierServiceImpl) GetAdapterCapabilities(ctx context.Context, adapterName string) (*supplierv1.AdapterCapabilities, error) {
	s.logger.Debug("GetAdapterCapabilities",
		zap.String("adapterName", adapterName),
	)
	
	capabilities, err := s.client.GetAdapterCapabilities(ctx, adapterName)
	if err != nil {
		s.logger.Error("Failed to get adapter capabilities",
			zap.String("adapterName", adapterName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get adapter capabilities: %w", err)
	}
	return capabilities, nil
}

// TestAdapterConnection tests connection to a supplier adapter
func (s *SupplierServiceImpl) TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error {
	s.logger.Debug("TestAdapterConnection",
		zap.String("adapterName", adapterName),
	)
	
	err := s.client.TestAdapterConnection(ctx, adapterName, config)
	if err != nil {
		s.logger.Error("Failed to test adapter connection",
			zap.String("adapterName", adapterName),
			zap.Error(err),
		)
		return fmt.Errorf("failed to test adapter connection: %w", err)
	}
	return nil
}

// SyncProducts synchronizes products from a supplier
func (s *SupplierServiceImpl) SyncProducts(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error) {
	s.logger.Debug("SyncProducts",
		zap.String("supplierID", supplierID),
	)
	
	jobID, err := s.client.SyncProducts(ctx, supplierID, options)
	if err != nil {
		s.logger.Error("Failed to sync products",
			zap.String("supplierID", supplierID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to sync products: %w", err)
	}
	return jobID, nil
}

// SyncInventory synchronizes inventory from a supplier
func (s *SupplierServiceImpl) SyncInventory(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error) {
	s.logger.Debug("SyncInventory",
		zap.String("supplierID", supplierID),
	)
	
	jobID, err := s.client.SyncInventory(ctx, supplierID, options)
	if err != nil {
		s.logger.Error("Failed to sync inventory",
			zap.String("supplierID", supplierID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to sync inventory: %w", err)
	}
	return jobID, nil
}
