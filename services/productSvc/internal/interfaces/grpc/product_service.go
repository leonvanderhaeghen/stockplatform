package grpc

import (
	"context"
	"errors"

	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductServiceServer implements the ProductService gRPC service
type ProductServiceServer struct {
	productv1.UnimplementedProductServiceServer
	productService  *application.ProductService
	categoryService *application.CategoryService
	logger          *zap.Logger
}

// NewProductServiceServer creates a new ProductServiceServer with required dependencies
func NewProductServiceServer(
	productService *application.ProductService,
	logger *zap.Logger,
) *ProductServiceServer {
	return &ProductServiceServer{
		productService: productService,
		logger:         logger.Named("grpc_product_service"),
	}
}

// SetCategoryService sets the category service for the server
// This allows for optional dependency injection after server creation
func (s *ProductServiceServer) SetCategoryService(categoryService *application.CategoryService) {
	s.categoryService = categoryService
}

// CreateProduct creates a new product
func (s *ProductServiceServer) CreateProduct(
	ctx context.Context,
	req *productv1.CreateProductRequest,
) (*productv1.CreateProductResponse, error) {
	// Convert metadata from map[string]string to map[string]interface{}
	metadata := make(map[string]interface{}, len(req.Metadata))
	for k, v := range req.Metadata {
		metadata[k] = v
	}

	product := &domain.Product{
		Name:         req.Name,
		Description:   req.Description,
		CostPrice:     req.CostPrice,
		SellingPrice:  req.SellingPrice,
		Currency:      req.Currency,
		SKU:           req.Sku,
		Barcode:       req.Barcode,
		CategoryIDs:   req.CategoryIds,
		SupplierID:    req.SupplierId,
		IsActive:      req.IsActive,
		InStock:       req.InStock,
		StockQty:      req.StockQty,
		LowStockAt:    req.LowStockAt,
		ImageURLs:     req.ImageUrls,
		VideoURLs:     req.VideoUrls,
		Metadata:      metadata,
	}

	created, err := s.productService.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		code, msg := convertToGRPCError(err)
		return nil, status.Error(code, msg)
	}

	return &productv1.CreateProductResponse{
		Product: toProtoProduct(created),
	}, nil
}

// GetProduct retrieves a product by ID
func (s *ProductServiceServer) GetProduct(
	ctx context.Context,
	req *productv1.GetProductRequest,
) (*productv1.GetProductResponse, error) {
	product, err := s.productService.GetProduct(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get product", 
			zap.String("id", req.Id), 
			zap.Error(err))
		return nil, status.Error(convertToGRPCError(err))
	}

	return &productv1.GetProductResponse{
		Product: toProtoProduct(product),
	}, nil
}

// convertToGRPCError converts domain errors to gRPC status errors
func convertToGRPCError(err error) (code codes.Code, msg string) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		return codes.InvalidArgument, "validation failed"
	case errors.Is(err, domain.ErrNotFound):
		return codes.NotFound, "resource not found"
	case errors.Is(err, domain.ErrAlreadyExists):
		return codes.AlreadyExists, "resource already exists"
	case errors.Is(err, domain.ErrParentCategoryNotFound):
		return codes.FailedPrecondition, "parent category not found"
	default:
		return codes.Internal, "internal server error"
	}
}

// ListCategories retrieves a list of categories
func (s *ProductServiceServer) ListCategories(
	ctx context.Context,
	req *productv1.ListCategoriesRequest,
) (*productv1.ListCategoriesResponse, error) {
	// Convert request parameters
	var parentID string
	if req.ParentId != "" {
		parentID = req.ParentId
	}

	var depth int32
	if req.Depth > 0 {
		depth = req.Depth
	}

	// Call the service
	categories, err := s.categoryService.ListCategories(ctx, parentID, depth)
	if err != nil {
		s.logger.Error("Failed to list categories", zap.Error(err))
		code, msg := convertToGRPCError(err)
		return nil, status.Error(code, msg)
	}

	// Convert domain categories to protobuf categories
	pbCategories := make([]*productv1.Category, 0, len(categories))
	for _, cat := range categories {
		pbCategories = append(pbCategories, toProtoCategory(cat))
	}

	return &productv1.ListCategoriesResponse{
		Categories: pbCategories,
	}, nil
}

// toProtoCategory converts a domain Category to a protobuf Category
func toProtoCategory(c *domain.Category) *productv1.Category {
	return &productv1.Category{
		Id:          c.ID.Hex(),
		Name:        c.Name,
		Description: c.Description,
		ParentId:    c.ParentID,
		Level:       c.Level,
		Path:        c.Path,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}

// toProtoProduct converts a domain Product to a protobuf Product
func toProtoProduct(p *domain.Product) *productv1.Product {
	pbProduct := &productv1.Product{
		Id:           p.ID.Hex(),
		Name:         p.Name,
		Description:   p.Description,
		CostPrice:     p.CostPrice,
		SellingPrice:  p.SellingPrice,
		Currency:      p.Currency,
		Sku:           p.SKU,
		Barcode:       p.Barcode,
		CategoryIds:   p.CategoryIDs,
		SupplierId:    p.SupplierID,
		IsActive:      p.IsActive,
		InStock:       p.InStock,
		StockQty:      p.StockQty,
		LowStockAt:    p.LowStockAt,
		ImageUrls:     p.ImageURLs,
		VideoUrls:     p.VideoURLs,
		Metadata:      convertMetadata(p.Metadata),
	}

	// Only set timestamps if they are not zero
	if !p.CreatedAt.IsZero() {
		pbProduct.CreatedAt = timestamppb.New(p.CreatedAt)
	}
	if !p.UpdatedAt.IsZero() {
		pbProduct.UpdatedAt = timestamppb.New(p.UpdatedAt)
	}

	return pbProduct
}
