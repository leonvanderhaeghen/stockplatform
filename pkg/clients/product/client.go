package product

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

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
func (c *Client) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductResponse, error) {
	c.logger.Debug("Creating product", zap.String("name", req.Name))
	
	resp, err := c.client.CreateProduct(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create product", zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	c.logger.Debug("Product created successfully", zap.String("id", resp.Product.Id))
	return resp, nil
}

// GetProduct retrieves a product by ID
func (c *Client) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
	c.logger.Debug("Getting product", zap.String("id", req.Id))
	
	resp, err := c.client.GetProduct(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get product", zap.Error(err))
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	
	return resp, nil
}

// ListProducts lists products with filtering and sorting
func (c *Client) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	c.logger.Debug("Listing products")
	
	resp, err := c.client.ListProducts(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list products", zap.Error(err))
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	
	return resp, nil
}

// ListCategories lists all product categories
func (c *Client) ListCategories(ctx context.Context, req *productv1.ListCategoriesRequest) (*productv1.ListCategoriesResponse, error) {
	c.logger.Debug("Listing categories")
	
	resp, err := c.client.ListCategories(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list categories", zap.Error(err))
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	
	return resp, nil
}

// CreateCategory creates a new product category
func (c *Client) CreateCategory(ctx context.Context, req *productv1.CreateCategoryRequest) (*productv1.CreateCategoryResponse, error) {
	c.logger.Debug("Creating category", zap.String("name", req.Name))
	
	resp, err := c.client.CreateCategory(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create category", zap.Error(err))
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	c.logger.Debug("Category created successfully", zap.String("id", resp.Category.Id))
	return resp, nil
}
