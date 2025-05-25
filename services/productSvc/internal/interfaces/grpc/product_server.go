package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"../../../../gen/go/product/v1"
	"../../application"
	"../../domain"
)

// ProductServer handles gRPC requests for the Product service
type ProductServer struct {
	productv1.UnimplementedProductServiceServer
	service *application.ProductService
	logger   *zap.Logger
}

// NewProductServer creates a new ProductServer
func NewProductServer(service *application.ProductService, logger *zap.Logger) *ProductServer {
	return &ProductServer{
		service: service,
		logger:  logger.Named("grpc_product_server"),
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

	// Convert protobuf message to domain model
	product := &domain.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		SKU:         req.GetSku(),
		CategoryID:  req.GetCategoryId(),
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
	return &productv1.CreateProductResponse{
		Product: &productv1.Product{
			Id:          created.ID.Hex(),
			Name:        created.Name,
			Description: created.Description,
			Price:       created.Price,
			Sku:         created.SKU,
			CategoryId:  created.CategoryID,
			CreatedAt:   timestamppb.New(created.CreatedAt),
			UpdatedAt:   timestamppb.New(created.UpdatedAt),
		},
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
	return &productv1.GetProductResponse{
		Product: &productv1.Product{
			Id:          product.ID.Hex(),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Sku:         product.SKU,
			CategoryId:  product.CategoryID,
			ImageUrls:   product.ImageURLs,
			CreatedAt:   timestamppb.New(product.CreatedAt),
			UpdatedAt:   timestamppb.New(product.UpdatedAt),
		},
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

	// Call the application service
	products, total, err := s.service.ListProducts(ctx, opts)
	if err != nil {
		s.logError(log, err, "Failed to list products")
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	// Convert domain models to protobuf
	pbProducts := make([]*productv1.Product, 0, len(products))
	for _, p := range products {
		pbProducts = append(pbProducts, &productv1.Product{
			Id:          p.ID.Hex(),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Sku:         p.SKU,
			CategoryId:  p.CategoryID,
			ImageUrls:   p.ImageURLs,
			CreatedAt:   timestamppb.New(p.CreatedAt),
			UpdatedAt:   timestamppb.New(p.UpdatedAt),
		})
	}

	// Log successful operation
	log.Info("Products listed successfully",
		zap.Int("count", len(products)),
		zap.Int64("total", total),
		zap.Duration("duration", time.Since(start)),
	)

	return &productv1.ListProductsResponse{
		Products:   pbProducts,
		TotalCount: total,
		Page:       int32(opts.Pagination.Page),
		PageSize:   int32(opts.Pagination.PageSize),
	}, nil
}

// logError logs errors with additional context
func (s *ProductServer) logError(log *zap.Logger, err error, msg string) {
	log.Error(msg,
		zap.String("error", err.Error()),
		zap.String("error_type", fmt.Sprintf("%T", err)),
		zap.Stack("stack"),
	)
}
