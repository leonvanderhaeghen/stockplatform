package domain

// POSProtoTypes defines the types needed to interact with proto messages for POS transactions
// These are temporary replacements until proper proto definitions are generated

// ProcessPOSTransactionRequest represents a request to process a POS transaction
type ProcessPOSTransactionRequest struct {
	TransactionType  string
	LocationID       string
	StaffID          string
	ReferenceOrderID string
	Items            []POSTransactionItem
	Payment          *PaymentInfo
}

// ProcessPOSTransactionResponse represents a response from processing a POS transaction
type ProcessPOSTransactionResponse struct {
	TransactionID  string
	Success        bool
	ProcessedItems []TransactionItemResult
	CompletedAt    string
	ReceiptURL     string
}

// TransactionItemResult represents the result of processing a transaction item
type TransactionItemResult struct {
	ProductID    string
	SKU          string
	Quantity     int32
	Price        float32
	Success      bool
	Description  string
	ErrorMessage string
}
