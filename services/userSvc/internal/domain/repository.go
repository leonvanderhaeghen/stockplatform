package domain

import "context"

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// Create adds a new user
	Create(ctx context.Context, user *User) error
	
	// GetByID finds a user by ID
	GetByID(ctx context.Context, id string) (*User, error)
	
	// GetByEmail finds a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *User) error
	
	// Delete removes a user
	Delete(ctx context.Context, id string) error
	
	// List returns all users with optional filtering and pagination
	List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*User, error)
}

// AddressRepository defines the interface for address persistence
type AddressRepository interface {
	// Create adds a new address
	Create(ctx context.Context, address *Address) error
	
	// GetByID finds an address by ID
	GetByID(ctx context.Context, id string) (*Address, error)
	
	// GetByUserID finds addresses for a user
	GetByUserID(ctx context.Context, userID string) ([]*Address, error)
	
	// GetDefaultByUserID finds the default address for a user
	GetDefaultByUserID(ctx context.Context, userID string) (*Address, error)
	
	// Update updates an existing address
	Update(ctx context.Context, address *Address) error
	
	// Delete removes an address
	Delete(ctx context.Context, id string) error
	
	// SetDefaultAddress sets an address as the default for a user
	SetDefaultAddress(ctx context.Context, userID string, addressID string) error
}
