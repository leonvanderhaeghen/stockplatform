package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	productv1 "stockplatform/pkg/gen/product/v1"
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
		CategoryId: categoryID,
		Query:      query,
		Active:     active,
		Limit:      int32(limit),
		Offset:     int32(offset),
		SortBy:     sortBy,
		Ascending:  ascending,
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
		Product: &productv1.Product{
			Name:        name,
			Description: description,
			Sku:         sku,
			Categories:  categories,
			Price:       price,
			Cost:        cost,
			Active:      active,
			Images:      images,
			Attributes:  attributes,
		},
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
func (s *ProductServiceImpl) UpdateProduct(
	ctx context.Context,
	id, name, description, sku string,
	categories []string,
	price, cost float64,
	active bool,
	images []string,
	attributes map[string]string,
) error {
	s.logger.Debug("UpdateProduct",
		zap.String("id", id),
		zap.String("name", name),
		zap.String("sku", sku),
	)

	req := &productv1.UpdateProductRequest{
		Product: &productv1.Product{
			Id:          id,
			Name:        name,
			Description: description,
			Sku:         sku,
			Categories:  categories,
			Price:       price,
			Cost:        cost,
			Active:      active,
			Images:      images,
			Attributes:  attributes,
		},
	}

	_, err := s.client.UpdateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update product",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// DeleteProduct deletes a product
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	s.logger.Debug("DeleteProduct",
		zap.String("id", id),
	)

	req := &productv1.DeleteProductRequest{
		Id: id,
	}

	_, err := s.client.DeleteProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete product",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
