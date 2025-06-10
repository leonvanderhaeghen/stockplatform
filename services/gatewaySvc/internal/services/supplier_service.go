package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/supplier/v1"
)

type supplierService struct {
	client supplierv1.SupplierServiceClient
	conn   *grpc.ClientConn
}

// NewSupplierService creates a new supplier service client
func NewSupplierService(addr string, logger *zap.Logger) (SupplierService, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := supplierv1.NewSupplierServiceClient(conn)

	return &supplierService{
		client: client,
		conn:   conn,
	}, nil
}

func (s *supplierService) CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.Supplier, error) {
	resp, err := s.client.CreateSupplier(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Supplier, nil
}

func (s *supplierService) GetSupplier(ctx context.Context, id string) (*supplierv1.Supplier, error) {
	req := &supplierv1.GetSupplierRequest{Id: id}
	resp, err := s.client.GetSupplier(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Supplier, nil
}

func (s *supplierService) UpdateSupplier(ctx context.Context, req *supplierv1.UpdateSupplierRequest) (*supplierv1.Supplier, error) {
	resp, err := s.client.UpdateSupplier(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Supplier, nil
}

func (s *supplierService) DeleteSupplier(ctx context.Context, id string) error {
	req := &supplierv1.DeleteSupplierRequest{Id: id}
	_, err := s.client.DeleteSupplier(ctx, req)
	return err
}

func (s *supplierService) ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*supplierv1.Supplier, int32, error) {
	req := &supplierv1.ListSuppliersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}
	resp, err := s.client.ListSuppliers(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Suppliers, resp.Total, nil
}

func (s *supplierService) Close() error {
	return s.conn.Close()
}

func (s *supplierService) ListAdapters(ctx context.Context) ([]*supplierv1.SupplierAdapter, error) {
	req := &supplierv1.ListAdaptersRequest{}
	resp, err := s.client.ListAdapters(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Adapters, nil
}

func (s *supplierService) GetAdapterCapabilities(ctx context.Context, adapterName string) (*supplierv1.AdapterCapabilities, error) {
	req := &supplierv1.GetAdapterCapabilitiesRequest{
		AdapterName: adapterName,
	}
	resp, err := s.client.GetAdapterCapabilities(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Capabilities, nil
}

func (s *supplierService) TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error {
	req := &supplierv1.TestAdapterConnectionRequest{
		AdapterName: adapterName,
		Config:      config,
	}
	resp, err := s.client.TestAdapterConnection(ctx, req)
	if err != nil {
		return err
	}
	
	if !resp.Success {
		return fmt.Errorf("connection test failed: %s", resp.Message)
	}
	return nil
}

func (s *supplierService) SyncProducts(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error) {
	req := &supplierv1.SyncProductsRequest{
		SupplierId: supplierID,
		Options:    options,
	}
	resp, err := s.client.SyncProducts(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.JobId, nil
}

func (s *supplierService) SyncInventory(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error) {
	req := &supplierv1.SyncInventoryRequest{
		SupplierId: supplierID,
		Options:    options,
	}
	resp, err := s.client.SyncInventory(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.JobId, nil
}
