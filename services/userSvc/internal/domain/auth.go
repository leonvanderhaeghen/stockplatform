package domain

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Permission represents system permissions
type Permission string

const (
	// User permissions
	PermissionUserRead   Permission = "user:read"
	PermissionUserWrite  Permission = "user:write"
	PermissionUserDelete Permission = "user:delete"
	
	// Order permissions
	PermissionOrderRead   Permission = "order:read"
	PermissionOrderWrite  Permission = "order:write"
	PermissionOrderDelete Permission = "order:delete"
	
	// Product permissions
	PermissionProductRead   Permission = "product:read"
	PermissionProductWrite  Permission = "product:write"
	PermissionProductDelete Permission = "product:delete"
	
	// Inventory permissions
	PermissionInventoryRead   Permission = "inventory:read"
	PermissionInventoryWrite  Permission = "inventory:write"
	PermissionInventoryDelete Permission = "inventory:delete"
	
	// Store permissions
	PermissionStoreRead   Permission = "store:read"
	PermissionStoreWrite  Permission = "store:write"
	PermissionStoreDelete Permission = "store:delete"
	
	// Supplier permissions
	PermissionSupplierRead   Permission = "supplier:read"
	PermissionSupplierWrite  Permission = "supplier:write"
	PermissionSupplierDelete Permission = "supplier:delete"
	
	// Admin permissions
	PermissionAdminAll Permission = "admin:all"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role][]Permission{
	RoleCustomer: {
		PermissionUserRead,
		PermissionOrderRead,
		PermissionOrderWrite,
		PermissionProductRead,
		PermissionInventoryRead,
	},
	RoleStaff: {
		PermissionUserRead,
		PermissionOrderRead,
		PermissionOrderWrite,
		PermissionOrderDelete,
		PermissionProductRead,
		PermissionProductWrite,
		PermissionInventoryRead,
		PermissionInventoryWrite,
		PermissionStoreRead,
		PermissionStoreWrite,
		PermissionSupplierRead,
	},
	RoleAdmin: {
		PermissionAdminAll, // Admin has all permissions
	},
}

// AuthToken represents a JWT authentication token
type AuthToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	TokenType string    `json:"token_type"`
}

// Claims represents JWT claims
type Claims struct {
	UserID    string   `json:"user_id"`
	Email     string   `json:"email"`
	Role      Role     `json:"role"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	jwt.RegisteredClaims
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	User  *User      `json:"user"`
	Token *AuthToken `json:"token"`
}

// AuthService interface defines authentication operations
type AuthService interface {
	// Login authenticates a user and returns a JWT token
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	
	// ValidateToken validates a JWT token and returns claims
	ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
	
	// RefreshToken refreshes an expired token
	RefreshToken(ctx context.Context, tokenString string) (*AuthToken, error)
	
	// HasPermission checks if a role has a specific permission
	HasPermission(role Role, permission Permission) bool
	
	// Logout invalidates a token (placeholder for token blacklisting)
	Logout(ctx context.Context, tokenString string) error
}

// Common authentication errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrUserNotActive      = errors.New("user account is not active")
)

// HasPermission checks if a role has a specific permission
func HasPermission(role Role, permission Permission) bool {
	permissions, exists := RolePermissions[role]
	if !exists {
		return false
	}
	
	// Admin role has all permissions
	if role == RoleAdmin {
		return true
	}
	
	// Check specific permissions
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	
	return false
}

// GetRolePermissions returns all permissions for a role
func GetRolePermissions(role Role) []Permission {
	if role == RoleAdmin {
		// Return all permissions for admin
		var allPermissions []Permission
		for _, permissions := range RolePermissions {
			allPermissions = append(allPermissions, permissions...)
		}
		// Add admin-specific permission
		allPermissions = append(allPermissions, PermissionAdminAll)
		return allPermissions
	}
	
	permissions, exists := RolePermissions[role]
	if !exists {
		return []Permission{}
	}
	
	return permissions
}
