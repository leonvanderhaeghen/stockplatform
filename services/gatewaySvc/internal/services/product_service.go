package services

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	productclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/product"
)

// ProductServiceImpl implements the ProductService interface
type ProductServiceImpl struct {
	client *productclient.Client
	logger *zap.Logger
}

// NewProductService creates a new instance of ProductServiceImpl
func NewProductService(productServiceAddr string, logger *zap.Logger) (ProductService, error) {
	// Create a new gRPC client
	prodCfg := productclient.Config{Address: productServiceAddr}
	client, err := productclient.New(prodCfg, logger)
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


	resp, err := s.client.CreateCategory(ctx, name, description, parentID)
	if err != nil {
		s.logger.Error("Failed to create category",
			zap.String("name", name),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return resp, nil
}

// ListCategories lists all product categories
// parentID: Optional parent category ID to filter by
// depth: Maximum depth of subcategories to return (0 for all)
func (s *ProductServiceImpl) ListCategories(
	ctx context.Context,
) (interface{}, error) {
	s.logger.Debug("ListCategories")

	// Call the gRPC service with default parameters
	resp, err := s.client.ListCategories(ctx, "", 100, 0) // parentID="", limit=100, offset=0
	if err != nil {
		s.logger.Error("Failed to list categories",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	// Return the client response directly
	return resp, nil
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

	// Use simplified client interface - pass primitive parameters directly
	var isActivePtr *bool
	if active {
		isActivePtr = &active
	}

	// Call the gRPC service via client abstraction
	resp, err := s.client.ListProducts(ctx, categoryID, query, isActivePtr, int32(limit), int32(offset))
	if err != nil {
		s.logger.Error("Failed to list products",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Return the client response directly
	return resp, nil
}

// GetProductByID gets a product by ID
func (s *ProductServiceImpl) GetProductByID(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetProductByID",
		zap.String("id", id),
	)

	product, err := s.client.GetProduct(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get product",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
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

	// Parse price strings to floats
	costPriceFloat, err := strconv.ParseFloat(costPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid cost price: %w", err)
	}
	sellingPriceFloat, err := strconv.ParseFloat(sellingPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid selling price: %w", err)
	}

	// Call the gRPC service using refactored client
	resp, err := s.client.CreateProduct(ctx, name, description, sku, supplierID, costPriceFloat, sellingPriceFloat, isActive, categoryIDs)
	if err != nil {
		s.logger.Error("Failed to create product",
			zap.Error(err),
			zap.String("name", name),
			zap.String("sku", sku),
		)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return resp, nil
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

	// Note: UpdateProduct is not implemented in the current gRPC client abstraction
	// This method currently returns an error indicating the limitation
	s.logger.Warn("UpdateProduct not implemented in client abstraction",
		zap.String("id", id),
	)
	return fmt.Errorf("update product operation not supported by client abstraction")
}

// DeleteProduct marks a product as inactive (soft delete)
// Since the gRPC service doesn't have a delete method, we implement this
// by fetching the existing product and updating its IsActive status to false.
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	s.logger.Debug("DeleteProduct",
		zap.String("id", id),
	)

	// Note: DeleteProduct is not implemented in the current gRPC client abstraction
	// This method currently returns an error indicating the limitation
	s.logger.Warn("DeleteProduct not implemented in client abstraction",
		zap.String("id", id),
	)
	return fmt.Errorf("delete product operation not supported by client abstraction")
}
