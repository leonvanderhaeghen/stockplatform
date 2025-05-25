package domain

// ProductFilter defines the filter criteria for listing products
type ProductFilter struct {
	IDs         []string
	CategoryIDs []string
	MinPrice    float64
	MaxPrice    float64
	SearchTerm  string
}

// SortField defines the field to sort by
type SortField int

const (
	SortFieldUnspecified SortField = iota
	SortFieldName
	SortFieldPrice
	SortFieldCreatedAt
	SortFieldUpdatedAt
)

// SortOrder defines the sort order
type SortOrder int

const (
	SortOrderUnspecified SortOrder = iota
	SortOrderAsc
	SortOrderDesc
)

// SortOption defines the sorting options for listing products
type SortOption struct {
	Field SortField
	Order SortOrder
}

// Pagination defines the pagination options for listing products
type Pagination struct {
	Page     int
	PageSize int
}

// ListOptions contains all options for listing products
type ListOptions struct {
	Filter     *ProductFilter
	Sort       *SortOption
	Pagination *Pagination
}
