package domain

// POSTransactionItem represents an item in a point-of-sale transaction
type POSTransactionItem struct {
	ProductID       string
	SKU             string
	Quantity        int32
	Price           float32
	Description     string
	Reason          string
	Direction       string // For exchanges: "in" or "out"
	InventoryItemID string
}

// PaymentInfo represents payment information for a POS transaction
type PaymentInfo struct {
	Method          string
	Amount          float32
	CurrencyCode    string
	PaymentReference string
	CardLast4       string
}
