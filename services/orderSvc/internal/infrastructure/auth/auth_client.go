package auth

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	userpb "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
)

// AuthClient handles authentication and authorization with the user service
type AuthClient struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// UserClaims represents the authenticated user information
type UserClaims struct {
	UserID    string
	Email     string
	Role      string
	FirstName string
	LastName  string
}

// NewAuthClient creates a new authentication client
func NewAuthClient(userServiceAddr string, logger *zap.Logger) (*AuthClient, error) {
	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userpb.NewUserServiceClient(conn)

	return &AuthClient{
		client: client,
		conn:   conn,
		logger: logger.Named("auth_client"),
	}, nil
}

// Close closes the gRPC connection
func (a *AuthClient) Close() error {
	return a.conn.Close()
}

// ValidateToken validates a JWT token with the user service
func (a *AuthClient) ValidateToken(ctx context.Context, token string) (*UserClaims, error) {
	req := &userpb.ValidateTokenRequest{
		Token: token,
	}

	resp, err := a.client.ValidateToken(ctx, req)
	if err != nil {
		a.logger.Warn("Token validation failed", zap.Error(err))
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return &UserClaims{
		UserID:    resp.UserId,
		Email:     resp.Email,
		Role:      resp.Role,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
	}, nil
}

// CheckPermission checks if a user has a specific permission
func (a *AuthClient) CheckPermission(ctx context.Context, role, permission string) (bool, error) {
	req := &userpb.CheckPermissionRequest{
		Role:       role,
		Permission: permission,
	}

	resp, err := a.client.CheckPermission(ctx, req)
	if err != nil {
		a.logger.Warn("Permission check failed", zap.Error(err))
		return false, fmt.Errorf("permission check failed: %w", err)
	}

	return resp.HasPermission, nil
}

// ExtractTokenFromContext extracts JWT token from gRPC metadata
func ExtractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return "", fmt.Errorf("missing authorization header")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// AuthMiddleware provides authentication middleware for the order service
type AuthMiddleware struct {
	authClient *AuthClient
	logger     *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authClient *AuthClient, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
		logger:     logger.Named("auth_middleware"),
	}
}

// UnaryServerInterceptor returns a gRPC unary server interceptor for authentication
func (m *AuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for health checks
		if strings.Contains(info.FullMethod, "Health") {
			return handler(ctx, req)
		}

		// Extract token from metadata
		token, err := ExtractTokenFromContext(ctx)
		if err != nil {
			m.logger.Warn("Failed to extract token", 
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return nil, fmt.Errorf("authentication required: %w", err)
		}

		// Validate token with user service
		claims, err := m.authClient.ValidateToken(ctx, token)
		if err != nil {
			m.logger.Warn("Token validation failed", 
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return nil, fmt.Errorf("invalid token: %w", err)
		}

		// Add claims to context
		ctx = context.WithValue(ctx, "user_claims", claims)

		return handler(ctx, req)
	}
}

// RequirePermission returns a middleware that checks for specific permissions
func (m *AuthMiddleware) RequirePermission(permission string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims := GetUserClaimsFromContext(ctx)
		if claims == nil {
			return nil, fmt.Errorf("authentication required")
		}

		hasPermission, err := m.authClient.CheckPermission(ctx, claims.Role, permission)
		if err != nil {
			m.logger.Error("Permission check failed", zap.Error(err))
			return nil, fmt.Errorf("permission check failed: %w", err)
		}

		if !hasPermission {
			m.logger.Warn("Permission denied", 
				zap.String("method", info.FullMethod),
				zap.String("user_id", claims.UserID),
				zap.String("role", claims.Role),
				zap.String("required_permission", permission),
			)
			return nil, fmt.Errorf("insufficient permissions")
		}

		return handler(ctx, req)
	}
}

// GetUserClaimsFromContext retrieves user claims from the context
func GetUserClaimsFromContext(ctx context.Context) *UserClaims {
	claims, ok := ctx.Value("user_claims").(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}
