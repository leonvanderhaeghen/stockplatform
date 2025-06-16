package user

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"go.uber.org/zap"

	userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
)

// Client provides a high-level interface for interacting with the User service
type Client struct {
	conn   *grpc.ClientConn
	client userv1.UserServiceClient
	logger *zap.Logger
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

	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the connection to the User service
func (c *Client) Close() error {
	return c.conn.Close()
}

// RegisterUser registers a new user
func (c *Client) RegisterUser(ctx context.Context, req *userv1.RegisterUserRequest) (*userv1.RegisterUserResponse, error) {
	c.logger.Debug("Registering user", zap.String("email", req.Email))
	
	resp, err := c.client.RegisterUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to register user", zap.Error(err))
		return nil, fmt.Errorf("failed to register user: %w", err)
	}
	
	c.logger.Debug("User registered successfully", zap.String("id", resp.User.Id))
	return resp, nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func (c *Client) AuthenticateUser(ctx context.Context, req *userv1.AuthenticateUserRequest) (*userv1.AuthenticateUserResponse, error) {
	c.logger.Debug("Authenticating user", zap.String("email", req.Email))
	
	resp, err := c.client.AuthenticateUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to authenticate user", zap.Error(err))
		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}
	
	return resp, nil
}

// GetUser retrieves a user by ID
func (c *Client) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	c.logger.Debug("Getting user", zap.String("id", req.Id))
	
	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return resp, nil
}

// GetUserByEmail retrieves a user by email
func (c *Client) GetUserByEmail(ctx context.Context, req *userv1.GetUserByEmailRequest) (*userv1.GetUserResponse, error) {
	c.logger.Debug("Getting user by email", zap.String("email", req.Email))
	
	resp, err := c.client.GetUserByEmail(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return resp, nil
}

// UpdateUserProfile updates a user's profile information
func (c *Client) UpdateUserProfile(ctx context.Context, req *userv1.UpdateUserProfileRequest) (*userv1.UpdateUserProfileResponse, error) {
	c.logger.Debug("Updating user profile", zap.String("id", req.Id))
	
	resp, err := c.client.UpdateUserProfile(ctx, req)
	if err != nil {
		c.logger.Error("Failed to update user profile", zap.Error(err))
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}
	
	return resp, nil
}

// ListUsers lists all users with optional filtering and pagination
func (c *Client) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	c.logger.Debug("Listing users")
	
	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		c.logger.Error("Failed to list users", zap.Error(err))
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	
	return resp, nil
}
