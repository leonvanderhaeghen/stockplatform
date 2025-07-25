package user

import (
	"time"
	
	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
)

// convertToUser converts protobuf User to domain User
func (c *Client) convertToUser(proto *userv1.User) *models.User {
	if proto == nil {
		return nil
	}

	user := &models.User{
		ID:        proto.Id,
		Email:     proto.Email,
		FirstName: proto.FirstName,
		LastName:  proto.LastName,
		Role:      roleFromProto(proto.Role),
		IsActive:  proto.Active,
	}

	// Handle timestamps (convert from strings to time.Time)
	if proto.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, proto.CreatedAt); err == nil {
			user.CreatedAt = t
		}
	}
	if proto.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, proto.UpdatedAt); err == nil {
			user.UpdatedAt = t
		}
	}

	return user
}

// convertFromUser converts domain User to protobuf User
func (c *Client) convertFromUser(user *models.User) *userv1.User {
	if user == nil {
		return nil
	}

	proto := &userv1.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      roleToProto(user.Role),
		Active:    user.IsActive,
	}

	// Handle timestamps (convert from time.Time to strings)
	if !user.CreatedAt.IsZero() {
		proto.CreatedAt = user.CreatedAt.Format(time.RFC3339)
	}
	if !user.UpdatedAt.IsZero() {
		proto.UpdatedAt = user.UpdatedAt.Format(time.RFC3339)
	}

	return proto
}

// convertToRegisterUserResponse converts protobuf RegisterUserResponse to domain RegisterUserResponse
func (c *Client) convertToRegisterUserResponse(proto *userv1.RegisterUserResponse) *models.RegisterUserResponse {
	if proto == nil {
		return nil
	}

	return &models.RegisterUserResponse{
		User:    c.convertToUser(proto.User),
		Message: "", // protobuf doesn't have message field
	}
}

// convertToAuthenticateUserResponse converts protobuf AuthenticateUserResponse to domain AuthenticateUserResponse
func (c *Client) convertToAuthenticateUserResponse(proto *userv1.AuthenticateUserResponse) *models.AuthenticateUserResponse {
	if proto == nil {
		return nil
	}

	return &models.AuthenticateUserResponse{
		Token:     proto.Token,
		ExpiresAt: 0, // protobuf doesn't have expires_at field
		User:      c.convertToUser(proto.User),
	}
}

// convertToListUsersResponse converts protobuf ListUsersResponse to domain ListUsersResponse
func (c *Client) convertToListUsersResponse(proto *userv1.ListUsersResponse) *models.ListUsersResponse {
	if proto == nil {
		return nil
	}

	users := make([]*models.User, len(proto.Users))
	for i, protoUser := range proto.Users {
		users[i] = c.convertToUser(protoUser)
	}

	return &models.ListUsersResponse{
		Users:      users,
		TotalCount: int32(len(users)), // protobuf doesn't have total_count field
	}
}

// convertToUpdateUserProfileResponse converts protobuf UpdateUserProfileResponse to domain UpdateUserProfileResponse
func (c *Client) convertToUpdateUserProfileResponse(proto *userv1.UpdateUserProfileResponse) *models.UpdateUserProfileResponse {
	if proto == nil {
		return nil
	}

	return &models.UpdateUserProfileResponse{
		User:    nil,    // protobuf doesn't return user in response
		Message: "User profile updated successfully",
	}
}

// Note: ValidateTokenResponse and CheckPermissionResponse are not used in current client
// They would be part of AuthService which may be implemented separately



// Helper functions to convert between role types
func roleFromProto(protoRole userv1.Role) string {
	switch protoRole {
	case userv1.Role_ROLE_CUSTOMER:
		return "customer"
	case userv1.Role_ROLE_ADMIN:
		return "admin"
	case userv1.Role_ROLE_STAFF:
		return "staff"
	// Note: MANAGER and SUPPLIER roles not defined in current protobuf
	default:
		return "customer"
	}
}

func roleToProto(role string) userv1.Role {
	switch role {
	case "customer":
		return userv1.Role_ROLE_CUSTOMER
	case "admin":
		return userv1.Role_ROLE_ADMIN
	case "staff":
		return userv1.Role_ROLE_STAFF
	// Note: MANAGER and SUPPLIER roles not defined in current protobuf
	default:
		return userv1.Role_ROLE_CUSTOMER
	}
}
