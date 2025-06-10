package domain

import (
	"context"
	"time"
)

// SupplierProductData represents a product from an external supplier
type SupplierProductData struct {
	ExternalID        string             // ID from supplier's system
	Name              string             // Product name
	SKU               string             // Stock Keeping Unit
	Description       string             // Product description
	Price             float64            // Price (in supplier's currency)
	Currency          string             // Currency code (USD, EUR, etc.)
	StockQuantity     int32              // Current stock quantity
	MinimumOrderQty   int32              // Minimum order quantity
	LeadTimeDays      int32              // Lead time for delivery in days
	Categories        []string           // Product categories
	Barcode           string             // Barcode (EAN, UPC, etc.)
	Weight            float64            // Weight (in kg)
	Dimensions        *ProductDimensions // Product dimensions
	Images            []string           // URLs of product images
	Active            bool               // Whether the product is active
	LastUpdated       time.Time          // Last updated timestamp
	AdditionalDetails map[string]string  // Any additional product details
}

// ProductDimensions represents the physical dimensions of a product
type ProductDimensions struct {
	Length float64
	Width  float64
	Height float64
	Unit   string // cm, inches, etc.
}

// SupplierInventoryData represents inventory data from an external supplier
type SupplierInventoryData struct {
	ProductID      string    // Internal product ID
	ExternalID     string    // ID from supplier's system
	StockQuantity  int32     // Current stock quantity
	ReservedQty    int32     // Quantity reserved for orders
	AvailableQty   int32     // Available quantity (StockQty - ReservedQty)
	IncomingQty    int32     // Quantity expected to arrive
	ExpectedDate   time.Time // Expected date for incoming stock
	LocationCode   string    // Location/warehouse code
	LastUpdateTime time.Time // Last inventory update timestamp
}

// SupplierSyncStats provides statistics about a sync operation
type SupplierSyncStats struct {
	StartTime          time.Time
	EndTime            time.Time
	ProductsProcessed  int
	ProductsCreated    int
	ProductsUpdated    int
	ProductsErrored    int
	InventoryProcessed int
	InventoryUpdated   int
	InventoryErrored   int
	Errors             []string
}

// SupplierSyncOptions provides configuration options for a sync operation
type SupplierSyncOptions struct {
	SyncProducts  bool // Whether to sync product data
	SyncInventory bool // Whether to sync inventory data
	FullSync      bool // Whether to do a full sync vs incremental
	SyncImages    bool // Whether to sync product images
	BatchSize     int  // Batch size for processing
	FromDate      time.Time
	ToDate        time.Time
}

// SupplierAdapter defines the interface for supplier data integrations
type SupplierAdapter interface {
	// Name returns the name of this adapter
	Name() string

	// Initialize initializes the adapter with supplier-specific configuration
	Initialize(ctx context.Context, config map[string]string) error

	// TestConnection tests the connection to the supplier's system
	TestConnection(ctx context.Context) error

	// GetProducts fetches products from the supplier
	GetProducts(ctx context.Context, options SupplierSyncOptions) ([]SupplierProductData, error)

	// GetInventory fetches inventory data from the supplier
	GetInventory(ctx context.Context, externalIDs []string, options SupplierSyncOptions) ([]SupplierInventoryData, error)

	// SyncProducts syncs product data from the supplier
	SyncProducts(ctx context.Context, options SupplierSyncOptions) (*SupplierSyncStats, error)

	// SyncInventory syncs inventory data from the supplier
	SyncInventory(ctx context.Context, options SupplierSyncOptions) (*SupplierSyncStats, error)

	// GetCapabilities returns the capabilities of this adapter
	GetCapabilities(ctx context.Context) map[string]bool
}

// AdapterRegistry manages supplier adapters
type AdapterRegistry interface {
	// Register registers a new supplier adapter
	Register(adapter SupplierAdapter) error

	// Get returns an adapter by name
	Get(name string) (SupplierAdapter, error)

	// List returns all registered adapters
	List() []SupplierAdapter
}
