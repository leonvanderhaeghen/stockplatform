package supplier

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
)

// Client provides a high-level interface for interacting with the Supplier service
type Client struct {
	conn   *grpc.ClientConn
	client supplierv1.SupplierServiceClient
	logger *zap.Logger
}

// Config holds configuration for the Supplier client
type Config struct {
	Address string
	Timeout time.Duration
}

// New creates a new Supplier service client
func New(config Config, logger *zap.Logger) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	conn, err := grpc.Dial(config.Address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to supplier service: %w", err)
	}

	client := supplierv1.NewSupplierServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the connection to the Supplier service
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateSupplier creates a new supplier
func (c *Client) CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.CreateSupplierResponse, error) {
	c.logger.Debug("Creating supplier", zap.String("name", req.Name))
	
	resp, err := c.client.CreateSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create supplier", zap.Error(err))
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}
	
	c.logger.Debug("Supplier created successfully", zap.String("id", resp.Supplier.Id))
	return resp, nil
}

// GetSupplier retrieves a supplier by ID
func (c *Client) GetSupplier(ctx context.Context, req *supplierv1.GetSupplierRequest) (*supplierv1.GetSupplierResponse, error) {
	c.logger.Debug("Getting supplier", zap.String("id", req.Id))
	
	resp, err := c.client.GetSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get supplier", zap.Error(err))
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}
	
	return resp, nil
}

// UpdateSupplier updates an existing supplier
func (c *Client) UpdateSupplier(ctx context.Context, req *supplierv1.UpdateSupplierRequest) (*supplierv1.UpdateSupplierResponse, error) {
	c.logger.Debug("Updating supplier", zap.String("id", req.Id))
	
	resp, err := c.client.UpdateSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update supplier", zap.Error(err))
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}
	
	return resp, nil
}

// DeleteSupplier deletes a supplier by ID
func (c *Client) DeleteSupplier(ctx context.Context, req *supplierv1.DeleteSupplierRequest) (*supplierv1.DeleteSupplierResponse, error) {
	c.logger.Debug("Deleting supplier", zap.String("id", req.Id))
	
	resp, err := c.client.DeleteSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to delete supplier", zap.Error(err))
		return nil, fmt.Errorf("failed to delete supplier: %w", err)
	}
	
	return resp, nil
}

// ListSuppliers lists suppliers with pagination
func (c *Client) ListSuppliers(ctx context.Context, req *supplierv1.ListSuppliersRequest) (*supplierv1.ListSuppliersResponse, error) {
	c.logger.Debug("Listing suppliers", zap.Int32("page_size", req.PageSize))
	
	resp, err := c.client.ListSuppliers(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list suppliers", zap.Error(err))
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}
	
	return resp, nil
}
