package supplier

import (
    "context"
    "fmt"

    "go.uber.org/zap"
    "github.com/leonvanderhaeghen/stockplatform/pkg/models"
    supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
)

// ListAdapters lists available supplier adapters
func (c *Client) ListAdapters(ctx context.Context) (*models.ListAdaptersResponse, error) {
    c.logger.Debug("Listing adapters")
    
    req := &supplierv1.ListAdaptersRequest{}
    
    resp, err := c.client.ListAdapters(ctx, req)
    if err != nil {
        c.logger.Error("Failed to list adapters", zap.Error(err))
        return nil, fmt.Errorf("failed to list adapters: %w", err)
    }
    
    return c.convertToListAdaptersResponse(resp), nil
}

// GetAdapterCapabilities gets capabilities of a supplier adapter
func (c *Client) GetAdapterCapabilities(ctx context.Context, adapterName string) (*models.AdapterCapabilities, error) {
    c.logger.Debug("Getting adapter capabilities", zap.String("adapter_name", adapterName))
    
    req := &supplierv1.GetAdapterCapabilitiesRequest{
        AdapterName: adapterName,
    }
    
    resp, err := c.client.GetAdapterCapabilities(ctx, req)
    if err != nil {
        c.logger.Error("Failed to get adapter capabilities", zap.Error(err))
        return nil, fmt.Errorf("failed to get adapter capabilities: %w", err)
    }
    
    return c.convertToAdapterCapabilities(resp.Capabilities), nil
}

// TestAdapterConnection tests connection to a supplier adapter
func (c *Client) TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) (*models.TestConnectionResponse, error) {
    c.logger.Debug("Testing adapter connection", zap.String("adapter_name", adapterName))
    
    req := &supplierv1.TestAdapterConnectionRequest{
        AdapterName: adapterName,
        Config:      config,
    }
    
    resp, err := c.client.TestAdapterConnection(ctx, req)
    if err != nil {
        c.logger.Error("Failed to test adapter connection", zap.Error(err))
        return nil, fmt.Errorf("failed to test adapter connection: %w", err)
    }
    
    return c.convertToTestConnectionResponse(resp), nil
}

// SyncProducts synchronizes products from a supplier
func (c *Client) SyncProducts(ctx context.Context, supplierID string, fullSync, dryRun bool, batchSize int32) (*models.SyncResponse, error) {
    c.logger.Debug("Syncing products", 
        zap.String("supplier_id", supplierID),
        zap.Bool("full_sync", fullSync),
        zap.Bool("dry_run", dryRun),
        zap.Int32("batch_size", batchSize))
    
    req := &supplierv1.SyncProductsRequest{
        SupplierId: supplierID,
        Options: &supplierv1.SyncOptions{
            FullSync:  fullSync,
            BatchSize: batchSize,
        },
    }
    
    resp, err := c.client.SyncProducts(ctx, req)
    if err != nil {
        c.logger.Error("Failed to sync products", zap.Error(err))
        return nil, fmt.Errorf("failed to sync products: %w", err)
    }
    
    return c.convertToSyncResponse(resp), nil
}

// SyncInventory synchronizes inventory from a supplier
func (c *Client) SyncInventory(ctx context.Context, supplierID string, fullSync, dryRun bool, batchSize int32) (*models.SyncResponse, error) {
    c.logger.Debug("Syncing inventory", 
        zap.String("supplier_id", supplierID),
        zap.Bool("full_sync", fullSync),
        zap.Bool("dry_run", dryRun),
        zap.Int32("batch_size", batchSize))
    
    req := &supplierv1.SyncInventoryRequest{
        SupplierId: supplierID,
        Options: &supplierv1.SyncOptions{
            FullSync:  fullSync,
            BatchSize: batchSize,
            // Note: dryRun is not in protobuf SyncOptions
        },
    }
    
    resp, err := c.client.SyncInventory(ctx, req)
    if err != nil {
        c.logger.Error("Failed to sync inventory", zap.Error(err))
        return nil, fmt.Errorf("failed to sync inventory: %w", err)
    }
    
    return c.convertToSyncInventoryResponse(resp), nil
}
