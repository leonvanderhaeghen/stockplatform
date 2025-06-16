package mocks

import (
	"context"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockLocationRepository is a mock implementation of the LocationRepository interface
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) Create(ctx context.Context, location *domain.StoreLocation) error {
	args := m.Called(ctx, location)
	return args.Error(0)
}

func (m *MockLocationRepository) GetByID(ctx context.Context, id string) (*domain.StoreLocation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StoreLocation), args.Error(1)
}

func (m *MockLocationRepository) GetByName(ctx context.Context, name string) (*domain.StoreLocation, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StoreLocation), args.Error(1)
}

func (m *MockLocationRepository) Update(ctx context.Context, location *domain.StoreLocation) error {
	args := m.Called(ctx, location)
	return args.Error(0)
}

func (m *MockLocationRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLocationRepository) List(ctx context.Context, limit int, offset int, includeInactive bool) ([]*domain.StoreLocation, error) {
	args := m.Called(ctx, limit, offset, includeInactive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StoreLocation), args.Error(1)
}

func (m *MockLocationRepository) ListByType(ctx context.Context, locationType string, limit int, offset int) ([]*domain.StoreLocation, error) {
	args := m.Called(ctx, locationType, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StoreLocation), args.Error(1)
}
