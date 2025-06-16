package order

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	orderv1 "github.com/leonvanderhaeghen/stockplatform/services/orderSvc/api/gen/go/proto/order/v1"
)

// Client provides a high-level interface for interacting with the Order service
type Client struct {
	conn   *grpc.ClientConn
	client orderv1.OrderServiceClient
	logger *zap.Logger
}

// Config holds configuration for the Order client
type Config struct {
	Address string
	Timeout time.Duration
}

// New creates a new Order service client
func New(config Config, logger *zap.Logger) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	conn, err := grpc.Dial(config.Address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	client := orderv1.NewOrderServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the connection to the Order service
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateOrder creates a new order
func (c *Client) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	c.logger.Debug("Creating order", zap.String("user_id", req.UserId))
	
	resp, err := c.client.CreateOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create order", zap.Error(err))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	c.logger.Debug("Order created successfully", zap.String("id", resp.Order.Id))
	return resp, nil
}

// GetOrder retrieves an order by ID
func (c *Client) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	c.logger.Debug("Getting order", zap.String("id", req.Id))
	
	resp, err := c.client.GetOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get order", zap.Error(err))
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return resp, nil
}

// GetUserOrders retrieves orders for a specific user
func (c *Client) GetUserOrders(ctx context.Context, req *orderv1.GetUserOrdersRequest) (*orderv1.GetUserOrdersResponse, error) {
	c.logger.Debug("Getting user orders", zap.String("user_id", req.UserId))
	
	resp, err := c.client.GetUserOrders(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user orders", zap.Error(err))
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	
	return resp, nil
}

// UpdateOrder updates an existing order
func (c *Client) UpdateOrder(ctx context.Context, req *orderv1.UpdateOrderRequest) (*orderv1.UpdateOrderResponse, error) {
	c.logger.Debug("Updating order", zap.String("id", req.Order.Id))
	
	resp, err := c.client.UpdateOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update order", zap.Error(err))
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	
	return resp, nil
}

// ListOrders lists all orders with optional filtering and pagination
func (c *Client) ListOrders(ctx context.Context, req *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	c.logger.Debug("Listing orders")
	
	resp, err := c.client.ListOrders(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list orders", zap.Error(err))
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	
	return resp, nil
}

// UpdateOrderStatus updates the status of an order
func (c *Client) UpdateOrderStatus(ctx context.Context, req *orderv1.UpdateOrderStatusRequest) (*orderv1.UpdateOrderStatusResponse, error) {
	c.logger.Debug("Updating order status", zap.String("id", req.Id))
	
	resp, err := c.client.UpdateOrderStatus(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update order status", zap.Error(err))
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}
	
	return resp, nil
}

// CancelOrder cancels an order
func (c *Client) CancelOrder(ctx context.Context, req *orderv1.CancelOrderRequest) (*orderv1.CancelOrderResponse, error) {
	c.logger.Debug("Cancelling order", zap.String("id", req.Id))
	
	resp, err := c.client.CancelOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to cancel order", zap.Error(err))
		return nil, fmt.Errorf("failed to cancel order: %w", err)
	}
	
	return resp, nil
}
