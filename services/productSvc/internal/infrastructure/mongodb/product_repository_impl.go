package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"stockplatform/services/productSvc/internal/domain"
)

// ProductRepository is a MongoDB implementation of the ProductRepository interface
type ProductRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewProductRepository creates a new MongoDB product repository
func NewProductRepository(db *mongo.Database, logger *zap.Logger) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
		logger:     logger.With(zap.String("component", "mongodb.ProductRepository")),
	}
}

// Create creates a new product in the database
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	// Set timestamps if not set
	now := time.Now()
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	if product.UpdatedAt.IsZero() {
		product.UpdatedAt = now
	}

	// Insert the product
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, domain.ErrProductAlreadyExists
		}
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		product.ID = oid
	} else {
		return nil, errors.New("failed to get inserted ID")
	}

	return product, nil
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	// Parse the ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	// Find the product
	var product domain.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	if product.ID.IsZero() {
		return domain.ErrInvalidID
	}

	// Set updated timestamp
	product.UpdatedAt = time.Now()

	// Prepare update document
	update := bson.M{
		"$set": bson.M{
			"name":           product.Name,
			"description":    product.Description,
			"cost_price":     product.CostPrice,
			"selling_price":  product.SellingPrice,
			"currency":       product.Currency,
			"sku":            product.SKU,
			"barcode":        product.Barcode,
			"category_ids":   product.CategoryIDs,
			"supplier_id":    product.SupplierID,
			"is_active":      product.IsActive,
			"is_visible":     product.IsVisible,
			"variants":       product.Variants,
			"in_stock":       product.InStock,
			"stock_qty":      product.StockQty,
			"low_stock_at":   product.LowStockAt,
			"image_urls":     product.ImageURLs,
			"video_urls":     product.VideoURLs,
			"metadata":       product.Metadata,
			"updated_at":     product.UpdatedAt,
		},
	}

	// Execute update
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": product.ID, "deleted_at": bson.M{"$exists": false}},
		update,
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrProductAlreadyExists
		}
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// Delete permanently deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	// Parse the ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	// Delete the product
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// SoftDelete marks a product as deleted without removing it from the database
func (r *ProductRepository) SoftDelete(ctx context.Context, id string) error {
	// Parse the ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	// Update the product to set deleted_at
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now(),
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to soft delete product: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// List retrieves a list of products with pagination and filtering
func (r *ProductRepository) List(ctx context.Context, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// Build the filter
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	// Apply search filter if provided
	if opts != nil && opts.Search != "" {
		filter["$text"] = bson.M{"$search": opts.Search}
	}

	// Apply additional filters
	if opts != nil && opts.Filters != nil {
		for key, value := range opts.Filters {
			// Handle special cases like price ranges, categories, etc.
			switch key {
			case "category_id":
				filter["category_ids"] = value
			case "supplier_id":
				filter["supplier_id"] = value
			case "min_price":
				filter["selling_price"] = bson.M{"$gte": value}
			case "max_price":
				if _, exists := filter["selling_price"]; exists {
					filter["selling_price"].(bson.M)["$lte"] = value
				} else {
					filter["selling_price"] = bson.M{"$lte": value}
				}
			case "in_stock":
				filter["in_stock"] = value
			default:
				filter[key] = value
			}
		}
	}

	// Set up find options
	findOptions := options.Find()

	// Apply pagination
	if opts != nil {
		if opts.PageSize > 0 {
			findOptions.SetLimit(int64(opts.PageSize))
			if opts.Page > 0 {
				findOptions.SetSkip(int64((opts.Page - 1) * opts.PageSize))
			}
		}

		// Apply sorting
		if opts.SortBy != "" {
			sortOrder := 1
			if opts.SortOrder == "desc" {
				sortOrder = -1
			}
			findOptions.SetSort(bson.D{{Key: opts.SortBy, Value: sortOrder}})
		}
	}

	// Count total matching documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Find products
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode products
	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, total, nil
}

// Search searches for products by query
func (r *ProductRepository) Search(ctx context.Context, query string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	if query == "" {
		return r.List(ctx, opts)
	}

	// Create a text search filter
	filter := bson.M{
		"$text":    bson.M{"$search": query},
		"deleted_at": bson.M{"$exists": false},
	}

	// Apply additional filters
	if opts != nil && opts.Filters != nil {
		for key, value := range opts.Filters {
			filter[key] = value
		}
	}

	// Set up find options
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: bson.M{"$meta": "textScore"}}})
	findOptions.SetProjection(bson.M{"score": bson.M{"$meta": "textScore"}})

	// Apply pagination
	if opts != nil && opts.PageSize > 0 {
		findOptions.SetLimit(int64(opts.PageSize))
		if opts.Page > 0 {
			findOptions.SetSkip(int64((opts.Page - 1) * opts.PageSize))
		}
	}

	// Count total matching documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Execute search
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search results: %w", err)
	}

	return products, total, nil
}

// UpdateStock updates the stock quantity of a product
func (r *ProductRepository) UpdateStock(ctx context.Context, id string, quantity int32) error {
	// Parse the ID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	// Update the stock quantity
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}},
		bson.M{
			"$set": bson.M{
				"stock_qty":  quantity,
				"in_stock":   quantity > 0,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

// BulkUpdateVisibility updates the visibility of multiple products for a supplier
func (r *ProductRepository) BulkUpdateVisibility(ctx context.Context, supplierID string, productIDs []string, isVisible bool) error {
	if supplierID == "" {
		return domain.ErrInvalidSupplierID
	}

	if len(productIDs) == 0 {
		return domain.ErrNoProductsProvided
	}

	// Convert string IDs to ObjectIDs
	var objIDs []primitive.ObjectID
	for _, id := range productIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return domain.ErrInvalidID
		}
		objIDs = append(objIDs, objID)
	}

	// Create the update operation
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("is_visible.%s", supplierID): isVisible,
			"updated_at": time.Now(),
		},
	}

	// Execute bulk update
	result, err := r.collection.UpdateMany(
		ctx,
		bson.M{
			"_id":        bson.M{"$in": objIDs},
			"deleted_at": bson.M{"$exists": false},
		},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to bulk update visibility: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNoProductsUpdated
	}

	return nil
}

// The following methods are stubs that need to be implemented

func (r *ProductRepository) GetBySupplier(ctx context.Context, supplierID string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// TODO: Implement GetBySupplier
	return nil, 0, nil
}

func (r *ProductRepository) GetByCategory(ctx context.Context, categoryID string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// TODO: Implement GetByCategory
	return nil, 0, nil
}

func (r *ProductRepository) BulkUpdateStock(ctx context.Context, updates map[string]int32) error {
	// TODO: Implement BulkUpdateStock
	return nil
}

func (r *ProductRepository) GetLowStockProducts(ctx context.Context, threshold int32, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	// TODO: Implement GetLowStockProducts
	return nil, 0, nil
}

func (r *ProductRepository) PublishProducts(ctx context.Context, productIDs []string, publish bool) error {
	// TODO: Implement PublishProducts
	return nil
}

func (r *ProductRepository) UpdateVariantStock(ctx context.Context, productID, variantID string, quantity int32) error {
	// TODO: Implement UpdateVariantStock
	return nil
}
