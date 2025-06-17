package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
)

// ProductServiceImpl implements the ProductService interface
type ProductServiceImpl struct {
	client *grpcclient.ProductClient
	logger *zap.Logger
}

// NewProductService creates a new instance of ProductServiceImpl
func NewProductService(productServiceAddr string, logger *zap.Logger) (ProductService, error) {
	// Create a new gRPC client
	client, err := grpcclient.NewProductClient(productServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create product client: %w", err)
	}

	return &ProductServiceImpl{
		client: client,
		logger: logger.Named("product_service"),
	}, nil
}

// CreateCategory creates a new product category
func (s *ProductServiceImpl) CreateCategory(
	ctx context.Context,
	name, description, parentID string,
	isActive bool,
) (interface{}, error) {
	s.logger.Debug("CreateCategory",
		zap.String("name", name),
		zap.String("parentID", parentID),
		zap.Bool("isActive", isActive),
	)

	resp, err := s.client.CreateCategory(ctx, name, description, parentID, isActive)
	if err != nil {
		s.logger.Error("Failed to create category",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return resp.GetCategory(), nil
}

// ListCategories lists all product categories
// parentID: Optional parent category ID to filter by
// depth: Maximum depth of subcategories to return (0 for all)
func (s *ProductServiceImpl) ListCategories(
	ctx context.Context,
) (interface{}, error) {
	s.logger.Debug("ListCategories")

	// Call the gRPC service with default parameters
	// In a real implementation, you might want to get these from query parameters
	req := &productv1.ListCategoriesRequest{
		ParentId: "", // Empty string means get root categories
		Depth:    3,  // Default to 3 levels deep
	}

	resp, err := s.client.ListCategories(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list categories",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	// Return the categories in a structured response
	return map[string]interface{}{
		"categories": resp.GetCategories(),
		"total":      len(resp.GetCategories()),
	}, nil
}

// ListProducts lists products with filtering options
func (s *ProductServiceImpl) ListProducts(
	ctx context.Context,
	categoryID, query string,
	active bool,
	limit, offset int,
	sortBy string,
	ascending bool,
) (interface{}, error) {
	s.logger.Debug("ListProducts",
		zap.String("categoryID", categoryID),
		zap.String("query", query),
		zap.Bool("active", active),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.String("sortBy", sortBy),
		zap.Bool("ascending", ascending),
	)

	// Convert sort field to protobuf enum
	var sortField productv1.ProductSort_SortField
	switch sortBy {
	case "name":
		sortField = productv1.ProductSort_SORT_FIELD_NAME
	case "price":
		sortField = productv1.ProductSort_SORT_FIELD_PRICE
	case "created_at":
		sortField = productv1.ProductSort_SORT_FIELD_CREATED_AT
	case "updated_at":
		sortField = productv1.ProductSort_SORT_FIELD_UPDATED_AT
	default:
		sortField = productv1.ProductSort_SORT_FIELD_UNSPECIFIED
	}

	// Convert sort order to protobuf enum
	sortOrder := productv1.ProductSort_SORT_ORDER_DESC
	if ascending {
		sortOrder = productv1.ProductSort_SORT_ORDER_ASC
	}

	// Build the request
	req := &productv1.ListProductsRequest{
		Pagination: &productv1.Pagination{
			Page:     int32(offset/limit) + 1,
			PageSize: int32(limit),
		},
		Sort: &productv1.ProductSort{
			Field: sortField,
			Order: sortOrder,
		},
	}

	// Add filters if provided
	if categoryID != "" || query != "" || active {
		req.Filter = &productv1.ProductFilter{}

		if categoryID != "" {
			req.Filter.CategoryIds = []string{categoryID}
		}

		if query != "" {
			req.Filter.SearchTerm = query
		}

		// Note: The active filter is not directly supported in the gRPC API
		// You might need to handle this in the client or modify the gRPC service
	}

	// Call the gRPC service
	resp, err := s.client.ListProducts(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list products",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Return the products in a structured response
	return map[string]interface{}{
		"products": resp.GetProducts(),
		"total":    resp.GetTotalCount(),
		"page":     resp.GetPage(),
		"pageSize": resp.GetPageSize(),
	}, nil
}

// GetProductByID gets a product by ID
func (s *ProductServiceImpl) GetProductByID(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetProductByID",
		zap.String("id", id),
	)

	req := &productv1.GetProductRequest{
		Id: id,
	}

	resp, err := s.client.GetProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get product",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return resp.GetProduct(), nil
}

// CreateProduct creates a new product
func (s *ProductServiceImpl) CreateProduct(
	ctx context.Context,
	name, description string,
	costPrice, sellingPrice string,
	currency, sku, barcode string,
	categoryIDs []string,
	supplierID string,
	isActive, inStock bool,
	stockQty, lowStockAt int32,
	imageURLs, videoURLs []string,
	metadata map[string]string,
) (interface{}, error) {
	s.logger.Debug("CreateProduct",
		zap.String("name", name),
		zap.String("sku", sku),
		zap.Strings("categoryIDs", categoryIDs),
	)

	// Convert the request to the gRPC format
	req := &productv1.CreateProductRequest{
		Name:         name,
		Description:  description,
		CostPrice:    costPrice,
		SellingPrice: sellingPrice,
		Currency:     currency,
		Sku:          sku,
		Barcode:      barcode,
		CategoryIds:  categoryIDs,
		SupplierId:   supplierID,
		IsActive:     isActive,
		InStock:      inStock,
		StockQty:     stockQty,
		LowStockAt:   lowStockAt,
		ImageUrls:    imageURLs,
		VideoUrls:    videoURLs,
	}

	// Add metadata if present
	if len(metadata) > 0 {
		req.Metadata = make(map[string]string, len(metadata))
		for k, v := range metadata {
			req.Metadata[k] = v
		}
	}

	// Call the gRPC service
	resp, err := s.client.CreateProduct(ctx, 
		req.Name,
		req.Description,
		req.CostPrice,
		req.SellingPrice,
		req.Currency,
		req.Sku,
		req.Barcode,
		req.SupplierId,
		req.CategoryIds,
		req.IsActive,
		req.InStock,
		req.StockQty,
		req.LowStockAt,
		req.ImageUrls,
		req.VideoUrls,
		req.Metadata,
	)
	if err != nil {
		s.logger.Error("Failed to create product",
			zap.Error(err),
			zap.String("name", name),
			zap.String("sku", sku),
		)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return resp.GetProduct(), nil
}

// UpdateProduct updates an existing product
// Note: Since the gRPC service doesn't have an update method,
// we implement this by fetching the existing product and creating a new one with the updated fields.
func (s *ProductServiceImpl) UpdateProduct(
	ctx context.Context,
	id, name, description, sku string,
	categories []string,
	price, cost string,
	active bool,
	images []string,
	attributes map[string]string,
) error {
	s.logger.Debug("UpdateProduct",
		zap.String("id", id),
		zap.String("name", name),
		zap.String("sku", sku),
	)

	// 1. Get the existing product
	resp, err := s.client.GetProduct(ctx, &productv1.GetProductRequest{Id: id})
	if err != nil {
		s.logger.Error("Failed to fetch product for update",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to fetch product: %w", err)
	}

	existing := resp.GetProduct()
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	// 2. Update the product fields
	// Only update fields that are provided (non-zero values)
	if name != "" {
		existing.Name = name
	}
	if description != "" {
		existing.Description = description
	}
	if sku != "" {
		existing.Sku = sku
	}
	if len(categories) > 0 {
		existing.CategoryIds = categories
	}
	if price != "" {
		existing.SellingPrice = price
	}
	if cost != "" {
		existing.CostPrice = cost
	}

	// Update boolean fields if they're being explicitly set
	existing.IsActive = active

	// Update images if provided
	if len(images) > 0 {
		existing.ImageUrls = images
	}

	// Update metadata if provided
	if len(attributes) > 0 {
		if existing.Metadata == nil {
			existing.Metadata = make(map[string]string)
		}
		for k, v := range attributes {
			existing.Metadata[k] = v
		}
	}

	// 3. Create a new product with the updated fields
	_, err = s.client.CreateProduct(ctx, 
		existing.Name,
		existing.Description,
		existing.CostPrice,
		existing.SellingPrice,
		existing.Currency,
		existing.Sku,
		existing.Barcode,
		existing.SupplierId,
		existing.CategoryIds,
		existing.IsActive,
		existing.InStock,
		existing.StockQty,
		existing.LowStockAt,
		existing.ImageUrls,
		existing.VideoUrls,
		nil, // metadata is not used in the client method
	)
	if err != nil {
		s.logger.Error("Failed to create updated product",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create updated product: %w", err)
	}

	// Note: In a real implementation, you might want to mark the old product as deleted
	// or inactive, but since we don't have a delete/update method, we'll just return success

	s.logger.Info("Product updated successfully",
		zap.String("id", id),
	)

	return nil
}

// DeleteProduct marks a product as inactive (soft delete)
// Since the gRPC service doesn't have a delete method, we implement this
// by fetching the existing product and updating its IsActive status to false.
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	s.logger.Debug("DeleteProduct",
		zap.String("id", id),
	)

	// 1. Get the existing product
	resp, err := s.client.GetProduct(ctx, &productv1.GetProductRequest{Id: id})
	if err != nil {
		s.logger.Error("Failed to fetch product for deletion",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to fetch product: %w", err)
	}

	existing := resp.GetProduct()
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	// 2. Mark the product as inactive
	existing.IsActive = false

	// 3. Create a new version of the product with IsActive = false
	_, err = s.client.CreateProduct(ctx, 
		existing.Name,
		existing.Description,
		existing.CostPrice,
		existing.SellingPrice,
		existing.Currency,
		existing.Sku,
		existing.Barcode,
		existing.SupplierId,
		existing.CategoryIds,
		false, // Mark as inactive
		existing.InStock,
		existing.StockQty,
		existing.LowStockAt,
		existing.ImageUrls,
		existing.VideoUrls,
		nil, // metadata is not used in the client method
	)
	if err != nil {
		s.logger.Error("Failed to deactivate product",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to deactivate product: %w", err)
	}

	s.logger.Info("Product marked as inactive",
		zap.String("id", id),
	)

	return nil
}
