package grpc

import (
	"context"
	"fmt"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/inventory/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckAvailability checks real-time inventory availability at store locations
func (s *InventoryServer) CheckAvailability(ctx context.Context, req *inventoryv1.CheckAvailabilityRequest) (*inventoryv1.CheckAvailabilityResponse, error) {
	logger := s.logger.With(zap.String("handler", "CheckAvailability"))
	logger.Info("Processing CheckAvailability request")

	// Validate request
	if req.LocationId == "" {
		return nil, status.Error(codes.InvalidArgument, "location_id is required")
	}

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one inventory item is required")
	}

	locationID := req.LocationId
	
	// Get location info
	location, err := s.locationService.GetLocation(ctx, locationID)
	if err != nil {
		logger.Error("Failed to get location", zap.Error(err), zap.String("location_id", locationID))
		return nil, status.Errorf(codes.NotFound, "location not found: %v", err)
	}

	// Process each requested item
	results := make([]*inventoryv1.ItemAvailability, 0, len(req.Items))
	allAvailable := true

	for _, item := range req.Items {
		// Create the result structure for this item
		result := &inventoryv1.ItemAvailability{
			ProductId: item.ProductId,
			Sku:       item.Sku,
			InStock:   false,
			Status:    "out_of_stock",
		}

		// First try to get by product ID if provided
		var protoItem *inventoryv1.InventoryItem
		if item.ProductId != "" {
			// Get inventory for this product at the specified location
			domainItem, err := s.service.GetInventoryItemByProductAndLocation(ctx, item.ProductId, locationID)
			if err == nil && domainItem != nil {
				// Convert domain inventory item to proto inventory item
				protoItem = convertDomainToProtoInventoryItem(domainItem)
				result.InventoryItemId = protoItem.Id
			}
		} else if item.Sku != "" {
			// If no product ID or no results by product ID, try by SKU
			domainItem, err := s.service.GetInventoryItemBySKUAndLocation(ctx, item.Sku, locationID)
			if err == nil && domainItem != nil {
				// Convert domain inventory item to proto inventory item
				protoItem = convertDomainToProtoInventoryItem(domainItem)
				result.InventoryItemId = protoItem.Id
			}
		}

		// Calculate availability
		if protoItem != nil {
			// Calculate available quantity (total minus reserved)
			availableQty := protoItem.Quantity - protoItem.Reserved
			result.AvailableQuantity = availableQty

			// Determine if the requested quantity is available
			requiredQty := item.Quantity
			if availableQty >= requiredQty {
				result.InStock = true
				// Check if it's low stock
				reorderPoint := int32(0)
				
				if availableQty <= 0 {
					result.Status = "out_of_stock"
					allAvailable = false
				} else if availableQty <= reorderPoint && reorderPoint > 0 {
					result.Status = "low_stock"
				} else {
					result.Status = "in_stock"
				}
			} else {
				// Not enough stock
				allAvailable = false
				if availableQty <= 0 {
					result.Status = "out_of_stock"
					result.Message = "Item is out of stock"
				} else {
					result.Status = "insufficient_stock"
					result.Message = fmt.Sprintf("Only %d available, %d requested", availableQty, requiredQty)
				}
			}
		} else {
			// No inventory found
			result.Status = "not_found"
			result.Message = "Item not found in this location"
			allAvailable = false
		}

		results = append(results, result)
	}

	logger.Info("CheckAvailability request completed", zap.Int("results_count", len(results)))

	return &inventoryv1.CheckAvailabilityResponse{
		LocationId:    locationID,
		LocationName:  location.Name,
		Items:         results,
		AllAvailable:  allAvailable,
	}, nil
}

// Helper function to convert domain inventory item to proto inventory item
func convertDomainToProtoInventoryItem(item *domain.InventoryItem) *inventoryv1.InventoryItem {
	if item == nil {
		return &inventoryv1.InventoryItem{}
	}
	
	return &inventoryv1.InventoryItem{
		Id:               item.ID,
		ProductId:        item.ProductID,
		Quantity:         item.Quantity,
		Reserved:         item.Reserved,
		Sku:              item.SKU,
		LocationId:       item.LocationID,
		ShelfLocation:    item.ShelfLocation,
		// Map other fields as needed
		ReorderThreshold: int32(item.ReorderPoint),
		ReorderAmount:    int32(item.ReorderQuantity),
	}
}
