package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
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
func (s *InventoryService) CreateInventoryItem(ctx context.Context, productID string, quantity int32, sku string, locationID string) (*domain.InventoryItem, error) {
	s.logger.Info("Creating inventory item",
		zap.String("product_id", productID),
		zap.Int32("quantity", quantity),
		zap.String("sku", sku),
		zap.String("location_id", locationID),
	)

	// Check if inventory item with this SKU already exists at this location
	existingItem, err := s.repo.GetBySKUAndLocation(ctx, sku, locationID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	if existingItem != nil {
		s.logger.Warn("Inventory item with this SKU already exists at this location",
			zap.String("sku", sku),
			zap.String("location_id", locationID),
		)
		return nil, errors.New("inventory item with this SKU already exists at this location")
	}

	item := domain.NewInventoryItem(productID, quantity, sku, locationID)
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

// GetInventoryItemsByProductID retrieves inventory items by product ID across all locations
func (s *InventoryService) GetInventoryItemsByProductID(ctx context.Context, productID string) ([]*domain.InventoryItem, error) {
	s.logger.Debug("Getting inventory items by product ID", 
		zap.String("product_id", productID),
	)
	
	items, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	if len(items) == 0 {
		// Return empty slice rather than nil to maintain consistent return type
		return []*domain.InventoryItem{}, nil
	}
	
	return items, nil
}

// GetInventoryItemByProductAndLocation retrieves an inventory item by product ID and location
func (s *InventoryService) GetInventoryItemByProductAndLocation(ctx context.Context, productID string, locationID string) (*domain.InventoryItem, error) {
	s.logger.Debug("Getting inventory item by product ID and location", 
		zap.String("product_id", productID),
		zap.String("location_id", locationID),
	)
	
	item, err := s.repo.GetByProductAndLocation(ctx, productID, locationID)
	if err != nil {
		return nil, err
	}
	
	if item == nil {
		return nil, errors.New("inventory item not found")
	}
	
	return item, nil
}

// GetInventoryItemsBySKU retrieves inventory items by SKU across all locations
func (s *InventoryService) GetInventoryItemsBySKU(ctx context.Context, sku string) ([]*domain.InventoryItem, error) {
	s.logger.Debug("Getting inventory items by SKU", zap.String("sku", sku))
	
	items, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, err
	}
	
	if len(items) == 0 {
		// Return empty slice rather than nil to maintain consistent return type
		return []*domain.InventoryItem{}, nil
	}
	
	return items, nil
}

// GetInventoryItemBySKUAndLocation retrieves an inventory item by SKU and location
func (s *InventoryService) GetInventoryItemBySKUAndLocation(ctx context.Context, sku string, locationID string) (*domain.InventoryItem, error) {
	s.logger.Debug("Getting inventory item by SKU and location", 
		zap.String("sku", sku),
		zap.String("location_id", locationID),
	)
	
	item, err := s.repo.GetBySKUAndLocation(ctx, sku, locationID)
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

// ListInventoryItemsByLocation returns inventory items for a specific location
func (s *InventoryService) ListInventoryItemsByLocation(ctx context.Context, locationID string, limit, offset int) ([]*domain.InventoryItem, error) {
	s.logger.Debug("Listing inventory items by location",
		zap.String("location_id", locationID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	return s.repo.ListByLocation(ctx, locationID, limit, offset)
}

// ListLowStockItems returns inventory items that are below their reorder point
func (s *InventoryService) ListLowStockItems(ctx context.Context, limit, offset int) ([]*domain.InventoryItem, error) {
	s.logger.Debug("Listing low stock items",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	return s.repo.ListLowStock(ctx, limit, offset)
}

// SetReorderParameters sets inventory reordering parameters
func (s *InventoryService) SetReorderParameters(ctx context.Context, id string, minimumStock, maximumStock, reorderPoint, reorderQuantity int32) error {
	s.logger.Info("Setting reorder parameters",
		zap.String("id", id),
		zap.Int32("minimum_stock", minimumStock),
		zap.Int32("maximum_stock", maximumStock),
		zap.Int32("reorder_point", reorderPoint),
		zap.Int32("reorder_quantity", reorderQuantity),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	item.SetReorderParameters(minimumStock, maximumStock, reorderPoint, reorderQuantity)
	return s.repo.Update(ctx, item)
}

// SetShelfLocation sets the precise shelf location within a store
func (s *InventoryService) SetShelfLocation(ctx context.Context, id string, shelfLocation string) error {
	s.logger.Info("Setting shelf location",
		zap.String("id", id),
		zap.String("shelf_location", shelfLocation),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	item.SetShelfLocation(shelfLocation)
	return s.repo.Update(ctx, item)
}

// ScheduleInventoryCount schedules the next inventory count date
func (s *InventoryService) ScheduleInventoryCount(ctx context.Context, id string, nextCountDate string) error {
	s.logger.Info("Scheduling inventory count",
		zap.String("id", id),
		zap.String("next_count_date", nextCountDate),
	)
	
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if item == nil {
		return errors.New("inventory item not found")
	}
	
	// Parse the date string
	parsedDate, err := time.Parse(time.RFC3339, nextCountDate)
	if err != nil {
		return errors.New("invalid date format, use RFC3339")
	}
	
	item.ScheduleInventoryCount(parsedDate)
	return s.repo.Update(ctx, item)
}

// AdjustStock adjusts inventory quantity and records reason
func (s *InventoryService) AdjustStock(ctx context.Context, id string, quantity int32, reason string, performedBy string) error {
	s.logger.Info("Adjusting stock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
		zap.String("reason", reason),
		zap.String("performed_by", performedBy),
	)
	
	return s.repo.AdjustStock(ctx, id, quantity, reason, performedBy)
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

// CancelPickupReservation cancels a pickup reservation for an order at a specific location
func (s *InventoryService) CancelPickupReservation(ctx context.Context, orderID string, locationID string, reason string) error {
	s.logger.Info("Cancelling pickup reservation",
		zap.String("order_id", orderID),
		zap.String("location_id", locationID),
		zap.String("reason", reason),
	)

	// Get inventory items associated with this order at this location
	inventoryItems, err := s.repo.GetByOrderAndLocation(ctx, orderID, locationID)
	if err != nil {
		return err
	}

	if len(inventoryItems) == 0 {
		return errors.New("no reserved inventory items found for this order at this location")
	}

	// Release the reservations for all items
	for _, item := range inventoryItems {
		item.SetReservationStatus(domain.ReservationStatusCancelled)
		err := s.repo.Update(ctx, item)
		if err != nil {
			return err
		}
	}

	return nil
}

// CompletePickup completes an in-store pickup reservation
func (s *InventoryService) CompletePickup(ctx context.Context, orderID string, locationID string, staffID string, notes string) error {
	s.logger.Info("Completing pickup reservation",
		zap.String("order_id", orderID),
		zap.String("location_id", locationID),
		zap.String("staff_id", staffID),
	)

	// Get inventory items associated with this order at this location
	inventoryItems, err := s.repo.GetByOrderAndLocation(ctx, orderID, locationID)
	if err != nil {
		return err
	}

	if len(inventoryItems) == 0 {
		return errors.New("no reserved inventory items found for this order at this location")
	}

	// Process each item for fulfillment
	for _, item := range inventoryItems {
		// Mark as fulfilled and update the inventory
		if !item.FulfillReservation(item.Reserved) {
			return errors.New("insufficient reserved quantity for item: " + item.ID)
		}
		
		// Set reservation status to fulfilled
		item.SetReservationStatus(domain.ReservationStatusFulfilled)
		if notes != "" {
			item.AddNote(fmt.Sprintf("Pickup completed by staff %s: %s", staffID, notes))
		}
		
		// Update the item in repository
		err := s.repo.Update(ctx, item)
		if err != nil {
			return err
		}
	}

	return nil
}

// PerformPOSInventoryCheck checks inventory availability for POS transactions
func (s *InventoryService) PerformPOSInventoryCheck(
	ctx context.Context, 
	locationID string, 
	items []domain.InventoryCheckItem,
) ([]*domain.InventoryCheckResult, error) {
	s.logger.Info("Performing POS inventory check",
		zap.String("location_id", locationID),
		zap.Int("item_count", len(items)),
	)

	if locationID == "" {
		return nil, errors.New("location ID is required")
	}

	results := make([]*domain.InventoryCheckResult, 0, len(items))

	for _, checkItem := range items {
		result := &domain.InventoryCheckResult{
			ProductID: checkItem.ProductID,
			SKU:       checkItem.SKU,
			Quantity:  checkItem.Quantity,
			Available: false,
		}

		// First try to find by product ID if provided
		var inventoryItem *domain.InventoryItem
		var err error

		if checkItem.ProductID != "" {
			inventoryItem, err = s.GetInventoryItemByProductAndLocation(ctx, checkItem.ProductID, locationID)
		} else if checkItem.SKU != "" {
			inventoryItem, err = s.GetInventoryItemBySKUAndLocation(ctx, checkItem.SKU, locationID)
		} else {
			result.Error = "Either product ID or SKU must be provided"
			results = append(results, result)
			continue
		}

		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				result.Error = "Item not found"
			} else {
				result.Error = err.Error()
			}
			results = append(results, result)
			continue
		}

		// Check availability
		result.AvailableQuantity = inventoryItem.GetAvailable()
		result.Location = inventoryItem.LocationID
		result.Available = inventoryItem.HasSufficientStock(checkItem.Quantity)

		results = append(results, result)
	}

	return results, nil
}

// ReserveForPOSTransaction reserves inventory items for a POS transaction
func (s *InventoryService) ReserveForPOSTransaction(
	ctx context.Context, 
	orderID string,
	locationID string,
	items []domain.ReservationItem,
) error {
	s.logger.Info("Reserving inventory for POS transaction",
		zap.String("order_id", orderID),
		zap.String("location_id", locationID),
		zap.Int("item_count", len(items)),
	)

	if orderID == "" {
		return errors.New("order ID is required")
	}
	if locationID == "" {
		return errors.New("location ID is required")
	}
	if len(items) == 0 {
		return errors.New("at least one item must be specified")
	}

	for _, reserveItem := range items {
		// Find the inventory item
		var inventoryItem *domain.InventoryItem
		var err error

		if reserveItem.ProductID != "" {
			inventoryItem, err = s.GetInventoryItemByProductAndLocation(ctx, reserveItem.ProductID, locationID)
		} else if reserveItem.SKU != "" {
			inventoryItem, err = s.GetInventoryItemBySKUAndLocation(ctx, reserveItem.SKU, locationID)
		} else {
			return errors.New("either product ID or SKU must be provided for each item")
		}

		if err != nil {
			return err
		}

		// Reserve the stock
		if !inventoryItem.ReserveForOrder(reserveItem.Quantity, orderID) {
			return errors.New("insufficient stock available for item: " + inventoryItem.ProductID)
		}

		// Update the inventory
		if err := s.repo.Update(ctx, inventoryItem); err != nil {
			return err
		}
	}

	return nil
}

// DeductForDirectPOSTransaction immediately deducts inventory for direct POS sales (no reservation)
func (s *InventoryService) DeductForDirectPOSTransaction(
	ctx context.Context,
	locationID string,
	staffID string,
	items []domain.InventoryCheckItem,
) error {
	s.logger.Info("Deducting inventory for direct POS transaction",
		zap.String("location_id", locationID),
		zap.String("staff_id", staffID),
		zap.Int("item_count", len(items)),
	)

	if locationID == "" {
		return errors.New("location ID is required")
	}
	if staffID == "" {
		return errors.New("staff ID is required")
	}
	if len(items) == 0 {
		return errors.New("at least one item must be specified")
	}

	for _, checkItem := range items {
		// Find the inventory item
		var inventoryItem *domain.InventoryItem
		var err error

		if checkItem.ProductID != "" {
			inventoryItem, err = s.GetInventoryItemByProductAndLocation(ctx, checkItem.ProductID, locationID)
		} else if checkItem.SKU != "" {
			inventoryItem, err = s.GetInventoryItemBySKUAndLocation(ctx, checkItem.SKU, locationID)
		} else {
			return errors.New("either product ID or SKU must be provided for each item")
		}

		if err != nil {
			return err
		}

		// Directly remove the stock (no reservation)
		if !inventoryItem.RemoveStock(checkItem.Quantity) {
			return errors.New("insufficient stock available for item: " + inventoryItem.ProductID)
		}

		// Add a note about the POS transaction
		inventoryItem.AddNote(fmt.Sprintf("POS transaction: %d units sold by staff %s", checkItem.Quantity, staffID))

		// Update the inventory
		if err := s.repo.Update(ctx, inventoryItem); err != nil {
			return err
		}
	}

	return nil
}

// End of service methods
