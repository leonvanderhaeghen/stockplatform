package domain

import (
	"time"

	"github.com/google/uuid"
)

// InventoryItem represents a product's inventory information
type InventoryItem struct {
	ID          string    `bson:"_id,omitempty"`
	ProductID   string    `bson:"product_id"`
	Quantity    int32     `bson:"quantity"`
	Reserved    int32     `bson:"reserved"`
	SKU         string    `bson:"sku"`
	Location    string    `bson:"location,omitempty"`
	LastUpdated time.Time `bson:"last_updated"`
	CreatedAt   time.Time `bson:"created_at"`
}

// NewInventoryItem creates a new inventory item
func NewInventoryItem(productID string, quantity int32, sku string, location string) *InventoryItem {
	now := time.Now()
	return &InventoryItem{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Quantity:    quantity,
		Reserved:    0,
		SKU:         sku,
		Location:    location,
		LastUpdated: now,
		CreatedAt:   now,
	}
}

// IsAvailable checks if there's enough quantity available
func (i *InventoryItem) IsAvailable(requestedQuantity int32) bool {
	return (i.Quantity - i.Reserved) >= requestedQuantity
}

// Reserve attempts to reserve the specified quantity
// Returns true if successful, false if not enough inventory
func (i *InventoryItem) Reserve(quantity int32) bool {
	if !i.IsAvailable(quantity) {
		return false
	}
	i.Reserved += quantity
	i.LastUpdated = time.Now()
	return true
}

// ReleaseReservation releases a reservation
func (i *InventoryItem) ReleaseReservation(quantity int32) {
	if quantity > i.Reserved {
		i.Reserved = 0
	} else {
		i.Reserved -= quantity
	}
	i.LastUpdated = time.Now()
}

// AddStock adds stock to inventory
func (i *InventoryItem) AddStock(quantity int32) {
	i.Quantity += quantity
	i.LastUpdated = time.Now()
}

// RemoveStock removes stock from inventory
// Returns true if successful, false if not enough inventory
func (i *InventoryItem) RemoveStock(quantity int32) bool {
	if quantity > i.Quantity {
		return false
	}
	i.Quantity -= quantity
	i.LastUpdated = time.Now()
	return true
}

// FulfillReservation converts a reservation to a completed transaction
// Returns true if successful, false if not enough reserved
func (i *InventoryItem) FulfillReservation(quantity int32) bool {
	if quantity > i.Reserved {
		return false
	}
	i.Reserved -= quantity
	i.Quantity -= quantity
	i.LastUpdated = time.Now()
	return true
}
