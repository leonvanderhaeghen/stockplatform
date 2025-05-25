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

// UserRepository implements the domain.UserRepository interface
type UserRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewUserRepository creates a new MongoDB user repository
func NewUserRepository(db *mongo.Database, collectionName string, logger *zap.Logger) domain.UserRepository {
	collection := db.Collection(collectionName)
	
	// Create indexes for improved query performance
	indexModels := []mongo.IndexModel{
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		mongo.IndexModel{
			Keys:    bson.D{{Key: "role", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		logger.Warn("Failed to create indexes", zap.Error(err))
	}
	
	return &UserRepository{
		collection: collection,
		logger:     logger.Named("user_repository"),
	}
}

// Create adds a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.logger.Debug("Creating user", 
		zap.String("id", user.ID),
		zap.String("email", user.Email),
	)
	
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		r.logger.Error("Failed to create user", 
			zap.Error(err),
			zap.String("email", user.Email),
		)
		return err
	}
	
	return nil
}

// GetByID finds a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.logger.Debug("Getting user by ID", zap.String("id", id))
	
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("User not found", zap.String("id", id))
			return nil, nil
		}
		r.logger.Error("Failed to get user", 
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, err
	}
	
	return &user, nil
}

// GetByEmail finds a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.logger.Debug("Getting user by email", zap.String("email", email))
	
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Debug("User not found", zap.String("email", email))
			return nil, nil
		}
		r.logger.Error("Failed to get user", 
			zap.Error(err),
			zap.String("email", email),
		)
		return nil, err
	}
	
	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	r.logger.Debug("Updating user", 
		zap.String("id", user.ID),
		zap.String("email", user.Email),
	)
	
	user.UpdatedAt = time.Now()
	
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		r.logger.Error("Failed to update user", 
			zap.Error(err),
			zap.String("id", user.ID),
		)
		return err
	}
	
	if result.MatchedCount == 0 {
		r.logger.Warn("No user was updated", zap.String("id", user.ID))
		return errors.New("user not found")
	}
	
	return nil
}

// Delete removes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Deleting user", zap.String("id", id))
	
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		r.logger.Error("Failed to delete user", 
			zap.Error(err),
			zap.String("id", id),
		)
		return err
	}
	
	if result.DeletedCount == 0 {
		r.logger.Warn("No user was deleted", zap.String("id", id))
		return errors.New("user not found")
	}
	
	return nil
}

// List returns all users with optional filtering and pagination
func (r *UserRepository) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*domain.User, error) {
	r.logger.Debug("Listing users", 
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
		r.logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var users []*domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			r.logger.Error("Failed to decode user", zap.Error(err))
			return nil, err
		}
		users = append(users, &user)
	}
	
	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor error while listing users", zap.Error(err))
		return nil, err
	}
	
	return users, nil
}
