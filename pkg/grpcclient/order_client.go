package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	orderpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/order/v1"
)

type OrderClient struct {
	conn   *grpc.ClientConn
	client orderpb.OrderServiceClient
}

// NewOrderClient creates a new gRPC client for the Order service
func NewOrderClient(addr string) (*OrderClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := orderpb.NewOrderServiceClient(conn)
	return &OrderClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *OrderClient) Close() error {
	return c.conn.Close()
}

// CreateOrder creates a new order
func (c *OrderClient) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	return c.client.CreateOrder(ctx, req)
}

// GetOrder retrieves an order by ID
func (c *OrderClient) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	return c.client.GetOrder(ctx, req)
}

// GetUserOrders retrieves orders for a specific user
func (c *OrderClient) GetUserOrders(ctx context.Context, req *orderpb.GetUserOrdersRequest) (*orderpb.GetUserOrdersResponse, error) {
	return c.client.GetUserOrders(ctx, req)
}

// UpdateOrder updates an existing order
func (c *OrderClient) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	return c.client.UpdateOrder(ctx, req)
}

// DeleteOrder deletes an order by ID
func (c *OrderClient) DeleteOrder(ctx context.Context, req *orderpb.DeleteOrderRequest) (*orderpb.DeleteOrderResponse, error) {
	return c.client.DeleteOrder(ctx, req)
}

// ListOrders lists all orders with optional filtering and pagination
func (c *OrderClient) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	return c.client.ListOrders(ctx, req)
}

// UpdateOrderStatus updates the status of an order
func (c *OrderClient) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.UpdateOrderStatusResponse, error) {
	return c.client.UpdateOrderStatus(ctx, req)
}

// AddPayment adds payment information to an order
func (c *OrderClient) AddPayment(ctx context.Context, req *orderpb.AddPaymentRequest) (*orderpb.AddPaymentResponse, error) {
	return c.client.AddPayment(ctx, req)
}

// AddTrackingCode adds a tracking code to an order
func (c *OrderClient) AddTrackingCode(ctx context.Context, req *orderpb.AddTrackingCodeRequest) (*orderpb.AddTrackingCodeResponse, error) {
	return c.client.AddTrackingCode(ctx, req)
}

// CancelOrder cancels an order
func (c *OrderClient) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	return c.client.CancelOrder(ctx, req)
}
