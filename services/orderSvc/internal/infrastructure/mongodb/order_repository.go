package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// OrderRepository implements the domain.OrderRepository interface
type OrderRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewOrderRepository creates a new MongoDB order repository
func NewOrderRepository(db *mongo.Database, collectionName string, logger *zap.Logger) domain.OrderRepository {
	collection := db.Collection(collectionName)
	
	// Create indexes for improved query performance
	indexModels := []mongo.IndexModel{
		mongo.IndexModel{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		mongo.IndexModel{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		mongo.IndexModel{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetUnique(false),
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		logger.Warn("Failed to create indexes", zap.Error(err))
	}
	
	return &OrderRepository{
		collection: collection,
		logger:     logger.Named("order_repository"),
	}
}

// Create adds a new order
func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	r.logger.Debug("Creating order", 
		zap.String("id", order.ID),
		zap.String("user_id", order.UserID),
	)
	
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		r.logger.Error("Failed to create order", 
			zap.Error(err),
			zap.String("id", order.ID),
		)
		return err
	}
	
	return nil
}

// GetByID finds an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	r.logger.Debug("Getting order by ID", zap.String("id", id))
	
	var order domain.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Order not found", zap.String("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get order", 
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, err
	}
	
	return &order, nil
}

// GetByUserID finds orders for a specific user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Order, error) {
	r.logger.Debug("Getting orders by user ID", 
		zap.String("user_id", userID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, findOptions)
	if err != nil {
		r.logger.Error("Failed to get orders by user ID", 
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var orders []*domain.Order
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			r.logger.Error("Failed to decode order", zap.Error(err))
			return nil, err
		}
		orders = append(orders, &order)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while getting orders by user ID", zap.Error(err))
		return nil, err
	}
	
	return orders, nil
}

// Update updates an existing order
func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	r.logger.Debug("Updating order", 
		zap.String("id", order.ID),
		zap.String("user_id", order.UserID),
	)
	
	order.UpdatedAt = time.Now()
	
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": order.ID}, order)
	if err != nil {
		r.logger.Error("Failed to update order", 
			zap.Error(err),
			zap.String("id", order.ID),
		)
		return err
	}
	
	if result.MatchedCount == 0 {
		r.logger.Warn("No order was updated", zap.String("id", order.ID))
		return errors.New("order not found")
	}
	
	return nil
}

// UpdateWithOptimisticLock updates an order with version checking to prevent concurrent modifications
func (r *OrderRepository) UpdateWithOptimisticLock(ctx context.Context, order *domain.Order, expectedVersion int32) error {
	r.logger.Debug("Updating order with optimistic lock", 
		zap.String("id", order.ID),
		zap.String("user_id", order.UserID),
		zap.Int32("expected_version", expectedVersion),
		zap.Int32("new_version", order.Version),
	)
	
	// Update the order's timestamp
	order.UpdatedAt = time.Now()
	
	// Use MongoDB's findOneAndReplace with version filter to ensure atomic update
	filter := bson.M{
		"_id":     order.ID,
		"version": expectedVersion,
	}
	
	result, err := r.collection.ReplaceOne(ctx, filter, order)
	if err != nil {
		r.logger.Error("Failed to update order with optimistic lock", 
			zap.Error(err),
			zap.String("id", order.ID),
			zap.Int32("expected_version", expectedVersion),
		)
		return err
	}
	
	// Check if no documents were matched (either order doesn't exist or version mismatch)
	if result.MatchedCount == 0 {
		// Check if order exists to determine the specific error
		var existingOrder domain.Order
		err := r.collection.FindOne(ctx, bson.M{"_id": order.ID}).Decode(&existingOrder)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				r.logger.Warn("Order not found during optimistic lock update", zap.String("id", order.ID))
				return errors.New("order not found")
			}
			r.logger.Error("Failed to check order existence during optimistic lock", zap.Error(err))
			return err
		}
		
		// Order exists but version mismatch
		r.logger.Warn("Optimistic lock failed due to version mismatch", 
			zap.String("id", order.ID),
			zap.Int32("expected_version", expectedVersion),
			zap.Int32("actual_version", existingOrder.Version),
		)
		return domain.ErrOptimisticLockFailed
	}
	
	r.logger.Debug("Successfully updated order with optimistic lock", 
		zap.String("id", order.ID),
		zap.Int32("new_version", order.Version),
	)
	
	return nil
}

// Delete removes an order
func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting order", zap.String("id", id))
	
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete order", 
			zap.Error(err),
			zap.String("id", id),
		)
		return err
	}
	
	if result.DeletedCount == 0 {
		r.logger.Warn("No order was deleted", zap.String("id", id))
		return errors.New("order not found")
	}
	
	return nil
}

// List returns all orders with optional filtering and pagination
func (r *OrderRepository) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*domain.Order, error) {
	r.logger.Debug("Listing orders", 
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))
	
	// Convert map to bson.M
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}
	
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		r.logger.Error("Failed to list orders", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var orders []*domain.Order
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			r.logger.Error("Failed to decode order", zap.Error(err))
			return nil, err
		}
		orders = append(orders, &order)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing orders", zap.Error(err))
		return nil, err
	}
	
	return orders, nil
}

// Count returns the number of orders matching a filter
func (r *OrderRepository) Count(ctx context.Context, filter map[string]interface{}) (int64, error) {
	r.logger.Debug("Counting orders")
	
	// Convert map to bson.M
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}
	
	count, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		r.logger.Error("Failed to count orders", zap.Error(err))
		return 0, err
	}
	
	return count, nil
}
