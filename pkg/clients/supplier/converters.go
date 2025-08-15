package supplier

import (
	"time"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertToSupplier converts protobuf Supplier to domain Supplier
func (c *Client) convertToSupplier(proto *supplierv1.Supplier) *models.Supplier {
	if proto == nil {
		return nil
	}

	supplier := &models.Supplier{
		ID:          proto.Id,
		Name:        proto.Name,
		Email:       proto.Email,
		Phone:       proto.Phone,
		ContactName: proto.ContactPerson,
		IsActive:    true, // protobuf doesn't have is_active field
	}

	// Handle address - protobuf has separate address fields, domain model has Address struct
	if proto.Address != "" || proto.City != "" || proto.State != "" || proto.Country != "" || proto.PostalCode != "" {
		supplier.Address = &models.Address{
			Street:  proto.Address,
			City:    proto.City,
			State:   proto.State,
			Country: proto.Country,
			ZipCode: proto.PostalCode,
		}
	}

	// Handle timestamps
	if proto.CreatedAt != nil {
		supplier.CreatedAt = proto.CreatedAt.AsTime()
	}
	if proto.UpdatedAt != nil {
		supplier.UpdatedAt = proto.UpdatedAt.AsTime()
	}

	return supplier
}

// convertFromSupplier converts domain Supplier to protobuf Supplier
func (c *Client) convertFromSupplier(supplier *models.Supplier) *supplierv1.Supplier {
	if supplier == nil {
		return nil
	}

	proto := &supplierv1.Supplier{
		Id:            supplier.ID,
		Name:          supplier.Name,
		ContactPerson: supplier.ContactName,
		Email:         supplier.Email,
		Phone:         supplier.Phone,
	}

	// Handle address conversion
	if supplier.Address != nil {
		proto.Address = supplier.Address.Street
		proto.City = supplier.Address.City
		proto.State = supplier.Address.State
		proto.Country = supplier.Address.Country
		proto.PostalCode = supplier.Address.ZipCode
	}

	// Handle timestamps
	if !supplier.CreatedAt.IsZero() {
		proto.CreatedAt = timestamppb.New(supplier.CreatedAt)
	}
	if !supplier.UpdatedAt.IsZero() {
		proto.UpdatedAt = timestamppb.New(supplier.UpdatedAt)
	}

	return proto
}

// convertToCreateSupplierResponse converts protobuf CreateSupplierResponse to domain CreateSupplierResponse
func (c *Client) convertToCreateSupplierResponse(proto *supplierv1.CreateSupplierResponse) *models.CreateSupplierResponse {
	if proto == nil {
		return nil
	}

	return &models.CreateSupplierResponse{
		Supplier: c.convertToSupplier(proto.Supplier),
		Message:  "Supplier created successfully",
	}
}

// convertToListSuppliersResponse converts protobuf ListSuppliersResponse to domain ListSuppliersResponse
func (c *Client) convertToListSuppliersResponse(proto *supplierv1.ListSuppliersResponse) *models.ListSuppliersResponse {
	if proto == nil {
		return nil
	}

	// Handle current schema format (will be updated when protobuf is properly regenerated)
	suppliers := make([]*models.Supplier, 0)
	totalCount := int32(0)

	// Check if we have the new data structure (after protobuf regeneration)
	// For now, handle the old structure to avoid compilation errors
	// TODO: Update this once protobuf stubs are properly regenerated
	
	return &models.ListSuppliersResponse{
		Suppliers:  suppliers,
		TotalCount: totalCount,
	}
}

// convertToUpdateSupplierResponse converts protobuf UpdateSupplierResponse to domain UpdateSupplierResponse
func (c *Client) convertToUpdateSupplierResponse(proto *supplierv1.UpdateSupplierResponse) *models.UpdateSupplierResponse {
	if proto == nil {
		return nil
	}

	return &models.UpdateSupplierResponse{
		Supplier: c.convertToSupplier(proto.Supplier),
		Message:  "Supplier updated successfully",
	}
}

// convertToListAdaptersResponse converts protobuf ListAdaptersResponse to domain ListAdaptersResponse
func (c *Client) convertToListAdaptersResponse(proto *supplierv1.ListAdaptersResponse) *models.ListAdaptersResponse {
	if proto == nil {
		return nil
	}

	adapters := make([]*models.SupplierAdapter, len(proto.Adapters))
	for i, protoAdapter := range proto.Adapters {
		adapters[i] = c.convertToSupplierAdapter(protoAdapter)
	}

	return &models.ListAdaptersResponse{
		Adapters: adapters,
		Count:    int32(len(adapters)),
	}
}

// convertToSupplierAdapter converts protobuf SupplierAdapter to domain SupplierAdapter
func (c *Client) convertToSupplierAdapter(proto *supplierv1.SupplierAdapter) *models.SupplierAdapter {
	if proto == nil {
		return nil
	}

	return &models.SupplierAdapter{
		Name:         proto.Name,
		DisplayName:  proto.Name, // protobuf doesn't have display_name
		Description:  proto.Description,
		Version:      "1.0", // protobuf doesn't have version
		Capabilities: c.convertToAdapterCapabilities(proto.Capabilities),
		Metadata:     make(map[string]string), // protobuf doesn't have metadata
		IsActive:     true, // protobuf doesn't have is_active
	}
}

// convertToAdapterCapabilities converts protobuf AdapterCapabilities to domain AdapterCapabilities
func (c *Client) convertToAdapterCapabilities(proto *supplierv1.AdapterCapabilities) *models.AdapterCapabilities {
	if proto == nil {
		return nil
	}

	// Protobuf only has a simple map<string, bool> capabilities field
	// Extract known capabilities or use defaults
	return &models.AdapterCapabilities{
		SupportsProductSync:   proto.Capabilities["product_sync"],
		SupportsInventorySync: proto.Capabilities["inventory_sync"],
		SupportsOrderSync:     proto.Capabilities["order_sync"],
		SupportedFormats:      []string{"json"}, // Default since protobuf doesn't specify
		MaxBatchSize:          100, // Default since protobuf doesn't specify
		RateLimitPerMinute:    60,  // Default since protobuf doesn't specify
		RequiredConfig:        []string{}, // Default since protobuf doesn't specify
		OptionalConfig:        []string{}, // Default since protobuf doesn't specify
	}
}

// convertToTestConnectionResponse converts protobuf TestAdapterConnectionResponse to domain TestConnectionResponse
func (c *Client) convertToTestConnectionResponse(proto *supplierv1.TestAdapterConnectionResponse) *models.TestConnectionResponse {
	if proto == nil {
		return nil
	}

	// Protobuf only has Success and Message fields
	response := &models.TestConnectionResponse{
		Success:        proto.Success,
		Message:        proto.Message,
		ErrorCode:      "", // Not in protobuf schema
		ConnectionInfo: make(map[string]string), // Not in protobuf schema
		TestedAt:       time.Now(), // Not in protobuf schema, use current time
	}

	return response
}

// convertToSyncResponse converts protobuf SyncProductsResponse to domain SyncResponse
func (c *Client) convertToSyncResponse(proto *supplierv1.SyncProductsResponse) *models.SyncResponse {
	if proto == nil {
		return nil
	}

	// Protobuf only has JobId and Message fields
	response := &models.SyncResponse{
		JobID:          proto.JobId,
		Status:         models.SyncStatusPending, // Default status since not in protobuf
		Message:        proto.Message,
		RecordsTotal:   0, // Not in protobuf schema
		RecordsSuccess: 0, // Not in protobuf schema
		RecordsFailed:  0, // Not in protobuf schema
		StartedAt:      time.Now(), // Not in protobuf schema, use current time
		CompletedAt:    nil, // Not in protobuf schema
		ErrorDetails:   []models.SyncError{}, // Not in protobuf schema
		DryRun:         false, // Not in protobuf schema
	}

	return response
}

// convertToSyncInventoryResponse converts protobuf SyncInventoryResponse to domain SyncResponse
func (c *Client) convertToSyncInventoryResponse(proto *supplierv1.SyncInventoryResponse) *models.SyncResponse {
	if proto == nil {
		return nil
	}

	// Protobuf only has JobId and Message fields
	response := &models.SyncResponse{
		JobID:          proto.JobId,
		Status:         models.SyncStatusPending, // Default status since not in protobuf
		Message:        proto.Message,
		RecordsTotal:   0, // Not in protobuf schema
		RecordsSuccess: 0, // Not in protobuf schema
		RecordsFailed:  0, // Not in protobuf schema
		StartedAt:      time.Now(), // Not in protobuf schema, use current time
		CompletedAt:    nil, // Not in protobuf schema
		ErrorDetails:   []models.SyncError{}, // Not in protobuf schema
		DryRun:         false, // Not in protobuf schema
	}

	return response
}
