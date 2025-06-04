package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/order/v1"
)

// OrderServiceImpl implements the OrderService interface
type OrderServiceImpl struct {
	client orderv1.OrderServiceClient
	logger *zap.Logger
}

// NewOrderService creates a new instance of OrderServiceImpl
func NewOrderService(orderServiceAddr string, logger *zap.Logger) (OrderService, error) {
	// Create a gRPC connection to the order service
	conn, err := grpc.Dial(
		orderServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	// Create a client
	client := orderv1.NewOrderServiceClient(conn)

	return &OrderServiceImpl{
		client: client,
		logger: logger.Named("order_service"),
	}, nil
}

// GetUserOrders gets orders for a specific user
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

	return resp.Orders, nil
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

	// First get the order by ID
	req := &orderv1.GetOrderRequest{
		Id: orderID,
	}

	resp, err := s.client.GetOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get order",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify the order belongs to the user
	if resp.Order.UserId != userID {
		s.logger.Error("Order does not belong to user",
			zap.String("orderID", orderID),
			zap.String("userID", userID),
		)
		return nil, fmt.Errorf("order not found")
	}

	return resp.Order, nil
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
		zap.Int("itemsCount", len(items)),
		zap.String("addressID", addressID),
	)

	// Convert items to OrderItems
	orderItems := make([]*orderv1.OrderItem, 0, len(items))
	totalAmount := 0.0

	for _, item := range items {
		price, _ := item["price"].(float64)
		quantity, _ := item["quantity"].(float64)
		subtotal := price * quantity

		orderItem := &orderv1.OrderItem{
			ProductId:   item["product_id"].(string),
			ProductSku:  item["sku"].(string),
			Name:        item["name"].(string),
			Quantity:    int32(quantity),
			Price:       price,
			Subtotal:    subtotal,
		}
		orderItems = append(orderItems, orderItem)
		totalAmount += subtotal
	}

	// TODO: Get address details from user service
	// For now, create a placeholder address
	address := &orderv1.Address{
		Street:     "123 Main St",
		City:       "Anytown",
		State:      "CA",
		PostalCode: "12345",
		Country:    "USA",
	}

	req := &orderv1.CreateOrderRequest{
		UserId:          userID,
		Items:           orderItems,
		ShippingAddress: address,
		BillingAddress:  address,
	}

	// Add payment information if provided
	if paymentType != "" && paymentData != nil {
		// Convert payment data to a Payment message
		payment := &orderv1.Payment{
			Method:        paymentType,
			TransactionId: paymentData["transaction_id"],
			Amount:        totalAmount,
			Status:        "pending",
			Timestamp:     time.Now().Format(time.RFC3339),
		}

		// Create a new order with payment information
		order := &orderv1.Order{
			UserId:          userID,
			Items:           orderItems,
			TotalAmount:     totalAmount,
			Status:          orderv1.OrderStatus_ORDER_STATUS_CREATED,
			ShippingAddress: address,
			BillingAddress:  address,
			Payment:         payment,
			Notes:           notes,
			CreatedAt:       time.Now().Format(time.RFC3339),
		}

		// Use UpdateOrder instead of CreateOrder since we have payment info
		updateReq := &orderv1.UpdateOrderRequest{
			Order: order,
		}

		_, err := s.client.UpdateOrder(ctx, updateReq)
		if err != nil {
			s.logger.Error("Failed to update order with payment",
				zap.String("userID", userID),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to update order with payment: %w", err)
		}

		return order, nil
	}

	// Create order without payment information
	resp, err := s.client.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create order",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return resp.Order, nil
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

	// If a user ID was provided, filter the results
	if userID != "" {
		filteredOrders := make([]*orderv1.Order, 0)
		for _, order := range resp.Orders {
			if order.UserId == userID {
				filteredOrders = append(filteredOrders, order)
			}
		}
		return filteredOrders, nil
	}

	return resp.Orders, nil
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
		s.logger.Error("Failed to get order",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return resp.Order, nil
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

// AddOrderPayment adds payment to an order (admin/staff)
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
		zap.String("type", paymentType),
	)

	// First, get the current order
	getReq := &orderv1.GetOrderRequest{Id: orderID}
	getResp, err := s.client.GetOrder(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get order for adding payment",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Create a new payment
	payment := &orderv1.Payment{
		Method:        paymentType,
		TransactionId: reference,
		Amount:        amount,
		Status:        status,
		Timestamp:     date.Format(time.RFC3339),
	}

	// Update the order with the new payment
	order := getResp.Order
	order.Payment = payment

	// Update the order status if needed
	if status == "completed" || status == "paid" {
		order.Status = orderv1.OrderStatus_ORDER_STATUS_PAID
	}

	// Add a note about the payment
	if order.Notes == "" {
		order.Notes = fmt.Sprintf("Payment added: %s - %s - %.2f", paymentType, reference, amount)
	} else {
		order.Notes = fmt.Sprintf("%s\nPayment added: %s - %s - %.2f", 
			order.Notes, paymentType, reference, amount)
	}

	// Save the updated order
	updateReq := &orderv1.UpdateOrderRequest{
		Order: order,
	}

	_, err = s.client.UpdateOrder(ctx, updateReq)
	if err != nil {
		s.logger.Error("Failed to update order with payment",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update order with payment: %w", err)
	}

	return nil
}

// AddOrderTracking adds tracking info to an order (admin/staff)
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
	)

	// First, get the current order
	getReq := &orderv1.GetOrderRequest{Id: orderID}
	getResp, err := s.client.GetOrder(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get order for adding tracking",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Update the order with tracking information
	order := getResp.Order
	order.TrackingCode = trackingNum

	// Add a note about the tracking update
	trackingNote := fmt.Sprintf("Tracking added - Carrier: %s, Tracking #: %s", carrier, trackingNum)
	if notes != "" {
		trackingNote = fmt.Sprintf("%s\nNotes: %s", trackingNote, notes)
	}

	if order.Notes == "" {
		order.Notes = trackingNote
	} else {
		order.Notes = fmt.Sprintf("%s\n%s", order.Notes, trackingNote)
	}

	// Update the order status to shipped if it's not already
	if order.Status < orderv1.OrderStatus_ORDER_STATUS_SHIPPED {
		order.Status = orderv1.OrderStatus_ORDER_STATUS_SHIPPED
	}

	// Save the updated order
	updateReq := &orderv1.UpdateOrderRequest{
		Order: order,
	}

	_, err = s.client.UpdateOrder(ctx, updateReq)
	if err != nil {
		s.logger.Error("Failed to update order with tracking",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update order with tracking: %w", err)
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

	// First, get the current order
	getReq := &orderv1.GetOrderRequest{Id: orderID}
	getResp, err := s.client.GetOrder(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get order for cancellation",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Update the order status to cancelled
	order := getResp.Order
	order.Status = orderv1.OrderStatus_ORDER_STATUS_CANCELLED

	// Add a note about the cancellation
	cancellationNote := "Order cancelled"
	if reason != "" {
		cancellationNote = fmt.Sprintf("%s. Reason: %s", cancellationNote, reason)
	}

	if order.Notes == "" {
		order.Notes = cancellationNote
	} else {
		order.Notes = fmt.Sprintf("%s\n%s", order.Notes, cancellationNote)
	}

	// Save the updated order
	updateReq := &orderv1.UpdateOrderRequest{
		Order: order,
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
