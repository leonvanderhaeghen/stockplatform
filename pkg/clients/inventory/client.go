package inventory

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
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
func (c *Client) CreateInventory(ctx context.Context, productID, sku, locationID string, quantity int32) (*models.InventoryItem, error) {
	c.logger.Debug("Creating inventory", zap.String("product_id", productID))

	req := &inventoryv1.CreateInventoryRequest{
		ProductId:  productID,
		Sku:        sku,
		Quantity:   quantity,
		LocationId: locationID,
	}

	resp, err := c.client.CreateInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	c.logger.Debug("Inventory created successfully", zap.String("id", resp.Inventory.Id))
	return c.convertToInventoryItem(resp.Inventory), nil
}

// GetInventory retrieves an inventory item by ID
func (c *Client) GetInventory(ctx context.Context, id string) (*models.InventoryItem, error) {
	c.logger.Debug("Getting inventory", zap.String("id", id))

	req := &inventoryv1.GetInventoryRequest{Id: id}

	resp, err := c.client.GetInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	return c.convertToInventoryItem(resp.Inventory), nil
}

// GetInventoryByProductID retrieves an inventory item by product ID
func (c *Client) GetInventoryByProductID(ctx context.Context, productID, locationID string) (*models.InventoryItem, error) {
	c.logger.Debug("Getting inventory by product ID", zap.String("product_id", productID))

	req := &inventoryv1.GetInventoryByProductIDRequest{ProductId: productID, LocationId: locationID}

	resp, err := c.client.GetInventoryByProductID(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get inventory by product ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get inventory by product ID: %w", err)
	}

	return c.convertToInventoryItem(resp.Inventory), nil
}

// UpdateInventory updates an existing inventory item
func (c *Client) UpdateInventory(ctx context.Context, item *models.InventoryItem) (bool, error) {
	c.logger.Debug("Updating inventory", zap.String("id", item.ID))

	req := &inventoryv1.UpdateInventoryRequest{
		Inventory: c.convertFromInventoryItem(item),
	}

	resp, err := c.client.UpdateInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update inventory", zap.Error(err))
		return false, fmt.Errorf("failed to update inventory: %w", err)
	}

	return resp.Success, nil
}

// ListInventory lists all inventory items with pagination
func (c *Client) ListInventory(ctx context.Context, limit, offset int32) ([]*models.InventoryItem, error) {
	c.logger.Debug("Listing inventory")

	req := &inventoryv1.ListInventoryRequest{
		Limit:  limit,
		Offset: offset,
	}

	resp, err := c.client.ListInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}

	items := make([]*models.InventoryItem, len(resp.Inventories))
	for i, item := range resp.Inventories {
		items[i] = c.convertToInventoryItem(item)
	}
	return items, nil
}

// AddStock adds stock to an inventory item
func (c *Client) AddStock(ctx context.Context, id string, quantity int32, reason, performedBy string) (bool, error) {
	c.logger.Debug("Adding stock", zap.String("id", id), zap.Int32("quantity", quantity))

	req := &inventoryv1.AddStockRequest{
		Id:          id,
		Quantity:    quantity,
		Reason:      reason,
		PerformedBy: performedBy,
	}

	resp, err := c.client.AddStock(ctx, req)
	if err != nil {
		c.logger.Error("Failed to add stock", zap.Error(err))
		return false, fmt.Errorf("failed to add stock: %w", err)
	}

	return resp.Success, nil
}

// RemoveStock removes stock from an inventory item
func (c *Client) RemoveStock(ctx context.Context, id string, quantity int32, reason, performedBy string) (bool, error) {
	c.logger.Debug("Removing stock", zap.String("id", id), zap.Int32("quantity", quantity))

	req := &inventoryv1.RemoveStockRequest{
		Id:          id,
		Quantity:    quantity,
		Reason:      reason,
		PerformedBy: performedBy,
	}

	resp, err := c.client.RemoveStock(ctx, req)
	if err != nil {
		c.logger.Error("Failed to remove stock", zap.Error(err))
		return false, fmt.Errorf("failed to remove stock: %w", err)
	}

	return resp.Success, nil
}



// ReserveStock reserves stock for an order
func (c *Client) ReserveStock(ctx context.Context, id string, quantity int32) (bool, error) {
	c.logger.Debug("Reserving stock", zap.String("id", id), zap.Int32("quantity", quantity))

	req := &inventoryv1.ReserveStockRequest{
		Id:       id,
		Quantity: quantity,
	}

	resp, err := c.client.ReserveStock(ctx, req)
	if err != nil {
		c.logger.Error("Failed to reserve stock", zap.Error(err))
		return false, fmt.Errorf("failed to reserve stock: %w", err)
	}

	return resp.Success, nil
}

// CheckAvailability checks item availability at a specific location
func (c *Client) CheckAvailability(ctx context.Context, locationID string, items []*models.InventoryRequestItem) (*models.CheckAvailabilityResponse, error) {
	c.logger.Debug("Checking availability", zap.String("location_id", locationID), zap.Int("items_count", len(items)))

	// Convert domain items to protobuf items
	protoItems := make([]*inventoryv1.InventoryRequestItem, len(items))
	for i, item := range items {
		protoItems[i] = &inventoryv1.InventoryRequestItem{
			ProductId: item.ProductID,
			Sku:       item.SKU,
			Quantity:  item.Quantity,
		}
	}

	req := &inventoryv1.CheckAvailabilityRequest{
		LocationId: locationID,
		Items:      protoItems,
	}

	resp, err := c.client.CheckAvailability(ctx, req)
	if err != nil {
		c.logger.Error("Failed to check availability", zap.Error(err))
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}

	return c.convertToCheckAvailabilityResponse(resp), nil
}

// GetInventoryBySKU retrieves an inventory item by SKU
func (c *Client) GetInventoryBySKU(ctx context.Context, sku string) (*models.InventoryItem, error) {
	c.logger.Debug("Getting inventory by SKU", zap.String("sku", sku))

	req := &inventoryv1.GetInventoryBySKURequest{Sku: sku}

	resp, err := c.client.GetInventoryBySKU(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get inventory by SKU", zap.Error(err))
		return nil, fmt.Errorf("failed to get inventory by SKU: %w", err)
	}

	return c.convertToInventoryItem(resp.Inventory), nil
}

// DeleteInventory deletes an inventory item
func (c *Client) DeleteInventory(ctx context.Context, id string) error {
	c.logger.Debug("Deleting inventory", zap.String("id", id))

	req := &inventoryv1.DeleteInventoryRequest{Id: id}

	_, err := c.client.DeleteInventory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to delete inventory", zap.Error(err))
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	return nil
}

// CompletePickup marks a pickup as complete
func (c *Client) CompletePickup(ctx context.Context, reservationID, staffID, notes string) error {
	c.logger.Debug("Completing pickup", zap.String("reservation_id", reservationID))

	req := &inventoryv1.CompletePickupRequest{
		ReservationId: reservationID,
		StaffId:       staffID,
		Notes:         notes,
	}

	_, err := c.client.CompletePickup(ctx, req)
	if err != nil {
		c.logger.Error("Failed to complete pickup", zap.Error(err))
		return fmt.Errorf("failed to complete pickup: %w", err)
	}

	return nil
}

// GetLowStockItems gets inventory items that are low in stock
// Note: This is a placeholder implementation since the protobuf service doesn't have this method yet
func (c *Client) GetLowStockItems(ctx context.Context, location string, threshold, limit, offset int) ([]*models.InventoryItem, error) {
	c.logger.Debug("Getting low stock items", zap.String("location", location), zap.Int("threshold", threshold))

	// For now, we'll use ListInventory and filter client-side
	// In the future, this should be implemented as a proper RPC method
	allItems, err := c.ListInventory(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory for low stock check: %w", err)
	}

	// Filter items that have quantity <= threshold
	var lowStockItems []*models.InventoryItem
	for _, item := range allItems {
		if item.Quantity <= int32(threshold) {
			if location == "" || item.LocationID == location {
				lowStockItems = append(lowStockItems, item)
			}
		}
	}

	return lowStockItems, nil
}
