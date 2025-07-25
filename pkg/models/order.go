package models

import "time"

// Order represents an order in the domain
type Order struct {
	ID          string      `json:"id"`
	CustomerID  string      `json:"customer_id"`
	Status      OrderStatus `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	Items       []*OrderItem `json:"items"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`
	BillingAddress  *Address `json:"billing_address,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	SKU       string  `json:"sku"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
	Total     float64 `json:"total"`
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Address represents a shipping or billing address
type Address struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
	Country  string `json:"country"`
}

// CreateOrderResponse represents the response from creating an order
type CreateOrderResponse struct {
	Order   *Order `json:"order"`
	Message string `json:"message"`
}

// ListOrdersResponse represents the response from listing orders
type ListOrdersResponse struct {
	Orders     []*Order `json:"orders"`
	TotalCount int32    `json:"total_count"`
}

// UpdateOrderResponse represents the response from updating an order
type UpdateOrderResponse struct {
	Order   *Order `json:"order"`
	Message string `json:"message"`
}
