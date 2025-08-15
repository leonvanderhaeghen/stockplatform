package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	inventoryclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/inventory"
)

// InventoryServiceImpl implements the InventoryService interface
type InventoryServiceImpl struct {
	client *inventoryclient.Client
	logger *zap.Logger
}

// NewInventoryService creates a new instance of InventoryServiceImpl
func NewInventoryService(inventoryServiceAddr string, logger *zap.Logger) (InventoryService, error) {
	// Create a gRPC client
	// Note: NewInventoryClient doesn't take a logger parameter
	invCfg := inventoryclient.Config{Address: inventoryServiceAddr}
	client, err := inventoryclient.New(invCfg, logger)
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

	if lowStock {
		// Use GetLowStockItems if lowStock is requested
		resp, err := s.client.GetLowStockItems(ctx, location, 10, limit, offset) // threshold = 10
		if err != nil {
			s.logger.Error("Failed to get low stock inventory", zap.Error(err))
			return nil, fmt.Errorf("failed to get low stock inventory: %w", err)
		}
		return resp, nil
	}

	resp, err := s.client.ListInventory(ctx, int32(limit), int32(offset))
	if err != nil {
		s.logger.Error("Failed to list inventory", zap.Error(err))
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}

	return resp, nil
}

// GetInventoryItemByID gets an inventory item by ID
func (s *InventoryServiceImpl) GetInventoryItemByID(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemByID",
		zap.String("id", id),
	)

	resp, err := s.client.GetInventory(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get inventory item",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory item: %w", err)
	}

	return resp, nil
}

// GetInventoryItemsByProductID gets inventory items by product ID
func (s *InventoryServiceImpl) GetInventoryItemsByProductID(ctx context.Context, productID string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemsByProductID",
		zap.String("productID", productID),
	)

	// Use empty locationID since interface doesn't provide it
	resp, err := s.client.GetInventoryByProductID(ctx, productID, "")
	if err != nil {
		s.logger.Error("Failed to get inventory by product ID",
			zap.String("productID", productID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory by product ID: %w", err)
	}

	return resp, nil
}

// GetInventoryItemBySKU gets an inventory item by SKU
func (s *InventoryServiceImpl) GetInventoryItemBySKU(ctx context.Context, sku string) (interface{}, error) {
	s.logger.Debug("GetInventoryItemBySKU",
		zap.String("sku", sku),
	)

	resp, err := s.client.GetInventoryBySKU(ctx, sku)
	if err != nil {
		s.logger.Error("Failed to get inventory by SKU",
			zap.String("sku", sku),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory by SKU: %w", err)
	}

	return resp, nil
}

// CreateInventoryItem creates a new inventory item
func (s *InventoryServiceImpl) CreateInventoryItem(ctx context.Context, productID, sku string, quantity int32, location string, reorderAt, reorderQty int32, cost float64) (interface{}, error) {
	s.logger.Debug("CreateInventoryItem",
		zap.String("productID", productID),
		zap.String("sku", sku),
		zap.Int32("quantity", quantity),
		zap.String("location", location),
		zap.Int32("reorderAt", reorderAt),
		zap.Int32("reorderQty", reorderQty),
		zap.Float64("cost", cost),
	)

	// Create inventory item with basic fields (client doesn't support all fields yet)
	resp, err := s.client.CreateInventory(ctx, productID, sku, location, quantity)
	if err != nil {
		s.logger.Error("Failed to create inventory item",
			zap.String("productID", productID),
			zap.String("sku", sku),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create inventory item: %w", err)
	}

	// Set reorder fields and cost in the returned inventory item if provided
	if reorderAt > 0 {
		resp.ReorderAt = reorderAt
	}
	if reorderQty > 0 {
		resp.ReorderQty = reorderQty
	}
	if cost > 0 {
		resp.Cost = cost
	}
	
	return resp, nil
}

// UpdateInventoryItem updates an existing inventory item
func (s *InventoryServiceImpl) UpdateInventoryItem(ctx context.Context, id, productID, sku string, quantity int32, location string, reorderAt, reorderQty int32, cost float64) error {
	s.logger.Debug("UpdateInventoryItem",
		zap.String("id", id),
		zap.String("productID", productID),
		zap.String("sku", sku),
		zap.String("location", location),
		zap.Int32("reorderAt", reorderAt),
		zap.Int32("reorderQty", reorderQty),
		zap.Float64("cost", cost),
	)

	item := &models.InventoryItem{
		ID:        id,
		ProductID: productID,
		SKU:       sku,
		LocationID: location,
		Quantity:  quantity,
		ReorderAt: reorderAt,
		ReorderQty: reorderQty,
		Cost:      cost,
	}
	_, err := s.client.UpdateInventory(ctx, item)
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

	err := s.client.DeleteInventory(ctx, id)
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
func (s *InventoryServiceImpl) AddStock(ctx context.Context, id string, quantity int32, reason, performedBy string) (interface{}, error) {
	s.logger.Debug("AddStock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)

	resp, err := s.client.AddStock(ctx, id, quantity, reason, performedBy)
	if err != nil {
		s.logger.Error("Failed to add stock",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}

	return resp, nil
}

// RemoveStock removes stock from an inventory item
func (s *InventoryServiceImpl) RemoveStock(ctx context.Context, id string, quantity int32, reason, performedBy string) (interface{}, error) {
	s.logger.Debug("RemoveStock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)

	resp, err := s.client.RemoveStock(ctx, id, quantity, reason, performedBy)
	if err != nil {
		s.logger.Error("Failed to remove stock",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to remove stock: %w", err)
	}

	return resp, nil
}

// Note: POS inventory operations are now handled via standard inventory endpoints:
// - POS inventory check: GetInventoryItemBySKU with availability parameters
// - POS reservations: Standard reservation methods with source parameter
// - POS deductions: RemoveStock with source parameter

// CompletePickup marks a pickup as complete
func (s *InventoryServiceImpl) CompletePickup(
	ctx context.Context,
	reservationID string,
	staffID string,
	notes string,
) (interface{}, error) {
	s.logger.Debug("CompletePickup",
		zap.String("reservationID", reservationID),
		zap.String("staffID", staffID),
	)

	err := s.client.CompletePickup(ctx, reservationID, staffID, notes)
	if err != nil {
		s.logger.Error("Failed to complete pickup",
			zap.String("reservationID", reservationID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to complete pickup: %w", err)
	}

	return reservationID, nil
}

// Note: POS inventory deductions are now handled via RemoveStock with source parameter
// All POS functionality has been consolidated into standard inventory endpoints

// GetInventoryReservations gets inventory reservations with optional filters
// Note: The inventory service doesn't currently have a method to list reservations,
// so this returns an empty list as a placeholder
func (s *InventoryServiceImpl) GetInventoryReservations(
	ctx context.Context,
	orderId, productId, status string,
	limit, offset int,
) (interface{}, error) {
	s.logger.Debug("GetInventoryReservations",
		zap.String("orderId", orderId),
		zap.String("productId", productId),
		zap.String("status", status),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	// NOTE: The inventory service doesn't have a ListReservations method yet.
	// This is a placeholder implementation that returns an empty list.
	// FEATURE ENHANCEMENT: Add ListReservations method to inventory service for full reservation tracking
	s.logger.Info("GetInventoryReservations called - returning empty list (method not implemented in inventory service)")

	// Return empty reservations list as placeholder
	return []interface{}{}, nil
}

// CreateInventoryReservation creates a new inventory reservation (supports POS source tracking)
func (s *InventoryServiceImpl) CreateInventoryReservation(
	ctx context.Context,
	productID string,
	quantity int32,
	orderID string,
) (interface{}, error) {
	s.logger.Debug("CreateInventoryReservation",
		zap.String("productId", productID),
		zap.Int32("quantity", quantity),
		zap.String("orderId", orderID),
	)

	// First, get the inventory item by product ID to get the inventory ID needed for reservation
	inventoryItem, err := s.client.GetInventoryByProductID(ctx, productID, "")
	if err != nil {
		s.logger.Error("Failed to find inventory item for reservation",
			zap.String("productId", productID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to find inventory item for product %s: %w", productID, err)
	}

	// Reserve stock using the inventory service
	success, err := s.client.ReserveStock(ctx, inventoryItem.ID, quantity)
	if err != nil {
		s.logger.Error("Failed to reserve stock",
			zap.String("inventoryId", inventoryItem.ID),
			zap.String("productId", productID),
			zap.Int32("quantity", quantity),
			zap.String("orderId", orderID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	if !success {
		s.logger.Warn("Stock reservation was not successful",
			zap.String("inventoryId", inventoryItem.ID),
			zap.String("productId", productID),
			zap.Int32("quantity", quantity),
		)
		return nil, fmt.Errorf("insufficient stock available for reservation")
	}

	// Generate reservation ID and create response
	reservationID := fmt.Sprintf("res_%s_%s_%d_%d", productID, orderID, quantity, time.Now().Unix())
	resp := map[string]interface{}{
		"reservationId": reservationID,
		"inventoryId": inventoryItem.ID,
		"productId": productID,
		"sku": inventoryItem.SKU,
		"quantity": quantity,
		"orderId": orderID,
		"status": "RESERVED",
		"createdAt": time.Now().Format(time.RFC3339),
		"locationId": inventoryItem.LocationID,
		"reservedStock": success,
	}

	s.logger.Info("Inventory reservation created successfully",
		zap.String("productId", productID),
		zap.Int32("quantity", quantity),
		zap.String("orderId", orderID),
	)

	return resp, nil
}

// GetLowStockItems gets inventory items that are low in stock with threshold and location filtering
func (s *InventoryServiceImpl) GetLowStockItems(
	ctx context.Context,
	location string,
	threshold, limit, offset int,
) (interface{}, error) {
	s.logger.Debug("GetLowStockItems",
		zap.String("location", location),
		zap.Int("threshold", threshold),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	resp, err := s.client.GetLowStockItems(ctx, location, threshold, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get low stock items",
			zap.Int("threshold", threshold),
			zap.String("location", location),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get low stock items: %w", err)
	}

	return resp, nil
}


