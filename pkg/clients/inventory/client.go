package inventory

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
)

// Client provides a high-level interface for interacting with the Inventory service
type Client struct {
	conn   *grpc.ClientConn
	client inventoryv1.InventoryServiceClient
	logger *zap.Logger
}

// Config holds configuration for the Inventory client
type Config struct {
	Address string
	Timeout time.Duration
}

// New creates a new Inventory service client
func New(config Config, logger *zap.Logger) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	conn, err := grpc.Dial(config.Address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	client := inventoryv1.NewInventoryServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the connection to the Inventory service
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateInventory creates a new inventory item
func (c *Client) CreateInventory(ctx context.Context, req *inventoryv1.CreateInventoryRequest) (*inventoryv1.CreateInventoryResponse, error) {
	c.logger.Debug("Creating inventory", zap.String("product_id", req.ProductId))
	
	resp, err := c.client.CreateInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}
	
	c.logger.Debug("Inventory created successfully", zap.String("id", resp.Inventory.Id))
	return resp, nil
}

// GetInventory retrieves an inventory item by ID
func (c *Client) GetInventory(ctx context.Context, req *inventoryv1.GetInventoryRequest) (*inventoryv1.GetInventoryResponse, error) {
	c.logger.Debug("Getting inventory", zap.String("id", req.Id))
	
	resp, err := c.client.GetInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	
	return resp, nil
}

// GetInventoryByProductID retrieves an inventory item by product ID
func (c *Client) GetInventoryByProductID(ctx context.Context, req *inventoryv1.GetInventoryByProductIDRequest) (*inventoryv1.GetInventoryResponse, error) {
	c.logger.Debug("Getting inventory by product ID", zap.String("product_id", req.ProductId))
	
	resp, err := c.client.GetInventoryByProductID(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get inventory by product ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get inventory by product ID: %w", err)
	}
	
	return resp, nil
}

// UpdateInventory updates an existing inventory item
func (c *Client) UpdateInventory(ctx context.Context, req *inventoryv1.UpdateInventoryRequest) (*inventoryv1.UpdateInventoryResponse, error) {
	c.logger.Debug("Updating inventory", zap.String("id", req.Inventory.Id))
	
	resp, err := c.client.UpdateInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}
	
	return resp, nil
}

// ListInventory lists all inventory items with pagination
func (c *Client) ListInventory(ctx context.Context, req *inventoryv1.ListInventoryRequest) (*inventoryv1.ListInventoryResponse, error) {
	c.logger.Debug("Listing inventory")
	
	resp, err := c.client.ListInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	
	return resp, nil
}

// AddStock adds stock to an inventory item
func (c *Client) AddStock(ctx context.Context, req *inventoryv1.AddStockRequest) (*inventoryv1.AddStockResponse, error) {
	c.logger.Debug("Adding stock", zap.String("id", req.Id), zap.Int32("quantity", req.Quantity))
	
	resp, err := c.client.AddStock(ctx, req)
	if err != nil {
		c.logger.Error("Failed to add stock", zap.Error(err))
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}
	
	return resp, nil
}

// RemoveStock removes stock from an inventory item
func (c *Client) RemoveStock(ctx context.Context, req *inventoryv1.RemoveStockRequest) (*inventoryv1.RemoveStockResponse, error) {
	c.logger.Debug("Removing stock", zap.String("id", req.Id), zap.Int32("quantity", req.Quantity))
	
	resp, err := c.client.RemoveStock(ctx, req)
	if err != nil {
		c.logger.Error("Failed to remove stock", zap.Error(err))
		return nil, fmt.Errorf("failed to remove stock: %w", err)
	}
	
	return resp, nil
}

// ReserveStock reserves stock for an order
func (c *Client) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	c.logger.Debug("Reserving stock", zap.String("id", req.Id), zap.Int32("quantity", req.Quantity))
	
	resp, err := c.client.ReserveStock(ctx, req)
	if err != nil {
		c.logger.Error("Failed to reserve stock", zap.Error(err))
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}
	
	return resp, nil
}

// CheckAvailability checks if a product is available in inventory
func (c *Client) CheckAvailability(ctx context.Context, req *inventoryv1.CheckAvailabilityRequest) (*inventoryv1.CheckAvailabilityResponse, error) {
	c.logger.Debug("Checking availability", zap.String("product_id", req.ProductId))
	
	resp, err := c.client.CheckAvailability(ctx, req)
	if err != nil {
		c.logger.Error("Failed to check availability", zap.Error(err))
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}
	
	return resp, nil
}

// AdjustInventory adjusts inventory levels
func (c *Client) AdjustInventory(ctx context.Context, req *inventoryv1.AdjustInventoryRequest) (*inventoryv1.AdjustInventoryResponse, error) {
	c.logger.Debug("Adjusting inventory", zap.String("id", req.Id), zap.Int32("quantity", req.Quantity))
	
	resp, err := c.client.AdjustInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to adjust inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to adjust inventory: %w", err)
	}
	
	return resp, nil
}
