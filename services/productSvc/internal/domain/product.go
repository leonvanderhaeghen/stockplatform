package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product in the system
type Product struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string                 `bson:"name" json:"name" validate:"required,min=3,max=100"`
	Description   string                 `bson:"description,omitempty" json:"description,omitempty"`
	CostPrice     string                 `bson:"cost_price" json:"cost_price" validate:"required,decimal"`
	SellingPrice  string                 `bson:"selling_price" json:"selling_price" validate:"required,decimal"`
	Currency      string                 `bson:"currency" json:"currency" validate:"required,iso4217"`
	SKU           string                 `bson:"sku" json:"sku" validate:"required,alphanum"`
	Barcode       string                 `bson:"barcode,omitempty" json:"barcode,omitempty"`
	CategoryIDs   []string               `bson:"category_ids" json:"category_ids" validate:"required,min=1,dive,required"`
	SupplierID    string                 `bson:"supplier_id" json:"supplier_id" validate:"required"`
	IsActive      bool                   `bson:"is_active" json:"is_active"`
	IsVisible     map[string]bool        `bson:"is_visible,omitempty" json:"is_visible,omitempty"`
	Variants      []Variant              `bson:"variants,omitempty" json:"variants,omitempty"`

	ImageURLs     []string               `bson:"image_urls,omitempty" json:"image_urls,omitempty"`
	VideoURLs     []string               `bson:"video_urls,omitempty" json:"video_urls,omitempty"`
	Metadata      map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt     time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time              `bson:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time             `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}



// Variant represents product variations (size, color, etc.)
type Variant struct {
	ID           string                 `bson:"id" json:"id"`
	Name         string                 `bson:"name" json:"name" validate:"required"`
	Options      []VariantOption        `bson:"options" json:"options" validate:"required,min=1"`
	CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time              `bson:"updated_at" json:"updated_at"`
}

// VariantOption represents a specific option for a variant
type VariantOption struct {
	ID             string    `bson:"id" json:"id"`
	Name           string    `bson:"name" json:"name" validate:"required"`
	Value          string    `bson:"value" json:"value" validate:"required"`
	SKU            string    `bson:"sku,omitempty" json:"sku,omitempty"`
	Barcode        string    `bson:"barcode,omitempty" json:"barcode,omitempty"`
	PriceAdjustment string    `bson:"price_adjustment,omitempty" json:"price_adjustment,omitempty"`
	IsDefault      bool      `bson:"is_default,omitempty" json:"is_default,omitempty"`
	ImageURL       string    `bson:"image_url,omitempty" json:"image_url,omitempty"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at" json:"updated_at"`
}

// ParsePrice parses a price string into a decimal value
func ParsePrice(priceStr string) (decimal.Decimal, error) {
	if priceStr == "" {
		return decimal.Zero, errors.New("price cannot be empty")
	}

	// Remove any currency symbols and thousands separators
	priceStr = strings.ReplaceAll(priceStr, "$", "")
	priceStr = strings.ReplaceAll(priceStr, ",", "")

	// Parse the decimal value
	price, err := decimal.NewFromString(priceStr)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid price format: %v", err)
	}

	// Ensure the price is non-negative
	if price.IsNegative() {
		return decimal.Zero, errors.New("price cannot be negative")
	}

	// Round to 2 decimal places (cents)
	return price.Round(2), nil
}

// FormatPrice formats a decimal price as a string with 2 decimal places
func FormatPrice(price decimal.Decimal) string {
	return price.StringFixed(2)
}

// CalculateProfitMargin calculates the profit margin percentage
// Returns (profitMargin, profitAmount, error)
func CalculateProfitMargin(costPrice, sellingPrice decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	if costPrice.IsNegative() || sellingPrice.IsNegative() {
		return decimal.Zero, decimal.Zero, errors.New("prices cannot be negative")
	}
	
	if costPrice.IsZero() {
		return decimal.Zero, decimal.Zero, errors.New("cost price cannot be zero")
	}

	profit := sellingPrice.Sub(costPrice)
	margin := profit.Div(costPrice).Mul(decimal.NewFromInt(100))

	return margin.Round(2), profit, nil
}

// CalculateMarkup calculates the markup percentage
// Returns (markup, markupAmount, error)
func CalculateMarkup(costPrice, sellingPrice decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	if costPrice.IsNegative() || sellingPrice.IsNegative() {
		return decimal.Zero, decimal.Zero, errors.New("prices cannot be negative")
	}
	
	if costPrice.IsZero() {
		return decimal.Zero, decimal.Zero, errors.New("cost price cannot be zero")
	}

	markupAmount := sellingPrice.Sub(costPrice)
	markup := markupAmount.Div(costPrice).Mul(decimal.NewFromInt(100))

	return markup.Round(2), markupAmount, nil
}

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error

	// Listing and searching
	List(ctx context.Context, opts *ListOptions) ([]*Product, int64, error)
	Search(ctx context.Context, query string, opts *ListOptions) ([]*Product, int64, error)
	GetBySupplier(ctx context.Context, supplierID string, opts *ListOptions) ([]*Product, int64, error)
	GetByCategory(ctx context.Context, categoryID string, opts *ListOptions) ([]*Product, int64, error)

	// Inventory operations
	UpdateStock(ctx context.Context, id string, quantity int32) error
	BulkUpdateStock(ctx context.Context, updates map[string]int32) error
	GetLowStockProducts(ctx context.Context, threshold int32, opts *ListOptions) ([]*Product, int64, error)

	// Visibility and publishing
	BulkUpdateVisibility(ctx context.Context, supplierID string, productIDs []string, isVisible bool) error
	PublishProducts(ctx context.Context, productIDs []string, publish bool) error

	// Variant operations
	UpdateVariantStock(ctx context.Context, productID, variantID string, quantity int32) error
}

// ProductUseCase defines the business logic for product operations
type ProductUseCase interface {
	// Basic CRUD operations
	CreateProduct(ctx context.Context, product *Product) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*Product, error)
	GetProductByBarcode(ctx context.Context, barcode string) (*Product, error)
	UpdateProduct(ctx context.Context, product *Product) error
	DeleteProduct(ctx context.Context, id string) error
	SoftDeleteProduct(ctx context.Context, id string) error

	// Listing and searching
	ListProducts(ctx context.Context, opts *ListOptions) ([]*Product, int64, error)
	SearchProducts(ctx context.Context, query string, opts *ListOptions) ([]*Product, int64, error)
	GetProductsBySupplier(ctx context.Context, supplierID string, opts *ListOptions) ([]*Product, int64, error)
	GetProductsByCategory(ctx context.Context, categoryID string, opts *ListOptions) ([]*Product, int64, error)

	// Inventory management
	UpdateProductStock(ctx context.Context, id string, quantity int32) error
	BulkUpdateStock(ctx context.Context, updates map[string]int32) error
	AdjustStock(ctx context.Context, id string, adjustment int32, note string) error
	GetLowStockProducts(ctx context.Context, threshold int32, opts *ListOptions) ([]*Product, int64, error)

	// Pricing operations
	UpdateProductPricing(ctx context.Context, id string, costPrice, sellingPrice string) error
	BulkUpdatePricing(ctx context.Context, updates map[string]struct{ CostPrice, SellingPrice string }) error
	CalculateProfitMargin(productID string) (decimal.Decimal, decimal.Decimal, error)
	CalculateMarkup(productID string) (decimal.Decimal, decimal.Decimal, error)

	// Variant management
	AddVariant(ctx context.Context, productID string, variant *Variant) error
	UpdateVariant(ctx context.Context, productID string, variant *Variant) error
	RemoveVariant(ctx context.Context, productID, variantID string) error
	UpdateVariantStock(ctx context.Context, productID, variantID string, quantity int32) error

	// Visibility and publishing
	BulkUpdateProductVisibility(ctx context.Context, supplierID string, productIDs []string, isVisible bool) error
	PublishProducts(ctx context.Context, productIDs []string, publish bool) error

	// Validation and utilities
	ValidateProduct(product *Product) error
	GenerateProductReport(ctx context.Context, format string) ([]byte, error)
}
