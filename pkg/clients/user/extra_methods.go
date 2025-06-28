package user

import (
    "context"
    "fmt"

    "go.uber.org/zap"

    userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
)

// ChangeUserPassword changes a user's password
func (c *Client) ChangeUserPassword(ctx context.Context, req *userv1.ChangeUserPasswordRequest) (*userv1.ChangeUserPasswordResponse, error) {
    c.logger.Debug("Changing user password", zap.String("user_id", req.GetId()))

    resp, err := c.client.ChangeUserPassword(ctx, req)
    if err != nil {
        c.logger.Error("Failed to change user password", zap.Error(err))
        return nil, fmt.Errorf("failed to change user password: %w", err)
    }

    return resp, nil
}

// GetUserAddresses retrieves addresses for a user
func (c *Client) GetUserAddresses(ctx context.Context, req *userv1.GetUserAddressesRequest) (*userv1.GetUserAddressesResponse, error) {
    c.logger.Debug("Getting user addresses", zap.String("user_id", req.GetUserId()))

    resp, err := c.client.GetUserAddresses(ctx, req)
    if err != nil {
        c.logger.Error("Failed to get user addresses", zap.Error(err))
        return nil, fmt.Errorf("failed to get user addresses: %w", err)
    }

    return resp, nil
}

// CreateUserAddress adds a new address for a user
func (c *Client) CreateUserAddress(ctx context.Context, req *userv1.CreateUserAddressRequest) (*userv1.CreateUserAddressResponse, error) {
    c.logger.Debug("Creating user address", zap.String("user_id", req.GetUserId()))

    resp, err := c.client.CreateUserAddress(ctx, req)
    if err != nil {
        c.logger.Error("Failed to create user address", zap.Error(err))
        return nil, fmt.Errorf("failed to create user address: %w", err)
    }

    return resp, nil
}

// GetUserDefaultAddress retrieves default address for a user
func (c *Client) GetUserDefaultAddress(ctx context.Context, req *userv1.GetUserDefaultAddressRequest) (*userv1.GetUserDefaultAddressResponse, error) {
    c.logger.Debug("Getting default address", zap.String("user_id", req.GetUserId()))

    resp, err := c.client.GetUserDefaultAddress(ctx, req)
    if err != nil {
        c.logger.Error("Failed to get default address", zap.Error(err))
        return nil, fmt.Errorf("failed to get default address: %w", err)
    }

    return resp, nil
}

// UpdateUserAddress updates a user address
func (c *Client) UpdateUserAddress(ctx context.Context, req *userv1.UpdateUserAddressRequest) (*userv1.UpdateUserAddressResponse, error) {
    c.logger.Debug("Updating user address", zap.String("address_id", req.GetId()))

    resp, err := c.client.UpdateUserAddress(ctx, req)
    if err != nil {
        c.logger.Error("Failed to update user address", zap.Error(err))
        return nil, fmt.Errorf("failed to update user address: %w", err)
    }

    return resp, nil
}

// DeleteUserAddress deletes a user address
func (c *Client) DeleteUserAddress(ctx context.Context, req *userv1.DeleteUserAddressRequest) (*userv1.DeleteUserAddressResponse, error) {
    c.logger.Debug("Deleting user address", zap.String("address_id", req.GetId()))

    resp, err := c.client.DeleteUserAddress(ctx, req)
    if err != nil {
        c.logger.Error("Failed to delete user address", zap.Error(err))
        return nil, fmt.Errorf("failed to delete user address: %w", err)
    }

    return resp, nil
}

// SetDefaultUserAddress sets a default address for a user
func (c *Client) SetDefaultUserAddress(ctx context.Context, req *userv1.SetDefaultUserAddressRequest) (*userv1.SetDefaultUserAddressResponse, error) {
    c.logger.Debug("Setting default address", zap.String("address_id", req.GetId()))

    resp, err := c.client.SetDefaultUserAddress(ctx, req)
    if err != nil {
        c.logger.Error("Failed to set default address", zap.Error(err))
        return nil, fmt.Errorf("failed to set default address: %w", err)
    }

    return resp, nil
}

// ActivateUser activates a user account
func (c *Client) ActivateUser(ctx context.Context, req *userv1.ActivateUserRequest) (*userv1.ActivateUserResponse, error) {
    c.logger.Debug("Activating user", zap.String("user_id", req.GetId()))

    resp, err := c.client.ActivateUser(ctx, req)
    if err != nil {
        c.logger.Error("Failed to activate user", zap.Error(err))
        return nil, fmt.Errorf("failed to activate user: %w", err)
    }

    return resp, nil
}

// DeactivateUser deactivates a user account
func (c *Client) DeactivateUser(ctx context.Context, req *userv1.DeactivateUserRequest) (*userv1.DeactivateUserResponse, error) {
    c.logger.Debug("Deactivating user", zap.String("user_id", req.GetId()))

    resp, err := c.client.DeactivateUser(ctx, req)
    if err != nil {
        c.logger.Error("Failed to deactivate user", zap.Error(err))
        return nil, fmt.Errorf("failed to deactivate user: %w", err)
    }

    return resp, nil
}
