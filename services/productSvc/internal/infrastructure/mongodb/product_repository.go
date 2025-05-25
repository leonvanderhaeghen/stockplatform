package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"stockplatform/services/productSvc/internal/domain"
)

type productRepository struct {
	collection *mongo.Collection
	logger    *zap.Logger
}

// NewProductRepository creates a new MongoDB product repository
func NewProductRepository(db *mongo.Database, collectionName string, logger *zap.Logger) domain.ProductRepository {
	collection := db.Collection(collectionName)
	return &productRepository{
		collection: collection,
		logger:     logger.Named("product_repository"),
	}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		r.logger.Error("Failed to create product", zap.Error(err))
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		product.ID = oid
	}

	return product, nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var product domain.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to get product by ID", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	if product.ID.IsZero() {
		return domain.ErrInvalidID
	}

	product.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"sku":         product.SKU,
			"category_id": product.CategoryID,
			"image_urls":  product.ImageURLs,
			"updated_at":  product.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateByID(
		ctx,
		product.ID,
		update,
	)

	if err != nil {
		r.logger.Error("Failed to update product", zap.String("id", product.ID.Hex()), zap.Error(err))
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		r.logger.Error("Failed to delete product", zap.String("id", id), zap.Error(err))
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *productRepository) List(
	ctx context.Context,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	// Build the filter
	filter := bson.M{}

	if opts != nil && opts.Filter != nil {
		// Apply ID filter
		if len(opts.Filter.IDs) > 0 {
			var objectIDs []primitive.ObjectID
			for _, id := range opts.Filter.IDs {
				objID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					continue // Skip invalid IDs
				}
				objectIDs = append(objectIDs, objID)
			}
			if len(objectIDs) > 0 {
				filter["_id"] = bson.M{"$in": objectIDs}
			}
		}

		// Apply category filter
		if len(opts.Filter.CategoryIDs) > 0 {
			filter["category_id"] = bson.M{"$in": opts.Filter.CategoryIDs}
		}

		// Apply price range filter
		priceFilter := bson.M{}
		if opts.Filter.MinPrice > 0 {
			priceFilter["$gte"] = opts.Filter.MinPrice
		}
		if opts.Filter.MaxPrice > 0 {
			priceFilter["$lte"] = opts.Filter.MaxPrice
		}
		if len(priceFilter) > 0 {
			filter["price"] = priceFilter
		}

		// Apply search term (case-insensitive search in name and description)
		if opts.Filter.SearchTerm != "" {
			filter["$or"] = []bson.M{
				{"name": bson.M{"$regex": primitive.Regex{Pattern: opts.Filter.SearchTerm, Options: "i"}}},
				{"description": bson.M{"$regex": primitive.Regex{Pattern: opts.Filter.SearchTerm, Options: "i"}}},
			}
		}
	}

	// Get total count for pagination
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count products", zap.Error(err))
		return nil, 0, err
	}

	// Set up find options
	findOptions := options.Find()

	// Apply pagination if provided
	if opts != nil && opts.Pagination != nil {
		if opts.Pagination.PageSize > 0 {
			findOptions.SetLimit(int64(opts.Pagination.PageSize))
			if opts.Pagination.Page > 1 {
				findOptions.SetSkip(int64((opts.Pagination.Page - 1) * opts.Pagination.PageSize))
			}
		}
	}

	// Apply sorting if provided
	if opts != nil && opts.Sort != nil {
		sortField := ""
		switch opts.Sort.Field {
		case domain.SortFieldName:
			sortField = "name"
		case domain.SortFieldPrice:
			sortField = "price"
		case domain.SortFieldCreatedAt:
			sortField = "created_at"
		case domain.SortFieldUpdatedAt:
			sortField = "updated_at"
		default:
			sortField = "created_at"
		}

		sortOrder := 1 // Default to ascending
		if opts.Sort.Order == domain.SortOrderDesc {
			sortOrder = -1
		}

		findOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})
	} else {
		// Default sorting by created_at desc
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("Failed to list products", zap.Error(err))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		r.logger.Error("Failed to decode products", zap.Error(err))
		return nil, 0, err
	}

	if products == nil {
		products = []*domain.Product{}
	}

	return products, total, nil
}
