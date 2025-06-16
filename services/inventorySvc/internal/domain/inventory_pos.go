package domain

import (
	"fmt"
	"time"
)

// InventoryCheckItem represents an item to check in inventory
type InventoryCheckItem struct {
	ProductID string `json:"product_id,omitempty"`
	SKU       string `json:"sku,omitempty"`
	Quantity  int32  `json:"quantity"`
}

// InventoryCheckResult represents the result of an inventory check
type InventoryCheckResult struct {
	ProductID         string `json:"product_id,omitempty"`
	SKU               string `json:"sku,omitempty"`
	Quantity          int32  `json:"quantity"`
	Available         bool   `json:"available"`
	AvailableQuantity int32  `json:"available_quantity"`
	Location          string `json:"location,omitempty"`
	ShelfLocation     string `json:"shelf_location,omitempty"`
	Error             string `json:"error,omitempty"`
}

// ReservationItem represents an item to be reserved in inventory
type ReservationItem struct {
	ProductID string `json:"product_id,omitempty"`
	SKU       string `json:"sku,omitempty"`
	Quantity  int32  `json:"quantity"`
}

// AddNote adds a note to the inventory item with timestamp
func (i *InventoryItem) AddNote(note string) {
	if i.ReservationNotes != "" {
		i.ReservationNotes += "\n"
	}
	timestamp := time.Now().Format(time.RFC3339)
	i.ReservationNotes += fmt.Sprintf("[%s] %s", timestamp, note)
	i.LastUpdated = time.Now()
}

// ReserveForOrder reserves inventory for a specific order ID
// Returns true if successful, false if not enough inventory
func (i *InventoryItem) ReserveForOrder(quantity int32, orderID string) bool {
	if !i.IsAvailable(quantity) {
		return false
	}
	
	i.Reserved += quantity
	i.ReservationStatus = ReservationStatusActive
	i.OrderID = orderID
	i.LastUpdated = time.Now()
	
	// Add a note about the reservation
	i.AddNote(fmt.Sprintf("Reserved %d units for order %s", quantity, orderID))
	
	return true
}

// GetAvailable returns the available quantity (total minus reserved)
func (i *InventoryItem) GetAvailable() int32 {
	return i.Quantity - i.Reserved
}

// Available is a shorthand property to get available quantity
// This is added so the code can use inventoryItem.Available
func (i *InventoryItem) Available() int32 {
	return i.GetAvailable()
}

// HasSufficientStock checks if there's enough quantity available for the requested amount
func (i *InventoryItem) HasSufficientStock(requestedQuantity int32) bool {
	return i.GetAvailable() >= requestedQuantity
}
