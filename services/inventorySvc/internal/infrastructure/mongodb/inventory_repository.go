package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
)

// InventoryRepository implements the domain.InventoryRepository interface
type InventoryRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewInventoryRepository creates a new MongoDB inventory repository
func NewInventoryRepository(db *mongo.Database, collectionName string, logger *zap.Logger) domain.InventoryRepository {
	collection := db.Collection(collectionName)
	
	// Create indexes for improved query performance
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "product_id", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.D{{Key: "sku", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "inventory_id", Value: 1}, {Key: "created_at", Value: -1}},
			Options: options.Index().SetUnique(false),
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		logger.Warn("Failed to create indexes", zap.Error(err))
	}
	
	return &InventoryRepository{
		collection: collection,
		logger:     logger.Named("inventory_repository"),
	}
}

// Create adds a new inventory item
func (r *InventoryRepository) Create(ctx context.Context, item *domain.InventoryItem) error {
	r.logger.Debug("Creating inventory item", 
		zap.String("product_id", item.ProductID),
		zap.String("sku", item.SKU),
	)
	
	_, err := r.collection.InsertOne(ctx, item)
	if err != nil {
		r.logger.Error("Failed to create inventory item", 
			zap.Error(err),
			zap.String("product_id", item.ProductID),
		)
		return err
	}
	
	return nil
}

// GetByID finds an inventory item by its ID
func (r *InventoryRepository) GetByID(ctx context.Context, id string) (*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory item by ID", zap.String("id", id))
	
	var item domain.InventoryItem
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Inventory item not found", zap.String("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get inventory item", 
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, err
	}
	
	return &item, nil
}

// GetByProductID finds inventory items by product ID
func (r *InventoryRepository) GetByProductID(ctx context.Context, productID string) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory items by product ID", 
		zap.String("product_id", productID),
	)
	
	findOptions := options.Find()
	
	cursor, err := r.collection.Find(ctx, bson.M{"product_id": productID}, findOptions)
	if err != nil {
		r.logger.Error("Failed to get inventory items by product ID", 
			zap.Error(err),
			zap.String("product_id", productID),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while getting inventory items by product ID", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// GetByProductAndLocation finds inventory items by product ID and location ID
func (r *InventoryRepository) GetByProductAndLocation(ctx context.Context, productID, locationID string) (*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory item by product and location", 
		zap.String("product_id", productID),
		zap.String("location_id", locationID),
	)
	
	var item domain.InventoryItem
	filter := bson.M{"product_id": productID, "location_id": locationID}

	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Inventory item not found by product and location", 
				zap.String("product_id", productID),
				zap.String("location_id", locationID),
			)
			return nil, nil
		}
		r.logger.Error("Failed to get inventory item by product and location", 
			zap.Error(err),
			zap.String("product_id", productID),
			zap.String("location_id", locationID),
		)
		return nil, err
	}
	
	return &item, nil
}

// Deprecated: ListInventoryItemsByProduct is replaced by GetByProductID and kept for compatibility
func (r *InventoryRepository) ListInventoryItemsByProduct(ctx context.Context, productID string) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Listing inventory items by product", 
		zap.String("product_id", productID),
	)
	
	findOptions := options.Find()
	
	cursor, err := r.collection.Find(ctx, bson.M{"product_id": productID}, findOptions)
	if err != nil {
		r.logger.Error("Failed to list inventory items by product", 
			zap.Error(err),
			zap.String("product_id", productID),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing inventory items by product", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// GetBySKU finds inventory items by SKU (across all locations)
func (r *InventoryRepository) GetBySKU(ctx context.Context, sku string) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory items by SKU", zap.String("sku", sku))
	
	findOptions := options.Find()
	
	cursor, err := r.collection.Find(ctx, bson.M{"sku": sku}, findOptions)
	if err != nil {
		r.logger.Error("Failed to get inventory items by SKU", 
			zap.Error(err),
			zap.String("sku", sku),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while getting inventory items by SKU", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// GetBySKUAndLocation finds an inventory item by SKU and location ID
func (r *InventoryRepository) GetBySKUAndLocation(ctx context.Context, sku, locationID string) (*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory item by SKU and location",
		zap.String("sku", sku),
		zap.String("location_id", locationID),
	)
	
	var item domain.InventoryItem
	filter := bson.M{"sku": sku, "location_id": locationID}

	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Inventory item not found by SKU and location",
				zap.String("sku", sku),
				zap.String("location_id", locationID),
			)
			return nil, nil
		}
		r.logger.Error("Failed to get inventory item by SKU and location",
			zap.Error(err),
			zap.String("sku", sku),
			zap.String("location_id", locationID),
		)
		return nil, err
	}
	
	return &item, nil
}

// Deprecated: ListInventoryItemsBySKU is replaced by GetBySKU and kept for compatibility
func (r *InventoryRepository) ListInventoryItemsBySKU(ctx context.Context, sku string) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Listing inventory items by SKU",
		zap.String("sku", sku),
	)
	
	findOptions := options.Find()
	
	cursor, err := r.collection.Find(ctx, bson.M{"sku": sku}, findOptions)
	if err != nil {
		r.logger.Error("Failed to list inventory items by SKU",
			zap.Error(err),
			zap.String("sku", sku),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing inventory items by SKU", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// Update updates an existing inventory item
func (r *InventoryRepository) Update(ctx context.Context, item *domain.InventoryItem) error {
	r.logger.Debug("Updating inventory item", 
		zap.String("id", item.ID),
		zap.String("product_id", item.ProductID),
	)
	
	item.LastUpdated = time.Now()
	
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": item.ID}, item)
	if err != nil {
		r.logger.Error("Failed to update inventory item", 
			zap.Error(err),
			zap.String("id", item.ID),
		)
		return err
	}
	
	if result.MatchedCount == 0 {
		r.logger.Warn("No inventory item was updated", zap.String("id", item.ID))
		return errors.New("inventory item not found")
	}
	
	return nil
}

// Delete removes an inventory item
func (r *InventoryRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting inventory item", zap.String("id", id))
	
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete inventory item", 
			zap.Error(err),
			zap.String("id", id),
		)
		return err
	}
	
	if result.DeletedCount == 0 {
		r.logger.Warn("No inventory item was deleted", zap.String("id", id))
		return errors.New("inventory item not found")
	}
	
	return nil
}

// ListByLocation returns inventory items for a specific location with pagination
func (r *InventoryRepository) ListByLocation(ctx context.Context, locationID string, limit, offset int) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Listing inventory items by location",
		zap.String("location_id", locationID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	
	cursor, err := r.collection.Find(ctx, bson.M{"location_id": locationID}, findOptions)
	if err != nil {
		r.logger.Error("Failed to list inventory items by location", 
			zap.Error(err),
			zap.String("location_id", locationID),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing inventory items by location", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// List returns all inventory items with optional pagination
func (r *InventoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Listing inventory items", 
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		r.logger.Error("Failed to list inventory items", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing inventory items", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// ListLowStock returns inventory items that are below their reorder point
func (r *InventoryRepository) ListLowStock(ctx context.Context, limit, offset int) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Listing low stock inventory items", 
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	
	filter := bson.M{
		"$expr": bson.M{
			"$lt": []interface{}{
				"$quantity", 
				"$reorder_threshold",
			},
		},
	}
	
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("Failed to list low stock inventory items", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing low stock inventory items", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// ListByStockStatus returns inventory items based on stock status
func (r *InventoryRepository) ListByStockStatus(ctx context.Context, status string, limit, offset int) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Listing inventory items by stock status", 
		zap.String("status", status),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	
	var filter bson.M
	
	switch status {
	case "in_stock":
		// Items with quantity > 0
		filter = bson.M{"quantity": bson.M{"$gt": 0}}
		
	case "low_stock":
		// Items below reorder threshold but not zero
		filter = bson.M{
			"quantity": bson.M{"$gt": 0},
			"$expr": bson.M{
				"$lt": []interface{}{
					"$quantity", 
					"$reorder_threshold",
				},
			},
		}
		
	case "out_of_stock":
		// Items with zero quantity
		filter = bson.M{"quantity": 0}
		
	default:
		r.logger.Warn("Invalid stock status provided", zap.String("status", status))
		return nil, domain.ErrInvalidInput
	}
	
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("Failed to list inventory items by stock status", 
			zap.Error(err),
			zap.String("status", status),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing inventory items by stock status", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// GetByOrderAndLocation finds inventory items reserved for a specific order at a location
func (r *InventoryRepository) GetByOrderAndLocation(ctx context.Context, orderID, locationID string) ([]*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory items by order and location",
		zap.String("order_id", orderID),
		zap.String("location_id", locationID),
	)
	
	// Find items that have reservations for this order at this location
	filter := bson.M{
		"location_id": locationID,
		"reservations.order_id": orderID,
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to get inventory items by order and location",
			zap.Error(err),
			zap.String("order_id", orderID),
			zap.String("location_id", locationID),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var items []*domain.InventoryItem
	for cursor.Next(ctx) {
		var item domain.InventoryItem
		if err := cursor.Decode(&item); err != nil {
			r.logger.Error("Failed to decode inventory item", zap.Error(err))
			return nil, err
		}
		items = append(items, &item)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while getting inventory items by order and location", zap.Error(err))
		return nil, err
	}
	
	return items, nil
}

// GetHistory retrieves the history of changes for a specific inventory item
func (r *InventoryRepository) GetHistory(ctx context.Context, inventoryID string, limit, offset int32) ([]*domain.InventoryHistory, int32, error) {
	r.logger.Debug("Getting inventory history", 
		zap.String("inventory_id", inventoryID),
		zap.Int32("limit", limit),
		zap.Int32("offset", offset),
	)

	// First, get the total count for pagination
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{"inventory_id": inventoryID})
	if err != nil {
		r.logger.Error("Failed to count inventory history", 
			zap.String("inventory_id", inventoryID),
			zap.Error(err),
		)
		return nil, 0, err
	}

	// Then get the paginated results
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Most recent first
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	if offset > 0 {
		opts.SetSkip(int64(offset))
	}

	cursor, err := r.collection.Find(ctx, bson.M{"inventory_id": inventoryID}, opts)
	if err != nil {
		r.logger.Error("Failed to find inventory history", 
			zap.String("inventory_id", inventoryID),
			zap.Error(err),
		)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var history []*domain.InventoryHistory
	if err := cursor.All(ctx, &history); err != nil {
		r.logger.Error("Failed to decode inventory history", 
			zap.String("inventory_id", inventoryID),
			zap.Error(err),
		)
		return nil, 0, err
	}

	return history, int32(totalCount), nil
}

// RecordHistory adds a new history entry for an inventory item
func (r *InventoryRepository) RecordHistory(ctx context.Context, history *domain.InventoryHistory) error {
	history.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, history)
	if err != nil {
		r.logger.Error("Failed to record inventory history", 
			zap.String("inventory_id", history.InventoryID,
			zap.String("change_type", history.ChangeType),
			zap.Error(err)),
		)
		return err
	}

	r.logger.Debug("Recorded inventory history", 
		zap.String("inventory_id", history.InventoryID),
		zap.String("change_type", history.ChangeType),
	)
	return nil
}

// AdjustStock adjusts stock with a reason and user identification
func (r *InventoryRepository) AdjustStock(ctx context.Context, id string, quantity int32, reason, performedBy string) error {
	r.logger.Debug("Adjusting stock",
		zap.String("id", id),
		zap.Int32("quantity", quantity),
		zap.String("reason", reason),
		zap.String("performed_by", performedBy),
	)
	
	// First, get the current inventory item
	item, err := r.GetByID(ctx, id)
	if err != nil {
		r.logger.Error("Failed to get inventory item for stock adjustment",
			zap.Error(err),
			zap.String("id", id),
		)
		return err
	}
	
	if item == nil {
		r.logger.Error("Inventory item not found for stock adjustment",
			zap.String("id", id),
		)
		return domain.ErrNotFound
	}
	
	// Apply the adjustment
	prevQuantity := item.Quantity
	item.Quantity += quantity
	item.LastUpdated = time.Now()
	
	// Don't allow negative stock
	if item.Quantity < 0 {
		r.logger.Warn("Stock adjustment would result in negative quantity",
			zap.String("id", id),
			zap.Int32("prev_quantity", prevQuantity),
			zap.Int32("adjustment", quantity),
		)
		return domain.ErrInsufficientStock
	}
	
	// Create update document
	update := bson.M{
		"$set": bson.M{
			"quantity": item.Quantity,
			"last_updated": item.LastUpdated,
		},
		"$push": bson.M{
			"stock_adjustments": bson.M{
				"adjustment": quantity,
				"previous_quantity": prevQuantity,
				"new_quantity": item.Quantity,
				"reason": reason,
				"performed_by": performedBy,
				"timestamp": item.LastUpdated,
			},
		},
	}
	
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		r.logger.Error("Failed to update inventory for stock adjustment",
			zap.Error(err),
			zap.String("id", id),
		)
		return err
	}
	
	if result.MatchedCount == 0 {
		r.logger.Warn("No inventory item was updated during stock adjustment",
			zap.String("id", id),
		)
		return domain.ErrNotFound
	}
	
	r.logger.Info("Stock adjustment completed",
		zap.String("id", id),
		zap.Int32("prev_quantity", prevQuantity),
		zap.Int32("new_quantity", item.Quantity),
		zap.String("reason", reason),
	)
	
	return nil
}
