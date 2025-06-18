package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// OrderEventHandler handles incoming order events
type OrderEventHandler struct {
	orderRepo domain.OrderRepository
	logger    *zap.Logger
}

// NewOrderEventHandler creates a new order event handler
func NewOrderEventHandler(orderRepo domain.OrderRepository, logger *zap.Logger) *OrderEventHandler {
	return &OrderEventHandler{
		orderRepo: orderRepo,
		logger:    logger.Named("order_event_handler"),
	}
}

// HandleEvent processes an order event
func (h *OrderEventHandler) HandleEvent(ctx context.Context, eventData []byte) error {
	var event domain.OrderEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		h.logger.Error("Failed to unmarshal order event", zap.Error(err))
		return fmt.Errorf("failed to unmarshal order event: %w", err)
	}

	h.logger.Info("Processing order event",
		zap.String("event_type", string(event.Type)),
		zap.String("order_id", event.OrderID),
		zap.String("user_id", event.UserID),
		zap.Int32("version", event.Version),
	)

	switch event.Type {
	case domain.EventOrderCreated:
		return h.handleOrderCreated(ctx, event)
	case domain.EventOrderStatusChanged:
		return h.handleOrderStatusChanged(ctx, event)
	case domain.EventPaymentProcessed:
		return h.handlePaymentProcessed(ctx, event)
	case domain.EventInventoryReserved:
		return h.handleInventoryReserved(ctx, event)
	case domain.EventInventoryReleased:
		return h.handleInventoryReleased(ctx, event)
	default:
		h.logger.Warn("Unknown event type", zap.String("event_type", string(event.Type)))
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

// handleOrderCreated processes order created events
func (h *OrderEventHandler) handleOrderCreated(ctx context.Context, event domain.OrderEvent) error {
	h.logger.Info("Handling order created event",
		zap.String("order_id", event.OrderID),
		zap.String("user_id", event.UserID),
	)

	// Could trigger additional processing like:
	// - Send order confirmation email
	// - Update analytics/reporting
	// - Trigger fulfillment process
	// - Update customer loyalty points

	return nil
}

// handleOrderStatusChanged processes order status change events
func (h *OrderEventHandler) handleOrderStatusChanged(ctx context.Context, event domain.OrderEvent) error {
	h.logger.Info("Handling order status changed event",
		zap.String("order_id", event.OrderID),
		zap.Any("data", event.Data),
	)

	// Could trigger additional processing like:
	// - Send status update notifications
	// - Update external systems
	// - Trigger shipping processes
	// - Update customer facing status

	return nil
}

// handlePaymentProcessed processes payment processed events
func (h *OrderEventHandler) handlePaymentProcessed(ctx context.Context, event domain.OrderEvent) error {
	h.logger.Info("Handling payment processed event",
		zap.String("order_id", event.OrderID),
		zap.Any("data", event.Data),
	)

	// Could trigger additional processing like:
	// - Send payment confirmation
	// - Update financial records
	// - Trigger order fulfillment
	// - Update fraud detection systems

	return nil
}

// handleInventoryReserved processes inventory reserved events
func (h *OrderEventHandler) handleInventoryReserved(ctx context.Context, event domain.OrderEvent) error {
	h.logger.Info("Handling inventory reserved event",
		zap.String("order_id", event.OrderID),
		zap.Any("data", event.Data),
	)

	// Could trigger additional processing like:
	// - Update inventory service
	// - Reserve stock in warehouse
	// - Update availability displays
	// - Trigger reorder processes

	return nil
}

// handleInventoryReleased processes inventory released events
func (h *OrderEventHandler) handleInventoryReleased(ctx context.Context, event domain.OrderEvent) error {
	h.logger.Info("Handling inventory released event",
		zap.String("order_id", event.OrderID),
		zap.Any("data", event.Data),
	)

	// Could trigger additional processing like:
	// - Release reserved stock
	// - Update inventory availability
	// - Trigger restock notifications
	// - Update warehouse systems

	return nil
}
