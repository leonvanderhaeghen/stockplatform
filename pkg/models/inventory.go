package models

import "time"

// InventoryItem represents an inventory item in the domain
type InventoryItem struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"product_id"`
	SKU        string    `json:"sku"`
	LocationID string    `json:"location_id"`
	Quantity   int32     `json:"quantity"`
	Reserved   int32     `json:"reserved"`
	Available  int32     `json:"available"`
	ReorderAt  int32     `json:"reorder_at"`
	ReorderQty int32     `json:"reorder_qty"`
	Cost       float64   `json:"cost"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CheckAvailabilityResponse represents availability check results
type CheckAvailabilityResponse struct {
	Available bool                       `json:"available"`
	Items     []InventoryAvailabilityItem `json:"items"`
}

// InventoryAvailabilityItem represents item availability info
type InventoryAvailabilityItem struct {
	ProductID     string `json:"product_id"`
	SKU           string `json:"sku"`
	RequestedQty  int32  `json:"requested_qty"`
	AvailableQty  int32  `json:"available_qty"`
	Available     bool   `json:"available"`
}

// InventoryRequestItem represents a request for inventory availability
type InventoryRequestItem struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Quantity  int32  `json:"quantity"`
}
