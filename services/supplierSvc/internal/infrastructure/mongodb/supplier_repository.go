package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

type supplierRepository struct {
	collection *mongo.Collection
}

// NewSupplierRepository creates a new MongoDB supplier repository
func NewSupplierRepository(db *mongo.Database, collectionName string) domain.SupplierRepository {
	return &supplierRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error) {
	supplier.CreatedAt = time.Now()
	supplier.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, supplier)
	if err != nil {
		return nil, err
	}

	supplier.ID = result.InsertedID.(primitive.ObjectID)
	return supplier, nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id string) (*domain.Supplier, error) {
	var supplier domain.Supplier
	
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidInput
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&supplier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &supplier, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	if supplier.ID.IsZero() {
		return domain.ErrInvalidInput
	}

	objectID := supplier.ID

	supplier.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": objectID},
		supplier,
	)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ErrNotFound
		}
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidInput
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *supplierRepository) List(ctx context.Context, page, pageSize int32, search string) ([]*domain.Supplier, int32, error) {
	// Set up pagination
	skip := (page - 1) * pageSize
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(pageSize))

	// Build the filter
	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
			{"email": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
			{"contact_person": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
		}
	}

	// Count total matching documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find documents
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var suppliers []*domain.Supplier
	if err := cursor.All(ctx, &suppliers); err != nil {
		return nil, 0, err
	}

	return suppliers, int32(total), nil
}
