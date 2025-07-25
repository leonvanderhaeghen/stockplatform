package models

import "time"

// Supplier represents a supplier in the domain
type Supplier struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     *Address  `json:"address,omitempty"`
	ContactName string    `json:"contact_name"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSupplierResponse represents the response from creating a supplier
type CreateSupplierResponse struct {
	Supplier *Supplier `json:"supplier"`
	Message  string    `json:"message"`
}

// ListSuppliersResponse represents the response from listing suppliers
type ListSuppliersResponse struct {
	Suppliers  []*Supplier `json:"suppliers"`
	TotalCount int32       `json:"total_count"`
}

// UpdateSupplierResponse represents the response from updating a supplier
type UpdateSupplierResponse struct {
	Supplier *Supplier `json:"supplier"`
	Message  string    `json:"message"`
}

// SupplierSearchResult represents a supplier search result
type SupplierSearchResult struct {
	Suppliers  []*Supplier `json:"suppliers"`
	TotalCount int32       `json:"total_count"`
	Query      string      `json:"query"`
}
