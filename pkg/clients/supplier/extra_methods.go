package supplier

import (
    "context"
    "fmt"

    "go.uber.org/zap"

    supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
)

// ListAdapters lists supplier adapters
func (c *Client) ListAdapters(ctx context.Context, req *supplierv1.ListAdaptersRequest) (*supplierv1.ListAdaptersResponse, error) {
    c.logger.Debug("Listing supplier adapters")

    resp, err := c.client.ListAdapters(ctx, req)
    if err != nil {
        c.logger.Error("Failed to list adapters", zap.Error(err))
        return nil, fmt.Errorf("failed to list adapters: %w", err)
    }

    return resp, nil
}

// GetAdapterCapabilities returns capabilities of a supplier adapter
func (c *Client) GetAdapterCapabilities(ctx context.Context, req *supplierv1.GetAdapterCapabilitiesRequest) (*supplierv1.GetAdapterCapabilitiesResponse, error) {
    c.logger.Debug("Getting adapter capabilities", zap.String("adapter", req.GetAdapterName()))

    resp, err := c.client.GetAdapterCapabilities(ctx, req)
    if err != nil {
        c.logger.Error("Failed to get adapter capabilities", zap.Error(err))
        return nil, fmt.Errorf("failed to get adapter capabilities: %w", err)
    }

    return resp, nil
}

// TestAdapterConnection tests a supplier adapter connection
func (c *Client) TestAdapterConnection(ctx context.Context, req *supplierv1.TestAdapterConnectionRequest) (*supplierv1.TestAdapterConnectionResponse, error) {
    c.logger.Debug("Testing adapter connection", zap.String("adapter", req.GetAdapterName()))

    resp, err := c.client.TestAdapterConnection(ctx, req)
    if err != nil {
        c.logger.Error("Failed to test adapter connection", zap.Error(err))
        return nil, fmt.Errorf("failed to test adapter connection: %w", err)
    }

    return resp, nil
}

// SyncProducts syncs products from a supplier
func (c *Client) SyncProducts(ctx context.Context, req *supplierv1.SyncProductsRequest) (*supplierv1.SyncProductsResponse, error) {
    c.logger.Debug("Syncing products", zap.String("supplier_id", req.SupplierId))

    resp, err := c.client.SyncProducts(ctx, req)
    if err != nil {
        c.logger.Error("Failed to sync products", zap.Error(err))
        return nil, fmt.Errorf("failed to sync products: %w", err)
    }

    return resp, nil
}

// SyncInventory syncs inventory from a supplier
func (c *Client) SyncInventory(ctx context.Context, req *supplierv1.SyncInventoryRequest) (*supplierv1.SyncInventoryResponse, error) {
    c.logger.Debug("Syncing inventory", zap.String("supplier_id", req.SupplierId))

    resp, err := c.client.SyncInventory(ctx, req)
    if err != nil {
        c.logger.Error("Failed to sync inventory", zap.Error(err))
        return nil, fmt.Errorf("failed to sync inventory: %w", err)
    }

    return resp, nil
}
