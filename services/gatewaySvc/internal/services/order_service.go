package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	orderv1 "stockplatform/pkg/gen/order/v1"
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
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	if startDate != "" {
		startTime, err := time.Parse(time.RFC3339, startDate)
		if err == nil {
			req.StartDate = timestamppb.New(startTime)
		}
	}

	if endDate != "" {
		endTime, err := time.Parse(time.RFC3339, endDate)
		if err == nil {
			req.EndDate = timestamppb.New(endTime)
		}
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

	req := &orderv1.GetUserOrderRequest{
		OrderId: orderID,
		UserId:  userID,
	}

	resp, err := s.client.GetUserOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get user order",
			zap.String("orderID", orderID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user order: %w", err)
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
		zap.String("addressID", addressID),
		zap.String("paymentType", paymentType),
		zap.String("shippingType", shippingType),
		zap.Int("itemCount", len(items)),
	)

	// Convert items to proto format
	orderItems := make([]*orderv1.OrderItem, 0, len(items))
	for _, item := range items {
		orderItem := &orderv1.OrderItem{}
		
		if productID, ok := item["productId"].(string); ok {
			orderItem.ProductId = productID
		}
		
		if sku, ok := item["sku"].(string); ok {
			orderItem.Sku = sku
		}
		
		if quantity, ok := item["quantity"].(int32); ok {
			orderItem.Quantity = quantity
		} else if quantityFloat, ok := item["quantity"].(float64); ok {
			orderItem.Quantity = int32(quantityFloat)
		}
		
		if price, ok := item["price"].(float64); ok {
			orderItem.Price = price
		}
		
		orderItems = append(orderItems, orderItem)
	}

	req := &orderv1.CreateOrderRequest{
		UserId:       userID,
		Items:        orderItems,
		AddressId:    addressID,
		PaymentType:  paymentType,
		PaymentData:  paymentData,
		ShippingType: shippingType,
		Notes:        notes,
	}

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
		UserId: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	if startDate != "" {
		startTime, err := time.Parse(time.RFC3339, startDate)
		if err == nil {
			req.StartDate = timestamppb.New(startTime)
		}
	}

	if endDate != "" {
		endTime, err := time.Parse(time.RFC3339, endDate)
		if err == nil {
			req.EndDate = timestamppb.New(endTime)
		}
	}

	resp, err := s.client.ListOrders(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list orders",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list orders: %w", err)
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
		zap.String("description", description),
	)

	req := &orderv1.UpdateOrderStatusRequest{
		Id:          orderID,
		Status:      status,
		Description: description,
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
		zap.String("paymentType", paymentType),
		zap.String("status", status),
	)

	req := &orderv1.AddPaymentRequest{
		OrderId:     orderID,
		Amount:      amount,
		Type:        paymentType,
		Reference:   reference,
		Status:      status,
		Date:        timestamppb.New(date),
		Description: description,
		Metadata:    metadata,
	}

	_, err := s.client.AddPayment(ctx, req)
	if err != nil {
		s.logger.Error("Failed to add order payment",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to add order payment: %w", err)
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

	req := &orderv1.AddTrackingRequest{
		OrderId:     orderID,
		Carrier:     carrier,
		TrackingNum: trackingNum,
		Notes:       notes,
	}

	if !shipDate.IsZero() {
		req.ShipDate = timestamppb.New(shipDate)
	}

	if !estDelivery.IsZero() {
		req.EstDelivery = timestamppb.New(estDelivery)
	}

	_, err := s.client.AddTracking(ctx, req)
	if err != nil {
		s.logger.Error("Failed to add order tracking",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to add order tracking: %w", err)
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
		Id:     orderID,
		Reason: reason,
	}

	_, err := s.client.CancelOrder(ctx, req)
	if err != nil {
		s.logger.Error("Failed to cancel order",
			zap.String("orderID", orderID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}
