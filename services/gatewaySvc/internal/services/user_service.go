package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	userclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/user"
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

	// Call the gRPC service via client abstraction
	resp, err := s.client.RegisterUser(ctx, email, password, firstName, lastName, role)
	if err != nil {
		s.logger.Error("Failed to register user",
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return resp, nil
}

// AuthenticateUser authenticates a user
func (s *UserServiceImpl) AuthenticateUser(
	ctx context.Context,
	email, password string,
) (interface{}, error) {
	s.logger.Debug("AuthenticateUser",
		zap.String("email", email),
	)

	// Call the gRPC service via client abstraction
	resp, err := s.client.AuthenticateUser(ctx, email, password)
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

	resp, err := s.client.GetUser(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user by ID",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp, nil
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

	_, err := s.client.UpdateUserProfile(ctx, userID, firstName, lastName, phone)
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

	err := s.client.ChangeUserPassword(ctx, userID, currentPassword, newPassword)
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

	addresses, err := s.client.GetUserAddresses(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user addresses",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}

	return addresses, nil
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
	)

	address, err := s.client.CreateUserAddress(ctx, userID, name, street, city, state, postalCode, country, phone, isDefault)
	if err != nil {
		s.logger.Error("Failed to create user address",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create user address: %w", err)
	}

	return address, nil
}

// GetUserDefaultAddress gets the default address for a user
func (s *UserServiceImpl) GetUserDefaultAddress(
	ctx context.Context,
	userID string,
) (interface{}, error) {
	s.logger.Debug("GetUserDefaultAddress",
		zap.String("userID", userID),
	)

	address, err := s.client.GetUserDefaultAddress(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get user default address",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user default address: %w", err)
	}

	return address, nil
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

	err := s.client.UpdateUserAddress(ctx, addressID, userID, name, street, city, state, postalCode, country, phone, isDefault)
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

	err := s.client.DeleteUserAddress(ctx, addressID, userID)
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

	err := s.client.SetDefaultUserAddress(ctx, addressID, userID)
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
	logFields := []zap.Field{
		zap.String("role", role),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	}
	
	if active != nil {
		logFields = append(logFields, zap.Bool("active", *active))
	}
	
	s.logger.Debug("ListUsers", logFields...)

	// Fix type conversion issues - convert parameters to expected types
	activeValue := false
	if active != nil {
		activeValue = *active
	}
	resp, err := s.client.ListUsers(ctx, role, activeValue, int32(limit), int32(offset))
	if err != nil {
		s.logger.Error("Failed to list users",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return resp, nil
}

// ActivateUser activates a user account
func (s *UserServiceImpl) ActivateUser(ctx context.Context, userID string) error {
	s.logger.Debug("ActivateUser", zap.String("userID", userID))

	err := s.client.ActivateUser(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to activate user",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to activate user: %w", err)
	}

	return nil
}

// DeactivateUser deactivates a user account
func (s *UserServiceImpl) DeactivateUser(ctx context.Context, userID string) error {
	s.logger.Debug("DeactivateUser", zap.String("userID", userID))

	err := s.client.DeactivateUser(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to deactivate user",
			zap.String("userID", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}
