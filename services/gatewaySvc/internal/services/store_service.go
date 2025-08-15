package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	storeclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/store"
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
	s.logger.Debug("ListStores", zap.Int("limit", limit), zap.Int("offset", offset))

	resp, err := s.client.ListStores(ctx, int32(limit), int32(offset))
	if err != nil {
		s.logger.Error("Failed to list stores", zap.Error(err))
		return nil, fmt.Errorf("failed to list stores: %w", err)
	}

	return resp, nil
}

// GetStore retrieves a store by ID
func (s *StoreServiceImpl) GetStore(ctx context.Context, id string) (interface{}, error) {
	s.logger.Debug("GetStore", zap.String("id", id))

	resp, err := s.client.GetStore(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get store", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	return resp, nil
}

// CreateStore creates a new store
func (s *StoreServiceImpl) CreateStore(ctx context.Context, name, description, street, city, state, country, postalCode, phone, email string) (interface{}, error) {
	s.logger.Debug("CreateStore",
		zap.String("name", name),
		zap.String("description", description),
		zap.String("street", street),
		zap.String("city", city))

	resp, err := s.client.CreateStore(ctx, name, description, street, city, state, country, postalCode, phone, email)
	if err != nil {
		s.logger.Error("Failed to create store", zap.Error(err))
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return resp, nil
}
