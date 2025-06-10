package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// StoreLocation represents a physical store or warehouse location
type StoreLocation struct {
	ID          string    `bson:"_id,omitempty"`
	Name        string    `bson:"name"`
	Type        string    `bson:"type"` // e.g., "store", "warehouse", "fulfillment_center"
	Address     string    `bson:"address"`
	City        string    `bson:"city"`
	State       string    `bson:"state"`
	PostalCode  string    `bson:"postal_code"`
	Country     string    `bson:"country"`
	IsActive    bool      `bson:"is_active"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

// LocationType constants for different location types
const (
	LocationTypeStore             = "store"
	LocationTypeWarehouse         = "warehouse"
	LocationTypeFulfillmentCenter = "fulfillment_center"
	LocationTypeOnline            = "online"
)

// NewStoreLocation creates a new store location
func NewStoreLocation(name, locationType, address, city, state, postalCode, country string) *StoreLocation {
	now := time.Now()
	return &StoreLocation{
		ID:         uuid.New().String(),
		Name:       name,
		Type:       locationType,
		Address:    address,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// InventoryTransfer represents a transfer of inventory between locations
type InventoryTransfer struct {
	ID                 string    `bson:"_id,omitempty"`
	ProductID          string    `bson:"product_id"`
	SKU                string    `bson:"sku"`
	SourceLocationID   string    `bson:"source_location_id"`
	DestLocationID     string    `bson:"dest_location_id"`
	Quantity           int32     `bson:"quantity"`
	Status             string    `bson:"status"` // "pending", "in_transit", "completed", "cancelled"
	RequestedBy        string    `bson:"requested_by"`
	ApprovedBy         string    `bson:"approved_by,omitempty"`
	Notes              string    `bson:"notes,omitempty"`
	EstimatedArrival   time.Time `bson:"estimated_arrival,omitempty"`
	CompletedAt        time.Time `bson:"completed_at,omitempty"`
	CreatedAt          time.Time `bson:"created_at"`
	UpdatedAt          time.Time `bson:"updated_at"`
}

// Transfer status constants
const (
	TransferStatusPending    = "pending"
	TransferStatusInTransit  = "in_transit"
	TransferStatusCompleted  = "completed"
	TransferStatusCancelled  = "cancelled"
)

// NewInventoryTransfer creates a new inventory transfer request
func NewInventoryTransfer(productID, sku, sourceLocationID, destLocationID string, quantity int32, requestedBy string) *InventoryTransfer {
	now := time.Now()
	return &InventoryTransfer{
		ID:               uuid.New().String(),
		ProductID:        productID,
		SKU:              sku,
		SourceLocationID: sourceLocationID,
		DestLocationID:   destLocationID,
		Quantity:         quantity,
		Status:           TransferStatusPending,
		RequestedBy:      requestedBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// UpdateTransferStatus updates the status of a transfer
func (t *InventoryTransfer) UpdateTransferStatus(status string) {
	t.Status = status
	t.UpdatedAt = time.Now()
	
	if status == TransferStatusCompleted {
		t.CompletedAt = time.Now()
	}
}

// InventoryItem represents a product's inventory information
type InventoryItem struct {
	ID                string    `bson:"_id,omitempty"`
	ProductID         string    `bson:"product_id"`
	Quantity          int32     `bson:"quantity"`
	Reserved          int32     `bson:"reserved"`
	SKU               string    `bson:"sku"`
	LocationID        string    `bson:"location_id"`
	ShelfLocation     string    `bson:"shelf_location,omitempty"` // For precise in-store location (aisle/shelf/bin)
	MinimumStock      int32     `bson:"minimum_stock,omitempty"`
	MaximumStock      int32     `bson:"maximum_stock,omitempty"`
	ReorderPoint      int32     `bson:"reorder_point,omitempty"`
	ReorderQuantity   int32     `bson:"reorder_quantity,omitempty"`
	LastCountDate     time.Time `bson:"last_count_date,omitempty"`
	NextCountDate     time.Time `bson:"next_count_date,omitempty"`
	LastUpdated       time.Time `bson:"last_updated"`
	CreatedAt         time.Time `bson:"created_at"`
}

// NewInventoryItem creates a new inventory item
func NewInventoryItem(productID string, quantity int32, sku string, locationID string) *InventoryItem {
	now := time.Now()
	return &InventoryItem{
		ID:                uuid.New().String(),
		ProductID:         productID,
		Quantity:          quantity,
		Reserved:          0,
		SKU:               sku,
		LocationID:        locationID,
		MinimumStock:      0,
		MaximumStock:      0,
		ReorderPoint:      0,
		ReorderQuantity:   0,
		LastUpdated:       now,
		CreatedAt:         now,
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

// SetShelfLocation sets the precise shelf location within a store
func (i *InventoryItem) SetShelfLocation(shelfLocation string) {
	i.ShelfLocation = shelfLocation
	i.LastUpdated = time.Now()
}

// SetReorderParameters sets inventory reordering parameters
func (i *InventoryItem) SetReorderParameters(minimumStock, maximumStock, reorderPoint, reorderQuantity int32) {
	i.MinimumStock = minimumStock
	i.MaximumStock = maximumStock
	i.ReorderPoint = reorderPoint
	i.ReorderQuantity = reorderQuantity
	i.LastUpdated = time.Now()
}

// NeedsReorder checks if the item needs reordering based on defined parameters
func (i *InventoryItem) NeedsReorder() bool {
	return i.ReorderPoint > 0 && i.Quantity <= i.ReorderPoint
}

// ScheduleInventoryCount schedules the next inventory count date
func (i *InventoryItem) ScheduleInventoryCount(nextCountDate time.Time) {
	i.LastCountDate = time.Now()
	i.NextCountDate = nextCountDate
	i.LastUpdated = time.Now()
}

// TransferStock transfers quantity to another inventory item
// Returns error if not enough available stock
func (i *InventoryItem) TransferStock(quantity int32, destination *InventoryItem) error {
	if quantity <= 0 {
		return errors.New("transfer quantity must be positive")
	}
	
	if !i.IsAvailable(quantity) {
		return errors.New("not enough available stock for transfer")
	}
	
	// Remove from source
	i.Quantity -= quantity
	i.LastUpdated = time.Now()
	
	// Add to destination
	destination.Quantity += quantity
	destination.LastUpdated = time.Now()
	
	return nil
}
