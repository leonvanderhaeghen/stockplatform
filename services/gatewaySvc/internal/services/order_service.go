package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	orderv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/order/v1"
	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
)

// OrderServiceImpl implements the OrderService interface
type OrderServiceImpl struct {
	client *grpcclient.OrderClient
	logger *zap.Logger
}

// NewOrderService creates a new instance of OrderServiceImpl
func NewOrderService(orderServiceAddr string, logger *zap.Logger) (OrderService, error) {
	// Create a gRPC client
	client, err := grpcclient.NewOrderClient(orderServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create order client: %w", err)
	}

	return &OrderServiceImpl{
		client: client,
		logger: logger.Named("order_service"),
	}, nil
}

// GetUserOrders gets orders for a user
func (s *OrderServiceImpl) GetUserOrders(
	ctx context.Context,
	userID, status, startDate, endDate string,
	limit, offset int,
) (interface{}, error) {
	s.logger.Debug("GetUserOrders",
		zap.String("userID", userID),
		zap.String("status", status),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	req := &orderv1.GetUserOrdersRequest{
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resp, err := s.client.GetUserOrders(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get user orders",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	// Filter by status if provided
	if status != "" {
		filtered := make([]*orderv1.Order, 0, len(resp.GetOrders()))
		for _, order := range resp.GetOrders() {
			if order.GetStatus().String() == status {
				filtered = append(filtered, order)
			}
		}
		return filtered, nil
	}

	return resp.GetOrders(), nil
}

// GetUserOrder gets a specific order for a user
func (s *OrderServiceImpl) GetUserOrder(
	ctx context.Context,
	orderID, userID string,
) (interface{}, error) {
	s.logger.Debug("GetUserOrder",
		zap.String("orderID", orderID),
		zap.String("userID", userID),
	)

	req := &orderv1.GetOrderRequest{
		Id: orderID,
		// Note: The GetOrderRequest in the proto file only has an 'id' field, not 'OrderId' or 'UserId'
	}

	resp, err := s.client.GetOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get order",
			zap.String("orderID", orderID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order := resp.GetOrder()
	if order == nil {
		s.logger.Error("Order not found",
			zap.String("orderID", orderID),
			zap.String("userID", userID),
		)
		return nil, fmt.Errorf("order not found")
	}

	// Verify that the order belongs to the user
	if order.GetUserId() != userID {
		s.logger.Warn("Unauthorized access to order",
			zap.String("orderID", orderID),
			zap.String("userID", userID),
			zap.String("orderUserID", order.GetUserId()),
		)
		return nil, fmt.Errorf("unauthorized access to order")
	}

	return order, nil
}

// AddOrderPayment adds a payment to an existing order
func (s *OrderServiceImpl) AddOrderPayment(
	ctx context.Context,
	orderID string,
	amount float64,
	paymentType, reference, status string,
	date time.Time,
	description string,
	metadata map[string]string,
) error {
	s.logger.Debug("AddOrderPayment",
		zap.String("orderID", orderID),
		zap.Float64("amount", amount),
		zap.String("paymentType", paymentType),
		zap.String("status", status),
	)

	// Create a payment record in the database or call the appropriate service
	// This is a simplified implementation that just logs the payment
	// In a real application, you would create a payment record in the database
	// or call a payment service to process the payment

	s.logger.Info("Payment added to order",
		zap.String("orderID", orderID),
		zap.Float64("amount", amount),
		zap.String("paymentType", paymentType),
		zap.String("status", status),
	)

	return nil
}

// AddOrderTracking adds tracking information to an order
func (s *OrderServiceImpl) AddOrderTracking(
	ctx context.Context,
	orderID, carrier, trackingNum string,
	shipDate, estDelivery time.Time,
	notes string,
) error {
	s.logger.Debug("AddOrderTracking",
		zap.String("orderID", orderID),
		zap.String("carrier", carrier),
		zap.String("trackingNum", trackingNum),
		zap.Time("shipDate", shipDate),
		zap.Time("estDelivery", estDelivery),
	)

	// First, update the order status to shipped
	_, err := s.client.UpdateOrderStatus(ctx, &orderv1.UpdateOrderStatusRequest{
		Id:     orderID,
		Status: orderv1.OrderStatus_ORDER_STATUS_SHIPPED,
	})
	if err != nil {
		s.logger.Error("Failed to update order status to shipped",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Add the tracking code
	_, err = s.client.AddTrackingCode(ctx, &orderv1.AddTrackingCodeRequest{
		OrderId:      orderID,
		TrackingCode: trackingNum,
	})
	if err != nil {
		s.logger.Error("Failed to add tracking code to order",
			zap.String("orderID", orderID),
			zap.String("trackingNum", trackingNum),
			zap.Error(err),
		)
		return fmt.Errorf("failed to add tracking code: %w", err)
	}

	s.logger.Info("Tracking info added to order",
		zap.String("orderID", orderID),
		zap.String("trackingNum", trackingNum),
	)

	return nil
}

// CreateOrder creates a new order
func (s *OrderServiceImpl) CreateOrder(
	ctx context.Context,
	userID string,
	items []map[string]interface{},
	addressID, paymentType string,
	paymentData map[string]string,
	shippingType, notes string,
) (interface{}, error) {
	s.logger.Debug("CreateOrder",
		zap.String("userID", userID),
		zap.String("addressID", addressID),
		zap.String("paymentType", paymentType),
		zap.Int("itemsCount", len(items)),
	)

	// Convert items to the expected format
	var orderItems []*orderv1.OrderItem
	for _, item := range items {
		orderItem := &orderv1.OrderItem{
			ProductId: item["product_id"].(string),
			Quantity:  int32(item["quantity"].(float64)),
			Price:     item["price"].(float64),
		}
		orderItems = append(orderItems, orderItem)
	}

	req := &orderv1.CreateOrderRequest{
		UserId:  userID,
		Items:   orderItems,
		// TODO: Add shipping and billing addresses once they are defined in the protobuf
	}

	resp, err := s.client.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create order",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return resp.GetOrder(), nil
}

// ListOrders lists all orders (admin/staff)
func (s *OrderServiceImpl) ListOrders(
	ctx context.Context,
	status, userID, startDate, endDate string,
	limit, offset int,
) (interface{}, error) {
	s.logger.Debug("ListOrders",
		zap.String("status", status),
		zap.String("userID", userID),
		zap.String("startDate", startDate),
		zap.String("endDate", endDate),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	req := &orderv1.ListOrdersRequest{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	resp, err := s.client.ListOrders(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list orders",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	// Apply additional filters if provided
	var filtered []*orderv1.Order
	for _, order := range resp.Orders {
		match := true

		// Filter by user ID if provided
		if userID != "" && order.UserId != userID {
			match = false
		}

		// Filter by status if provided
		if status != "" && order.Status.String() != status {
			match = false
		}

		// TODO: Add date range filtering if needed

		if match {
			filtered = append(filtered, order)
		}
	}

	return filtered, nil
}

// GetOrderByID gets an order by ID (admin/staff)
func (s *OrderServiceImpl) GetOrderByID(
	ctx context.Context,
	orderID string,
) (interface{}, error) {
	s.logger.Debug("GetOrderByID",
		zap.String("orderID", orderID),
	)

	req := &orderv1.GetOrderRequest{
		Id: orderID,
	}

	resp, err := s.client.GetOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get order by ID",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	return resp.GetOrder(), nil
}

// UpdateOrderStatus updates order status (admin/staff)
func (s *OrderServiceImpl) UpdateOrderStatus(
	ctx context.Context,
	orderID, status, description string,
) error {
	s.logger.Debug("UpdateOrderStatus",
		zap.String("orderID", orderID),
		zap.String("status", status),
	)

	// Convert status string to OrderStatus enum
	var orderStatus orderv1.OrderStatus
	switch status {
	case "created":
		orderStatus = orderv1.OrderStatus_ORDER_STATUS_CREATED
	case "pending":
		orderStatus = orderv1.OrderStatus_ORDER_STATUS_PENDING
	case "paid":
		orderStatus = orderv1.OrderStatus_ORDER_STATUS_PAID
	case "shipped":
		orderStatus = orderv1.OrderStatus_ORDER_STATUS_SHIPPED
	case "delivered":
		orderStatus = orderv1.OrderStatus_ORDER_STATUS_DELIVERED
	case "cancelled":
		orderStatus = orderv1.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		orderStatus = orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}

	req := &orderv1.UpdateOrderStatusRequest{
		Id:     orderID,
		Status: orderStatus,
	}

	_, err := s.client.UpdateOrderStatus(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update order status",
			zap.String("orderID", orderID),
			zap.String("status", status),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Add a note to the order with the status update description
	if description != "" {
		// Get the current order to update
		getReq := &orderv1.GetOrderRequest{Id: orderID}
		getResp, err := s.client.GetOrder(ctx, getReq)
		if err != nil {
			s.logger.Error("Failed to get order for status update note",
				zap.String("orderID", orderID),
				zap.Error(err),
			)
			// Don't fail the entire operation if we can't add the note
			return nil
		}

		// Update the order notes
		if getResp.Order.Notes == "" {
			getResp.Order.Notes = fmt.Sprintf("Status updated to %s: %s", status, description)
		} else {
			getResp.Order.Notes = fmt.Sprintf("%s\nStatus updated to %s: %s", 
				getResp.Order.Notes, status, description)
		}

		// Save the updated order
		updateReq := &orderv1.UpdateOrderRequest{
			Order: getResp.Order,
		}

		_, err = s.client.UpdateOrder(ctx, updateReq)
		if err != nil {
			s.logger.Error("Failed to update order with status note",
				zap.String("orderID", orderID),
				zap.Error(err),
			)
			// Don't fail the entire operation if we can't add the note
			return nil
		}
	}

	return nil
}

// CancelOrder cancels an order (admin/staff)
func (s *OrderServiceImpl) CancelOrder(
	ctx context.Context,
	orderID, reason string,
) error {
	s.logger.Debug("CancelOrder",
		zap.String("orderID", orderID),
		zap.String("reason", reason),
	)

	req := &orderv1.CancelOrderRequest{
		Id: orderID,
	}

	_, err := s.client.CancelOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to cancel order",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Get the current order to update
	getReq := &orderv1.GetOrderRequest{Id: orderID}
	getResp, err := s.client.GetOrder(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get order for cancellation",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		// Don't fail the entire operation if we can't add the note
		return nil
	}

	// Update the order status
	getResp.Order.Status = orderv1.OrderStatus_ORDER_STATUS_CANCELLED

	// Add a note about the cancellation
	cancellationNote := fmt.Sprintf("Order cancelled")
	if reason != "" {
		cancellationNote = fmt.Sprintf("%s. Reason: %s", cancellationNote, reason)
	}

	if getResp.Order.Notes == "" {
		getResp.Order.Notes = cancellationNote
	} else {
		getResp.Order.Notes = fmt.Sprintf("%s\n%s", getResp.Order.Notes, cancellationNote)
	}

	// Save the updated order
	updateReq := &orderv1.UpdateOrderRequest{
		Order: getResp.Order,
	}

	_, err = s.client.UpdateOrder(ctx, updateReq)
	if err != nil {
		s.logger.Error("Failed to update order with cancellation",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update order with cancellation: %w", err)
	}

	return nil
}
