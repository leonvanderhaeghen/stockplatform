package application

import (
	"context"
	"fmt"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

type supplierService struct {
	repo domain.SupplierRepository
	adapterRegistry domain.AdapterRegistry
}

// NewSupplierService creates a new supplier service
func NewSupplierService(repo domain.SupplierRepository) domain.SupplierUseCase {
	return &supplierService{
		repo: repo,
		adapterRegistry: NewAdapterRegistry(),
	}
}

func (s *supplierService) CreateSupplier(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error) {
	// Validate supplier data
	if supplier.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Create the supplier
	return s.repo.Create(ctx, supplier)
}

func (s *supplierService) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *supplierService) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error) {
	// Check if supplier exists
	existing, err := s.repo.GetByID(ctx, supplier.ID.Hex())
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.Name = supplier.Name
	existing.ContactPerson = supplier.ContactPerson
	existing.Email = supplier.Email
	existing.Phone = supplier.Phone
	existing.Address = supplier.Address
	existing.City = supplier.City
	existing.State = supplier.State
	existing.Country = supplier.Country
	existing.PostalCode = supplier.PostalCode
	existing.TaxID = supplier.TaxID
	existing.Website = supplier.Website
	existing.Currency = supplier.Currency
	existing.LeadTimeDays = supplier.LeadTimeDays
	existing.PaymentTerms = supplier.PaymentTerms
	existing.Metadata = supplier.Metadata

	// Update in repository
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *supplierService) DeleteSupplier(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *supplierService) ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*domain.Supplier, int32, error) {
	// Ensure page and pageSize are within reasonable bounds
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.List(ctx, page, pageSize, search)
}

// RegisterAdapter registers a new supplier adapter
func (s *supplierService) RegisterAdapter(ctx context.Context, adapter domain.SupplierAdapter) error {
	return s.adapterRegistry.Register(adapter)
}

// ListAdapters lists all registered adapters
func (s *supplierService) ListAdapters(ctx context.Context) ([]string, error) {
	adapters := s.adapterRegistry.List()
	adapterNames := make([]string, len(adapters))
	for i, adapter := range adapters {
		adapterNames[i] = adapter.Name()
	}
	return adapterNames, nil
}

// GetAdapterCapabilities returns the capabilities of a specific adapter
func (s *supplierService) GetAdapterCapabilities(ctx context.Context, adapterName string) (map[string]bool, error) {
	adapter, err := s.adapterRegistry.Get(adapterName)
	if err != nil {
		return nil, err
	}
	return adapter.GetCapabilities(ctx), nil
}

// TestAdapterConnection tests the connection to a supplier's system
func (s *supplierService) TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error {
	adapter, err := s.adapterRegistry.Get(adapterName)
	if err != nil {
		return err
	}

	// Initialize the adapter with the provided configuration
	if err := adapter.Initialize(ctx, config); err != nil {
		return fmt.Errorf("failed to initialize adapter: %w", err)
	}

	// Test the connection
	return adapter.TestConnection(ctx)
}

// SyncAdapterProducts syncs product data from the supplier
func (s *supplierService) SyncAdapterProducts(ctx context.Context, adapterName string, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error) {
	adapter, err := s.adapterRegistry.Get(adapterName)
	if err != nil {
		return nil, err
	}

	// Sync products
	return adapter.SyncProducts(ctx, options)
}

// SyncAdapterInventory syncs inventory data from the supplier
func (s *supplierService) SyncAdapterInventory(ctx context.Context, adapterName string, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error) {
	adapter, err := s.adapterRegistry.Get(adapterName)
	if err != nil {
		return nil, err
	}

	// Sync inventory
	return adapter.SyncInventory(ctx, options)
}
