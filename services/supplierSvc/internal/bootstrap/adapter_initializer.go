package bootstrap

import (
	"context"
	"log"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/adapters"
	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

// RegisterAdapters registers all available supplier adapters with the registry
func RegisterAdapters(supplierUseCase domain.SupplierUseCase) {
	ctx := context.Background()

	// Create all supplier adapters
	sampleAdapter := adapters.NewSampleAdapter()

	// Register each adapter
	if err := supplierUseCase.RegisterAdapter(ctx, sampleAdapter); err != nil {
		log.Printf("Failed to register sample adapter: %v", err)
	} else {
		log.Printf("Successfully registered adapter: %s", sampleAdapter.Name())
	}

	// Add more adapters here as they are implemented
	// Example:
	// shopifyAdapter := adapters.NewShopifyAdapter()
	// if err := supplierUseCase.RegisterAdapter(ctx, shopifyAdapter); err != nil {
	//     log.Printf("Failed to register Shopify adapter: %v", err)
	// } else {
	//     log.Printf("Successfully registered adapter: %s", shopifyAdapter.Name())
	// }
}
