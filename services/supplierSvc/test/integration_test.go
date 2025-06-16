package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/supplier/v1"
)

const (
	supplierSvcAddr = "localhost:50053" // Update this to match your supplier service port
)

// TestSupplierServiceIntegration performs integration tests against the running supplier service
func TestSupplierServiceIntegration(t *testing.T) {
	// Skip if integration tests are not enabled
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to the supplier service
	conn, err := grpc.Dial(supplierSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := supplierv1.NewSupplierServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test creating a supplier
	supplier := &supplierv1.CreateSupplierRequest{
		Name:         "Test Supplier Integration",
		ContactPerson: "John Doe",
		Email:        "john.doe@testsupplier.com",
		Phone:        "1234567890",
		Address:      "123 Integration St",
		City:         "Test City",
		State:        "Test State",
		Country:      "Test Country",
		PostalCode:   "12345",
		Website:      "https://testsupplier.com",
		Currency:     "USD",
		LeadTimeDays: 5,
		PaymentTerms: "Net 30",
	}

	// Create the supplier
	createResp, err := client.CreateSupplier(ctx, supplier)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.NotNil(t, createResp.Supplier)
	require.NotEmpty(t, createResp.Supplier.Id)

	// Store the created supplier ID for subsequent tests
	supplierID := createResp.Supplier.Id

	// Test getting the created supplier
	getResp, err := client.GetSupplier(ctx, &supplierv1.GetSupplierRequest{
		Id: supplierID,
	})
	require.NoError(t, err)
	require.NotNil(t, getResp)
	require.NotNil(t, getResp.Supplier)
	assert.Equal(t, supplier.Name, getResp.Supplier.Name)
	assert.Equal(t, supplier.Email, getResp.Supplier.Email)

	// Test updating the supplier
	updateReq := &supplierv1.UpdateSupplierRequest{
		Id:           supplierID,
		Name:         "Updated Test Supplier",
		ContactPerson: supplier.ContactPerson,
		Email:        supplier.Email,
		Phone:        supplier.Phone,
		Address:      supplier.Address,
		City:         supplier.City,
		State:        supplier.State,
		Country:      supplier.Country,
		PostalCode:   supplier.PostalCode,
		Website:      supplier.Website,
		Currency:     supplier.Currency,
		LeadTimeDays: supplier.LeadTimeDays,
		PaymentTerms: supplier.PaymentTerms,
	}

	updateResp, err := client.UpdateSupplier(ctx, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.NotNil(t, updateResp.Supplier)
	assert.Equal(t, "Updated Test Supplier", updateResp.Supplier.Name)

	// Test listing suppliers
	listResp, err := client.ListSuppliers(ctx, &supplierv1.ListSuppliersRequest{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	assert.GreaterOrEqual(t, len(listResp.Suppliers), 1)
	found := false
	for _, s := range listResp.Suppliers {
		if s.Id == supplierID {
			found = true
			assert.Equal(t, "Updated Test Supplier", s.Name)
			break
		}
	}
	assert.True(t, found, "The created supplier should be in the list")

	// Test listing adapter capabilities
	// First, check if we have any adapters
	adaptersResp, err := client.ListAdapters(ctx, &supplierv1.ListAdaptersRequest{})
	require.NoError(t, err)
	
	if len(adaptersResp.Adapters) > 0 {
		// Test the first adapter's capabilities
		adapterName := adaptersResp.Adapters[0].Name
		capabilitiesResp, err := client.GetAdapterCapabilities(ctx, &supplierv1.GetAdapterCapabilitiesRequest{
			AdapterName: adapterName,
		})
		require.NoError(t, err)
		require.NotNil(t, capabilitiesResp)
		// At least some capabilities should be defined
		assert.NotNil(t, capabilitiesResp.Capabilities)
	}

	// Test deleting the supplier
	deleteResp, err := client.DeleteSupplier(ctx, &supplierv1.DeleteSupplierRequest{
		Id: supplierID,
	})
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
	assert.True(t, deleteResp.Success)

	// Verify the supplier is deleted
	_, err = client.GetSupplier(ctx, &supplierv1.GetSupplierRequest{
		Id: supplierID,
	})
	assert.Error(t, err, "Getting a deleted supplier should return an error")
}
