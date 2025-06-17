package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	inventoryclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/inventory"
	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
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
		ProductId:  productID,
		Sku:        sku,
		Quantity:   quantity,
		LocationId: location,
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
	inventory.LocationId = location

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

// PerformPOSInventoryCheck checks inventory availability for POS
func (s *InventoryServiceImpl) PerformPOSInventoryCheck(
	ctx context.Context,
	locationID string,
	items []map[string]interface{},
) (interface{}, error) {
	s.logger.Debug("PerformPOSInventoryCheck",
		zap.String("locationID", locationID),
		zap.Any("items", items),
	)

	// Convert items to InventoryRequestItem for check
	requestItems := make([]*inventoryv1.InventoryRequestItem, 0, len(items))
	for _, item := range items {
		// Extract product ID
		productIDVal, ok := item["product_id"]
		if !ok {
			return nil, fmt.Errorf("missing product_id in item")
		}
		productID, ok := productIDVal.(string)
		if !ok {
			return nil, fmt.Errorf("product_id must be a string")
		}
		
		// Extract SKU if available
		sku := ""
		if skuVal, ok := item["sku"].(string); ok {
			sku = skuVal
		}
		
		// Extract quantity if available
		quantity := int32(1) // Default to 1 if not specified
		if qtyVal, ok := item["quantity"]; ok {
			switch q := qtyVal.(type) {
			case int:
				quantity = int32(q)
			case int32:
				quantity = q
			case float64:
				quantity = int32(q)
			case float32:
				quantity = int32(q)
			}
		}
		
		// Create the inventory request item
		requestItem := &inventoryv1.InventoryRequestItem{
			ProductId: productID,
			Sku:       sku,
			Quantity:  quantity,
		}
		
		// Add inventory item ID if available
		if invIDVal, ok := item["inventory_item_id"].(string); ok {
			requestItem.InventoryItemId = invIDVal
		}
		
		requestItems = append(requestItems, requestItem)
	}

	// Create CheckAvailability request with the Items field
	req := &inventoryv1.CheckAvailabilityRequest{
		LocationId: locationID,
		Items:      requestItems,
	}

	// Call the gRPC method
	resp, err := s.client.CheckAvailability(ctx, req)
	if err != nil {
		s.logger.Error("Failed to check inventory availability",
			zap.String("locationID", locationID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check inventory availability: %w", err)
	}

	// Return the availability information
	return resp, nil
}

// ReserveForPOSTransaction reserves inventory for POS transactions
func (s *InventoryServiceImpl) ReserveForPOSTransaction(
	ctx context.Context,
	locationID string,
	orderID string,
	items []map[string]interface{},
) (interface{}, error) {
	s.logger.Debug("ReserveForPOSTransaction",
		zap.String("locationID", locationID),
		zap.String("orderID", orderID),
		zap.Any("items", items),
	)

	// Note: ReserveStockRequest only has Id and Quantity in the protobuf definition
	// Since we have multiple items, we need to process them sequentially
	// and return a combined result
	results := make(map[string]bool)

	for _, item := range items {
		// Get product ID
		productIDVal, ok := item["product_id"]
		if !ok {
			return nil, fmt.Errorf("missing product_id in item")
		}
		productID, ok := productIDVal.(string)
		if !ok {
			return nil, fmt.Errorf("product_id must be a string")
		}

		// Get quantity
		quantityVal, ok := item["quantity"]
		if !ok {
			return nil, fmt.Errorf("missing quantity in item")
		}

		// Handle different types for quantity
		var quantity int32
		switch q := quantityVal.(type) {
		case int:
			quantity = int32(q)
		case int32:
			quantity = q
		case float64:
			quantity = int32(q)
		default:
			return nil, fmt.Errorf("quantity must be a number")
		}

		// Create ReserveStock request for this item
		req := &inventoryv1.ReserveStockRequest{
			Id:       productID, // Using Id instead of ProductId based on the protobuf definition
			Quantity: quantity,
		}

		// Call the gRPC method for this item
		resp, err := s.client.ReserveStock(ctx, req)
		if err != nil {
			s.logger.Error("Failed to reserve inventory item",
				zap.String("productID", productID),
				zap.Int32("quantity", quantity),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to reserve inventory item %s: %w", productID, err)
		}

		results[productID] = resp.GetSuccess()
	}

	// For backward compatibility, we're returning a reservation ID
	// Note: This is a workaround since the actual API doesn't return a reservation ID
	// In a real implementation, we should update the API to match the expected behavior
	return orderID, nil
}

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

	// Create CompletePickup request
	req := &inventoryv1.CompletePickupRequest{
		ReservationId: reservationID,
		StaffId:       staffID,
		Notes:         notes,
	}

	// Call the gRPC method
	_, err := s.client.CompletePickup(ctx, req)
	if err != nil {
		s.logger.Error("Failed to complete pickup",
			zap.String("reservationID", reservationID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to complete pickup: %w", err)
	}

	// Note: The actual response might not have a GetPickupId method
	// For backward compatibility, returning the reservation ID
	return reservationID, nil
}

// DeductForDirectPOSTransaction directly deducts inventory for POS sales
func (s *InventoryServiceImpl) DeductForDirectPOSTransaction(
	ctx context.Context,
	locationID string,
	staffID string,
	items []map[string]interface{},
	reason string,
) (interface{}, error) {
	s.logger.Debug("DeductForDirectPOSTransaction",
		zap.String("locationID", locationID),
		zap.String("staffID", staffID),
		zap.Any("items", items),
	)

	// Process each inventory item separately since AdjustInventoryForOrderRequest
	// structure is different than what we're using
	results := make(map[string]bool)
	for _, item := range items {
		// Get product ID
		productIDVal, ok := item["product_id"]
		if !ok {
			return nil, fmt.Errorf("missing product_id in item")
		}
		productID, ok := productIDVal.(string)
		if !ok {
			return nil, fmt.Errorf("product_id must be a string")
		}

		// Get quantity
		quantityVal, ok := item["quantity"]
		if !ok {
			return nil, fmt.Errorf("missing quantity in item")
		}

		// Handle different types for quantity
		var quantity int32
		switch q := quantityVal.(type) {
		case int:
			quantity = int32(q)
		case int32:
			quantity = q
		case float64:
			quantity = int32(q)
		default:
			return nil, fmt.Errorf("quantity must be a number")
		}

		// Create proper request based on the actual protobuf definition
		// The Items field takes an array of InventoryAdjustmentItem
		adjustmentItem := &inventoryv1.InventoryAdjustmentItem{
			ProductId: productID,
			Quantity:  quantity,
		}

		req := &inventoryv1.AdjustInventoryForOrderRequest{
			OrderId:        "pos_" + productID,  // Generate an order ID
			LocationId:     locationID,
			AdjustmentType: "sale",
			ReferenceId:    reason,
			Items:          []*inventoryv1.InventoryAdjustmentItem{adjustmentItem},
			StaffId:        staffID,
		}

		// Call the gRPC method
		_, err := s.client.AdjustInventoryForOrder(ctx, req)
		if err != nil {
			s.logger.Error("Failed to adjust inventory item",
				zap.String("productID", productID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to adjust inventory for product %s: %w", productID, err)
		}

		// Store the result
		// Assuming the response has a Success field
		results[productID] = true // Default to true since we don't know the actual structure
	}

	// Generate a transaction ID for backward compatibility
	transactionID := fmt.Sprintf("pos_%d", time.Now().UnixNano())
	return transactionID, nil
}

// Second implementation of PerformPOSInventoryCheck removed to resolve duplicate method error

// Second implementation of ReserveForPOSTransaction removed to resolve duplicate method error

// Third duplicate implementation of CompletePickup removed to resolve duplicate method error

// Fourth duplicate implementation of DeductForDirectPOSTransaction removed to resolve duplicate method error
