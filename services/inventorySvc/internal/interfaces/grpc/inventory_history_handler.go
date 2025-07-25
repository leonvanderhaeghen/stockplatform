package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/application"
	inventorypb "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
)

// InventoryHistoryServer handles inventory history gRPC requests
type InventoryHistoryServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	service *application.InventoryService
	logger  *zap.Logger
}

// NewInventoryHistoryServer creates a new inventory history server
func NewInventoryHistoryServer(service *application.InventoryService, logger *zap.Logger) *InventoryHistoryServer {
	return &InventoryHistoryServer{
		service: service,
		logger:  logger.Named("inventory_history_handler"),
	}
}

// GetInventoryHistory retrieves the history of changes for a specific inventory item
func (s *InventoryHistoryServer) GetInventoryHistory(
	ctx context.Context,
	req *inventorypb.GetInventoryHistoryRequest,
) (*inventorypb.GetInventoryHistoryResponse, error) {
	// Validate request
	if req.InventoryId == "" {
		return nil, status.Error(codes.InvalidArgument, "inventory_id is required")
	}

	// Set default values for pagination
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	s.logger.Debug("Getting inventory history",
		zap.String("inventory_id", req.InventoryId),
		zap.Int32("limit", limit),
		zap.Int32("offset", offset),
	)

	// Call the service layer
	history, total, err := s.service.GetInventoryHistory(ctx, req.InventoryId, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get inventory history",
			zap.String("inventory_id", req.InventoryId),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to get inventory history")
	}

	// Convert domain models to protobuf
	entries := make([]*inventorypb.InventoryHistoryEntry, 0, len(history))
	for _, h := range history {
		entries = append(entries, &inventorypb.InventoryHistoryEntry{
			Id:            h.ID,
			InventoryId:   h.InventoryID,
			ChangeType:    h.ChangeType,
			Description:   h.Description,
			QuantityBefore: h.QuantityBefore,
			QuantityAfter:  h.QuantityAfter,
			ReferenceId:   h.ReferenceID,
			ReferenceType: h.ReferenceType,
			PerformedBy:   h.PerformedBy,
			CreatedAt:     h.CreatedAt.Format(time.RFC3339),
		})
	}

	return &inventorypb.GetInventoryHistoryResponse{
		Entries: entries,
		Total:   total,
	}, nil
}
