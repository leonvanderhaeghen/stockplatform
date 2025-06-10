package application

import (
	"context"
	"fmt"
	"time"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

// BaseAdapter provides common functionality for supplier adapters
type BaseAdapter struct {
	name         string
	config       map[string]string
	capabilities map[string]bool
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(name string) *BaseAdapter {
	return &BaseAdapter{
		name:         name,
		config:       make(map[string]string),
		capabilities: make(map[string]bool),
	}
}

// Name returns the name of this adapter
func (a *BaseAdapter) Name() string {
	return a.name
}

// Initialize initializes the adapter with supplier-specific configuration
func (a *BaseAdapter) Initialize(ctx context.Context, config map[string]string) error {
	a.config = config
	return nil
}

// TestConnection provides a default implementation for testing connections
func (a *BaseAdapter) TestConnection(ctx context.Context) error {
	return fmt.Errorf("TestConnection not implemented for adapter %s", a.name)
}

// GetProducts provides a default implementation for fetching products
func (a *BaseAdapter) GetProducts(ctx context.Context, options domain.SupplierSyncOptions) ([]domain.SupplierProductData, error) {
	return nil, fmt.Errorf("GetProducts not implemented for adapter %s", a.name)
}

// GetInventory provides a default implementation for fetching inventory
func (a *BaseAdapter) GetInventory(ctx context.Context, externalIDs []string, options domain.SupplierSyncOptions) ([]domain.SupplierInventoryData, error) {
	return nil, fmt.Errorf("GetInventory not implemented for adapter %s", a.name)
}

// SyncProducts provides a base implementation for synchronizing products
func (a *BaseAdapter) SyncProducts(ctx context.Context, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error) {
	stats := &domain.SupplierSyncStats{
		StartTime: time.Now(),
	}

	// Default implementation could be extended in concrete adapters
	products, err := a.GetProducts(ctx, options)
	if err != nil {
		stats.EndTime = time.Now()
		stats.Errors = append(stats.Errors, fmt.Sprintf("Error fetching products: %v", err))
		return stats, err
	}

	stats.ProductsProcessed = len(products)
	stats.EndTime = time.Now()
	return stats, nil
}

// SyncInventory provides a base implementation for synchronizing inventory
func (a *BaseAdapter) SyncInventory(ctx context.Context, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error) {
	stats := &domain.SupplierSyncStats{
		StartTime: time.Now(),
	}

	// Default implementation could be extended in concrete adapters
	// This would typically involve getting product IDs and then fetching inventory for them
	stats.EndTime = time.Now()
	return stats, nil
}

// GetCapabilities returns the capabilities of this adapter
func (a *BaseAdapter) GetCapabilities(ctx context.Context) map[string]bool {
	return a.capabilities
}

// SetCapability sets a specific capability for this adapter
func (a *BaseAdapter) SetCapability(capability string, supported bool) {
	a.capabilities[capability] = supported
}

// GetConfig returns the configuration value for a key
func (a *BaseAdapter) GetConfig(key string) (string, bool) {
	value, ok := a.config[key]
	return value, ok
}

// ValidateRequiredConfig checks if all required config keys are present
func (a *BaseAdapter) ValidateRequiredConfig(requiredKeys []string) error {
	missingKeys := []string{}

	for _, key := range requiredKeys {
		if _, ok := a.config[key]; !ok {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required configuration keys: %v", missingKeys)
	}

	return nil
}
