package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"stockplatform/services/userSvc/internal/domain"
)

// AddressRepository implements the domain.AddressRepository interface
type AddressRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewAddressRepository creates a new MongoDB address repository
func NewAddressRepository(db *mongo.Database, collectionName string, logger *zap.Logger) domain.AddressRepository {
	collection := db.Collection(collectionName)
	
	// Create indexes for improved query performance
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.D{{Key: "is_default", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		logger.Warn("Failed to create indexes", zap.Error(err))
	}
	
	return &AddressRepository{
		collection: collection,
		logger:     logger.Named("address_repository"),
	}
}

// Create adds a new address
func (r *AddressRepository) Create(ctx context.Context, address *domain.Address) error {
	r.logger.Debug("Creating address", 
		zap.String("id", address.ID),
		zap.String("user_id", address.UserID),
	)
	
	// If this is the default address, ensure all other addresses are not default
	if address.IsDefault {
		if err := r.clearDefaultAddresses(ctx, address.UserID); err != nil {
			r.logger.Error("Failed to clear default addresses", 
				zap.Error(err),
				zap.String("user_id", address.UserID),
			)
			return err
		}
	}
	
	_, err := r.collection.InsertOne(ctx, address)
	if err != nil {
		r.logger.Error("Failed to create address", 
			zap.Error(err),
			zap.String("id", address.ID),
			zap.String("user_id", address.UserID),
		)
		return err
	}
	
	return nil
}

// GetByID finds an address by ID
func (r *AddressRepository) GetByID(ctx context.Context, id string) (*domain.Address, error) {
	r.logger.Debug("Getting address by ID", zap.String("id", id))
	
	var address domain.Address
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&address)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Address not found", zap.String("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get address", 
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, err
	}
	
	return &address, nil
}

// GetByUserID finds addresses for a user
func (r *AddressRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Address, error) {
	r.logger.Debug("Getting addresses by user ID", zap.String("user_id", userID))
	
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "is_default", Value: -1}, {Key: "created_at", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, findOptions)
	if err != nil {
		r.logger.Error("Failed to get addresses by user ID", 
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var addresses []*domain.Address
	for cursor.Next(ctx) {
		var address domain.Address
		if err := cursor.Decode(&address); err != nil {
			r.logger.Error("Failed to decode address", zap.Error(err))
			return nil, err
		}
		addresses = append(addresses, &address)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while getting addresses by user ID", zap.Error(err))
		return nil, err
	}
	
	return addresses, nil
}

// GetDefaultByUserID finds the default address for a user
func (r *AddressRepository) GetDefaultByUserID(ctx context.Context, userID string) (*domain.Address, error) {
	r.logger.Debug("Getting default address by user ID", zap.String("user_id", userID))
	
	var address domain.Address
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID, "is_default": true}).Decode(&address)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Default address not found", zap.String("user_id", userID))
			return nil, nil
		}
		r.logger.Error("Failed to get default address", 
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return nil, err
	}
	
	return &address, nil
}

// Update updates an existing address
func (r *AddressRepository) Update(ctx context.Context, address *domain.Address) error {
	r.logger.Debug("Updating address", 
		zap.String("id", address.ID),
		zap.String("user_id", address.UserID),
	)
	
	// If this is being set as the default address, ensure all other addresses are not default
	if address.IsDefault {
		if err := r.clearDefaultAddresses(ctx, address.UserID); err != nil {
			r.logger.Error("Failed to clear default addresses", 
				zap.Error(err),
				zap.String("user_id", address.UserID),
			)
			return err
		}
	}
	
	address.UpdatedAt = time.Now()
	
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": address.ID}, address)
	if err != nil {
		r.logger.Error("Failed to update address", 
			zap.Error(err),
			zap.String("id", address.ID),
		)
		return err
	}
	
	if result.MatchedCount == 0 {
		r.logger.Warn("No address was updated", zap.String("id", address.ID))
		return errors.New("address not found")
	}
	
	return nil
}

// Delete removes an address
func (r *AddressRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting address", zap.String("id", id))
	
	// First get the address to check if it's the default
	var address domain.Address
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&address)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("Address not found", zap.String("id", id))
			return errors.New("address not found")
		}
		r.logger.Error("Failed to get address before deletion", 
			zap.Error(err),
			zap.String("id", id),
		)
		return err
	}
	
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete address", 
			zap.Error(err),
			zap.String("id", id),
		)
		return err
	}
	
	if result.DeletedCount == 0 {
		r.logger.Warn("No address was deleted", zap.String("id", id))
		return errors.New("address not found")
	}
	
	// If this was the default address, try to set a new default
	if address.IsDefault {
		r.logger.Debug("Default address deleted, attempting to set a new default",
			zap.String("user_id", address.UserID),
		)
		
		// Find another address for this user
		var newDefault domain.Address
		err := r.collection.FindOne(ctx, bson.M{"user_id": address.UserID}).Decode(&newDefault)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				// No other addresses, nothing to do
				r.logger.Debug("No other addresses to set as default", 
					zap.String("user_id", address.UserID),
				)
				return nil
			}
			r.logger.Error("Failed to find another address to set as default", 
				zap.Error(err),
				zap.String("user_id", address.UserID),
			)
			return nil // Don't fail the deletion just because we couldn't set a new default
		}
		
		// Set this address as the new default
		newDefault.IsDefault = true
		newDefault.UpdatedAt = time.Now()
		
		_, err = r.collection.UpdateOne(
			ctx,
			bson.M{"_id": newDefault.ID},
			bson.M{"$set": bson.M{"is_default": true, "updated_at": newDefault.UpdatedAt}},
		)
		if err != nil {
			r.logger.Error("Failed to set new default address", 
				zap.Error(err),
				zap.String("id", newDefault.ID),
				zap.String("user_id", newDefault.UserID),
			)
			// Don't fail the deletion just because we couldn't set a new default
		}
	}
	
	return nil
}

// SetDefaultAddress sets an address as the default for a user
func (r *AddressRepository) SetDefaultAddress(ctx context.Context, userID string, addressID string) error {
	r.logger.Debug("Setting default address", 
		zap.String("user_id", userID),
		zap.String("address_id", addressID),
	)
	
	// First clear all default addresses for this user
	if err := r.clearDefaultAddresses(ctx, userID); err != nil {
		return err
	}
	
	// Now set the specified address as default
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": addressID, "user_id": userID},
		bson.M{"$set": bson.M{"is_default": true, "updated_at": time.Now()}},
	)
	if err != nil {
		r.logger.Error("Failed to set default address", 
			zap.Error(err),
			zap.String("user_id", userID),
			zap.String("address_id", addressID),
		)
		return err
	}
	
	if result.MatchedCount == 0 {
		r.logger.Warn("No address was set as default", 
			zap.String("user_id", userID),
			zap.String("address_id", addressID),
		)
		return errors.New("address not found")
	}
	
	return nil
}

// clearDefaultAddresses clears the default flag from all addresses for a user
func (r *AddressRepository) clearDefaultAddresses(ctx context.Context, userID string) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID, "is_default": true},
		bson.M{"$set": bson.M{"is_default": false, "updated_at": time.Now()}},
	)
	if err != nil {
		r.logger.Error("Failed to clear default addresses", 
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return err
	}
	
	return nil
}
