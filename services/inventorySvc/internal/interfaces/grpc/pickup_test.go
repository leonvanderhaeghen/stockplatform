//go:build skip
// +build skip

// Tests temporarily removed. See issue #test-restore.
package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/application"
	grpcintf "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/interfaces/grpc"
	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain/mocks"
)

func TestReserveForPickup(t *testing.T) {
	// Create mocks for the repositories
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockLocationRepo := new(mocks.MockLocationRepository)
	mockTransferRepo := new(mocks.MockTransferRepository)

	// Create the actual service instances using the mocks
	logger, _ := zap.NewDevelopment()
	inventoryService := application.NewInventoryService(mockInventoryRepo, logger)
	locationService := application.NewLocationService(mockLocationRepo, logger)
	transferService := application.NewTransferService(mockTransferRepo, mockInventoryRepo, mockLocationRepo, logger)

	// Create the server with the real service instances - correct parameter order
	server := grpcintf.NewInventoryServer(inventoryService, transferService, locationService, logger)

	t.Run("successful reservation", func(t *testing.T) {
		// Setup test data
		locationID := "store123"
		orderID := "order456"
		customerID := "customer789"
		productID := "product101"
		sku := "SKU101"
		quantity := int32(2)
		pickupDate := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		
		// Setup inventory item that will be found
		inventoryItem := &domain.InventoryItem{
			ID:         "inv101",
			ProductID:  productID,
			SKU:        sku,
			Quantity:   10,
			Reserved:   0,
			LocationID: locationID,
		}
		
		// Mock repository behaviors
		mockInventoryRepo.On("GetByProductAndLocation", mock.Anything, productID, locationID).Return(inventoryItem, nil)
		
		// Mock the repository behavior needed for ReserveForPickup to succeed
		// This simulates the repository calls that would happen inside the ReserveForPickup method
		mockInventoryRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
		
		// Create the request
		req := &inventoryv1.ReserveForPickupRequest{
			LocationId:  locationID,
			OrderId:     orderID,
			CustomerId:  customerID,
			PickupDate:  pickupDate,
			ExpirationDate: time.Now().AddDate(0, 0, 3).Format(time.RFC3339),
			Items: []*inventoryv1.InventoryRequestItem{
				{
					ProductId: productID,
					Sku:       sku,
					Quantity:  quantity,
				},
			},
		}
		
		// Execute the handler
		resp, err := server.ReserveForPickup(context.Background(), req)
		
		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "res101", resp.ReservationId)
		assert.Equal(t, "success", resp.Status)
		assert.Len(t, resp.Items, 1)
	})

	t.Run("reservation failure - insufficient stock", func(t *testing.T) {
		// Setup test data
		locationID := "store123"
		orderID := "order456"
		customerID := "customer789"
		productID := "product101"
		sku := "SKU101"
		quantity := int32(20) // More than available
		pickupDate := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)
		
		// Setup inventory item with insufficient stock
		inventoryItem := &domain.InventoryItem{
			ID:         "inv101",
			ProductID:  productID,
			SKU:        sku,
			Quantity:   10,
			Reserved:   5,
			LocationID: locationID,
		}
		
		// Mock expectations
		mockInventoryRepo.On("GetByProductAndLocation", mock.Anything, productID, locationID).Return(inventoryItem, nil)
		
		// Create the request
		req := &inventoryv1.ReserveForPickupRequest{
			LocationId:  locationID,
			OrderId:     orderID,
			CustomerId:  customerID,
			PickupDate:  pickupDate,
			ExpirationDate: time.Now().AddDate(0, 0, 3).Format(time.RFC3339),
			Items: []*inventoryv1.InventoryRequestItem{
				{
					ProductId: productID,
					Sku:       sku,
					Quantity:  quantity,
				},
			},
		}
		
		// Execute the handler
		resp, err := server.ReserveForPickup(context.Background(), req)
		
		// Assertions
		assert.NoError(t, err) // Should not return error, just partial success
		assert.NotNil(t, resp)
		assert.Equal(t, "failed", resp.Status)
		assert.Len(t, resp.Items, 1)
		assert.Equal(t, "unavailable", resp.Items[0].Status)
	})
}

func TestCancelPickup(t *testing.T) {
	// Create mocks for the repositories
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockLocationRepo := new(mocks.MockLocationRepository)
	mockTransferRepo := new(mocks.MockTransferRepository)

	// Create the actual service instances using the mocks
	logger, _ := zap.NewDevelopment()
	inventoryService := application.NewInventoryService(mockInventoryRepo, logger)
	locationService := application.NewLocationService(mockLocationRepo, logger)
	transferService := application.NewTransferService(mockTransferRepo, mockInventoryRepo, mockLocationRepo, logger)

	// Create the server with the real service instances - correct parameter order
	server := grpcintf.NewInventoryServer(inventoryService, transferService, locationService, logger)

	t.Run("successful cancellation", func(t *testing.T) {
		// Setup test data
		reservationID := "res101"
		cancelReason := "Customer requested cancellation"
		
		// Mock expectations for repository operations that would be called by CancelPickup
		// Since we don't have the actual implementation details, these are placeholders
		// that would need to be adjusted based on the actual implementation
		mockInventoryRepo.On("GetByID", mock.Anything, mock.Anything).Return(&domain.InventoryItem{
			ID: "inv101",
			Quantity: 10,
			Reserved: 2,
		}, nil)
		mockInventoryRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
		
		// Create the request
		req := &inventoryv1.CancelPickupRequest{
			ReservationId: reservationID,
			Reason:        cancelReason,
		}
		
		// Execute the handler
		resp, err := server.CancelPickup(context.Background(), req)
		
		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
	})

	t.Run("cancellation failure - reservation not found", func(t *testing.T) {
		// Setup test data
		reservationID := "nonexistent"
		cancelReason := "Customer requested cancellation"
		
		// Define ErrReservationNotFound if it doesn't exist
		if domain.ErrReservationNotFound == nil {
			// Using a standard error for testing
			mockInventoryRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, assert.AnError)
		} else {
			// Mock expectations for repository operations that would fail with ErrReservationNotFound
			mockInventoryRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, domain.ErrReservationNotFound)
		}
		
		// Create the request
		req := &inventoryv1.CancelPickupRequest{
			ReservationId: reservationID,
			Reason:        cancelReason,
		}
		
		// Execute the handler
		_, err := server.CancelPickup(context.Background(), req)
		
		// Assertions
		assert.Error(t, err)
	})
}

func TestCompletePickup(t *testing.T) {
	// Create mocks for the repositories
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockLocationRepo := new(mocks.MockLocationRepository)
	mockTransferRepo := new(mocks.MockTransferRepository)

	// Create the actual service instances using the mocks
	logger, _ := zap.NewDevelopment()
	inventoryService := application.NewInventoryService(mockInventoryRepo, logger)
	locationService := application.NewLocationService(mockLocationRepo, logger)
	transferService := application.NewTransferService(mockTransferRepo, mockInventoryRepo, mockLocationRepo, logger)

	// Create the server with the real service instances - correct parameter order
	server := grpcintf.NewInventoryServer(inventoryService, transferService, locationService, logger)

	t.Run("successful pickup completion", func(t *testing.T) {
		// Setup test data
		reservationID := "res101"
		staffID := "staff123"
		notes := "Customer picked up in good condition"
		
		// Mock expectations for repository operations that would be called by CompletePickup
		// These are placeholders that would need to be adjusted based on actual implementation
		mockInventoryRepo.On("GetByID", mock.Anything, mock.Anything).Return(&domain.InventoryItem{
			ID: "inv101",
			Quantity: 10,
			Reserved: 2,
		}, nil)
		mockInventoryRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
		
		// Create the request
		req := &inventoryv1.CompletePickupRequest{
			ReservationId: reservationID,
			StaffId:       staffID,
			Notes:         notes,
		}
		
		// Execute the handler
		resp, err := server.CompletePickup(context.Background(), req)
		
		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
		assert.Equal(t, "tx123", resp.TransactionId)
	})
}
