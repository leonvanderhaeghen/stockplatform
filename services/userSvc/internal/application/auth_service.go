package application

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
)

// AuthServiceImpl implements the AuthService interface
type AuthServiceImpl struct {
	userRepo   domain.UserRepository
	jwtSecret  []byte
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	logger     *zap.Logger
	config     *AuthConfig
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	JWTSecret      []byte
	TokenDuration  time.Duration
	RefreshDuration time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo domain.UserRepository, config *AuthConfig, logger *zap.Logger) (*AuthServiceImpl, error) {
	// Generate RSA key pair for JWT signing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	service := &AuthServiceImpl{
		userRepo:   userRepo,
		jwtSecret:  config.JWTSecret,
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		logger:     logger.Named("auth_service"),
		config:     config,
	}

	return service, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	s.logger.Info("User login attempt", zap.String("email", email))

	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, domain.ErrInvalidCredentials
	}

	if user == nil {
		s.logger.Warn("Login attempt with non-existent email", zap.String("email", email))
		return nil, domain.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.Active {
		s.logger.Warn("Login attempt with inactive user", zap.String("email", email))
		return nil, domain.ErrUserNotActive
	}

	// Verify password
	if !user.CheckPassword(password) {
		s.logger.Warn("Invalid password attempt", zap.String("email", email))
		return nil, domain.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, err
	}

	// Record login
	user.RecordLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Warn("Failed to record login", zap.Error(err))
		// Don't fail the login, just log the warning
	}

	s.logger.Info("User login successful", 
		zap.String("email", email),
		zap.String("user_id", user.ID),
		zap.String("role", string(user.Role)),
	)

	return &domain.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure token's signature algorithm is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, domain.ErrInvalidToken
		}
		return s.publicKey, nil
	})

	if err != nil {
		s.logger.Warn("Token validation failed", zap.Error(err))
		return nil, domain.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		// Check if token is expired
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, domain.ErrTokenExpired
		}
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}

// RefreshToken refreshes an expired token
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, tokenString string) (*domain.AuthToken, error) {
	// Parse token without validation to extract claims
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	// Get user from database to ensure they're still active
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil || !user.Active {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate new token
	newToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

// HasPermission checks if a user has a specific permission
func (s *AuthServiceImpl) HasPermission(role domain.Role, permission domain.Permission) bool {
	return domain.HasPermission(role, permission)
}

// HasPermissions checks if a user has all specified permissions
func (s *AuthServiceImpl) HasPermissions(role domain.Role, permissions []domain.Permission) bool {
	for _, permission := range permissions {
		if !s.HasPermission(role, permission) {
			return false
		}
	}
	return true
}

// GetUserPermissions returns all permissions for a role
func (s *AuthServiceImpl) GetUserPermissions(role domain.Role) []domain.Permission {
	return domain.GetRolePermissions(role)
}

// Logout invalidates a token (placeholder for token blacklisting)
func (s *AuthServiceImpl) Logout(ctx context.Context, tokenString string) error {
	// In a production system, you would implement token blacklisting here
	// For now, we'll just log the logout event
	s.logger.Info("User logout", zap.String("token_prefix", tokenString[:min(len(tokenString), 10)]))
	return nil
}

// generateToken creates a new JWT token for a user
func (s *AuthServiceImpl) generateToken(user *domain.User) (*domain.AuthToken, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hour expiry

	claims := &domain.Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "stockplatform-userservice",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return nil, err
	}

	return &domain.AuthToken{
		Token:     tokenString,
		ExpiresAt: expirationTime,
		TokenType: "Bearer",
	}, nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
