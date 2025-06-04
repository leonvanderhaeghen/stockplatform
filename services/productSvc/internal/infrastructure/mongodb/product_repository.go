package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
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
	err = r.collection.FindOne(ctx, bson.M{
		"_id":        objID,
		"deleted_at": nil,
	}).Decode(&product)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrProductNotFound
		}
		r.logger.Error("Failed to get product by ID", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Ensure empty slices are not nil for the API
	if product.CategoryIDs == nil {
		product.CategoryIDs = []string{}
	}
	if product.ImageURLs == nil {
		product.ImageURLs = []string{}
	}
	if product.VideoURLs == nil {
		product.VideoURLs = []string{}
	}
	if product.Variants == nil {
		product.Variants = []domain.Variant{}
	}

	return &product, nil
}

func (r *productRepository) GetBySupplier(ctx context.Context, supplierID string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	if opts == nil {
		opts = &domain.ListOptions{}
	}

	filter := bson.M{
		"supplier_id": supplierID,
		"deleted_at":  nil,
	}

	// Add search filter if query is provided
	if opts.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": opts.Search, "$options": "i"}},
			{"sku": bson.M{"$regex": opts.Search, "$options": "i"}},
			{"barcode": bson.M{"$regex": opts.Search, "$options": "i"}},
		}
	}

	return r.findProducts(ctx, filter, opts)
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	if product.ID == primitive.NilObjectID {
		return domain.ErrInvalidID
	}

	// Ensure timestamps are set
	now := time.Now()
	product.UpdatedAt = now

	// If this is a new variant, set created_at
	for i := range product.Variants {
		if product.Variants[i].CreatedAt.IsZero() {
			product.Variants[i].CreatedAt = now
		}
		product.Variants[i].UpdatedAt = now

		// Set timestamps for variant options
		for j := range product.Variants[i].Options {
			if product.Variants[i].Options[j].CreatedAt.IsZero() {
				product.Variants[i].Options[j].CreatedAt = now
			}
			product.Variants[i].Options[j].UpdatedAt = now
		}
	}

	// Update stock status based on quantity
	product.InStock = product.StockQty > 0

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{
			"_id":        product.ID,
			"deleted_at": nil,
		},
		product,
	)

	if err != nil {
		r.logger.Error("Failed to update product", 
			zap.String("id", product.ID.Hex()), 
			zap.Error(err))
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
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
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) findProducts(ctx context.Context, filter bson.M, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// Build find options
	findOptions := options.Find()

	// Apply pagination
	if opts.Page > 0 && opts.PageSize > 0 {
		findOptions.SetSkip(int64((opts.Page - 1) * opts.PageSize))
		findOptions.SetLimit(int64(opts.PageSize))
	}

	// Apply sorting
	sortField := "created_at"
	sortOrder := -1 // Default to descending

	if opts.SortBy != "" {
		sortField = strings.ToLower(opts.SortBy)
	}

	if opts.SortOrder == "asc" {
		sortOrder = 1
	}

	findOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	// Get total count for pagination
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Error("Failed to list products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		r.logger.Error("Failed to decode products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to decode products: %w", err)
	}

	// Ensure we never return nil slices
	for _, p := range products {
		if p.CategoryIDs == nil {
			p.CategoryIDs = []string{}
		}
		if p.ImageURLs == nil {
			p.ImageURLs = []string{}
		}
		if p.VideoURLs == nil {
			p.VideoURLs = []string{}
		}
		if p.Variants == nil {
			p.Variants = []domain.Variant{}
		}
	}

	return products, total, nil
}

// List retrieves a paginated list of products with optional filtering
func (r *productRepository) List(
	ctx context.Context,
	opts *domain.ListOptions,
) ([]*domain.Product, int64, error) {
	if opts == nil {
		opts = &domain.ListOptions{}
	}

	// Build base filter
	filter := bson.M{"deleted_at": nil}

	// Apply search if provided
	if opts.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": opts.Search, "$options": "i"}},
			{"sku": bson.M{"$regex": opts.Search, "$options": "i"}},
			{"barcode": bson.M{"$regex": opts.Search, "$options": "i"}},
			{"description": bson.M{"$regex": opts.Search, "$options": "i"}},
		}
	}

	// Apply additional filters if provided
	for k, v := range opts.Filters {
		// Handle special filter types
		switch k {
		case "category_id":
			filter["category_ids"] = v
		case "supplier_id":
			// If filtering by supplier, include both direct supplier and visible products
			if supplierID, ok := v.(string); ok && supplierID != "" {
				filter["$or"] = []bson.M{
					{"supplier_id": supplierID},
					{fmt.Sprintf("is_visible.%s", supplierID): true},
				}
			}
		case "in_stock":
			filter["in_stock"] = v
		case "min_price":
			if minPrice, ok := v.(float64); ok {
				filter["selling_price"] = bson.M{"$gte": fmt.Sprintf("%.2f", minPrice)}
			}
		case "max_price":
			if maxPrice, ok := v.(float64); ok {
				if _, exists := filter["selling_price"]; exists {
					filter["selling_price"].(bson.M)["$lte"] = fmt.Sprintf("%.2f", maxPrice)
				} else {
					filter["selling_price"] = bson.M{"$lte": fmt.Sprintf("%.2f", maxPrice)}
				}
			}
		default:
			filter[k] = v
		}
	}

	return r.findProducts(ctx, filter, opts)
}
