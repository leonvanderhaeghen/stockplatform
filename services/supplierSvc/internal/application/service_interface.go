package application

import (
	"context"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

// SupplierService is the interface that defines the application layer for supplier operations
// This follows clean architecture principles by providing an application layer interface
// that can be used by the interfaces layer (e.g., gRPC, HTTP) without direct domain dependency
type SupplierService interface {
	// Basic CRUD operations
	CreateSupplier(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error)
	GetSupplier(ctx context.Context, id string) (*domain.Supplier, error)
	UpdateSupplier(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error)
	DeleteSupplier(ctx context.Context, id string) error
	ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*domain.Supplier, int32, error)
	
	// Adapter-related operations
	RegisterAdapter(ctx context.Context, adapter domain.SupplierAdapter) error
	ListAdapters(ctx context.Context) ([]string, error)
	GetAdapterCapabilities(ctx context.Context, adapterName string) (map[string]bool, error)
	TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error
	SyncAdapterProducts(ctx context.Context, adapterName string, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error)
	SyncAdapterInventory(ctx context.Context, adapterName string, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error)
}
