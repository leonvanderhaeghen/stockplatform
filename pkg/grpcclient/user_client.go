package grpcclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	userpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/user/v1"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

// NewUserClient creates a new gRPC client for the User service
func NewUserClient(addr string) (*UserClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := userpb.NewUserServiceClient(conn)
	return &UserClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// RegisterUser registers a new user
func (c *UserClient) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.RegisterUserResponse, error) {
	return c.client.RegisterUser(ctx, req)
}

// AuthenticateUser authenticates a user and returns a JWT token
func (c *UserClient) AuthenticateUser(ctx context.Context, req *userpb.AuthenticateUserRequest) (*userpb.AuthenticateUserResponse, error) {
	return c.client.AuthenticateUser(ctx, req)
}

// GetUser retrieves a user by ID
func (c *UserClient) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	return c.client.GetUser(ctx, req)
}

// GetUserByEmail retrieves a user by email
func (c *UserClient) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserResponse, error) {
	return c.client.GetUserByEmail(ctx, req)
}

// UpdateUserProfile updates a user's profile information
func (c *UserClient) UpdateUserProfile(ctx context.Context, req *userpb.UpdateUserProfileRequest) (*userpb.UpdateUserProfileResponse, error) {
	return c.client.UpdateUserProfile(ctx, req)
}

// ChangeUserPassword changes a user's password
func (c *UserClient) ChangeUserPassword(ctx context.Context, req *userpb.ChangeUserPasswordRequest) (*userpb.ChangeUserPasswordResponse, error) {
	return c.client.ChangeUserPassword(ctx, req)
}

// DeactivateUser deactivates a user account
func (c *UserClient) DeactivateUser(ctx context.Context, req *userpb.DeactivateUserRequest) (*userpb.DeactivateUserResponse, error) {
	return c.client.DeactivateUser(ctx, req)
}

// ActivateUser activates a user account
func (c *UserClient) ActivateUser(ctx context.Context, req *userpb.ActivateUserRequest) (*userpb.ActivateUserResponse, error) {
	return c.client.ActivateUser(ctx, req)
}

// ListUsers lists all users with optional filtering and pagination
func (c *UserClient) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	return c.client.ListUsers(ctx, req)
}

// CreateUserAddress creates a new address for a user
func (c *UserClient) CreateUserAddress(ctx context.Context, req *userpb.CreateUserAddressRequest) (*userpb.CreateUserAddressResponse, error) {
	return c.client.CreateUserAddress(ctx, req)
}

// GetUserAddresses retrieves all addresses for a user
func (c *UserClient) GetUserAddresses(ctx context.Context, req *userpb.GetUserAddressesRequest) (*userpb.GetUserAddressesResponse, error) {
	return c.client.GetUserAddresses(ctx, req)
}

// GetUserDefaultAddress retrieves the default address for a user
func (c *UserClient) GetUserDefaultAddress(ctx context.Context, req *userpb.GetUserDefaultAddressRequest) (*userpb.GetUserDefaultAddressResponse, error) {
	return c.client.GetUserDefaultAddress(ctx, req)
}

// UpdateUserAddress updates a user address
func (c *UserClient) UpdateUserAddress(ctx context.Context, req *userpb.UpdateUserAddressRequest) (*userpb.UpdateUserAddressResponse, error) {
	return c.client.UpdateUserAddress(ctx, req)
}

// DeleteUserAddress deletes a user address
func (c *UserClient) DeleteUserAddress(ctx context.Context, req *userpb.DeleteUserAddressRequest) (*userpb.DeleteUserAddressResponse, error) {
	return c.client.DeleteUserAddress(ctx, req)
}

// SetDefaultUserAddress sets an address as the default for a user
func (c *UserClient) SetDefaultUserAddress(ctx context.Context, req *userpb.SetDefaultUserAddressRequest) (*userpb.SetDefaultUserAddressResponse, error) {
	return c.client.SetDefaultUserAddress(ctx, req)
}
