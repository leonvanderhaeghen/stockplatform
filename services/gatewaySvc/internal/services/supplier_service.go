package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	supplierclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/supplier"
)

// SupplierServiceImpl implements the SupplierService interface
type SupplierServiceImpl struct {
	client *supplierclient.Client
	logger *zap.Logger
}

// NewSupplierService creates a new instance of SupplierServiceImpl
func NewSupplierService(supplierServiceAddr string, logger *zap.Logger) (SupplierService, error) {
	// Create a gRPC client via new abstraction
	supCfg := supplierclient.Config{Address: supplierServiceAddr}
	client, err := supplierclient.New(supCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier client: %w", err)
	}

	return &SupplierServiceImpl{
		client: client,
		logger: logger.Named("supplier_service"),
	}, nil
}

// CreateSupplier creates a new supplier
func (s *SupplierServiceImpl) CreateSupplier(ctx context.Context, name, contactPerson, email, phone, address, city, state, country, postalCode, taxID, website, currency, paymentTerms string, leadTimeDays int32, metadata map[string]string) (interface{}, error) {
	s.logger.Debug("CreateSupplier", 
		zap.String("name", name),
		zap.String("email", email),
	)
	
	resp, err := s.client.CreateSupplier(ctx, name, contactPerson, email, phone, address, city, state, country, postalCode, taxID, website, currency, paymentTerms, leadTimeDays, metadata)
	if err != nil {
		s.logger.Error("Failed to create supplier",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}
	return resp, nil
}

// GetSupplier gets a supplier by ID
func (s *SupplierServiceImpl) GetSupplier(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetSupplier",
		zap.String("id", id),
	)
	
	resp, err := s.client.GetSupplier(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get supplier",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}
	return resp, nil
}

// UpdateSupplier updates a supplier
func (s *SupplierServiceImpl) UpdateSupplier(ctx context.Context, id, name, contactPerson, email, phone, address, city, state, country, postalCode, taxID, website, currency, paymentTerms string, leadTimeDays int32, metadata map[string]string) (interface{}, error) {
	s.logger.Debug("UpdateSupplier",
		zap.String("id", id),
		zap.String("name", name),
	)
	
	resp, err := s.client.UpdateSupplier(ctx, id, name, contactPerson, email, phone, address, city, state, country, postalCode, taxID, website, currency, paymentTerms, leadTimeDays, metadata)
	if err != nil {
		s.logger.Error("Failed to update supplier",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}
	return resp, nil
}

// DeleteSupplier deletes a supplier by ID
func (s *SupplierServiceImpl) DeleteSupplier(ctx context.Context, id string) error {
	s.logger.Debug("DeleteSupplier",
		zap.String("id", id),
	)
	
	err := s.client.DeleteSupplier(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete supplier",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete supplier: %w", err)
	}
	
	
	return nil
}

// ListSuppliers lists suppliers with pagination and search
func (s *SupplierServiceImpl) ListSuppliers(ctx context.Context, page, pageSize int32, search string) (interface{}, error) {
	s.logger.Debug("ListSuppliers",
		zap.Int32("page", page),
		zap.Int32("pageSize", pageSize),
		zap.String("search", search),
	)
	
	// Client interface expects (pageSize, pageToken/search) - adjust parameters
	resp, err := s.client.ListSuppliers(ctx, pageSize, search)
	if err != nil {
		s.logger.Error("Failed to list suppliers",
			zap.Int32("page", page),
			zap.Int32("pageSize", pageSize),
			zap.String("search", search),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}
	return resp, nil
}

// Close closes the client connection
func (s *SupplierServiceImpl) Close() error {
	s.logger.Debug("Closing supplier service connection")
	return s.client.Close()
}

// ListAdapters lists all available supplier adapters
func (s *SupplierServiceImpl) ListAdapters(ctx context.Context) (interface{}, error) {
	s.logger.Debug("ListAdapters")
	
	resp, err := s.client.ListAdapters(ctx)
	if err != nil {
		s.logger.Error("Failed to list adapters",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list adapters: %w", err)
	}
	return resp, nil
}

// GetAdapterCapabilities gets capabilities of a supplier adapter
func (s *SupplierServiceImpl) GetAdapterCapabilities(ctx context.Context, adapterName string) (interface{}, error) {
	s.logger.Debug("GetAdapterCapabilities",
		zap.String("adapterName", adapterName),
	)
	
	resp, err := s.client.GetAdapterCapabilities(ctx, adapterName)
	if err != nil {
		s.logger.Error("Failed to get adapter capabilities",
			zap.String("adapterName", adapterName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get adapter capabilities: %w", err)
	}
	return resp, nil
}

// TestAdapterConnection tests connection to a supplier adapter
func (s *SupplierServiceImpl) TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error {
	s.logger.Debug("TestAdapterConnection",
		zap.String("adapterName", adapterName),
	)
	
	resp, err := s.client.TestAdapterConnection(ctx, adapterName, config)
	if err != nil {
		s.logger.Error("Failed to test adapter connection",
			zap.String("adapterName", adapterName),
			zap.Error(err),
		)
		return fmt.Errorf("failed to test adapter connection: %w", err)
	}
	
	// Check if connection test was successful
	if resp != nil && !resp.Success {
		return fmt.Errorf("connection test failed: %s", resp.Message)
	}
	
	return nil
}

// SyncProducts synchronizes products from a supplier
func (s *SupplierServiceImpl) SyncProducts(ctx context.Context, supplierID string, fullSync, dryRun bool, batchSize int32) (string, error) {
	s.logger.Debug("SyncProducts",
		zap.String("supplierID", supplierID),
		zap.Bool("fullSync", fullSync),
		zap.Bool("dryRun", dryRun),
	)
	
	resp, err := s.client.SyncProducts(ctx, supplierID, fullSync, dryRun, batchSize)
	if err != nil {
		s.logger.Error("Failed to sync products",
			zap.String("supplierID", supplierID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to sync products: %w", err)
	}

	// Extract JobID from response
	if resp != nil {
		return resp.JobID, nil
	}
	return "", fmt.Errorf("received nil response from sync products")
}

// SyncInventory synchronizes inventory from a supplier
func (s *SupplierServiceImpl) SyncInventory(ctx context.Context, supplierID string, fullSync, dryRun bool, batchSize int32) (string, error) {
	s.logger.Debug("SyncInventory",
		zap.String("supplierID", supplierID),
		zap.Bool("fullSync", fullSync),
		zap.Bool("dryRun", dryRun),
	)
	
	resp, err := s.client.SyncInventory(ctx, supplierID, fullSync, dryRun, batchSize)
	if err != nil {
		s.logger.Error("Failed to sync inventory",
			zap.String("supplierID", supplierID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to sync inventory: %w", err)
	}

	// Extract JobID from response
	if resp != nil {
		return resp.JobID, nil
	}
	return "", fmt.Errorf("received nil response from sync inventory")
}
