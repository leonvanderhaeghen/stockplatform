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
func (r *InventoryRepository) GetByProductID(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory item by product ID", 
		zap.String("product_id", productID),
	)
	
	var item domain.InventoryItem
	err := r.collection.FindOne(ctx, bson.M{"product_id": productID}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Inventory item not found", 
				zap.String("product_id", productID),
			)
			return nil, nil
		}
		r.logger.Error("Failed to get inventory item", 
			zap.Error(err),
			zap.String("product_id", productID),
		)
		return nil, err
	}
	
	return &item, nil
}

// GetBySKU finds an inventory item by SKU
func (r *InventoryRepository) GetBySKU(ctx context.Context, sku string) (*domain.InventoryItem, error) {
	r.logger.Debug("Getting inventory item by SKU", zap.String("sku", sku))
	
	var item domain.InventoryItem
	err := r.collection.FindOne(ctx, bson.M{"sku": sku}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Inventory item not found", zap.String("sku", sku))
			return nil, nil
		}
		r.logger.Error("Failed to get inventory item", 
			zap.Error(err),
			zap.String("sku", sku),
		)
		return nil, err
	}
	
	return &item, nil
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
