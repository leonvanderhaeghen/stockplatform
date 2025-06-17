package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/leonvanderhaeghen/stockplatform/services/orderSvc/api/gen/go/proto/order/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// OrderServer implements the gRPC interface for order service
type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer
	service              *application.OrderService
	posTransactionService *application.POSTransactionService
	logger               *zap.Logger
}

// NewOrderServer creates a new order gRPC server
func NewOrderServer(service *application.OrderService, posService *application.POSTransactionService, logger *zap.Logger) orderv1.OrderServiceServer {
	return &OrderServer{
		service:              service,
		posTransactionService: posService,
		logger:               logger.Named("order_grpc_server"),
	}
}

// CreateOrder creates a new order
func (s *OrderServer) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	s.logger.Info("gRPC CreateOrder called",
		zap.String("user_id", req.UserId),
		zap.Int("item_count", len(req.Items)),
	)

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one item is required")
	}

	// Convert proto items to domain items
	items := make([]domain.OrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.OrderItem{
			ProductID:  item.ProductId,
			ProductSKU: item.ProductSku,
			Name:       item.Name,
			Quantity:   item.Quantity,
			Price:      item.Price,
			Subtotal:   item.Subtotal,
		})
	}

	// Convert proto addresses to domain addresses
	shippingAddr := domain.Address{
		Street:     req.ShippingAddress.Street,
		City:       req.ShippingAddress.City,
		State:      req.ShippingAddress.State,
		PostalCode: req.ShippingAddress.PostalCode,
		Country:    req.ShippingAddress.Country,
	}

	billingAddr := domain.Address{
		Street:     req.BillingAddress.Street,
		City:       req.BillingAddress.City,
		State:      req.BillingAddress.State,
		PostalCode: req.BillingAddress.PostalCode,
		Country:    req.BillingAddress.Country,
	}

	order, err := s.service.CreateOrder(ctx, req.UserId, items, shippingAddr, billingAddr)
	if err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create order: "+err.Error())
	}

	return &orderv1.CreateOrderResponse{
		Order: toProtoOrder(order),
	}, nil
}

// GetOrder retrieves an order by ID
func (s *OrderServer) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	s.logger.Debug("gRPC GetOrder called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	order, err := s.service.GetOrder(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get order", zap.Error(err))
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return &orderv1.GetOrderResponse{
		Order: toProtoOrder(order),
	}, nil
}

// GetUserOrders retrieves orders for a specific user
func (s *OrderServer) GetUserOrders(ctx context.Context, req *orderv1.GetUserOrdersRequest) (*orderv1.GetUserOrdersResponse, error) {
	s.logger.Debug("gRPC GetUserOrders called",
		zap.String("user_id", req.UserId),
		zap.Int32("limit", req.Limit),
		zap.Int32("offset", req.Offset),
	)

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10 // Default limit
	}

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	orders, err := s.service.GetUserOrders(ctx, req.UserId, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get user orders", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get user orders: "+err.Error())
	}

	// Convert domain orders to proto orders
	protoOrders := make([]*orderv1.Order, 0, len(orders))
	for _, order := range orders {
		protoOrders = append(protoOrders, toProtoOrder(order))
	}

	return &orderv1.GetUserOrdersResponse{
		Orders: protoOrders,
	}, nil
}

// UpdateOrder updates an existing order
func (s *OrderServer) UpdateOrder(ctx context.Context, req *orderv1.UpdateOrderRequest) (*orderv1.UpdateOrderResponse, error) {
	s.logger.Info("gRPC UpdateOrder called", zap.String("id", req.Order.Id))

	if req.Order == nil {
		return nil, status.Error(codes.InvalidArgument, "order is required")
	}
	if req.Order.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "order.id is required")
	}

	// First get the existing order
	existingOrder, err := s.service.GetOrder(ctx, req.Order.Id)
	if err != nil {
		s.logger.Error("Failed to get order for update", zap.Error(err))
		return nil, status.Error(codes.NotFound, "order not found")
	}

	// Update order fields
	existingOrder.Status = domain.OrderStatus(req.Order.Status.String())
	existingOrder.Notes = req.Order.Notes
	existingOrder.TrackingCode = req.Order.TrackingCode

	// Update addresses if provided
	if req.Order.ShippingAddress != nil {
		existingOrder.ShippingAddr = domain.Address{
			Street:     req.Order.ShippingAddress.Street,
			City:       req.Order.ShippingAddress.City,
			State:      req.Order.ShippingAddress.State,
			PostalCode: req.Order.ShippingAddress.PostalCode,
			Country:    req.Order.ShippingAddress.Country,
		}
	}

	if req.Order.BillingAddress != nil {
		existingOrder.BillingAddr = domain.Address{
			Street:     req.Order.BillingAddress.Street,
			City:       req.Order.BillingAddress.City,
			State:      req.Order.BillingAddress.State,
			PostalCode: req.Order.BillingAddress.PostalCode,
			Country:    req.Order.BillingAddress.Country,
		}
	}

	if err := s.service.UpdateOrder(ctx, existingOrder); err != nil {
		s.logger.Error("Failed to update order", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update order: "+err.Error())
	}

	return &orderv1.UpdateOrderResponse{
		Success: true,
	}, nil
}

// DeleteOrder removes an order
func (s *OrderServer) DeleteOrder(ctx context.Context, req *orderv1.DeleteOrderRequest) (*orderv1.DeleteOrderResponse, error) {
	s.logger.Info("gRPC DeleteOrder called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.service.DeleteOrder(ctx, req.Id); err != nil {
		s.logger.Error("Failed to delete order", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete order: "+err.Error())
	}

	return &orderv1.DeleteOrderResponse{
		Success: true,
	}, nil
}

// ListOrders lists all orders with optional filtering and pagination
func (s *OrderServer) ListOrders(ctx context.Context, req *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	s.logger.Debug("gRPC ListOrders called",
		zap.String("status", req.Status),
		zap.Int32("limit", req.Limit),
		zap.Int32("offset", req.Offset),
	)

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10 // Default limit
	}

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	orders, err := s.service.ListOrders(ctx, req.Status, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list orders", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list orders: "+err.Error())
	}

	// Convert domain orders to proto orders
	protoOrders := make([]*orderv1.Order, 0, len(orders))
	for _, order := range orders {
		protoOrders = append(protoOrders, toProtoOrder(order))
	}

	return &orderv1.ListOrdersResponse{
		Orders: protoOrders,
	}, nil
}

// UpdateOrderStatus updates the status of an order
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderv1.UpdateOrderStatusRequest) (*orderv1.UpdateOrderStatusResponse, error) {
	s.logger.Info("gRPC UpdateOrderStatus called",
		zap.String("id", req.Id),
		zap.String("status", req.Status.String()),
	)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Convert proto status to domain status
	var domainStatus domain.OrderStatus
	switch req.Status {
	case orderv1.OrderStatus_ORDER_STATUS_CREATED:
		domainStatus = domain.StatusCreated
	case orderv1.OrderStatus_ORDER_STATUS_PENDING:
		domainStatus = domain.StatusPending
	case orderv1.OrderStatus_ORDER_STATUS_PAID:
		domainStatus = domain.StatusPaid
	case orderv1.OrderStatus_ORDER_STATUS_SHIPPED:
		domainStatus = domain.StatusShipped
	case orderv1.OrderStatus_ORDER_STATUS_DELIVERED:
		domainStatus = domain.StatusDelivered
	case orderv1.OrderStatus_ORDER_STATUS_CANCELLED:
		domainStatus = domain.StatusCancelled
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}

	if err := s.service.UpdateOrderStatus(ctx, req.Id, domainStatus); err != nil {
		s.logger.Error("Failed to update order status", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update order status: "+err.Error())
	}

	return &orderv1.UpdateOrderStatusResponse{
		Success: true,
	}, nil
}

// AddPayment adds payment information to an order
func (s *OrderServer) AddPayment(ctx context.Context, req *orderv1.AddPaymentRequest) (*orderv1.AddPaymentResponse, error) {
	s.logger.Info("gRPC AddPayment called",
		zap.String("order_id", req.OrderId),
		zap.String("method", req.Method),
		zap.String("transaction_id", req.TransactionId),
		zap.Float64("amount", req.Amount),
	)

	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.Method == "" {
		return nil, status.Error(codes.InvalidArgument, "method is required")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be positive")
	}

	if err := s.service.AddPaymentToOrder(ctx, req.OrderId, req.Method, req.TransactionId, req.Amount); err != nil {
		s.logger.Error("Failed to add payment", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to add payment: "+err.Error())
	}

	return &orderv1.AddPaymentResponse{
		Success: true,
	}, nil
}

// AddTrackingCode adds a tracking code to an order
func (s *OrderServer) AddTrackingCode(ctx context.Context, req *orderv1.AddTrackingCodeRequest) (*orderv1.AddTrackingCodeResponse, error) {
	s.logger.Info("gRPC AddTrackingCode called",
		zap.String("order_id", req.OrderId),
		zap.String("tracking_code", req.TrackingCode),
	)

	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.TrackingCode == "" {
		return nil, status.Error(codes.InvalidArgument, "tracking_code is required")
	}

	if err := s.service.AddTrackingCodeToOrder(ctx, req.OrderId, req.TrackingCode); err != nil {
		s.logger.Error("Failed to add tracking code", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to add tracking code: "+err.Error())
	}

	return &orderv1.AddTrackingCodeResponse{
		Success: true,
	}, nil
}

// CancelOrder cancels an order
func (s *OrderServer) CancelOrder(ctx context.Context, req *orderv1.CancelOrderRequest) (*orderv1.CancelOrderResponse, error) {
	s.logger.Info("gRPC CancelOrder called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.service.CancelOrder(ctx, req.Id); err != nil {
		s.logger.Error("Failed to cancel order", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to cancel order: "+err.Error())
	}

	return &orderv1.CancelOrderResponse{
		Success: true,
	}, nil
}

// toProtoOrder converts a domain order to a proto order
func toProtoOrder(order *domain.Order) *orderv1.Order {
	protoOrder := &orderv1.Order{
		Id:           order.ID,
		UserId:       order.UserID,
		TotalAmount:  order.TotalAmount,
		Notes:        order.Notes,
		TrackingCode: order.TrackingCode,
		CreatedAt:    order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    order.UpdatedAt.Format(time.RFC3339),
	}

	// Convert status
	switch order.Status {
	case domain.StatusCreated:
		protoOrder.Status = orderv1.OrderStatus_ORDER_STATUS_CREATED
	case domain.StatusPending:
		protoOrder.Status = orderv1.OrderStatus_ORDER_STATUS_PENDING
	case domain.StatusPaid:
		protoOrder.Status = orderv1.OrderStatus_ORDER_STATUS_PAID
	case domain.StatusShipped:
		protoOrder.Status = orderv1.OrderStatus_ORDER_STATUS_SHIPPED
	case domain.StatusDelivered:
		protoOrder.Status = orderv1.OrderStatus_ORDER_STATUS_DELIVERED
	case domain.StatusCancelled:
		protoOrder.Status = orderv1.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		protoOrder.Status = orderv1.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}

	// Convert items
	protoOrder.Items = make([]*orderv1.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		protoOrder.Items = append(protoOrder.Items, &orderv1.OrderItem{
			ProductId:  item.ProductID,
			ProductSku: item.ProductSKU,
			Name:       item.Name,
			Quantity:   item.Quantity,
			Price:      item.Price,
			Subtotal:   item.Subtotal,
		})
	}

	// Convert addresses
	protoOrder.ShippingAddress = &orderv1.Address{
		Street:     order.ShippingAddr.Street,
		City:       order.ShippingAddr.City,
		State:      order.ShippingAddr.State,
		PostalCode: order.ShippingAddr.PostalCode,
		Country:    order.ShippingAddr.Country,
	}

	protoOrder.BillingAddress = &orderv1.Address{
		Street:     order.BillingAddr.Street,
		City:       order.BillingAddr.City,
		State:      order.BillingAddr.State,
		PostalCode: order.BillingAddr.PostalCode,
		Country:    order.BillingAddr.Country,
	}

	// Convert payment if it exists
	if order.Payment.Method != "" {
		protoOrder.Payment = &orderv1.Payment{
			Method:        order.Payment.Method,
			TransactionId: order.Payment.TransactionID,
			Amount:        order.Payment.Amount,
			Status:        order.Payment.Status,
		}
		if !order.Payment.Timestamp.IsZero() {
			protoOrder.Payment.Timestamp = order.Payment.Timestamp.Format(time.RFC3339)
		}
	}

	// Set completed at if it exists
	if !order.CompletedAt.IsZero() {
		protoOrder.CompletedAt = order.CompletedAt.Format(time.RFC3339)
	}

	return protoOrder
}
