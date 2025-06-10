package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/inventory/v1"
)

// InventoryServiceImpl implements the InventoryService interface
type InventoryServiceImpl struct {
	client *grpcclient.InventoryClient
	logger *zap.Logger
}

// NewInventoryService creates a new instance of InventoryServiceImpl
func NewInventoryService(inventoryServiceAddr string, logger *zap.Logger) (InventoryService, error) {
	// Create a gRPC client
	client, err := grpcclient.NewInventoryClient(inventoryServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory client: %w", err)
	}

	return &InventoryServiceImpl{
		client: client,
		logger: logger.Named("inventory_service"),
	}, nil
}

// ListInventory lists all inventory items with pagination
func (s *InventoryServiceImpl) ListInventory(
	ctx context.Context,
	location string,
	lowStock bool,
	limit, offset int,
) (interface{}, error) {
	s.logger.Debug("ListInventory",
		zap.String("location", location),
		zap.Bool("lowStock", lowStock),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	req := &inventoryv1.ListInventoryRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resp, err := s.client.ListInventory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list inventory",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}

	return resp.GetInventories(), nil
}

// GetInventoryItemByID gets an inventory item by ID
func (s *InventoryServiceImpl) GetInventoryItemByID(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemByID",
		zap.String("id", id),
	)

	req := &inventoryv1.GetInventoryRequest{
		Id: id,
	}

	resp, err := s.client.GetInventory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get inventory item",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	return resp.GetInventory(), nil
}

// GetInventoryItemsByProductID gets inventory items by product ID
func (s *InventoryServiceImpl) GetInventoryItemsByProductID(ctx context.Context, productID string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemsByProductID",
		zap.String("productID", productID),
	)

	req := &inventoryv1.GetInventoryByProductIDRequest{
		ProductId: productID,
	}

	resp, err := s.client.GetInventoryByProductID(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get inventory by product ID",
			zap.String("productID", productID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory by product ID: %w", err)
	}

	return resp.GetInventory(), nil
}

// GetInventoryItemBySKU gets an inventory item by SKU
func (s *InventoryServiceImpl) GetInventoryItemBySKU(ctx context.Context, sku string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemBySKU",
		zap.String("sku", sku),
	)

	req := &inventoryv1.GetInventoryBySKURequest{
		Sku: sku,
	}

	resp, err := s.client.GetInventoryBySKU(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get inventory by SKU",
			zap.String("sku", sku),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory by SKU: %w", err)
	}

	return resp.GetInventory(), nil
}

// CreateInventoryItem creates a new inventory item
func (s *InventoryServiceImpl) CreateInventoryItem(
	ctx context.Context,
	productID, sku string,
	quantity int32,
	location string,
	reorderAt, reorderQty int32,
	cost float64,
) (interface{}, error) {
	s.logger.Debug("CreateInventoryItem",
		zap.String("productID", productID),
		zap.String("sku", sku),
		zap.Int32("quantity", quantity),
		zap.String("location", location),
	)

	req := &inventoryv1.CreateInventoryRequest{
		ProductId: productID,
		Sku:       sku,
		Quantity:  quantity,
		Location:  location,
	}

	resp, err := s.client.CreateInventory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create inventory item",
			zap.String("productID", productID),
			zap.String("sku", sku),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create inventory item: %w", err)
	}

	return resp.GetInventory(), nil
}

// UpdateInventoryItem updates an existing inventory item
func (s *InventoryServiceImpl) UpdateInventoryItem(
	ctx context.Context,
	id, productID, sku string,
	quantity int32,
	location string,
	reorderAt, reorderQty int32,
	cost float64,
) error {
	s.logger.Debug("UpdateInventoryItem",
		zap.String("id", id),
		zap.String("productID", productID),
		zap.String("sku", sku),
		zap.Int32("quantity", quantity),
		zap.String("location", location),
	)

	// First, get the current inventory item
	getReq := &inventoryv1.GetInventoryRequest{
		Id: id,
	}

	getResp, err := s.client.GetInventory(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get inventory item for update",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get inventory item for update: %w", err)
	}

	// Update the fields
	inventory := getResp.GetInventory()
	inventory.ProductId = productID
	inventory.Quantity = quantity
	inventory.Sku = sku
	inventory.Location = location

	req := &inventoryv1.UpdateInventoryRequest{
		Inventory: inventory,
	}

	_, err = s.client.UpdateInventory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update inventory item",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update inventory item: %w", err)
	}

	return nil
}

// DeleteInventoryItem deletes an inventory item
func (s *InventoryServiceImpl) DeleteInventoryItem(ctx context.Context, id string) error {
	s.logger.Debug("DeleteInventoryItem",
		zap.String("id", id),
	)

	req := &inventoryv1.DeleteInventoryRequest{
		Id: id,
	}

	_, err := s.client.DeleteInventory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete inventory item",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete inventory item: %w", err)
	}

	return nil
}

// AddStock adds stock to an inventory item
func (s *InventoryServiceImpl) AddStock(
	ctx context.Context,
	id string,
	quantity int32,
	reason, reference string,
) (interface{}, error) {
	s.logger.Debug("AddStock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)

	req := &inventoryv1.AddStockRequest{
		Id:       id,
		Quantity: quantity,
	}

	_, err := s.client.AddStock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to add stock",
			zap.String("id", id),
			zap.Int32("quantity", quantity),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}

	// Get the updated inventory item to return
	getReq := &inventoryv1.GetInventoryRequest{
		Id: id,
	}

	getResp, err := s.client.GetInventory(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get updated inventory item",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get updated inventory item: %w", err)
	}

	return getResp.GetInventory(), nil
}

// RemoveStock removes stock from an inventory item
func (s *InventoryServiceImpl) RemoveStock(
	ctx context.Context,
	id string,
	quantity int32,
	reason, reference string,
) (interface{}, error) {
	s.logger.Debug("RemoveStock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)

	req := &inventoryv1.RemoveStockRequest{
		Id:       id,
		Quantity: quantity,
	}

	_, err := s.client.RemoveStock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to remove stock",
			zap.String("id", id),
			zap.Int32("quantity", quantity),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to remove stock: %w", err)
	}

	// Get the updated inventory item to return
	getReq := &inventoryv1.GetInventoryRequest{
		Id: id,
	}

	getResp, err := s.client.GetInventory(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get updated inventory item",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get updated inventory item: %w", err)
	}

	return getResp.GetInventory(), nil
}
