package store

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	storev1 "github.com/leonvanderhaeghen/stockplatform/services/storeSvc/api/gen/go/proto/store/v1"
)

// Client provides an abstraction for the store service
type Client struct {
	conn   *grpc.ClientConn
	client storev1.StoreServiceClient
}

// NewClient creates a new store service client
func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to store service: %w", err)
	}

	client := storev1.NewStoreServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Store Management Methods
func (c *Client) CreateStore(ctx context.Context, req *storev1.CreateStoreRequest) (*storev1.CreateStoreResponse, error) {
	return c.client.CreateStore(ctx, req)
}

func (c *Client) GetStore(ctx context.Context, req *storev1.GetStoreRequest) (*storev1.GetStoreResponse, error) {
	return c.client.GetStore(ctx, req)
}

func (c *Client) ListStores(ctx context.Context, req *storev1.ListStoresRequest) (*storev1.ListStoresResponse, error) {
	return c.client.ListStores(ctx, req)
}

func (c *Client) UpdateStore(ctx context.Context, req *storev1.UpdateStoreRequest) (*storev1.UpdateStoreResponse, error) {
	return c.client.UpdateStore(ctx, req)
}

func (c *Client) DeleteStore(ctx context.Context, req *storev1.DeleteStoreRequest) (*storev1.DeleteStoreResponse, error) {
	return c.client.DeleteStore(ctx, req)
}

// Store Inventory Management Methods
func (c *Client) AddProductToStore(ctx context.Context, req *storev1.AddProductToStoreRequest) (*storev1.AddProductToStoreResponse, error) {
	return c.client.AddProductToStore(ctx, req)
}

func (c *Client) UpdateStoreProductStock(ctx context.Context, req *storev1.UpdateStoreProductStockRequest) (*storev1.UpdateStoreProductStockResponse, error) {
	return c.client.UpdateStoreProductStock(ctx, req)
}

func (c *Client) RemoveProductFromStore(ctx context.Context, req *storev1.RemoveProductFromStoreRequest) (*storev1.RemoveProductFromStoreResponse, error) {
	return c.client.RemoveProductFromStore(ctx, req)
}

func (c *Client) GetStoreProducts(ctx context.Context, req *storev1.GetStoreProductsRequest) (*storev1.GetStoreProductsResponse, error) {
	return c.client.GetStoreProducts(ctx, req)
}

func (c *Client) GetProductStoreLocations(ctx context.Context, req *storev1.GetProductStoreLocationsRequest) (*storev1.GetProductStoreLocationsResponse, error) {
	return c.client.GetProductStoreLocations(ctx, req)
}

// Product Reservation Methods
func (c *Client) ReserveProduct(ctx context.Context, req *storev1.ReserveProductRequest) (*storev1.ReserveProductResponse, error) {
	return c.client.ReserveProduct(ctx, req)
}

func (c *Client) CancelReservation(ctx context.Context, req *storev1.CancelReservationRequest) (*storev1.CancelReservationResponse, error) {
	return c.client.CancelReservation(ctx, req)
}

func (c *Client) GetReservations(ctx context.Context, req *storev1.GetReservationsRequest) (*storev1.GetReservationsResponse, error) {
	return c.client.GetReservations(ctx, req)
}

func (c *Client) CompleteReservation(ctx context.Context, req *storev1.CompleteReservationRequest) (*storev1.CompleteReservationResponse, error) {
	return c.client.CompleteReservation(ctx, req)
}

// Store User Management Methods
func (c *Client) AssignUserToStore(ctx context.Context, req *storev1.AssignUserToStoreRequest) (*storev1.AssignUserToStoreResponse, error) {
	return c.client.AssignUserToStore(ctx, req)
}

func (c *Client) RemoveUserFromStore(ctx context.Context, req *storev1.RemoveUserFromStoreRequest) (*storev1.RemoveUserFromStoreResponse, error) {
	return c.client.RemoveUserFromStore(ctx, req)
}

func (c *Client) GetStoreUsers(ctx context.Context, req *storev1.GetStoreUsersRequest) (*storev1.GetStoreUsersResponse, error) {
	return c.client.GetStoreUsers(ctx, req)
}

func (c *Client) GetUserStores(ctx context.Context, req *storev1.GetUserStoresRequest) (*storev1.GetUserStoresResponse, error) {
	return c.client.GetUserStores(ctx, req)
}

// Sales Tracking Methods
func (c *Client) RecordSale(ctx context.Context, req *storev1.RecordSaleRequest) (*storev1.RecordSaleResponse, error) {
	return c.client.RecordSale(ctx, req)
}

func (c *Client) GetStoreSales(ctx context.Context, req *storev1.GetStoreSalesRequest) (*storev1.GetStoreSalesResponse, error) {
	return c.client.GetStoreSales(ctx, req)
}

// Export Methods
func (c *Client) ExportStoreProducts(ctx context.Context, req *storev1.ExportStoreProductsRequest) (*storev1.ExportStoreProductsResponse, error) {
	return c.client.ExportStoreProducts(ctx, req)
}

func (c *Client) ExportStoreSales(ctx context.Context, req *storev1.ExportStoreSalesRequest) (*storev1.ExportStoreSalesResponse, error) {
	return c.client.ExportStoreSales(ctx, req)
}
