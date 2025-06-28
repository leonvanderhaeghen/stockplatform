package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
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
func (s *SupplierServiceImpl) CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.Supplier, error) {
	s.logger.Debug("CreateSupplier", 
		zap.String("name", req.GetName()),
	)
	
	resp, err := s.client.CreateSupplier(ctx, req)
	if err != nil {
		// error handled below
	}
	supplier := resp.GetSupplier()
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
	
	resp, err := s.client.GetSupplier(ctx, &supplierv1.GetSupplierRequest{Id: id})
	supplier := resp.GetSupplier()
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
	
	resp, err := s.client.UpdateSupplier(ctx, req)
	supplier := resp.GetSupplier()
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
	
	_, err := s.client.DeleteSupplier(ctx, &supplierv1.DeleteSupplierRequest{Id: id})
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
func (s *SupplierServiceImpl) ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*supplierv1.Supplier, int32, error) {
	s.logger.Debug("ListSuppliers",
		zap.Int32("page", page),
		zap.Int32("pageSize", pageSize),
		zap.String("search", search),
	)
	
	resp, err := s.client.ListSuppliers(ctx, &supplierv1.ListSuppliersRequest{
		Page: page,
		PageSize: pageSize,
		Search: search,
	})
	suppliers := resp.GetSuppliers()
	total := resp.GetTotal()
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
	
	req := &supplierv1.ListAdaptersRequest{}
	adaptersResp, err := s.client.ListAdapters(ctx, req)
	adapters := adaptersResp.GetAdapters()
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
	
	req := &supplierv1.GetAdapterCapabilitiesRequest{AdapterName: adapterName}
	capResp, err := s.client.GetAdapterCapabilities(ctx, req)
	capabilities := capResp.GetCapabilities()
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
	
	req := &supplierv1.TestAdapterConnectionRequest{AdapterName: adapterName, Config: config}
	_, err := s.client.TestAdapterConnection(ctx, req)
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
	
	req := &supplierv1.SyncProductsRequest{SupplierId: supplierID, Options: options}
	prodResp, err := s.client.SyncProducts(ctx, req)
	jobID := prodResp.GetJobId()
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
	
	req := &supplierv1.SyncInventoryRequest{SupplierId: supplierID, Options: options}
	invResp, err := s.client.SyncInventory(ctx, req)
	jobID := invResp.GetJobId()
	if err != nil {
		s.logger.Error("Failed to sync inventory",
			zap.String("supplierID", supplierID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to sync inventory: %w", err)
	}
	return jobID, nil
}
