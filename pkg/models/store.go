package models

import "time"

// Store represents a physical store location
type Store struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Address     *Address  `json:"address" bson:"address"`
	Phone       string    `json:"phone" bson:"phone"`
	Email       string    `json:"email" bson:"email"`
	IsActive    bool      `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// CreateStoreResponse represents the response after creating a store
type CreateStoreResponse struct {
	Store   *Store `json:"store"`
	Message string `json:"message"`
}

// ListStoresResponse represents the response from listing stores
type ListStoresResponse struct {
	Stores      []*Store `json:"stores"`
	TotalCount  int32    `json:"total_count"`
	HasNextPage bool     `json:"has_next_page"`
}
