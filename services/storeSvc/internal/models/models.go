package models

import "time"

// Store represents a physical store location
type Store struct {
	ID            string            `bson:"_id" json:"id"`
	Name          string            `bson:"name" json:"name"`
	Description   string            `bson:"description" json:"description"`
	Address       Address           `bson:"address" json:"address"`
	Phone         string            `bson:"phone" json:"phone"`
	Email         string            `bson:"email" json:"email"`
	ManagerUserID string            `bson:"manager_user_id" json:"manager_user_id"`
	IsActive      bool              `bson:"is_active" json:"is_active"`
	Hours         StoreHours        `bson:"hours" json:"hours"`
	Metadata      map[string]string `bson:"metadata" json:"metadata"`
	CreatedAt     time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time         `bson:"updated_at" json:"updated_at"`
}

// Address represents a physical address
type Address struct {
	Street     string  `bson:"street" json:"street"`
	City       string  `bson:"city" json:"city"`
	State      string  `bson:"state" json:"state"`
	PostalCode string  `bson:"postal_code" json:"postal_code"`
	Country    string  `bson:"country" json:"country"`
	Latitude   float64 `bson:"latitude" json:"latitude"`
	Longitude  float64 `bson:"longitude" json:"longitude"`
}

// StoreHours represents operating hours for a store
type StoreHours struct {
	Days []DayHours `bson:"days" json:"days"`
}

// DayHours represents hours for a specific day
type DayHours struct {
	Day       string `bson:"day" json:"day"`         // Monday, Tuesday, etc.
	OpenTime  string `bson:"open_time" json:"open_time"`   // HH:MM format
	CloseTime string `bson:"close_time" json:"close_time"` // HH:MM format
	IsClosed  bool   `bson:"is_closed" json:"is_closed"`   // True if store is closed on this day
}

// StoreProduct represents a product available in a specific store
type StoreProduct struct {
	StoreID           string    `bson:"store_id" json:"store_id"`
	ProductID         string    `bson:"product_id" json:"product_id"`
	StockQuantity     int32     `bson:"stock_quantity" json:"stock_quantity"`
	ReservedQuantity  int32     `bson:"reserved_quantity" json:"reserved_quantity"`
	AvailableQuantity int32     `bson:"available_quantity" json:"available_quantity"`
	StorePrice        string    `bson:"store_price" json:"store_price"` // Store-specific pricing (optional)
	IsAvailable       bool      `bson:"is_available" json:"is_available"`
	LastUpdated       time.Time `bson:"last_updated" json:"last_updated"`
}

// ProductReservation represents a reserved product
type ProductReservation struct {
	ID          string    `bson:"_id" json:"id"`
	StoreID     string    `bson:"store_id" json:"store_id"`
	ProductID   string    `bson:"product_id" json:"product_id"`
	UserID      string    `bson:"user_id" json:"user_id"` // Customer who made the reservation
	Quantity    int32     `bson:"quantity" json:"quantity"`
	Status      string    `bson:"status" json:"status"` // ACTIVE, EXPIRED, COMPLETED, CANCELLED
	ReservedAt  time.Time `bson:"reserved_at" json:"reserved_at"`
	ExpiresAt   time.Time `bson:"expires_at" json:"expires_at"`
	CompletedAt time.Time `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	Notes       string    `bson:"notes" json:"notes"`
}

// StoreUser represents the relationship between a user and a store
type StoreUser struct {
	StoreID    string    `bson:"store_id" json:"store_id"`
	UserID     string    `bson:"user_id" json:"user_id"`
	Role       string    `bson:"role" json:"role"` // EMPLOYEE, MANAGER, ADMIN
	AssignedAt time.Time `bson:"assigned_at" json:"assigned_at"`
}

// StoreSale represents a sale made at a physical store
type StoreSale struct {
	ID             string            `bson:"_id" json:"id"`
	StoreID        string            `bson:"store_id" json:"store_id"`
	OrderID        string            `bson:"order_id,omitempty" json:"order_id,omitempty"` // Link to order service
	SalesUserID    string            `bson:"sales_user_id" json:"sales_user_id"`           // Employee who made the sale
	CustomerUserID string            `bson:"customer_user_id,omitempty" json:"customer_user_id,omitempty"` // Customer (optional)
	Items          []StoreSaleItem   `bson:"items" json:"items"`
	TotalAmount    string            `bson:"total_amount" json:"total_amount"`
	Currency       string            `bson:"currency" json:"currency"`
	SaleType       string            `bson:"sale_type" json:"sale_type"` // WALK_IN, RESERVATION, ONLINE_PICKUP
	SaleDate       time.Time         `bson:"sale_date" json:"sale_date"`
	ReservationID  string            `bson:"reservation_id,omitempty" json:"reservation_id,omitempty"`
	Metadata       map[string]string `bson:"metadata" json:"metadata"`
}

// StoreSaleItem represents an item in a store sale
type StoreSaleItem struct {
	ProductID   string `bson:"product_id" json:"product_id"`
	ProductName string `bson:"product_name" json:"product_name"`
	ProductSKU  string `bson:"product_sku" json:"product_sku"`
	Quantity    int32  `bson:"quantity" json:"quantity"`
	UnitPrice   string `bson:"unit_price" json:"unit_price"`
	Subtotal    string `bson:"subtotal" json:"subtotal"`
}
