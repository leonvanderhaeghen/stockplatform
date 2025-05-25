package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"stockplatform/services/productSvc/internal/domain"
)

// ProductService implements the business logic for product operations
type ProductService struct {
	repo   domain.ProductRepository
	logger *zap.Logger
}

// NewProductService creates a new product service
func NewProductService(repo domain.ProductRepository, logger *zap.Logger) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logger.Named("product_service"),
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, input *domain.Product) (*domain.Product, error) {
	// Generate a new UUID for the product if not provided
	if input.SKU == "" {
		input.SKU = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now

	// Validate the product
	if err := validateProduct(input); err != nil {
		s.logger.Error("Product validation failed", zap.Error(err))
		return nil, domain.ErrValidation
	}

	// Check if product with same SKU already exists
	existing, _, err := s.repo.List(ctx, map[string]interface{}{"sku": input.SKU}, 1, 0)
	if err != nil {
		s.logger.Error("Failed to check for existing product", zap.Error(err))
		return nil, err
	}
	if len(existing) > 0 {
		return nil, domain.ErrAlreadyExists
	}

	// Create the product
	product, err := s.repo.Create(ctx, input)
	if err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Product created successfully", zap.String("id", product.ID.Hex()))
	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get product", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id string, input *domain.Product) (*domain.Product, error) {
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	// Get existing product
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	existing.Name = input.Name
	existing.Description = input.Description
	existing.Price = input.Price
	existing.CategoryID = input.CategoryID
	existing.ImageURLs = input.ImageURLs
	existing.UpdatedAt = time.Now()

	// Validate the updated product
	if err := validateProduct(existing); err != nil {
		s.logger.Error("Product validation failed", zap.Error(err))
		return nil, domain.ErrValidation
	}

	// Update the product
	if err := s.repo.Update(ctx, existing); err != nil {
		s.logger.Error("Failed to update product", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Product updated successfully", zap.String("id", id))
	return existing, nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete product", zap.String("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Product deleted successfully", zap.String("id", id))
	return nil
}

// ListProducts retrieves a paginated and filtered list of products
func (s *ProductService) ListProducts(
	ctx context.Context,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	// Set default values if options are nil
	if opts == nil {
		opts = &domain.ListOptions{}
	}

	// Initialize pagination if not provided
	if opts.Pagination == nil {
		opts.Pagination = &domain.Pagination{
			Page:     1,
			PageSize: 10,
		}
	}

	// Validate pagination
	if opts.Pagination.Page < 1 {
		opts.Pagination.Page = 1
	}

	if opts.Pagination.PageSize < 1 || opts.Pagination.PageSize > 100 {
		opts.Pagination.PageSize = 10
	}

	// Execute the query
	products, total, err := s.repo.List(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to list products", 
			zap.Any("options", opts), 
			zap.Error(err))
		return nil, 0, err
	}

	if products == nil {
		products = []*domain.Product{}
	}

	return products, total, nil
}

// validateProduct validates the product fields
func validateProduct(p *domain.Product) error {
	if p.Name == "" {
		return domain.ErrValidation
	}

	if p.Price < 0 {
		return domain.ErrValidation
	}

	if p.CategoryID == "" {
		return domain.ErrValidation
	}

	return nil
}
