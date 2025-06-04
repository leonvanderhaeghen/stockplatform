package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	productv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/product/v1"
)

// ProductServiceImpl implements the ProductService interface
type ProductServiceImpl struct {
	client productv1.ProductServiceClient
	logger *zap.Logger
}

// NewProductService creates a new instance of ProductServiceImpl
func NewProductService(productServiceAddr string, logger *zap.Logger) (ProductService, error) {
	// Create a gRPC connection to the product service
	conn, err := grpc.Dial(
		productServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	// Create a client
	client := productv1.NewProductServiceClient(conn)

	return &ProductServiceImpl{
		client: client,
		logger: logger.Named("product_service"),
	}, nil
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

	// Convert parameters to the appropriate types for gRPC
	req := &productv1.ListProductsRequest{
		Filter: &productv1.ProductFilter{
			CategoryIds: []string{categoryID},
			SearchTerm:  query,
		},
		Sort: &productv1.ProductSort{
			Field: productv1.ProductSort_SortField(productv1.ProductSort_SORT_FIELD_UNSPECIFIED), // TODO: Map sortBy to SortField
			Order: productv1.ProductSort_SortOrder(productv1.ProductSort_SORT_ORDER_ASC),      // TODO: Use ascending parameter
		},
		Pagination: &productv1.Pagination{
			Page:     int32(offset/limit + 1),
			PageSize: int32(limit),
		},
	}

	// Call the gRPC service
	resp, err := s.client.ListProducts(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list products",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return resp, nil
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

	return resp.Product, nil
}

// CreateProduct creates a new product
func (s *ProductServiceImpl) CreateProduct(
	ctx context.Context,
	name, description, sku string,
	categories []string,
	price, cost float64,
	active bool,
	images []string,
	attributes map[string]string,
) (interface{}, error) {
	s.logger.Debug("CreateProduct",
		zap.String("name", name),
		zap.String("sku", sku),
		zap.Float64("price", price),
		zap.Bool("active", active),
	)

	req := &productv1.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
		Sku:         sku,
		CategoryId:   categories[0], // Using first category as category_id
		ImageUrls:    images,
	}

	resp, err := s.client.CreateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create product",
			zap.String("name", name),
			zap.String("sku", sku),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return resp.Product, nil
}

// UpdateProduct updates an existing product
// Note: Not implemented in the product service
func (s *ProductServiceImpl) UpdateProduct(
	ctx context.Context,
	id, name, description, sku string,
	categories []string,
	price, cost float64,
	active bool,
	images []string,
	attributes map[string]string,
) error {
	return fmt.Errorf("UpdateProduct is not implemented in the product service")
}

// DeleteProduct deletes a product
// Note: Not implemented in the product service
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	return fmt.Errorf("DeleteProduct is not implemented in the product service")
}
