package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	productpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/product/v1"
	"go.uber.org/zap"
)

type ProductClient struct {
	conn   *grpc.ClientConn
	client productpb.ProductServiceClient
	logger *zap.Logger
}

// NewProductClient creates a new gRPC client for the Product service
func NewProductClient(addr string) (*ProductClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := productpb.NewProductServiceClient(conn)
	return &ProductClient{
		conn:   conn,
		client: client,
		logger: zap.NewNop(), // Initialize with no-op logger by default
	}, nil
}

// Close closes the gRPC connection
func (c *ProductClient) Close() error {
	return c.conn.Close()
}

// CreateProduct creates a new product
func (c *ProductClient) CreateProduct(ctx context.Context, name, description, costPrice, sellingPrice, currency, sku, barcode, supplierID string, categoryIDs []string, isActive, inStock bool, stockQty, lowStockAt int32, imageURLs, videoURLs []string, metadata map[string]string) (*productpb.CreateProductResponse, error) {
	req := &productpb.CreateProductRequest{
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
		Metadata:     metadata,
	}
	return c.client.CreateProduct(ctx, req)
}

// GetProduct retrieves a product by ID
func (c *ProductClient) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	return c.client.GetProduct(ctx, req)
}

// ListProducts lists all products with pagination
func (c *ProductClient) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
	return c.client.ListProducts(ctx, req)
}

// ListCategories lists all product categories
func (c *ProductClient) ListCategories(ctx context.Context, req *productpb.ListCategoriesRequest) (*productpb.ListCategoriesResponse, error) {
	return c.client.ListCategories(ctx, req)
}

// CreateCategory creates a new product category
func (c *ProductClient) CreateCategory(ctx context.Context, name, description, parentID string, isActive bool) (*productpb.CreateCategoryResponse, error) {
	req := &productpb.CreateCategoryRequest{
		Name:        name,
		Description: description,
		ParentId:    parentID,
		IsActive:    isActive,
	}
	return c.client.CreateCategory(ctx, req)
}
