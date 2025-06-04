package application

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
)

// CategoryService implements the business logic for category operations
type CategoryService struct {
	repo   domain.CategoryRepository
	logger *zap.Logger
}

// NewCategoryService creates a new category service
func NewCategoryService(repo domain.CategoryRepository, logger *zap.Logger) *CategoryService {
	return &CategoryService{
		repo:   repo,
		logger: logger.Named("category_service"),
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	// Set timestamps
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	// If parent ID is provided, get the parent and set the level and path
	if category.ParentID != "" {
		parent, err := s.repo.GetByID(ctx, category.ParentID)
		if err != nil {
			s.logger.Error("Failed to get parent category", zap.Error(err))
			return nil, domain.ErrParentCategoryNotFound
		}

		category.Level = parent.Level + 1
		if parent.Path != "" {
			category.Path = parent.Path + "/" + parent.ID.Hex()
		} else {
			category.Path = parent.ID.Hex()
		}
	} else {
		category.Level = 0
		category.Path = ""
	}

	return s.repo.Create(ctx, category)
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) error {
	// Get the existing category to preserve some fields
	existing, err := s.repo.GetByID(ctx, category.ID.Hex())
	if err != nil {
		return err
	}

	// Preserve created_at and update updated_at
	category.CreatedAt = existing.CreatedAt
	category.UpdatedAt = time.Now()

	return s.repo.Update(ctx, category)
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	// Check if category has any products
	// This would require a product repository to check
	// For now, we'll just try to delete and let the repository handle any constraints

	return s.repo.Delete(ctx, id)
}

// ListCategories retrieves a list of categories, optionally filtered by parent ID and depth
func (s *CategoryService) ListCategories(ctx context.Context, parentID string, depth int32) ([]*domain.Category, error) {
	return s.repo.List(ctx, parentID, depth)
}
