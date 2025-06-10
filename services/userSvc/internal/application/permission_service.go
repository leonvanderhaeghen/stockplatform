package application

import (
	"context"

	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
)

// PermissionService handles permission-related business logic
type PermissionService struct {
	repo domain.PermissionRepository
}

// NewPermissionService creates a new permission service
func NewPermissionService(repo domain.PermissionRepository) *PermissionService {
	return &PermissionService{
		repo: repo,
	}
}

// GrantPermission grants a permission to a user for a specific resource
func (s *PermissionService) GrantPermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string, permission domain.Permission) error {
	return s.repo.GrantPermission(ctx, userID, resourceType, resourceID, permission)
}

// RevokePermission revokes a permission from a user for a specific resource
func (s *PermissionService) RevokePermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string) error {
	return s.repo.RevokePermission(ctx, userID, resourceType, resourceID)
}

// GetUserPermissions returns all permissions for a user
func (s *PermissionService) GetUserPermissions(ctx context.Context, userID string, resourceType *domain.ResourceType) ([]*domain.UserPermission, error) {
	return s.repo.GetUserPermissions(ctx, userID, resourceType)
}

// GetResourcePermissions returns all users who have permissions on a resource
func (s *PermissionService) GetResourcePermissions(ctx context.Context, resourceType domain.ResourceType, resourceID string) ([]*domain.UserPermission, error) {
	return s.repo.GetResourcePermissions(ctx, resourceType, resourceID)
}

// HasPermission checks if a user has a specific permission on a resource
func (s *PermissionService) HasPermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string, permission domain.Permission) (bool, error) {
	return s.repo.HasPermission(ctx, userID, resourceType, resourceID, permission)
}

// GetUserResources returns all resources of a specific type that a user has access to
func (s *PermissionService) GetUserResources(ctx context.Context, userID string, resourceType domain.ResourceType) ([]*domain.UserResource, error) {
	return s.repo.GetUserResources(ctx, userID, resourceType)
}

// UpdatePermission updates a user's permission on a resource
func (s *PermissionService) UpdatePermission(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string, permission domain.Permission) error {
	return s.repo.UpdatePermission(ctx, userID, resourceType, resourceID, permission)
}

// CheckUserAccess checks if a user has the required permission on a resource
func (s *PermissionService) CheckUserAccess(ctx context.Context, userID string, resourceType domain.ResourceType, resourceID string, requiredPermission domain.Permission) (bool, error) {
	return s.repo.HasPermission(ctx, userID, resourceType, resourceID, requiredPermission)
}

// GetUserResourceAccess returns a list of resources the user has access to with their permission level
func (s *PermissionService) GetUserResourceAccess(ctx context.Context, userID string, resourceType domain.ResourceType) ([]*domain.ResourceAccess, error) {
	resources, err := s.repo.GetUserResources(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.ResourceAccess, 0, len(resources))
	for _, r := range resources {
		result = append(result, &domain.ResourceAccess{
			ResourceID: r.ResourceID,
			Permission: r.Permission,
		})
	}

	return result, nil
}
