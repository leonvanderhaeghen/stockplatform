package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo    domain.UserRepository
	addressRepo domain.AddressRepository
	jwtSecret   string
	logger      *zap.Logger
}

// NewUserService creates a new user service
func NewUserService(
	userRepo domain.UserRepository,
	addressRepo domain.AddressRepository,
	jwtSecret string,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		addressRepo: addressRepo,
		jwtSecret:   jwtSecret,
		logger:      logger.Named("user_service"),
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	s.logger.Info("Registering new user",
		zap.String("email", email),
		zap.String("first_name", firstName),
		zap.String("last_name", lastName),
	)

	// Validate input
	if email == "" || password == "" || firstName == "" || lastName == "" {
		return nil, errors.New("all fields are required")
	}

	// Check if email is already in use
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email is already in use")
	}

	// Create new user with customer role
	user, err := domain.NewUser(email, password, firstName, lastName, domain.RoleCustomer)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	s.logger.Info("Authenticating user", zap.String("email", email))

	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.Active {
		return "", errors.New("account is deactivated")
	}

	// Verify password
	if !user.CheckPassword(password) {
		return "", errors.New("invalid email or password")
	}

	// Update last login time
	user.RecordLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Warn("Failed to update last login time",
			zap.Error(err),
			zap.String("id", user.ID),
		)
		// Continue anyway as this is not critical
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.FullName(),
		"email": user.Email,
		"role": string(user.Role),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("Failed to generate JWT token", zap.Error(err))
		return "", errors.New("failed to generate authentication token")
	}

	return tokenString, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	s.logger.Debug("Getting user by ID", zap.String("id", id))
	
	if id == "" {
		return nil, errors.New("user ID is required")
	}
	
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	s.logger.Debug("Getting user by email", zap.String("email", email))
	
	if email == "" {
		return nil, errors.New("email is required")
	}
	
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(ctx context.Context, id, firstName, lastName, phone string) error {
	s.logger.Info("Updating user profile",
		zap.String("id", id),
		zap.String("first_name", firstName),
		zap.String("last_name", lastName),
	)
	
	if id == "" {
		return errors.New("user ID is required")
	}
	
	if firstName == "" || lastName == "" {
		return errors.New("first name and last name are required")
	}
	
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if user == nil {
		return errors.New("user not found")
	}
	
	user.UpdateProfile(firstName, lastName, phone)
	return s.userRepo.Update(ctx, user)
}

// ChangeUserPassword changes a user's password
func (s *UserService) ChangeUserPassword(ctx context.Context, id, currentPassword, newPassword string) error {
	s.logger.Info("Changing user password", zap.String("id", id))
	
	if id == "" {
		return errors.New("user ID is required")
	}
	
	if currentPassword == "" || newPassword == "" {
		return errors.New("current password and new password are required")
	}
	
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}
	
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if user == nil {
		return errors.New("user not found")
	}
	
	if !user.CheckPassword(currentPassword) {
		return errors.New("current password is incorrect")
	}
	
	if err := user.UpdatePassword(newPassword); err != nil {
		s.logger.Error("Failed to update password", zap.Error(err))
		return errors.New("failed to update password")
	}
	
	return s.userRepo.Update(ctx, user)
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(ctx context.Context, id string) error {
	s.logger.Info("Deactivating user", zap.String("id", id))
	
	if id == "" {
		return errors.New("user ID is required")
	}
	
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if user == nil {
		return errors.New("user not found")
	}
	
	user.Deactivate()
	return s.userRepo.Update(ctx, user)
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(ctx context.Context, id string) error {
	s.logger.Info("Activating user", zap.String("id", id))
	
	if id == "" {
		return errors.New("user ID is required")
	}
	
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if user == nil {
		return errors.New("user not found")
	}
	
	user.Activate()
	return s.userRepo.Update(ctx, user)
}

// CreateUserAddress creates a new address for a user
func (s *UserService) CreateUserAddress(ctx context.Context, userID, name, street, city, state, postalCode, country, phone string, isDefault bool) (*domain.Address, error) {
	s.logger.Info("Creating user address",
		zap.String("user_id", userID),
		zap.String("city", city),
		zap.String("country", country),
		zap.Bool("is_default", isDefault),
	)
	
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	// Validate required fields
	if name == "" || street == "" || city == "" || postalCode == "" || country == "" {
		return nil, errors.New("name, street, city, postal code, and country are required")
	}
	
	// Check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, errors.New("user not found")
	}
	
	address := domain.NewAddress(userID, name, street, city, state, postalCode, country, phone, isDefault)
	if err := s.addressRepo.Create(ctx, address); err != nil {
		return nil, err
	}
	
	return address, nil
}

// GetUserAddresses retrieves all addresses for a user
func (s *UserService) GetUserAddresses(ctx context.Context, userID string) ([]*domain.Address, error) {
	s.logger.Debug("Getting user addresses", zap.String("user_id", userID))
	
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	return s.addressRepo.GetByUserID(ctx, userID)
}

// GetUserDefaultAddress retrieves the default address for a user
func (s *UserService) GetUserDefaultAddress(ctx context.Context, userID string) (*domain.Address, error) {
	s.logger.Debug("Getting user default address", zap.String("user_id", userID))
	
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	
	return s.addressRepo.GetDefaultByUserID(ctx, userID)
}

// UpdateUserAddress updates a user address
func (s *UserService) UpdateUserAddress(ctx context.Context, addressID, userID, name, street, city, state, postalCode, country, phone string, isDefault bool) error {
	s.logger.Info("Updating user address",
		zap.String("address_id", addressID),
		zap.String("user_id", userID),
	)
	
	if addressID == "" || userID == "" {
		return errors.New("address ID and user ID are required")
	}
	
	// Validate required fields
	if name == "" || street == "" || city == "" || postalCode == "" || country == "" {
		return errors.New("name, street, city, postal code, and country are required")
	}
	
	// Check if address exists and belongs to the user
	address, err := s.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return err
	}
	
	if address == nil {
		return errors.New("address not found")
	}
	
	if address.UserID != userID {
		return errors.New("address does not belong to the user")
	}
	
	address.Update(name, street, city, state, postalCode, country, phone)
	address.SetDefault(isDefault)
	
	return s.addressRepo.Update(ctx, address)
}

// DeleteUserAddress deletes a user address
func (s *UserService) DeleteUserAddress(ctx context.Context, addressID, userID string) error {
	s.logger.Info("Deleting user address",
		zap.String("address_id", addressID),
		zap.String("user_id", userID),
	)
	
	if addressID == "" || userID == "" {
		return errors.New("address ID and user ID are required")
	}
	
	// Check if address exists and belongs to the user
	address, err := s.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return err
	}
	
	if address == nil {
		return errors.New("address not found")
	}
	
	if address.UserID != userID {
		return errors.New("address does not belong to the user")
	}
	
	return s.addressRepo.Delete(ctx, addressID)
}

// SetDefaultUserAddress sets an address as the default for a user
func (s *UserService) SetDefaultUserAddress(ctx context.Context, addressID, userID string) error {
	s.logger.Info("Setting default user address",
		zap.String("address_id", addressID),
		zap.String("user_id", userID),
	)
	
	if addressID == "" || userID == "" {
		return errors.New("address ID and user ID are required")
	}
	
	// Check if address exists and belongs to the user
	address, err := s.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return err
	}
	
	if address == nil {
		return errors.New("address not found")
	}
	
	if address.UserID != userID {
		return errors.New("address does not belong to the user")
	}
	
	return s.addressRepo.SetDefaultAddress(ctx, userID, addressID)
}

// CreateAdminUser creates a new admin user
func (s *UserService) CreateAdminUser(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	s.logger.Info("Creating admin user",
		zap.String("email", email),
		zap.String("first_name", firstName),
		zap.String("last_name", lastName),
	)

	// Validate input
	if email == "" || password == "" || firstName == "" || lastName == "" {
		return nil, errors.New("all fields are required")
	}

	// Check if email is already in use
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email is already in use")
	}

	// Create new user with admin role
	user, err := domain.NewUser(email, password, firstName, lastName, domain.RoleAdmin)
	if err != nil {
		s.logger.Error("Failed to create admin user", zap.Error(err))
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers lists all users with optional filtering and pagination
func (s *UserService) ListUsers(ctx context.Context, role string, active *bool, limit, offset int) ([]*domain.User, error) {
	s.logger.Debug("Listing users",
		zap.String("role", role),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	if limit <= 0 {
		limit = 10 // Default limit
	}
	
	filter := make(map[string]interface{})
	
	// Apply role filter if provided
	if role != "" {
		filter["role"] = strings.ToUpper(role)
	}
	
	// Apply active filter if provided
	if active != nil {
		filter["active"] = *active
	}
	
	return s.userRepo.List(ctx, filter, limit, offset)
}
