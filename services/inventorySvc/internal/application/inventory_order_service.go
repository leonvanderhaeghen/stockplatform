package application

import (
	"context"
	"fmt"

	"github.com/leonvanderhaeghen/stockplatform/pkg/grpcclient"
	orderpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/order/v1"
	"go.uber.org/zap"
)

// InventoryOrderService coordinates between inventory and order services
type InventoryOrderService struct {
	orderClient    *grpcclient.OrderClient
	inventoryService *InventoryService
	logger         *zap.Logger
}

// NewInventoryOrderService creates a new InventoryOrderService
func NewInventoryOrderService(
	inventoryService *InventoryService,
	orderServiceAddr string,
	logger *zap.Logger,
) (*InventoryOrderService, error) {
	// Initialize the order client
	orderClient, err := grpcclient.NewOrderClient(orderServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create order client: %w", err)
	}

	return &InventoryOrderService{
		orderClient:     orderClient,
		inventoryService: inventoryService,
		logger:          logger.Named("inventory_order_service"),
	}, nil
}

// Close closes any open connections
func (s *InventoryOrderService) Close() error {
	if s.orderClient != nil {
		return s.orderClient.Close()
	}
	return nil
}

// ProcessOrderUpdate handles inventory updates when an order status changes
func (s *InventoryOrderService) ProcessOrderUpdate(
	ctx context.Context,
	orderID string,
	newStatus string,
) error {
	// Get the order details
	orderResp, err := s.orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: orderID})
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	order := orderResp.GetOrder()
	
	// Handle different order statuses
	switch newStatus {
	case "processing":
		// Reserve inventory when order is being processed
		for _, item := range order.GetItems() {
			err := s.inventoryService.ReserveStock(ctx, item.GetProductId(), int32(item.GetQuantity()))
			if err != nil {
				return fmt.Errorf("failed to reserve stock for product %s: %w", 
					item.GetProductId(), err)
			}
		}

	case "shipped":
		// Fulfill reservation when order is shipped
		for _, item := range order.GetItems() {
			err := s.inventoryService.FulfillReservation(ctx, item.GetProductId(), int32(item.GetQuantity()))
			if err != nil {
				return fmt.Errorf("failed to fulfill reservation for product %s: %w", 
					item.GetProductId(), err)
			}
		}

	case "cancelled":
		// Release reservation when order is cancelled
		for _, item := range order.GetItems() {
			err := s.inventoryService.ReleaseReservation(ctx, item.GetProductId(), int32(item.GetQuantity()))
			if err != nil {
				return fmt.Errorf("failed to release reservation for product %s: %w", 
					item.GetProductId(), err)
			}
		}
	}

	return nil
}

// CheckInventoryForOrder checks if there's enough inventory for an order
func (s *InventoryOrderService) CheckInventoryForOrder(
	ctx context.Context,
	orderID string,
) (bool, error) {
	// Get the order details
	orderResp, err := s.orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: orderID})
	if err != nil {
		return false, fmt.Errorf("failed to get order: %w", err)
	}

	// Check inventory for each item
	for _, item := range orderResp.GetOrder().GetItems() {
		inventory, err := s.inventoryService.GetInventoryItemByProductID(ctx, item.GetProductId())
		if err != nil {
			return false, fmt.Errorf("failed to get inventory for product %s: %w", 
				item.GetProductId(), err)
		}

		// Check if there's enough available stock
		if !inventory.IsAvailable(int32(item.GetQuantity())) {
			s.logger.Warn("Insufficient stock",
				zap.String("product_id", item.GetProductId()),
				zap.Int32("available", inventory.Quantity - inventory.Reserved),
				zap.Int32("requested", int32(item.GetQuantity())))
			return false, nil
		}
	}

	return true, nil
}
