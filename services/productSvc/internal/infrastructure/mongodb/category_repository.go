package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
)

type categoryRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewCategoryRepository creates a new MongoDB category repository
func NewCategoryRepository(db *mongo.Database, logger *zap.Logger) domain.CategoryRepository {
	return &categoryRepository{
		collection: db.Collection("categories"),
		logger:     logger.Named("mongodb_category_repository"),
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		r.logger.Error("Failed to create category", zap.Error(err))
		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var category domain.Category
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to get category", zap.Error(err))
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	category.UpdatedAt = time.Now()

	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": category.ID}, category)
	if err != nil {
		r.logger.Error("Failed to update category", zap.Error(err))
		return err
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("Failed to delete category", zap.Error(err))
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *categoryRepository) List(ctx context.Context, parentID string, depth int32) ([]*domain.Category, error) {
	filter := bson.M{}

	if parentID != "" {
		objectID, err := primitive.ObjectIDFromHex(parentID)
		if err != nil {
			return nil, domain.ErrInvalidID
		}
		filter["parent_id"] = objectID
	} else {
		filter["$or"] = []bson.M{{"parent_id": ""}, {"parent_id": bson.M{"$exists": false}}}
	}

	if depth > 0 {
		filter["level"] = bson.M{"$lte": depth}
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to list categories", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*domain.Category
	if err := cursor.All(ctx, &categories); err != nil {
		r.logger.Error("Failed to decode categories", zap.Error(err))
		return nil, err
	}

	return categories, nil
}
