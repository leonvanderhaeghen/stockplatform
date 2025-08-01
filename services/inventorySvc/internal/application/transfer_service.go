package application

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
)

// TransferService handles business logic for inventory transfers between locations
type TransferService struct {
	transferRepo domain.TransferRepository
	inventoryRepo domain.InventoryRepository
	locationRepo domain.LocationRepository
	logger      *zap.Logger
}

// NewTransferService creates a new transfer service
func NewTransferService(
	transferRepo domain.TransferRepository,
	inventoryRepo domain.InventoryRepository,
	locationRepo domain.LocationRepository,
	logger *zap.Logger,
) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		inventoryRepo: inventoryRepo,
		locationRepo:  locationRepo,
		logger:      logger.Named("transfer_service"),
	}
}

// RequestTransfer creates a new inventory transfer request
func (s *TransferService) RequestTransfer(
	ctx context.Context,
	productID string,
	sku string,
	sourceLocationID string,
	destLocationID string,
	quantity int32,
	requestedBy string,
) (*domain.Transfer, error) {
	s.logger.Info("Requesting inventory transfer",
		zap.String("product_id", productID),
		zap.String("sku", sku),
		zap.String("source_location", sourceLocationID),
		zap.String("dest_location", destLocationID),
		zap.Int32("quantity", quantity),
	)

	// Validate source location exists
	sourceLocation, err := s.locationRepo.GetByID(ctx, sourceLocationID)
	if err != nil {
		return nil, err
	}
	if sourceLocation == nil || !sourceLocation.IsActive {
		return nil, errors.New("source location not found or inactive")
	}

	// Validate destination location exists
	destLocation, err := s.locationRepo.GetByID(ctx, destLocationID)
	if err != nil {
		return nil, err
	}
	if destLocation == nil || !destLocation.IsActive {
		return nil, errors.New("destination location not found or inactive")
	}

	// Validate inventory exists at source location with sufficient quantity
	sourceInventory, err := s.inventoryRepo.GetByProductAndLocation(ctx, productID, sourceLocationID)
	if err != nil {
		return nil, err
	}
	if sourceInventory == nil {
		return nil, errors.New("product not found at source location")
	}
	if !sourceInventory.IsAvailable(quantity) {
		return nil, errors.New("insufficient stock available at source location")
	}

	// Create the transfer request
	transferItems := []domain.TransferItem{
		{
			ProductID: productID,
			SKU:       sku,
			Quantity:  quantity,
			Status:    "pending",
		},
	}
	transfer := domain.NewTransfer(
		sourceLocationID,
		destLocationID,
		transferItems,
		requestedBy,
		"stock_transfer",
	)

	if err := s.transferRepo.Create(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

// GetTransfer retrieves a transfer by ID
func (s *TransferService) GetTransfer(ctx context.Context, id string) (*domain.Transfer, error) {
	s.logger.Debug("Getting transfer", zap.String("id", id))

	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if transfer == nil {
		return nil, errors.New("transfer not found")
	}

	return transfer, nil
}

// ApproveTransfer approves a transfer request and updates its status
func (s *TransferService) ApproveTransfer(ctx context.Context, id string, approvedBy string) error {
	s.logger.Info("Approving transfer", zap.String("id", id))

	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if transfer == nil {
		return errors.New("transfer not found")
	}

	if transfer.Status != domain.TransferStatusRequested {
		return errors.New("transfer is not in requested status")
	}

	// Set estimated arrival time - example calculation
	estimatedArrival := time.Now().Add(24 * time.Hour)
	
	// Use the built-in Approve method from the Transfer model
	transfer.Approve(approvedBy, &estimatedArrival)
	
	return s.transferRepo.Update(ctx, transfer)
}

// CompleteTransfer completes a transfer and moves inventory between locations
func (s *TransferService) CompleteTransfer(ctx context.Context, id string) error {
	s.logger.Info("Completing transfer", zap.String("id", id))

	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if transfer == nil {
		return errors.New("transfer not found")
	}

	if transfer.Status != domain.TransferStatusApproved && transfer.Status != domain.TransferStatusShipped {
		return errors.New("transfer is not approved or shipped")
	}

	// Make sure transfer has items
	if len(transfer.Items) == 0 {
		return errors.New("transfer has no items")
	}
	
	// Use the first item's data (typically transfers handle one product at a time)
	transferItem := transfer.Items[0]
	productID := transferItem.ProductID
	sku := transferItem.SKU
	quantity := transferItem.Quantity

	// Get source inventory
	sourceInventory, err := s.inventoryRepo.GetByProductAndLocation(ctx, productID, transfer.SourceLocationID)
	if err != nil {
		return err
	}
	
	if sourceInventory == nil {
		return errors.New("source inventory not found")
	}

	// Check if there's sufficient stock
	if !sourceInventory.IsAvailable(quantity) {
		transfer.Status = domain.TransferStatusCancelled
		if err := s.transferRepo.Update(ctx, transfer); err != nil {
			return err
		}
		return errors.New("insufficient stock at source location")
	}

	// Get or create destination inventory
	destInventory, err := s.inventoryRepo.GetByProductAndLocation(ctx, productID, transfer.DestinationLocationID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	// If destination inventory doesn't exist, create it
	if destInventory == nil {
		destInventory = domain.NewInventoryItem(
			productID,
			0, // Start with 0 quantity
			sku,
			transfer.DestinationLocationID,
		)
		
		if err := s.inventoryRepo.Create(ctx, destInventory); err != nil {
			return err
		}
	}

	// Perform the stock transfer
	if err := sourceInventory.TransferStock(quantity, destInventory); err != nil {
		return err
	}

	// Update both inventories
	if err := s.inventoryRepo.Update(ctx, sourceInventory); err != nil {
		return err
	}

	if err := s.inventoryRepo.Update(ctx, destInventory); err != nil {
		// Try to roll back the source inventory if destination update fails
		sourceInventory.AddStock(quantity)
		s.inventoryRepo.Update(ctx, sourceInventory) // best effort, ignore error
		return err
	}

	// Mark transfer as completed
	transfer.Complete("system")
	return s.transferRepo.Update(ctx, transfer)
}

// CancelTransfer cancels a pending or in-transit transfer
func (s *TransferService) CancelTransfer(ctx context.Context, id string) error {
	s.logger.Info("Cancelling transfer", zap.String("id", id))

	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if transfer == nil {
		return errors.New("transfer not found")
	}

	if transfer.Status == domain.TransferStatusCompleted {
		return errors.New("completed transfers cannot be cancelled")
	}

	if transfer.Status == domain.TransferStatusCancelled {
		return errors.New("transfer is already cancelled")
	}

	// Use Cancel method from Transfer domain model
	transfer.Status = domain.TransferStatusCancelled
	return s.transferRepo.Update(ctx, transfer)
}

// ListPendingTransfers retrieves all pending transfers
func (s *TransferService) ListPendingTransfers(ctx context.Context, limit, offset int) ([]*domain.Transfer, error) {
	s.logger.Debug("Listing pending transfers",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	return s.transferRepo.ListPendingTransfers(ctx, limit, offset)
}

// ListTransfersByStatus retrieves transfers with a specific status
func (s *TransferService) ListTransfersByStatus(ctx context.Context, status domain.TransferStatus, limit, offset int) ([]*domain.Transfer, error) {
	s.logger.Debug("Listing transfers by status",
		zap.String("status", string(status)),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	return s.transferRepo.ListByStatus(ctx, status, limit, offset)
}

// ListTransfersByLocation retrieves transfers for a specific location (source or destination)
func (s *TransferService) ListTransfersByLocation(ctx context.Context, locationID string, isSource bool, limit, offset int) ([]*domain.Transfer, error) {
	s.logger.Debug("Listing transfers by location",
		zap.String("location_id", locationID),
		zap.Bool("is_source", isSource),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	if isSource {
		return s.transferRepo.ListBySourceLocation(ctx, locationID, limit, offset)
	}
	return s.transferRepo.ListByDestLocation(ctx, locationID, limit, offset)
}

// ListTransfersByProduct retrieves transfers for a specific product
func (s *TransferService) ListTransfersByProduct(ctx context.Context, productID string, limit, offset int) ([]*domain.Transfer, error) {
	s.logger.Debug("Listing transfers by product",
		zap.String("product_id", productID),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	return s.transferRepo.ListByProduct(ctx, productID, limit, offset)
}
