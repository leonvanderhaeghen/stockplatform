package application

import (
	"context"
	"fmt"
	"time"

	orderclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/order"
	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
	"go.uber.org/zap"
)

// UserAuthService handles user authentication and authorization
type UserAuthService struct {
	orderClient *orderclient.Client
	userService *UserService
	logger      *zap.Logger
}

// NewUserAuthService creates a new UserAuthService
func NewUserAuthService(
	userService *UserService,
	orderServiceAddr string,
	logger *zap.Logger,
) (*UserAuthService, error) {
	// Initialize the order client
	ordCfg := orderclient.Config{Address: orderServiceAddr}
	orderClient, err := orderclient.New(ordCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create order client: %w", err)
	}

	return &UserAuthService{
		orderClient: orderClient,
		userService: userService,
		logger:      logger.Named("user_auth_service"),
	}, nil
}

// Close closes any open connections
func (s *UserAuthService) Close() error {
	if s.orderClient != nil {
		return s.orderClient.Close()
	}
	return nil
}

// AuthenticateUser authenticates a user with email and password
func (s *UserAuthService) AuthenticateUser(
	ctx context.Context,
	email, password string,
) (*domain.User, string, error) {
	// Get user by email
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !user.CheckPassword(password) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := generateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// GetUserOrders retrieves orders for a specific user
func (s *UserAuthService) GetUserOrders(
	ctx context.Context,
	userID string,
) ([]*orderpb.Order, error) {
	// Call the order service to get user orders
	resp, err := s.orderClient.GetUserOrders(ctx, &orderpb.GetUserOrdersRequest{
		UserId: userID,
	})
	if err != nil {
		s.logger.Error("Failed to get user orders",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	return resp.Orders, nil
}

// UpdateUserProfile updates user profile and handles related operations
func (s *UserAuthService) UpdateUserProfile(
	ctx context.Context,
	userID string,
	update *domain.User,
) (*domain.User, error) {
	// Update user profile using the user service
	err := s.userService.UpdateUserProfile(ctx, userID, update.FirstName, update.LastName, update.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Get the updated user
	updatedUser, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	return updatedUser, nil
}

// DeleteUserAccount handles user account deletion and cleanup
func (s *UserAuthService) DeleteUserAccount(
	ctx context.Context,
	userID string,
) error {
	// Deactivate the user account instead of deleting for data retention
	err := s.userService.DeactivateUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Note: In a real application, you might want to:
	// 1. Cancel any pending orders
	// 2. Anonymize personal data
	// 3. Log the deletion for audit purposes

	s.logger.Info("User account deactivated",
		zap.String("user_id", userID),
	)

	return nil
}

// generateJWT generates a JWT token for the user
func generateJWT(userID string) (string, error) {
	// This is a simplified example. In a real system, you would use a JWT library
	// and include proper claims, expiration, and signing.
	// For example, using github.com/golang-jwt/jwt/v5:
	//
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//     "user_id": userID,
	//     "exp":     time.Now().Add(24 * time.Hour).Unix(),
	// })
	// return token.SignedString([]byte("your-secret-key"))

	// For now, return a simple token
	return fmt.Sprintf("token-%s-%d", userID, time.Now().Unix()), nil
}
