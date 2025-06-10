package domain

import "context"

// PermissionRepository defines the interface for permission data operations
type PermissionRepository interface {
	// GrantPermission grants a permission to a user for a specific resource
	GrantPermission(ctx context.Context, userID string, resourceType ResourceType, resourceID string, permission Permission) error
	
	// RevokePermission revokes a permission from a user for a specific resource
	RevokePermission(ctx context.Context, userID string, resourceType ResourceType, resourceID string) error
	
	// GetUserPermissions returns all permissions for a user
	GetUserPermissions(ctx context.Context, userID string, resourceType *ResourceType) ([]*UserPermission, error)
	
	// GetResourcePermissions returns all users who have permissions on a resource
	GetResourcePermissions(ctx context.Context, resourceType ResourceType, resourceID string) ([]*UserPermission, error)
	
	// HasPermission checks if a user has a specific permission on a resource
	HasPermission(ctx context.Context, userID string, resourceType ResourceType, resourceID string, permission Permission) (bool, error)
	
	// GetUserResources returns all resources of a specific type that a user has access to
	GetUserResources(ctx context.Context, userID string, resourceType ResourceType) ([]*UserResource, error)
	
	// UpdatePermission updates a user's permission on a resource
	UpdatePermission(ctx context.Context, userID string, resourceType ResourceType, resourceID string, permission Permission) error
}
