package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "stockplatform/pkg/gen/user/v1"
	"stockplatform/services/userSvc/internal/application"
	"stockplatform/services/userSvc/internal/domain"
)

// UserServer implements the gRPC interface for user service
type UserServer struct {
	userv1.UnimplementedUserServiceServer
	service *application.UserService
	logger  *zap.Logger
}

// NewUserServer creates a new user gRPC server
func NewUserServer(service *application.UserService, logger *zap.Logger) userv1.UserServiceServer {
	return &UserServer{
		service: service,
		logger:  logger.Named("user_grpc_server"),
	}
}

// RegisterUser registers a new user
func (s *UserServer) RegisterUser(ctx context.Context, req *userv1.RegisterUserRequest) (*userv1.RegisterUserResponse, error) {
	s.logger.Info("gRPC RegisterUser called",
		zap.String("email", req.Email),
		zap.String("first_name", req.FirstName),
		zap.String("last_name", req.LastName),
	)

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first_name is required")
	}
	if req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "last_name is required")
	}

	user, err := s.service.RegisterUser(ctx, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		s.logger.Error("Failed to register user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to register user: "+err.Error())
	}

	return &userv1.RegisterUserResponse{
		User: toProtoUser(user),
	}, nil
}

// AuthenticateUser authenticates a user and returns a JWT token
func (s *UserServer) AuthenticateUser(ctx context.Context, req *userv1.AuthenticateUserRequest) (*userv1.AuthenticateUserResponse, error) {
	s.logger.Info("gRPC AuthenticateUser called", zap.String("email", req.Email))

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.service.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Error("Failed to authenticate user", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "authentication failed: "+err.Error())
	}

	// Get user details to include in response
	user, err := s.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Failed to get user details after authentication", zap.Error(err))
		// Still return the token even if we can't get user details
		return &userv1.AuthenticateUserResponse{
			Token: token,
		}, nil
	}

	return &userv1.AuthenticateUserResponse{
		Token: token,
		User:  toProtoUser(user),
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	s.logger.Debug("gRPC GetUser called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	user, err := s.service.GetUserByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get user", zap.Error(err))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userv1.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserServer) GetUserByEmail(ctx context.Context, req *userv1.GetUserByEmailRequest) (*userv1.GetUserResponse, error) {
	s.logger.Debug("gRPC GetUserByEmail called", zap.String("email", req.Email))

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	user, err := s.service.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Failed to get user by email", zap.Error(err))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userv1.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserServer) UpdateUserProfile(ctx context.Context, req *userv1.UpdateUserProfileRequest) (*userv1.UpdateUserProfileResponse, error) {
	s.logger.Info("gRPC UpdateUserProfile called",
		zap.String("id", req.Id),
		zap.String("first_name", req.FirstName),
		zap.String("last_name", req.LastName),
	)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first_name is required")
	}
	if req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "last_name is required")
	}

	if err := s.service.UpdateUserProfile(ctx, req.Id, req.FirstName, req.LastName, req.Phone); err != nil {
		s.logger.Error("Failed to update user profile", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update user profile: "+err.Error())
	}

	return &userv1.UpdateUserProfileResponse{
		Success: true,
	}, nil
}

// ChangeUserPassword changes a user's password
func (s *UserServer) ChangeUserPassword(ctx context.Context, req *userv1.ChangeUserPasswordRequest) (*userv1.ChangeUserPasswordResponse, error) {
	s.logger.Info("gRPC ChangeUserPassword called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.CurrentPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "current_password is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}

	if err := s.service.ChangeUserPassword(ctx, req.Id, req.CurrentPassword, req.NewPassword); err != nil {
		s.logger.Error("Failed to change user password", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to change user password: "+err.Error())
	}

	return &userv1.ChangeUserPasswordResponse{
		Success: true,
	}, nil
}

// DeactivateUser deactivates a user account
func (s *UserServer) DeactivateUser(ctx context.Context, req *userv1.DeactivateUserRequest) (*userv1.DeactivateUserResponse, error) {
	s.logger.Info("gRPC DeactivateUser called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.service.DeactivateUser(ctx, req.Id); err != nil {
		s.logger.Error("Failed to deactivate user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to deactivate user: "+err.Error())
	}

	return &userv1.DeactivateUserResponse{
		Success: true,
	}, nil
}

// ActivateUser activates a user account
func (s *UserServer) ActivateUser(ctx context.Context, req *userv1.ActivateUserRequest) (*userv1.ActivateUserResponse, error) {
	s.logger.Info("gRPC ActivateUser called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.service.ActivateUser(ctx, req.Id); err != nil {
		s.logger.Error("Failed to activate user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to activate user: "+err.Error())
	}

	return &userv1.ActivateUserResponse{
		Success: true,
	}, nil
}

// ListUsers lists all users with optional filtering and pagination
func (s *UserServer) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	s.logger.Debug("gRPC ListUsers called",
		zap.String("role", req.Role),
		zap.Bool("active", req.Active),
		zap.Int32("limit", req.Limit),
		zap.Int32("offset", req.Offset),
	)

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10 // Default limit
	}

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	active := &req.Active
	
	users, err := s.service.ListUsers(ctx, req.Role, active, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list users", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list users: "+err.Error())
	}

	// Convert domain users to proto users
	protoUsers := make([]*userv1.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, toProtoUser(user))
	}

	return &userv1.ListUsersResponse{
		Users: protoUsers,
	}, nil
}

// CreateUserAddress creates a new address for a user
func (s *UserServer) CreateUserAddress(ctx context.Context, req *userv1.CreateUserAddressRequest) (*userv1.CreateUserAddressResponse, error) {
	s.logger.Info("gRPC CreateUserAddress called",
		zap.String("user_id", req.UserId),
		zap.String("city", req.City),
		zap.String("country", req.Country),
	)

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Street == "" {
		return nil, status.Error(codes.InvalidArgument, "street is required")
	}
	if req.City == "" {
		return nil, status.Error(codes.InvalidArgument, "city is required")
	}
	if req.PostalCode == "" {
		return nil, status.Error(codes.InvalidArgument, "postal_code is required")
	}
	if req.Country == "" {
		return nil, status.Error(codes.InvalidArgument, "country is required")
	}

	address, err := s.service.CreateUserAddress(ctx, req.UserId, req.Name, req.Street, req.City, req.State, req.PostalCode, req.Country, req.Phone, req.IsDefault)
	if err != nil {
		s.logger.Error("Failed to create user address", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create user address: "+err.Error())
	}

	return &userv1.CreateUserAddressResponse{
		Address: toProtoAddress(address),
	}, nil
}

// GetUserAddresses retrieves all addresses for a user
func (s *UserServer) GetUserAddresses(ctx context.Context, req *userv1.GetUserAddressesRequest) (*userv1.GetUserAddressesResponse, error) {
	s.logger.Debug("gRPC GetUserAddresses called", zap.String("user_id", req.UserId))

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	addresses, err := s.service.GetUserAddresses(ctx, req.UserId)
	if err != nil {
		s.logger.Error("Failed to get user addresses", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get user addresses: "+err.Error())
	}

	// Convert domain addresses to proto addresses
	protoAddresses := make([]*userv1.Address, 0, len(addresses))
	for _, address := range addresses {
		protoAddresses = append(protoAddresses, toProtoAddress(address))
	}

	return &userv1.GetUserAddressesResponse{
		Addresses: protoAddresses,
	}, nil
}

// GetUserDefaultAddress retrieves the default address for a user
func (s *UserServer) GetUserDefaultAddress(ctx context.Context, req *userv1.GetUserDefaultAddressRequest) (*userv1.GetUserDefaultAddressResponse, error) {
	s.logger.Debug("gRPC GetUserDefaultAddress called", zap.String("user_id", req.UserId))

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	address, err := s.service.GetUserDefaultAddress(ctx, req.UserId)
	if err != nil {
		s.logger.Error("Failed to get user default address", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get user default address: "+err.Error())
	}

	if address == nil {
		return &userv1.GetUserDefaultAddressResponse{}, nil
	}

	return &userv1.GetUserDefaultAddressResponse{
		Address: toProtoAddress(address),
	}, nil
}

// UpdateUserAddress updates a user address
func (s *UserServer) UpdateUserAddress(ctx context.Context, req *userv1.UpdateUserAddressRequest) (*userv1.UpdateUserAddressResponse, error) {
	s.logger.Info("gRPC UpdateUserAddress called",
		zap.String("id", req.Id),
		zap.String("user_id", req.UserId),
	)

	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "id and user_id are required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Street == "" {
		return nil, status.Error(codes.InvalidArgument, "street is required")
	}
	if req.City == "" {
		return nil, status.Error(codes.InvalidArgument, "city is required")
	}
	if req.PostalCode == "" {
		return nil, status.Error(codes.InvalidArgument, "postal_code is required")
	}
	if req.Country == "" {
		return nil, status.Error(codes.InvalidArgument, "country is required")
	}

	if err := s.service.UpdateUserAddress(ctx, req.Id, req.UserId, req.Name, req.Street, req.City, req.State, req.PostalCode, req.Country, req.Phone, req.IsDefault); err != nil {
		s.logger.Error("Failed to update user address", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update user address: "+err.Error())
	}

	return &userv1.UpdateUserAddressResponse{
		Success: true,
	}, nil
}

// DeleteUserAddress deletes a user address
func (s *UserServer) DeleteUserAddress(ctx context.Context, req *userv1.DeleteUserAddressRequest) (*userv1.DeleteUserAddressResponse, error) {
	s.logger.Info("gRPC DeleteUserAddress called",
		zap.String("id", req.Id),
		zap.String("user_id", req.UserId),
	)

	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "id and user_id are required")
	}

	if err := s.service.DeleteUserAddress(ctx, req.Id, req.UserId); err != nil {
		s.logger.Error("Failed to delete user address", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete user address: "+err.Error())
	}

	return &userv1.DeleteUserAddressResponse{
		Success: true,
	}, nil
}

// SetDefaultUserAddress sets an address as the default for a user
func (s *UserServer) SetDefaultUserAddress(ctx context.Context, req *userv1.SetDefaultUserAddressRequest) (*userv1.SetDefaultUserAddressResponse, error) {
	s.logger.Info("gRPC SetDefaultUserAddress called",
		zap.String("id", req.Id),
		zap.String("user_id", req.UserId),
	)

	if req.Id == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "id and user_id are required")
	}

	if err := s.service.SetDefaultUserAddress(ctx, req.Id, req.UserId); err != nil {
		s.logger.Error("Failed to set default user address", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to set default user address: "+err.Error())
	}

	return &userv1.SetDefaultUserAddressResponse{
		Success: true,
	}, nil
}

// toProtoUser converts a domain user to a proto user
func toProtoUser(user *domain.User) *userv1.User {
	protoUser := &userv1.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	// Convert role
	switch user.Role {
	case domain.RoleCustomer:
		protoUser.Role = userv1.Role_ROLE_CUSTOMER
	case domain.RoleAdmin:
		protoUser.Role = userv1.Role_ROLE_ADMIN
	case domain.RoleStaff:
		protoUser.Role = userv1.Role_ROLE_STAFF
	default:
		protoUser.Role = userv1.Role_ROLE_UNSPECIFIED
	}

	// Set last login if it exists
	if !user.LastLogin.IsZero() {
		protoUser.LastLogin = user.LastLogin.Format(time.RFC3339)
	}

	return protoUser
}

// toProtoAddress converts a domain address to a proto address
func toProtoAddress(address *domain.Address) *userv1.Address {
	return &userv1.Address{
		Id:         address.ID,
		UserId:     address.UserID,
		Name:       address.Name,
		Street:     address.Street,
		City:       address.City,
		State:      address.State,
		PostalCode: address.PostalCode,
		Country:    address.Country,
		IsDefault:  address.IsDefault,
		Phone:      address.Phone,
		CreatedAt:  address.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  address.UpdatedAt.Format(time.RFC3339),
	}
}
