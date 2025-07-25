package models

import "time"

// User represents a user in the domain
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterUserResponse represents the response from user registration
type RegisterUserResponse struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
}

// AuthenticateUserResponse represents the response from user authentication
type AuthenticateUserResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      *User  `json:"user"`
}

// ListUsersResponse represents the response from listing users
type ListUsersResponse struct {
	Users      []*User `json:"users"`
	TotalCount int32   `json:"total_count"`
}

// UpdateUserProfileResponse represents the response from updating user profile
type UpdateUserProfileResponse struct {
	User    *User  `json:"user"`
	Message string `json:"message"`
}

// ValidateTokenResponse represents the response from token validation
type ValidateTokenResponse struct {
	Valid     bool   `json:"valid"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"expires_at"`
}

// CheckPermissionResponse represents the response from permission check
type CheckPermissionResponse struct {
	HasPermission bool   `json:"has_permission"`
	Role          string `json:"role"`
	Permission    string `json:"permission"`
}
