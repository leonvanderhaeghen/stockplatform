package order

import (
	"time"
	
	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	orderv1 "github.com/leonvanderhaeghen/stockplatform/services/orderSvc/api/gen/go/proto/order/v1"
)

// convertToOrder converts protobuf Order to domain Order
func (c *Client) convertToOrder(proto *orderv1.Order) *models.Order {
	if proto == nil {
		return nil
	}

	order := &models.Order{
		ID:          proto.Id,
		CustomerID:  proto.UserId,
		Status:      convertOrderStatusFromProto(proto.Status),
		TotalAmount: proto.TotalAmount,
	}

	// Convert order items
	order.Items = make([]*models.OrderItem, len(proto.Items))
	for i, protoItem := range proto.Items {
		order.Items[i] = c.convertToOrderItem(protoItem)
	}

	// Handle timestamps (protobuf uses strings)
	if proto.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, proto.CreatedAt); err == nil {
			order.CreatedAt = t
		}
	}
	if proto.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, proto.UpdatedAt); err == nil {
			order.UpdatedAt = t
		}
	}

	// Handle shipping address if present
	if proto.ShippingAddress != nil {
		order.ShippingAddress = c.convertToAddress(proto.ShippingAddress)
	}

	return order
}

// convertToOrderItem converts protobuf OrderItem to domain OrderItem
func (c *Client) convertToOrderItem(proto *orderv1.OrderItem) *models.OrderItem {
	if proto == nil {
		return nil
	}

	return &models.OrderItem{
		ID:        "", // protobuf doesn't have ID field for OrderItem
		ProductID: proto.ProductId,
		SKU:       proto.ProductSku,
		Quantity:  proto.Quantity,
		Price:     proto.Price,
		Total:     proto.Subtotal,
	}
}

// convertToAddress converts protobuf Address to domain Address
func (c *Client) convertToAddress(proto *orderv1.Address) *models.Address {
	if proto == nil {
		return nil
	}

	return &models.Address{
		Street:  proto.Street,
		City:    proto.City,
		State:   proto.State,
		ZipCode: proto.PostalCode,
		Country: proto.Country,
	}
}

// convertFromOrder converts domain Order to protobuf Order
func (c *Client) convertFromOrder(order *models.Order) *orderv1.Order {
	if order == nil {
		return nil
	}

	proto := &orderv1.Order{
		Id:          order.ID,
		UserId:      order.CustomerID,
		Status:      convertOrderStatusToProto(order.Status),
		TotalAmount: order.TotalAmount,
	}

	// Convert order items
	proto.Items = make([]*orderv1.OrderItem, len(order.Items))
	for i, item := range order.Items {
		proto.Items[i] = c.convertFromOrderItem(item)
	}

	// Handle timestamps (convert to strings)
	if !order.CreatedAt.IsZero() {
		proto.CreatedAt = order.CreatedAt.Format(time.RFC3339)
	}
	if !order.UpdatedAt.IsZero() {
		proto.UpdatedAt = order.UpdatedAt.Format(time.RFC3339)
	}

	// Handle shipping address if present
	if order.ShippingAddress != nil {
		proto.ShippingAddress = c.convertFromAddress(order.ShippingAddress)
	}

	return proto
}

// convertFromOrderItem converts domain OrderItem to protobuf OrderItem
func (c *Client) convertFromOrderItem(item *models.OrderItem) *orderv1.OrderItem {
	if item == nil {
		return nil
	}

	return &orderv1.OrderItem{
		ProductId:  item.ProductID,
		ProductSku: item.SKU,
		Quantity:   item.Quantity,
		Price:      item.Price,
		Subtotal:   item.Total,
	}
}

// convertFromAddress converts domain Address to protobuf Address
func (c *Client) convertFromAddress(addr *models.Address) *orderv1.Address {
	if addr == nil {
		return nil
	}

	return &orderv1.Address{
		Street:     addr.Street,
		City:       addr.City,
		State:      addr.State,
		PostalCode: addr.ZipCode,
		Country:    addr.Country,
	}
}

// convertToCreateOrderResponse converts protobuf CreateOrderResponse to domain CreateOrderResponse
func (c *Client) convertToCreateOrderResponse(proto *orderv1.CreateOrderResponse) *models.CreateOrderResponse {
	if proto == nil {
		return nil
	}

	return &models.CreateOrderResponse{
		Order:   c.convertToOrder(proto.Order),
		Message: "Order created successfully",
	}
}

// convertToListOrdersResponse converts protobuf ListOrdersResponse to domain ListOrdersResponse
func (c *Client) convertToListOrdersResponse(proto *orderv1.ListOrdersResponse) *models.ListOrdersResponse {
	if proto == nil {
		return nil
	}

	orders := make([]*models.Order, len(proto.Orders))
	for i, protoOrder := range proto.Orders {
		orders[i] = c.convertToOrder(protoOrder)
	}

	return &models.ListOrdersResponse{
		Orders:     orders,
		TotalCount: int32(len(orders)), // protobuf doesn't have total_count field
	}
}

// convertToUpdateOrderResponse converts protobuf UpdateOrderResponse to domain UpdateOrderResponse
func (c *Client) convertToUpdateOrderResponse(proto *orderv1.UpdateOrderResponse) *models.UpdateOrderResponse {
	if proto == nil {
		return nil
	}

	return &models.UpdateOrderResponse{
		Order:   nil, // protobuf doesn't return order in response
		Message: "Order updated successfully",
	}
}

// Helper functions to convert between order status types
func convertOrderStatusFromProto(protoStatus orderv1.OrderStatus) models.OrderStatus {
	switch protoStatus {
	case orderv1.OrderStatus_ORDER_STATUS_CREATED:
		return models.OrderStatusPending
	case orderv1.OrderStatus_ORDER_STATUS_PENDING:
		return models.OrderStatusPending
	case orderv1.OrderStatus_ORDER_STATUS_PAID:
		return models.OrderStatusConfirmed
	case orderv1.OrderStatus_ORDER_STATUS_SHIPPED:
		return models.OrderStatusShipped
	case orderv1.OrderStatus_ORDER_STATUS_DELIVERED:
		return models.OrderStatusDelivered
	case orderv1.OrderStatus_ORDER_STATUS_CANCELLED:
		return models.OrderStatusCancelled
	default:
		return models.OrderStatusPending
	}
}

func convertOrderStatusToProto(status models.OrderStatus) orderv1.OrderStatus {
	switch status {
	case models.OrderStatusPending:
		return orderv1.OrderStatus_ORDER_STATUS_PENDING
	case models.OrderStatusConfirmed:
		return orderv1.OrderStatus_ORDER_STATUS_PAID
	case models.OrderStatusProcessing:
		return orderv1.OrderStatus_ORDER_STATUS_PENDING
	case models.OrderStatusShipped:
		return orderv1.OrderStatus_ORDER_STATUS_SHIPPED
	case models.OrderStatusDelivered:
		return orderv1.OrderStatus_ORDER_STATUS_DELIVERED
	case models.OrderStatusCancelled:
		return orderv1.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return orderv1.OrderStatus_ORDER_STATUS_PENDING
	}
}
