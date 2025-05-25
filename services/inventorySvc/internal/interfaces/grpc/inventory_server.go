package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryv1 "stockplatform/pkg/gen/inventory/v1"
	"stockplatform/services/inventorySvc/internal/application"
	"stockplatform/services/inventorySvc/internal/domain"
)

// InventoryServer implements the gRPC interface for inventory service
type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	service *application.InventoryService
	logger  *zap.Logger
}

// NewInventoryServer creates a new inventory gRPC server
func NewInventoryServer(service *application.InventoryService, logger *zap.Logger) inventoryv1.InventoryServiceServer {
	return &InventoryServer{
		service: service,
		logger:  logger.Named("inventory_grpc_server"),
	}
}

// CreateInventory creates a new inventory item
func (s *InventoryServer) CreateInventory(ctx context.Context, req *inventoryv1.CreateInventoryRequest) (*inventoryv1.CreateInventoryResponse, error) {
	s.logger.Info("gRPC CreateInventory called",
		zap.String("product_id", req.ProductId),
		zap.Int32("quantity", req.Quantity),
		zap.String("sku", req.Sku),
	)

	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}
	if req.Quantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be non-negative")
	}
	if req.Sku == "" {
		return nil, status.Error(codes.InvalidArgument, "sku is required")
	}

	item, err := s.service.CreateInventoryItem(ctx, req.ProductId, req.Quantity, req.Sku, req.Location)
	if err != nil {
		s.logger.Error("Failed to create inventory item", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create inventory item: "+err.Error())
	}

	return &inventoryv1.CreateInventoryResponse{
		Inventory: toProtoInventoryItem(item),
	}, nil
}

// GetInventory retrieves an inventory item by ID
func (s *InventoryServer) GetInventory(ctx context.Context, req *inventoryv1.GetInventoryRequest) (*inventoryv1.GetInventoryResponse, error) {
	s.logger.Debug("gRPC GetInventory called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	item, err := s.service.GetInventoryItem(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get inventory item", zap.Error(err))
		return nil, status.Error(codes.NotFound, "inventory item not found")
	}

	return &inventoryv1.GetInventoryResponse{
		Inventory: toProtoInventoryItem(item),
	}, nil
}

// GetInventoryByProductID retrieves an inventory item by product ID
func (s *InventoryServer) GetInventoryByProductID(ctx context.Context, req *inventoryv1.GetInventoryByProductIDRequest) (*inventoryv1.GetInventoryResponse, error) {
	s.logger.Debug("gRPC GetInventoryByProductID called", zap.String("product_id", req.ProductId))

	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	item, err := s.service.GetInventoryItemByProductID(ctx, req.ProductId)
	if err != nil {
		s.logger.Error("Failed to get inventory item by product ID", zap.Error(err))
		return nil, status.Error(codes.NotFound, "inventory item not found")
	}

	return &inventoryv1.GetInventoryResponse{
		Inventory: toProtoInventoryItem(item),
	}, nil
}

// GetInventoryBySKU retrieves an inventory item by SKU
func (s *InventoryServer) GetInventoryBySKU(ctx context.Context, req *inventoryv1.GetInventoryBySKURequest) (*inventoryv1.GetInventoryResponse, error) {
	s.logger.Debug("gRPC GetInventoryBySKU called", zap.String("sku", req.Sku))

	if req.Sku == "" {
		return nil, status.Error(codes.InvalidArgument, "sku is required")
	}

	item, err := s.service.GetInventoryItemBySKU(ctx, req.Sku)
	if err != nil {
		s.logger.Error("Failed to get inventory item by SKU", zap.Error(err))
		return nil, status.Error(codes.NotFound, "inventory item not found")
	}

	return &inventoryv1.GetInventoryResponse{
		Inventory: toProtoInventoryItem(item),
	}, nil
}

// UpdateInventory updates an existing inventory item
func (s *InventoryServer) UpdateInventory(ctx context.Context, req *inventoryv1.UpdateInventoryRequest) (*inventoryv1.UpdateInventoryResponse, error) {
	s.logger.Info("gRPC UpdateInventory called", zap.String("id", req.Inventory.Id))

	if req.Inventory == nil {
		return nil, status.Error(codes.InvalidArgument, "inventory is required")
	}
	if req.Inventory.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory.id is required")
	}

	// Fetch the existing item first
	existingItem, err := s.service.GetInventoryItem(ctx, req.Inventory.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "inventory item not found")
	}

	// Update the item with the new values
	existingItem.ProductID = req.Inventory.ProductId
	existingItem.Quantity = req.Inventory.Quantity
	existingItem.Reserved = req.Inventory.Reserved
	existingItem.SKU = req.Inventory.Sku
	existingItem.Location = req.Inventory.Location
	existingItem.LastUpdated = time.Now()

	if err := s.service.UpdateInventoryItem(ctx, existingItem); err != nil {
		s.logger.Error("Failed to update inventory item", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update inventory item: "+err.Error())
	}

	return &inventoryv1.UpdateInventoryResponse{
		Success: true,
	}, nil
}

// DeleteInventory removes an inventory item
func (s *InventoryServer) DeleteInventory(ctx context.Context, req *inventoryv1.DeleteInventoryRequest) (*inventoryv1.DeleteInventoryResponse, error) {
	s.logger.Info("gRPC DeleteInventory called", zap.String("id", req.Id))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.service.DeleteInventoryItem(ctx, req.Id); err != nil {
		s.logger.Error("Failed to delete inventory item", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete inventory item: "+err.Error())
	}

	return &inventoryv1.DeleteInventoryResponse{
		Success: true,
	}, nil
}

// ListInventory lists all inventory items with pagination
func (s *InventoryServer) ListInventory(ctx context.Context, req *inventoryv1.ListInventoryRequest) (*inventoryv1.ListInventoryResponse, error) {
	s.logger.Debug("gRPC ListInventory called",
		zap.Int32("limit", req.Limit),
		zap.Int32("offset", req.Offset),
	)

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10 // Default limit
	}
	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	items, err := s.service.ListInventoryItems(ctx, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list inventory items", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list inventory items: "+err.Error())
	}

	response := &inventoryv1.ListInventoryResponse{
		Inventories: make([]*inventoryv1.InventoryItem, 0, len(items)),
	}

	for _, item := range items {
		response.Inventories = append(response.Inventories, toProtoInventoryItem(item))
	}

	return response, nil
}

// AddStock adds stock to an inventory item
func (s *InventoryServer) AddStock(ctx context.Context, req *inventoryv1.AddStockRequest) (*inventoryv1.AddStockResponse, error) {
	s.logger.Info("gRPC AddStock called",
		zap.String("id", req.Id),
		zap.Int32("quantity", req.Quantity),
	)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	if err := s.service.AddStock(ctx, req.Id, req.Quantity); err != nil {
		s.logger.Error("Failed to add stock", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to add stock: "+err.Error())
	}

	return &inventoryv1.AddStockResponse{
		Success: true,
	}, nil
}

// RemoveStock removes stock from an inventory item
func (s *InventoryServer) RemoveStock(ctx context.Context, req *inventoryv1.RemoveStockRequest) (*inventoryv1.RemoveStockResponse, error) {
	s.logger.Info("gRPC RemoveStock called",
		zap.String("id", req.Id),
		zap.Int32("quantity", req.Quantity),
	)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	if err := s.service.RemoveStock(ctx, req.Id, req.Quantity); err != nil {
		s.logger.Error("Failed to remove stock", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to remove stock: "+err.Error())
	}

	return &inventoryv1.RemoveStockResponse{
		Success: true,
	}, nil
}

// ReserveStock reserves stock for an order
func (s *InventoryServer) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	s.logger.Info("gRPC ReserveStock called",
		zap.String("id", req.Id),
		zap.Int32("quantity", req.Quantity),
	)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	if err := s.service.ReserveStock(ctx, req.Id, req.Quantity); err != nil {
		s.logger.Error("Failed to reserve stock", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to reserve stock: "+err.Error())
	}

	return &inventoryv1.ReserveStockResponse{
		Success: true,
	}, nil
}

// ReleaseReservation releases a reservation without fulfilling it
func (s *InventoryServer) ReleaseReservation(ctx context.Context, req *inventoryv1.ReleaseReservationRequest) (*inventoryv1.ReleaseReservationResponse, error) {
	s.logger.Info("gRPC ReleaseReservation called",
		zap.String("id", req.Id),
		zap.Int32("quantity", req.Quantity),
	)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	if err := s.service.ReleaseReservation(ctx, req.Id, req.Quantity); err != nil {
		s.logger.Error("Failed to release reservation", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to release reservation: "+err.Error())
	}

	return &inventoryv1.ReleaseReservationResponse{
		Success: true,
	}, nil
}

// FulfillReservation completes a reservation and deducts from stock
func (s *InventoryServer) FulfillReservation(ctx context.Context, req *inventoryv1.FulfillReservationRequest) (*inventoryv1.FulfillReservationResponse, error) {
	s.logger.Info("gRPC FulfillReservation called",
		zap.String("id", req.Id),
		zap.Int32("quantity", req.Quantity),
	)

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	if err := s.service.FulfillReservation(ctx, req.Id, req.Quantity); err != nil {
		s.logger.Error("Failed to fulfill reservation", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to fulfill reservation: "+err.Error())
	}

	return &inventoryv1.FulfillReservationResponse{
		Success: true,
	}, nil
}

// toProtoInventoryItem converts a domain inventory item to a proto inventory item
func toProtoInventoryItem(item *domain.InventoryItem) *inventoryv1.InventoryItem {
	return &inventoryv1.InventoryItem{
		Id:          item.ID,
		ProductId:   item.ProductID,
		Quantity:    item.Quantity,
		Reserved:    item.Reserved,
		Sku:         item.SKU,
		Location:    item.Location,
		LastUpdated: item.LastUpdated.Format(time.RFC3339),
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
	}
}
