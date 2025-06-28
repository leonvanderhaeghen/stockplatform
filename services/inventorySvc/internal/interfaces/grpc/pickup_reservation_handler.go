package grpc

import (
	"context"

	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

// ReserveForPickup handles reservation requests for in-store pickup
func (s *InventoryServer) ReserveForPickup(ctx context.Context, req *inventoryv1.ReserveForPickupRequest) (*inventoryv1.ReserveForPickupResponse, error) {
	logger := s.logger.With(
		zap.String("handler", "ReserveForPickup"),
		zap.String("order_id", req.OrderId),
		zap.String("location_id", req.LocationId),
		zap.Int("items_count", len(req.Items)),
	)

	logger.Info("Processing ReserveForPickup request")

	// Validate request
	if req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	if req.LocationId == "" {
		return nil, status.Error(codes.InvalidArgument, "location ID is required")
	}

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one item is required")
	}

	// Process reservation for each item
	reservationResults := make([]*inventoryv1.InventoryReservationResult, 0, len(req.Items))
	allSuccess := true
	resStatus := "success"

	for _, item := range req.Items {
		result := &inventoryv1.InventoryReservationResult{
			ProductId:         item.ProductId,
			Sku:               item.Sku,
			RequestedQuantity: item.Quantity,
			ReservedQuantity:  0,
			Status:            "unavailable",
			InventoryItemId:   item.InventoryItemId,
		}

		// Validate item
		if item.Quantity <= 0 {
			result.ErrorMessage = "quantity must be positive"
			reservationResults = append(reservationResults, result)
			allSuccess = false
			resStatus = "failed"
			continue
		}

		// Reserve the inventory
		err := s.service.ReserveStock(ctx, item.InventoryItemId, item.Quantity)
		if err != nil {
			logger.Error("Failed to reserve stock", 
				zap.Error(err),
				zap.String("product_id", item.ProductId),
				zap.String("inventory_item_id", item.InventoryItemId),
				zap.Int32("quantity", item.Quantity),
			)

			result.Status = "unavailable"
			result.ErrorMessage = err.Error()
			allSuccess = false
			resStatus = "failed"
		} else {
			result.Status = "reserved"
			result.ReservedQuantity = item.Quantity
		}

		reservationResults = append(reservationResults, result)
	}

	logger.Info("ReserveForPickup request completed",
		zap.Bool("all_success", allSuccess),
	)

	// Generate a unique reservation ID
	reservationID := generateReservationID(req.OrderId, req.LocationId)
	
	// Calculate expiration date (24 hours from now)
	expirationDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	if req.ExpirationDate != "" {
		expirationDate = req.ExpirationDate
	}

	return &inventoryv1.ReserveForPickupResponse{
		ReservationId:  reservationID,
		Status:         resStatus,
		ExpirationDate: expirationDate,
		Items:          reservationResults,
	}, nil
}

// generateReservationID creates a unique reservation ID
func generateReservationID(orderID, locationID string) string {
	return orderID + "-" + locationID + "-" + time.Now().Format("20060102150405")
}

// CancelPickup handles cancellation of in-store pickup reservations
func (s *InventoryServer) CancelPickup(ctx context.Context, req *inventoryv1.CancelPickupRequest) (*inventoryv1.CancelPickupResponse, error) {
	logger := s.logger.With(
		zap.String("handler", "CancelPickup"),
		zap.String("reservation_id", req.ReservationId),
	)

	logger.Info("Processing CancelPickup request")

	// Validate request
	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation ID is required")
	}

	// Parse the reservation ID to extract order and location
	parts := strings.Split(req.ReservationId, "-")
	if len(parts) < 2 {
		return nil, status.Error(codes.InvalidArgument, "invalid reservation ID format")
	}

	orderID := parts[0]
	locationID := parts[1]

	// Release the reservation in the service
	err := s.service.CancelPickupReservation(ctx, orderID, locationID, req.Reason)
	if err != nil {
		logger.Error("Failed to cancel pickup reservation", 
			zap.Error(err),
			zap.String("reservation_id", req.ReservationId),
		)
		return nil, status.Errorf(codes.Internal, "failed to cancel reservation: %v", err)
	}

	logger.Info("CancelPickup request completed successfully")

	return &inventoryv1.CancelPickupResponse{
		Success:     true,
		CancelledAt: time.Now().Format(time.RFC3339),
	}, nil
}

// CompletePickup handles fulfillment of in-store pickup reservations
func (s *InventoryServer) CompletePickup(ctx context.Context, req *inventoryv1.CompletePickupRequest) (*inventoryv1.CompletePickupResponse, error) {
	logger := s.logger.With(
		zap.String("handler", "CompletePickup"),
		zap.String("reservation_id", req.ReservationId),
		zap.String("staff_id", req.StaffId),
	)

	logger.Info("Processing CompletePickup request")

	// Validate request
	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation ID is required")
	}

	// Parse the reservation ID to extract order and location
	parts := strings.Split(req.ReservationId, "-")
	if len(parts) < 2 {
		return nil, status.Error(codes.InvalidArgument, "invalid reservation ID format")
	}

	orderID := parts[0]
	locationID := parts[1]

	// Complete the pickup in the service
	err := s.service.CompletePickup(ctx, orderID, locationID, req.StaffId, req.Notes)
	if err != nil {
		logger.Error("Failed to complete pickup", 
			zap.Error(err),
			zap.String("reservation_id", req.ReservationId),
		)
		return nil, status.Errorf(codes.Internal, "failed to complete pickup: %v", err)
	}

	logger.Info("CompletePickup request completed successfully")

	transactionID := "TRX-" + req.ReservationId + "-" + time.Now().Format("20060102150405")

	return &inventoryv1.CompletePickupResponse{
		Success:      true,
		TransactionId: transactionID,
		CompletedAt:  time.Now().Format(time.RFC3339),
	}, nil
}
