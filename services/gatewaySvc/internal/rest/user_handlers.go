package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	userv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/user/v1"
)

// UserRegisterRequest represents the register request body
type UserRegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Role      string `json:"role" enums:"CUSTOMER,ADMIN,STAFF"`
}

// UserLoginRequest represents the login request body
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the profile update request body
type UpdateProfileRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Phone     string `json:"phone"`
}

// ChangePasswordRequest represents the password change request body
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// AddressRequest represents the address request body
type AddressRequest struct {
	Name       string `json:"name" binding:"required"`
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode" binding:"required"`
	Country    string `json:"country" binding:"required"`
	Phone      string `json:"phone"`
	IsDefault  bool   `json:"isDefault"`
}

// registerUser handles user registration
func (s *Server) registerUser(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Default to CUSTOMER role if not specified
	if req.Role == "" {
		req.Role = "CUSTOMER"
	}

	// Validate role
	validRoles := map[string]bool{
		"CUSTOMER": true,
		"ADMIN":    true,
		"STAFF":    true,
	}

	if !validRoles[req.Role] {
		respondWithError(c, http.StatusBadRequest, "Invalid role. Must be one of: CUSTOMER, ADMIN, STAFF")
		return
	}

	user, err := s.userSvc.RegisterUser(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName, req.Role)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "User registration")
		return
	}

	respondWithSuccess(c, http.StatusCreated, user)
}

// loginUser handles user authentication
func (s *Server) loginUser(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	result, err := s.userSvc.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		s.logger.Debug("Authentication failed", 
			zap.String("email", req.Email),
			zap.Error(err),
		)
		respondWithError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Transform the gRPC response to include string role
	resp, ok := result.(*userv1.AuthenticateUserResponse)
	if !ok {
		s.logger.Error("Unexpected response type from user service")
		respondWithError(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert the role enum to string
	var roleStr string
	switch resp.User.Role {
	case userv1.Role_ROLE_CUSTOMER:
		roleStr = "CUSTOMER"
	case userv1.Role_ROLE_ADMIN:
		roleStr = "ADMIN"
	case userv1.Role_ROLE_STAFF:
		roleStr = "STAFF"
	default:
		roleStr = "UNKNOWN"
	}

	// Create a new response with the string role
	transformedResp := map[string]interface{}{
		"token": resp.Token,
		"user": map[string]interface{}{
			"id":          resp.User.Id,
			"email":       resp.User.Email,
			"first_name":  resp.User.FirstName,
			"last_name":   resp.User.LastName,
			"role":        roleStr,
			"phone":       resp.User.Phone,
			"active":      resp.User.Active,
			"last_login":  resp.User.LastLogin,
			"created_at":  resp.User.CreatedAt,
			"updated_at":  resp.User.UpdatedAt,
		},
	}

	respondWithSuccess(c, http.StatusOK, transformedResp)
}

// getCurrentUser returns the current authenticated user
func (s *Server) getCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	user, err := s.userSvc.GetUserByID(c.Request.Context(), userIDStr)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get current user")
		return
	}

	respondWithSuccess(c, http.StatusOK, user)
}

// updateUserProfile updates the current user's profile
func (s *Server) updateUserProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := s.userSvc.UpdateUserProfile(c.Request.Context(), userIDStr, req.FirstName, req.LastName, req.Phone)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Update user profile")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// changeUserPassword changes the current user's password
func (s *Server) changeUserPassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := s.userSvc.ChangeUserPassword(c.Request.Context(), userIDStr, req.CurrentPassword, req.NewPassword)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Change user password")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// getUserAddresses returns all addresses for the current user
func (s *Server) getUserAddresses(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	addresses, err := s.userSvc.GetUserAddresses(c.Request.Context(), userIDStr)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get user addresses")
		return
	}

	respondWithSuccess(c, http.StatusOK, addresses)
}

// createUserAddress creates a new address for the current user
func (s *Server) createUserAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	address, err := s.userSvc.CreateUserAddress(
		c.Request.Context(),
		userIDStr,
		req.Name,
		req.Street,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.Phone,
		req.IsDefault,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Create user address")
		return
	}

	respondWithSuccess(c, http.StatusCreated, address)
}

// getUserDefaultAddress returns the default address for the current user
func (s *Server) getUserDefaultAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	address, err := s.userSvc.GetUserDefaultAddress(c.Request.Context(), userIDStr)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get user default address")
		return
	}

	if address == nil {
		respondWithSuccess(c, http.StatusOK, nil)
		return
	}

	respondWithSuccess(c, http.StatusOK, address)
}

// updateUserAddress updates an address for the current user
func (s *Server) updateUserAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		respondWithError(c, http.StatusBadRequest, "Address ID is required")
		return
	}

	var req AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := s.userSvc.UpdateUserAddress(
		c.Request.Context(),
		addressID,
		userIDStr,
		req.Name,
		req.Street,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
		req.Phone,
		req.IsDefault,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Update user address")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Address updated successfully"})
}

// deleteUserAddress deletes an address for the current user
func (s *Server) deleteUserAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		respondWithError(c, http.StatusBadRequest, "Address ID is required")
		return
	}

	err := s.userSvc.DeleteUserAddress(c.Request.Context(), addressID, userIDStr)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Delete user address")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

// setDefaultUserAddress sets an address as the default for the current user
func (s *Server) setDefaultUserAddress(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	addressID := c.Param("id")
	if addressID == "" {
		respondWithError(c, http.StatusBadRequest, "Address ID is required")
		return
	}

	err := s.userSvc.SetDefaultUserAddress(c.Request.Context(), addressID, userIDStr)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Set default user address")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Default address set successfully"})
}

// listUsers returns all users (admin only)
func (s *Server) listUsers(c *gin.Context) {
	role := c.Query("role")
	activeStr := c.Query("active")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Parse parameters
	var active *bool
	if activeStr != "" {
		activeBool := activeStr == "true"
		active = &activeBool
	}

	limit, err := parseIntParam(limitStr, 10)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	offset, err := parseIntParam(offsetStr, 0)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	users, err := s.userSvc.ListUsers(c.Request.Context(), role, active, limit, offset)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "List users")
		return
	}

	respondWithSuccess(c, http.StatusOK, users)
}

// getUserByID returns a specific user by ID (admin only)
func (s *Server) getUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		respondWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	user, err := s.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get user by ID")
		return
	}

	respondWithSuccess(c, http.StatusOK, user)
}

// activateUser activates a user account (admin only)
func (s *Server) activateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		respondWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	err := s.userSvc.ActivateUser(c.Request.Context(), userID)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Activate user")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "User activated successfully"})
}

// deactivateUser deactivates a user account (admin only)
func (s *Server) deactivateUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		respondWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	err := s.userSvc.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Deactivate user")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "User deactivated successfully"})
}
