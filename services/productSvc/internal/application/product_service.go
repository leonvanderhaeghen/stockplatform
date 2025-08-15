package application

import (
	"context"

	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
	supplierclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/supplier"
	inventoryclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/inventory"
)

// ProductService implements the business logic for product operations
type ProductService struct {
	repo           domain.ProductRepository
	supplierClient *supplierclient.Client
	inventoryClient *inventoryclient.Client
	logger         *zap.Logger
}

// Ensure ProductService implements ProductUseCase
var _ domain.ProductUseCase = (*ProductService)(nil)

// GetLowStockProducts retrieves products with stock below the specified threshold
func (s *ProductService) GetLowStockProducts(ctx context.Context, threshold int32, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	if threshold < 0 {
		return nil, 0, fmt.Errorf("threshold must be non-negative")
	}

	// Set default options if nil
	if opts == nil {
		opts = &domain.ListOptions{
			Pagination: &domain.Pagination{
				Page:     1,
				PageSize: 20,
			},
		}
	}

	// Get low stock products from repository
	products, total, err := s.repo.GetLowStockProducts(ctx, threshold, opts)
	if err != nil {
		s.logger.Error("Failed to get low stock products",
			zap.Int32("threshold", threshold),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to get low stock products: %w", err)
	}

	s.logger.Debug("Retrieved low stock products",
		zap.Int("count", len(products)),
		zap.Int64("total", total))

	return products, total, nil
}

// ListProducts retrieves a paginated list of products with optional filtering
func (s *ProductService) ListProducts(ctx context.Context, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	if opts == nil {
		opts = &domain.ListOptions{
			Pagination: &domain.Pagination{
				Page:     1,
				PageSize: 20, // Default page size
			},
		}
	}

	// Initialize pagination if not provided
	if opts.Pagination == nil {
		opts.Pagination = &domain.Pagination{
			Page:     1,
			PageSize: 20,
		}
	}

	// Ensure page and page size are within reasonable bounds
	if opts.Pagination.Page < 1 {
		opts.Pagination.Page = 1
	}
	if opts.Pagination.PageSize < 1 || opts.Pagination.PageSize > 100 {
		opts.Pagination.PageSize = 20
	}

	// Call the repository to get the paginated list of products
	products, total, err := s.repo.List(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to list products", 
			zap.Any("options", opts), 
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	s.logger.Info("Listed products successfully", 
		zap.Int("count", len(products)),
		zap.Int64("total", total))

	return products, total, nil
}

// NewProductService creates a new product service
func NewProductService(repo domain.ProductRepository, supplierClient *supplierclient.Client, inventoryClient *inventoryclient.Client, logger *zap.Logger) *ProductService {
	return &ProductService{
		repo:           repo,
		supplierClient: supplierClient,
		inventoryClient: inventoryClient,
		logger:         logger.Named("product_service"),
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, input *domain.Product) (*domain.Product, error) {
	// Validate supplier exists
	if input.SupplierID == "" {
		return nil, domain.ErrSupplierRequired
	}

	// Use client abstraction instead of direct protobuf
	_, err := s.supplierClient.GetSupplier(ctx, input.SupplierID)
	if err != nil {
		s.logger.Error("Invalid supplier ID", 
			zap.String("supplierID", input.SupplierID), 
			zap.Error(err))
		return nil, domain.ErrSupplierNotFound
	}

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



	// Set default visibility for the supplier
	if input.SupplierID != "" {
		if input.IsVisible == nil {
			input.IsVisible = make(map[string]bool)
		}
		// By default, the product is visible to its own supplier
		input.IsVisible[input.SupplierID] = true
	}

	// Validate the product
	if err := s.ValidateProduct(input); err != nil {
		s.logger.Error("Product validation failed", 
			zap.String("sku", input.SKU), 
			zap.Error(err))
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Rely on MongoDB unique indexes for duplicate SKU / barcode detection.
	// Attempting a pre-check with text search fails before the collection and its text index exist.
	// Duplicate errors will surface as mongo.IsDuplicateKeyError during the insert below.

	// Create the product
	product, err := s.repo.Create(ctx, input)
	if err != nil {
		s.logger.Error("Failed to create product", 
			zap.String("sku", input.SKU), 
			zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Create inventory item for the product using client abstraction
	_, err = s.inventoryClient.CreateInventory(ctx, product.ID.Hex(), product.SKU, "default", 0)
	if err != nil {
		s.logger.Error("Failed to create inventory item", 
			zap.String("product_id", product.ID.Hex()),
			zap.Error(err))
		// Don't fail product creation if inventory creation fails
		// Log the error and continue
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
func (s *ProductService) UpdateProduct(ctx context.Context, input *domain.Product) error {
	if input == nil || input.ID.IsZero() {
		return fmt.Errorf("invalid product")
	}

	id := input.ID.Hex()

	// Get existing product
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get product for update", 
			zap.String("id", id), 
			zap.Error(err))
		return fmt.Errorf("failed to get product: %w", err)
	}

	// If supplier ID is being updated, validate the new supplier exists
	if input.SupplierID != "" && input.SupplierID != existing.SupplierID {
		// Use client abstraction instead of direct protobuf
		_, err := s.supplierClient.GetSupplier(ctx, input.SupplierID)
		if err != nil {
			s.logger.Error("Invalid supplier ID", 
				zap.String("supplierID", input.SupplierID), 
				zap.Error(err))
			return domain.ErrSupplierNotFound
		}
		existing.SupplierID = input.SupplierID
	} else {
		// Preserve existing supplier ID if not being changed
		input.SupplierID = existing.SupplierID
	}

	// Preserve immutable fields
	input.ID = existing.ID
	input.CreatedAt = existing.CreatedAt

	// Update timestamps
	input.UpdatedAt = time.Now()

	// Ensure SKU is uppercase and trimmed
	input.SKU = strings.TrimSpace(strings.ToUpper(input.SKU))

	// Validate the product
	if err := s.ValidateProduct(input); err != nil {
		s.logger.Error("Product validation failed", 
			zap.String("id", id), 
			zap.Error(err))
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check for duplicate SKU or barcode if they are being updated
	if input.SKU != existing.SKU || input.Barcode != existing.Barcode {
		// Check for duplicates using the repository's search functionality
		existingProducts, _, err := s.repo.Search(ctx, input.SKU, &domain.ListOptions{
			Pagination: &domain.Pagination{
				Page:     1,
				PageSize: 1,
			},
		})

		if err != nil {
			s.logger.Error("Failed to check for duplicate products", 
				zap.String("id", id), 
				zap.Error(err))
			return fmt.Errorf("failed to check for duplicate products: %w", err)
		}

		// Check if any product with the same SKU or barcode exists (excluding current product)
		for _, p := range existingProducts {
			// Convert both IDs to strings for comparison
			existingID := p.ID.Hex()
			if existingID != id && (p.SKU == input.SKU || (input.Barcode != "" && p.Barcode == input.Barcode)) {
				s.logger.Warn("Product with same SKU or barcode already exists", 
					zap.String("id", id),
					zap.String("existing_id", existingID),
					zap.String("sku", input.SKU),
					zap.String("barcode", input.Barcode))
				return fmt.Errorf("product with same SKU or barcode already exists")
			}
		}
	}

	// Update variants timestamps
	now := time.Now()
	for i := range input.Variants {
		// If variant has no ID, it's a new variant
		if input.Variants[i].ID == "" {
			input.Variants[i].CreatedAt = now
		}
		input.Variants[i].UpdatedAt = now
	}

	// Update product fields
	existing.Name = input.Name
	existing.Description = input.Description
	existing.CostPrice = input.CostPrice
	existing.SellingPrice = input.SellingPrice
	existing.Currency = input.Currency
	existing.Barcode = input.Barcode
	existing.CategoryIDs = input.CategoryIDs
	existing.IsActive = input.IsActive
	existing.ImageURLs = input.ImageURLs
	existing.VideoURLs = input.VideoURLs
	existing.Metadata = input.Metadata

	// Update timestamps
	existing.UpdatedAt = now

	// Update the product in the repository
	err = s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.Error("Failed to update product", 
			zap.String("id", id), 
			zap.Error(err))
		return fmt.Errorf("failed to update product: %w", err)
	}

	s.logger.Info("Product updated successfully", 
		zap.String("id", id),
		zap.String("sku", input.SKU))

	return nil
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

// List retrieves a paginated list of products with optional filtering
func (s *ProductService) List(
	ctx context.Context,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	// Set default values
	if opts == nil {
		opts = &domain.ListOptions{
			Pagination: &domain.Pagination{
				Page:     1,
				PageSize: 20,
			},
		}
	}

	// Initialize pagination if nil
	if opts.Pagination == nil {
		opts.Pagination = &domain.Pagination{
			Page:     1,
			PageSize: 20,
		}
	}

	// Validate pagination
	if opts.Pagination.Page < 1 {
		opts.Pagination.Page = 1
	}
	if opts.Pagination.PageSize < 1 || opts.Pagination.PageSize > 100 {
		opts.Pagination.PageSize = 20
	}

	// Call repository
	products, total, err := s.repo.List(ctx, opts)
	if err != nil {
		s.logger.Error("Failed to list products",
			zap.Any("options", opts),
			zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	s.logger.Debug("Listed products successfully",
		zap.Int("count", len(products)),
		zap.Int64("total", total))

	return products, total, nil
}

// SoftDeleteProduct marks a product as deleted without removing it from the database
func (s *ProductService) SoftDeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	// Check if product exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(ctx, id)
}

// UpdateProductStock updates the stock quantity for a product
func (s *ProductService) UpdateProductStock(ctx context.Context, id string, quantity int32) error {
	if id == "" {
		return domain.ErrInvalidID
	}
	if quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	// Check if product exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.UpdateStock(ctx, id, quantity)
}

// BulkUpdateStock updates stock quantities for multiple products
func (s *ProductService) BulkUpdateStock(ctx context.Context, updates map[string]int32) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	// Validate all quantities are non-negative
	for productID, qty := range updates {
		if qty < 0 {
			return fmt.Errorf("negative quantity for product %s", productID)
		}
	}

	return s.repo.BulkUpdateStock(ctx, updates)
}

// AdjustStock adjusts the stock quantity for a product by a delta value
func (s *ProductService) AdjustStock(ctx context.Context, id string, adjustment int32, note string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	// Inventory operations are handled by inventorySvc
	return fmt.Errorf("inventory operations are handled by inventorySvc")
}

// UpdateProductPricing updates the pricing for a product
func (s *ProductService) UpdateProductPricing(ctx context.Context, id string, costPrice, sellingPrice string) error {
	if id == "" {
		return domain.ErrInvalidID
	}

	// Parse and validate prices
	cost, err := domain.ParsePrice(costPrice)
	if err != nil {
		return fmt.Errorf("invalid cost price: %w", err)
	}

	selling, err := domain.ParsePrice(sellingPrice)
	if err != nil {
		return fmt.Errorf("invalid selling price: %w", err)
	}

	if selling.LessThan(cost) {
		return fmt.Errorf("selling price cannot be less than cost price")
	}

	// Get existing product
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update prices
	product.CostPrice = costPrice
	product.SellingPrice = sellingPrice

	return s.repo.Update(ctx, product)
}

// BulkUpdatePricing updates pricing for multiple products
func (s *ProductService) BulkUpdatePricing(ctx context.Context, updates map[string]struct{ CostPrice, SellingPrice string }) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	// Process updates in a transaction
	for productID, prices := range updates {
		err := s.UpdateProductPricing(ctx, productID, prices.CostPrice, prices.SellingPrice)
		if err != nil {
			return fmt.Errorf("failed to update product %s: %w", productID, err)
		}
	}

	return nil
}

// CalculateProfitMargin calculates the profit margin for a product
func (s *ProductService) CalculateProfitMargin(productID string) (decimal.Decimal, decimal.Decimal, error) {
	if productID == "" {
		return decimal.Zero, decimal.Zero, domain.ErrInvalidID
	}

	product, err := s.repo.GetByID(context.Background(), productID)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	costPrice, err := domain.ParsePrice(product.CostPrice)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("invalid cost price: %w", err)
	}

	sellingPrice, err := domain.ParsePrice(product.SellingPrice)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("invalid selling price: %w", err)
	}

	return domain.CalculateProfitMargin(costPrice, sellingPrice)
}

// CalculateMarkup calculates the markup percentage for a product
func (s *ProductService) CalculateMarkup(productID string) (decimal.Decimal, decimal.Decimal, error) {
	if productID == "" {
		return decimal.Zero, decimal.Zero, domain.ErrInvalidID
	}

	product, err := s.repo.GetByID(context.Background(), productID)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	costPrice, err := domain.ParsePrice(product.CostPrice)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("invalid cost price: %w", err)
	}

	sellingPrice, err := domain.ParsePrice(product.SellingPrice)
	if err != nil {
		return decimal.Zero, decimal.Zero, fmt.Errorf("invalid selling price: %w", err)
	}

	return domain.CalculateMarkup(costPrice, sellingPrice)
}

// AddVariant adds a variant to a product
func (s *ProductService) AddVariant(ctx context.Context, productID string, variant *domain.Variant) error {
	if productID == "" {
		return domain.ErrInvalidID
	}
	if variant == nil {
		return fmt.Errorf("variant is required")
	}

	// Get existing product
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Initialize variants slice if nil
	if product.Variants == nil {
		product.Variants = []domain.Variant{}
	}

	// Set timestamps
	now := time.Now()
	variant.ID = primitive.NewObjectID().Hex()
	variant.CreatedAt = now
	variant.UpdatedAt = now

	// Add variant
	product.Variants = append(product.Variants, *variant)

	return s.repo.Update(ctx, product)
}

// UpdateVariant updates a variant for a product
func (s *ProductService) UpdateVariant(ctx context.Context, productID string, variant *domain.Variant) error {
	if productID == "" || variant == nil || variant.ID == "" {
		return fmt.Errorf("product ID and variant ID are required")
	}

	// Get existing product
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Find and update variant
	found := false
	for i, v := range product.Variants {
		if v.ID == variant.ID {
			variant.CreatedAt = v.CreatedAt
			variant.UpdatedAt = time.Now()
			product.Variants[i] = *variant
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("variant not found")
	}

	return s.repo.Update(ctx, product)
}

// RemoveVariant removes a variant from a product
func (s *ProductService) RemoveVariant(ctx context.Context, productID, variantID string) error {
	if productID == "" || variantID == "" {
		return fmt.Errorf("product ID and variant ID are required")
	}

	// Get existing product
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Find and remove variant
	found := false
	for i, v := range product.Variants {
		if v.ID == variantID {
			// Remove variant by slicing
			product.Variants = append(product.Variants[:i], product.Variants[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("variant not found")
	}

	return s.repo.Update(ctx, product)
}

// UpdateVariantStock updates the stock for a specific variant
func (s *ProductService) UpdateVariantStock(ctx context.Context, productID, variantID string, quantity int32) error {
	if productID == "" || variantID == "" {
		return fmt.Errorf("product ID and variant ID are required")
	}

	// Get existing product
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Find variant and update its stock
	variantFound := false
	for i, v := range product.Variants {
		if v.ID == variantID {
			// Update the variant's updated_at timestamp
			product.Variants[i].UpdatedAt = time.Now()
			variantFound = true
			break
		}
	}

	if !variantFound {
		return fmt.Errorf("variant not found")
	}

	// Update the product's stock quantity
	// Inventory operations are handled by inventorySvc

	// Update the product with modified variant
	return s.repo.Update(ctx, product)
}

// BulkUpdateProductVisibility updates visibility for multiple products
func (s *ProductService) BulkUpdateProductVisibility(ctx context.Context, supplierID string, productIDs []string, isVisible bool) error {
	if supplierID == "" {
		return domain.ErrSupplierRequired
	}
	if len(productIDs) == 0 {
		return fmt.Errorf("no product IDs provided")
	}

	// Check if supplier exists using client abstraction
	_, err := s.supplierClient.GetSupplier(ctx, supplierID)
	if err != nil {
		s.logger.Error("Invalid supplier ID",
			zap.String("supplierID", supplierID),
			zap.Error(err))
		return domain.ErrSupplierNotFound
	}

	return s.repo.BulkUpdateVisibility(ctx, supplierID, productIDs, isVisible)
}

// PublishProducts publishes or unpublishes multiple products
func (s *ProductService) PublishProducts(ctx context.Context, productIDs []string, publish bool) error {
	if len(productIDs) == 0 {
		return fmt.Errorf("no product IDs provided")
	}

	return s.repo.PublishProducts(ctx, productIDs, publish)
}

// SearchProducts searches for products by query
func (s *ProductService) SearchProducts(
	ctx context.Context,
	query string,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	if query == "" {
		return s.List(ctx, opts)
	}
	return s.repo.Search(ctx, query, opts)
}

// GetProductBySKU retrieves a product by SKU
func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	if sku == "" {
		return nil, fmt.Errorf("SKU is required")
	}
	products, _, err := s.repo.Search(ctx, sku, &domain.ListOptions{
		Pagination: &domain.Pagination{
			Page:     1,
			PageSize: 1,
		},
	})
	if err != nil {
		s.logger.Error("Failed to search product by SKU",
			zap.String("sku", sku),
			zap.Error(err))
		return nil, fmt.Errorf("failed to search product by SKU: %w", err)
	}

	if len(products) == 0 {
		return nil, domain.ErrProductNotFound
	}

	return products[0], nil
}

// GetProductByBarcode retrieves a product by barcode
func (s *ProductService) GetProductByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	if barcode == "" {
		return nil, fmt.Errorf("barcode is required")
	}

	products, _, err := s.repo.Search(ctx, barcode, &domain.ListOptions{
		Pagination: &domain.Pagination{
			Page:     1,
			PageSize: 1,
		},
	})

	if err != nil {
		s.logger.Error("Failed to search product by barcode",
			zap.String("barcode", barcode),
			zap.Error(err))
		return nil, fmt.Errorf("failed to search product by barcode: %w", err)
	}

	if len(products) == 0 {
		return nil, domain.ErrProductNotFound
	}

	return products[0], nil
}

// GetProductsBySupplier retrieves products by supplier ID
func (s *ProductService) GetProductsBySupplier(
	ctx context.Context,
	supplierID string,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	if supplierID == "" {
		return nil, 0, domain.ErrSupplierRequired
	}

	// Check if supplier exists using client abstraction
	_, err := s.supplierClient.GetSupplier(ctx, supplierID)
	if err != nil {
		s.logger.Error("Invalid supplier ID",
			zap.String("supplierID", supplierID),
			zap.Error(err))
		return nil, 0, domain.ErrSupplierNotFound
	}

	return s.repo.GetBySupplier(ctx, supplierID, opts)
}

// GetProductsByCategory retrieves products by category ID
func (s *ProductService) GetProductsByCategory(
	ctx context.Context,
	categoryID string,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	if categoryID == "" {
		return nil, 0, fmt.Errorf("category ID is required")
	}

	// Note: We assume category existence is validated by the category service
	return s.repo.GetByCategory(ctx, categoryID, opts)
}

// ValidateProduct validates the product fields
func (s *ProductService) ValidateProduct(p *domain.Product) error {
	// Basic validations
	if p.Name == "" {
		return fmt.Errorf("product name is required")
	}

	if p.SKU == "" {
		return fmt.Errorf("product SKU is required")
	}

	// Validate prices
	if p.CostPrice != "" {
		_, err := domain.ParsePrice(p.CostPrice)
		if err != nil {
			return fmt.Errorf("invalid cost price: %w", err)
		}
	}

	if p.SellingPrice == "" {
		return fmt.Errorf("selling price is required")
	}

	_, err := domain.ParsePrice(p.SellingPrice)
	if err != nil {
		return fmt.Errorf("invalid selling price: %w", err)
	}

	// Validate stock quantity
	// Inventory validation is handled by inventorySvc

	// Validate currency (ISO 4217)
	if p.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if len(p.Currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter ISO 4217 code")
	}

	// Validate variants
	variantNames := make(map[string]bool)
	for i, variant := range p.Variants {
		if variant.Name == "" {
			return fmt.Errorf("variant at index %d: name is required", i)
		}

		// Check for duplicate variant names
		if _, exists := variantNames[variant.Name]; exists {
			return fmt.Errorf("duplicate variant name: %s", variant.Name)
		}
		variantNames[variant.Name] = true

		// Validate variant options
		optionNames := make(map[string]bool)
		for j, option := range variant.Options {
			if option.Name == "" {
				return fmt.Errorf("variant %s: option at index %d: name is required", 
					variant.Name, j)
			}

			// Check for duplicate option names
			if _, exists := optionNames[option.Name]; exists {
				return fmt.Errorf("variant %s: duplicate option name: %s", 
					variant.Name, option.Name)
			}
			optionNames[option.Name] = true

			// Validate option values
			if option.Value == "" {
				return fmt.Errorf("variant %s: option %s: value is required", 
					variant.Name, option.Name)
			}

			// Validate option price adjustment if present
			if option.PriceAdjustment != "" {
				if _, err := domain.ParsePrice(option.PriceAdjustment); err != nil {
					return fmt.Errorf("variant %s: option %s: invalid price adjustment: %w", 
						variant.Name, option.Name, err)
				}
			}
		}
	}

	return nil
}

// GenerateProductReport generates a report for products
func (s *ProductService) GenerateProductReport(ctx context.Context, format string) ([]byte, error) {
	// This is a placeholder implementation
	// In a real implementation, this would generate a report in the specified format (e.g., CSV, PDF, Excel)
	// containing product details, stock levels, etc.
	return []byte("Product report in " + format + " format"), nil
}
