package application

import (
	"context"
	"fmt"

	inventoryclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/inventory"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
	"go.uber.org/zap"
)

// OrderInventoryService coordinates between order and inventory services
type OrderInventoryService struct {
	inventoryClient *inventoryclient.Client
	orderService    *OrderService
	logger          *zap.Logger
}

// NewOrderInventoryService creates a new OrderInventoryService
func NewOrderInventoryService(
	orderService *OrderService,
	inventoryServiceAddr string,
	logger *zap.Logger,
) (*OrderInventoryService, error) {
	// Initialize the inventory client using new abstraction
	invCfg := inventoryclient.Config{Address: inventoryServiceAddr}
	inventoryClient, err := inventoryclient.New(invCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory client: %w", err)
	}

	return &OrderInventoryService{
		inventoryClient: inventoryClient,
		orderService:    orderService,
		logger:          logger.Named("order_inventory_service"),
	}, nil
}

// Close closes any open connections
func (s *OrderInventoryService) Close() error {
	if s.inventoryClient != nil {
		return s.inventoryClient.Close()
	}
	return nil
}

// CreateOrderWithInventoryCheck creates a new order after checking inventory
func (s *OrderInventoryService) CreateOrderWithInventoryCheck(
	ctx context.Context,
	input *domain.Order,
) (*domain.Order, error) {
	// Check inventory for all items
	for _, item := range input.Items {
		inventory, err := s.inventoryClient.GetInventoryByProductID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to check inventory for product %s: %w", item.ProductID, err)
		}

		// Check available stock
		available := inventory.GetInventory().GetQuantity() - inventory.GetInventory().GetReserved()
		if available < int32(item.Quantity) {
			return nil, fmt.Errorf("insufficient stock for product %s: available %d, requested %d",
				item.ProductID, available, item.Quantity)
		}
	}

	// Create the order with the provided items and addresses
	order, err := s.orderService.CreateOrder(ctx, input.UserID, input.Items, input.ShippingAddr, input.BillingAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Update order status to indicate inventory is being processed
	order.Status = domain.StatusPending // Using the constant from the domain package
	err = s.orderService.UpdateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	return order, nil
}

// ProcessOrderFulfillment handles the fulfillment of an order
func (s *OrderInventoryService) ProcessOrderFulfillment(
	ctx context.Context,
	orderID string,
) error {
	// Get the order
	order, err := s.orderService.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Update order status to indicate fulfillment is in progress
	order.Status = domain.StatusPaid
	err = s.orderService.UpdateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Process each item in the order
	for _, item := range order.Items {
		// Reserve the inventory for this item
		_, err := s.inventoryClient.ReserveStock(ctx, &inventorypb.ReserveStockRequest{
			Id:       item.ProductID,
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			// If we can't reserve stock, update the order status and return an error
			order.Status = domain.StatusCancelled
			err = s.orderService.UpdateOrder(ctx, order)
			if err != nil {
				s.logger.Error("Failed to update order status after inventory reservation failure", 
					zap.Error(err))
			}
			return fmt.Errorf("failed to reserve inventory for product %s: %w", 
				item.ProductID, err)
		}

		// Simulate order processing (in a real system, this would involve shipping, etc.)
		s.logger.Info("Processing order item",
			zap.String("order_id", orderID),
			zap.String("product_id", item.ProductID),
			zap.Int32("quantity", int32(item.Quantity)))
	}

	// Update order status to completed
	order.Status = domain.StatusShipped
	err = s.orderService.UpdateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to complete order: %w", err)
	}

	return nil
}
