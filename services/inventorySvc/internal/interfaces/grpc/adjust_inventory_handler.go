package grpc

import (
	"context"

	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/inventory/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdjustInventoryForOrder handles inventory adjustments from order operations (sales, returns, exchanges)
func (s *InventoryServer) AdjustInventoryForOrder(ctx context.Context, req *inventoryv1.AdjustInventoryForOrderRequest) (*inventoryv1.AdjustInventoryForOrderResponse, error) {
	logger := s.logger.With(
		zap.String("handler", "AdjustInventoryForOrder"),
		zap.String("order_id", req.OrderId),
		zap.String("location_id", req.LocationId),
		zap.String("adjustment_type", req.AdjustmentType),
		zap.String("reference_id", req.ReferenceId),
		zap.Int("items_count", len(req.Items)),
	)

	logger.Info("Processing AdjustInventoryForOrder request")

	// Validate request
	if req.LocationId == "" {
		return nil, status.Error(codes.InvalidArgument, "location ID is required")
	}

	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one item is required")
	}

	// Process each item
	results := make([]*inventoryv1.InventoryAdjustmentResult, 0, len(req.Items))
	allSuccess := true

	for _, item := range req.Items {
		result := &inventoryv1.InventoryAdjustmentResult{
			ProductId:       item.ProductId,
			Sku:             item.Sku,
			Quantity:        item.Quantity,
			Success:         false,
			InventoryItemId: item.InventoryItemId,
		}

		// Handle quantity adjustment based on item.Quantity (positive for additions, negative for removals)
		var err error
		if item.Quantity > 0 {
			// Add stock (e.g., return)
			err = s.service.AddStock(ctx, item.InventoryItemId, item.Quantity)
			if err != nil {
				logger.Error("Failed to add stock", 
					zap.Error(err),
					zap.String("product_id", item.ProductId),
					zap.String("inventory_item_id", item.InventoryItemId),
					zap.Int32("quantity", item.Quantity),
				)

				result.Success = false
				result.ErrorMessage = err.Error()
				allSuccess = false
			} else {
				result.Success = true
			}
		} else if item.Quantity < 0 {
			// Remove stock (e.g., sale)
			// Convert negative quantity to positive for RemoveStock
			removeQuantity := -item.Quantity

			err = s.service.RemoveStock(ctx, item.InventoryItemId, removeQuantity)
			if err != nil {
				logger.Error("Failed to remove stock", 
					zap.Error(err),
					zap.String("product_id", item.ProductId),
					zap.String("inventory_item_id", item.InventoryItemId),
					zap.Int32("quantity", removeQuantity),
				)

				result.Success = false
				result.ErrorMessage = err.Error()
				allSuccess = false
			} else {
				result.Success = true
			}
		} else {
			// Quantity is zero, nothing to do
			result.Success = true
			result.ErrorMessage = "No quantity adjustment needed"
		}

		results = append(results, result)
	}

	logger.Info("AdjustInventoryForOrder request completed", 
		zap.Bool("all_success", allSuccess),
	)

	return &inventoryv1.AdjustInventoryForOrderResponse{
		Success: allSuccess,
		Items:   results,
	}, nil
}

// getAdjustmentReason formats a reason string for inventory adjustments
func getAdjustmentReason(adjustmentType, itemReason string) string {
	baseReason := adjustmentType
	if itemReason != "" {
		return baseReason + ": " + itemReason
	}
	return baseReason
}
