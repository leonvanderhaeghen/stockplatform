package application

import (
	"context"
	"fmt"

	inventoryclient "github.com/leonvanderhaeghen/stockplatform/pkg/clients/inventory"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// ProductInventoryService coordinates between product and inventory services
type ProductInventoryService struct {
	inventoryClient *inventoryclient.Client
	productService  *ProductService
	logger         *zap.Logger
}

// NewProductInventoryService creates a new ProductInventoryService
func NewProductInventoryService(
	productService *ProductService,
	inventoryAddr string,
	logger *zap.Logger,
) (*ProductInventoryService, error) {
	// Initialize the inventory client
	invCfg := inventoryclient.Config{Address: inventoryAddr}
	inventoryClient, err := inventoryclient.New(invCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory client: %w", err)
	}

	return &ProductInventoryService{
		inventoryClient: inventoryClient,
		productService:  productService,
		logger:         logger.Named("product_inventory_service"),
	}, nil
}

// Close closes any open connections
func (s *ProductInventoryService) Close() error {
	if s.inventoryClient != nil {
		return s.inventoryClient.Close()
	}
	return nil
}

// CreateProductWithInventory creates a new product with initial inventory
func (s *ProductInventoryService) CreateProductWithInventory(
	ctx context.Context,
	input *domain.Product,
	initialStock int32,
) (*domain.Product, error) {
	// First, create the product
	product, err := s.productService.CreateProduct(ctx, input)
	if err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Then, create inventory for the product
	_, err = s.inventoryClient.CreateInventory(ctx, &inventorypb.CreateInventoryRequest{
		ProductId: product.ID.Hex(),
		Quantity:  initialStock,
		Sku:       product.SKU,
	})

	if err != nil {
		s.logger.Error("Failed to create inventory", 
			zap.String("product_id", product.ID.Hex()),
			zap.Error(err))
		
		// Try to clean up the product if inventory creation fails
		if delErr := s.productService.DeleteProduct(ctx, product.ID.Hex()); delErr != nil {
			s.logger.Error("Failed to clean up product after inventory creation failure",
				zap.String("product_id", product.ID.Hex()),
				zap.Error(delErr))
		}
		
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	s.logger.Info("Successfully created product with inventory",
		zap.String("product_id", product.ID.Hex()),
		zap.Int32("initial_stock", initialStock))

	return product, nil
}

// GetProductWithInventory retrieves a product along with its inventory information
func (s *ProductInventoryService) GetProductWithInventory(
	ctx context.Context,
	productID string,
) (*domain.Product, *inventorypb.InventoryItem, error) {
	// Get the product
	product, err := s.productService.GetProduct(ctx, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Get the inventory for this product
	inventoryResp, err := s.inventoryClient.GetInventoryByProductID(ctx, &inventorypb.GetInventoryByProductIDRequest{
		ProductId: productID,
	})

	// Handle case where inventory doesn't exist
	if err != nil {
		s.logger.Warn("No inventory found for product", 
			zap.String("product_id", productID),
			zap.Error(err))
		return product, nil, nil
	}

	return product, inventoryResp.GetInventory(), nil
}

// UpdateProductInventory updates a product and its inventory
func (s *ProductInventoryService) UpdateProductInventory(
	ctx context.Context,
	productID string,
	productUpdate *domain.Product,
	stockAdjustment *int32,
) (*domain.Product, *inventorypb.InventoryItem, error) {
	// Update the product
	productUpdate.ID, _ = primitive.ObjectIDFromHex(productID)
	err := s.productService.UpdateProduct(ctx, productUpdate)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update product: %w", err)
	}
	
	// Get the updated product
	updatedProduct, err := s.productService.GetProduct(ctx, productID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get updated product: %w", err)
	}

	var inventory *inventorypb.InventoryItem

	// If stock adjustment is provided, update the inventory
	if stockAdjustment != nil {
		// Get current inventory
		inventoryResp, err := s.inventoryClient.GetInventoryByProductID(ctx, &inventorypb.GetInventoryByProductIDRequest{
			ProductId: productID,
		})

		// If inventory doesn't exist and we have a positive stock adjustment, create it
		if err != nil && *stockAdjustment > 0 {
			createResp, err := s.inventoryClient.CreateInventory(ctx, &inventorypb.CreateInventoryRequest{
				ProductId: productID,
				Quantity:  *stockAdjustment,
				Sku:       updatedProduct.SKU,
			})
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create inventory: %w", err)
			}
			inventory = createResp.GetInventory()
		} else if err == nil {
			// Update existing inventory
			newQuantity := inventoryResp.GetInventory().GetQuantity() + *stockAdjustment
			if newQuantity < 0 {
				return nil, nil, fmt.Errorf("insufficient stock")
			}

			_, err = s.inventoryClient.UpdateInventory(ctx, &inventorypb.UpdateInventoryRequest{
				Inventory: &inventorypb.InventoryItem{
					Id:        inventoryResp.GetInventory().GetId(),
					ProductId: productID,
					Quantity:  newQuantity,
					Sku:       updatedProduct.SKU,
				},
			})
			if err != nil {
				return nil, nil, fmt.Errorf("failed to update inventory: %w", err)
			}

			inventory = &inventorypb.InventoryItem{
				Id:        inventoryResp.GetInventory().GetId(),
				ProductId: productID,
				Quantity:  newQuantity,
				Sku:       updatedProduct.SKU,
			}
		}
	}

	return updatedProduct, inventory, nil
}
