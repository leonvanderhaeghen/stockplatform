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

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
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
	// Build the base filter
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}

	// Apply filters from ListOptions if provided
	if opts != nil && opts.Filter != nil {
		// Apply category filter
		if len(opts.Filter.CategoryIDs) > 0 {
			filter["category_ids"] = bson.M{"$in": opts.Filter.CategoryIDs}
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
			filter["selling_price"] = priceFilter
		}

		// Apply search term
		if opts.Filter.SearchTerm != "" {
			filter["$text"] = bson.M{"$search": opts.Filter.SearchTerm}
		}

		// Apply IDs filter
		if len(opts.Filter.IDs) > 0 {
			var objectIDs []primitive.ObjectID
			for _, id := range opts.Filter.IDs {
				objID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return nil, 0, fmt.Errorf("invalid product ID: %v", id)
				}
				objectIDs = append(objectIDs, objID)
			}
			filter["_id"] = bson.M{"$in": objectIDs}
		}
	}

	// Set up find options
	findOptions := options.Find()

	// Apply pagination if provided
	if opts != nil && opts.Pagination != nil {
		if opts.Pagination.PageSize > 0 {
			findOptions.SetLimit(int64(opts.Pagination.PageSize))
			if opts.Pagination.Page > 0 {
				findOptions.SetSkip(int64((opts.Pagination.Page - 1) * opts.Pagination.PageSize))
			}
		}

		// Apply sorting if provided
		if opts.Sort != nil {
			sortField := "created_at" // Default sort field
			switch opts.Sort.Field {
			case domain.SortFieldName:
				sortField = "name"
			case domain.SortFieldPrice:
				sortField = "selling_price"
			case domain.SortFieldCreatedAt:
				sortField = "created_at"
			case domain.SortFieldUpdatedAt:
				sortField = "updated_at"
			}

			sortOrder := 1 // Default to ascending
			if opts.Sort.Order == domain.SortOrderDesc {
				sortOrder = -1
			}

			findOptions.SetSort(bson.D{{Key: sortField, Value: sortOrder}})
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

	// Create a copy of the options to avoid modifying the original
	searchOpts := &domain.ListOptions{}
	if opts != nil {
		// Copy the filter if it exists
		if opts.Filter != nil {
			searchOpts.Filter = &domain.ProductFilter{
				IDs:         opts.Filter.IDs,
				CategoryIDs: opts.Filter.CategoryIDs,
				MinPrice:    opts.Filter.MinPrice,
				MaxPrice:    opts.Filter.MaxPrice,
				SearchTerm:  query, // Override the search term with the query parameter
			}
		} else {
			searchOpts.Filter = &domain.ProductFilter{
				SearchTerm: query,
			}
		}

		// Copy pagination and sort options
		if opts.Pagination != nil {
			searchOpts.Pagination = &domain.Pagination{
				Page:     opts.Pagination.Page,
				PageSize: opts.Pagination.PageSize,
			}
		}

		if opts.Sort != nil {
			searchOpts.Sort = &domain.SortOption{
				Field: opts.Sort.Field,
				Order: opts.Sort.Order,
			}
		}
	} else {
		// If no options provided, just set the search term
		searchOpts.Filter = &domain.ProductFilter{
			SearchTerm: query,
		}
	}

	// Use the List method with the updated options
	return r.List(ctx, searchOpts)
}

// UpdateStock is deprecated - inventory operations are handled by inventorySvc
func (r *ProductRepository) UpdateStock(ctx context.Context, id string, quantity int32) error {
	return fmt.Errorf("inventory operations are handled by inventorySvc")
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
	r.logger.Debug("Getting products by supplier",
		zap.String("supplier_id", supplierID),
		zap.Any("options", opts))

	// Build filter for supplier ID
	filter := bson.M{"supplier_id": supplierID}

	// Apply additional filters from options if provided
	if opts != nil && opts.Filter != nil {
		// Add search filter if provided
		if opts.Filter.SearchTerm != "" {
			filter["$or"] = []bson.M{
				{"name": bson.M{"$regex": opts.Filter.SearchTerm, "$options": "i"}},
				{"description": bson.M{"$regex": opts.Filter.SearchTerm, "$options": "i"}},
				{"sku": bson.M{"$regex": opts.Filter.SearchTerm, "$options": "i"}},
			}
		}

		// Add price range filters if specified
		if opts.Filter.MinPrice > 0 {
			filter["price"] = bson.M{"$gte": opts.Filter.MinPrice}
		}
		if opts.Filter.MaxPrice > 0 {
			if priceFilter, exists := filter["price"]; exists {
				priceFilter.(bson.M)["$lte"] = opts.Filter.MaxPrice
			} else {
				filter["price"] = bson.M{"$lte": opts.Filter.MaxPrice}
			}
		}
	}

	// Get total count
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products by supplier: %w", err)
	}

	// Build find options with pagination
	findOptions := options.Find()
	if opts != nil && opts.Pagination != nil {
		if opts.Pagination.PageSize > 0 {
			findOptions.SetLimit(int64(opts.Pagination.PageSize))
			if opts.Pagination.Page > 1 {
				skip := (opts.Pagination.Page - 1) * opts.Pagination.PageSize
				findOptions.SetSkip(int64(skip))
			}
		}
	}

	// Apply sorting
	if opts != nil && opts.Sort != nil {
		var sortField string
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
			sortField = "name" // default to name
		}

		sortOrder := 1
		if opts.Sort.Order == domain.SortOrderDesc {
			sortOrder = -1
		}
		findOptions.SetSort(bson.D{{sortField, sortOrder}})
	} else {
		// Default sorting by name ascending
		findOptions.SetSort(bson.D{{"name", 1}})
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products by supplier: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var products []*domain.Product
	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, 0, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	r.logger.Debug("Successfully retrieved products by supplier",
		zap.String("supplier_id", supplierID),
		zap.Int("count", len(products)),
		zap.Int64("total", totalCount))

	return products, totalCount, nil
}

func (r *ProductRepository) GetByCategory(ctx context.Context, categoryID string, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	r.logger.Debug("Getting products by category",
		zap.String("category_id", categoryID),
		zap.Any("options", opts))

	// Build filter for category ID
	filter := bson.M{"category_id": categoryID}

	// Apply additional filters from options if provided
	if opts != nil && opts.Filter != nil {
		// Add search filter if provided
		if opts.Filter.SearchTerm != "" {
			filter["$or"] = []bson.M{
				{"name": bson.M{"$regex": opts.Filter.SearchTerm, "$options": "i"}},
				{"description": bson.M{"$regex": opts.Filter.SearchTerm, "$options": "i"}},
				{"sku": bson.M{"$regex": opts.Filter.SearchTerm, "$options": "i"}},
			}
		}

		// Add price range filters if specified
		if opts.Filter.MinPrice > 0 {
			filter["price"] = bson.M{"$gte": opts.Filter.MinPrice}
		}
		if opts.Filter.MaxPrice > 0 {
			if priceFilter, exists := filter["price"]; exists {
				priceFilter.(bson.M)["$lte"] = opts.Filter.MaxPrice
			} else {
				filter["price"] = bson.M{"$lte": opts.Filter.MaxPrice}
			}
		}
	}

	// Get total count
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products by category: %w", err)
	}

	// Build find options with pagination
	findOptions := options.Find()
	if opts != nil && opts.Pagination != nil {
		if opts.Pagination.PageSize > 0 {
			findOptions.SetLimit(int64(opts.Pagination.PageSize))
			if opts.Pagination.Page > 1 {
				skip := (opts.Pagination.Page - 1) * opts.Pagination.PageSize
				findOptions.SetSkip(int64(skip))
			}
		}
	}

	// Apply sorting
	if opts != nil && opts.Sort != nil {
		var sortField string
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
			sortField = "name" // default to name
		}

		sortOrder := 1
		if opts.Sort.Order == domain.SortOrderDesc {
			sortOrder = -1
		}
		findOptions.SetSort(bson.D{{sortField, sortOrder}})
	} else {
		// Default sorting by name ascending
		findOptions.SetSort(bson.D{{"name", 1}})
	}

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products by category: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results
	var products []*domain.Product
	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, 0, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, &product)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("cursor error: %w", err)
	}

	r.logger.Debug("Successfully retrieved products by category",
		zap.String("category_id", categoryID),
		zap.Int("count", len(products)),
		zap.Int64("total", totalCount))

	return products, totalCount, nil
}

// BulkUpdateStock is deprecated - inventory operations are handled by inventorySvc
func (r *ProductRepository) BulkUpdateStock(ctx context.Context, updates map[string]int32) error {
	return fmt.Errorf("inventory operations are handled by inventorySvc")
}

// GetLowStockProducts is deprecated - inventory operations are handled by inventorySvc
func (r *ProductRepository) GetLowStockProducts(ctx context.Context, threshold int32, opts *domain.ListOptions) ([]*domain.Product, int64, error) {
	return nil, 0, fmt.Errorf("inventory operations are handled by inventorySvc")
}

// PublishProducts updates the published status of multiple products
func (r *ProductRepository) PublishProducts(ctx context.Context, productIDs []string, publish bool) error {
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
			"is_published": publish,
			"updated_at":  time.Now(),
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
		return fmt.Errorf("failed to update product publish status: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.ErrNoProductsUpdated
	}

	return nil
}

// UpdateVariantStock updates the stock quantity of a product variant
func (r *ProductRepository) UpdateVariantStock(ctx context.Context, productID, variantID string, quantity int32) error {
	// Validate input
	if productID == "" {
		return domain.ErrInvalidID
	}

	if variantID == "" {
		return errors.New("variant ID is required")
	}

	if quantity < 0 {
		return domain.ErrInvalidQuantity
	}

	// Convert product ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return domain.ErrInvalidID
	}

	// First, find the product to get the current variant
	var product domain.Product
	err = r.collection.FindOne(
		ctx,
		bson.M{
			"_id":        objID,
			"deleted_at": bson.M{"$exists": false},
		},
	).Decode(&product)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("failed to find product: %w", err)
	}

	// Find the variant and update its stock
	variantFound := false
	for _, variant := range product.Variants {
		if variant.ID == variantID {
			// Update the variant's stock in the product
			// Note: Since we can't directly update the variant in the array,
			// we'll use an arrayFilters update
			update := bson.M{
				"$set": bson.M{
					"variants.$[elem].stock_quantity": quantity,
					"variants.$[elem].in_stock":       quantity > 0,
					"variants.$[elem].updated_at":     time.Now(),
					"updated_at":                      time.Now(),
				},
			}

			// Execute the update with arrayFilters
			result, err := r.collection.UpdateOne(
				ctx,
				bson.M{
					"_id":        objID,
					"deleted_at": bson.M{"$exists": false},
				},
				update,
				options.Update().SetArrayFilters(options.ArrayFilters{
					Filters: []interface{}{bson.M{"elem.id": variantID}},
				}),
			)

			if err != nil {
				return fmt.Errorf("failed to update variant stock: %w", err)
			}

			if result.MatchedCount == 0 {
				return domain.ErrProductNotFound
			}

			variantFound = true
			break
		}
	}

	if !variantFound {
		return domain.ErrVariantNotFound
	}

	return nil
}
