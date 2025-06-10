package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryRequest represents the category request body
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id"`
	IsActive    bool   `json:"is_active"`
}

// ProductRequest represents the product request body
type ProductRequest struct {
	Name         string            `json:"name" binding:"required"`
	Description  string            `json:"description"`
	CostPrice    string            `json:"cost_price" binding:"required"`
	SellingPrice string            `json:"selling_price" binding:"required"`
	Currency     string            `json:"currency"`
	SKU          string            `json:"sku" binding:"required"`
	Barcode      string            `json:"barcode"`
	CategoryIDs  []string          `json:"category_ids"`
	SupplierID   string            `json:"supplier_id"`
	IsActive     bool              `json:"is_active"`
	InStock      bool              `json:"in_stock"`
	StockQty     int32             `json:"stock_qty"`
	LowStockAt   int32             `json:"low_stock_at"`
	ImageURLs    []string          `json:"image_urls"`
	VideoURLs    []string          `json:"video_urls"`
	Metadata     map[string]string `json:"metadata"`
}

// listCategories returns a list of product categories
func (s *Server) listCategories(c *gin.Context) {
	categories, err := s.productSvc.ListCategories(c.Request.Context())
	if err != nil {
		genericErrorHandler(c, err, s.logger, "List categories")
		return
	}

	respondWithSuccess(c, http.StatusOK, categories)
}

// createCategory creates a new product category
func (s *Server) createCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	category, err := s.productSvc.CreateCategory(c.Request.Context(), req.Name, req.Description, req.ParentID, req.IsActive)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Create category")
		return
	}

	respondWithSuccess(c, http.StatusCreated, category)
}

// listProducts returns a list of products
func (s *Server) listProducts(c *gin.Context) {
	categoryID := c.Query("category")
	query := c.Query("q")
	activeStr := c.DefaultQuery("active", "true")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	sortBy := c.DefaultQuery("sort", "name")
	orderStr := c.DefaultQuery("order", "asc")

	// Parse parameters
	active := activeStr == "true"
	limit, err := parseIntParam(limitStr, 10)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	offset, err := parseIntParam(offsetStr, 0)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	order := orderStr == "asc"

	products, err := s.productSvc.ListProducts(c.Request.Context(), categoryID, query, active, limit, offset, sortBy, order)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "List products")
		return
	}

	respondWithSuccess(c, http.StatusOK, products)
}

// getProduct returns a product by ID
func (s *Server) getProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	product, err := s.productSvc.GetProductByID(c.Request.Context(), id)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get product")
		return
	}

	respondWithSuccess(c, http.StatusOK, product)
}

// createProduct creates a new product
func (s *Server) createProduct(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Validate price formats
	if _, err := strconv.ParseFloat(req.CostPrice, 64); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid cost price format")
		return
	}

	if _, err := strconv.ParseFloat(req.SellingPrice, 64); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid selling price format")
		return
	}

	// Convert the request to the format expected by the product service
	product, err := s.productSvc.CreateProduct(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.CostPrice,
		req.SellingPrice,
		req.Currency,
		req.SKU,
		req.Barcode,
		req.CategoryIDs,
		req.SupplierID,
		req.IsActive,
		req.InStock,
		req.StockQty,
		req.LowStockAt,
		req.ImageURLs,
		req.VideoURLs,
		req.Metadata,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Create product")
		return
	}

	respondWithSuccess(c, http.StatusCreated, product)
}

// updateProduct updates an existing product
func (s *Server) updateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Validate price formats
	if _, err := strconv.ParseFloat(req.SellingPrice, 64); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid selling price format")
		return
	}

	if _, err := strconv.ParseFloat(req.CostPrice, 64); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid cost price format")
		return
	}

	// Call the service to update the product
	err := s.productSvc.UpdateProduct(
		c.Request.Context(),
		id,
		req.Name,
		req.Description,
		req.SKU,
		req.CategoryIDs,
		req.SellingPrice,
		req.CostPrice,
		req.IsActive,
		req.ImageURLs,
		req.Metadata,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Update product")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// deleteProduct deletes a product
func (s *Server) deleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	err := s.productSvc.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Delete product")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// parseIntParam parses a string parameter to an integer with a default value
func parseIntParam(param string, defaultValue int) (int, error) {
	if param == "" {
		return defaultValue, nil
	}
	
	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue, err
	}
	
	return value, nil
}
