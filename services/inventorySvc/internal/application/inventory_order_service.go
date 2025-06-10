package application

import (
	"context"
	"fmt"

	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
	orderpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/order/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"go.uber.org/zap"
)

// InventoryOrderService coordinates between inventory and order services
type InventoryOrderService struct {
	orderClient      *grpcclient.OrderClient
	inventoryService *InventoryService
	locationService  *LocationService
	localLocationID  string // Default location ID for order fulfillment
	logger           *zap.Logger
}

// NewInventoryOrderService creates a new InventoryOrderService
func NewInventoryOrderService(
	inventoryService *InventoryService,
	locationService *LocationService,
	orderServiceAddr string,
	localLocationID string, // ID of the default location for fulfillment
	logger *zap.Logger,
) (*InventoryOrderService, error) {
	// Initialize the order client
	orderClient, err := grpcclient.NewOrderClient(orderServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create order client: %w", err)
	}

	return &InventoryOrderService{
		orderClient:      orderClient,
		inventoryService:  inventoryService,
		locationService:   locationService,
		localLocationID:   localLocationID,
		logger:           logger.Named("inventory_order_service"),
	}, nil
}

// Close closes any open connections
func (s *InventoryOrderService) Close() error {
	if s.orderClient != nil {
		return s.orderClient.Close()
	}
	return nil
}

// ProcessOrderUpdate handles inventory updates when an order status changes
func (s *InventoryOrderService) ProcessOrderUpdate(
	ctx context.Context,
	orderID string,
	newStatus string,
) error {
	// Get the order details
	orderResp, err := s.orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: orderID})
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	order := orderResp.GetOrder()
	
	// Get fulfillment location from order metadata or shipping address - fallback to local location if not specified
	fulfillmentLocationID := s.localLocationID
	// Check order metadata for location (this assumes the service might be updated to include location)
	// For now we'll use the default location but this can be extended when the order service is updated
	
	// Handle different order statuses
	switch newStatus {
	case "processing":
		// Reserve inventory when order is being processed
		for _, item := range order.GetItems() {
			// Get inventory item for this product at the fulfillment location
			inventory, err := s.inventoryService.GetInventoryItemByProductAndLocation(ctx, 
				item.GetProductId(), fulfillmentLocationID)
			if err != nil {
				return fmt.Errorf("failed to get inventory for product %s at location %s: %w",
					item.GetProductId(), fulfillmentLocationID, err)
			}
			
			// Reserve stock
			err = s.inventoryService.ReserveStock(ctx, inventory.ID, int32(item.GetQuantity()))
			if err != nil {
				return fmt.Errorf("failed to reserve stock for product %s: %w", 
					item.GetProductId(), err)
			}
		}

	case "shipped":
		// Fulfill reservation when order is shipped
		for _, item := range order.GetItems() {
			// Get inventory item for this product at the fulfillment location
			inventory, err := s.inventoryService.GetInventoryItemByProductAndLocation(ctx, 
				item.GetProductId(), fulfillmentLocationID)
			if err != nil {
				return fmt.Errorf("failed to get inventory for product %s at location %s: %w",
					item.GetProductId(), fulfillmentLocationID, err)
			}
			
			// Fulfill reservation
			err = s.inventoryService.FulfillReservation(ctx, inventory.ID, int32(item.GetQuantity()))
			if err != nil {
				return fmt.Errorf("failed to fulfill reservation for product %s: %w", 
					item.GetProductId(), err)
			}
		}

	case "cancelled":
		// Release reservation when order is cancelled
		for _, item := range order.GetItems() {
			// Get inventory item for this product at the fulfillment location
			inventory, err := s.inventoryService.GetInventoryItemByProductAndLocation(ctx, 
				item.GetProductId(), fulfillmentLocationID)
			if err != nil {
				return fmt.Errorf("failed to get inventory for product %s at location %s: %w",
					item.GetProductId(), fulfillmentLocationID, err)
			}
			
			// Release reservation
			err = s.inventoryService.ReleaseReservation(ctx, inventory.ID, int32(item.GetQuantity()))
			if err != nil {
				return fmt.Errorf("failed to release reservation for product %s: %w", 
					item.GetProductId(), err)
			}
		}
	}

	return nil
}

// CheckInventoryForOrder checks if there's enough inventory for an order
func (s *InventoryOrderService) CheckInventoryForOrder(
	ctx context.Context,
	orderID string,
) (bool, error) {
	// Get the order details
	orderResp, err := s.orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: orderID})
	if err != nil {
		return false, fmt.Errorf("failed to get order: %w", err)
	}

	order := orderResp.GetOrder()
	
	// Get fulfillment location from order metadata or shipping address - fallback to local location if not specified
	fulfillmentLocationID := s.localLocationID
	// Check order metadata for location (this assumes the service might be updated to include location)
	// For now we'll use the default location but this can be extended when the order service is updated
	
	// Check inventory for each item at the specified fulfillment location
	for _, item := range order.GetItems() {
		// Get inventory for product at this specific location
		inventory, err := s.inventoryService.GetInventoryItemByProductAndLocation(ctx, 
			item.GetProductId(), fulfillmentLocationID)
		if err != nil {
			// If not found at the specified location, we could check other locations or just report not available
			s.logger.Warn("Product not available at requested fulfillment location", 
				zap.String("product_id", item.GetProductId()),
				zap.String("location_id", fulfillmentLocationID))
			return false, nil
		}

		// Check if there's enough available stock
		if !inventory.IsAvailable(int32(item.GetQuantity())) {
			s.logger.Warn("Insufficient stock",
				zap.String("product_id", item.GetProductId()),
				zap.String("location_id", fulfillmentLocationID),
				zap.Int32("available", inventory.Quantity - inventory.Reserved),
				zap.Int32("requested", int32(item.GetQuantity())))
			return false, nil
		}
	}

	return true, nil
}

// FindBestFulfillmentLocation determines the best location to fulfill an order
// based on inventory availability and proximity to delivery address
func (s *InventoryOrderService) FindBestFulfillmentLocation(
	ctx context.Context,
	orderID string,
) (string, error) {
	// Get the order details
	orderResp, err := s.orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: orderID})
	if err != nil {
		return "", fmt.Errorf("failed to get order: %w", err)
	}
	
	order := orderResp.GetOrder()

	// Get all active locations
	locations, err := s.locationService.ListLocations(ctx, 100, 0, false)
	if err != nil {
		return "", fmt.Errorf("failed to list locations: %w", err)
	}
	
	// Filter to store-type locations first (for customer-facing orders)
	storeLocations := make([]*domain.StoreLocation, 0)
	for _, loc := range locations {
		if loc.Type == "store" { // LocationTypeStore constant
			storeLocations = append(storeLocations, loc)
		}
	}
	
	// Simple implementation - find first location that has all items in stock
	// A more advanced implementation would consider proximity to delivery address
	// and optimize for shipping costs or delivery time
	for _, location := range storeLocations {
		hasAllItems := true
		
		// Check if this location has enough inventory for all items
		for _, item := range order.GetItems() {
			// Get inventory for this product at this location
			inventory, err := s.inventoryService.GetInventoryItemByProductAndLocation(
				ctx, item.GetProductId(), location.ID)
			
			if err != nil || inventory == nil || !inventory.IsAvailable(int32(item.GetQuantity())) {
				hasAllItems = false
				break
			}
		}
		
		if hasAllItems {
			return location.ID, nil
		}
	}
	
	// If no store has all items, try warehouses next
	warehouseLocations := make([]*domain.StoreLocation, 0)
	for _, loc := range locations {
		if loc.Type == "warehouse" || loc.Type == "fulfillment_center" { // LocationType constants
			warehouseLocations = append(warehouseLocations, loc)
		}
	}
	
	for _, location := range warehouseLocations {
		hasAllItems := true
		
		// Check if this location has enough inventory for all items
		for _, item := range order.GetItems() {
			inventory, err := s.inventoryService.GetInventoryItemByProductAndLocation(
				ctx, item.GetProductId(), location.ID)
				
			if err != nil || inventory == nil || !inventory.IsAvailable(int32(item.GetQuantity())) {
				hasAllItems = false
				break
			}
		}
		
		if hasAllItems {
			return location.ID, nil
		}
	}
	
	// No single location can fulfill the entire order
	// For a more advanced system, we could split the order across locations
	s.logger.Warn("No location has sufficient inventory for all items in order", 
		zap.String("order_id", orderID))
	return "", fmt.Errorf("no location has sufficient inventory for all items in order")
}
