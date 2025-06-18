package middleware

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/leonvanderhaeghen/stockplatform/services/userSvc/internal/domain"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	authService domain.AuthService
	logger      *zap.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService domain.AuthService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger.Named("auth_middleware"),
	}
}

// UnaryServerInterceptor returns a gRPC unary server interceptor for authentication
func (m *AuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for certain endpoints
		if m.shouldSkipAuth(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract token from metadata
		token, err := m.extractToken(ctx)
		if err != nil {
			m.logger.Warn("Failed to extract token", 
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// Validate token
		claims, err := m.authService.ValidateToken(ctx, token)
		if err != nil {
			m.logger.Warn("Token validation failed", 
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add claims to context
		ctx = m.addClaimsToContext(ctx, claims)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor for authentication
func (m *AuthMiddleware) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Skip authentication for certain endpoints
		if m.shouldSkipAuth(info.FullMethod) {
			return handler(srv, ss)
		}

		// Extract token from metadata
		token, err := m.extractToken(ss.Context())
		ctx := ss.Context()
		if err != nil {
			m.logger.Warn("Failed to extract token", 
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// Validate token
		claims, err := m.authService.ValidateToken(ctx, token)
		if err != nil {
			m.logger.Warn("Token validation failed", 
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add claims to context
		newCtx := m.addClaimsToContext(ss.Context(), claims)

		// Create new stream with updated context
		wrappedStream := &wrappedStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		return handler(srv, wrappedStream)
	}
}

// RequirePermission returns a middleware that checks for specific permissions
func (m *AuthMiddleware) RequirePermission(permission domain.Permission) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims := m.getClaimsFromContext(ctx)
		if claims == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		if !m.authService.HasPermission(claims.Role, permission) {
			m.logger.Warn("Permission denied", 
				zap.String("method", info.FullMethod),
				zap.String("user_id", claims.UserID),
				zap.String("role", string(claims.Role)),
				zap.String("required_permission", string(permission)),
			)
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		return handler(ctx, req)
	}
}

// RequirePermissions returns a middleware that checks for multiple permissions
func (m *AuthMiddleware) RequirePermissions(permissions ...domain.Permission) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims := m.getClaimsFromContext(ctx)
		if claims == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		allowed := true
		for _, perm := range permissions {
			if !m.authService.HasPermission(claims.Role, perm) {
				allowed = false
				break
			}
		}
		if !allowed {
			m.logger.Warn("Permissions denied", 
				zap.String("method", info.FullMethod),
				zap.String("user_id", claims.UserID),
				zap.String("role", string(claims.Role)),
				zap.Any("required_permissions", permissions),
			)
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		return handler(ctx, req)
	}
}

// RequireRole returns a middleware that checks for specific roles
func (m *AuthMiddleware) RequireRole(role domain.Role) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims := m.getClaimsFromContext(ctx)
		if claims == nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		if claims.Role != role && claims.Role != domain.RoleAdmin {
			m.logger.Warn("Role requirement not met", 
				zap.String("method", info.FullMethod),
				zap.String("user_id", claims.UserID),
				zap.String("user_role", string(claims.Role)),
				zap.String("required_role", string(role)),
			)
			return nil, status.Error(codes.PermissionDenied, "insufficient role")
		}

		return handler(ctx, req)
	}
}

// extractToken extracts JWT token from gRPC metadata
func (m *AuthMiddleware) extractToken(ctx context.Context) (string, error) {
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

// shouldSkipAuth determines if authentication should be skipped for certain endpoints
func (m *AuthMiddleware) shouldSkipAuth(method string) bool {
	// Skip auth for these methods
	skipMethods := []string{
		"/user.v1.UserService/Login",
		"/user.v1.UserService/Register",
		"/grpc.health.v1.Health/Check",
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo",
	}

	for _, skipMethod := range skipMethods {
		if method == skipMethod {
			return true
		}
	}

	return false
}

// addClaimsToContext adds JWT claims to the context
func (m *AuthMiddleware) addClaimsToContext(ctx context.Context, claims *domain.Claims) context.Context {
	return context.WithValue(ctx, "claims", claims)
}

// getClaimsFromContext retrieves JWT claims from the context
func (m *AuthMiddleware) getClaimsFromContext(ctx context.Context) *domain.Claims {
	claims, ok := ctx.Value("claims").(*domain.Claims)
	if !ok {
		return nil
	}
	return claims
}

// GetCurrentUser retrieves the current user from context
func GetCurrentUser(ctx context.Context) *domain.Claims {
	claims, ok := ctx.Value("claims").(*domain.Claims)
	if !ok {
		return nil
	}
	return claims
}

// wrappedStream wraps the server stream with a new context
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}
