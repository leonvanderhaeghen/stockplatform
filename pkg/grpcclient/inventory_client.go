package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	inventorypb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/inventory/v1"
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

// GetInventoryBySKU retrieves an inventory item by SKU
func (c *InventoryClient) GetInventoryBySKU(ctx context.Context, req *inventorypb.GetInventoryBySKURequest) (*inventorypb.GetInventoryResponse, error) {
	return c.client.GetInventoryBySKU(ctx, req)
}

// AddStock adds stock to an inventory item
func (c *InventoryClient) AddStock(ctx context.Context, req *inventorypb.AddStockRequest) (*inventorypb.AddStockResponse, error) {
	return c.client.AddStock(ctx, req)
}

// RemoveStock removes stock from an inventory item
func (c *InventoryClient) RemoveStock(ctx context.Context, req *inventorypb.RemoveStockRequest) (*inventorypb.RemoveStockResponse, error) {
	return c.client.RemoveStock(ctx, req)
}

// ReleaseReservation releases a reservation without fulfilling it
func (c *InventoryClient) ReleaseReservation(ctx context.Context, req *inventorypb.ReleaseReservationRequest) (*inventorypb.ReleaseReservationResponse, error) {
	return c.client.ReleaseReservation(ctx, req)
}

// FulfillReservation completes a reservation and deducts from stock
func (c *InventoryClient) FulfillReservation(ctx context.Context, req *inventorypb.FulfillReservationRequest) (*inventorypb.FulfillReservationResponse, error) {
	return c.client.FulfillReservation(ctx, req)
}

// ListInventoryByLocation lists inventory items for a specific location
func (c *InventoryClient) ListInventoryByLocation(ctx context.Context, req *inventorypb.ListInventoryByLocationRequest) (*inventorypb.ListInventoryResponse, error) {
	return c.client.ListInventoryByLocation(ctx, req)
}

// CreateLocation creates a new store location
func (c *InventoryClient) CreateLocation(ctx context.Context, req *inventorypb.CreateLocationRequest) (*inventorypb.CreateLocationResponse, error) {
	return c.client.CreateLocation(ctx, req)
}

// GetLocation retrieves a store location by ID
func (c *InventoryClient) GetLocation(ctx context.Context, req *inventorypb.GetLocationRequest) (*inventorypb.GetLocationResponse, error) {
	return c.client.GetLocation(ctx, req)
}

// UpdateLocation updates an existing store location
func (c *InventoryClient) UpdateLocation(ctx context.Context, req *inventorypb.UpdateLocationRequest) (*inventorypb.UpdateLocationResponse, error) {
	return c.client.UpdateLocation(ctx, req)
}

// DeleteLocation removes a store location
func (c *InventoryClient) DeleteLocation(ctx context.Context, req *inventorypb.DeleteLocationRequest) (*inventorypb.DeleteLocationResponse, error) {
	return c.client.DeleteLocation(ctx, req)
}

// ListLocations lists all store locations with pagination
func (c *InventoryClient) ListLocations(ctx context.Context, req *inventorypb.ListLocationsRequest) (*inventorypb.ListLocationsResponse, error) {
	return c.client.ListLocations(ctx, req)
}

// CreateTransfer creates a new inventory transfer between locations
func (c *InventoryClient) CreateTransfer(ctx context.Context, req *inventorypb.CreateTransferRequest) (*inventorypb.CreateTransferResponse, error) {
	return c.client.CreateTransfer(ctx, req)
}

// GetTransfer retrieves an inventory transfer by ID
func (c *InventoryClient) GetTransfer(ctx context.Context, req *inventorypb.GetTransferRequest) (*inventorypb.GetTransferResponse, error) {
	return c.client.GetTransfer(ctx, req)
}

// UpdateTransferStatus updates the status of a transfer
func (c *InventoryClient) UpdateTransferStatus(ctx context.Context, req *inventorypb.UpdateTransferStatusRequest) (*inventorypb.UpdateTransferStatusResponse, error) {
	return c.client.UpdateTransferStatus(ctx, req)
}

// ListTransfers lists transfers with pagination and filters
func (c *InventoryClient) ListTransfers(ctx context.Context, req *inventorypb.ListTransfersRequest) (*inventorypb.ListTransfersResponse, error) {
	return c.client.ListTransfers(ctx, req)
}
