package application

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// OrderService handles business logic for order operations
type OrderService struct {
	repo         domain.OrderRepository
	eventService *EventService
	logger       *zap.Logger
}

// NewOrderService creates a new order service
func NewOrderService(repo domain.OrderRepository, eventService *EventService, logger *zap.Logger) *OrderService {
	return &OrderService{
		repo:         repo,
		eventService: eventService,
		logger:       logger.Named("order_service"),
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

	// Publish order created event
	if s.eventService != nil {
		if err := s.eventService.PublishOrderCreated(ctx, order); err != nil {
			s.logger.Warn("Failed to publish order created event", zap.Error(err))
			// Don't fail the operation, just log the warning
		}

		// Publish inventory reservation event
		if err := s.eventService.PublishInventoryReserved(ctx, order); err != nil {
			s.logger.Warn("Failed to publish inventory reserved event", zap.Error(err))
		}
	}

	return order, nil
}

// CreatePOSOrder creates a new order from a Point of Sale terminal
func (s *OrderService) CreatePOSOrder(ctx context.Context, userID string, items []domain.OrderItem, 
	shippingAddr, billingAddr domain.Address, locationID, staffID string) (*domain.Order, error) {
	s.logger.Info("Creating POS order",
		zap.String("user_id", userID),
		zap.Int("item_count", len(items)),
		zap.String("location_id", locationID),
		zap.String("staff_id", staffID),
	)

	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	if locationID == "" {
		return nil, errors.New("location ID is required for POS orders")
	}

	order := domain.NewOrderWithSource(userID, items, shippingAddr, billingAddr, domain.SourcePOS, locationID, staffID)
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Publish order created event
	if s.eventService != nil {
		if err := s.eventService.PublishOrderCreated(ctx, order); err != nil {
			s.logger.Warn("Failed to publish order created event", zap.Error(err))
		}

		// Publish inventory reservation event
		if err := s.eventService.PublishInventoryReserved(ctx, order); err != nil {
			s.logger.Warn("Failed to publish inventory reserved event", zap.Error(err))
		}
	}

	return order, nil
}

// ProcessQuickPOSTransaction processes a quick checkout at a POS terminal (create order and add payment in one operation)
func (s *OrderService) ProcessQuickPOSTransaction(ctx context.Context, userID string, items []domain.OrderItem, 
	shippingAddr, billingAddr domain.Address, locationID, staffID string, 
	paymentMethod, paymentTransactionID string, paymentAmount float64) (*domain.Order, error) {
	
	s.logger.Info("Processing quick POS transaction",
		zap.String("user_id", userID),
		zap.Int("item_count", len(items)),
		zap.String("location_id", locationID),
		zap.String("payment_method", paymentMethod),
	)

	// Create the POS order
	order, err := s.CreatePOSOrder(ctx, userID, items, shippingAddr, billingAddr, locationID, staffID)
	if err != nil {
		return nil, err
	}

	// Add payment immediately
	order.AddPayment(paymentMethod, paymentTransactionID, paymentAmount)
	
	// Update the order with payment information
	if err := s.repo.Update(ctx, order); err != nil {
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
	
	// Ensure updated timestamp is set and increment version
	order.IncrementVersion()
	
	// Use optimistic locking for concurrent updates
	expectedVersion := order.Version - 1 // Version was incremented by IncrementVersion
	return s.repo.UpdateWithOptimisticLock(ctx, order, expectedVersion)
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
	
	previousStatus := order.Status
	err = order.UpdateStatus(status)
	if err != nil {
		return err
	}
	
	// Use optimistic locking for concurrent updates
	expectedVersion := order.Version - 1 // Version was incremented by UpdateStatus
	err = s.repo.UpdateWithOptimisticLock(ctx, order, expectedVersion)
	if err != nil {
		return err
	}
	
	// Publish status change event
	if s.eventService != nil {
		if err := s.eventService.PublishOrderStatusChanged(ctx, order, previousStatus); err != nil {
			s.logger.Warn("Failed to publish order status changed event", zap.Error(err))
		}
	}
	
	return nil
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
	
	err = order.AddPayment(method, transactionID, amount)
	if err != nil {
		return err
	}
	
	// Use optimistic locking for concurrent updates
	expectedVersion := order.Version - 1 // Version was incremented by AddPayment
	err = s.repo.UpdateWithOptimisticLock(ctx, order, expectedVersion)
	if err != nil {
		return err
	}
	
	// Publish payment processed event
	if s.eventService != nil {
		if err := s.eventService.PublishPaymentProcessed(ctx, order); err != nil {
			s.logger.Warn("Failed to publish payment processed event", zap.Error(err))
		}
	}
	
	return nil
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
	
	err = order.AddTrackingCode(trackingCode)
	if err != nil {
		return err
	}
	
	// Use optimistic locking for concurrent updates
	expectedVersion := order.Version - 1 // Version was incremented by AddTrackingCode
	err = s.repo.UpdateWithOptimisticLock(ctx, order, expectedVersion)
	if err != nil {
		return err
	}
	
	// Publish tracking code added event
	if s.eventService != nil {
		if err := s.eventService.PublishOrderStatusChanged(ctx, order, order.Status); err != nil {
			s.logger.Warn("Failed to publish order status changed event", zap.Error(err))
		}
	}
	
	return nil
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
	
	previousStatus := order.Status
	err = order.Cancel()
	if err != nil {
		return err
	}
	
	// Use optimistic locking for concurrent updates
	expectedVersion := order.Version - 1 // Version was incremented by Cancel
	err = s.repo.UpdateWithOptimisticLock(ctx, order, expectedVersion)
	if err != nil {
		return err
	}
	
	// Publish order cancellation events
	if s.eventService != nil {
		// Publish status changed event
		if err := s.eventService.PublishOrderStatusChanged(ctx, order, previousStatus); err != nil {
			s.logger.Warn("Failed to publish order status changed event", zap.Error(err))
		}
		
		// Publish inventory release event
		if err := s.eventService.PublishInventoryReleased(ctx, order); err != nil {
			s.logger.Warn("Failed to publish inventory released event", zap.Error(err))
		}
	}
	
	return nil
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
