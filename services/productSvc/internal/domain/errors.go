package domain

import (
	"errors"
	"fmt"
)

// Common domain errors
var (
	// General errors
	ErrInternal        = errors.New("internal server error")
	ErrNotFound        = errors.New("resource not found")
	ErrInvalidID       = errors.New("invalid ID format")
	ErrValidation      = errors.New("validation error")
	ErrAlreadyExists   = errors.New("resource already exists")
	ErrInvalidArgument = errors.New("invalid argument")

	// Product errors
	ErrProductNotFound           = fmt.Errorf("%w: product not found", ErrNotFound)
	ErrProductAlreadyExists      = fmt.Errorf("%w: product with same SKU or barcode already exists", ErrAlreadyExists)
	ErrProductNameRequired       = fmt.Errorf("%w: product name is required", ErrValidation)
	ErrProductSKURequired        = fmt.Errorf("%w: product SKU is required", ErrValidation)
	ErrInvalidCostPrice         = fmt.Errorf("%w: invalid cost price", ErrValidation)
	ErrSellingPriceRequired     = fmt.Errorf("%w: selling price is required", ErrValidation)
	ErrInvalidSellingPrice      = fmt.Errorf("%w: invalid selling price", ErrValidation)
	ErrInvalidStockQuantity     = fmt.Errorf("%w: stock quantity cannot be negative", ErrValidation)
	ErrInvalidCurrency          = fmt.Errorf("%w: currency must be a 3-letter ISO code", ErrValidation)
	ErrProductNotActive         = errors.New("product is not active")
	ErrInsufficientStock        = errors.New("insufficient stock")

	// Variant errors
	ErrVariantNotFound          = fmt.Errorf("%w: variant not found", ErrNotFound)
	ErrVariantAlreadyExists     = fmt.Errorf("%w: variant already exists", ErrAlreadyExists)
	ErrVariantSKURequired       = fmt.Errorf("%w: variant SKU is required", ErrValidation)
	ErrVariantOptionRequired    = fmt.Errorf("%w: at least one variant option is required", ErrValidation)

	// Variant option errors
	ErrOptionNameRequired       = fmt.Errorf("%w: option name is required", ErrValidation)
	ErrOptionValueRequired      = fmt.Errorf("%w: option value is required", ErrValidation)
	ErrInvalidPriceAdjustment   = fmt.Errorf("%w: invalid price adjustment", ErrValidation)

	// Category errors
	ErrCategoryNotFound         = fmt.Errorf("%w: category not found", ErrNotFound)
	ErrCategoryNameRequired     = fmt.Errorf("%w: category name is required", ErrValidation)
	ErrParentCategoryNotFound   = fmt.Errorf("%w: parent category not found", ErrValidation)
	ErrCategoryInUse            = errors.New("category is in use by one or more products")
	ErrInvalidCategoryHierarchy = errors.New("invalid category hierarchy")

	// Supplier errors
	ErrSupplierNotFound         = fmt.Errorf("%w: supplier not found", ErrNotFound)
	ErrSupplierNameRequired     = fmt.Errorf("%w: supplier name is required", ErrValidation)

	// Inventory errors
	ErrInvalidQuantity          = fmt.Errorf("%w: quantity must be greater than zero", ErrValidation)
	ErrStockAdjustmentFailed    = errors.New("failed to adjust stock")
	ErrStockTransferFailed      = errors.New("failed to transfer stock")

	// Search errors
	ErrInvalidSearchQuery       = fmt.Errorf("%w: invalid search query", ErrValidation)
)
