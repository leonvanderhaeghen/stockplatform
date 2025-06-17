package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
	userclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/user"
)

var (
	_ = userv1.UserService_ServiceDesc // Ensure the service is linked
)

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	client *userclient.Client
	logger *zap.Logger
}

// NewUserService creates a new instance of UserServiceImpl
func NewUserService(userServiceAddr string, logger *zap.Logger) (UserService, error) {
	// Create a gRPC client via new abstraction
	usrCfg := userclient.Config{Address: userServiceAddr}
	client, err := userclient.New(usrCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create user client: %w", err)
	}

	return &UserServiceImpl{
		client: client,
		logger: logger.Named("user_service"),
	}, nil
}

// RegisterUser registers a new user
func (s *UserServiceImpl) RegisterUser(
	ctx context.Context,
	email, password, firstName, lastName, role string,
) (interface{}, error) {
	s.logger.Debug("RegisterUser",
		zap.String("email", email),
		zap.String("firstName", firstName),
		zap.String("lastName", lastName),
		zap.String("role", role),
	)

	// Create the request
	req := &userv1.RegisterUserRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	}

	// Call the gRPC service
	resp, err := s.client.RegisterUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to register user",
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return resp.GetUser(), nil
}

// AuthenticateUser authenticates a user
func (s *UserServiceImpl) AuthenticateUser(
	ctx context.Context,
	email, password string,
) (interface{}, error) {
	s.logger.Debug("AuthenticateUser",
		zap.String("email", email),
	)

	// Create the request
	req := &userv1.AuthenticateUserRequest{
		Email:    email,
		Password: password,
	}

	// Call the gRPC service
	resp, err := s.client.AuthenticateUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to authenticate user",
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}

	return resp, nil
}

// GetUserByID gets a user by ID
func (s *UserServiceImpl) GetUserByID(
	ctx context.Context,
	userID string,
) (interface{}, error) {
	s.logger.Debug("GetUserByID",
		zap.String("userID", userID),
	)

	req := &userv1.GetUserRequest{
		Id: userID,
	}

	resp, err := s.client.GetUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get user by ID",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp.GetUser(), nil
}

// UpdateUserProfile updates a user's profile
func (s *UserServiceImpl) UpdateUserProfile(
	ctx context.Context,
	userID, firstName, lastName, phone string,
) error {
	s.logger.Debug("UpdateUserProfile",
		zap.String("userID", userID),
		zap.String("firstName", firstName),
		zap.String("lastName", lastName),
	)

	req := &userv1.UpdateUserProfileRequest{
		Id:        userID,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
	}

	_, err := s.client.UpdateUserProfile(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update user profile",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	return nil
}

// ChangeUserPassword changes a user's password
func (s *UserServiceImpl) ChangeUserPassword(
	ctx context.Context,
	userID, currentPassword, newPassword string,
) error {
	s.logger.Debug("ChangeUserPassword",
		zap.String("userID", userID),
	)

	req := &userv1.ChangeUserPasswordRequest{
		Id:              userID,
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	_, err := s.client.ChangeUserPassword(ctx, req)
	if err != nil {
		s.logger.Error("Failed to change user password",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to change user password: %w", err)
	}

	return nil
}

// GetUserAddresses gets addresses for a user
func (s *UserServiceImpl) GetUserAddresses(
	ctx context.Context,
	userID string,
) (interface{}, error) {
	s.logger.Debug("GetUserAddresses",
		zap.String("userID", userID),
	)

	req := &userv1.GetUserAddressesRequest{
		UserId: userID,
	}

	resp, err := s.client.GetUserAddresses(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get user addresses",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}

	return resp.GetAddresses(), nil
}

// CreateUserAddress creates a new address for a user
func (s *UserServiceImpl) CreateUserAddress(
	ctx context.Context,
	userID, name, street, city, state, postalCode, country, phone string,
	isDefault bool,
) (interface{}, error) {
	s.logger.Debug("CreateUserAddress",
		zap.String("userID", userID),
		zap.String("name", name),
		zap.String("city", city),
		zap.String("country", country),
	)

	req := &userv1.CreateUserAddressRequest{
		UserId:     userID,
		Name:       name,
		Street:     street,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
		Phone:      phone,
		IsDefault:  isDefault,
	}

	resp, err := s.client.CreateUserAddress(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create user address",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create user address: %w", err)
	}

	return resp.GetAddress(), nil
}

// GetUserDefaultAddress gets the default address for a user
func (s *UserServiceImpl) GetUserDefaultAddress(
	ctx context.Context,
	userID string,
) (interface{}, error) {
	s.logger.Debug("GetUserDefaultAddress",
		zap.String("userID", userID),
	)

	req := &userv1.GetUserDefaultAddressRequest{
		UserId: userID,
	}

	resp, err := s.client.GetUserDefaultAddress(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get user default address",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user default address: %w", err)
	}

	return resp.GetAddress(), nil
}

// UpdateUserAddress updates an address for a user
func (s *UserServiceImpl) UpdateUserAddress(
	ctx context.Context,
	addressID, userID, name, street, city, state, postalCode, country, phone string,
	isDefault bool,
) error {
	s.logger.Debug("UpdateUserAddress",
		zap.String("addressID", addressID),
		zap.String("userID", userID),
	)

	req := &userv1.UpdateUserAddressRequest{
		Id:          addressID,
		UserId:      userID,
		Name:        name,
		Street:      street,
		City:        city,
		State:       state,
		PostalCode:  postalCode,
		Country:     country,
		Phone:       phone,
		IsDefault:   isDefault,
	}

	_, err := s.client.UpdateUserAddress(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update user address",
			zap.String("addressID", addressID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update user address: %w", err)
	}

	return nil
}

// DeleteUserAddress deletes an address for a user
func (s *UserServiceImpl) DeleteUserAddress(
	ctx context.Context,
	addressID, userID string,
) error {
	s.logger.Debug("DeleteUserAddress",
		zap.String("addressID", addressID),
		zap.String("userID", userID),
	)

	req := &userv1.DeleteUserAddressRequest{
		Id:     addressID,
		UserId: userID,
	}

	_, err := s.client.DeleteUserAddress(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete user address",
			zap.String("addressID", addressID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete user address: %w", err)
	}

	return nil
}

// SetDefaultUserAddress sets an address as default for a user
func (s *UserServiceImpl) SetDefaultUserAddress(
	ctx context.Context,
	addressID, userID string,
) error {
	s.logger.Debug("SetDefaultUserAddress",
		zap.String("addressID", addressID),
		zap.String("userID", userID),
	)

	req := &userv1.SetDefaultUserAddressRequest{
		Id:     addressID,
		UserId: userID,
	}

	_, err := s.client.SetDefaultUserAddress(ctx, req)
	if err != nil {
		s.logger.Error("Failed to set default user address",
			zap.String("addressID", addressID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to set default user address: %w", err)
	}

	return nil
}

// ListUsers lists all users (admin only)
func (s *UserServiceImpl) ListUsers(
	ctx context.Context,
	role string,
	active *bool,
	limit, offset int,
) (interface{}, error) {
	s.logger.Debug("ListUsers",
		zap.String("role", role),
		zap.Bool("active", *active),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	req := &userv1.ListUsersRequest{
		Role:   role,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	if active != nil {
		req.Active = *active
	}

	resp, err := s.client.ListUsers(ctx, req)
	if err != nil {
		s.logger.Error("Failed to list users",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return resp.GetUsers(), nil
}

// ActivateUser activates a user (admin only)
func (s *UserServiceImpl) ActivateUser(
	ctx context.Context,
	userID string,
) error {
	s.logger.Debug("ActivateUser",
		zap.String("userID", userID),
	)

	req := &userv1.ActivateUserRequest{
		Id: userID,
	}

	_, err := s.client.ActivateUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to activate user",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// DeactivateUser deactivates a user (admin only)
func (s *UserServiceImpl) DeactivateUser(
	ctx context.Context,
	userID string,
) error {
	s.logger.Debug("DeactivateUser",
		zap.String("userID", userID),
	)

	req := &userv1.DeactivateUserRequest{
		Id: userID,
	}

	_, err := s.client.DeactivateUser(ctx, req)
	if err != nil {
		s.logger.Error("Failed to deactivate user",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}
