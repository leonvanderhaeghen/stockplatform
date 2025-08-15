package product

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
)

// Client provides a high-level interface for interacting with the Product service
type Client struct {
	conn   *grpc.ClientConn
	client productv1.ProductServiceClient
	logger *zap.Logger
}

// Config holds configuration for the Product client
type Config struct {
	Address string
	Timeout time.Duration
}

// New creates a new Product service client
func New(config Config, logger *zap.Logger) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	conn, err := grpc.Dial(config.Address, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := productv1.NewProductServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the connection to the Product service
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateProduct creates a new product
func (c *Client) CreateProduct(ctx context.Context, name, description, sku, supplierID string, costPrice, sellingPrice float64, isActive bool, categoryIDs []string) (*models.CreateProductResponse, error) {
	c.logger.Debug("Creating product", zap.String("name", name))
	
	req := convertToCreateProductRequest(name, description, sku, supplierID, costPrice, sellingPrice, isActive, categoryIDs)
	
	resp, err := c.client.CreateProduct(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create product", zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	c.logger.Debug("Product created successfully", zap.String("id", resp.Product.Id))
	return convertToCreateProductResponse(resp), nil
}

// GetProduct retrieves a product by ID
func (c *Client) GetProduct(ctx context.Context, id string) (*models.Product, error) {
	c.logger.Debug("Getting product", zap.String("id", id))
	
	req := &productv1.GetProductRequest{
		Id: id,
	}
	
	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get product", zap.Error(err))
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	
	return convertToProduct(resp.Product), nil
}

// ListProducts lists products with filtering and sorting
func (c *Client) ListProducts(ctx context.Context, categoryID, supplierID string, isActive *bool, limit, offset int32) (*models.ListProductsResponse, error) {
	c.logger.Debug("Listing products")
	
	req := &productv1.ListProductsRequest{
		Pagination: &productv1.Pagination{
			Page:     (offset / limit) + 1,
			PageSize: limit,
		},
	}
	
	// Add filters if provided
	if categoryID != "" || supplierID != "" {
		req.Filter = &productv1.ProductFilter{}
		if categoryID != "" {
			req.Filter.CategoryIds = []string{categoryID}
		}
		// Note: SupplierId not available in ProductFilter, ignoring for now
	}
	
	resp, err := c.client.ListProducts(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list products", zap.Error(err))
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	
	return convertToListProductsResponse(resp), nil
}

// ListCategories lists all product categories
func (c *Client) ListCategories(ctx context.Context, parentID string, limit, offset int32) ([]*models.Category, error) {
	c.logger.Debug("Listing categories")
	
	req := &productv1.ListCategoriesRequest{
		ParentId: parentID,
		Depth:    0, // 0 for all depths
	}
	
	resp, err := c.client.ListCategories(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list categories", zap.Error(err))
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	
	return convertToCategories(resp.Categories), nil
}

// CreateCategory creates a new product category
func (c *Client) CreateCategory(ctx context.Context, name, description, parentID string) (*models.Category, error) {
	c.logger.Debug("Creating category", zap.String("name", name))
	
	req := &productv1.CreateCategoryRequest{
		Name:        name,
		Description: description,
		ParentId:    parentID,
	}
	
	resp, err := c.client.CreateCategory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create category", zap.Error(err))
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	c.logger.Debug("Category created successfully", zap.String("id", resp.Category.Id))
	category := convertProtoCategory(resp.Category)
	return &category, nil
}
