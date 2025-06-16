package domain

import (
	"context"
	"errors"
)

// ErrLocationNotFound is returned when a location is not found
var ErrLocationNotFound = errors.New("location not found")

// LocationRepository defines operations for store location persistence
type LocationRepository interface {
	// Create creates a new store location
	Create(ctx context.Context, location *StoreLocation) error
	
	// GetByID gets a store location by ID
	GetByID(ctx context.Context, id string) (*StoreLocation, error)
	
	// GetByName gets a store location by name
	GetByName(ctx context.Context, name string) (*StoreLocation, error)
	
	// Update updates an existing store location
	Update(ctx context.Context, location *StoreLocation) error
	
	// Delete deletes a store location
	Delete(ctx context.Context, id string) error
	
	// List lists store locations with pagination
	// If includeInactive is false, returns only active locations
	List(ctx context.Context, limit int, offset int, includeInactive bool) ([]*StoreLocation, error)

	// ListByType lists store locations of a specific type
	ListByType(ctx context.Context, locationType string, limit int, offset int) ([]*StoreLocation, error)
}
