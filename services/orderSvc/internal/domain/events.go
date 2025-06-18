package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of order event
type EventType string

const (
	// Order lifecycle events
	EventOrderCreated   EventType = "order.created"
	EventOrderPaid      EventType = "order.paid"
	EventOrderShipped   EventType = "order.shipped"
	EventOrderDelivered EventType = "order.delivered"
	EventOrderCancelled EventType = "order.cancelled"
	EventOrderFailed    EventType = "order.failed"
	
	// Order status event
	EventOrderStatusChanged EventType = "order.status_changed"

	// Payment events
	EventPaymentProcessed EventType = "payment.processed"
	EventPaymentFailed    EventType = "payment.failed"
	
	// Inventory events
	EventInventoryReserved EventType = "inventory.reserved"
	EventInventoryReleased EventType = "inventory.released"
)

// OrderEvent represents an event in the order lifecycle
type OrderEvent struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	OrderID     string                 `json:"order_id"`
	UserID      string                 `json:"user_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     int32                  `json:"version"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Metadata    map[string]string      `json:"metadata,omitempty"`
}

// NewOrderEvent creates a new order event
func NewOrderEvent(eventType EventType, orderID, userID string, version int32, data map[string]interface{}) *OrderEvent {
	return &OrderEvent{
		ID:        generateEventID(),
		Type:      eventType,
		OrderID:   orderID,
		UserID:    userID,
		Timestamp: time.Now(),
		Version:   version,
		Data:      data,
		Metadata:  make(map[string]string),
	}
}

// ToJSON serializes the event to JSON
func (e *OrderEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON deserializes an event from JSON
func FromJSON(data []byte) (*OrderEvent, error) {
	var event OrderEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// AddMetadata adds metadata to the event
func (e *OrderEvent) AddMetadata(key, value string) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
}

// generateEventID generates a unique event ID
func generateEventID() string {
	// Using UUID for event IDs
	return uuid.New().String()
}

// EventPublisher defines the interface for publishing order events
type EventPublisher interface {
	// PublishOrderEvent publishes an order lifecycle event
	PublishOrderEvent(event *OrderEvent) error
	
	// PublishInventoryEvent publishes an inventory-related event
	PublishInventoryEvent(event *OrderEvent) error
	
	// PublishPaymentEvent publishes a payment-related event
	PublishPaymentEvent(event *OrderEvent) error
}

// EventConsumer defines the interface for consuming order events
type EventConsumer interface {
	// ConsumeOrderEvents starts consuming order events
	ConsumeOrderEvents(handler OrderEventHandler) error
	
	// ConsumeInventoryEvents starts consuming inventory events
	ConsumeInventoryEvents(handler OrderEventHandler) error
	
	// ConsumePaymentEvents starts consuming payment events
	ConsumePaymentEvents(handler OrderEventHandler) error
	
	// Close closes the consumer
	Close() error
}

// OrderEventHandler defines the interface for handling order events
type OrderEventHandler interface {
	HandleEvent(event *OrderEvent) error
}
