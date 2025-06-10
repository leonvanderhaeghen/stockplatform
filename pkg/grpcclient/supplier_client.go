package grpcclient

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/supplier/v1"
)

// SupplierClient is a client for interacting with the SupplierService
// It's a wrapper around the gRPC client that provides a more convenient interface
// and handles connection management.
type SupplierClient struct {
	client supplierv1.SupplierServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// NewSupplierClient creates a new SupplierClient with the given address and logger
func NewSupplierClient(address string, logger *zap.Logger) (*SupplierClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := supplierv1.NewSupplierServiceClient(conn)

	return &SupplierClient{
		client: client,
		conn:   conn,
		logger: logger.Named("supplier_client"),
	}, nil
}

// Close closes the underlying gRPC connection
func (c *SupplierClient) Close() error {
	return c.conn.Close()
}

// CreateSupplier creates a new supplier
func (c *SupplierClient) CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.Supplier, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.CreateSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create supplier", zap.Error(err))
		return nil, err
	}

	return resp.GetSupplier(), nil
}

// GetSupplier retrieves a supplier by ID
func (c *SupplierClient) GetSupplier(ctx context.Context, id string) (*supplierv1.Supplier, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &supplierv1.GetSupplierRequest{
		Id: id,
	}

	resp, err := c.client.GetSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get supplier", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return resp.GetSupplier(), nil
}

// UpdateSupplier updates an existing supplier
func (c *SupplierClient) UpdateSupplier(ctx context.Context, req *supplierv1.UpdateSupplierRequest) (*supplierv1.Supplier, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.UpdateSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update supplier", zap.String("id", req.GetId()), zap.Error(err))
		return nil, err
	}

	return resp.GetSupplier(), nil
}

// DeleteSupplier deletes a supplier by ID
func (c *SupplierClient) DeleteSupplier(ctx context.Context, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &supplierv1.DeleteSupplierRequest{
		Id: id,
	}

	resp, err := c.client.DeleteSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to delete supplier", zap.String("id", id), zap.Error(err))
		return false, err
	}

	return resp != nil, nil
}

// ListSuppliers lists suppliers with pagination and optional search
func (c *SupplierClient) ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*supplierv1.Supplier, int32, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &supplierv1.ListSuppliersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	resp, err := c.client.ListSuppliers(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list suppliers", zap.Error(err))
		return nil, 0, err
	}

	return resp.GetSuppliers(), resp.GetTotal(), nil
}

// ListAdapters returns all available supplier adapters
func (c *SupplierClient) ListAdapters(ctx context.Context) ([]*supplierv1.SupplierAdapter, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &supplierv1.ListAdaptersRequest{}

	resp, err := c.client.ListAdapters(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list adapters", zap.Error(err))
		return nil, err
	}

	return resp.GetAdapters(), nil
}

// GetAdapterCapabilities returns the capabilities of a specific adapter
func (c *SupplierClient) GetAdapterCapabilities(ctx context.Context, adapterName string) (*supplierv1.AdapterCapabilities, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &supplierv1.GetAdapterCapabilitiesRequest{
		AdapterName: adapterName,
	}

	resp, err := c.client.GetAdapterCapabilities(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get adapter capabilities", zap.String("adapter", adapterName), zap.Error(err))
		return nil, err
	}

	return resp.GetCapabilities(), nil
}

// TestAdapterConnection tests the connection to a supplier's system using the specified adapter
func (c *SupplierClient) TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // Longer timeout for connection tests
	defer cancel()

	req := &supplierv1.TestAdapterConnectionRequest{
		AdapterName: adapterName,
		Config:      config,
	}

	resp, err := c.client.TestAdapterConnection(ctx, req)
	if err != nil {
		c.logger.Error("Failed to test adapter connection", zap.String("adapter", adapterName), zap.Error(err))
		return err
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("connection test failed: %s", resp.GetMessage())
	}

	return nil
}

// SyncProducts synchronizes products from a supplier using their configured adapter
func (c *SupplierClient) SyncProducts(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second) // Longer timeout for sync operations
	defer cancel()

	req := &supplierv1.SyncProductsRequest{
		SupplierId: supplierID,
		Options:    options,
	}

	resp, err := c.client.SyncProducts(ctx, req)
	if err != nil {
		c.logger.Error("Failed to sync products", zap.String("supplier_id", supplierID), zap.Error(err))
		return "", err
	}

	return resp.GetJobId(), nil
}

// SyncInventory synchronizes inventory from a supplier using their configured adapter
func (c *SupplierClient) SyncInventory(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second) // Longer timeout for sync operations
	defer cancel()

	req := &supplierv1.SyncInventoryRequest{
		SupplierId: supplierID,
		Options:    options,
	}

	resp, err := c.client.SyncInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to sync inventory", zap.String("supplier_id", supplierID), zap.Error(err))
		return "", err
	}

	return resp.GetJobId(), nil
}
