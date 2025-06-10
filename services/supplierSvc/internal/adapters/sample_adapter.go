package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

// SampleAdapter is an example adapter implementation for a fictional supplier API
type SampleAdapter struct {
	name         string
	client       *http.Client
	apiURL       string
	apiKey       string
	config       map[string]string
	capabilities map[string]bool
}

// NewSampleAdapter creates a new sample adapter
func NewSampleAdapter() *SampleAdapter {
	adapter := &SampleAdapter{
		name:         "sample_supplier",
		client:       &http.Client{Timeout: 30 * time.Second},
		config:       make(map[string]string),
		capabilities: make(map[string]bool),
	}

	// Set capabilities
	adapter.capabilities["sync_products"] = true
	adapter.capabilities["sync_inventory"] = true
	adapter.capabilities["real_time_inventory"] = false

	return adapter
}

// Name returns the name of the adapter
func (a *SampleAdapter) Name() string {
	return a.name
}

// Initialize initializes the adapter with supplier-specific configuration
func (a *SampleAdapter) Initialize(ctx context.Context, config map[string]string) error {
	// Store the configuration
	a.config = config

	// Sample adapter requires API key and URL
	requiredKeys := []string{"api_key", "api_url"}
	if err := a.validateRequiredConfig(requiredKeys); err != nil {
		return err
	}

	// Extract API URL and key from config
	a.apiURL = a.config["api_url"]
	a.apiKey = a.config["api_key"]

	return nil
}

// validateRequiredConfig checks if all required config keys are present
func (a *SampleAdapter) validateRequiredConfig(requiredKeys []string) error {
	missingKeys := []string{}

	for _, key := range requiredKeys {
		if _, ok := a.config[key]; !ok {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required configuration keys: %v", missingKeys)
	}

	return nil
}

// GetCapabilities returns the capabilities of this adapter
func (a *SampleAdapter) GetCapabilities(ctx context.Context) map[string]bool {
	return a.capabilities
}

// TestConnection tests the connection to the supplier's system
func (a *SampleAdapter) TestConnection(ctx context.Context) error {

	// Create a test request to the supplier's health endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", a.apiURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	// Execute the request
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection test failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// GetProducts fetches products from the supplier
func (a *SampleAdapter) GetProducts(ctx context.Context, options domain.SupplierSyncOptions) ([]domain.SupplierProductData, error) {

	// Construct the URL with query parameters based on options
	url := fmt.Sprintf("%s/products?full=%t&batch_size=%d",
		a.apiURL,
		options.FullSync,
		options.BatchSize,
	)

	// Add date filters if specified
	if !options.FromDate.IsZero() {
		url += "&from_date=" + options.FromDate.Format(time.RFC3339)
	}
	if !options.ToDate.IsZero() {
		url += "&to_date=" + options.ToDate.Format(time.RFC3339)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	// Execute the request
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch products with status code: %d", resp.StatusCode)
	}

	// Parse the response
	var supplierResponse struct {
		Products []struct {
			ID          string            `json:"id"`
			Name        string            `json:"name"`
			SKU         string            `json:"sku"`
			Description string            `json:"description"`
			Price       float64           `json:"price"`
			Currency    string            `json:"currency"`
			Stock       int32             `json:"stock"`
			MinOrder    int32             `json:"min_order"`
			LeadTime    int32             `json:"lead_time"`
			Categories  []string          `json:"categories"`
			Barcode     string            `json:"barcode"`
			Weight      float64           `json:"weight"`
			Length      float64           `json:"length"`
			Width       float64           `json:"width"`
			Height      float64           `json:"height"`
			DimUnit     string            `json:"dim_unit"`
			Images      []string          `json:"images"`
			Active      bool              `json:"active"`
			UpdatedAt   string            `json:"updated_at"`
			Extra       map[string]string `json:"extra"`
		} `json:"products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&supplierResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to our domain model
	products := make([]domain.SupplierProductData, 0, len(supplierResponse.Products))
	for _, p := range supplierResponse.Products {
		// Parse the updated timestamp
		updatedAt, err := time.Parse(time.RFC3339, p.UpdatedAt)
		if err != nil {
			// Use current time as fallback
			updatedAt = time.Now()
		}

		product := domain.SupplierProductData{
			ExternalID:      p.ID,
			Name:            p.Name,
			SKU:             p.SKU,
			Description:     p.Description,
			Price:           p.Price,
			Currency:        p.Currency,
			StockQuantity:   p.Stock,
			MinimumOrderQty: p.MinOrder,
			LeadTimeDays:    p.LeadTime,
			Categories:      p.Categories,
			Barcode:         p.Barcode,
			Weight:          p.Weight,
			Dimensions: &domain.ProductDimensions{
				Length: p.Length,
				Width:  p.Width,
				Height: p.Height,
				Unit:   p.DimUnit,
			},
			Images:            p.Images,
			Active:            p.Active,
			LastUpdated:       updatedAt,
			AdditionalDetails: p.Extra,
		}
		products = append(products, product)
	}

	return products, nil
}

// GetInventory fetches inventory data from the supplier
func (a *SampleAdapter) GetInventory(ctx context.Context, externalIDs []string, options domain.SupplierSyncOptions) ([]domain.SupplierInventoryData, error) {

	// Create request body with product IDs
	idsJSON, err := json.Marshal(externalIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal product IDs: %w", err)
	}

	// Create the request with JSON body
	url := fmt.Sprintf("%s/inventory", a.apiURL)
	reqBody := bytes.NewBuffer(idsJSON)
	req, err := http.NewRequestWithContext(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add request parameters
	q := req.URL.Query()
	if options.FullSync {
		q.Add("full", "true")
	}
	req.URL.RawQuery = q.Encode()

	// Add authentication
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch inventory: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch inventory with status code: %d", resp.StatusCode)
	}

	// Parse the response
	var supplierResponse struct {
		Inventory []struct {
			ProductID    string `json:"product_id"`
			StockLevel   int32  `json:"stock_level"`
			Reserved     int32  `json:"reserved"`
			Available    int32  `json:"available"`
			Incoming     int32  `json:"incoming"`
			ExpectedDate string `json:"expected_date"`
			Location     string `json:"location"`
			UpdatedAt    string `json:"updated_at"`
		} `json:"inventory"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&supplierResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to domain model
	inventory := make([]domain.SupplierInventoryData, 0, len(supplierResponse.Inventory))
	for _, i := range supplierResponse.Inventory {
		// Parse the dates
		var expectedDate time.Time
		var lastUpdate time.Time

		if i.ExpectedDate != "" {
			expectedDate, _ = time.Parse(time.RFC3339, i.ExpectedDate)
		}

		if i.UpdatedAt != "" {
			lastUpdate, _ = time.Parse(time.RFC3339, i.UpdatedAt)
		} else {
			lastUpdate = time.Now()
		}

		item := domain.SupplierInventoryData{
			ExternalID:     i.ProductID,
			StockQuantity:  i.StockLevel,
			ReservedQty:    i.Reserved,
			AvailableQty:   i.Available,
			IncomingQty:    i.Incoming,
			ExpectedDate:   expectedDate,
			LocationCode:   i.Location,
			LastUpdateTime: lastUpdate,
		}
		inventory = append(inventory, item)
	}

	return inventory, nil
}

// SyncProducts syncs product data from the supplier
func (a *SampleAdapter) SyncProducts(ctx context.Context, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error) {
	stats := &domain.SupplierSyncStats{
		StartTime: time.Now(),
	}

	// Fetch products from supplier
	products, err := a.GetProducts(ctx, options)
	if err != nil {
		stats.EndTime = time.Now()
		stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to fetch products: %v", err))
		return stats, err
	}

	// Process each product
	for _, product := range products {
		stats.ProductsProcessed++

		// Here you would typically:
		// 1. Check if product already exists in your system
		// 2. Create or update the product in your system
		// 3. Handle any mapping between supplier and system product data

		// Simulating success/failure for demonstration
		if product.Active {
			stats.ProductsUpdated++
		} else {
			stats.ProductsErrored++
			stats.Errors = append(stats.Errors, fmt.Sprintf("Product %s is inactive", product.ExternalID))
		}
	}

	stats.EndTime = time.Now()
	return stats, nil
}

// SyncInventory syncs inventory data from the supplier
func (a *SampleAdapter) SyncInventory(ctx context.Context, options domain.SupplierSyncOptions) (*domain.SupplierSyncStats, error) {
	stats := &domain.SupplierSyncStats{
		StartTime: time.Now(),
	}

	// In a real implementation, you would:
	// 1. Fetch all product IDs that need inventory updates
	// 2. Get inventory data for those products
	// 3. Update inventory in your system

	// For this example, we'll simulate a simple inventory sync
	externalIDs := []string{"sample-product-1", "sample-product-2"}

	inventory, err := a.GetInventory(ctx, externalIDs, options)
	if err != nil {
		stats.EndTime = time.Now()
		stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to fetch inventory: %v", err))
		return stats, err
	}

	// Process inventory data
	for _, item := range inventory {
		stats.InventoryProcessed++

		// Here you would update inventory in your system
		// Simulating success for demonstration
		if item.StockQuantity > 0 {
			stats.InventoryUpdated++
		} else {
			stats.InventoryErrored++
			stats.Errors = append(stats.Errors, fmt.Sprintf("Product %s has zero stock", item.ExternalID))
		}
	}

	stats.EndTime = time.Now()
	return stats, nil
}
