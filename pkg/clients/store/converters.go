package store

import (
	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	storev1 "github.com/leonvanderhaeghen/stockplatform/services/storeSvc/api/gen/go/api/proto/store/v1"
)

// convertStoreFromProto converts a protobuf Store to a domain model Store
func convertStoreFromProto(protoStore *storev1.Store) *models.Store {
	if protoStore == nil {
		return nil
	}

	store := &models.Store{
		ID:          protoStore.GetId(),
		Name:        protoStore.GetName(),
		Description: protoStore.GetDescription(),
		Phone:       protoStore.GetPhone(),
		Email:       protoStore.GetEmail(),
		IsActive:    protoStore.GetIsActive(),
	}

	// Convert address
	if protoAddress := protoStore.GetAddress(); protoAddress != nil {
		store.Address = &models.Address{
			Street:     protoAddress.GetStreet(),
			City:       protoAddress.GetCity(),
			State:      protoAddress.GetState(),
			PostalCode: protoAddress.GetPostalCode(),
			Country:    protoAddress.GetCountry(),
			Latitude:   protoAddress.GetLatitude(),
			Longitude:  protoAddress.GetLongitude(),
		}
	}

	// Convert timestamps
	if createdAt := protoStore.GetCreatedAt(); createdAt != nil {
		store.CreatedAt = createdAt.AsTime()
	}
	if updatedAt := protoStore.GetUpdatedAt(); updatedAt != nil {
		store.UpdatedAt = updatedAt.AsTime()
	}

	return store
}

// convertStoreToProto converts a domain model Store to a protobuf Store
func convertStoreToProto(store *models.Store) *storev1.Store {
	if store == nil {
		return nil
	}

	protoStore := &storev1.Store{
		Id:          store.ID,
		Name:        store.Name,
		Description: store.Description,
		Phone:       store.Phone,
		Email:       store.Email,
		IsActive:    store.IsActive,
	}

	// Convert address
	if store.Address != nil {
		protoStore.Address = &storev1.Address{
			Street:     store.Address.Street,
			City:       store.Address.City,
			State:      store.Address.State,
			PostalCode: store.Address.PostalCode,
			Country:    store.Address.Country,
			Latitude:   store.Address.Latitude,
			Longitude:  store.Address.Longitude,
		}
	}

	return protoStore
}

// convertListStoresResponseFromProto converts a protobuf ListStoresResponse to a domain model response
func convertListStoresResponseFromProto(protoResponse *storev1.ListStoresResponse) *models.ListStoresResponse {
	if protoResponse == nil {
		return nil
	}

	// Calculate HasNextPage based on the number of stores returned and total count
	stores := make([]*models.Store, len(protoResponse.GetStores()))
	hasNextPage := len(protoResponse.GetStores()) > 0 && int32(len(protoResponse.GetStores())) < protoResponse.GetTotalCount()

	response := &models.ListStoresResponse{
		Stores:      stores,
		TotalCount:  protoResponse.GetTotalCount(),
		HasNextPage: hasNextPage,
	}

	for i, protoStore := range protoResponse.GetStores() {
		response.Stores[i] = convertStoreFromProto(protoStore)
	}

	return response
}
