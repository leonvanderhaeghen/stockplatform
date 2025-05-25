package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Claims represents the JWT claims
type Claims struct {
	UserID    string `json:"sub"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

// authMiddleware creates a middleware for JWT authentication
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondWithError(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.jwtSecret), nil
		})

		if err != nil {
			s.logger.Debug("JWT validation failed", zap.Error(err))
			respondWithError(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Extract the claims
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Store the claims in the context for later use
			c.Set("userID", claims.UserID)
			c.Set("name", claims.Name)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			respondWithError(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}
	}
}

// adminMiddleware checks if the user has the admin role
func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "User role not found")
			c.Abort()
			return
		}

		if role != "ADMIN" {
			respondWithError(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// staffMiddleware checks if the user has the staff or admin role
func (s *Server) staffMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			respondWithError(c, http.StatusUnauthorized, "User role not found")
			c.Abort()
			return
		}

		if role != "ADMIN" && role != "STAFF" {
			respondWithError(c, http.StatusForbidden, "Staff access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
