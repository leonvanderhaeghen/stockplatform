package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

type SupplierServer struct {
	supplierv1.UnimplementedSupplierServiceServer
	service application.SupplierService
	logger  *zap.Logger
}

// NewSupplierServer creates a new gRPC supplier server
func NewSupplierServer(service application.SupplierService, logger *zap.Logger) *SupplierServer {
	return &SupplierServer{
		service: service,
		logger:  logger.Named("supplier_grpc_server"),
	}
}

func (s *SupplierServer) CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.CreateSupplierResponse, error) {
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

	created, err := s.service.CreateSupplier(ctx, supplier)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &supplierv1.CreateSupplierResponse{
		Supplier: domainToPb(created),
	}, nil
}

func (s *SupplierServer) GetSupplier(ctx context.Context, req *supplierv1.GetSupplierRequest) (*supplierv1.GetSupplierResponse, error) {
	supplier, err := s.service.GetSupplier(ctx, req.GetId())
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "supplier not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &supplierv1.GetSupplierResponse{
		Supplier: domainToPb(supplier),
	}, nil
}

func (s *SupplierServer) UpdateSupplier(ctx context.Context, req *supplierv1.UpdateSupplierRequest) (*supplierv1.UpdateSupplierResponse, error) {
	objectID, err := domain.ObjectIDFromString(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid supplier ID")
	}

	supplier := &domain.Supplier{
		ID:            objectID,
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

	updated, err := s.service.UpdateSupplier(ctx, supplier)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "supplier not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &supplierv1.UpdateSupplierResponse{
		Supplier: domainToPb(updated),
	}, nil
}

func (s *SupplierServer) DeleteSupplier(ctx context.Context, req *supplierv1.DeleteSupplierRequest) (*supplierv1.DeleteSupplierResponse, error) {
	if err := s.service.DeleteSupplier(ctx, req.GetId()); err != nil {
		if err == domain.ErrNotFound {
			return nil, status.Error(codes.NotFound, "supplier not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &supplierv1.DeleteSupplierResponse{
		Success: true,
	}, nil
}

func (s *SupplierServer) ListSuppliers(ctx context.Context, req *supplierv1.ListSuppliersRequest) (*supplierv1.ListSuppliersResponse, error) {
	suppliers, total, err := s.service.ListSuppliers(
		ctx,
		req.GetPage(),
		req.GetPageSize(),
		req.GetSearch(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbSuppliers := make([]*supplierv1.Supplier, 0, len(suppliers))
	for _, s := range suppliers {
		pbSuppliers = append(pbSuppliers, domainToPb(s))
	}

	return &supplierv1.ListSuppliersResponse{
		Suppliers: pbSuppliers,
		Total:     total,
	}, nil
}

// domainToPb converts a domain Supplier to a protobuf Supplier
func domainToPb(s *domain.Supplier) *supplierv1.Supplier {
	if s == nil {
		return nil
	}

	return &supplierv1.Supplier{
		Id:            s.ID.Hex(),
		Name:          s.Name,
		ContactPerson: s.ContactPerson,
		Email:         s.Email,
		Phone:         s.Phone,
		Address:       s.Address,
		City:          s.City,
		State:         s.State,
		Country:       s.Country,
		PostalCode:    s.PostalCode,
		TaxId:         s.TaxID,
		Website:       s.Website,
		Currency:      s.Currency,
		LeadTimeDays:  s.LeadTimeDays,
		PaymentTerms:  s.PaymentTerms,
		Metadata:      s.Metadata,
		CreatedAt:     timestamppb.New(s.CreatedAt),
		UpdatedAt:     timestamppb.New(s.UpdatedAt),
	}
}
