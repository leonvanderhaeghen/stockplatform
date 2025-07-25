package domain

import (
	"context"
	"time"
)

// InventoryHistory represents a historical record of changes to an inventory item
type InventoryHistory struct {
	ID           string    `bson:"_id,omitempty"`
	InventoryID  string    `bson:"inventory_id"`
	ChangeType   string    `bson:"change_type"`   // e.g., QUANTITY_CHANGE, STATUS_UPDATE, etc.
	Description  string    `bson:"description"`   // Human-readable description of the change
	QuantityBefore int32   `bson:"quantity_before"`
	QuantityAfter  int32   `bson:"quantity_after"`
	ReferenceID   string    `bson:"reference_id,omitempty"`   // e.g., order ID, transfer ID, etc.
	ReferenceType string    `bson:"reference_type,omitempty"` // e.g., ORDER, TRANSFER, ADJUSTMENT, etc.
	PerformedBy  string    `bson:"performed_by"`  // User ID who performed the change
	CreatedAt    time.Time `bson:"created_at"`    // Timestamp of the change
}

// TransferRepository defines the interface for inventory transfer persistence
type TransferRepository interface {
	// Create adds a new inventory transfer
	Create(ctx context.Context, transfer *Transfer) error
	
	// GetByID finds an inventory transfer by its ID
	GetByID(ctx context.Context, id string) (*Transfer, error)
	
	// Update updates an existing inventory transfer
	Update(ctx context.Context, transfer *Transfer) error
	
	// ListBySourceLocation lists transfers from a specific source location
	ListBySourceLocation(ctx context.Context, sourceLocationID string, limit, offset int) ([]*Transfer, error)
	
	// ListByDestLocation lists transfers to a specific destination location
	ListByDestLocation(ctx context.Context, destLocationID string, limit, offset int) ([]*Transfer, error)
	
	// ListByProduct lists transfers for a specific product
	ListByProduct(ctx context.Context, productID string, limit, offset int) ([]*Transfer, error)
	
	// ListByStatus lists transfers with a specific status
	ListByStatus(ctx context.Context, status TransferStatus, limit, offset int) ([]*Transfer, error)
	
	// ListPendingTransfers lists all pending transfers
	ListPendingTransfers(ctx context.Context, limit, offset int) ([]*Transfer, error)
}

// InventoryRepository defines the interface for inventory persistence
type InventoryRepository interface {
	// Create adds a new inventory item
	Create(ctx context.Context, item *InventoryItem) error
	
	// GetByID finds an inventory item by its ID
	GetByID(ctx context.Context, id string) (*InventoryItem, error)
	
	// GetByProductID finds inventory items by product ID
	GetByProductID(ctx context.Context, productID string) ([]*InventoryItem, error)
	
	// GetBySKU finds an inventory item by SKU (across all locations)
	GetBySKU(ctx context.Context, sku string) ([]*InventoryItem, error)
	
	// GetByProductAndLocation finds inventory items by product ID and location
	GetByProductAndLocation(ctx context.Context, productID, locationID string) (*InventoryItem, error)
	
	// GetBySKUAndLocation finds an inventory item by SKU and location
	GetBySKUAndLocation(ctx context.Context, sku, locationID string) (*InventoryItem, error)
	
	// Update updates an existing inventory item
	Update(ctx context.Context, item *InventoryItem) error
	
	// Delete removes an inventory item
	Delete(ctx context.Context, id string) error
	
	// List returns all inventory items with optional pagination
	List(ctx context.Context, limit, offset int) ([]*InventoryItem, error)
	
	// ListByLocation returns all inventory items for a specific location
	ListByLocation(ctx context.Context, locationID string, limit, offset int) ([]*InventoryItem, error)
	
	// ListLowStock returns inventory items that are below their reorder point
	ListLowStock(ctx context.Context, limit, offset int) ([]*InventoryItem, error)
	
	// ListByStockStatus returns inventory items based on stock status (in stock, low stock, out of stock)
	ListByStockStatus(ctx context.Context, status string, limit, offset int) ([]*InventoryItem, error)
	
	// GetByOrderAndLocation finds inventory items reserved for a specific order at a specific location
	GetByOrderAndLocation(ctx context.Context, orderID, locationID string) ([]*InventoryItem, error)
	
	// AdjustStock adjusts inventory quantity and records reason
	AdjustStock(ctx context.Context, itemID string, quantity int32, reason string, performedBy string) error
	
	// GetHistory retrieves the history of changes for a specific inventory item
	GetHistory(ctx context.Context, inventoryID string, limit, offset int32) ([]*InventoryHistory, int32, error)
	
	// RecordHistory adds a new history entry for an inventory item
	RecordHistory(ctx context.Context, history *InventoryHistory) error
}
