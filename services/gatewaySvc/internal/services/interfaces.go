package services

import (
	"context"
	"time"
)

// ProductService defines the interface for product operations
type ProductService interface {
	// List products with filtering options
	ListProducts(ctx context.Context, categoryID, query string, active bool, limit, offset int, sortBy string, ascending bool) (interface{}, error)
	
	// List all product categories
	ListCategories(ctx context.Context) (interface{}, error)
	
	// Get a product by ID
	GetProductByID(ctx context.Context, id string) (interface{}, error)
	
	// Create a new product
	CreateProduct(ctx context.Context, name, description, sku string, categories []string, price, cost float64, active bool, images []string, attributes map[string]string) (interface{}, error)
	
	// Update an existing product
	UpdateProduct(ctx context.Context, id, name, description, sku string, categories []string, price, cost float64, active bool, images []string, attributes map[string]string) error
	
	// Delete a product
	DeleteProduct(ctx context.Context, id string) error
}

// InventoryService defines the interface for inventory operations
type InventoryService interface {
	// List inventory items with filtering options
	ListInventory(ctx context.Context, location string, lowStock bool, limit, offset int) (interface{}, error)
	
	// Get an inventory item by ID
	GetInventoryItemByID(ctx context.Context, id string) (interface{}, error)
	
	// Get inventory items by product ID
	GetInventoryItemsByProductID(ctx context.Context, productID string) (interface{}, error)
	
	// Get an inventory item by SKU
	GetInventoryItemBySKU(ctx context.Context, sku string) (interface{}, error)
	
	// Create a new inventory item
	CreateInventoryItem(ctx context.Context, productID, sku string, quantity int32, location string, reorderAt, reorderQty int32, cost float64) (interface{}, error)
	
	// Update an existing inventory item
	UpdateInventoryItem(ctx context.Context, id, productID, sku string, quantity int32, location string, reorderAt, reorderQty int32, cost float64) error
	
	// Delete an inventory item
	DeleteInventoryItem(ctx context.Context, id string) error
	
	// Add stock to an inventory item
	AddStock(ctx context.Context, id string, quantity int32, reason, reference string) (interface{}, error)
	
	// Remove stock from an inventory item
	RemoveStock(ctx context.Context, id string, quantity int32, reason, reference string) (interface{}, error)
}

// OrderService defines the interface for order operations
type OrderService interface {
	// Get orders for a specific user
	GetUserOrders(ctx context.Context, userID, status, startDate, endDate string, limit, offset int) (interface{}, error)
	
	// Get a specific order for a user
	GetUserOrder(ctx context.Context, orderID, userID string) (interface{}, error)
	
	// Create a new order
	CreateOrder(ctx context.Context, userID string, items []map[string]interface{}, addressID, paymentType string, paymentData map[string]string, shippingType, notes string) (interface{}, error)
	
	// List all orders (admin/staff)
	ListOrders(ctx context.Context, status, userID, startDate, endDate string, limit, offset int) (interface{}, error)
	
	// Get an order by ID (admin/staff)
	GetOrderByID(ctx context.Context, orderID string) (interface{}, error)
	
	// Update order status (admin/staff)
	UpdateOrderStatus(ctx context.Context, orderID, status, description string) error
	
	// Add payment to an order (admin/staff)
	AddOrderPayment(ctx context.Context, orderID string, amount float64, paymentType, reference, status string, date time.Time, description string, metadata map[string]string) error
	
	// Add tracking info to an order (admin/staff)
	AddOrderTracking(ctx context.Context, orderID, carrier, trackingNum string, shipDate, estDelivery time.Time, notes string) error
	
	// Cancel an order (admin/staff)
	CancelOrder(ctx context.Context, orderID, reason string) error
}

// UserService defines the interface for user operations
type UserService interface {
	// Register a new user
	RegisterUser(ctx context.Context, email, password, firstName, lastName, role string) (interface{}, error)
	
	// Authenticate a user
	AuthenticateUser(ctx context.Context, email, password string) (interface{}, error)
	
	// Get a user by ID
	GetUserByID(ctx context.Context, userID string) (interface{}, error)
	
	// Update user profile
	UpdateUserProfile(ctx context.Context, userID, firstName, lastName, phone string) error
	
	// Change user password
	ChangeUserPassword(ctx context.Context, userID, currentPassword, newPassword string) error
	
	// Get addresses for a user
	GetUserAddresses(ctx context.Context, userID string) (interface{}, error)
	
	// Create a new address for a user
	CreateUserAddress(ctx context.Context, userID, name, street, city, state, postalCode, country, phone string, isDefault bool) (interface{}, error)
	
	// Get default address for a user
	GetUserDefaultAddress(ctx context.Context, userID string) (interface{}, error)
	
	// Update an address for a user
	UpdateUserAddress(ctx context.Context, addressID, userID, name, street, city, state, postalCode, country, phone string, isDefault bool) error
	
	// Delete an address for a user
	DeleteUserAddress(ctx context.Context, addressID, userID string) error
	
	// Set an address as default for a user
	SetDefaultUserAddress(ctx context.Context, addressID, userID string) error
	
	// List all users (admin only)
	ListUsers(ctx context.Context, role string, active *bool, limit, offset int) (interface{}, error)
	
	// Activate a user (admin only)
	ActivateUser(ctx context.Context, userID string) error
	
	// Deactivate a user (admin only)
	DeactivateUser(ctx context.Context, userID string) error
}
