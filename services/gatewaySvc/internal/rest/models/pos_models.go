package models

import "time"

// POSTransactionRequest represents a request for POS transaction
type POSTransactionRequest struct {
	TransactionType  string                 `json:"transaction_type" validate:"required,oneof=sale return exchange"`
	LocationID       string                 `json:"location_id" validate:"required"`
	StaffID          string                 `json:"staff_id" validate:"required"`
	ReferenceOrderID string                 `json:"reference_order_id,omitempty"`
	Items            []POSTransactionItem   `json:"items" validate:"required,min=1,dive"`
	PaymentInfo      POSTransactionPayment  `json:"payment_info" validate:"required"`
}

// POSTransactionItem represents an item in a POS transaction
type POSTransactionItem struct {
	ProductID       string  `json:"product_id,omitempty"`
	SKU             string  `json:"sku,omitempty" validate:"required_without=ProductID"`
	Quantity        int32   `json:"quantity" validate:"required,gt=0"`
	Price           float32 `json:"price" validate:"required"`
	Description     string  `json:"description,omitempty"`
	Reason          string  `json:"reason,omitempty"`
	Direction       string  `json:"direction,omitempty" validate:"omitempty,oneof=in out"`
	InventoryItemID string  `json:"inventory_item_id,omitempty"`
}

// POSTransactionPayment represents payment information for a POS transaction
type POSTransactionPayment struct {
	Method           string  `json:"method" validate:"required,oneof=cash credit debit gift_card store_credit"`
	Amount           float32 `json:"amount" validate:"required,gt=0"`
	CurrencyCode     string  `json:"currency_code" validate:"required,len=3"`
	PaymentReference string  `json:"payment_reference,omitempty"`
	CardLast4        string  `json:"card_last_4,omitempty"`
}

// POSTransactionResponse represents a response for a POS transaction
type POSTransactionResponse struct {
	Success        bool                     `json:"success"`
	TransactionID  string                   `json:"transaction_id"`
	OrderID        string                   `json:"order_id,omitempty"`
	TotalAmount    float32                  `json:"total_amount"`
	CurrencyCode   string                   `json:"currency_code"`
	ProcessedItems []POSTransactionItemResp `json:"processed_items"`
	Message        string                   `json:"message,omitempty"`
	Errors         []string                 `json:"errors,omitempty"`
}

// POSTransactionItemResp represents a processed item in a POS transaction response
type POSTransactionItemResp struct {
	ProductID       string  `json:"product_id,omitempty"`
	SKU             string  `json:"sku,omitempty"`
	Quantity        int32   `json:"quantity"`
	ProcessedPrice  float32 `json:"processed_price"`
	InventoryItemID string  `json:"inventory_item_id,omitempty"`
	Success         bool    `json:"success"`
	Message         string  `json:"message,omitempty"`
}

// CheckAvailabilityRequest represents a request to check inventory availability
type CheckAvailabilityRequest struct {
	LocationID string                `json:"location_id" validate:"required"`
	Items      []InventoryCheckItem  `json:"items" validate:"required,min=1,dive"`
}

// InventoryCheckItem represents an item in an availability check
type InventoryCheckItem struct {
	ProductID string `json:"product_id,omitempty"`
	SKU       string `json:"sku,omitempty" validate:"required_without=ProductID"`
	Quantity  int32  `json:"quantity" validate:"required,gt=0"`
}

// CheckAvailabilityResponse represents a response for an availability check
type CheckAvailabilityResponse struct {
	AllAvailable bool                     `json:"all_available"`
	ItemResults  []InventoryAvailability  `json:"item_results"`
}

// InventoryAvailability represents availability information for an item
type InventoryAvailability struct {
	ProductID       string `json:"product_id,omitempty"`
	SKU             string `json:"sku,omitempty"`
	InventoryItemID string `json:"inventory_item_id,omitempty"`
	Requested       int32  `json:"requested"`
	Available       int32  `json:"available"`
	IsAvailable     bool   `json:"is_available"`
}

// NearbyInventoryRequest represents a request to find nearby inventory
type NearbyInventoryRequest struct {
	ProductID string  `json:"product_id,omitempty"`
	SKU       string  `json:"sku,omitempty" validate:"required_without=ProductID"`
	Quantity  int32   `json:"quantity" validate:"required,gt=0"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	RadiusKm  int32   `json:"radius_km" validate:"required,gt=0"`
}

// NearbyInventoryResponse represents a response for nearby inventory
type NearbyInventoryResponse struct {
	Locations []LocationInventory `json:"locations"`
}

// LocationInventory represents inventory information for a location
type LocationInventory struct {
	LocationID   string  `json:"location_id"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	Distance     float32 `json:"distance"`
	Available    int32   `json:"available"`
	InStock      bool    `json:"in_stock"`
	PickupReady  bool    `json:"pickup_ready"`
	PickupTime   string  `json:"pickup_time,omitempty"`
	OpeningHours string  `json:"opening_hours,omitempty"`
}

// PickupReservationRequest represents a request to reserve items for pickup
type PickupReservationRequest struct {
	UserID     string              `json:"user_id" validate:"required"`
	LocationID string              `json:"location_id" validate:"required"`
	Items      []ReservationItem   `json:"items" validate:"required,min=1,dive"`
	PickupTime time.Time           `json:"pickup_time" validate:"required"`
	Notes      string              `json:"notes,omitempty"`
}

// ReservationItem represents an item in a pickup reservation
type ReservationItem struct {
	ProductID       string `json:"product_id,omitempty"`
	SKU             string `json:"sku,omitempty" validate:"required_without=ProductID"`
	Quantity        int32  `json:"quantity" validate:"required,gt=0"`
	InventoryItemID string `json:"inventory_item_id,omitempty"`
}

// PickupReservationResponse represents a response for a pickup reservation
type PickupReservationResponse struct {
	Success       bool                        `json:"success"`
	ReservationID string                      `json:"reservation_id"`
	ExpiresAt     time.Time                   `json:"expires_at"`
	Items         []ReservationItemResult     `json:"items"`
	Errors        []string                    `json:"errors,omitempty"`
}

// ReservationItemResult represents a result for a reservation item
type ReservationItemResult struct {
	ProductID       string `json:"product_id,omitempty"`
	SKU             string `json:"sku,omitempty"`
	InventoryItemID string `json:"inventory_item_id,omitempty"`
	Quantity        int32  `json:"quantity"`
	Reserved        bool   `json:"reserved"`
	Message         string `json:"message,omitempty"`
}

// CompletePickupRequest represents a request to complete a pickup
type CompletePickupRequest struct {
	ReservationID string `json:"reservation_id" validate:"required"`
	StaffID       string `json:"staff_id" validate:"required"`
	Notes         string `json:"notes,omitempty"`
}

// CancelPickupRequest represents a request to cancel a pickup
type CancelPickupRequest struct {
	ReservationID string `json:"reservation_id" validate:"required"`
	Reason        string `json:"reason" validate:"required"`
}

// PickupActionResponse represents a response for a pickup action
type PickupActionResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}
