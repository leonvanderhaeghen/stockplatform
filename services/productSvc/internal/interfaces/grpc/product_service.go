package grpc

import (
	"context"
	"errors"

	"../../logger"
	"../../../../gen/go/product/v1"
	"../../application"
	"../../domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProductServiceServer implements the ProductService gRPC service
type ProductServiceServer struct {
	productv1.UnimplementedProductServiceServer
	service *application.ProductService
	logger   *logger.Logger
	logger   *zap.Logger
}

// NewProductServiceServer creates a new ProductServiceServer
func NewProductServiceServer(
	service *application.ProductService,
	logger *zap.Logger,
) *ProductServiceServer {
	return &ProductServiceServer{
		service: service,
		logger:  logger.Named("grpc_product_service"),
	}
}

// CreateProduct creates a new product
func (s *ProductServiceServer) CreateProduct(
	ctx context.Context,
	req *productv1.CreateProductRequest,
) (*productv1.CreateProductResponse, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.Sku,
		CategoryID:  req.CategoryId,
		ImageURLs:   req.ImageUrls,
	}

	created, err := s.service.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, status.Error(convertToGRPCError(err))
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
	product, err := s.service.GetProduct(ctx, req.Id)
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
	case errors.Is(err, domain.ErrNotFound):
		return codes.NotFound, "product not found"
	case errors.Is(err, domain.ErrInvalidID):
		return codes.InvalidArgument, "invalid product ID"
	case errors.Is(err, domain.ErrValidation):
		return codes.InvalidArgument, "invalid product data"
	case errors.Is(err, domain.ErrAlreadyExists):
		return codes.AlreadyExists, "product already exists"
	default:
		return codes.Internal, "internal server error"
	}
}

// toProtoProduct converts a domain Product to a protobuf Product
func toProtoProduct(p *domain.Product) *productv1.Product {
	return &productv1.Product{
		Id:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Sku:         p.SKU,
		CategoryId:  p.CategoryID,
		ImageUrls:   p.ImageURLs,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}
