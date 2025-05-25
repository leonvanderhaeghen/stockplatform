package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryv1 "stockplatform/pkg/gen/inventory/v1"
)

// InventoryServiceImpl implements the InventoryService interface
type InventoryServiceImpl struct {
	client inventoryv1.InventoryServiceClient
	logger *zap.Logger
}

// NewInventoryService creates a new instance of InventoryServiceImpl
func NewInventoryService(inventoryServiceAddr string, logger *zap.Logger) (InventoryService, error) {
	// Create a gRPC connection to the inventory service
	conn, err := grpc.Dial(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	// Create a client
	client := inventoryv1.NewInventoryServiceClient(conn)

	return &InventoryServiceImpl{
		client: client,
		logger: logger.Named("inventory_service"),
	}, nil
}

// ListInventory lists inventory items with filtering options
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
		Location: location,
		LowStock: lowStock,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	resp, err := s.client.ListInventory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list inventory",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}

	return resp, nil
}

// GetInventoryItemByID gets an inventory item by ID
func (s *InventoryServiceImpl) GetInventoryItemByID(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemByID",
		zap.String("id", id),
	)

	req := &inventoryv1.GetInventoryItemRequest{
		Id: id,
	}

	resp, err := s.client.GetInventoryItem(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get inventory item",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	return resp.Item, nil
}

// GetInventoryItemsByProductID gets inventory items by product ID
func (s *InventoryServiceImpl) GetInventoryItemsByProductID(ctx context.Context, productID string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemsByProductID",
		zap.String("productID", productID),
	)

	req := &inventoryv1.GetInventoryByProductRequest{
		ProductId: productID,
	}

	resp, err := s.client.GetInventoryByProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get inventory items by product ID",
			zap.String("productID", productID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory items by product ID: %w", err)
	}

	return resp.Items, nil
}

// GetInventoryItemBySKU gets an inventory item by SKU
func (s *InventoryServiceImpl) GetInventoryItemBySKU(ctx context.Context, sku string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemBySKU",
		zap.String("sku", sku),
	)

	req := &inventoryv1.GetInventoryBySkuRequest{
		Sku: sku,
	}

	resp, err := s.client.GetInventoryBySku(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get inventory item by SKU",
			zap.String("sku", sku),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory item by SKU: %w", err)
	}

	return resp.Item, nil
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

	req := &inventoryv1.CreateInventoryItemRequest{
		Item: &inventoryv1.InventoryItem{
			ProductId:  productID,
			Sku:        sku,
			Quantity:   quantity,
			Location:   location,
			ReorderAt:  reorderAt,
			ReorderQty: reorderQty,
			Cost:       cost,
		},
	}

	resp, err := s.client.CreateInventoryItem(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create inventory item",
			zap.String("productID", productID),
			zap.String("sku", sku),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create inventory item: %w", err)
	}

	return resp.Item, nil
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
	)

	req := &inventoryv1.UpdateInventoryItemRequest{
		Item: &inventoryv1.InventoryItem{
			Id:         id,
			ProductId:  productID,
			Sku:        sku,
			Quantity:   quantity,
			Location:   location,
			ReorderAt:  reorderAt,
			ReorderQty: reorderQty,
			Cost:       cost,
		},
	}

	_, err := s.client.UpdateInventoryItem(ctx, req)
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

	req := &inventoryv1.DeleteInventoryItemRequest{
		Id: id,
	}

	_, err := s.client.DeleteInventoryItem(ctx, req)
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
		zap.String("reason", reason),
		zap.String("reference", reference),
	)

	req := &inventoryv1.AddStockRequest{
		Id:        id,
		Quantity:  quantity,
		Reason:    reason,
		Reference: reference,
	}

	resp, err := s.client.AddStock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to add stock",
			zap.String("id", id),
			zap.Int32("quantity", quantity),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}

	return resp.Item, nil
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
		zap.String("reason", reason),
		zap.String("reference", reference),
	)

	req := &inventoryv1.RemoveStockRequest{
		Id:        id,
		Quantity:  quantity,
		Reason:    reason,
		Reference: reference,
	}

	resp, err := s.client.RemoveStock(ctx, req)
	if err != nil {
		s.logger.Error("Failed to remove stock",
			zap.String("id", id),
			zap.Int32("quantity", quantity),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to remove stock: %w", err)
	}

	return resp.Item, nil
}
