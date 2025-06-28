package inventory

import (
    "context"
    "fmt"

    "go.uber.org/zap"

    inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
)

// GetInventoryBySKU retrieves an inventory item by SKU
func (c *Client) GetInventoryBySKU(ctx context.Context, req *inventoryv1.GetInventoryBySKURequest) (*inventoryv1.GetInventoryResponse, error) {
    c.logger.Debug("Getting inventory by SKU", zap.String("sku", req.Sku))

    resp, err := c.client.GetInventoryBySKU(ctx, req)
    if err != nil {
        c.logger.Error("Failed to get inventory by SKU", zap.Error(err))
        return nil, fmt.Errorf("failed to get inventory by SKU: %w", err)
    }

    return resp, nil
}

// DeleteInventory deletes an inventory item by ID
func (c *Client) DeleteInventory(ctx context.Context, req *inventoryv1.DeleteInventoryRequest) (*inventoryv1.DeleteInventoryResponse, error) {
    c.logger.Debug("Deleting inventory", zap.String("id", req.Id))

    resp, err := c.client.DeleteInventory(ctx, req)
    if err != nil {
        c.logger.Error("Failed to delete inventory", zap.Error(err))
        return nil, fmt.Errorf("failed to delete inventory: %w", err)
    }

    return resp, nil
}

// CompletePickup completes an in-store pickup reservation
func (c *Client) CompletePickup(ctx context.Context, req *inventoryv1.CompletePickupRequest) (*inventoryv1.CompletePickupResponse, error) {
    c.logger.Debug("Completing pickup", zap.String("reservation_id", req.ReservationId))

    resp, err := c.client.CompletePickup(ctx, req)
    if err != nil {
        c.logger.Error("Failed to complete pickup", zap.Error(err))
        return nil, fmt.Errorf("failed to complete pickup: %w", err)
    }

    return resp, nil
}
