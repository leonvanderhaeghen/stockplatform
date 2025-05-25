package application

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"stockplatform/services/inventorySvc/internal/domain"
)

// InventoryService handles business logic for inventory operations
type InventoryService struct {
	repo   domain.InventoryRepository
	logger *zap.Logger
}

// NewInventoryService creates a new inventory service
func NewInventoryService(repo domain.InventoryRepository, logger *zap.Logger) *InventoryService {
	return &InventoryService{
		repo:   repo,
		logger: logger.Named("inventory_service"),
	}
}

// CreateInventoryItem creates a new inventory item
func (s *InventoryService) CreateInventoryItem(ctx context.Context, productID string, quantity int32, sku string, location string) (*domain.InventoryItem, error) {
	s.logger.Info("Creating inventory item",
		zap.String("product_id", productID),
		zap.Int32("quantity", quantity),
		zap.String("sku", sku),
	)

	// Check if inventory item with this SKU already exists
	existingItem, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}

	if existingItem != nil {
		s.logger.Warn("Inventory item with this SKU already exists",
			zap.String("sku", sku),
		)
		return nil, errors.New("inventory item with this SKU already exists")
	}

	item := domain.NewInventoryItem(productID, quantity, sku, location)
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetInventoryItem retrieves an inventory item by ID
func (s *InventoryService) GetInventoryItem(ctx context.Context, id string) (*domain.InventoryItem, error) {
	s.logger.Debug("Getting inventory item", zap.String("id", id))
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if item == nil {
		return nil, errors.New("inventory item not found")
	}
	
	return item, nil
}

// GetInventoryItemByProductID retrieves an inventory item by product ID
func (s *InventoryService) GetInventoryItemByProductID(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	s.logger.Debug("Getting inventory item by product ID", 
		zap.String("product_id", productID),
	)
	
	item, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	if item == nil {
		return nil, errors.New("inventory item not found")
	}
	
	return item, nil
}

// GetInventoryItemBySKU retrieves an inventory item by SKU
func (s *InventoryService) GetInventoryItemBySKU(ctx context.Context, sku string) (*domain.InventoryItem, error) {
	s.logger.Debug("Getting inventory item by SKU", zap.String("sku", sku))
	
	item, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}
	
	if item == nil {
		return nil, errors.New("inventory item not found")
	}
	
	return item, nil
}

// UpdateInventoryItem updates an existing inventory item
func (s *InventoryService) UpdateInventoryItem(ctx context.Context, item *domain.InventoryItem) error {
	s.logger.Info("Updating inventory item",
		zap.String("id", item.ID),
		zap.String("product_id", item.ProductID),
	)
	
	return s.repo.Update(ctx, item)
}

// DeleteInventoryItem removes an inventory item
func (s *InventoryService) DeleteInventoryItem(ctx context.Context, id string) error {
	s.logger.Info("Deleting inventory item", zap.String("id", id))
	
	return s.repo.Delete(ctx, id)
}

// ListInventoryItems returns all inventory items with pagination
func (s *InventoryService) ListInventoryItems(ctx context.Context, limit, offset int) ([]*domain.InventoryItem, error) {
	s.logger.Debug("Listing inventory items",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	return s.repo.List(ctx, limit, offset)
}

// AddStock increases the quantity of an inventory item
func (s *InventoryService) AddStock(ctx context.Context, id string, quantity int32) error {
	s.logger.Info("Adding stock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	item.AddStock(quantity)
	return s.repo.Update(ctx, item)
}

// RemoveStock decreases the quantity of an inventory item
func (s *InventoryService) RemoveStock(ctx context.Context, id string, quantity int32) error {
	s.logger.Info("Removing stock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	if !item.RemoveStock(quantity) {
		return errors.New("insufficient stock")
	}
	
	return s.repo.Update(ctx, item)
}

// ReserveStock reserves stock for an order
func (s *InventoryService) ReserveStock(ctx context.Context, id string, quantity int32) error {
	s.logger.Info("Reserving stock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	if !item.Reserve(quantity) {
		return errors.New("insufficient stock available")
	}
	
	return s.repo.Update(ctx, item)
}

// ReleaseReservation releases a reservation without fulfilling it
func (s *InventoryService) ReleaseReservation(ctx context.Context, id string, quantity int32) error {
	s.logger.Info("Releasing reservation",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	item.ReleaseReservation(quantity)
	return s.repo.Update(ctx, item)
}

// FulfillReservation completes a reservation and deducts from stock
func (s *InventoryService) FulfillReservation(ctx context.Context, id string, quantity int32) error {
	s.logger.Info("Fulfilling reservation",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	if !item.FulfillReservation(quantity) {
		return errors.New("insufficient reserved quantity")
	}
	
	return s.repo.Update(ctx, item)
}
