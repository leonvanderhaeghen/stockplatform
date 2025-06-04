package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"stockplatform/services/productSvc/internal/domain"
)

type productService struct {
	repo   domain.ProductRepository
	logger *zap.Logger
}

// NewProductService creates a new product service
func NewProductService(repo domain.ProductRepository, logger *zap.Logger) domain.ProductUseCase {
	return &productService{
		repo:   repo,
		logger: logger.With(zap.String("service", "product")),
	}
}

// CreateProduct creates a new product
func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Set default values
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	// Validate the product
	if err := s.ValidateProduct(product); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate IDs for variants and options
	s.generateVariantIDs(product)

	// Set default visibility if not provided
	if product.IsVisible == nil {
		product.IsVisible = make(map[string]bool)
	}
	if _, exists := product.IsVisible[product.SupplierID]; !exists {
		product.IsVisible[product.SupplierID] = true
	}

	// Set stock status
	product.InStock = product.StockQty > 0

	// Create the product in the repository
	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		s.logger.Error("failed to create product", zap.Error(err), zap.String("sku", product.SKU))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	s.logger.Info("product created successfully", zap.String("id", createdProduct.ID.Hex()))
	return createdProduct, nil
}

// GetProduct retrieves a product by ID
func (s *productService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	if id == "" {
		return nil, domain.ErrInvalidID
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get product", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if product.ID.IsZero() {
		return domain.ErrInvalidID
	}

	// Get the existing product
	existing, err := s.repo.GetByID(ctx, product.ID.Hex())
	if err != nil {
		return fmt.Errorf("failed to get existing product: %w", err)
	}

	// Preserve created_at and set updated_at
	product.CreatedAt = existing.CreatedAt
	product.UpdatedAt = time.Now()

	// Validate the product
	if err := s.ValidateProduct(product); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update the product in the repository
	if err := s.repo.Update(ctx, product); err != nil {
		s.logger.Error("failed to update product", zap.Error(err), zap.String("id", product.ID.Hex()))
		return fmt.Errorf("failed to update product: %w", err)
	}

	s.logger.Info("product updated successfully", zap.String("id", product.ID.Hex()))
	return nil
}

// DeleteProduct deletes a product
func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete product", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to delete product: %w", err)
	}

	s.logger.Info("product deleted successfully", zap.String("id", id))
	return nil
}

// SoftDeleteProduct soft deletes a product
func (s *productService) SoftDeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		s.logger.Error("failed to soft delete product", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to soft delete product: %w", err)
	}

	s.logger.Info("product soft deleted successfully", zap.String("id", id))
	return nil
}

// ListProducts retrieves a list of products with pagination and filtering
func (s *productService) ListProducts(ctx context.Context, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	products, total, err := s.repo.List(ctx, opts)
	if err != nil {
		s.logger.Error("failed to list products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, total, nil
}

// SearchProducts searches for products by query
func (s *productService) SearchProducts(ctx context.Context, query string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	if query == "" {
		return s.ListProducts(ctx, opts)
	}

	products, total, err := s.repo.Search(ctx, query, opts)
	if err != nil {
		s.logger.Error("failed to search products", zap.Error(err), zap.String("query", query))
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	return products, total, nil
}

// UpdateProductStock updates the stock quantity of a product
func (s *productService) UpdateProductStock(ctx context.Context, id string, quantity int32) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	if quantity < 0 {
		return domain.ErrInvalidQuantity
	}

	if err := s.repo.UpdateStock(ctx, id, quantity); err != nil {
		s.logger.Error("failed to update product stock", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	s.logger.Info("product stock updated successfully", zap.String("id", id), zap.Int32("quantity", quantity))
	return nil
}

// BulkUpdateProductVisibility updates the visibility of multiple products for a supplier
func (s *productService) BulkUpdateProductVisibility(ctx context.Context, supplierID string, productIDs []string, isVisible bool) error {
	if supplierID == "" {
		return domain.ErrInvalidSupplierID
	}

	if len(productIDs) == 0 {
		return domain.ErrNoProductsProvided
	}

	if err := s.repo.BulkUpdateVisibility(ctx, supplierID, productIDs, isVisible); err != nil {
		s.logger.Error("failed to bulk update product visibility", 
			zap.Error(err), 
			zap.String("supplierID", supplierID),
			zap.Bool("isVisible", isVisible))
		return fmt.Errorf("failed to bulk update product visibility: %w", err)
	}

	s.logger.Info("bulk product visibility updated successfully", 
		zap.String("supplierID", supplierID),
		zap.Int("count", len(productIDs)),
		zap.Bool("isVisible", isVisible))

	return nil
}

// ValidateProduct validates a product
func (s *productService) ValidateProduct(product *domain.Product) error {
	// Basic validation
	if product.Name == "" {
		return domain.ErrProductNameRequired
	}

	if product.SKU == "" {
		return domain.ErrProductSKURequired
	}

	// Normalize SKU
	product.SKU = strings.TrimSpace(strings.ToUpper(product.SKU))

	// Validate prices
	if _, err := domain.ParsePrice(product.SellingPrice); err != nil {
		return domain.ErrInvalidSellingPrice
	}

	if product.CostPrice != "" {
		if _, err := domain.ParsePrice(product.CostPrice); err != nil {
			return domain.ErrInvalidCostPrice
		}
	}

	// Validate currency (ISO 4217)
	if len(product.Currency) != 3 {
		return domain.ErrInvalidCurrency
	}

	// Validate variants
	variantNames := make(map[string]bool)
	variantSKUs := make(map[string]bool)

	for i, variant := range product.Variants {
		if variant.Name == "" {
			return fmt.Errorf("variant %d: %w", i, domain.ErrVariantNameRequired)
		}

		// Check for duplicate variant names
		if _, exists := variantNames[variant.Name]; exists {
			return fmt.Errorf("duplicate variant name: %s", variant.Name)
		}
		variantNames[variant.Name] = true

		// Check for duplicate variant SKUs
		if variant.SKU != "" {
			if _, exists := variantSKUs[variant.SKU]; exists {
				return fmt.Errorf("duplicate variant SKU: %s", variant.SKU)
			}
			variantSKUs[variant.SKU] = true
		}

		// Validate variant options
		optionValues := make(map[string]bool)
		for j, option := range variant.Options {
			if option.Name == "" {
				return fmt.Errorf("variant %s: option %d: %w", variant.Name, j, domain.ErrOptionNameRequired)
			}

			if option.Value == "" {
				return fmt.Errorf("variant %s: option %s: %w", variant.Name, option.Name, domain.ErrOptionValueRequired)
			}

			// Check for duplicate option values
			optionKey := fmt.Sprintf("%s:%s", option.Name, option.Value)
			if _, exists := optionValues[optionKey]; exists {
				return fmt.Errorf("duplicate option value for %s: %s", option.Name, option.Value)
			}
			optionValues[optionKey] = true

			// Validate price adjustment if present
			if option.PriceAdjustment != "" {
				if _, err := domain.ParsePrice(option.PriceAdjustment); err != nil {
					return fmt.Errorf("variant %s: option %s: %w: %v", 
						variant.Name, option.Name, domain.ErrInvalidPriceAdjustment, err)
				}
			}
		}
	}

	return nil
}

// generateVariantIDs generates unique IDs for variants and their options
func (s *productService) generateVariantIDs(product *domain.Product) {
	now := time.Now()

	for i := range product.Variants {
		// Generate variant ID if not set
		if product.Variants[i].ID == "" {
			product.Variants[i].ID = uuid.New().String()
		}

		// Set timestamps for variant
		if product.Variants[i].CreatedAt.IsZero() {
			product.Variants[i].CreatedAt = now
		}
		product.Variants[i].UpdatedAt = now

		// Process variant options
		for j := range product.Variants[i].Options {
			// Generate option ID if not set
			if product.Variants[i].Options[j].ID == "" {
				product.Variants[i].Options[j].ID = uuid.New().String()
			}

			// Set timestamps for option
			if product.Variants[i].Options[j].CreatedAt.IsZero() {
				product.Variants[i].Options[j].CreatedAt = now
			}
			product.Variants[i].Options[j].UpdatedAt = now
		}
	}
}

// The following methods are stubs that need to be implemented

func (s *productService) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	// TODO: Implement GetProductBySKU
	return nil, nil
}

func (s *productService) GetProductByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	// TODO: Implement GetProductByBarcode
	return nil, nil
}

func (s *productService) GetProductsBySupplier(ctx context.Context, supplierID string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// TODO: Implement GetProductsBySupplier
	return nil, 0, nil
}

func (s *productService) GetProductsByCategory(ctx context.Context, categoryID string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// TODO: Implement GetProductsByCategory
	return nil, 0, nil
}

func (s *productService) BulkUpdateStock(ctx context.Context, updates map[string]int32) error {
	// TODO: Implement BulkUpdateStock
	return nil
}

func (s *productService) AdjustStock(ctx context.Context, id string, adjustment int32, note string) error {
	// TODO: Implement AdjustStock
	return nil
}

func (s *productService) GetLowStockProducts(ctx context.Context, threshold int32, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// TODO: Implement GetLowStockProducts
	return nil, 0, nil
}

// UpdateProductPricing updates the pricing information for a product
func (s *productService) UpdateProductPricing(ctx context.Context, id string, costPrice, sellingPrice string) error {
	// Validate the product ID
	if id == "" {
		return domain.ErrInvalidID
	}

	// Get the existing product
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get existing product", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Update prices if provided
	if costPrice != "" {
		if _, err := domain.ParsePrice(costPrice); err != nil {
			return domain.ErrInvalidCostPrice
		}
		existing.CostPrice = costPrice
	}

	if sellingPrice != "" {
		if _, err := domain.ParsePrice(sellingPrice); err != nil {
			return domain.ErrInvalidSellingPrice
		}
		existing.SellingPrice = sellingPrice
	}

	// Validate the updated product
	if err := s.ValidateProduct(existing); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Update the product in the repository
	existing.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, existing); err != nil {
		s.logger.Error("failed to update product pricing", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to update product pricing: %w", err)
	}

	s.logger.Info("product pricing updated successfully", zap.String("id", id))
	return nil
}

// BulkUpdatePricing updates pricing information for multiple products
func (s *productService) BulkUpdatePricing(ctx context.Context, updates map[string]struct{ CostPrice, SellingPrice string }) error {
	if len(updates) == 0 {
		return nil // Nothing to update
	}

	// Process updates in batches to avoid overwhelming the database
	batchSize := 50
	batch := make(map[string]*domain.Product, batchSize)

	for productID, prices := range updates {
		// Get the existing product
		existing, err := s.repo.GetByID(ctx, productID)
		if err != nil {
			s.logger.Error("failed to get product for bulk update", 
				zap.Error(err), 
				zap.String("productID", productID))
			continue // Skip this product but continue with others
		}

		// Update prices if provided
		if prices.CostPrice != "" {
			if _, err := domain.ParsePrice(prices.CostPrice); err != nil {
				s.logger.Error("invalid cost price in bulk update",
					zap.String("productID", productID),
					zap.String("costPrice", prices.CostPrice),
					zap.Error(err))
				continue // Skip this product but continue with others
			}
			existing.CostPrice = prices.CostPrice
		}

		if prices.SellingPrice != "" {
			if _, err := domain.ParsePrice(prices.SellingPrice); err != nil {
				s.logger.Error("invalid selling price in bulk update",
					zap.String("productID", productID),
					zap.String("sellingPrice", prices.SellingPrice),
					zap.Error(err))
				continue // Skip this product but continue with others
			}
			existing.SellingPrice = prices.SellingPrice
		}

		// Update timestamps
		existing.UpdatedAt = time.Now()

		// Add to batch
		batch[productID] = existing

		// Process batch if it reaches the batch size
		if len(batch) >= batchSize {
			if err := s.processPricingBatch(ctx, batch); err != nil {
				s.logger.Error("error processing pricing batch", zap.Error(err))
			}
			// Reset batch
			batch = make(map[string]*domain.Product, batchSize)
		}
	}

	// Process any remaining items in the last batch
	if len(batch) > 0 {
		if err := s.processPricingBatch(ctx, batch); err != nil {
			s.logger.Error("error processing final pricing batch", zap.Error(err))
			return fmt.Errorf("error processing final batch: %w", err)
		}
	}

	s.logger.Info("bulk pricing update completed", 
		zap.Int("total_products", len(updates)),
		zap.Int("successful_updates", len(updates)-len(batch)))

	return nil
}

// processPricingBatch processes a batch of product pricing updates
func (s *productService) processPricingBatch(ctx context.Context, batch map[string]*domain.Product) error {
	// In a real implementation, you would update all products in a single database operation
	// For now, we'll update them one by one
	for productID, product := range batch {
		if err := s.repo.Update(ctx, product); err != nil {
			s.logger.Error("failed to update product in batch",
				zap.String("productID", productID),
				zap.Error(err))
			// Continue with other products even if one fails
			continue
		}
	}
	return nil
}

func (s *productService) CalculateProfitMargin(productID string) (decimal.Decimal, decimal.Decimal, error) {
	// TODO: Implement CalculateProfitMargin
	return decimal.Zero, decimal.Zero, nil
}

func (s *productService) CalculateMarkup(productID string) (decimal.Decimal, decimal.Decimal, error) {
	// TODO: Implement CalculateMarkup
	return decimal.Zero, decimal.Zero, nil
}

func (s *productService) AddVariant(ctx context.Context, productID string, variant *domain.Variant) error {
	// TODO: Implement AddVariant
	return nil
}

func (s *productService) UpdateVariant(ctx context.Context, productID string, variant *domain.Variant) error {
	// TODO: Implement UpdateVariant
	return nil
}

func (s *productService) RemoveVariant(ctx context.Context, productID, variantID string) error {
	// TODO: Implement RemoveVariant
	return nil
}

func (s *productService) UpdateVariantStock(ctx context.Context, productID, variantID string, quantity int32) error {
	// TODO: Implement UpdateVariantStock
	return nil
}

func (s *productService) PublishProducts(ctx context.Context, productIDs []string, publish bool) error {
	// TODO: Implement PublishProducts
	return nil
}

func (s *productService) GenerateProductReport(ctx context.Context, format string) ([]byte, error) {
	// TODO: Implement GenerateProductReport
	return nil, nil
}
