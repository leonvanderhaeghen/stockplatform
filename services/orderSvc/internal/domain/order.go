package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// OrderSourceType represents the source of an order
type OrderSourceType string

const (
	// SourceOnline represents an order placed online
	SourceOnline OrderSourceType = "ONLINE"
	// SourcePOS represents an order placed at a point of sale terminal
	SourcePOS OrderSourceType = "POS"
	// SourceMobile represents an order placed through a mobile app
	SourceMobile OrderSourceType = "MOBILE"
	// SourceAPI represents an order placed through API integration
	SourceAPI OrderSourceType = "API"
)

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	// StatusCreated represents a newly created order
	StatusCreated OrderStatus = "CREATED"
	// StatusPending represents an order that is being processed
	StatusPending OrderStatus = "PENDING"
	// StatusPaid represents an order that has been paid
	StatusPaid OrderStatus = "PAID"
	// StatusShipped represents an order that has been shipped
	StatusShipped OrderStatus = "SHIPPED"
	// StatusDelivered represents an order that has been delivered
	StatusDelivered OrderStatus = "DELIVERED"
	// StatusCancelled represents an order that has been cancelled
	StatusCancelled OrderStatus = "CANCELLED"
	// StatusFailed represents an order that has failed processing
	StatusFailed OrderStatus = "FAILED"
)

// ValidateStatusTransition checks if a status transition is valid
func ValidateStatusTransition(current, new OrderStatus) error {
	// Define valid status transitions
	validTransitions := map[OrderStatus][]OrderStatus{
		StatusCreated:   {StatusPending, StatusPaid, StatusCancelled, StatusFailed},
		StatusPending:   {StatusPaid, StatusCancelled, StatusFailed},
		StatusPaid:      {StatusShipped, StatusCancelled},
		StatusShipped:   {StatusDelivered, StatusFailed},
		StatusDelivered: {}, // Terminal state
		StatusCancelled: {}, // Terminal state
		StatusFailed:    {StatusPending}, // Can retry failed orders
	}

	allowed, exists := validTransitions[current]
	if !exists {
		return errors.New("invalid current status")
	}

	for _, allowedStatus := range allowed {
		if allowedStatus == new {
			return nil
		}
	}

	return errors.New("invalid status transition from " + string(current) + " to " + string(new))
}

// IsTerminalStatus returns true if the status is a terminal state
func IsTerminalStatus(status OrderStatus) bool {
	return status == StatusDelivered || status == StatusCancelled
}

// OrderItem represents a product in an order
type OrderItem struct {
	ProductID  string  `bson:"product_id"`
	ProductSKU string  `bson:"product_sku"`
	Name       string  `bson:"name"`
	Quantity   int32   `bson:"quantity"`
	Price      float64 `bson:"price"`
	Subtotal   float64 `bson:"subtotal"`
}

// Address represents a shipping or billing address
type Address struct {
	Street     string `bson:"street"`
	City       string `bson:"city"`
	State      string `bson:"state"`
	PostalCode string `bson:"postal_code"`
	Country    string `bson:"country"`
}

// Payment represents payment information for an order
type Payment struct {
	Method        string    `bson:"method"`
	TransactionID string    `bson:"transaction_id,omitempty"`
	Amount        float64   `bson:"amount"`
	Status        string    `bson:"status"`
	Timestamp     time.Time `bson:"timestamp,omitempty"`
}

// Order represents a customer order
type Order struct {
	ID            string          `bson:"_id,omitempty"`
	UserID        string          `bson:"user_id"`
	Items         []OrderItem     `bson:"items"`
	TotalAmount   float64         `bson:"total_amount"`
	Status        OrderStatus     `bson:"status"`
	Source        OrderSourceType `bson:"source"`
	ShippingAddr  Address         `bson:"shipping_address"`
	BillingAddr   Address         `bson:"billing_address"`
	Payment       Payment         `bson:"payment,omitempty"`
	TrackingCode  string          `bson:"tracking_code,omitempty"`
	Notes         string          `bson:"notes,omitempty"`
	Version       int32           `bson:"version"`           // For optimistic locking
	CreatedAt     time.Time       `bson:"created_at"`
	UpdatedAt     time.Time       `bson:"updated_at"`
	CompletedAt   time.Time       `bson:"completed_at,omitempty"`
	LocationID    string          `bson:"location_id,omitempty"` // Store location for POS orders
	StaffID       string          `bson:"staff_id,omitempty"`    // Staff member who processed the POS order
}

// NewOrder creates a new order
func NewOrder(userID string, items []OrderItem, shippingAddr, billingAddr Address) *Order {
	return NewOrderWithSource(userID, items, shippingAddr, billingAddr, SourceOnline, "", "")
}

// NewOrderWithSource creates a new order with specified source
func NewOrderWithSource(
	userID string,
	items []OrderItem,
	shippingAddr, billingAddr Address,
	source OrderSourceType,
	locationID string,
	staffID string,
) *Order {
	now := time.Now()
	order := &Order{
		ID:           uuid.New().String(),
		UserID:       userID,
		Items:        items,
		TotalAmount:  calculateTotal(items),
		Status:       StatusCreated,
		Source:       source,
		ShippingAddr: shippingAddr,
		BillingAddr:  billingAddr,
		Version:      1, // Initialize version for optimistic locking
		CreatedAt:    now,
		UpdatedAt:    now,
		LocationID:   locationID,
		StaffID:      staffID,
	}
	return order
}

// CalculateTotal calculates the total amount for the order
func calculateTotal(items []OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Subtotal
	}
	return total
}

// IncrementVersion increments the version for optimistic locking
func (o *Order) IncrementVersion() {
	o.Version++
	o.UpdatedAt = time.Now()
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) error {
	if err := ValidateStatusTransition(o.Status, status); err != nil {
		return err
	}
	o.Status = status
	o.IncrementVersion()
	if status == StatusDelivered {
		o.CompletedAt = time.Now()
	}
	return nil
}

// AddPayment adds payment information to the order
func (o *Order) AddPayment(method string, transactionID string, amount float64) error {
	if err := ValidateStatusTransition(o.Status, StatusPaid); err != nil {
		return err
	}
	o.Payment = Payment{
		Method:        method,
		TransactionID: transactionID,
		Amount:        amount,
		Status:        "COMPLETED",
		Timestamp:     time.Now(),
	}
	if err := o.UpdateStatus(StatusPaid); err != nil {
		return err
	}
	return nil
}

// IsPOSOrder returns true if this is a Point of Sale order
func (o *Order) IsPOSOrder() bool {
	return o.Source == SourcePOS
}

// SetPOSInfo sets Point of Sale specific information
func (o *Order) SetPOSInfo(locationID, staffID string) {
	o.Source = SourcePOS
	o.LocationID = locationID
	o.StaffID = staffID
	o.IncrementVersion()
}

// AddTrackingCode adds a tracking code to the order
func (o *Order) AddTrackingCode(code string) error {
	if err := ValidateStatusTransition(o.Status, StatusShipped); err != nil {
		return err
	}
	o.TrackingCode = code
	if err := o.UpdateStatus(StatusShipped); err != nil {
		return err
	}
	return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
	if err := ValidateStatusTransition(o.Status, StatusCancelled); err != nil {
		return err
	}
	if err := o.UpdateStatus(StatusCancelled); err != nil {
		return err
	}
	return nil
}

// Recalculate recalculates the order total
func (o *Order) Recalculate() {
	o.TotalAmount = calculateTotal(o.Items)
	o.IncrementVersion()
}

// AddNote adds a note to the order
func (o *Order) AddNote(note string) {
	if o.Notes != "" {
		o.Notes += "\n"
	}
	o.Notes += note
	o.IncrementVersion()
}
