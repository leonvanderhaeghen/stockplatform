package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
	supplierv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/supplier/v1"
)

// SupplierServer handles gRPC requests for the Supplier service
type SupplierServer struct {
	supplierv1.UnimplementedSupplierServiceServer
	service *application.SupplierService
	logger   *zap.Logger
}

// NewSupplierServer creates a new SupplierServer
func NewSupplierServer(
	service *application.SupplierService,
	logger *zap.Logger,
) *SupplierServer {
	return &SupplierServer{
		service: service,
		logger:   logger.Named("grpc_supplier_server"),
	}
}

// CreateSupplier handles the CreateSupplier gRPC request
func (s *SupplierServer) CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.CreateSupplierResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "CreateSupplier"),
		zap.String("name", req.GetName()),
	)

	// Log incoming request
	log.Debug("Processing CreateSupplier request")

	// Convert request to domain model
	supplier := &domain.Supplier{
		Name:          req.GetName(),
		ContactPerson: req.GetContactPerson(),
		Email:         req.GetEmail(),
		Phone:         req.GetPhone(),
		Address:       req.GetAddress(),
		City:          req.GetCity(),
		State:         req.GetState(),
		Country:       req.GetCountry(),
		PostalCode:    req.GetPostalCode(),
		TaxID:         req.GetTaxId(),
		Website:       req.GetWebsite(),
		Currency:      req.GetCurrency(),
		LeadTimeDays:  req.GetLeadTimeDays(),
		PaymentTerms:  req.GetPaymentTerms(),
		Metadata:      req.GetMetadata(),
	}

	// Create supplier using the service
	createdSupplier, err := s.service.CreateSupplier(ctx, supplier)
	if err != nil {
		s.logError(log, err, "Failed to create supplier")
		switch err {
		case domain.ErrSupplierNameRequired:
			return nil, status.Error(codes.InvalidArgument, domain.ErrSupplierNameRequired.Error())
		case domain.ErrAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "supplier with this name already exists")
		default:
			return nil, status.Error(codes.Internal, "failed to create supplier")
		}
	}

	// Convert domain model to protobuf response
	supplierProto := s.toProto(createdSupplier)

	log.Info("Successfully created supplier",
		zap.String("supplierID", createdSupplier.ID.Hex()),
		zap.Duration("duration", time.Since(start)),
	)

	return &supplierv1.CreateSupplierResponse{
		Supplier: supplierProto,
	}, nil
}

// GetSupplier handles the GetSupplier gRPC request
func (s *SupplierServer) GetSupplier(ctx context.Context, req *supplierv1.GetSupplierRequest) (*supplierv1.GetSupplierResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "GetSupplier"),
		zap.String("id", req.GetId()),
	)

	// Log incoming request
	log.Debug("Processing GetSupplier request")

	// Get supplier using the service
	supplier, err := s.service.GetSupplier(ctx, req.GetId())
	if err != nil {
		s.logError(log, err, "Failed to get supplier")
		switch err {
		case domain.ErrSupplierNotFound:
			return nil, status.Error(codes.NotFound, domain.ErrSupplierNotFound.Error())
		case domain.ErrInvalidID:
			return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidID.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to get supplier")
		}
	}

	// Convert domain model to protobuf response
	supplierProto := s.toProto(supplier)

	log.Info("Successfully retrieved supplier",
		zap.String("supplierID", supplier.ID.Hex()),
		zap.Duration("duration", time.Since(start)),
	)

	return &supplierv1.GetSupplierResponse{
		Supplier: supplierProto,
	}, nil
}

// UpdateSupplier handles the UpdateSupplier gRPC request
func (s *SupplierServer) UpdateSupplier(ctx context.Context, req *supplierv1.UpdateSupplierRequest) (*supplierv1.UpdateSupplierResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "UpdateSupplier"),
		zap.String("id", req.GetId()),
	)

	// Log incoming request
	log.Debug("Processing UpdateSupplier request")

	// Convert request to domain model
	supplier := &domain.Supplier{
		Name:          req.GetName(),
		ContactPerson: req.GetContactPerson(),
		Email:         req.GetEmail(),
		Phone:         req.GetPhone(),
		Address:       req.GetAddress(),
		City:          req.GetCity(),
		State:         req.GetState(),
		Country:       req.GetCountry(),
		PostalCode:    req.GetPostalCode(),
		TaxID:         req.GetTaxId(),
		Website:       req.GetWebsite(),
		Currency:      req.GetCurrency(),
		LeadTimeDays:  req.GetLeadTimeDays(),
		PaymentTerms:  req.GetPaymentTerms(),
		Metadata:      req.GetMetadata(),
	}

	// Set ID from request
	objectID, err := domain.StringToObjectID(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidID.Error())
	}
	supplier.ID = objectID

	// Update supplier using the service
	err = s.service.UpdateSupplier(ctx, supplier)
	if err != nil {
		s.logError(log, err, "Failed to update supplier")
		switch err {
		case domain.ErrSupplierNotFound:
			return nil, status.Error(codes.NotFound, domain.ErrSupplierNotFound.Error())
		case domain.ErrInvalidID:
			return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidID.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to update supplier")
		}
	}

	// Get updated supplier to return
	updatedSupplier, err := s.service.GetSupplier(ctx, req.GetId())
	if err != nil {
		s.logError(log, err, "Failed to get updated supplier")
		return nil, status.Error(codes.Internal, "failed to get updated supplier")
	}

	// Convert domain model to protobuf response
	supplierProto := s.toProto(updatedSupplier)

	log.Info("Successfully updated supplier",
		zap.String("supplierID", updatedSupplier.ID.Hex()),
		zap.Duration("duration", time.Since(start)),
	)

	return &supplierv1.UpdateSupplierResponse{
		Supplier: supplierProto,
	}, nil
}

// DeleteSupplier handles the DeleteSupplier gRPC request
func (s *SupplierServer) DeleteSupplier(ctx context.Context, req *supplierv1.DeleteSupplierRequest) (*supplierv1.DeleteSupplierResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "DeleteSupplier"),
		zap.String("id", req.GetId()),
	)

	// Log incoming request
	log.Debug("Processing DeleteSupplier request")

	// Delete supplier using the service
	err := s.service.DeleteSupplier(ctx, req.GetId())
	if err != nil {
		s.logError(log, err, "Failed to delete supplier")
		switch err {
		case domain.ErrSupplierNotFound:
			return nil, status.Error(codes.NotFound, domain.ErrSupplierNotFound.Error())
		case domain.ErrInvalidID:
			return nil, status.Error(codes.InvalidArgument, domain.ErrInvalidID.Error())
		default:
			return nil, status.Error(codes.Internal, "failed to delete supplier")
		}
	}

	log.Info("Successfully deleted supplier",
		zap.String("supplierID", req.GetId()),
		zap.Duration("duration", time.Since(start)),
	)

	return &supplierv1.DeleteSupplierResponse{
		Success: true,
	}, nil
}

// ListSuppliers handles the ListSuppliers gRPC request
func (s *SupplierServer) ListSuppliers(ctx context.Context, req *supplierv1.ListSuppliersRequest) (*supplierv1.ListSuppliersResponse, error) {
	start := time.Now()
	log := s.logger.With(
		zap.String("method", "ListSuppliers"),
		zap.Int32("page", req.GetPage()),
		zap.Int32("page_size", req.GetPageSize()),
	)

	// Log incoming request
	log.Debug("Processing ListSuppliers request")

	// Validate pagination
	page := req.GetPage()
	if page < 1 {
		page = 1
	}

	pageSize := req.GetPageSize()
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	// Validate pagination parameters
	if page < 1 {
		err := domain.ErrInvalidArgument
		s.logError(log, err, "Invalid page number")
		return nil, status.Error(codes.InvalidArgument, "page must be greater than 0")
	}

	if pageSize < 1 || pageSize > 100 {
		err := domain.ErrInvalidArgument
		s.logError(log, err, "Invalid page size")
		return nil, status.Error(codes.InvalidArgument, "page_size must be between 1 and 100")
	}

	// List suppliers using the service
	suppliers, total, err := s.service.ListSuppliers(ctx, page, pageSize, req.GetSearch())
	if err != nil {
		s.logError(log, err, "Failed to list suppliers")
		return nil, status.Error(codes.Internal, "failed to list suppliers")
	}

	// Convert domain models to protobuf responses
	supplierProtos := make([]*supplierv1.Supplier, 0, len(suppliers))
	for _, supplier := range suppliers {
		supplierProtos = append(supplierProtos, s.toProto(supplier))
	}

	log.Info("Successfully listed suppliers",
		zap.Int("count", len(supplierProtos)),
		zap.Duration("duration", time.Since(start)),
	)

	return &supplierv1.ListSuppliersResponse{
		Suppliers: supplierProtos,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// toProto converts a domain Supplier to a protobuf Supplier
func (s *SupplierServer) toProto(supplier *domain.Supplier) *supplierv1.Supplier {
	if supplier == nil {
		return nil
	}

	return &supplierv1.Supplier{
		Id:            supplier.ID.Hex(),
		Name:          supplier.Name,
		ContactPerson: supplier.ContactPerson,
		Email:         supplier.Email,
		Phone:         supplier.Phone,
		Address:       supplier.Address,
		City:          supplier.City,
		State:         supplier.State,
		Country:       supplier.Country,
		PostalCode:    supplier.PostalCode,
		TaxId:         supplier.TaxID,
		Website:       supplier.Website,
		Currency:      supplier.Currency,
		LeadTimeDays:  supplier.LeadTimeDays,
		PaymentTerms:  supplier.PaymentTerms,
		Metadata:      supplier.Metadata,
		CreatedAt:     timestamppb.New(supplier.CreatedAt),
		UpdatedAt:     timestamppb.New(supplier.UpdatedAt),
	}
}

// logError logs errors with additional context
func (s *SupplierServer) logError(log *zap.Logger, err error, msg string) {
	log.Error(msg,
		zap.Error(err),
		zap.String("error_type", err.Error()),
	)
}
