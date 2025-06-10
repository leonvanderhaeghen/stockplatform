package domain

import "context"

// SupplierUseCase defines the interface for supplier business logic
type SupplierUseCase interface {
	// Supplier CRUD operations
	CreateSupplier(ctx context.Context, supplier *Supplier) (*Supplier, error)
	GetSupplier(ctx context.Context, id string) (*Supplier, error)
	UpdateSupplier(ctx context.Context, supplier *Supplier) (*Supplier, error)
	DeleteSupplier(ctx context.Context, id string) error
	ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*Supplier, int32, error)

	// Adapter-related operations
	RegisterAdapter(ctx context.Context, adapter SupplierAdapter) error
	ListAdapters(ctx context.Context) ([]string, error)
	GetAdapterCapabilities(ctx context.Context, adapterName string) (map[string]bool, error)
	TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error
	SyncAdapterProducts(ctx context.Context, adapterName string, options SupplierSyncOptions) (*SupplierSyncStats, error)
	SyncAdapterInventory(ctx context.Context, adapterName string, options SupplierSyncOptions) (*SupplierSyncStats, error)
}
