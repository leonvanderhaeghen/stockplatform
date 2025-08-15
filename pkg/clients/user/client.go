package user

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
)

// Client provides a high-level interface for interacting with the User service
type Client struct {
	conn       *grpc.ClientConn
	client     userv1.UserServiceClient
	authClient userv1.AuthServiceClient
	logger     *zap.Logger
}

// Config holds configuration for the User client
type Config struct {
	Address string
	Timeout time.Duration
}

// New creates a new User service client
func New(config Config, logger *zap.Logger) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	conn, err := grpc.Dial(config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(config.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userv1.NewUserServiceClient(conn)
	authClient := userv1.NewAuthServiceClient(conn)

	return &Client{
		conn:       conn,
		client:     client,
		authClient: authClient,
		logger:     logger,
	}, nil
}

// Close closes the connection to the User service
func (c *Client) Close() error {
	return c.conn.Close()
}

// RegisterUser registers a new user
func (c *Client) RegisterUser(ctx context.Context, email, password, firstName, lastName, role string) (*models.RegisterUserResponse, error) {
	c.logger.Debug("Registering user", zap.String("email", email))

	req := &userv1.RegisterUserRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
	}

	resp, err := c.client.RegisterUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to register user", zap.Error(err))
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	c.logger.Debug("User registered successfully", zap.String("id", resp.User.Id))
	return c.convertToRegisterUserResponse(resp), nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func (c *Client) AuthenticateUser(ctx context.Context, email, password string) (*models.AuthenticateUserResponse, error) {
	c.logger.Debug("Authenticating user", zap.String("email", email))

	req := &userv1.AuthenticateUserRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.client.AuthenticateUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to authenticate user", zap.Error(err))
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}

	return c.convertToAuthenticateUserResponse(resp), nil
}

// GetUser retrieves a user by ID
func (c *Client) GetUser(ctx context.Context, id string) (*models.User, error) {
	c.logger.Debug("Getting user", zap.String("id", id))

	req := &userv1.GetUserRequest{
		Id: id,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return c.convertToUser(resp.User), nil
}

// GetUserByEmail retrieves a user by email
func (c *Client) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	c.logger.Debug("Getting user by email", zap.String("email", email))

	req := &userv1.GetUserByEmailRequest{
		Email: email,
	}

	resp, err := c.client.GetUserByEmail(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return c.convertToUser(resp.User), nil
}

// UpdateUserProfile updates a user's profile information
func (c *Client) UpdateUserProfile(ctx context.Context, id, firstName, lastName, phone string) (*models.UpdateUserProfileResponse, error) {
	c.logger.Debug("Updating user profile", zap.String("id", id))

	req := &userv1.UpdateUserProfileRequest{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
	}

	resp, err := c.client.UpdateUserProfile(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update user profile", zap.Error(err))
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return c.convertToUpdateUserProfileResponse(resp), nil
}

// ListUsers lists all users with optional filtering and pagination
func (c *Client) ListUsers(ctx context.Context, role string, active bool, limit, offset int32) (*models.ListUsersResponse, error) {
	c.logger.Debug("Listing users", zap.String("role", role), zap.Bool("active", active))

	req := &userv1.ListUsersRequest{
		Role:   role,
		Active: active,
		Limit:  limit,
		Offset: offset,
	}

	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list users", zap.Error(err))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return c.convertToListUsersResponse(resp), nil
}

// ValidateToken validates a JWT token and returns user claims
func (c *Client) ValidateToken(ctx context.Context, token string) (*models.User, bool, error) {
	c.logger.Debug("Validating token")

	req := &userv1.ValidateTokenRequest{
		Token: token,
	}

	resp, err := c.authClient.ValidateToken(ctx, req)
	if err != nil {
		c.logger.Error("Failed to validate token", zap.Error(err))
		return nil, false, fmt.Errorf("failed to validate token: %w", err)
	}

	if !resp.GetValid() {
		return nil, false, nil
	}

	// Convert protobuf user to domain model
	user := c.convertToUser(resp.GetUser())
	return user, true, nil
}

// CheckPermission checks if a user role has a specific permission
func (c *Client) CheckPermission(ctx context.Context, role, permission string) (bool, error) {
	req := &userv1.CheckPermissionRequest{
		Role:       stringToRole(role),
		Permission: permission,
	}

	resp, err := c.authClient.CheckPermission(ctx, req)
	if err != nil {
		c.logger.Error("Failed to check permission", zap.Error(err))
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return resp.GetAllowed(), nil
}

// ChangeUserPassword changes a user's password
func (c *Client) ChangeUserPassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	req := &userv1.ChangeUserPasswordRequest{
		Id:              userID,
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}

	resp, err := c.client.ChangeUserPassword(ctx, req)
	if err != nil {
		c.logger.Error("Failed to change user password", zap.Error(err))
		return fmt.Errorf("failed to change user password: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("password change failed")
	}

	return nil
}

// CreateUserAddress creates a new address for a user
func (c *Client) CreateUserAddress(ctx context.Context, userID, name, street, city, state, postalCode, country, phone string, isDefault bool) (*models.UserAddress, error) {
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

	resp, err := c.client.CreateUserAddress(ctx, req)
	if err != nil {
		c.logger.Error("Failed to create user address", zap.Error(err))
		return nil, fmt.Errorf("failed to create user address: %w", err)
	}

	return convertAddressFromProto(resp.GetAddress()), nil
}

// GetUserAddresses retrieves all addresses for a user
func (c *Client) GetUserAddresses(ctx context.Context, userID string) ([]*models.UserAddress, error) {
	req := &userv1.GetUserAddressesRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUserAddresses(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user addresses", zap.Error(err))
		return nil, fmt.Errorf("failed to get user addresses: %w", err)
	}

	addresses := make([]*models.UserAddress, len(resp.GetAddresses()))
	for i, addr := range resp.GetAddresses() {
		addresses[i] = convertAddressFromProto(addr)
	}

	return addresses, nil
}

// GetUserDefaultAddress retrieves the default address for a user
func (c *Client) GetUserDefaultAddress(ctx context.Context, userID string) (*models.UserAddress, error) {
	req := &userv1.GetUserDefaultAddressRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUserDefaultAddress(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user default address", zap.Error(err))
		return nil, fmt.Errorf("failed to get user default address: %w", err)
	}

	return convertAddressFromProto(resp.GetAddress()), nil
}

// UpdateUserAddress updates a user address
func (c *Client) UpdateUserAddress(ctx context.Context, addressID, userID, name, street, city, state, postalCode, country, phone string, isDefault bool) error {
	req := &userv1.UpdateUserAddressRequest{
		Id:         addressID,
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

	resp, err := c.client.UpdateUserAddress(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update user address", zap.Error(err))
		return fmt.Errorf("failed to update user address: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("address update failed")
	}

	return nil
}

// DeleteUserAddress deletes a user address
func (c *Client) DeleteUserAddress(ctx context.Context, addressID, userID string) error {
	req := &userv1.DeleteUserAddressRequest{
		Id:     addressID,
		UserId: userID,
	}

	resp, err := c.client.DeleteUserAddress(ctx, req)
	if err != nil {
		c.logger.Error("Failed to delete user address", zap.Error(err))
		return fmt.Errorf("failed to delete user address: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("address deletion failed")
	}

	return nil
}

// SetDefaultUserAddress sets an address as the default for a user
func (c *Client) SetDefaultUserAddress(ctx context.Context, addressID, userID string) error {
	req := &userv1.SetDefaultUserAddressRequest{
		Id:     addressID,
		UserId: userID,
	}

	resp, err := c.client.SetDefaultUserAddress(ctx, req)
	if err != nil {
		c.logger.Error("Failed to set default user address", zap.Error(err))
		return fmt.Errorf("failed to set default user address: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("setting default address failed")
	}

	return nil
}

// ActivateUser activates a user account
func (c *Client) ActivateUser(ctx context.Context, userID string) error {
	req := &userv1.ActivateUserRequest{
		Id: userID,
	}

	resp, err := c.client.ActivateUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to activate user", zap.Error(err))
		return fmt.Errorf("failed to activate user: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("user activation failed")
	}

	return nil
}

// DeactivateUser deactivates a user account
func (c *Client) DeactivateUser(ctx context.Context, userID string) error {
	req := &userv1.DeactivateUserRequest{
		Id: userID,
	}

	resp, err := c.client.DeactivateUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to deactivate user", zap.Error(err))
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	if !resp.GetSuccess() {
		return fmt.Errorf("user deactivation failed")
	}

	return nil
}

// convertAddressFromProto converts a protobuf Address to a domain UserAddress
func convertAddressFromProto(protoAddr *userv1.Address) *models.UserAddress {
	if protoAddr == nil {
		return nil
	}

	return &models.UserAddress{
		ID:         protoAddr.GetId(),
		UserID:     protoAddr.GetUserId(),
		Name:       protoAddr.GetName(),
		Street:     protoAddr.GetStreet(),
		City:       protoAddr.GetCity(),
		State:      protoAddr.GetState(),
		PostalCode: protoAddr.GetPostalCode(),
		Country:    protoAddr.GetCountry(),
		Phone:      protoAddr.GetPhone(),
		IsDefault:  protoAddr.GetIsDefault(),
		// Note: CreatedAt and UpdatedAt would need to be converted from proto timestamps
		// if they were included in the proto definition
	}
}

// stringToRole converts a string role to the protobuf Role enum
func stringToRole(role string) userv1.Role {
	switch role {
	case "admin":
		return userv1.Role_ROLE_ADMIN
	case "customer":
		return userv1.Role_ROLE_CUSTOMER
	case "staff":
		return userv1.Role_ROLE_STAFF
	case "manager":
		return userv1.Role_ROLE_MANAGER
	case "supplier":
		return userv1.Role_ROLE_SUPPLIER
	default:
		return userv1.Role_ROLE_UNSPECIFIED
	}
}
