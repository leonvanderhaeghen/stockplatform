package mocks

import (
	"context"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockTransferRepository is a mock implementation of the TransferRepository interface
type MockTransferRepository struct {
	mock.Mock
}

func (m *MockTransferRepository) Create(ctx context.Context, transfer *domain.Transfer) error {
	args := m.Called(ctx, transfer)
	return args.Error(0)
}

func (m *MockTransferRepository) GetByID(ctx context.Context, id string) (*domain.Transfer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transfer), args.Error(1)
}

func (m *MockTransferRepository) Update(ctx context.Context, transfer *domain.Transfer) error {
	args := m.Called(ctx, transfer)
	return args.Error(0)
}

func (m *MockTransferRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTransferRepository) List(ctx context.Context, limit int, offset int) ([]*domain.Transfer, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transfer), args.Error(1)
}

func (m *MockTransferRepository) ListByStatus(ctx context.Context, status domain.TransferStatus, limit int, offset int) ([]*domain.Transfer, error) {
	args := m.Called(ctx, status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transfer), args.Error(1)
}

func (m *MockTransferRepository) ListPendingTransfers(ctx context.Context, limit int, offset int) ([]*domain.Transfer, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transfer), args.Error(1)
}

func (m *MockTransferRepository) ListBySourceLocation(ctx context.Context, locationID string, limit int, offset int) ([]*domain.Transfer, error) {
	args := m.Called(ctx, locationID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transfer), args.Error(1)
}

func (m *MockTransferRepository) ListByDestLocation(ctx context.Context, locationID string, limit int, offset int) ([]*domain.Transfer, error) {
	args := m.Called(ctx, locationID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transfer), args.Error(1)
}

func (m *MockTransferRepository) ListByProduct(ctx context.Context, productID string, limit int, offset int) ([]*domain.Transfer, error) {
	args := m.Called(ctx, productID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transfer), args.Error(1)
}
