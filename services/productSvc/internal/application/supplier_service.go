package application

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
)

// SupplierService implements the business logic for supplier operations
type SupplierService struct {
	repo   domain.SupplierRepository
	logger *zap.Logger
}

// NewSupplierService creates a new supplier service
func NewSupplierService(repo domain.SupplierRepository, logger *zap.Logger) *SupplierService {
	return &SupplierService{
		repo:   repo,
		logger: logger.Named("supplier_service"),
	}
}

// CreateSupplier creates a new supplier
func (s *SupplierService) CreateSupplier(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error) {
	// Set timestamps
	now := time.Now()
	supplier.CreatedAt = now
	supplier.UpdatedAt = now

	// Validate required fields
	if supplier.Name == "" {
		return nil, domain.ErrSupplierNameRequired
	}

	// Create the supplier
	createdSupplier, err := s.repo.Create(ctx, supplier)
	if err != nil {
		s.logger.Error("Failed to create supplier", zap.Error(err))
		return nil, err
	}

	return createdSupplier, nil
}

// GetSupplier retrieves a supplier by ID
func (s *SupplierService) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get supplier", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return supplier, nil
}

// UpdateSupplier updates an existing supplier
func (s *SupplierService) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	if supplier.ID.IsZero() {
		return domain.ErrInvalidID
	}

	// Get existing supplier to ensure it exists
	existing, err := s.repo.GetByID(ctx, supplier.ID.Hex())
	if err != nil {
		return err
	}

	// Update timestamps
	supplier.UpdatedAt = time.Now()
	supplier.CreatedAt = existing.CreatedAt // Preserve created at

	return s.repo.Update(ctx, supplier)
}

// DeleteSupplier deletes a supplier by ID
func (s *SupplierService) DeleteSupplier(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	// Check if supplier exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check if supplier is being used by any products before deleting

	return s.repo.Delete(ctx, id)
}

// ListSuppliers retrieves a list of suppliers with pagination and optional filtering
func (s *SupplierService) ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*domain.Supplier, int32, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	suppliers, total, err := s.repo.List(ctx, page, pageSize, search)
	if err != nil {
		s.logger.Error("Failed to list suppliers", zap.Error(err))
		return nil, 0, err
	}

	return suppliers, total, nil
}
