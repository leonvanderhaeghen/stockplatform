package application

import (
	"context"
	"errors"
	"testing"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInventoryItemsByProductID(t *testing.T) {
	// Set up mock repository
	mockRepo := new(mocks.MockInventoryRepository)
	logger, _ := NewLogger("test")
	service := NewInventoryService(mockRepo, logger)

	// Test data
	productID := "product123"
	locationID := "location456"
	inv1 := domain.InventoryItem{
		ID:        "inv1",
		ProductID: productID,
		LocationID: locationID,
		Quantity:  10,
	}
	inv2 := domain.InventoryItem{
		ID:        "inv2",
		ProductID: productID,
		LocationID: "anotherlocation",
		Quantity:  5,
	}

	// Test cases
	tests := []struct {
		name           string
		productID      string
		locationID     string
		mockReturnInvs []domain.InventoryItem
		mockReturnErr  error
		expectedInvs   []domain.InventoryItem
		expectedErr    error
	}{
		{
			name:           "success with no location filter",
			productID:      productID,
			locationID:     "",
			mockReturnInvs: []domain.InventoryItem{inv1, inv2},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{inv1, inv2},
			expectedErr:    nil,
		},
		{
			name:           "success with location filter",
			productID:      productID,
			locationID:     locationID,
			mockReturnInvs: []domain.InventoryItem{inv1},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{inv1},
			expectedErr:    nil,
		},
		{
			name:           "not found",
			productID:      "nonexistent",
			locationID:     "",
			mockReturnInvs: []domain.InventoryItem{},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{},
			expectedErr:    nil,
		},
		{
			name:           "repository error",
			productID:      productID,
			locationID:     "",
			mockReturnInvs: nil,
			mockReturnErr:  errors.New("database error"),
			expectedInvs:   nil,
			expectedErr:    errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock expectations
			if tc.locationID == "" {
				mockRepo.On("GetByProductID", mock.Anything, tc.productID).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()
			} else {
				mockRepo.On("GetByProductAndLocation", mock.Anything, tc.productID, tc.locationID).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()
			}

			// Call the service method
			invs, err := service.GetInventoryItemsByProductID(context.Background(), tc.productID, tc.locationID)

			// Check expectations
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInvs, invs)
			}

			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetInventoryItemsBySKU(t *testing.T) {
	// Set up mock repository
	mockRepo := new(mocks.MockInventoryRepository)
	logger, _ := NewLogger("test")
	service := NewInventoryService(mockRepo, logger)

	// Test data
	sku := "SKU123"
	locationID := "location456"
	inv1 := domain.InventoryItem{
		ID:        "inv1",
		SKU:       sku,
		LocationID: locationID,
		Quantity:  10,
	}
	inv2 := domain.InventoryItem{
		ID:        "inv2",
		SKU:       sku,
		LocationID: "anotherlocation",
		Quantity:  5,
	}

	// Test cases
	tests := []struct {
		name           string
		sku            string
		locationID     string
		mockReturnInvs []domain.InventoryItem
		mockReturnErr  error
		expectedInvs   []domain.InventoryItem
		expectedErr    error
	}{
		{
			name:           "success with no location filter",
			sku:            sku,
			locationID:     "",
			mockReturnInvs: []domain.InventoryItem{inv1, inv2},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{inv1, inv2},
			expectedErr:    nil,
		},
		{
			name:           "success with location filter",
			sku:            sku,
			locationID:     locationID,
			mockReturnInvs: []domain.InventoryItem{inv1},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{inv1},
			expectedErr:    nil,
		},
		{
			name:           "not found",
			sku:            "nonexistent",
			locationID:     "",
			mockReturnInvs: []domain.InventoryItem{},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{},
			expectedErr:    nil,
		},
		{
			name:           "repository error",
			sku:            sku,
			locationID:     "",
			mockReturnInvs: nil,
			mockReturnErr:  errors.New("database error"),
			expectedInvs:   nil,
			expectedErr:    errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock expectations
			if tc.locationID == "" {
				mockRepo.On("GetBySKU", mock.Anything, tc.sku).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()
			} else {
				mockRepo.On("GetBySKUAndLocation", mock.Anything, tc.sku, tc.locationID).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()
			}

			// Call the service method
			invs, err := service.GetInventoryItemsBySKU(context.Background(), tc.sku, tc.locationID)

			// Check expectations
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInvs, invs)
			}

			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAdjustStock(t *testing.T) {
	// Set up mock repository
	mockRepo := new(mocks.MockInventoryRepository)
	logger, _ := NewLogger("test")
	service := NewInventoryService(mockRepo, logger)

	// Test data
	invID := "inv123"
	quantity := 5
	reason := "Restock from supplier"
	performedBy := "user123"

	// Test cases
	tests := []struct {
		name          string
		id            string
		quantity      int
		reason        string
		performedBy   string
		operation     string
		mockReturnErr error
		expectedErr   error
	}{
		{
			name:          "add stock success",
			id:            invID,
			quantity:      quantity,
			reason:        reason,
			performedBy:   performedBy,
			operation:     "add",
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name:          "remove stock success",
			id:            invID,
			quantity:      quantity,
			reason:        reason,
			performedBy:   performedBy,
			operation:     "remove",
			mockReturnErr: nil,
			expectedErr:   nil,
		},
		{
			name:          "add stock error",
			id:            invID,
			quantity:      quantity,
			reason:        reason,
			performedBy:   performedBy,
			operation:     "add",
			mockReturnErr: errors.New("database error"),
			expectedErr:   errors.New("database error"),
		},
		{
			name:          "remove stock error",
			id:            invID,
			quantity:      quantity,
			reason:        reason,
			performedBy:   performedBy,
			operation:     "remove",
			mockReturnErr: errors.New("insufficient stock"),
			expectedErr:   errors.New("insufficient stock"),
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock expectations
			if tc.operation == "add" {
				mockRepo.On("AdjustStock", mock.Anything, tc.id, tc.quantity, tc.reason, tc.performedBy).Return(tc.mockReturnErr).Once()
			} else {
				mockRepo.On("AdjustStock", mock.Anything, tc.id, -tc.quantity, tc.reason, tc.performedBy).Return(tc.mockReturnErr).Once()
			}

			// Call the service method
			var err error
			if tc.operation == "add" {
				err = service.AddStock(context.Background(), tc.id, tc.quantity, tc.reason, tc.performedBy)
			} else {
				err = service.RemoveStock(context.Background(), tc.id, tc.quantity, tc.reason, tc.performedBy)
			}

			// Check expectations
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListInventoryByLocation(t *testing.T) {
	// Set up mock repository
	mockRepo := new(mocks.MockInventoryRepository)
	logger, _ := NewLogger("test")
	service := NewInventoryService(mockRepo, logger)

	// Test data
	locationID := "location456"
	inv1 := domain.InventoryItem{
		ID:         "inv1",
		LocationID: locationID,
		Quantity:   10,
	}
	inv2 := domain.InventoryItem{
		ID:         "inv2",
		LocationID: locationID,
		Quantity:   5,
	}

	// Test cases
	tests := []struct {
		name           string
		locationID     string
		stockStatus    string
		limit          int
		offset         int
		mockReturnInvs []domain.InventoryItem
		mockReturnErr  error
		expectedInvs   []domain.InventoryItem
		expectedErr    error
	}{
		{
			name:           "success all stock",
			locationID:     locationID,
			stockStatus:    "all",
			limit:          10,
			offset:         0,
			mockReturnInvs: []domain.InventoryItem{inv1, inv2},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{inv1, inv2},
			expectedErr:    nil,
		},
		{
			name:           "success low stock",
			locationID:     locationID,
			stockStatus:    "low_stock",
			limit:          10,
			offset:         0,
			mockReturnInvs: []domain.InventoryItem{inv2},
			mockReturnErr:  nil,
			expectedInvs:   []domain.InventoryItem{inv2},
			expectedErr:    nil,
		},
		{
			name:           "repository error",
			locationID:     locationID,
			stockStatus:    "all",
			limit:          10,
			offset:         0,
			mockReturnInvs: nil,
			mockReturnErr:  errors.New("database error"),
			expectedInvs:   nil,
			expectedErr:    errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock expectations
			if tc.stockStatus == "low_stock" {
				mockRepo.On("ListLowStock", mock.Anything, tc.locationID, tc.limit, tc.offset).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()
			} else {
				mockRepo.On("ListByLocation", mock.Anything, tc.locationID, tc.limit, tc.offset).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()
			}

			// Call the service method
			invs, err := service.ListInventoryByLocation(context.Background(), tc.locationID, tc.stockStatus, tc.limit, tc.offset)

			// Check expectations
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInvs, invs)
			}

			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
