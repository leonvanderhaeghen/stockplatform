package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
)

// SupplierRepository is a MongoDB implementation of the SupplierRepository interface
type SupplierRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewSupplierRepository creates a new SupplierRepository
func NewSupplierRepository(db *mongo.Database, logger *zap.Logger) *SupplierRepository {
	return &SupplierRepository{
		collection: db.Collection("suppliers"),
		logger:     logger.Named("supplier_repository"),
	}
}

// Create creates a new supplier
func (r *SupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error) {
	now := time.Now()
	supplier.CreatedAt = now
	supplier.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, supplier)
	if err != nil {
		r.logger.Error("Failed to create supplier", zap.Error(err))
		return nil, err
	}

	supplier.ID = result.InsertedID.(primitive.ObjectID)
	return supplier, nil
}

// GetByID retrieves a supplier by ID
func (r *SupplierRepository) GetByID(ctx context.Context, id string) (*domain.Supplier, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	var supplier domain.Supplier
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&supplier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrSupplierNotFound
		}
		r.logger.Error("Failed to get supplier", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return &supplier, nil
}

// Update updates an existing supplier
func (r *SupplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	supplier.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":            supplier.Name,
			"contact_person":  supplier.ContactPerson,
			"email":           supplier.Email,
			"phone":           supplier.Phone,
			"address":         supplier.Address,
			"city":            supplier.City,
			"state":           supplier.State,
			"country":         supplier.Country,
			"postal_code":     supplier.PostalCode,
			"tax_id":          supplier.TaxID,
			"website":         supplier.Website,
			"currency":        supplier.Currency,
			"lead_time_days":  supplier.LeadTimeDays,
			"payment_terms":   supplier.PaymentTerms,
			"metadata":        supplier.Metadata,
			"updated_at":      supplier.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateByID(ctx, supplier.ID, update)
	if err != nil {
		r.logger.Error("Failed to update supplier", zap.String("id", supplier.ID.Hex()), zap.Error(err))
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrSupplierNotFound
	}

	return nil
}

// Delete deletes a supplier by ID
func (r *SupplierRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("Failed to delete supplier", zap.String("id", id), zap.Error(err))
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrSupplierNotFound
	}

	return nil
}

// List retrieves a list of suppliers with pagination and optional filtering
func (r *SupplierRepository) List(ctx context.Context, page, pageSize int32, search string) ([]*domain.Supplier, int32, error) {
	// Set up pagination
	skip := (page - 1) * pageSize
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(pageSize))

	// Build filter
	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
			{"email": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
			{"contact_person": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
		}
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Error("Failed to count suppliers", zap.Error(err))
		return nil, 0, err
	}

	// Find suppliers
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to list suppliers", zap.Error(err))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var suppliers []*domain.Supplier
	if err := cursor.All(ctx, &suppliers); err != nil {
		r.logger.Error("Failed to decode suppliers", zap.Error(err))
		return nil, 0, err
	}

	return suppliers, int32(total), nil
}
