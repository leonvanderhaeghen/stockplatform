package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
)

// convertMetadata converts a map[string]interface{} to map[string]string by converting all values to strings
func convertMetadata(metadata map[string]interface{}) map[string]string {
	if metadata == nil {
		return nil
	}

	result := make(map[string]string, len(metadata))
	for k, v := range metadata {
		switch v := v.(type) {
		case string:
			result[k] = v
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result[k] = fmt.Sprintf("%d", v)
		case float32, float64:
			result[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case bool:
			result[k] = strconv.FormatBool(v)
		default:
			// For complex types, use JSON marshaling
			if b, err := json.Marshal(v); err == nil {
				result[k] = string(b)
			} else {
				// Fallback to fmt.Sprint if JSON marshaling fails
				result[k] = fmt.Sprint(v)
			}
		}
	}
	return result
}

// ProductServer handles gRPC requests for the Product service
type ProductServer struct {
	productv1.UnimplementedProductServiceServer
	service        *application.ProductService
	categoryService *application.CategoryService
	logger         *zap.Logger
}

// NewProductServer creates a new ProductServer
func NewProductServer(
	service *application.ProductService,
	categoryService *application.CategoryService,
	logger *zap.Logger,
) *ProductServer {
	return &ProductServer{
		service:        service,
		categoryService: categoryService,
		logger:         logger.Named("grpc_product_server"),
	}
}

// CreateProduct handles the CreateProduct gRPC request
func (s *ProductServer) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "CreateProduct"),
		zap.String("name", req.GetName()),
		zap.String("sku", req.GetSku()),
	)

	// Log incoming request
	log.Debug("Processing CreateProduct request",
		zap.Any("request", req),
	)

	// Validate request
	if req.GetName() == "" {
		err := status.Error(codes.InvalidArgument, "product name is required")
		s.logError(log, err, "Validation failed")
		return nil, err
	}

	if req.GetSku() == "" {
		err := status.Error(codes.InvalidArgument, "product SKU is required")
		s.logError(log, err, "Validation failed")
		return nil, err
	}

	// Convert request metadata from map[string]string to map[string]interface{}
	metadata := make(map[string]interface{})
	for k, v := range req.GetMetadata() {
		metadata[k] = v
	}

	// Convert protobuf message to domain model
	product := &domain.Product{
		Name:         req.GetName(),
		Description:   req.GetDescription(),
		CostPrice:     req.GetCostPrice(),
		SellingPrice:  req.GetSellingPrice(),
		Currency:      req.GetCurrency(),
		SKU:           req.GetSku(),
		Barcode:       req.GetBarcode(),
		CategoryIDs:   req.GetCategoryIds(),
		SupplierID:    req.GetSupplierId(),
		IsActive:      req.GetIsActive(),



		ImageURLs:     req.GetImageUrls(),
		VideoURLs:     req.GetVideoUrls(),
		Metadata:      metadata,
	}

	// Call the application service
	created, err := s.service.CreateProduct(ctx, product)
	if err != nil {
		s.logError(log, err, "Failed to create product")
		
		// Convert domain errors to gRPC status errors
		switch {
		case errors.Is(err, domain.ErrValidation):
			return nil, status.Error(codes.InvalidArgument, "invalid product data: " + err.Error())
		case errors.Is(err, domain.ErrAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, "product with this SKU already exists")
		default:
			// Log the full error for debugging
			log.Error("Internal server error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	// Log successful operation
	log.Info("Product created successfully",
		zap.String("product_id", created.ID.Hex()),
		zap.Duration("duration", time.Since(start)),
	)

	// Convert domain model back to protobuf
	pbProduct := &productv1.Product{
		Id:           created.ID.Hex(),
		Name:         created.Name,
		Description:   created.Description,
		CostPrice:     created.CostPrice,
		SellingPrice:  created.SellingPrice,
		Currency:      created.Currency,
		Sku:           created.SKU,
		Barcode:       created.Barcode,
		CategoryIds:   created.CategoryIDs,
		SupplierId:    created.SupplierID,
		IsActive:      created.IsActive,
		ImageUrls:     created.ImageURLs,
		VideoUrls:     created.VideoURLs,
		Metadata:      convertMetadata(created.Metadata),
	}

	// Only set timestamps if they are not zero
	if !created.CreatedAt.IsZero() {
		pbProduct.CreatedAt = timestamppb.New(created.CreatedAt)
	}
	if !created.UpdatedAt.IsZero() {
		pbProduct.UpdatedAt = timestamppb.New(created.UpdatedAt)
	}

	return &productv1.CreateProductResponse{
		Product: pbProduct,
	}, nil
}

// GetProduct handles the GetProduct gRPC request
func (s *ProductServer) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "GetProduct"),
		zap.String("product_id", req.GetId()),
	)

	// Log incoming request
	log.Debug("Processing GetProduct request")

	// Validate request
	if req.GetId() == "" {
		err := status.Error(codes.InvalidArgument, "product ID is required")
		s.logError(log, err, "Validation failed")
		return nil, err
	}

	// Call the application service
	product, err := s.service.GetProduct(ctx, req.GetId())
	if err != nil {
		s.logError(log, err, "Failed to get product")
		
		// Convert domain errors to gRPC status errors
		switch {
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.Error(codes.NotFound, "product not found")
		case errors.Is(err, domain.ErrInvalidID):
			return nil, status.Error(codes.InvalidArgument, "invalid product ID format")
		default:
			// Log the full error for debugging
			log.Error("Internal server error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	// Log successful operation
	log.Info("Product retrieved successfully",
		zap.String("product_id", product.ID.Hex()),
		zap.Duration("duration", time.Since(start)),
	)

	// Convert domain model to protobuf
	pbProduct := &productv1.Product{
		Id:           product.ID.Hex(),
		Name:         product.Name,
		Description:   product.Description,
		CostPrice:     product.CostPrice,
		SellingPrice:  product.SellingPrice,
		Currency:      product.Currency,
		Sku:           product.SKU,
		Barcode:       product.Barcode,
		CategoryIds:   product.CategoryIDs,
		SupplierId:    product.SupplierID,
		IsActive:      product.IsActive,
		ImageUrls:     product.ImageURLs,
		VideoUrls:     product.VideoURLs,
		Metadata:      convertMetadata(product.Metadata),
	}

	// Only set timestamps if they are not zero
	if !product.CreatedAt.IsZero() {
		pbProduct.CreatedAt = timestamppb.New(product.CreatedAt)
	}
	if !product.UpdatedAt.IsZero() {
		pbProduct.UpdatedAt = timestamppb.New(product.UpdatedAt)
	}

	return &productv1.GetProductResponse{
		Product: pbProduct,
	}, nil
}

// ListProducts handles the ListProducts gRPC request
func (s *ProductServer) ListProducts(ctx context.Context, req *productv1.ListProductsRequest) (*productv1.ListProductsResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "ListProducts"),
	)

	// Log incoming request
	log.Debug("Processing ListProducts request", zap.Any("request", req))

	// Convert protobuf request to domain options
	opts := &domain.ListOptions{}

	// Apply filters if provided
	if req.GetFilter() != nil {
		opts.Filter = &domain.ProductFilter{
			IDs:         req.GetFilter().GetIds(),
			CategoryIDs: req.GetFilter().GetCategoryIds(),
			MinPrice:    req.GetFilter().GetMinPrice(),
			MaxPrice:    req.GetFilter().GetMaxPrice(),
			SearchTerm:  req.GetFilter().GetSearchTerm(),
		}
	}

	// Apply sorting if provided
	if req.GetSort() != nil {
		var sortField domain.SortField
		switch req.GetSort().GetField() {
		case productv1.ProductSort_SORT_FIELD_NAME:
			sortField = domain.SortFieldName
		case productv1.ProductSort_SORT_FIELD_PRICE:
			sortField = domain.SortFieldPrice
		case productv1.ProductSort_SORT_FIELD_CREATED_AT:
			sortField = domain.SortFieldCreatedAt
		case productv1.ProductSort_SORT_FIELD_UPDATED_AT:
			sortField = domain.SortFieldUpdatedAt
		default:
			sortField = domain.SortFieldCreatedAt
		}

		var sortOrder domain.SortOrder
		switch req.GetSort().GetOrder() {
		case productv1.ProductSort_SORT_ORDER_ASC:
			sortOrder = domain.SortOrderAsc
		case productv1.ProductSort_SORT_ORDER_DESC:
			sortOrder = domain.SortOrderDesc
		default:
			sortOrder = domain.SortOrderDesc
		}

		opts.Sort = &domain.SortOption{
			Field: sortField,
			Order: sortOrder,
		}
	}

	// Apply pagination
	if req.GetPagination() != nil {
		opts.Pagination = &domain.Pagination{
			Page:     int(req.GetPagination().GetPage()),
			PageSize: int(req.GetPagination().GetPageSize()),
		}
	}

	// Call the service to list products
	products, total, err := s.service.ListProducts(ctx, opts)
	if err != nil {
		s.logError(log, err, "Failed to list products")
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	// Convert domain models to protobuf
	pbProducts := make([]*productv1.Product, 0, len(products))
	for _, p := range products {
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

		pbProducts = append(pbProducts, pbProduct)
	}

	// Log successful operation
	log.Info("Products listed successfully",
		zap.Int("count", len(products)),
		zap.Int64("total", total),
		zap.Duration("duration", time.Since(start)),
	)

	return &productv1.ListProductsResponse{
		Products:   pbProducts,
		TotalCount: int32(total),
		Page:       int32(opts.Pagination.Page),
		PageSize:   int32(opts.Pagination.PageSize),
	}, nil
}

// logError logs errors with additional context
// ListCategories handles the ListCategories gRPC request
func (s *ProductServer) ListCategories(ctx context.Context, req *productv1.ListCategoriesRequest) (*productv1.ListCategoriesResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "ListCategories"),
	)

	// Log incoming request
	log.Debug("Processing ListCategories request", zap.Any("request", req))

	// Convert request parameters
	var parentID string
	if req.GetParentId() != "" {
		parentID = req.GetParentId()
	}

	// Default depth to 3 levels if not specified
	depth := int32(3)
	if req.GetDepth() > 0 {
		depth = req.GetDepth()
	}

	// Call the application service
	categories, err := s.categoryService.ListCategories(ctx, parentID, depth)
	if err != nil {
		s.logError(log, err, "Failed to list categories")
		return nil, status.Error(codes.Internal, "failed to list categories")
	}

	// Convert domain models to protobuf
	pbCategories := make([]*productv1.Category, 0, len(categories))
	for _, cat := range categories {
		pbCategories = append(pbCategories, &productv1.Category{
			Id:          cat.ID.Hex(),
			Name:        cat.Name,
			Description: cat.Description,
			ParentId:    cat.ParentID,
			CreatedAt:   timestamppb.New(cat.CreatedAt),
			UpdatedAt:   timestamppb.New(cat.UpdatedAt),
		})
	}

	// Log successful operation
	log.Info("Categories listed successfully",
		zap.Int("count", len(categories)),
		zap.Duration("duration", time.Since(start)),
	)

	return &productv1.ListCategoriesResponse{
		Categories: pbCategories,
	}, nil
}

// CreateCategory handles the CreateCategory gRPC request
func (s *ProductServer) CreateCategory(ctx context.Context, req *productv1.CreateCategoryRequest) (*productv1.CreateCategoryResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "CreateCategory"),
		zap.String("name", req.GetName()),
	)

	// Log incoming request
	log.Debug("Processing CreateCategory request",
		zap.Any("request", req),
	)

	// Validate request
	if req.GetName() == "" {
		err := status.Error(codes.InvalidArgument, "category name is required")
		s.logError(log, err, "Validation failed")
		return nil, err
	}

	// Convert request to domain model
	category := &domain.Category{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ParentID:    req.GetParentId(),
	}

	// Create category using the category service
	createdCategory, err := s.categoryService.CreateCategory(ctx, category)
	if err != nil {
		s.logError(log, err, "Failed to create category")
		return nil, status.Error(codes.Internal, "failed to create category")
	}

	// Convert domain model to protobuf response
	categoryProto := &productv1.Category{
		Id:          createdCategory.ID.Hex(),
		Name:        createdCategory.Name,
		Description: createdCategory.Description,
		ParentId:    createdCategory.ParentID,
		Level:       createdCategory.Level,
		Path:        createdCategory.Path,
		CreatedAt:   timestamppb.New(createdCategory.CreatedAt),
		UpdatedAt:   timestamppb.New(createdCategory.UpdatedAt),
	}

	log.Info("Successfully created category",
		zap.String("categoryID", createdCategory.ID.Hex()),
		zap.Duration("duration", time.Since(start)),
	)

	return &productv1.CreateCategoryResponse{
		Category: categoryProto,
	}, nil
}

// logError logs errors with additional context
func (s *ProductServer) logError(log *zap.Logger, err error, msg string) {
	log.Error(msg,
		zap.Error(err),
		zap.String("error_type", fmt.Sprintf("%T", err)),
		zap.Stack("stack"),
	)
}
