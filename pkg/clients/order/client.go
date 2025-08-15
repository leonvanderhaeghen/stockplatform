package order

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
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
func (c *Client) CreateOrder(ctx context.Context, userID string, items []*models.OrderItem, shippingAddress *models.Address, notes string) (*models.CreateOrderResponse, error) {
	c.logger.Debug("Creating order", zap.String("user_id", userID))
	
	// Convert domain items to protobuf items
	protoItems := make([]*orderv1.OrderItem, len(items))
	for i, item := range items {
		protoItems[i] = c.convertFromOrderItem(item)
	}
	
	req := &orderv1.CreateOrderRequest{
		UserId: userID,
		Items:  protoItems,
		// Notes field not available in protobuf schema
	}
	
	// Add shipping address if provided
	if shippingAddress != nil {
		req.ShippingAddress = c.convertFromAddress(shippingAddress)
	}
	
	resp, err := c.client.CreateOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create order", zap.Error(err))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	c.logger.Debug("Order created successfully", zap.String("id", resp.Order.Id))
	return c.convertToCreateOrderResponse(resp), nil
}

// GetOrder retrieves an order by ID
func (c *Client) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	c.logger.Debug("Getting order", zap.String("id", id))
	
	req := &orderv1.GetOrderRequest{
		Id: id,
	}
	
	resp, err := c.client.GetOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get order", zap.Error(err))
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	return c.convertToOrder(resp.Order), nil
}

// GetUserOrders retrieves orders for a specific user
func (c *Client) GetUserOrders(ctx context.Context, userID string, limit, offset int32) ([]*models.Order, error) {
	c.logger.Debug("Getting user orders", zap.String("user_id", userID))
	
	req := &orderv1.GetUserOrdersRequest{
		UserId: userID,
		Limit:  limit,
		Offset: offset,
	}
	
	resp, err := c.client.GetUserOrders(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user orders", zap.Error(err))
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}
	
	// Convert orders to domain models
	orders := make([]*models.Order, len(resp.Orders))
	for i, protoOrder := range resp.Orders {
		orders[i] = c.convertToOrder(protoOrder)
	}
	
	return orders, nil
}

// UpdateOrder updates an existing order
func (c *Client) UpdateOrder(ctx context.Context, order *models.Order) (*models.UpdateOrderResponse, error) {
	c.logger.Debug("Updating order", zap.String("id", order.ID))
	
	req := &orderv1.UpdateOrderRequest{
		Order: c.convertFromOrder(order),
	}
	
	resp, err := c.client.UpdateOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update order", zap.Error(err))
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	
	return c.convertToUpdateOrderResponse(resp), nil
}

// ListOrders lists all orders with optional filtering and pagination
func (c *Client) ListOrders(ctx context.Context, status string, userID string, limit, offset int32) (*models.ListOrdersResponse, error) {
	c.logger.Debug("Listing orders", zap.String("status", status), zap.String("user_id", userID))
	
	req := &orderv1.ListOrdersRequest{
		Status: status,
		// UserId field not available in ListOrdersRequest
		Limit:  limit,
		Offset: offset,
	}
	
	resp, err := c.client.ListOrders(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list orders", zap.Error(err))
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	
	return c.convertToListOrdersResponse(resp), nil
}

// UpdateOrderStatus updates the status of an order
func (c *Client) UpdateOrderStatus(ctx context.Context, id, status string) error {
	c.logger.Debug("Updating order status", zap.String("id", id), zap.String("status", status))
	
	req := &orderv1.UpdateOrderStatusRequest{
		Id:     id,
		Status: convertStringToOrderStatus(status),
	}
	
	_, err := c.client.UpdateOrderStatus(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update order status", zap.Error(err))
		return fmt.Errorf("failed to update order status: %w", err)
	}
	
	c.logger.Debug("Order status updated successfully", zap.String("id", id), zap.String("status", status))
	return nil
}

// CancelOrder cancels an order
func (c *Client) CancelOrder(ctx context.Context, id, reason string) error {
	c.logger.Debug("Cancelling order", zap.String("id", id))
	
	req := &orderv1.CancelOrderRequest{
		Id: id,
		// Reason field not available in protobuf schema
	}
	
	_, err := c.client.CancelOrder(ctx, req)
	if err != nil {
		c.logger.Error("Failed to cancel order", zap.Error(err))
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	
	c.logger.Debug("Order cancelled successfully", zap.String("id", id))
	return nil
}

// AddPayment adds payment information to an order
func (c *Client) AddPayment(ctx context.Context, orderID, method, transactionID string, amount float64) error {
	req := &orderv1.AddPaymentRequest{
		OrderId:       orderID,
		Method:        method,
		TransactionId: transactionID,
		Amount:        amount,
	}

	_, err := c.client.AddPayment(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add payment to order: %w", err)
	}

	return nil
}

// AddTrackingCode adds a tracking code to an order
func (c *Client) AddTrackingCode(ctx context.Context, orderID, trackingCode string) error {
	req := &orderv1.AddTrackingCodeRequest{
		OrderId:      orderID,
		TrackingCode: trackingCode,
	}

	_, err := c.client.AddTrackingCode(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add tracking code to order: %w", err)
	}

	return nil
}

// Helper function to convert string status to protobuf enum
func convertStringToOrderStatus(status string) orderv1.OrderStatus {
	switch status {
	case "CREATED":
		return orderv1.OrderStatus_ORDER_STATUS_CREATED
	case "PENDING":
		return orderv1.OrderStatus_ORDER_STATUS_PENDING
	case "PAID":
		return orderv1.OrderStatus_ORDER_STATUS_PAID
	case "SHIPPED":
		return orderv1.OrderStatus_ORDER_STATUS_SHIPPED
	case "DELIVERED":
		return orderv1.OrderStatus_ORDER_STATUS_DELIVERED
	case "CANCELLED":
		return orderv1.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
