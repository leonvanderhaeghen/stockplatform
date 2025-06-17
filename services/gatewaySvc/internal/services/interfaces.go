package services

import (
	"context"
	"time"

	supplierv1 "github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/api/gen/go/proto/supplier/v1"
)

// ProductService defines the interface for product operations
type ProductService interface {
	// List products with filtering options
	ListProducts(ctx context.Context, categoryID, query string, active bool, limit, offset int, sortBy string, ascending bool) (interface{}, error)
	
	// List all product categories
	ListCategories(ctx context.Context) (interface{}, error)
	
	// Create a new product category
	CreateCategory(ctx context.Context, name, description, parentID string, isActive bool) (interface{}, error)
	
	// Get a product by ID
	GetProductByID(ctx context.Context, id string) (interface{}, error)
	
	// Create a new product
	CreateProduct(
		ctx context.Context,
		name, description string,
		costPrice, sellingPrice string,
		currency, sku, barcode string,
		categoryIDs []string,
		supplierID string,
		isActive, inStock bool,
		stockQty, lowStockAt int32,
		imageURLs, videoURLs []string,
		metadata map[string]string,
	) (interface{}, error)
	
	// Update an existing product
	// Note: This is not fully implemented in the gRPC service
	UpdateProduct(
		ctx context.Context,
		id, name, description, sku string,
		categories []string,
		price, cost string,
		active bool,
		images []string,
		attributes map[string]string,
	) error
	
	// Delete a product
	// Note: This is not implemented in the gRPC service
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
	
	// POS-related inventory operations
	
	// PerformPOSInventoryCheck checks inventory availability for POS
	PerformPOSInventoryCheck(ctx context.Context, locationID string, items []map[string]interface{}) (interface{}, error)
	
	// ReserveForPOSTransaction reserves inventory for POS transactions
	ReserveForPOSTransaction(ctx context.Context, locationID string, orderID string, items []map[string]interface{}) (interface{}, error)
	
	// CompletePickup marks a pickup as complete
	CompletePickup(ctx context.Context, reservationID string, staffID string, notes string) (interface{}, error)
	
	// DeductForDirectPOSTransaction directly deducts inventory for POS sales
	DeductForDirectPOSTransaction(ctx context.Context, locationID string, staffID string, items []map[string]interface{}, reason string) (interface{}, error)
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
	
	// POS-related order operations
	
	// CreatePOSOrder creates an order from a POS terminal
	CreatePOSOrder(ctx context.Context, userID string, items []map[string]interface{}, locationID, staffID, paymentType string, paymentData map[string]string, notes string) (interface{}, error)
	
	// ProcessQuickPOSTransaction creates and processes a POS order in one step
	ProcessQuickPOSTransaction(ctx context.Context, locationID, staffID string, items []map[string]interface{}, paymentInfo map[string]interface{}) (interface{}, error)
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

// SupplierService defines the interface for supplier operations
type SupplierService interface {
	// Create a new supplier
	CreateSupplier(ctx context.Context, req *supplierv1.CreateSupplierRequest) (*supplierv1.Supplier, error)
	// Get a supplier by ID
	GetSupplier(ctx context.Context, id string) (*supplierv1.Supplier, error)
	// Update an existing supplier
	UpdateSupplier(ctx context.Context, req *supplierv1.UpdateSupplierRequest) (*supplierv1.Supplier, error)
	// Delete a supplier
	DeleteSupplier(ctx context.Context, id string) error
	// List suppliers with pagination and search
	ListSuppliers(ctx context.Context, page, pageSize int32, search string) ([]*supplierv1.Supplier, int32, error)
	// Close closes the connection to the supplier service
	Close() error
	
	// ListAdapters returns all available supplier adapters
	ListAdapters(ctx context.Context) ([]*supplierv1.SupplierAdapter, error)
	// GetAdapterCapabilities returns the capabilities of a specific adapter
	GetAdapterCapabilities(ctx context.Context, adapterName string) (*supplierv1.AdapterCapabilities, error)
	// TestAdapterConnection tests the connection to a supplier's system using the specified adapter
	TestAdapterConnection(ctx context.Context, adapterName string, config map[string]string) error
	// SyncProducts synchronizes products from a supplier using their configured adapter
	SyncProducts(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error)
	// SyncInventory synchronizes inventory from a supplier using their configured adapter
	SyncInventory(ctx context.Context, supplierID string, options *supplierv1.SyncOptions) (string, error)
}

// POSService defines the interface for point-of-sale operations
type POSService interface {
	// ProcessTransaction processes a point-of-sale transaction (sale, return, exchange)
	ProcessTransaction(
		ctx context.Context,
		transactionType string, // "sale", "return", or "exchange"
		locationID string,
		staffID string,
		referenceOrderID string,
		items []map[string]interface{},
		paymentInfo map[string]interface{},
	) (interface{}, error)
	
	// CheckInventoryAvailability checks if items are available at the specified location
	CheckInventoryAvailability(
		ctx context.Context,
		locationID string,
		items []map[string]interface{},
	) (interface{}, error)
	
	// GetNearbyInventory finds nearby locations with available inventory
	GetNearbyInventory(
		ctx context.Context,
		productID string,
		sku string,
		quantity int32,
		lat, lng float64,
		radiusKm int32,
	) (interface{}, error)
	
	// ReserveForPickup reserves inventory for in-store pickup
	ReserveForPickup(
		ctx context.Context,
		userID string,
		locationID string,
		items []map[string]interface{},
		pickupTime time.Time,
		notes string,
	) (interface{}, error)
	
	// CompletePickup marks a pickup reservation as completed
	CompletePickup(
		ctx context.Context,
		reservationID string,
		staffID string,
		notes string,
	) (interface{}, error)
	
	// CancelPickup cancels a pickup reservation
	CancelPickup(
		ctx context.Context,
		reservationID string,
		reason string,
	) (interface{}, error)
}
