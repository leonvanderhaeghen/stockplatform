package grpc

import (
	"context"
	"sort"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetNearbyInventory finds inventory availability at nearby locations
func (s *InventoryServer) GetNearbyInventory(ctx context.Context, req *inventoryv1.GetNearbyInventoryRequest) (*inventoryv1.GetNearbyInventoryResponse, error) {
	logger := s.logger.With(
		zap.String("handler", "GetNearbyInventory"),
		zap.String("location_id", req.LocationId),
		zap.Int32("radius_km", req.RadiusKm),
	)
	logger.Info("Processing GetNearbyInventory request")

	// Validate request
	if req.LocationId == "" {
		return nil, status.Error(codes.InvalidArgument, "location_id is required")
	}

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one inventory item is required")
	}

	if req.RadiusKm <= 0 {
		return nil, status.Error(codes.InvalidArgument, "radius_km must be greater than 0")
	}

	// Default max locations if not specified
	maxLocations := int(req.MaxLocations)
	if maxLocations <= 0 {
		maxLocations = 10 // Default
	}

	// Get current location for distance calculation
	originLocation, err := s.locationService.GetLocation(ctx, req.LocationId)
	if err != nil {
		logger.Error("Failed to get origin location", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "origin location not found: %v", err)
	}

	// Get all locations
	allLocations, err := s.locationService.ListLocations(ctx, 0, 100, true) // Pagination limits applied, include all locations
	if err != nil {
		logger.Error("Failed to list locations", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list locations: %v", err)
	}

	// Calculate distances and filter nearby locations
	type locationWithDistance struct {
		Location *domain.StoreLocation
		Distance float64
	}

	nearbyLocations := make([]locationWithDistance, 0)
	for _, loc := range allLocations {
		// Skip the origin location
		if loc.ID == req.LocationId {
			continue
		}

		// Calculate distance between locations
		distance := calculateDistance(
			originLocation.Latitude, originLocation.Longitude,
			loc.Latitude, loc.Longitude,
		)

		// Include location if within radius
		if distance <= float64(req.RadiusKm) {
			nearbyLocations = append(nearbyLocations, locationWithDistance{
				Location: loc,
				Distance: distance,
			})
		}
	}

	// Sort locations by distance (closest first)
	sort.Slice(nearbyLocations, func(i, j int) bool {
		return nearbyLocations[i].Distance < nearbyLocations[j].Distance
	})

	// Limit to max locations
	if len(nearbyLocations) > maxLocations {
		nearbyLocations = nearbyLocations[:maxLocations]
	}

	// Check inventory at each nearby location
	response := &inventoryv1.GetNearbyInventoryResponse{
		Locations: make([]*inventoryv1.NearbyLocationInventory, 0, len(nearbyLocations)),
	}

	for _, locWithDist := range nearbyLocations {
		loc := locWithDist.Location
		
		// Create availability check for this location
		availabilityItems := make([]*inventoryv1.ItemAvailability, 0, len(req.Items))
		
		for _, item := range req.Items {
			// Create result structure
			result := &inventoryv1.ItemAvailability{
				ProductId: item.ProductId,
				Sku:       item.Sku,
				InStock:   false,
				Status:    "unknown",
			}
			
			// Check availability by product ID or SKU
			var inventoryItem *domain.InventoryItem
			var err error
			
			if item.ProductId != "" {
				inventoryItem, err = s.service.GetInventoryItemByProductAndLocation(ctx, item.ProductId, loc.ID)
			} else if item.Sku != "" {
				inventoryItem, err = s.service.GetInventoryItemBySKUAndLocation(ctx, item.Sku, loc.ID)
			}
			
			if err != nil || inventoryItem == nil {
				// Item not found at this location
				result.Status = "not_found"
				availabilityItems = append(availabilityItems, result)
				continue
			}
			
			// Calculate available quantity
			availableQty := inventoryItem.Quantity - inventoryItem.Reserved
			result.AvailableQuantity = availableQty
			result.InventoryItemId = inventoryItem.ID
			
			// Determine stock status
			if availableQty <= 0 {
				result.Status = "out_of_stock"
			} else if availableQty < inventoryItem.ReorderPoint {
				result.Status = "low_stock"
				result.InStock = true
			} else {
				result.Status = "in_stock"
				result.InStock = true
			}
			
			// Check if it meets requested quantity
			if item.Quantity > 0 && availableQty >= item.Quantity {
				result.InStock = true
			} else if item.Quantity > 0 {
				result.InStock = false
			}
			
			availabilityItems = append(availabilityItems, result)
		}
		
		// Add this location to the response
		response.Locations = append(response.Locations, &inventoryv1.NearbyLocationInventory{
			LocationId:   loc.ID,
			LocationName: loc.Name,
			DistanceKm:   locWithDist.Distance,
			Items:        availabilityItems,
		})
	}

	logger.Info("GetNearbyInventory request completed",
		zap.Int("nearby_locations_count", len(response.Locations)))
	
	return response, nil
}

// calculateDistance calculates the distance between two points using the Haversine formula
// Returns distance in kilometers
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Implementation of the Haversine formula would go here
	// For now, we'll use a simple approximation based on latitude/longitude differences
	// This is just a placeholder - in a real implementation, use a proper geospatial library
	
	// Simple Euclidean distance scaled to approximate kilometers
	// This is very approximate and only works for small distances
	latDiff := lat2 - lat1
	lonDiff := lon2 - lon1
	
	// Very rough approximation (not for production use)
	// 1 degree of latitude ≈ 111 kilometers
	distance := (latDiff*latDiff + lonDiff*lonDiff) * 111.0
	
	return distance
}
