package product

import (
	"strconv"
	"time"

	"github.com/leonvanderhaeghen/stockplatform/pkg/models"
	productv1 "github.com/leonvanderhaeghen/stockplatform/services/productSvc/api/gen/go/proto/product/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertToProduct converts protobuf Product to domain Product
func convertToProduct(protoProduct *productv1.Product) *models.Product {
	if protoProduct == nil {
		return nil
	}

	// Parse cost and selling prices from strings
	costPrice, _ := strconv.ParseFloat(protoProduct.CostPrice, 64)
	sellingPrice, _ := strconv.ParseFloat(protoProduct.SellingPrice, 64)

	// Use first category if multiple categories exist
	category := ""
	if len(protoProduct.CategoryIds) > 0 {
		category = protoProduct.CategoryIds[0]
	}

	return &models.Product{
		ID:          protoProduct.Id,
		Name:        protoProduct.Name,
		Description: protoProduct.Description,
		SKU:         protoProduct.Sku,
		Price:       sellingPrice,
		Cost:        costPrice,
		Category:    category,
		Brand:       "", // Not available in protobuf schema
		Weight:      0,  // Not available in protobuf schema
		Dimensions:  nil, // Not available in protobuf schema
		IsActive:    protoProduct.IsActive,
		SupplierID:  protoProduct.SupplierId,
		CreatedAt:   convertTimestamp(protoProduct.CreatedAt),
		UpdatedAt:   convertTimestamp(protoProduct.UpdatedAt),
	}
}

// convertToProtoProduct converts domain Product to protobuf Product
func convertToProtoProduct(product *models.Product) *productv1.Product {
	if product == nil {
		return nil
	}

	return &productv1.Product{
		Id:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		CostPrice:    strconv.FormatFloat(product.Cost, 'f', 2, 64),
		SellingPrice: strconv.FormatFloat(product.Price, 'f', 2, 64),
		Currency:     "USD", // Default currency
		Sku:          product.SKU,
		CategoryIds:  []string{product.Category}, // Convert single category to slice
		SupplierId:   product.SupplierID,
		IsActive:     product.IsActive,
		CreatedAt:    timestamppb.New(product.CreatedAt),
		UpdatedAt:    timestamppb.New(product.UpdatedAt),
	}
}

// convertToCreateProductResponse converts protobuf CreateProductResponse to domain CreateProductResponse
func convertToCreateProductResponse(resp *productv1.CreateProductResponse) *models.CreateProductResponse {
	if resp == nil {
		return nil
	}

	return &models.CreateProductResponse{
		Product: convertToProduct(resp.Product),
		Message: "Product created successfully",
	}
}

// convertToListProductsResponse converts protobuf ListProductsResponse to domain ListProductsResponse
func convertToListProductsResponse(resp *productv1.ListProductsResponse) *models.ListProductsResponse {
	if resp == nil {
		return nil
	}

	products := make([]*models.Product, len(resp.Products))
	for i, protoProduct := range resp.Products {
		products[i] = convertToProduct(protoProduct)
	}

	return &models.ListProductsResponse{
		Products:   products,
		TotalCount: resp.TotalCount,
	}
}

// convertTimestamp converts protobuf timestamp to time.Time
func convertTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// convertToCreateProductRequest converts domain parameters to protobuf CreateProductRequest
func convertToCreateProductRequest(name, description, sku, supplierID string, costPrice, sellingPrice float64, isActive bool, categoryIDs []string) *productv1.CreateProductRequest {
	return &productv1.CreateProductRequest{
		Name:         name,
		Description:  description,
		CostPrice:    strconv.FormatFloat(costPrice, 'f', 2, 64),
		SellingPrice: strconv.FormatFloat(sellingPrice, 'f', 2, 64),
		Currency:     "USD", // Default currency
		Sku:          sku,
		CategoryIds:  categoryIDs,
		SupplierId:   supplierID,
		IsActive:     isActive,
	}
}
