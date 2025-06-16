package integration

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
	userv1 "github.com/leonvanderhaeghen/stockplatform/services/userSvc/api/gen/go/proto/user/v1"
	orderv1 "github.com/leonvanderhaeghen/stockplatform/services/orderSvc/api/gen/go/proto/order/v1"
)

// TestProductService tests basic product service functionality
func TestProductService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to connect to product service: %v", err)
	}
	defer conn.Close()

	client := productv1.NewProductServiceClient(conn)

	// Test creating a product
	productReq := &productv1.CreateProductRequest{
		Name:         "Test Product",
		Description:  "A test product",
		CostPrice:    "10.00",
		SellingPrice: "15.00",
		Currency:     "USD",
		Sku:          "TEST-001",
		CategoryIds:  []string{"test-category"},
		SupplierId:   "test-supplier",
		IsActive:     true,
		InStock:      true,
		StockQty:     100,
	}

	productResp, err := client.CreateProduct(ctx, productReq)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if productResp.Product == nil {
		t.Fatal("Created product is nil")
	}

	productID := productResp.Product.Id
	t.Logf("✅ Created product with ID: %s", productID)

	// Test getting the product
	getReq := &productv1.GetProductRequest{Id: productID}
	getResp, err := client.GetProduct(ctx, getReq)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if getResp.Product.Name != productReq.Name {
		t.Errorf("Expected product name %s, got %s", productReq.Name, getResp.Product.Name)
	}

	t.Logf("✅ Retrieved product: %s", getResp.Product.Name)

	// Test listing products
	listReq := &productv1.ListProductsRequest{
		Pagination: &productv1.Pagination{
			Page:     1,
			PageSize: 10,
		},
	}

	listResp, err := client.ListProducts(ctx, listReq)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	if len(listResp.Products) == 0 {
		t.Error("Expected at least one product in list")
	}

	t.Logf("✅ Listed %d products", len(listResp.Products))
}

// TestUserService tests basic user service functionality
func TestUserService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:50056",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to connect to user service: %v", err)
	}
	defer conn.Close()

	client := userv1.NewUserServiceClient(conn)

	// Test creating a user
	userReq := &userv1.RegisterUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	userResp, err := client.RegisterUser(ctx, userReq)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if userResp.User == nil {
		t.Fatal("Created user is nil")
	}

	userID := userResp.User.Id
	t.Logf("✅ Created user with ID: %s", userID)

	// Test getting the user
	getReq := &userv1.GetUserRequest{Id: userID}
	getResp, err := client.GetUser(ctx, getReq)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if getResp.User.Email != userReq.Email {
		t.Errorf("Expected user email %s, got %s", userReq.Email, getResp.User.Email)
	}

	t.Logf("✅ Retrieved user: %s", getResp.User.Email)
}

// TestOrderService tests basic order service functionality
func TestOrderService(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:50055",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("Failed to connect to order service: %v", err)
	}
	defer conn.Close()

	client := orderv1.NewOrderServiceClient(conn)

	// Test creating an order
	orderReq := &orderv1.CreateOrderRequest{
		UserId: "test-user-id",
		Items: []*orderv1.OrderItem{
			{
				ProductId: "test-product-id",
				Quantity:  2,
				Price:     "15.00",
			},
		},
		ShippingAddress: &orderv1.Address{
			Street:  "123 Test St",
			City:    "Test City",
			State:   "TS",
			ZipCode: "12345",
			Country: "US",
		},
		BillingAddress: &orderv1.Address{
			Street:  "123 Test St",
			City:    "Test City",
			State:   "TS",
			ZipCode: "12345",
			Country: "US",
		},
	}

	orderResp, err := client.CreateOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	if orderResp.Order == nil {
		t.Fatal("Created order is nil")
	}

	orderID := orderResp.Order.Id
	t.Logf("✅ Created order with ID: %s", orderID)

	// Test getting the order
	getReq := &orderv1.GetOrderRequest{Id: orderID}
	getResp, err := client.GetOrder(ctx, getReq)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}

	if getResp.Order.UserId != orderReq.UserId {
		t.Errorf("Expected order user ID %s, got %s", orderReq.UserId, getResp.Order.UserId)
	}

	t.Logf("✅ Retrieved order for user: %s", getResp.Order.UserId)
}
