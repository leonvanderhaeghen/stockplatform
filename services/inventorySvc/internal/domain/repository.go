package domain

import "context"

// InventoryRepository defines the interface for inventory persistence
type InventoryRepository interface {
	// Create adds a new inventory item
	Create(ctx context.Context, item *InventoryItem) error
	
	// GetByID finds an inventory item by its ID
	GetByID(ctx context.Context, id string) (*InventoryItem, error)
	
	// GetByProductID finds inventory items by product ID
	GetByProductID(ctx context.Context, productID string) (*InventoryItem, error)
	
	// GetBySKU finds an inventory item by SKU
	GetBySKU(ctx context.Context, sku string) (*InventoryItem, error)
	
	// Update updates an existing inventory item
	Update(ctx context.Context, item *InventoryItem) error
	
	// Delete removes an inventory item
	Delete(ctx context.Context, id string) error
	
	// List returns all inventory items with optional pagination
	List(ctx context.Context, limit, offset int) ([]*InventoryItem, error)
}
