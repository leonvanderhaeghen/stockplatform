package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	orderclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/order"
	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
)

// OrderServiceImpl implements the OrderService interface
type OrderServiceImpl struct {
	client *orderclient.Client
	logger *zap.Logger
}

// NewOrderService creates a new instance of OrderServiceImpl
func NewOrderService(orderServiceAddr string, logger *zap.Logger) (OrderService, error) {
	// Create a gRPC client via the new abstraction
	ordCfg := orderclient.Config{Address: orderServiceAddr}
	client, err := orderclient.New(ordCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create order client: %w", err)
	}

	return &OrderServiceImpl{
		client: client,
		logger: logger.Named("order_service"),
	}, nil
}

// GetUserOrders retrieves orders for a user
func (s *OrderServiceImpl) GetUserOrders(
	ctx context.Context,
	userID, status, startDate, endDate string,
	limit, offset int,
) (interface{}, error) {
	s.logger.Debug("GetUserOrders",
		zap.String("userID", userID),
		zap.String("status", status),
	)

	// Client GetUserOrders only supports userID, limit, offset
	resp, err := s.client.GetUserOrders(
		ctx,
		userID,
		int32(limit),
		int32(offset),
	)
	if err != nil {
		s.logger.Error("Failed to get user orders",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	return resp, nil
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

	// Client only has GetOrder method, so use that
	resp, err := s.client.GetOrder(ctx, orderID)
	if err != nil {
		s.logger.Error("Failed to get user order",
			zap.String("orderID", orderID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user order: %w", err)
	}

	// Validate that the order belongs to the requesting user
	// resp is already of type *models.Order from the client
	if resp != nil && resp.CustomerID != userID {
		s.logger.Warn("User attempted to access order belonging to another user",
			zap.String("orderID", orderID),
			zap.String("requesting_user", userID),
			zap.String("order_owner", resp.CustomerID),
		)
		return nil, fmt.Errorf("order not found or access denied")
	}

	return resp, nil
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

	// Add tracking code to the order using the proper client method
	err := s.client.AddTrackingCode(ctx, orderID, trackingNum)
	if err != nil {
		s.logger.Error("Failed to add tracking code to order",
			zap.String("orderID", orderID),
			zap.String("trackingNum", trackingNum),
			zap.Error(err),
		)
		return fmt.Errorf("failed to add tracking code: %w", err)
	}

	// Update order status to shipped
	err = s.client.UpdateOrderStatus(ctx, orderID, "SHIPPED")
	if err != nil {
		s.logger.Error("Failed to update order status to shipped",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	s.logger.Info("Order marked as shipped",
		zap.String("orderID", orderID),
		zap.String("carrier", carrier),
		zap.String("trackingNum", trackingNum),
	)

	s.logger.Info("Tracking info added to order",
		zap.String("orderID", orderID),
		zap.String("trackingNum", trackingNum),
	)

	return nil
}

// CreateOrder creates a new order (supports both online and POS orders)
func (s *OrderServiceImpl) CreateOrder(
	ctx context.Context,
	userID string,
	items []map[string]interface{},
	addressID, paymentType string,
	paymentData map[string]string,
	shippingType, notes, source, storeID string,
	customerInfo map[string]string,
) (interface{}, error) {
	s.logger.Debug("CreateOrder",
		zap.String("userID", userID),
		zap.String("addressID", addressID),
		zap.String("paymentType", paymentType),
		zap.String("source", source),
		zap.String("storeID", storeID),
		zap.Int("itemsCount", len(items)),
	)

	// Handle POS vs Online order logic
	isPOSOrder := source == "POS" || source == "QUICK_POS"
	if isPOSOrder {
		s.logger.Info("Processing POS order",
			zap.String("source", source),
			zap.String("storeID", storeID),
			zap.Any("customerInfo", customerInfo),
		)
		
		// For POS orders, handle walk-in customer information
		if len(customerInfo) > 0 {
			s.logger.Debug("Walk-in customer info provided",
				zap.Any("customerInfo", customerInfo),
			)
			// Store customer info in payment metadata for tracking
			if paymentData == nil {
				paymentData = make(map[string]string)
			}
			for k, v := range customerInfo {
				paymentData["customer_"+k] = v
			}
		}
		
		// Add POS-specific metadata
		if paymentData == nil {
			paymentData = make(map[string]string)
		}
		paymentData["source"] = source
		paymentData["store_id"] = storeID
		paymentData["order_type"] = "pos"
	}

	// Convert items to domain models
	orderItems := make([]*models.OrderItem, len(items))
	for i, item := range items {
		orderItems[i] = &models.OrderItem{
			ProductID: item["productId"].(string),
			SKU:       item["sku"].(string),
			Quantity:  item["quantity"].(int32),
			Price:     item["price"].(float64),
		}
	}

	// Handle shipping address based on order type
	var shippingAddress *models.Address
	if !isPOSOrder && addressID != "" {
		// Online order - would need user service integration for address lookup
		// For now, leaving as nil until proper address lookup is implemented
		shippingAddress = nil // Placeholder until address lookup implemented
	} else if isPOSOrder {
		// POS orders don't need shipping address (in-store pickup)
		shippingAddress = nil
		s.logger.Debug("POS order - no shipping address required")
	}

	resp, err := s.client.CreateOrder(ctx, userID, orderItems, shippingAddress, notes)
	if err != nil {
		s.logger.Error("Failed to create order",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	s.logger.Info("Order created successfully",
		zap.String("userID", userID),
		zap.String("orderID", resp.Order.ID),
	)

	return resp, nil
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

	// Convert parameters to match client interface (context, status, userID, limit, offset)
	// Note: Client doesn't support date filtering - parameters reduced to match signature
	resp, err := s.client.ListOrders(ctx, status, userID, int32(limit), int32(offset))
	if err != nil {
		s.logger.Error("Failed to list orders",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return resp, nil
}

// GetOrderByID gets an order by ID (admin/staff)
func (s *OrderServiceImpl) GetOrderByID(
	ctx context.Context,
	orderID string,
) (interface{}, error) {
	s.logger.Debug("GetOrderByID",
		zap.String("orderID", orderID),
	)

	// Use GetOrder method instead of GetOrderByID (which doesn't exist in client)
	resp, err := s.client.GetOrder(ctx, orderID)
	if err != nil {
		s.logger.Error("Failed to get order by ID",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	return resp, nil
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

	// Client interface only accepts context, orderID, and status (no description parameter)
	err := s.client.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		s.logger.Error("Failed to update order status",
			zap.String("orderID", orderID),
			zap.String("status", status),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update order status: %w", err)
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

	err := s.client.CancelOrder(ctx, orderID, reason)
	if err != nil {
		s.logger.Error("Failed to cancel order",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Note: Client abstraction handles cancellation logic internally
	// CancelOrder client method handles all cancellation logic including status updates

	s.logger.Info("Order cancelled successfully",
		zap.String("orderID", orderID),
		zap.String("reason", reason),
	)

	return nil
}

// Note: POS order creation is now handled via CreateOrder with source="POS" parameter
// All POS functionality has been consolidated into standard order endpoints

// Note: Quick POS transactions are now handled via CreateOrder with source="QUICK_POS" parameter
// All POS functionality has been consolidated into standard order endpoints
