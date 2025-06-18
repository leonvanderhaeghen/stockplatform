package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
)

type permissionRepository struct {
	collection *mongo.Collection
}

// NewPermissionRepository creates a new MongoDB permission repository
func NewPermissionRepository(db *mongo.Database, collectionName string) domain.PermissionRepository {
	return &permissionRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *permissionRepository) GrantPermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string, permission domain.Permission) error {
	// Check if permission already exists
	existing, err := r.collection.FindOne(ctx, bson.M{
		"user_id":       userID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}).DecodeBytes()

	if err == nil && existing != nil {
		// Update existing permission
		_, err = r.collection.UpdateOne(
			ctx,
			bson.M{
				"user_id":       userID,
				"resource_type": resourceType,
				"resource_id":   resourceID,
			},
			bson.M{
				"$set": bson.M{
					"permission": permission,
					"updated_at": time.Now(),
				},
			},
		)
		return err
	}

	// Create new permission
	_, err = r.collection.InsertOne(ctx, bson.M{
		"user_id":       userID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
		"permission":    permission,
		"created_at":    time.Now(),
		"updated_at":    time.Now(),
	})

	return err
}

func (r *permissionRepository) RevokePermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"user_id":       userID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
	})

	return err
}

func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID string, resourceType *domain.ResourceType) ([]*domain.UserPermission, error) {
	filter := bson.M{"user_id": userID}
	if resourceType != nil {
		filter["resource_type"] = *resourceType
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*domain.UserPermission
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *permissionRepository) GetResourcePermissions(ctx context.Context, resourceType domain.ResourceType, resourceID string) ([]*domain.UserPermission, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*domain.UserPermission
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *permissionRepository) HasPermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string, permission domain.Permission) (bool, error) {
	// Admin users have all permissions
	userRepo := NewUserRepository(r.collection.Database(), "users", nil) // TODO: Pass proper logger
	user, err := userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if user.IsAdmin() {
		return true, nil
	}

	// Check specific permission
	var perm domain.UserPermission
	err = r.collection.FindOne(ctx, bson.M{
		"user_id":       userID,
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}).Decode(&perm)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return perm.HasPermission(string(permission)), nil
}

func (r *permissionRepository) GetUserResources(ctx context.Context, userID string, resourceType domain.ResourceType) ([]*domain.UserResource, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":       userID,
		"resource_type": resourceType,
	})

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resources []*domain.UserResource
	if err := cursor.All(ctx, &resources); err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *permissionRepository) UpdatePermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string, permission domain.Permission) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"user_id":       userID,
			"resource_type": resourceType,
			"resource_id":   resourceID,
		},
		bson.M{
			"$set": bson.M{
				"permission": permission,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}
