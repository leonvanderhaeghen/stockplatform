package application

import (
	"context"
	"time"
	"math/rand"

	inventoryclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/inventory"
	inventorypb "github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/api/gen/go/proto/inventory/v1"
	"github.com/leonvanderhaeghen/stockplatform/services/orderSvc/internal/domain"
)

// POSTransactionService handles Point of Sale transactions
type POSTransactionService struct {
	orderService *OrderService
	config       *domain.ServiceConfig
}

// NewPOSTransactionService creates a new POS transaction service
func NewPOSTransactionService(orderService *OrderService, config *domain.ServiceConfig) *POSTransactionService {
	return &POSTransactionService{
		orderService: orderService,
		config:       config,
	}
}

// TransactionResult represents the result of a POS transaction
type TransactionResult struct {
	TransactionID  string
	Success        bool
	ProcessedItems []*ProcessedItem
	CompletedAt    time.Time
	ReceiptURL     string
}

// ProcessedItem represents a processed item in a POS transaction
type ProcessedItem struct {
	ProductID    string
	SKU          string
	Quantity     int32
	Price        float32
	Success      bool
	Description  string
	ErrorMessage string
}

// InventoryAdjustmentItem represents an item for inventory adjustment
type InventoryAdjustmentItem struct {
	ProductID       string
	SKU             string
	Quantity        int32
	Reason          string
	InventoryItemID string
}

// ProcessTransaction processes a point-of-sale transaction
func (s *POSTransactionService) ProcessTransaction(
	ctx context.Context,
	transactionType string,
	locationID string,
	staffID string,
	referenceOrderID string,
	items []domain.POSTransactionItem,
	payment *domain.PaymentInfo,
) (*TransactionResult, error) {
	// Generate transaction ID
	transactionID := s.generateTransactionID()

	// Convert order items to inventory adjustment items
	adjustmentItems := make([]*InventoryAdjustmentItem, 0, len(items))
	for _, item := range items {
		// For sales, quantity is negative; for returns, positive; for exchanges, it depends on the item direction
		quantity := item.Quantity
		if transactionType == "sale" {
			quantity = -quantity // Make negative for sales (inventory reduction)
		} else if transactionType == "exchange" && item.Direction == "out" {
			quantity = -quantity // Make negative for outgoing exchange items
		}
		// Returns and incoming exchange items stay positive (inventory addition)

		adjustmentItems = append(adjustmentItems, &InventoryAdjustmentItem{
			ProductID:       item.ProductID,
			SKU:             item.SKU,
			Quantity:        quantity,
			Reason:          item.Reason,
			InventoryItemID: item.InventoryItemID,
		})
	}

	// Call inventory service via gRPC client to adjust stock
	invCfg := inventoryclient.Config{Address: s.config.InventoryServiceAddr}
	inventoryClient, err := inventoryclient.New(invCfg, nil)
	if err != nil {
		return nil, err
	}
	defer inventoryClient.Close()

	// Create gRPC request items
	grpcItems := make([]*inventorypb.InventoryAdjustmentItem, 0, len(adjustmentItems))
	for _, item := range adjustmentItems {
		grpcItems = append(grpcItems, &inventorypb.InventoryAdjustmentItem{
			ProductId:       item.ProductID,
			Sku:             item.SKU,
			Quantity:        item.Quantity,
			Reason:          item.Reason,
			InventoryItemId: item.InventoryItemID,
		})
	}

	// Call inventory service to adjust stock
	adjustReq := &inventorypb.AdjustInventoryForOrderRequest{
		OrderId:        transactionID,
		LocationId:     locationID,
		AdjustmentType: transactionType,
		ReferenceId:    referenceOrderID,
		Items:          grpcItems,
		StaffId:        staffID,
	}

	adjustResp, err := inventoryClient.AdjustInventoryForOrder(ctx, adjustReq)
	if err != nil {
		return nil, err
	}

	// Process transaction results based on inventory adjustment
	processedItems := make([]*ProcessedItem, 0, len(items))
	for i, item := range items {
		var adjustResult *inventorypb.InventoryAdjustmentResult
		if i < len(adjustResp.Items) {
			adjustResult = adjustResp.Items[i]
		}

		result := &ProcessedItem{
			ProductID:   item.ProductID,
			SKU:         item.SKU,
			Quantity:    item.Quantity,
			Price:       item.Price,
			Success:     adjustResult != nil && adjustResult.Success,
			Description: item.Description,
		}

		if !result.Success {
			result.ErrorMessage = "Failed to process inventory adjustment"
			if adjustResult != nil && adjustResult.ErrorMessage != "" {
				result.ErrorMessage = adjustResult.ErrorMessage
			}
		}

		processedItems = append(processedItems, result)
	}

	// Create response
	result := &TransactionResult{
		TransactionID:  transactionID,
		Success:        adjustResp.Success,
		ProcessedItems: processedItems,
		CompletedAt:    time.Now(),
		ReceiptURL:     s.generateReceiptURL(transactionID),
	}

	return result, nil
}

// Helper function to generate a transaction ID
func (s *POSTransactionService) generateTransactionID() string {
	return "TX-" + time.Now().Format("20060102-150405") + "-" + s.generateRandomString(6)
}

// Helper function to generate a random string
func (s *POSTransactionService) generateRandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Helper function to generate a receipt URL
func (s *POSTransactionService) generateReceiptURL(transactionID string) string {
	return "/receipts/" + transactionID + ".pdf"
}
