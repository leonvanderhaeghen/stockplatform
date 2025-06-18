package application

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// EventService handles order event publishing and processing
type EventService struct {
	publisher domain.EventPublisher
	logger    *zap.Logger
}

// NewEventService creates a new event service
func NewEventService(publisher domain.EventPublisher, logger *zap.Logger) *EventService {
	return &EventService{
		publisher: publisher,
		logger:    logger.Named("event_service"),
	}
}

// PublishOrderCreated publishes an order created event
func (s *EventService) PublishOrderCreated(ctx context.Context, order *domain.Order) error {
	event := domain.NewOrderEvent(
		domain.EventOrderCreated,
		order.ID,
		order.UserID,
		order.Version,
		map[string]interface{}{
			"total_amount": order.TotalAmount,
			"items_count":  len(order.Items),
			"source":       string(order.Source),
			"location_id":  order.LocationID,
			"staff_id":     order.StaffID,
		},
	)

	err := s.publisher.PublishOrderEvent(event)
	if err != nil {
		s.logger.Error("Failed to publish order created event",
			zap.Error(err),
			zap.String("order_id", order.ID),
			zap.String("user_id", order.UserID),
		)
		return fmt.Errorf("failed to publish order created event: %w", err)
	}

	s.logger.Info("Published order created event",
		zap.String("order_id", order.ID),
		zap.String("user_id", order.UserID),
	)

	return nil
}

// PublishOrderStatusChanged publishes an order status change event
func (s *EventService) PublishOrderStatusChanged(ctx context.Context, order *domain.Order, previousStatus domain.OrderStatus) error {
	var eventType domain.EventType
	
	switch order.Status {
	case domain.StatusPaid:
		eventType = domain.EventOrderPaid
	case domain.StatusShipped:
		eventType = domain.EventOrderShipped
	case domain.StatusDelivered:
		eventType = domain.EventOrderDelivered
	case domain.StatusCancelled:
		eventType = domain.EventOrderCancelled
	case domain.StatusFailed:
		eventType = domain.EventOrderFailed
	default:
		// For other status changes, use a generic event type
		return s.publishGenericStatusChange(ctx, order, previousStatus)
	}

	event := domain.NewOrderEvent(
		eventType,
		order.ID,
		order.UserID,
		order.Version,
		map[string]interface{}{
			"previous_status": string(previousStatus),
			"new_status":      string(order.Status),
			"total_amount":    order.TotalAmount,
			"tracking_code":   order.TrackingCode,
		},
	)

	err := s.publisher.PublishOrderEvent(event)
	if err != nil {
		s.logger.Error("Failed to publish order status changed event",
			zap.Error(err),
			zap.String("order_id", order.ID),
			zap.String("previous_status", string(previousStatus)),
			zap.String("new_status", string(order.Status)),
		)
		return fmt.Errorf("failed to publish order status changed event: %w", err)
	}

	s.logger.Info("Published order status changed event",
		zap.String("order_id", order.ID),
		zap.String("event_type", string(eventType)),
		zap.String("previous_status", string(previousStatus)),
		zap.String("new_status", string(order.Status)),
	)

	return nil
}

// PublishPaymentProcessed publishes a payment processed event
func (s *EventService) PublishPaymentProcessed(ctx context.Context, order *domain.Order) error {
	event := domain.NewOrderEvent(
		domain.EventPaymentProcessed,
		order.ID,
		order.UserID,
		order.Version,
		map[string]interface{}{
			"payment_method":      order.Payment.Method,
			"payment_amount":      order.Payment.Amount,
			"transaction_id":      order.Payment.TransactionID,
			"payment_status":      order.Payment.Status,
			"payment_timestamp":   order.Payment.Timestamp,
		},
	)

	err := s.publisher.PublishPaymentEvent(event)
	if err != nil {
		s.logger.Error("Failed to publish payment processed event",
			zap.Error(err),
			zap.String("order_id", order.ID),
			zap.String("transaction_id", order.Payment.TransactionID),
		)
		return fmt.Errorf("failed to publish payment processed event: %w", err)
	}

	s.logger.Info("Published payment processed event",
		zap.String("order_id", order.ID),
		zap.String("transaction_id", order.Payment.TransactionID),
		zap.Float64("amount", order.Payment.Amount),
	)

	return nil
}

// PublishInventoryReserved publishes an inventory reserved event
func (s *EventService) PublishInventoryReserved(ctx context.Context, order *domain.Order) error {
	// Create inventory reservation data
	var reservations []map[string]interface{}
	for _, item := range order.Items {
		reservations = append(reservations, map[string]interface{}{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
			"price":      item.Price,
		})
	}

	event := domain.NewOrderEvent(
		domain.EventInventoryReserved,
		order.ID,
		order.UserID,
		order.Version,
		map[string]interface{}{
			"reservations": reservations,
			"total_amount": order.TotalAmount,
		},
	)

	err := s.publisher.PublishInventoryEvent(event)
	if err != nil {
		s.logger.Error("Failed to publish inventory reserved event",
			zap.Error(err),
			zap.String("order_id", order.ID),
			zap.Int("items_count", len(order.Items)),
		)
		return fmt.Errorf("failed to publish inventory reserved event: %w", err)
	}

	s.logger.Info("Published inventory reserved event",
		zap.String("order_id", order.ID),
		zap.Int("items_count", len(order.Items)),
	)

	return nil
}

// PublishInventoryReleased publishes an inventory released event
func (s *EventService) PublishInventoryReleased(ctx context.Context, order *domain.Order) error {
	// Create inventory release data
	var releases []map[string]interface{}
	for _, item := range order.Items {
		releases = append(releases, map[string]interface{}{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
		})
	}

	event := domain.NewOrderEvent(
		domain.EventInventoryReleased,
		order.ID,
		order.UserID,
		order.Version,
		map[string]interface{}{
			"releases": releases,
			"reason":   "order_cancelled",
		},
	)

	err := s.publisher.PublishInventoryEvent(event)
	if err != nil {
		s.logger.Error("Failed to publish inventory released event",
			zap.Error(err),
			zap.String("order_id", order.ID),
			zap.Int("items_count", len(order.Items)),
		)
		return fmt.Errorf("failed to publish inventory released event: %w", err)
	}

	s.logger.Info("Published inventory released event",
		zap.String("order_id", order.ID),
		zap.Int("items_count", len(order.Items)),
	)

	return nil
}

// publishGenericStatusChange publishes a generic status change event
func (s *EventService) publishGenericStatusChange(ctx context.Context, order *domain.Order, previousStatus domain.OrderStatus) error {
	event := domain.NewOrderEvent(
		"order.status_changed",
		order.ID,
		order.UserID,
		order.Version,
		map[string]interface{}{
			"previous_status": string(previousStatus),
			"new_status":      string(order.Status),
		},
	)

	return s.publisher.PublishOrderEvent(event)
}
