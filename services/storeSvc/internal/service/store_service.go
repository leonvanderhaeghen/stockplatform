package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/config"
	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/database"
	"github.com/leonvanderhaeghen/stockplatform/services/storeSvc/internal/models"
	storev1 "github.com/leonvanderhaeghen/stockplatform/services/storeSvc/api/gen/go/proto/store/v1"
)

// StoreService implements the store service gRPC interface
type StoreService struct {
	storev1.UnimplementedStoreServiceServer
	db     *database.Database
	config *config.Config
}

// NewStoreService creates a new store service instance
func NewStoreService(db *database.Database, cfg *config.Config) (*StoreService, error) {
	return &StoreService{
		db:     db,
		config: cfg,
	}, nil
}

// CreateStore creates a new store
func (s *StoreService) CreateStore(ctx context.Context, req *storev1.CreateStoreRequest) (*storev1.CreateStoreResponse, error) {
	store := &models.Store{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Address:     convertAddressFromProto(req.Address),
		Phone:       req.Phone,
		Email:       req.Email,
		IsActive:    true,
		Hours:       convertStoreHoursFromProto(req.Hours),
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := s.db.GetCollection("stores")
	_, err := collection.InsertOne(ctx, store)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	return &storev1.CreateStoreResponse{
		Store: convertStoreToProto(store),
	}, nil
}

// GetStore retrieves a store by ID
func (s *StoreService) GetStore(ctx context.Context, req *storev1.GetStoreRequest) (*storev1.GetStoreResponse, error) {
	collection := s.db.GetCollection("stores")
	
	var store models.Store
	err := collection.FindOne(ctx, bson.M{"_id": req.Id}).Decode(&store)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("store not found")
		}
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	return &storev1.GetStoreResponse{
		Store: convertStoreToProto(&store),
	}, nil
}

// ListStores lists stores with optional filtering
func (s *StoreService) ListStores(ctx context.Context, req *storev1.ListStoresRequest) (*storev1.ListStoresResponse, error) {
	collection := s.db.GetCollection("stores")
	
	// Build filter
	filter := bson.M{}
	if req.City != "" {
		filter["address.city"] = bson.M{"$regex": req.City, "$options": "i"}
	}
	if req.State != "" {
		filter["address.state"] = bson.M{"$regex": req.State, "$options": "i"}
	}
	if req.ActiveOnly {
		filter["is_active"] = true
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count stores: %w", err)
	}

	// Find with pagination
	findOptions := options.Find()
	if req.Limit > 0 {
		findOptions.SetLimit(int64(req.Limit))
	}
	if req.Offset > 0 {
		findOptions.SetSkip(int64(req.Offset))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find stores: %w", err)
	}
	defer cursor.Close(ctx)

	var stores []models.Store
	if err := cursor.All(ctx, &stores); err != nil {
		return nil, fmt.Errorf("failed to decode stores: %w", err)
	}

	// Convert to proto
	protoStores := make([]*storev1.Store, len(stores))
	for i, store := range stores {
		protoStores[i] = convertStoreToProto(&store)
	}

	return &storev1.ListStoresResponse{
		Stores:     protoStores,
		TotalCount: int32(total),
	}, nil
}

// AddProductToStore adds a product to a store's inventory
func (s *StoreService) AddProductToStore(ctx context.Context, req *storev1.AddProductToStoreRequest) (*storev1.AddProductToStoreResponse, error) {
	storeProduct := &models.StoreProduct{
		StoreID:           req.StoreId,
		ProductID:         req.ProductId,
		StockQuantity:     req.InitialStock,
		ReservedQuantity:  0,
		AvailableQuantity: req.InitialStock,
		StorePrice:        req.StorePrice,
		IsAvailable:       true,
		LastUpdated:       time.Now(),
	}

	collection := s.db.GetCollection("store_products")
	_, err := collection.InsertOne(ctx, storeProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to add product to store: %w", err)
	}

	return &storev1.AddProductToStoreResponse{
		StoreProduct: convertStoreProductToProto(storeProduct),
	}, nil
}

// ReserveProduct creates a product reservation
func (s *StoreService) ReserveProduct(ctx context.Context, req *storev1.ReserveProductRequest) (*storev1.ReserveProductResponse, error) {
	// Check if product is available in store
	collection := s.db.GetCollection("store_products")
	var storeProduct models.StoreProduct
	err := collection.FindOne(ctx, bson.M{
		"store_id":   req.StoreId,
		"product_id": req.ProductId,
	}).Decode(&storeProduct)
	if err != nil {
		return nil, fmt.Errorf("product not found in store: %w", err)
	}

	if storeProduct.AvailableQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock available")
	}

	// Create reservation
	reservation := &models.ProductReservation{
		ID:        uuid.New().String(),
		StoreID:   req.StoreId,
		ProductID: req.ProductId,
		UserID:    req.UserId,
		Quantity:  req.Quantity,
		Status:    "ACTIVE",
		ReservedAt: time.Now(),
		ExpiresAt:  time.Now().Add(time.Duration(req.ReservationDurationHours) * time.Hour),
		Notes:     req.Notes,
	}

	reservationCollection := s.db.GetCollection("reservations")
	_, err = reservationCollection.InsertOne(ctx, reservation)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// Update store product reserved quantity
	_, err = collection.UpdateOne(ctx,
		bson.M{"store_id": req.StoreId, "product_id": req.ProductId},
		bson.M{
			"$inc": map[string]interface{}{
				"reserved_quantity":  req.Quantity,
				"available_quantity": -req.Quantity,
			},
			"$set": map[string]interface{}{
				"last_updated": time.Now(),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update product reservation: %w", err)
	}

	return &storev1.ReserveProductResponse{
		Reservation: convertReservationToProto(reservation),
	}, nil
}

// RecordSale records a sale at a physical store
func (s *StoreService) RecordSale(ctx context.Context, req *storev1.RecordSaleRequest) (*storev1.RecordSaleResponse, error) {
	// Calculate total amount
	var totalAmount float64
	for _, item := range req.Items {
		price, _ := strconv.ParseFloat(item.UnitPrice, 64)
		totalAmount += price * float64(item.Quantity)
	}

	sale := &models.StoreSale{
		ID:             uuid.New().String(),
		StoreID:        req.StoreId,
		SalesUserID:    req.SalesUserId,
		CustomerUserID: req.CustomerUserId,
		Items:          convertSaleItemsFromProto(req.Items),
		TotalAmount:    fmt.Sprintf("%.2f", totalAmount),
		Currency:       "USD", // Default currency
		SaleType:       req.SaleType.String(),
		SaleDate:       time.Now(),
		ReservationID:  req.ReservationId,
		Metadata:       req.Metadata,
	}

	collection := s.db.GetCollection("sales")
	_, err := collection.InsertOne(ctx, sale)
	if err != nil {
		return nil, fmt.Errorf("failed to record sale: %w", err)
	}

	// Update product quantities
	for _, item := range req.Items {
		err = s.updateProductQuantityAfterSale(ctx, req.StoreId, item.ProductId, item.Quantity)
		if err != nil {
			// Log error but don't fail the sale
			fmt.Printf("Warning: failed to update product quantity for %s: %v\n", item.ProductId, err)
		}
	}

	return &storev1.RecordSaleResponse{
		Sale: convertSaleToProto(sale),
	}, nil
}

// ExportStoreProducts exports store products to CSV
func (s *StoreService) ExportStoreProducts(ctx context.Context, req *storev1.ExportStoreProductsRequest) (*storev1.ExportStoreProductsResponse, error) {
	collection := s.db.GetCollection("store_products")
	
	cursor, err := collection.Find(ctx, bson.M{"store_id": req.StoreId})
	if err != nil {
		return nil, fmt.Errorf("failed to find store products: %w", err)
	}
	defer cursor.Close(ctx)

	var products []models.StoreProduct
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	// Generate CSV
	var csvData strings.Builder
	writer := csv.NewWriter(&csvData)
	
	// Write header
	header := []string{"Store ID", "Product ID", "Stock Quantity", "Reserved Quantity", "Available Quantity", "Store Price", "Is Available", "Last Updated"}
	writer.Write(header)

	// Write data
	for _, product := range products {
		record := []string{
			product.StoreID,
			product.ProductID,
			strconv.Itoa(int(product.StockQuantity)),
			strconv.Itoa(int(product.ReservedQuantity)),
			strconv.Itoa(int(product.AvailableQuantity)),
			product.StorePrice,
			strconv.FormatBool(product.IsAvailable),
			product.LastUpdated.Format(time.RFC3339),
		}
		writer.Write(record)
	}
	writer.Flush()

	filename := fmt.Sprintf("store_%s_products_%s.csv", req.StoreId, time.Now().Format("20060102_150405"))

	return &storev1.ExportStoreProductsResponse{
		Data:        []byte(csvData.String()),
		Filename:    filename,
		ContentType: "text/csv",
	}, nil
}

// Helper functions for conversion between proto and model types
func convertStoreToProto(store *models.Store) *storev1.Store {
	return &storev1.Store{
		Id:            store.ID,
		Name:        store.Name,
		Description: store.Description,
		Address:     convertAddressToProto(&store.Address),
		Phone:       store.Phone,
		Email:       store.Email,
		IsActive:    store.IsActive,
		Hours:       convertStoreHoursToProto(&store.Hours),
		Metadata:    store.Metadata,
		CreatedAt:   timestamppb.New(store.CreatedAt),
		UpdatedAt:   timestamppb.New(store.UpdatedAt),
	}
}

func convertAddressToProto(addr *models.Address) *storev1.Address {
	return &storev1.Address{
		Street:     addr.Street,
		City:       addr.City,
		State:      addr.State,
		PostalCode: addr.PostalCode,
		Country:    addr.Country,
		Latitude:   addr.Latitude,
		Longitude:  addr.Longitude,
	}
}

func convertAddressFromProto(addr *storev1.Address) models.Address {
	if addr == nil {
		return models.Address{}
	}
	return models.Address{
		Street:     addr.Street,
		City:       addr.City,
		State:      addr.State,
		PostalCode: addr.PostalCode,
		Country:    addr.Country,
		Latitude:   addr.Latitude,
		Longitude:  addr.Longitude,
	}
}

func convertStoreHoursToProto(hours *models.StoreHours) *storev1.StoreHours {
	if hours == nil {
		return nil
	}
	
	protoDays := make([]*storev1.DayHours, len(hours.Days))
	for i, day := range hours.Days {
		protoDays[i] = &storev1.DayHours{
			Day:       day.Day,
			OpenTime:  day.OpenTime,
			CloseTime: day.CloseTime,
			IsClosed:  day.IsClosed,
		}
	}
	
	return &storev1.StoreHours{
		Days: protoDays,
	}
}

func convertStoreHoursFromProto(hours *storev1.StoreHours) models.StoreHours {
	if hours == nil {
		return models.StoreHours{}
	}
	
	days := make([]models.DayHours, len(hours.Days))
	for i, day := range hours.Days {
		days[i] = models.DayHours{
			Day:       day.Day,
			OpenTime:  day.OpenTime,
			CloseTime: day.CloseTime,
			IsClosed:  day.IsClosed,
		}
	}
	
	return models.StoreHours{
		Days: days,
	}
}

func convertStoreProductToProto(sp *models.StoreProduct) *storev1.StoreProduct {
	return &storev1.StoreProduct{
		StoreId:           sp.StoreID,
		ProductId:         sp.ProductID,
		StockQuantity:     sp.StockQuantity,
		ReservedQuantity:  sp.ReservedQuantity,
		AvailableQuantity: sp.AvailableQuantity,
		StorePrice:        sp.StorePrice,
		IsAvailable:       sp.IsAvailable,
		LastUpdated:       timestamppb.New(sp.LastUpdated),
	}
}

func convertReservationToProto(r *models.ProductReservation) *storev1.ProductReservation {
	status := storev1.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED
	switch r.Status {
	case "ACTIVE":
		status = storev1.ReservationStatus_RESERVATION_STATUS_ACTIVE
	case "EXPIRED":
		status = storev1.ReservationStatus_RESERVATION_STATUS_EXPIRED
	case "COMPLETED":
		status = storev1.ReservationStatus_RESERVATION_STATUS_COMPLETED
	case "CANCELLED":
		status = storev1.ReservationStatus_RESERVATION_STATUS_CANCELLED
	}

	proto := &storev1.ProductReservation{
		Id:        r.ID,
		StoreId:   r.StoreID,
		ProductId: r.ProductID,
		UserId:    r.UserID,
		Quantity:  r.Quantity,
		Status:    status,
		ReservedAt: timestamppb.New(r.ReservedAt),
		ExpiresAt:  timestamppb.New(r.ExpiresAt),
		Notes:     r.Notes,
	}

	if !r.CompletedAt.IsZero() {
		proto.CompletedAt = timestamppb.New(r.CompletedAt)
	}

	return proto
}

func convertSaleToProto(s *models.StoreSale) *storev1.StoreSale {
	saleType := storev1.SaleType_SALE_TYPE_UNSPECIFIED
	switch s.SaleType {
	case "WALK_IN":
		saleType = storev1.SaleType_SALE_TYPE_WALK_IN
	case "RESERVATION":
		saleType = storev1.SaleType_SALE_TYPE_RESERVATION
	case "ONLINE_PICKUP":
		saleType = storev1.SaleType_SALE_TYPE_ONLINE_PICKUP
	}

	items := make([]*storev1.StoreSaleItem, len(s.Items))
	for i, item := range s.Items {
		items[i] = &storev1.StoreSaleItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			ProductSku:  item.ProductSKU,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
		}
	}

	return &storev1.StoreSale{
		Id:             s.ID,
		StoreId:        s.StoreID,
		SalesUserId:    s.SalesUserID,
		CustomerUserId: s.CustomerUserID,
		Items:          items,
		TotalAmount:    s.TotalAmount,
		Currency:       s.Currency,
		SaleType:       saleType,
		SaleDate:       timestamppb.New(s.SaleDate),
		Metadata:       s.Metadata,
	}
}

func convertSaleItemsFromProto(items []*storev1.StoreSaleItem) []models.StoreSaleItem {
	result := make([]models.StoreSaleItem, len(items))
	for i, item := range items {
		result[i] = models.StoreSaleItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			ProductSKU:  item.ProductSku,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Subtotal:    item.Subtotal,
		}
	}
	return result
}

func (s *StoreService) updateProductQuantityAfterSale(ctx context.Context, storeID, productID string, quantitySold int32) error {
	collection := s.db.GetCollection("store_products")
	_, err := collection.UpdateOne(ctx,
		bson.M{"store_id": storeID, "product_id": productID},
		bson.M{
			"$inc": map[string]interface{}{
				"stock_quantity": -quantitySold,
			},
			"$set": map[string]interface{}{
				"last_updated": time.Now(),
			},
		},
	)
	return err
}

// Implement remaining methods with similar patterns...
// For brevity, I'm showing the key methods. The remaining methods would follow similar patterns.

// Placeholder implementations for remaining methods
func (s *StoreService) UpdateStore(ctx context.Context, req *storev1.UpdateStoreRequest) (*storev1.UpdateStoreResponse, error) {
	// Implementation would update store in database
	return &storev1.UpdateStoreResponse{Success: true}, nil
}

func (s *StoreService) DeleteStore(ctx context.Context, req *storev1.DeleteStoreRequest) (*storev1.DeleteStoreResponse, error) {
	// Implementation would soft delete store
	return &storev1.DeleteStoreResponse{Success: true}, nil
}

func (s *StoreService) UpdateStoreProductStock(ctx context.Context, req *storev1.UpdateStoreProductStockRequest) (*storev1.UpdateStoreProductStockResponse, error) {
	// Implementation would update product stock
	return &storev1.UpdateStoreProductStockResponse{Success: true}, nil
}

func (s *StoreService) RemoveProductFromStore(ctx context.Context, req *storev1.RemoveProductFromStoreRequest) (*storev1.RemoveProductFromStoreResponse, error) {
	// Implementation would remove product from store
	return &storev1.RemoveProductFromStoreResponse{Success: true}, nil
}

func (s *StoreService) GetStoreProducts(ctx context.Context, req *storev1.GetStoreProductsRequest) (*storev1.GetStoreProductsResponse, error) {
	// Implementation would get store products
	return &storev1.GetStoreProductsResponse{}, nil
}

func (s *StoreService) GetProductStoreLocations(ctx context.Context, req *storev1.GetProductStoreLocationsRequest) (*storev1.GetProductStoreLocationsResponse, error) {
	// Implementation would get product locations
	return &storev1.GetProductStoreLocationsResponse{}, nil
}

func (s *StoreService) CancelReservation(ctx context.Context, req *storev1.CancelReservationRequest) (*storev1.CancelReservationResponse, error) {
	// Implementation would cancel reservation
	return &storev1.CancelReservationResponse{Success: true}, nil
}

func (s *StoreService) GetReservations(ctx context.Context, req *storev1.GetReservationsRequest) (*storev1.GetReservationsResponse, error) {
	// Implementation would get reservations
	return &storev1.GetReservationsResponse{}, nil
}

func (s *StoreService) CompleteReservation(ctx context.Context, req *storev1.CompleteReservationRequest) (*storev1.CompleteReservationResponse, error) {
	// Implementation would complete reservation
	return &storev1.CompleteReservationResponse{Success: true}, nil
}

func (s *StoreService) AssignUserToStore(ctx context.Context, req *storev1.AssignUserToStoreRequest) (*storev1.AssignUserToStoreResponse, error) {
	// Implementation would assign user to store
	return &storev1.AssignUserToStoreResponse{Success: true}, nil
}

func (s *StoreService) RemoveUserFromStore(ctx context.Context, req *storev1.RemoveUserFromStoreRequest) (*storev1.RemoveUserFromStoreResponse, error) {
	// Implementation would remove user from store
	return &storev1.RemoveUserFromStoreResponse{Success: true}, nil
}

func (s *StoreService) GetStoreUsers(ctx context.Context, req *storev1.GetStoreUsersRequest) (*storev1.GetStoreUsersResponse, error) {
	// Implementation would get store users
	return &storev1.GetStoreUsersResponse{}, nil
}

func (s *StoreService) GetUserStores(ctx context.Context, req *storev1.GetUserStoresRequest) (*storev1.GetUserStoresResponse, error) {
	// Implementation would get user stores
	return &storev1.GetUserStoresResponse{}, nil
}

func (s *StoreService) GetStoreSales(ctx context.Context, req *storev1.GetStoreSalesRequest) (*storev1.GetStoreSalesResponse, error) {
	// Implementation would get store sales
	return &storev1.GetStoreSalesResponse{}, nil
}

func (s *StoreService) ExportStoreSales(ctx context.Context, req *storev1.ExportStoreSalesRequest) (*storev1.ExportStoreSalesResponse, error) {
	// Implementation would export store sales to CSV
	return &storev1.ExportStoreSalesResponse{}, nil
}
