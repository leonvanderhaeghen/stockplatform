package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
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
		input.SKU = strings.ToUpper(uuid.New().String()[:8]) // Use first 8 chars of UUID as SKU
	} else {
		// Ensure SKU is uppercase and trimmed
		input.SKU = strings.TrimSpace(strings.ToUpper(input.SKU))
	}

	// Set timestamps
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now

	// Set default values
	if input.Currency == "" {
		input.Currency = "USD" // Default currency
	}

	// Initialize slices if nil
	if input.CategoryIDs == nil {
		input.CategoryIDs = []string{}
	}
	if input.ImageURLs == nil {
		input.ImageURLs = []string{}
	}
	if input.VideoURLs == nil {
		input.VideoURLs = []string{}
	}
	if input.Variants == nil {
		input.Variants = []domain.Variant{}
	}

	// Set stock status
	input.InStock = input.StockQty > 0

	// Set default visibility for the supplier
	if input.SupplierID != "" {
		if input.IsVisible == nil {
			input.IsVisible = make(map[string]bool)
		}
		// By default, the product is visible to its own supplier
		input.IsVisible[input.SupplierID] = true
	}

	// Validate the product
	if err := s.validateProduct(input); err != nil {
		s.logger.Error("Product validation failed", 
			zap.String("sku", input.SKU), 
			zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if product with same SKU or barcode already exists
	existing, _, err := s.repo.List(ctx, &domain.ListOptions{
		Search: input.SKU,
		Filters: map[string]interface{}{
			"$or": []bson.M{
				{"sku": input.SKU},
				{"barcode": input.Barcode},
			},
		},
		PageSize: 1,
	})

	if err != nil {
		s.logger.Error("Failed to check for existing product", 
			zap.String("sku", input.SKU), 
			zap.Error(err))
		return nil, fmt.Errorf("failed to check for existing product: %w", err)
	}

	if len(existing) > 0 {
		s.logger.Warn("Product with same SKU or barcode already exists", 
			zap.String("sku", input.SKU), 
			zap.String("barcode", input.Barcode))
		return nil, domain.ErrProductAlreadyExists
	}

	// Create the product
	product, err := s.repo.Create(ctx, input)
	if err != nil {
		s.logger.Error("Failed to create product", 
			zap.String("sku", input.SKU), 
			zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	s.logger.Info("Product created successfully", 
		zap.String("id", product.ID.Hex()),
		zap.String("sku", product.SKU))

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
		s.logger.Error("Failed to get product for update", 
			zap.String("id", id), 
			zap.Error(err))
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Preserve immutable fields
	input.ID = existing.ID
	input.CreatedAt = existing.CreatedAt
	input.SupplierID = existing.SupplierID // Supplier cannot be changed

	// Update timestamps
	input.UpdatedAt = time.Now()

	// Ensure SKU is uppercase and trimmed
	input.SKU = strings.TrimSpace(strings.ToUpper(input.SKU))

	// Validate the product
	if err := s.validateProduct(input); err != nil {
		s.logger.Error("Product validation failed", 
			zap.String("id", id), 
			zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check for duplicate SKU or barcode if they are being updated
	if input.SKU != existing.SKU || input.Barcode != existing.Barcode {
		filters := map[string]interface{}{
			"_id": bson.M{"$ne": existing.ID},
			"$or": []bson.M{
				{"sku": input.SKU},
			},
		}

		// Only check barcode if it's not empty
		if input.Barcode != "" {
			filters["$or"] = append(filters["$or"].([]bson.M), bson.M{"barcode": input.Barcode})
		}

		// Check for duplicates
		existingProducts, _, err := s.repo.List(ctx, &domain.ListOptions{
			Filters:  filters,
			PageSize: 1,
		})

		if err != nil {
			s.logger.Error("Failed to check for duplicate products", 
				zap.String("id", id), 
				zap.Error(err))
			return nil, fmt.Errorf("failed to check for duplicate products: %w", err)
		}

		if len(existingProducts) > 0 {
			s.logger.Warn("Product with same SKU or barcode already exists", 
				zap.String("id", id),
				zap.String("sku", input.SKU),
				zap.String("barcode", input.Barcode))
			return nil, domain.ErrProductAlreadyExists
		}
	}

	// Update variants timestamps
	now := time.Now()
	for i := range input.Variants {
		// If variant has no ID, it's a new variant
		if input.Variants[i].ID == "" {
			input.Variants[i].ID = uuid.New().String()
			input.Variants[i].CreatedAt = now
		}
		input.Variants[i].UpdatedAt = now

		// Update variant options timestamps
		for j := range input.Variants[i].Options {
			if input.Variants[i].Options[j].ID == "" {
				input.Variants[i].Options[j].ID = uuid.New().String()
				input.Variants[i].Options[j].CreatedAt = now
			}
			input.Variants[i].Options[j].UpdatedAt = now
		}
	}

	// Update stock status
	input.InStock = input.StockQty > 0

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

	// Check if product exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get product for deletion", 
			zap.String("id", id), 
			zap.Error(err))
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Soft delete the product (set deleted_at timestamp)
	err = s.repo.SoftDelete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete product", 
			zap.String("id", id), 
			zap.Error(err))
		return fmt.Errorf("failed to delete product: %w", err)
	}

	s.logger.Info("Product soft deleted successfully", 
		zap.String("id", id))

	return nil
}

// ListProducts retrieves a paginated and filtered list of products
func (s *ProductService) ListProducts(
	ctx context.Context,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	// Apply default options if nil
	if opts == nil {
		opts = &domain.ListOptions{
			Page:     1,
			PageSize: 10,
		}
	}

	// Ensure page and page size are valid
	if opts.Page < 1 {
		opts.Page = 1
	}

	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 10
	}

	// Add default filters if not provided
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}

	// Only include non-deleted products by default
	if _, exists := opts.Filters["deleted_at"]; !exists {
		opts.Filters["deleted_at"] = nil
	}

	// Get products from repository
	products, total, err := s.repo.List(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to list products", 
			zap.Any("filters", opts.Filters),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	s.logger.Debug("Listed products successfully",
		zap.Int("count", len(products)),
		zap.Int64("total", total))

	return products, total, nil
}

// validateProduct validates the product fields
func (s *ProductService) validateProduct(p *domain.Product) error {
	// Basic validations
	if p.Name == "" {
		return domain.ErrProductNameRequired
	}

	if p.SKU == "" {
		return domain.ErrProductSKURequired
	}

	// Validate prices
	if p.CostPrice != "" {
		costPrice, err := domain.ParsePrice(p.CostPrice)
		if err != nil || costPrice < 0 {
			return domain.ErrInvalidCostPrice
		}
	}

	if p.SellingPrice == "" {
		return domain.ErrSellingPriceRequired
	}

	sellingPrice, err := domain.ParsePrice(p.SellingPrice)
	if err != nil || sellingPrice < 0 {
		return domain.ErrInvalidSellingPrice
	}

	// Validate stock quantity
	if p.StockQty < 0 {
		return domain.ErrInvalidStockQuantity
	}

	// Validate currency (ISO 4217)
	if len(p.Currency) != 3 {
		return domain.ErrInvalidCurrency
	}

	// Validate variants
	variantSKUs := make(map[string]bool)
	for i, variant := range p.Variants {
		if variant.SKU == "" {
			return fmt.Errorf("variant at index %d: %w", i, domain.ErrVariantSKURequired)
		}

		// Check for duplicate SKUs in variants
		if _, exists := variantSKUs[variant.SKU]; exists {
			return fmt.Errorf("duplicate variant SKU: %s", variant.SKU)
		}
		variantSKUs[variant.SKU] = true

		// Validate variant options
		optionNames := make(map[string]bool)
		for j, option := range variant.Options {
			if option.Name == "" {
				return fmt.Errorf("variant %s: option at index %d: %w", 
					variant.SKU, j, domain.ErrOptionNameRequired)
			}

			// Check for duplicate option names
			if _, exists := optionNames[option.Name]; exists {
				return fmt.Errorf("variant %s: duplicate option name: %s", 
					variant.SKU, option.Name)
			}
			optionNames[option.Name] = true

			// Validate option values
			if option.Value == "" {
				return fmt.Errorf("variant %s: option %s: %w", 
					variant.SKU, option.Name, domain.ErrOptionValueRequired)
			}

			// Validate option price adjustment if present
			if option.PriceAdjustment != "" {
				if _, err := domain.ParsePrice(option.PriceAdjustment); err != nil {
					return fmt.Errorf("variant %s: option %s: %w: %v", 
						variant.SKU, option.Name, domain.ErrInvalidPriceAdjustment, err)
				}
			}
		}
	}

	return nil
}
