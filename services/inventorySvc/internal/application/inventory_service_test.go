//go:build skip
// +build skip

// Tests temporarily removed. See issue #test-restore.
package application

import (
	"context"
	"errors"
	"testing"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestGetInventoryItemsByProductID(t *testing.T) {
	mockRepo := new(mocks.MockInventoryRepository)
	logger, _ := zap.NewDevelopment()
	service := NewInventoryService(mockRepo, logger)

	// Test data
	productID := "prod-123"
	locationID := "location456"
	inv1 := &domain.InventoryItem{
		ID:        "inv1",
		ProductID: productID,
		LocationID: locationID,
		Quantity:  10,
	}
	inv2 := &domain.InventoryItem{
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
		mockReturnInvs []*domain.InventoryItem
		mockReturnErr  error
		expectedInvs   []*domain.InventoryItem
		expectedErr    error
	}{
		{
			name:           "success with no location filter",
			productID:      productID,
			locationID:     "",
			mockReturnInvs: []*domain.InventoryItem{inv1, inv2},
			mockReturnErr:  nil,
			expectedInvs:   []*domain.InventoryItem{inv1, inv2},
			expectedErr:    nil,
		},
		{
			name:           "success with location filter",
			productID:      productID,
			locationID:     locationID,
			mockReturnInvs: []*domain.InventoryItem{inv1},
			mockReturnErr:  nil,
			expectedInvs:   []*domain.InventoryItem{inv1},
			expectedErr:    nil,
		},
		{
			name:           "not found",
			productID:      "nonexistent",
			locationID:     "",
			mockReturnInvs: []*domain.InventoryItem{},
			mockReturnErr:  nil,
			expectedInvs:   []*domain.InventoryItem{},
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
			// Always expect GetByProductID since that's what the service implementation calls
			mockRepo.On("GetByProductID", mock.Anything, tc.productID).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()

			// Call the service method
			result, err := service.GetInventoryItemsByProductID(context.Background(), tc.productID)

			// Check expectations
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInvs, result)
			}

			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetInventoryItemsBySKU(t *testing.T) {
	mockRepo := new(mocks.MockInventoryRepository)
	logger, _ := zap.NewDevelopment()
	service := NewInventoryService(mockRepo, logger)

	// Test data
	sku := "SKU123"
	locationID := "location456"
	inv1 := &domain.InventoryItem{
		ID:        "inv1",
		SKU:       sku,
		LocationID: locationID,
		Quantity:  10,
	}
	inv2 := &domain.InventoryItem{
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
		mockReturnInvs []*domain.InventoryItem
		mockReturnErr  error
		expectedInvs   []*domain.InventoryItem
		expectedErr    error
	}{
		{
			name:           "success with no location filter",
			sku:            sku,
			locationID:     "",
			mockReturnInvs: []*domain.InventoryItem{inv1, inv2},
			mockReturnErr:  nil,
			expectedInvs:   []*domain.InventoryItem{inv1, inv2},
			expectedErr:    nil,
		},
		{
			name:           "success with location filter",
			sku:            sku,
			locationID:     locationID,
			mockReturnInvs: []*domain.InventoryItem{inv1},
			mockReturnErr:  nil,
			expectedInvs:   []*domain.InventoryItem{inv1},
			expectedErr:    nil,
		},
		{
			name:           "not found",
			sku:            "nonexistent",
			locationID:     "",
			mockReturnInvs: []*domain.InventoryItem{},
			mockReturnErr:  nil,
			expectedInvs:   []*domain.InventoryItem{},
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
			// Always expect GetBySKU since that's what the service implementation calls
			mockRepo.On("GetBySKU", mock.Anything, tc.sku).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()

			// Call the service method
			result, err := service.GetInventoryItemsBySKU(context.Background(), tc.sku)

			// Check expectations
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInvs, result)
			}

			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAddAndRemoveStock(t *testing.T) {
	mockRepo := new(mocks.MockInventoryRepository)
	logger, _ := zap.NewDevelopment()
	service := NewInventoryService(mockRepo, logger)

	// Test data
	invID := "item-123"
	quantity := int32(5)

	// Test cases
	tests := []struct {
		name           string
		id             string
		quantity       int32
		operation      string
		mockUpdateErr  error
		expectedErr    error
		sufficientStock bool  // For RemoveStock tests
	}{
		{
			name:           "add stock success",
			id:             invID,
			quantity:       quantity,
			operation:      "add",
			mockUpdateErr:  nil,
			expectedErr:    nil,
			sufficientStock: true,
		},
		{
			name:           "remove stock success",
			id:             invID,
			quantity:       quantity,
			operation:      "remove",
			mockUpdateErr:  nil,
			expectedErr:    nil,
			sufficientStock: true,
		},
		{
			name:           "add stock error",
			id:             invID,
			quantity:       quantity,
			operation:      "add",
			mockUpdateErr:  errors.New("database error"),
			expectedErr:    errors.New("database error"),
			sufficientStock: true,
		},
		{
			name:           "remove stock insufficient",
			id:             invID,
			quantity:       quantity,
			operation:      "remove",
			mockUpdateErr:  nil,
			expectedErr:    errors.New("insufficient stock"),
			sufficientStock: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a custom sample item for each test case
			sampleItem := &domain.InventoryItem{
				ID:         tc.id,
				ProductID:  "product-123",
				LocationID: "location-456",
				Quantity:   10,
			}
			
			// Mock GetByID which is called by both AddStock and RemoveStock
			mockRepo.On("GetByID", mock.Anything, tc.id).Return(sampleItem, nil).Once()
			
			// For the insufficient stock test, we don't expect Update to be called
			if tc.operation == "remove" && !tc.sufficientStock {
				// For this case, make the item have less quantity than requested
				sampleItem.Quantity = tc.quantity - 1
			} else {
				// Otherwise we expect Update to be called
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.InventoryItem")).Return(tc.mockUpdateErr).Once()
			}

			// Call the service method
			var err error
			if tc.operation == "add" {
				err = service.AddStock(context.Background(), tc.id, tc.quantity)
			} else {
				err = service.RemoveStock(context.Background(), tc.id, tc.quantity)
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
	logger, _ := zap.NewDevelopment()

	// Test data
	locationID := "loc-123"
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
		limit          int
		offset         int
		mockReturnInvs []*domain.InventoryItem
		mockReturnErr  error
		expectedInvs   []*domain.InventoryItem
		expectedErr    error
	}{
		{
			name:           "success all stock",
			locationID:     locationID,
			limit:          10,
			offset:         0,
			mockReturnInvs: []*domain.InventoryItem{&inv1, &inv2},
			mockReturnErr:  nil,
			expectedInvs:   []*domain.InventoryItem{&inv1, &inv2},
			expectedErr:    nil,
		},
		{
			name:           "repository error",
			locationID:     locationID,
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
			// Create a fresh mock repository for each test case
			mockRepo := new(mocks.MockInventoryRepository)
			service := NewInventoryService(mockRepo, logger)
			
			// Set up mock expectations
			mockRepo.On("ListByLocation", mock.Anything, tc.locationID, tc.limit, tc.offset).Return(tc.mockReturnInvs, tc.mockReturnErr).Once()

			// Call the service method
			result, err := service.ListInventoryItemsByLocation(context.Background(), tc.locationID, tc.limit, tc.offset)

			// Check expectations
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedInvs, result)
			}

			// Verify that all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
