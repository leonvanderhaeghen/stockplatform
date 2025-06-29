package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	storeclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/store"
	storev1 "github.com/leonvanderhaeghen/stockplatform/services/storeSvc/api/gen/go/api/proto/store/v1"
)

// StoreServiceImpl implements the StoreService interface
type StoreServiceImpl struct {
	client *storeclient.Client
	logger *zap.Logger
}

// NewStoreService creates a new instance of StoreServiceImpl
func NewStoreService(storeServiceAddr string, logger *zap.Logger) (StoreService, error) {
	// Create a gRPC client
	client, err := storeclient.NewClient(storeServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create store client: %w", err)
	}

	return &StoreServiceImpl{
		client: client,
		logger: logger,
	}, nil
}

// ListStores lists all stores with pagination
func (s *StoreServiceImpl) ListStores(ctx context.Context, limit, offset int) (interface{}, error) {
	s.logger.Debug("ListStores",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	// Create the request for listing stores
	req := &storev1.ListStoresRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// Call the gRPC method
	resp, err := s.client.ListStores(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list stores",
			zap.Int("limit", limit),
			zap.Int("offset", offset),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list stores: %w", err)
	}

	return resp.GetStores(), nil
}

// GetStore retrieves a store by ID
func (s *StoreServiceImpl) GetStore(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetStore",
		zap.String("id", id),
	)

	// Create the request for getting a store
	req := &storev1.GetStoreRequest{
		Id: id,
	}

	// Call the gRPC method
	resp, err := s.client.GetStore(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get store",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	return resp.GetStore(), nil
}
