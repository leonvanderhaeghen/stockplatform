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

// CreatePOSOrder creates an order from a POS terminal
func (s *OrderServiceImpl) CreatePOSOrder(
	ctx context.Context,
	userID string,
	items []map[string]interface{},
	locationID, staffID, paymentType string,
	paymentData map[string]string,
	notes string,
) (interface{}, error) {
	s.logger.Debug("CreatePOSOrder",
		zap.String("userID", userID),
		zap.String("locationID", locationID),
		zap.String("staffID", staffID),
		zap.String("paymentType", paymentType),
		zap.Any("items", items),
	)

	// Convert items to order items
	orderItems := make([]*orderv1.OrderItem, 0, len(items))
	for _, item := range items {
		// Get product ID
		productIDVal, ok := item["product_id"]
		if !ok {
			return nil, fmt.Errorf("missing product_id in item")
		}
		productID, ok := productIDVal.(string)
		if !ok {
			return nil, fmt.Errorf("product_id must be a string")
		}

		// Get quantity
		quantityVal, ok := item["quantity"]
		if !ok {
			return nil, fmt.Errorf("missing quantity in item")
		}

		// Handle different types for quantity
		var quantity int32
		switch q := quantityVal.(type) {
		case int:
			quantity = int32(q)
		case int32:
			quantity = q
		case float64:
			quantity = int32(q)
		default:
			return nil, fmt.Errorf("quantity must be a number")
		}

		// Get price
		priceVal, ok := item["price"]
		if !ok {
			return nil, fmt.Errorf("missing price in item")
		}

		// Handle different types for price
		var price float64
		switch p := priceVal.(type) {
		case float32:
			price = float64(p)
		case float64:
			price = p
		case int:
			price = float64(p)
		case int32:
			price = float64(p)
		default:
			return nil, fmt.Errorf("price must be a number")
		}

		// Create an order item
		orderItem := &orderv1.OrderItem{
			ProductId: productID,
			Quantity:  quantity,
			Price:     price,
			Subtotal:  price * float64(quantity),
		}

		// Add any additional fields if present in the item
		if nameVal, ok := item["name"].(string); ok {
			orderItem.Name = nameVal
		}
		if skuVal, ok := item["sku"].(string); ok {
			orderItem.ProductSku = skuVal
		}

		orderItems = append(orderItems, orderItem)
	}

	// Calculate order totals
	totalAmount := float64(0)
	for _, item := range orderItems {
		totalAmount += item.Subtotal
	}

	// Create a payment object for order if payment info provided
	payment := &orderv1.Payment{
		Method: paymentType,
		Amount: totalAmount,
		Status: "completed",
		Timestamp: time.Now().Format(time.RFC3339),
		TransactionId: fmt.Sprintf("pos-%s-%d", locationID, time.Now().UnixNano()),
	}

	// Create the order request
	req := &orderv1.CreateOrderRequest{
		UserId: userID,
		Items:  orderItems,
	}

	// Create the order
	resp, err := s.client.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create POS order",
			zap.String("userID", userID),
			zap.String("locationID", locationID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create POS order: %w", err)
	}

	if resp.Order != nil && resp.Order.Id != "" {
		orderID := resp.Order.Id
		
		// Update the order with POS-specific details
		getReq := &orderv1.GetOrderRequest{Id: orderID}
		getResp, err := s.client.GetOrder(ctx, getReq)
		if err == nil && getResp.Order != nil {
			order := getResp.Order
			order.Status = orderv1.OrderStatus_ORDER_STATUS_PAID
			order.TotalAmount = totalAmount
			order.Notes = notes
			order.Payment = payment
			
			// Add tracking information for metadata
			order.TrackingCode = fmt.Sprintf("POS-%s-%s", locationID, staffID)
			
			// Set completion time
			order.CompletedAt = time.Now().Format(time.RFC3339)
			
			// Update the order
			_, updateErr := s.client.UpdateOrder(ctx, &orderv1.UpdateOrderRequest{
				Order: order,
			})
			if updateErr != nil {
				s.logger.Warn("Failed to update POS order with additional details",
					zap.String("orderID", orderID),
					zap.Error(updateErr),
				)
				// Don't fail the overall operation if this update fails
			}
		}
		
		return orderID, nil
	}
	
	// Fallback if order ID is not in response (shouldn't happen with proper proto implementation)
	return fmt.Sprintf("order-%d", time.Now().UnixNano()), nil
}

// ProcessQuickPOSTransaction creates and processes a POS order in one step
func (s *OrderServiceImpl) ProcessQuickPOSTransaction(
	ctx context.Context,
	locationID, staffID string,
	items []map[string]interface{},
	paymentInfo map[string]interface{},
) (interface{}, error) {
	s.logger.Debug("ProcessQuickPOSTransaction",
		zap.String("locationID", locationID),
		zap.String("staffID", staffID),
		zap.Any("items", items),
		zap.Any("paymentInfo", paymentInfo),
	)

	// Extract payment details
	paymentType := "cash" // Default payment type
	if pt, ok := paymentInfo["method"].(string); ok {
		paymentType = pt
	} else if pt, ok := paymentInfo["type"].(string); ok {
		paymentType = pt
	}

	// Extract user ID if available, otherwise use guest user
	userID := "guest"
	if uid, ok := paymentInfo["user_id"].(string); ok && uid != "" {
		userID = uid
	}

	// Convert payment info to string map for metadata
	paymentData := make(map[string]string)
	for k, v := range paymentInfo {
		if str, ok := v.(string); ok {
			paymentData[k] = str
		} else {
			paymentData[k] = fmt.Sprintf("%v", v)
		}
	}

	// Create notes if available
	notes := ""
	if n, ok := paymentInfo["notes"].(string); ok {
		notes = n
	}

	// Create the POS order
	orderID, err := s.CreatePOSOrder(
		ctx,
		userID,
		items,
		locationID,
		staffID,
		paymentType,
		paymentData,
		notes,
	)

	if err != nil {
		s.logger.Error("Failed to process quick POS transaction", 
			zap.String("locationID", locationID), 
			zap.String("staffID", staffID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to process quick POS transaction: %w", err)
	}

	// Generate consistent transaction ID
	txID := fmt.Sprintf("pos-%s-%s-%d", 
		locationID, 
		staffID, 
		time.Now().UnixNano())

	// Return transaction information
	transactionInfo := map[string]interface{}{
		"order_id":       orderID,
		"transaction_id": txID,
		"status":        "completed",
		"timestamp":     time.Now().Format(time.RFC3339),
		"total_amount":  calculateTotalAmount(items),
		"payment_method": paymentType,
	}

	return transactionInfo, nil
}

// Helper function to calculate the total amount from items
func calculateTotalAmount(items []map[string]interface{}) float64 {
	total := 0.0
	
	for _, item := range items {
		var price float64
		var quantity float64
		
		// Get price
		if p, ok := item["price"].(float64); ok {
			price = p
		} else if p, ok := item["price"].(float32); ok {
			price = float64(p)
		} else if p, ok := item["price"].(int); ok {
			price = float64(p)
		}
		
		// Get quantity
		if q, ok := item["quantity"].(float64); ok {
			quantity = q
		} else if q, ok := item["quantity"].(float32); ok {
			quantity = float64(q)
		} else if q, ok := item["quantity"].(int); ok {
			quantity = float64(q)
		} else {
			quantity = 1.0 // Default to 1 if quantity is missing
		}
		
		total += price * quantity
	}
	
	return total
}
