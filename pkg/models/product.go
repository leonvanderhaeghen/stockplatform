package models

import "time"

// Product represents a product in the domain
type Product struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SKU         string     `json:"sku"`
	Price       float64    `json:"price"`
	Cost        float64    `json:"cost"`
	Category    string     `json:"category"`
	Brand       string     `json:"brand"`
	Weight      float64    `json:"weight"`
	Dimensions  *Dimensions `json:"dimensions,omitempty"`
	IsActive    bool       `json:"is_active"`
	SupplierID  string     `json:"supplier_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Dimensions represents product dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"` // e.g., "cm", "in"
}

// CreateProductResponse represents the response from creating a product
type CreateProductResponse struct {
	Product *Product `json:"product"`
	Message string   `json:"message"`
}

// ListProductsResponse represents the response from listing products
type ListProductsResponse struct {
	Products   []*Product `json:"products"`
	TotalCount int32      `json:"total_count"`
}

// UpdateProductResponse represents the response from updating a product
type UpdateProductResponse struct {
	Product *Product `json:"product"`
	Message string   `json:"message"`
}

// ProductSearchResult represents a product search result
type ProductSearchResult struct {
	Products   []*Product `json:"products"`
	TotalCount int32      `json:"total_count"`
	Query      string     `json:"query"`
}
