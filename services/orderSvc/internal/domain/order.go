package domain

import (
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
)

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

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = status
	o.UpdatedAt = time.Now()
	if status == StatusDelivered {
		o.CompletedAt = time.Now()
	}
}

// AddPayment adds payment information to the order
func (o *Order) AddPayment(method string, transactionID string, amount float64) {
	o.Payment = Payment{
		Method:        method,
		TransactionID: transactionID,
		Amount:        amount,
		Status:        "COMPLETED",
		Timestamp:     time.Now(),
	}
	o.UpdateStatus(StatusPaid)
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
	o.UpdatedAt = time.Now()
}

// AddTrackingCode adds a tracking code to the order
func (o *Order) AddTrackingCode(code string) {
	o.TrackingCode = code
	o.UpdateStatus(StatusShipped)
}

// Cancel cancels the order
func (o *Order) Cancel() {
	o.UpdateStatus(StatusCancelled)
}

// Recalculate recalculates the order total
func (o *Order) Recalculate() {
	o.TotalAmount = calculateTotal(o.Items)
}

// AddNote adds a note to the order
func (o *Order) AddNote(note string) {
	if o.Notes != "" {
		o.Notes += "\n"
	}
	o.Notes += note
	o.UpdatedAt = time.Now()
}
