package grpc

import (
	"context"
	"time"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
	inventoryv1 "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateTransfer handles the CreateTransfer gRPC request
func (s *InventoryServer) CreateTransfer(ctx context.Context, req *inventoryv1.CreateTransferRequest) (*inventoryv1.CreateTransferResponse, error) {
	s.logger.Info("Creating inventory transfer",
		zap.String("source_location", req.SourceLocationId),
		zap.String("destination_location", req.DestinationLocationId),
		zap.String("product_id", req.ProductId),
		zap.Int32("quantity", req.Quantity),
	)

	// Create a transfer request with the provided information
	transfer, err := s.transferService.RequestTransfer(
		ctx,
		req.ProductId,
		req.Sku,
		req.SourceLocationId,
		req.DestinationLocationId,
		req.Quantity,
		req.RequestedBy,
	)
	if err != nil {
		s.logger.Error("Failed to create inventory transfer", zap.Error(err))
		return nil, err
	}

	// Map domain model to proto response
	return &inventoryv1.CreateTransferResponse{
		Transfer: mapDomainTransferToProto(transfer),
	}, nil
}

// GetTransfer handles the GetTransfer gRPC request
func (s *InventoryServer) GetTransfer(ctx context.Context, req *inventoryv1.GetTransferRequest) (*inventoryv1.GetTransferResponse, error) {
	s.logger.Info("Getting inventory transfer", zap.String("id", req.Id))

	transfer, err := s.transferService.GetTransfer(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get inventory transfer", 
			zap.String("id", req.Id), 
			zap.Error(err))
		return nil, err
	}

	return &inventoryv1.GetTransferResponse{
		Transfer: mapDomainTransferToProto(transfer),
	}, nil
}

// UpdateTransferStatus handles the UpdateTransferStatus gRPC request
func (s *InventoryServer) UpdateTransferStatus(ctx context.Context, req *inventoryv1.UpdateTransferStatusRequest) (*inventoryv1.UpdateTransferStatusResponse, error) {
	s.logger.Info("Updating inventory transfer status",
		zap.String("id", req.Id),
		zap.String("status", req.Status),
	)

	// Verify the transfer exists before attempting to update
	_, err := s.transferService.GetTransfer(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get inventory transfer for status update", zap.Error(err))
		return nil, err
	}

	switch req.Status {
	case "approved":
		err = s.transferService.ApproveTransfer(ctx, req.Id, req.ApprovedBy)
	case "completed":
		err = s.transferService.CompleteTransfer(ctx, req.Id)
	case "cancelled":
		err = s.transferService.CancelTransfer(ctx, req.Id)
	default:
		s.logger.Warn("Unsupported transfer status", zap.String("status", req.Status))
		return nil, status.Error(codes.InvalidArgument, "invalid transfer status")
	}

	if err != nil {
		s.logger.Error("Failed to update inventory transfer status", zap.Error(err))
		return nil, err
	}

	// Fetch the updated transfer to return in the response
	updatedTransfer, err := s.transferService.GetTransfer(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get updated inventory transfer", zap.Error(err))
		return nil, err
	}

	return &inventoryv1.UpdateTransferStatusResponse{
		Success:  true,
		Transfer: mapDomainTransferToProto(updatedTransfer),
	}, nil
}

// ListTransfers handles the ListTransfers gRPC request
func (s *InventoryServer) ListTransfers(ctx context.Context, req *inventoryv1.ListTransfersRequest) (*inventoryv1.ListTransfersResponse, error) {
	s.logger.Info("Listing inventory transfers",
		zap.Int32("limit", req.Limit),
		zap.Int32("offset", req.Offset),
		zap.String("status", req.Status),
	)

	limit := int(req.Limit)
	offset := int(req.Offset)

	var transfers []*domain.Transfer
	var err error

	// Apply filters if provided
	if req.Status != "" {
		// Convert string status to domain TransferStatus
		status := domain.TransferStatus(req.Status)
		transfers, err = s.transferService.ListTransfersByStatus(ctx, status, limit, offset)
	} else if req.SourceLocationId != "" {
		transfers, err = s.transferService.ListTransfersByLocation(ctx, req.SourceLocationId, true, limit, offset)
	} else if req.DestinationLocationId != "" {
		transfers, err = s.transferService.ListTransfersByLocation(ctx, req.DestinationLocationId, false, limit, offset)
	} else if req.ProductId != "" {
		transfers, err = s.transferService.ListTransfersByProduct(ctx, req.ProductId, limit, offset)
	} else {
		// Get pending transfers if no filter is specified
		transfers, err = s.transferService.ListPendingTransfers(ctx, limit, offset)
	}

	if err != nil {
		s.logger.Error("Failed to list inventory transfers", zap.Error(err))
		return nil, err
	}

	// Map domain transfers to proto transfers
	protoTransfers := make([]*inventoryv1.InventoryTransfer, 0, len(transfers))
	for _, transfer := range transfers {
		protoTransfers = append(protoTransfers, mapDomainTransferToProto(transfer))
	}

	return &inventoryv1.ListTransfersResponse{
		Transfers: protoTransfers,
		Total:     int32(len(protoTransfers)),
	}, nil
}

// Helper function to map domain Transfer to proto InventoryTransfer
func mapDomainTransferToProto(t *domain.Transfer) *inventoryv1.InventoryTransfer {
	if t == nil {
		return nil
	}
	
	// Extract the first item details for the proto model (since proto expects a single item)
	var productID, sku string
	var quantity int32
	if len(t.Items) > 0 {
		productID = t.Items[0].ProductID
		sku = t.Items[0].SKU
		quantity = int32(t.Items[0].Quantity)
	}
	
	result := &inventoryv1.InventoryTransfer{
		Id:                    t.ID,
		SourceLocationId:      t.SourceLocationID,
		DestinationLocationId: t.DestinationLocationID,
		ProductId:             productID,
		Sku:                   sku,
		Quantity:              quantity,
		Status:                string(t.Status),
		RequestedBy:           t.RequestedBy,
		ApprovedBy:            t.ApprovedBy,
		RequestedDate:         timestampToString(t.RequestedAt),
		CreatedAt:             timestampToString(t.RequestedAt),
		UpdatedAt:             timestampToString(t.RequestedAt), // Default to requested time
		Notes:                 "", // Not mapped in domain model
	}

	// Set optional fields if available
	if t.ApprovedAt != nil {
		result.UpdatedAt = timestampToString(*t.ApprovedAt)
	}

	if t.EstimatedArrival != nil {
		result.ExpectedDeliveryDate = timestampToString(*t.EstimatedArrival)
	}

	if t.ReceivedAt != nil {
		result.ActualDeliveryDate = timestampToString(*t.ReceivedAt)
	}

	return result
}

// Helper function to convert timestamp to string
func timestampToString(t interface{}) string {
	switch v := t.(type) {
	case timestamppb.Timestamp:
		return v.AsTime().Format(time.RFC3339)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return ""
	}
}
