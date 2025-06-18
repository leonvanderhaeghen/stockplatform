package domain

import (
	"context"
	"errors"
)

// ErrOptimisticLockFailed is returned when an optimistic lock update fails due to version mismatch
var ErrOptimisticLockFailed = errors.New("optimistic lock failed: order was modified by another process")

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	// Create adds a new order
	Create(ctx context.Context, order *Order) error
	
	// GetByID finds an order by its ID
	GetByID(ctx context.Context, id string) (*Order, error)
	
	// GetByUserID finds orders for a specific user
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Order, error)
	
	// Update updates an existing order
	Update(ctx context.Context, order *Order) error
	
	// UpdateWithOptimisticLock updates an order with version checking to prevent concurrent modifications
	UpdateWithOptimisticLock(ctx context.Context, order *Order, expectedVersion int32) error
	
	// Delete removes an order
	Delete(ctx context.Context, id string) error
	
	// List returns all orders with optional filtering and pagination
	List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*Order, error)
	
	// Count returns the number of orders matching a filter
	Count(ctx context.Context, filter map[string]interface{}) (int64, error)
}
