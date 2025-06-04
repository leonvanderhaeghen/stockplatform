package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	inventorypb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/inventory/v1"
)

type InventoryClient struct {
	conn   *grpc.ClientConn
	client inventorypb.InventoryServiceClient
}

// NewInventoryClient creates a new gRPC client for the Inventory service
func NewInventoryClient(addr string) (*InventoryClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := inventorypb.NewInventoryServiceClient(conn)
	return &InventoryClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *InventoryClient) Close() error {
	return c.conn.Close()
}

// CreateInventory creates a new inventory item
func (c *InventoryClient) CreateInventory(ctx context.Context, req *inventorypb.CreateInventoryRequest) (*inventorypb.CreateInventoryResponse, error) {
	return c.client.CreateInventory(ctx, req)
}

// GetInventory retrieves an inventory item by ID
func (c *InventoryClient) GetInventory(ctx context.Context, req *inventorypb.GetInventoryRequest) (*inventorypb.GetInventoryResponse, error) {
	return c.client.GetInventory(ctx, req)
}

// GetInventoryByProductID retrieves an inventory item by product ID
func (c *InventoryClient) GetInventoryByProductID(ctx context.Context, req *inventorypb.GetInventoryByProductIDRequest) (*inventorypb.GetInventoryResponse, error) {
	return c.client.GetInventoryByProductID(ctx, req)
}

// UpdateInventory updates an existing inventory item
func (c *InventoryClient) UpdateInventory(ctx context.Context, req *inventorypb.UpdateInventoryRequest) (*inventorypb.UpdateInventoryResponse, error) {
	return c.client.UpdateInventory(ctx, req)
}

// DeleteInventory deletes an inventory item by ID
func (c *InventoryClient) DeleteInventory(ctx context.Context, req *inventorypb.DeleteInventoryRequest) (*inventorypb.DeleteInventoryResponse, error) {
	return c.client.DeleteInventory(ctx, req)
}

// ListInventory lists all inventory items with pagination
func (c *InventoryClient) ListInventory(ctx context.Context, req *inventorypb.ListInventoryRequest) (*inventorypb.ListInventoryResponse, error) {
	return c.client.ListInventory(ctx, req)
}

// ReserveStock reserves inventory for an order
func (c *InventoryClient) ReserveStock(ctx context.Context, req *inventorypb.ReserveStockRequest) (*inventorypb.ReserveStockResponse, error) {
	return c.client.ReserveStock(ctx, req)
}
