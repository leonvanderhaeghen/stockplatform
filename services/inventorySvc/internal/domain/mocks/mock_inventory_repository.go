package mocks

import (
	"context"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockInventoryRepository is a mock implementation of the InventoryRepository interface
type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) Create(ctx context.Context, inventory domain.InventoryItem) (string, error) {
	args := m.Called(ctx, inventory)
	return args.String(0), args.Error(1)
}

func (m *MockInventoryRepository) GetByID(ctx context.Context, id string) (domain.InventoryItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) GetByProductID(ctx context.Context, productID string) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) GetByProductAndLocation(ctx context.Context, productID string, locationID string) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, productID, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) GetBySKU(ctx context.Context, sku string) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) GetBySKUAndLocation(ctx context.Context, sku string, locationID string) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, sku, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) Update(ctx context.Context, inventory domain.InventoryItem) error {
	args := m.Called(ctx, inventory)
	return args.Error(0)
}

func (m *MockInventoryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInventoryRepository) List(ctx context.Context, limit int, offset int) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) ListByStockStatus(ctx context.Context, status string, limit int, offset int) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) ListByLocation(ctx context.Context, locationID string, limit int, offset int) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, locationID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) ListLowStock(ctx context.Context, locationID string, limit int, offset int) ([]domain.InventoryItem, error) {
	args := m.Called(ctx, locationID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.InventoryItem), args.Error(1)
}

func (m *MockInventoryRepository) AdjustStock(ctx context.Context, id string, quantity int, reason string, performedBy string) error {
	args := m.Called(ctx, id, quantity, reason, performedBy)
	return args.Error(0)
}

func (m *MockInventoryRepository) ReserveStock(ctx context.Context, id string, quantity int, orderID string) error {
	args := m.Called(ctx, id, quantity, orderID)
	return args.Error(0)
}

func (m *MockInventoryRepository) ReleaseReservation(ctx context.Context, id string, quantity int, orderID string) error {
	args := m.Called(ctx, id, quantity, orderID)
	return args.Error(0)
}

func (m *MockInventoryRepository) FulfillReservation(ctx context.Context, id string, quantity int, orderID string) error {
	args := m.Called(ctx, id, quantity, orderID)
	return args.Error(0)
}
