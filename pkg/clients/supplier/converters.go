package supplier

import (
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

	suppliers := make([]*models.Supplier, len(proto.Suppliers))
	for i, protoSupplier := range proto.Suppliers {
		suppliers[i] = c.convertToSupplier(protoSupplier)
	}

	return &models.ListSuppliersResponse{
		Suppliers:  suppliers,
		TotalCount: int32(len(suppliers)), // protobuf doesn't have total_count field
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
