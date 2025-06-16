package domain

import (
	"time"
)

// TransferStatus represents the status of an inventory transfer
type TransferStatus string

const (
	// TransferStatusRequested indicates a transfer has been requested but not yet approved
	TransferStatusRequested TransferStatus = "requested"
	
	// TransferStatusApproved indicates a transfer has been approved but items haven't been shipped
	TransferStatusApproved TransferStatus = "approved"
	
	// TransferStatusShipped indicates items have been shipped but not received
	TransferStatusShipped TransferStatus = "shipped"
	
	// TransferStatusCompleted indicates items have been received at the destination
	TransferStatusCompleted TransferStatus = "completed"
	
	// TransferStatusCancelled indicates a transfer was cancelled
	TransferStatusCancelled TransferStatus = "cancelled"
	
	// TransferStatusRejected indicates a transfer was rejected
	TransferStatusRejected TransferStatus = "rejected"
)

// TransferItem represents an item in a transfer
type TransferItem struct {
	ProductID  string `bson:"product_id" json:"product_id"`
	SKU        string `bson:"sku" json:"sku,omitempty"`
	Quantity   int32  `bson:"quantity" json:"quantity"`
	Status     string `bson:"status" json:"status"`
}

// Transfer represents an inventory transfer between locations
type Transfer struct {
	ID                   string          `bson:"_id" json:"id"`
	SourceLocationID     string          `bson:"source_location_id" json:"source_location_id"`
	DestinationLocationID string         `bson:"destination_location_id" json:"destination_location_id"`
	Items                []TransferItem  `bson:"items" json:"items"`
	RequestedBy          string          `bson:"requested_by" json:"requested_by"`
	ApprovedBy           string          `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	ReceivedBy           string          `bson:"received_by,omitempty" json:"received_by,omitempty"`
	RequestedAt          time.Time       `bson:"requested_at" json:"requested_at"`
	ApprovedAt           *time.Time      `bson:"approved_at,omitempty" json:"approved_at,omitempty"`
	ShippedAt            *time.Time      `bson:"shipped_at,omitempty" json:"shipped_at,omitempty"`
	ReceivedAt           *time.Time      `bson:"received_at,omitempty" json:"received_at,omitempty"`
	EstimatedArrival     *time.Time      `bson:"estimated_arrival,omitempty" json:"estimated_arrival,omitempty"`
	Status               TransferStatus  `bson:"status" json:"status"`
	Notes                string          `bson:"notes,omitempty" json:"notes,omitempty"`
	Reason               string          `bson:"reason" json:"reason"`
}

// GenerateID generates a unique ID for domain objects
func GenerateID() string {
	return time.Now().Format("20060102-150405") + "-" + time.Now().Format("000000")
}

// NewTransfer creates a new inventory transfer
func NewTransfer(sourceLocationID, destinationLocationID string, items []TransferItem, requestedBy, reason string) *Transfer {
	return &Transfer{
		ID:                   GenerateID(),
		SourceLocationID:     sourceLocationID,
		DestinationLocationID: destinationLocationID,
		Items:                items,
		RequestedBy:          requestedBy,
		RequestedAt:          time.Now(),
		Status:               TransferStatusRequested,
		Reason:               reason,
	}
}

// Approve changes the transfer status to approved
func (t *Transfer) Approve(approvedBy string, estimatedArrival *time.Time) {
	now := time.Now()
	t.Status = TransferStatusApproved
	t.ApprovedBy = approvedBy
	t.ApprovedAt = &now
	t.EstimatedArrival = estimatedArrival
}

// Ship changes the transfer status to shipped
func (t *Transfer) Ship() {
	now := time.Now()
	t.Status = TransferStatusShipped
	t.ShippedAt = &now
}

// Complete changes the transfer status to completed
func (t *Transfer) Complete(receivedBy string) {
	now := time.Now()
	t.Status = TransferStatusCompleted
	t.ReceivedBy = receivedBy
	t.ReceivedAt = &now
}

// Cancel changes the transfer status to cancelled
func (t *Transfer) Cancel() {
	t.Status = TransferStatusCancelled
}

// Reject changes the transfer status to rejected
func (t *Transfer) Reject() {
	t.Status = TransferStatusRejected
}

// Note: TransferRepository for Transfer objects is defined in a separate file
// to avoid redeclaration conflicts with the InventoryTransfer repository
