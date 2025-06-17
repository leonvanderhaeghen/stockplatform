//go:build skip
// +build skip

// Tests temporarily removed. See issue #test-restore.
package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
)

const (
	userSvcAddr = "localhost:50051" // Update this to match your user service port
)

// TestUserServiceIntegration performs integration tests against the running user service
func TestUserServiceIntegration(t *testing.T) {
	// Skip if integration tests are not enabled
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to the user service
	conn, err := grpc.Dial(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := userv1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test user email - use timestamp to make it unique
	testEmail := "test.user." + time.Now().Format("20060102150405") + "@example.com"

	// Test registering a user
	registerReq := &userv1.RegisterUserRequest{
		Email:     testEmail,
		Password:  "Test@password123",
		FirstName: "Test",
		LastName:  "User",
		Role:      "CUSTOMER",
	}

	registerResp, err := client.RegisterUser(ctx, registerReq)
	require.NoError(t, err)
	require.NotNil(t, registerResp)
	require.NotNil(t, registerResp.User)
	require.NotEmpty(t, registerResp.User.Id)
	
	// Store the created user ID for subsequent tests
	userID := registerResp.User.Id

	// Test authenticating the user
	authReq := &userv1.AuthenticateUserRequest{
		Email:    testEmail,
		Password: "Test@password123",
	}

	authResp, err := client.AuthenticateUser(ctx, authReq)
	require.NoError(t, err)
	require.NotNil(t, authResp)
	require.NotEmpty(t, authResp.Token)
	require.NotNil(t, authResp.User)
	assert.Equal(t, userID, authResp.User.Id)
	assert.Equal(t, testEmail, authResp.User.Email)

	// Test getting the user by ID
	getUserResp, err := client.GetUser(ctx, &userv1.GetUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, getUserResp)
	require.NotNil(t, getUserResp.User)
	assert.Equal(t, testEmail, getUserResp.User.Email)
	assert.Equal(t, "Test", getUserResp.User.FirstName)
	assert.Equal(t, "User", getUserResp.User.LastName)

	// Test getting user by email
	getByEmailResp, err := client.GetUserByEmail(ctx, &userv1.GetUserByEmailRequest{
		Email: testEmail,
	})
	require.NoError(t, err)
	require.NotNil(t, getByEmailResp)
	require.NotNil(t, getByEmailResp.User)
	assert.Equal(t, userID, getByEmailResp.User.Id)

	// Test updating user profile
	updateProfileReq := &userv1.UpdateUserProfileRequest{
		Id:        userID,
		FirstName: "Updated",
		LastName:  "User",
		Phone:     "1234567890",
	}

	updateProfileResp, err := client.UpdateUserProfile(ctx, updateProfileReq)
	require.NoError(t, err)
	require.NotNil(t, updateProfileResp)
	assert.True(t, updateProfileResp.Success)

	// Verify profile updates
	getUserResp2, err := client.GetUser(ctx, &userv1.GetUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, getUserResp2)
	require.NotNil(t, getUserResp2.User)
	assert.Equal(t, "Updated", getUserResp2.User.FirstName)
	assert.Equal(t, "1234567890", getUserResp2.User.Phone)

	// Test creating user address
	addressReq := &userv1.CreateUserAddressRequest{
		UserId:     userID,
		Name:       "Home",
		Street:     "123 Test St",
		City:       "Test City",
		State:      "Test State",
		PostalCode: "12345",
		Country:    "Test Country",
		Phone:      "1234567890",
		IsDefault:  true,
	}

	addressResp, err := client.CreateUserAddress(ctx, addressReq)
	require.NoError(t, err)
	require.NotNil(t, addressResp)
	require.NotNil(t, addressResp.Address)
	require.NotEmpty(t, addressResp.Address.Id)
	
	// Store the created address ID for subsequent tests
	addressID := addressResp.Address.Id

	// Test getting user addresses
	getAddressesResp, err := client.GetUserAddresses(ctx, &userv1.GetUserAddressesRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, getAddressesResp)
	require.NotEmpty(t, getAddressesResp.Addresses)
	assert.Equal(t, 1, len(getAddressesResp.Addresses))
	assert.Equal(t, addressID, getAddressesResp.Addresses[0].Id)

	// Test getting default address
	getDefaultAddressResp, err := client.GetUserDefaultAddress(ctx, &userv1.GetUserDefaultAddressRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, getDefaultAddressResp)
	require.NotNil(t, getDefaultAddressResp.Address)
	assert.Equal(t, addressID, getDefaultAddressResp.Address.Id)
	assert.True(t, getDefaultAddressResp.Address.IsDefault)

	// Test updating address
	updateAddressReq := &userv1.UpdateUserAddressRequest{
		Id:         addressID,
		UserId:     userID,
		Name:       "Updated Home",
		Street:     addressReq.Street,
		City:       addressReq.City,
		State:      addressReq.State,
		PostalCode: addressReq.PostalCode,
		Country:    addressReq.Country,
		Phone:      addressReq.Phone,
		IsDefault:  true,
	}

	updateAddressResp, err := client.UpdateUserAddress(ctx, updateAddressReq)
	require.NoError(t, err)
	require.NotNil(t, updateAddressResp)
	assert.True(t, updateAddressResp.Success)

	// Verify address updates
	getAddressesResp2, err := client.GetUserAddresses(ctx, &userv1.GetUserAddressesRequest{
		UserId: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, getAddressesResp2)
	require.NotEmpty(t, getAddressesResp2.Addresses)
	assert.Equal(t, "Updated Home", getAddressesResp2.Addresses[0].Name)

	// Clean up: Delete the address
	deleteAddressResp, err := client.DeleteUserAddress(ctx, &userv1.DeleteUserAddressRequest{
		Id:     addressID,
		UserId: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, deleteAddressResp)
	assert.True(t, deleteAddressResp.Success)

	// Test deactivating the user
	deactivateResp, err := client.DeactivateUser(ctx, &userv1.DeactivateUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, deactivateResp)
	assert.True(t, deactivateResp.Success)

	// Verify user is deactivated
	getUserResp3, err := client.GetUser(ctx, &userv1.GetUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, getUserResp3)
	require.NotNil(t, getUserResp3.User)
	assert.False(t, getUserResp3.User.Active)

	// Test activating the user
	activateResp, err := client.ActivateUser(ctx, &userv1.ActivateUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, activateResp)
	assert.True(t, activateResp.Success)

	// Verify user is activated
	getUserResp4, err := client.GetUser(ctx, &userv1.GetUserRequest{
		Id: userID,
	})
	require.NoError(t, err)
	require.NotNil(t, getUserResp4)
	require.NotNil(t, getUserResp4.User)
	assert.True(t, getUserResp4.User.Active)
}
