package grpcclient

import (
	"context"
	"fmt"
	"log"
	"time"

	inventorypb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/inventory/v1"
	orderpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/order/v1"
	productpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/product/v1"
	userpb "github.com/leonvanderhaeghen/stockplatform/pkg/gen/go/user/v1"
)

// ExampleService demonstrates how to use the gRPC clients
// This is just an example and should be adapted to your actual service needs
type ExampleService struct {
	productClient   *ProductClient
	inventoryClient *InventoryClient
	userClient     *UserClient
	orderClient    *OrderClient
}

// NewExampleService creates a new example service with the given client addresses
func NewExampleService(productAddr, inventoryAddr, userAddr, orderAddr string) (*ExampleService, error) {
	// Initialize product client
	productClient, err := NewProductClient(productAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create product client: %w", err)
	}

	// Initialize inventory client
	inventoryClient, err := NewInventoryClient(inventoryAddr)
	if err != nil {
		productClient.Close()
		return nil, fmt.Errorf("failed to create inventory client: %w", err)
	}

	// Initialize user client
	userClient, err := NewUserClient(userAddr)
	if err != nil {
		productClient.Close()
		inventoryClient.Close()
		return nil, fmt.Errorf("failed to create user client: %w", err)
	}

	// Initialize order client
	orderClient, err := NewOrderClient(orderAddr)
	if err != nil {
		productClient.Close()
		inventoryClient.Close()
		userClient.Close()
		return nil, fmt.Errorf("failed to create order client: %w", err)
	}

	return &ExampleService{
		productClient:   productClient,
		inventoryClient: inventoryClient,
		userClient:     userClient,
		orderClient:    orderClient,
	}, nil
}

// Close closes all client connections
func (s *ExampleService) Close() {
	s.productClient.Close()
	s.inventoryClient.Close()
	s.userClient.Close()
	s.orderClient.Close()
}

// CreateProductWithInventory creates a new product and its initial inventory
func (s *ExampleService) CreateProductWithInventory(ctx context.Context, name, description, costPrice, sellingPrice, currency, sku, barcode, supplierID string, categoryIDs []string, imageURLs []string, initialStock int32) (string, error) {
	// Create the product
	createResp, err := s.productClient.CreateProduct(
		ctx,
		name,
		description,
		costPrice,
		sellingPrice,
		currency,
		sku,
		barcode,
		supplierID,
		categoryIDs,
		true,  // isActive
		true,  // inStock
		initialStock, // stockQty
		10,    // lowStockAt
		imageURLs,
		nil,   // videoURLs
		nil,   // metadata
	)
	if err != nil {
		return "", fmt.Errorf("failed to create product: %w", err)
	}

	// Create inventory for the product
	_, err = s.inventoryClient.CreateInventory(ctx, &inventorypb.CreateInventoryRequest{
		ProductId: createResp.GetProduct().GetId(),
		Quantity:  initialStock,
		Sku:       sku,
	})
	if err != nil {
		// In a real application, you might want to clean up the product if inventory creation fails
		return "", fmt.Errorf("failed to create inventory: %w", err)
	}

	log.Printf("Created product %s with initial stock of %d", createResp.GetProduct().GetId(), initialStock)
	return createResp.GetProduct().GetId(), nil
}

// GetProductWithInventory retrieves a product along with its inventory information
func (s *ExampleService) GetProductWithInventory(ctx context.Context, productID string) (*productpb.Product, *inventorypb.InventoryItem, error) {
	// Get the product
	productResp, err := s.productClient.GetProduct(ctx, &productpb.GetProductRequest{Id: productID})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Get the inventory for this product
	inventoryResp, err := s.inventoryClient.GetInventoryByProductID(ctx, &inventorypb.GetInventoryByProductIDRequest{
		ProductId: productID,
	})
	if err != nil {
		// In a real application, you might want to handle the case where inventory doesn't exist
		return productResp.Product, nil, nil
	}

	return productResp.Product, inventoryResp.Inventory, nil
}

// RegisterAndCreateOrder demonstrates a complete flow of user registration and order creation
func (s *ExampleService) RegisterAndCreateOrder(ctx context.Context, userEmail, userPassword, productID string, quantity int32) (string, error) {
	// 1. Register a new user
	registerResp, err := s.userClient.RegisterUser(ctx, &userpb.RegisterUserRequest{
		Email:     userEmail,
		Password:  userPassword,
		FirstName: "John",
		LastName:  "Doe",
	})
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}
	userID := registerResp.GetUser().GetId()
	log.Printf("Registered user with ID: %s", userID)

	// 2. Create a user address
	addressResp, err := s.userClient.CreateUserAddress(ctx, &userpb.CreateUserAddressRequest{
		UserId:     userID,
		Name:       "Home",
		Street:     "123 Main St",
		City:       "New York",
		State:      "NY",
		PostalCode: "10001",
		Country:    "USA",
		IsDefault:  true,
		Phone:      "+1234567890",
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user address: %w", err)
	}
	_ = addressResp.GetAddress().GetId() // Address ID is not used in this example

	// 3. Get the default address
	defaultAddrResp, err := s.userClient.GetUserDefaultAddress(ctx, &userpb.GetUserDefaultAddressRequest{
		UserId: userID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get default address: %w", err)
	}
	defaultAddress := defaultAddrResp.GetAddress()

	// 4. Create an order
	orderResp, err := s.orderClient.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:  userID,
		Items: []*orderpb.OrderItem{
			{
				ProductId:  productID,
				Quantity:   quantity,
			},
		},
		ShippingAddress: &orderpb.Address{
			Street:     defaultAddress.GetStreet(),
			City:       defaultAddress.GetCity(),
			State:      defaultAddress.GetState(),
			PostalCode: defaultAddress.GetPostalCode(),
			Country:    defaultAddress.GetCountry(),
		},
		BillingAddress: &orderpb.Address{
			Street:     defaultAddress.GetStreet(),
			City:       defaultAddress.GetCity(),
			State:      defaultAddress.GetState(),
			PostalCode: defaultAddress.GetPostalCode(),
			Country:    defaultAddress.GetCountry(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	orderID := orderResp.GetOrder().GetId()
	log.Printf("Created order with ID: %s", orderID)

	// 5. Add payment to the order
	_, err = s.orderClient.AddPayment(ctx, &orderpb.AddPaymentRequest{
		OrderId:       orderID,
		Method:        "credit_card",
		TransactionId: fmt.Sprintf("txn_%d", time.Now().Unix()),
		Amount:        orderResp.GetOrder().GetTotalAmount(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to add payment: %w", err)
	}

	// 6. Update order status to paid
	_, err = s.orderClient.UpdateOrderStatus(ctx, &orderpb.UpdateOrderStatusRequest{
		Id:     orderID,
		Status: orderpb.OrderStatus_ORDER_STATUS_PAID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	log.Printf("Order %s has been paid successfully", orderID)
	return orderID, nil
}

// GetOrderDetails retrieves an order with all its details
func (s *ExampleService) GetOrderDetails(ctx context.Context, orderID string) (*orderpb.Order, error) {
	// Get the order
	orderResp, err := s.orderClient.GetOrder(ctx, &orderpb.GetOrderRequest{Id: orderID})
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return orderResp.GetOrder(), nil
}

// ListUserOrders retrieves all orders for a specific user
func (s *ExampleService) ListUserOrders(ctx context.Context, userID string, limit, offset int32) ([]*orderpb.Order, error) {
	// Get user's orders
	ordersResp, err := s.orderClient.GetUserOrders(ctx, &orderpb.GetUserOrdersRequest{
		UserId: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user orders: %w", err)
	}

	return ordersResp.GetOrders(), nil
}
