package inventory

import (
	"time"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
)

// convertToInventoryItem converts protobuf InventoryItem to domain InventoryItem
func (c *Client) convertToInventoryItem(proto *inventoryv1.InventoryItem) *models.InventoryItem {
	if proto == nil {
		return nil
	}

	// Calculate available quantity (total - reserved)
	available := proto.Quantity - proto.Reserved
	
	return &models.InventoryItem{
		ID:          proto.Id,
		ProductID:   proto.ProductId,
		SKU:         proto.Sku,
		Quantity:    proto.Quantity,
		Reserved:    proto.Reserved,
		LocationID:  proto.LocationId,
		Available:   available,
		ReorderAt:   proto.ReorderThreshold,
		ReorderQty:  proto.ReorderAmount,
		Cost:        0.0, // Cost not available in protobuf schema
		
		// Handle timestamp conversion from string
		CreatedAt: parseTimestamp(proto.CreatedAt),
		UpdatedAt: parseTimestamp(proto.LastUpdated),
	}
}

// convertFromInventoryItem converts domain InventoryItem to protobuf InventoryItem
func (c *Client) convertFromInventoryItem(item *models.InventoryItem) *inventoryv1.InventoryItem {
	if item == nil {
		return nil
	}

	return &inventoryv1.InventoryItem{
		Id:               item.ID,
		ProductId:        item.ProductID,
		Sku:              item.SKU,
		Quantity:         item.Quantity,
		Reserved:         item.Reserved,
		LocationId:       item.LocationID,
		ReorderThreshold: item.ReorderAt,
		ReorderAmount:    item.ReorderQty,
		
		// Handle timestamp conversion to string
		CreatedAt:   formatTimestamp(item.CreatedAt),
		LastUpdated: formatTimestamp(item.UpdatedAt),
	}
}

// convertToCheckAvailabilityResponse converts protobuf CheckAvailabilityResponse to domain CheckAvailabilityResponse
func (c *Client) convertToCheckAvailabilityResponse(proto *inventoryv1.CheckAvailabilityResponse) *models.CheckAvailabilityResponse {
	if proto == nil {
		return nil
	}

	response := &models.CheckAvailabilityResponse{
		Available: len(proto.Items) > 0, // Determine availability based on items
		Items:     make([]models.InventoryAvailabilityItem, len(proto.Items)),
	}

	for i, protoItem := range proto.Items {
		response.Items[i] = protoToInventoryAvailabilityItem(protoItem)
	}

	return response
}

func protoToInventoryAvailabilityItem(protoItem *inventoryv1.ItemAvailability) models.InventoryAvailabilityItem {
	return models.InventoryAvailabilityItem{
		ProductID:    protoItem.ProductId,
		SKU:          protoItem.Sku,
		RequestedQty: 0, // Not available in current protobuf schema
		AvailableQty: protoItem.AvailableQuantity,
		Available:    protoItem.InStock,
	}
}

func inventoryAvailabilityItemToProto(item models.InventoryAvailabilityItem) *inventoryv1.ItemAvailability {
	return &inventoryv1.ItemAvailability{
		ProductId:         item.ProductID,
		Sku:               item.SKU,
		AvailableQuantity: item.AvailableQty,
		InStock:           item.Available,
		Status:            determineStockStatus(item.AvailableQty, item.Available),
	}
}

func determineStockStatus(availableQty int32, isAvailable bool) string {
	if !isAvailable || availableQty == 0 {
		return "out_of_stock"
	}
	if availableQty < 10 { // Arbitrary low stock threshold
		return "low_stock"
	}
	return "in_stock"
}

// parseTimestamp converts string timestamp to time.Time
func parseTimestamp(timestamp string) time.Time {
	if timestamp == "" {
		return time.Time{}
	}
	// Try common timestamp formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
	}
	for _, format := range formats {
		if parsed, err := time.Parse(format, timestamp); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

// formatTimestamp converts time.Time to string
func formatTimestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
