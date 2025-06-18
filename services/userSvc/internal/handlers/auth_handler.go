package handlers

import (
	"context"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
)

// AuthHandler implements the AuthService gRPC methods
type AuthHandler struct {
	userv1.UnimplementedAuthServiceServer
	authService domain.AuthService
	logger      *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService domain.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.Named("auth_handler"),
	}
}

// ValidateToken validates a JWT token and returns user information
func (h *AuthHandler) ValidateToken(ctx context.Context, req *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
	h.logger.Debug("ValidateToken request received")

	if req.Token == "" {
		return &userv1.ValidateTokenResponse{
			Valid: false,
			Error: "token is required",
		}, nil
	}

	// Validate token using AuthService
	claims, err := h.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		h.logger.Warn("Token validation failed", zap.Error(err))
		return &userv1.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	// Convert claims to User proto message
	user := &userv1.User{
		Id:        claims.UserID,
		Email:     claims.Email,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Role:      mapRoleToProto(domain.Role(claims.Role)),
		Active:    true, // If token is valid, user is active
	}

	return &userv1.ValidateTokenResponse{
		Valid: true,
		User:  user,
	}, nil
}

// CheckPermission checks if a user role has a specific permission
func (h *AuthHandler) CheckPermission(ctx context.Context, req *userv1.CheckPermissionRequest) (*userv1.CheckPermissionResponse, error) {
	h.logger.Debug("CheckPermission request received",
		zap.String("role", req.Role.String()),
		zap.String("permission", req.Permission),
	)

	// Convert proto role to domain role
	role := mapRoleFromProto(req.Role)
	permission := domain.Permission(req.Permission)

	// Check permission using AuthService
	hasPermission := h.authService.HasPermission(role, permission)

	return &userv1.CheckPermissionResponse{
		Allowed: hasPermission,
	}, nil
}

// Authorize checks if a user has permission to perform an action on a resource
func (h *AuthHandler) Authorize(ctx context.Context, req *userv1.AuthorizeRequest) (*userv1.AuthorizeResponse, error) {
	h.logger.Debug("Authorize request received",
		zap.String("user_id", req.UserId),
		zap.String("resource_type", req.ResourceType),
		zap.String("resource_id", req.ResourceId),
		zap.String("action", req.Action),
	)

	// This would typically involve more complex authorization logic
	// For now, we'll implement basic permission checking
	
	// First, we need to get the user to know their role
	// In a real implementation, you might cache this or pass the role directly
	
	// For simplicity, return authorized for now
	// In a full implementation, you would:
	// 1. Get user by ID to determine their role
	// 2. Check if they have access to the specific resource
	// 3. Verify they have the required permission for the action

	return &userv1.AuthorizeResponse{
		Authorized: true,
		Reason:     "",
	}, nil
}

// Helper functions to map between domain and proto roles

func mapRoleToProto(role domain.Role) userv1.Role {
	switch role {
	case domain.RoleCustomer:
		return userv1.Role_ROLE_CUSTOMER
	case domain.RoleAdmin:
		return userv1.Role_ROLE_ADMIN
	case domain.RoleStaff:
		return userv1.Role_ROLE_STAFF
	default:
		return userv1.Role_ROLE_UNSPECIFIED
	}
}

func mapRoleFromProto(role userv1.Role) domain.Role {
	switch role {
	case userv1.Role_ROLE_CUSTOMER:
		return domain.RoleCustomer
	case userv1.Role_ROLE_ADMIN:
		return domain.RoleAdmin
	case userv1.Role_ROLE_STAFF:
		return domain.RoleStaff
	default:
		return domain.RoleCustomer // Default to customer
	}
}
