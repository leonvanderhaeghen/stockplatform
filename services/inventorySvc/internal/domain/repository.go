package domain

import "context"

// LocationRepository defines the interface for store location persistence
type LocationRepository interface {
	// Create adds a new store location
	Create(ctx context.Context, location *StoreLocation) error
	
	// GetByID finds a store location by its ID
	GetByID(ctx context.Context, id string) (*StoreLocation, error)
	
	// GetByName finds a store location by name
	GetByName(ctx context.Context, name string) (*StoreLocation, error)
	
	// Update updates an existing store location
	Update(ctx context.Context, location *StoreLocation) error
	
	// Delete marks a store location as inactive
	Delete(ctx context.Context, id string) error
	
	// List returns all store locations with optional pagination and filters
	List(ctx context.Context, limit, offset int, includeInactive bool) ([]*StoreLocation, error)
	
	// ListByType returns all store locations of a specific type
	ListByType(ctx context.Context, locationType string, limit, offset int) ([]*StoreLocation, error)
}

// TransferRepository defines the interface for inventory transfer persistence
type TransferRepository interface {
	// Create adds a new inventory transfer
	Create(ctx context.Context, transfer *InventoryTransfer) error
	
	// GetByID finds an inventory transfer by its ID
	GetByID(ctx context.Context, id string) (*InventoryTransfer, error)
	
	// Update updates an existing inventory transfer
	Update(ctx context.Context, transfer *InventoryTransfer) error
	
	// ListBySourceLocation lists transfers from a specific source location
	ListBySourceLocation(ctx context.Context, sourceLocationID string, limit, offset int) ([]*InventoryTransfer, error)
	
	// ListByDestLocation lists transfers to a specific destination location
	ListByDestLocation(ctx context.Context, destLocationID string, limit, offset int) ([]*InventoryTransfer, error)
	
	// ListByProduct lists transfers for a specific product
	ListByProduct(ctx context.Context, productID string, limit, offset int) ([]*InventoryTransfer, error)
	
	// ListByStatus lists transfers with a specific status
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*InventoryTransfer, error)
	
	// ListPendingTransfers lists all pending transfers
	ListPendingTransfers(ctx context.Context, limit, offset int) ([]*InventoryTransfer, error)
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
	
	// AdjustStock adjusts inventory quantity and records reason
	AdjustStock(ctx context.Context, itemID string, quantity int32, reason string, performedBy string) error
}
