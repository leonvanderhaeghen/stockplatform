package supplier

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
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
func (c *Client) CreateSupplier(ctx context.Context, name, contactPerson, email, phone, address, city, state, country, postalCode, taxID, website, currency, paymentTerms string, leadTimeDays int32, metadata map[string]string) (*models.CreateSupplierResponse, error) {
	c.logger.Debug("Creating supplier", zap.String("name", name))
	
	req := &supplierv1.CreateSupplierRequest{
		Name:          name,
		ContactPerson: contactPerson,
		Email:         email,
		Phone:         phone,
		Address:       address,
		City:          city,
		State:         state,
		Country:       country,
		PostalCode:    postalCode,
		TaxId:         taxID,
		Website:       website,
		Currency:      currency,
		LeadTimeDays:  leadTimeDays,
		PaymentTerms:  paymentTerms,
		Metadata:      metadata,
	}
	
	resp, err := c.client.CreateSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create supplier", zap.Error(err))
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}
	
	c.logger.Debug("Supplier created successfully", zap.String("id", resp.Supplier.Id))
	return c.convertToCreateSupplierResponse(resp), nil
}

// GetSupplier retrieves a supplier by ID
func (c *Client) GetSupplier(ctx context.Context, id string) (*models.Supplier, error) {
	c.logger.Debug("Getting supplier", zap.String("id", id))
	
	req := &supplierv1.GetSupplierRequest{
		Id: id,
	}
	
	resp, err := c.client.GetSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get supplier", zap.Error(err))
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}
	
	return c.convertToSupplier(resp.Supplier), nil
}

// UpdateSupplier updates an existing supplier
func (c *Client) UpdateSupplier(ctx context.Context, id, name, contactPerson, email, phone, address, city, state, country, postalCode, taxID, website, currency, paymentTerms string, leadTimeDays int32, metadata map[string]string) (*models.UpdateSupplierResponse, error) {
	c.logger.Debug("Updating supplier", zap.String("id", id))
	
	req := &supplierv1.UpdateSupplierRequest{
		Id:            id,
		Name:          name,
		ContactPerson: contactPerson,
		Email:         email,
		Phone:         phone,
		Address:       address,
		City:          city,
		State:         state,
		Country:       country,
		PostalCode:    postalCode,
		TaxId:         taxID,
		Website:       website,
		Currency:      currency,
		LeadTimeDays:  leadTimeDays,
		PaymentTerms:  paymentTerms,
		Metadata:      metadata,
	}
	
	resp, err := c.client.UpdateSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update supplier", zap.Error(err))
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}
	
	return c.convertToUpdateSupplierResponse(resp), nil
}

// DeleteSupplier deletes a supplier by ID
func (c *Client) DeleteSupplier(ctx context.Context, id string) error {
	c.logger.Debug("Deleting supplier", zap.String("id", id))
	
	req := &supplierv1.DeleteSupplierRequest{
		Id: id,
	}
	
	_, err := c.client.DeleteSupplier(ctx, req)
	if err != nil {
		c.logger.Error("Failed to delete supplier", zap.Error(err))
		return fmt.Errorf("failed to delete supplier: %w", err)
	}
	
	c.logger.Debug("Supplier deleted successfully", zap.String("id", id))
	return nil
}

// ListSuppliers lists suppliers with pagination
func (c *Client) ListSuppliers(ctx context.Context, pageSize int32, pageToken string) (*models.ListSuppliersResponse, error) {
	c.logger.Debug("Listing suppliers", zap.Int32("page_size", pageSize))
	
	req := &supplierv1.ListSuppliersRequest{
		PageSize: pageSize,
		Search:   pageToken, // Using search parameter instead of pageToken
	}
	
	resp, err := c.client.ListSuppliers(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list suppliers", zap.Error(err))
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}
	
	return c.convertToListSuppliersResponse(resp), nil
}
