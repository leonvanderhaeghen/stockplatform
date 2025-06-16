package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// LocationDocument is the MongoDB document for store locations
type LocationDocument struct {
	ID         string    `bson:"_id,omitempty"`
	Name       string    `bson:"name"`
	Type       string    `bson:"type"`
	Address    string    `bson:"address"`
	City       string    `bson:"city"`
	State      string    `bson:"state"`
	PostalCode string    `bson:"postal_code"`
	Country    string    `bson:"country"`
	IsActive   bool      `bson:"is_active"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

// LocationRepository implements the domain repository interface for MongoDB
type LocationRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewLocationRepository creates a new MongoDB location repository
func NewLocationRepository(db *mongo.Database, collectionName string, logger *zap.Logger) domain.LocationRepository {
	return &LocationRepository{
		collection: db.Collection(collectionName),
		logger:     logger.Named("location_repository"),
	}
}

// mapToStoreLocation converts a MongoDB document to a domain model
func (r *LocationRepository) mapToStoreLocation(doc *LocationDocument) *domain.StoreLocation {
	return &domain.StoreLocation{
		ID:         doc.ID,
		Name:       doc.Name,
		Type:       doc.Type,
		Address:    doc.Address,
		City:       doc.City,
		State:      doc.State,
		PostalCode: doc.PostalCode,
		Country:    doc.Country,
		IsActive:   doc.IsActive,
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
	}
}

// mapToDocument converts a domain model to a MongoDB document
func (r *LocationRepository) mapToDocument(location *domain.StoreLocation) *LocationDocument {
	return &LocationDocument{
		ID:         location.ID,
		Name:       location.Name,
		Type:       location.Type,
		Address:    location.Address,
		City:       location.City,
		State:      location.State,
		PostalCode: location.PostalCode,
		Country:    location.Country,
		IsActive:   location.IsActive,
		CreatedAt:  location.CreatedAt,
		UpdatedAt:  location.UpdatedAt,
	}
}

// Create creates a new store location
func (r *LocationRepository) Create(ctx context.Context, location *domain.StoreLocation) error {
	if location.ID == "" {
		location.ID = primitive.NewObjectID().Hex()
	}
	
	now := time.Now()
	if location.CreatedAt.IsZero() {
		location.CreatedAt = now
	}
	location.UpdatedAt = now
	
	doc := r.mapToDocument(location)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		r.logger.Error("Failed to create location", 
			zap.Error(err),
			zap.String("location_id", location.ID),
		)
		return err
	}
	
	return nil
}

// GetByID gets a store location by ID
func (r *LocationRepository) GetByID(ctx context.Context, id string) (*domain.StoreLocation, error) {
	var doc LocationDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrLocationNotFound
		}
		r.logger.Error("Failed to get location", 
			zap.Error(err),
			zap.String("location_id", id),
		)
		return nil, err
	}
	
	return r.mapToStoreLocation(&doc), nil
}

// Update updates a store location
func (r *LocationRepository) Update(ctx context.Context, location *domain.StoreLocation) error {
	location.UpdatedAt = time.Now()
	doc := r.mapToDocument(location)
	
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": location.ID}, doc)
	if err != nil {
		r.logger.Error("Failed to update location", 
			zap.Error(err),
			zap.String("location_id", location.ID),
		)
		return err
	}
	
	if result.MatchedCount == 0 {
		return domain.ErrLocationNotFound
	}
	
	return nil
}

// Delete deletes a store location
func (r *LocationRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete location", 
			zap.Error(err),
			zap.String("location_id", id),
		)
		return err
	}
	
	if result.DeletedCount == 0 {
		return domain.ErrLocationNotFound
	}
	
	return nil
}

// GetByName gets a store location by name
func (r *LocationRepository) GetByName(ctx context.Context, name string) (*domain.StoreLocation, error) {
	var doc LocationDocument
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrLocationNotFound
		}
		r.logger.Error("Failed to get location by name", 
			zap.Error(err),
			zap.String("name", name),
		)
		return nil, err
	}
	
	return r.mapToStoreLocation(&doc), nil
}

// ListByType lists store locations of a specific type
func (r *LocationRepository) ListByType(ctx context.Context, locationType string, limit int, offset int) ([]*domain.StoreLocation, error) {
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	
	filter := bson.M{"type": locationType}
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to list locations by type", 
			zap.Error(err),
			zap.String("type", locationType),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*domain.StoreLocation
	for cursor.Next(ctx) {
		var doc LocationDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode location document", zap.Error(err))
			return nil, err
		}
		results = append(results, r.mapToStoreLocation(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing locations by type", zap.Error(err))
		return nil, err
	}
	
	return results, nil
}

// List lists all store locations with pagination
func (r *LocationRepository) List(ctx context.Context, limit int, offset int, includeInactive bool) ([]*domain.StoreLocation, error) {
	opts := options.Find()
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	
	filter := bson.M{}
	if !includeInactive {
		filter["is_active"] = true
	}
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Error("Failed to list locations", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var results []*domain.StoreLocation
	for cursor.Next(ctx) {
		var doc LocationDocument
		if err := cursor.Decode(&doc); err != nil {
			r.logger.Error("Failed to decode location document", zap.Error(err))
			return nil, err
		}
		results = append(results, r.mapToStoreLocation(&doc))
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing locations", zap.Error(err))
		return nil, err
	}
	
	return results, nil
}
