package user

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

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
func (c *Client) ValidateToken(ctx context.Context, token string) (bool, error) {
	c.logger.Debug("Validating token")
	
	req := &userv1.ValidateTokenRequest{
		Token: token,
	}
	
	resp, err := c.authClient.ValidateToken(ctx, req)
	if err != nil {
		c.logger.Error("Failed to validate token", zap.Error(err))
		return false, fmt.Errorf("failed to validate token: %w", err)
	}
	
	return resp.Valid, nil
}

// CheckPermission checks if a user role has a specific permission
func (c *Client) CheckPermission(ctx context.Context, role, permission string) (bool, error) {
	c.logger.Debug("Checking permission", zap.String("role", role), zap.String("permission", permission))
	
	req := &userv1.CheckPermissionRequest{
		Role:       role,
		Permission: permission,
	}
	
	resp, err := c.authClient.CheckPermission(ctx, req)
	if err != nil {
		c.logger.Error("Failed to check permission", zap.Error(err))
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	
	return resp.HasPermission, nil
}
