package application

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// OrderService handles business logic for order operations
type OrderService struct {
	repo   domain.OrderRepository
	logger *zap.Logger
}

// NewOrderService creates a new order service
func NewOrderService(repo domain.OrderRepository, logger *zap.Logger) *OrderService {
	return &OrderService{
		repo:   repo,
		logger: logger.Named("order_service"),
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem, shippingAddr, billingAddr domain.Address) (*domain.Order, error) {
	s.logger.Info("Creating order",
		zap.String("user_id", userID),
		zap.Int("item_count", len(items)),
	)

	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	order := domain.NewOrder(userID, items, shippingAddr, billingAddr)
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *OrderService) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	s.logger.Debug("Getting order", zap.String("id", id))
	
	if id == "" {
		return nil, errors.New("order ID is required")
	}
	
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if order == nil {
		return nil, errors.New("order not found")
	}
	
	return order, nil
}

// GetUserOrders retrieves orders for a specific user
func (s *OrderService) GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*domain.Order, error) {
	s.logger.Debug("Getting user orders", 
		zap.String("user_id", userID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	if limit <= 0 {
		limit = 10 // Default limit
	}
	
	return s.repo.GetByUserID(ctx, userID, limit, offset)
}

// UpdateOrder updates an existing order
func (s *OrderService) UpdateOrder(ctx context.Context, order *domain.Order) error {
	s.logger.Info("Updating order",
		zap.String("id", order.ID),
		zap.String("user_id", order.UserID),
		zap.String("status", string(order.Status)),
	)
	
	if order.ID == "" {
		return errors.New("order ID is required")
	}
	
	// Ensure updated timestamp is set
	order.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, order)
}

// DeleteOrder removes an order
func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	s.logger.Info("Deleting order", zap.String("id", id))
	
	if id == "" {
		return errors.New("order ID is required")
	}
	
	return s.repo.Delete(ctx, id)
}

// ListOrders returns all orders with optional filtering and pagination
func (s *OrderService) ListOrders(ctx context.Context, status string, limit, offset int) ([]*domain.Order, error) {
	s.logger.Debug("Listing orders",
		zap.String("status", status),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	if limit <= 0 {
		limit = 10 // Default limit
	}
	
	filter := make(map[string]interface{})
	if status != "" {
		filter["status"] = status
	}
	
	return s.repo.List(ctx, filter, limit, offset)
}

// UpdateOrderStatus updates the status of an order
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	s.logger.Info("Updating order status",
		zap.String("id", orderID),
		zap.String("status", string(status)),
	)
	
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	order.UpdateStatus(status)
	return s.repo.Update(ctx, order)
}

// AddPaymentToOrder adds payment information to an order
func (s *OrderService) AddPaymentToOrder(ctx context.Context, orderID, method, transactionID string, amount float64) error {
	s.logger.Info("Adding payment to order",
		zap.String("id", orderID),
		zap.String("method", method),
		zap.String("transaction_id", transactionID),
		zap.Float64("amount", amount),
	)
	
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	order.AddPayment(method, transactionID, amount)
	return s.repo.Update(ctx, order)
}

// AddTrackingCodeToOrder adds a tracking code to an order
func (s *OrderService) AddTrackingCodeToOrder(ctx context.Context, orderID, trackingCode string) error {
	s.logger.Info("Adding tracking code to order",
		zap.String("id", orderID),
		zap.String("tracking_code", trackingCode),
	)
	
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	order.AddTrackingCode(trackingCode)
	return s.repo.Update(ctx, order)
}

// CancelOrder cancels an order
func (s *OrderService) CancelOrder(ctx context.Context, orderID string) error {
	s.logger.Info("Cancelling order", zap.String("id", orderID))
	
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	
	if order == nil {
		return errors.New("order not found")
	}
	
	order.Cancel()
	return s.repo.Update(ctx, order)
}

// CountOrdersByStatus counts orders with a specific status
func (s *OrderService) CountOrdersByStatus(ctx context.Context, status string) (int64, error) {
	s.logger.Debug("Counting orders by status", zap.String("status", status))
	
	filter := make(map[string]interface{})
	if status != "" {
		filter["status"] = status
	}
	
	return s.repo.Count(ctx, filter)
}
