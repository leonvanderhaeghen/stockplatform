package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductRequest represents the product request body
type ProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	SKU         string   `json:"sku" binding:"required"`
	Categories  []string `json:"categories"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Cost        float64  `json:"cost" binding:"required,gte=0"`
	Active      bool     `json:"active"`
	Images      []string `json:"images"`
	Attributes  map[string]string `json:"attributes"`
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

	product, err := s.productSvc.CreateProduct(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.SKU,
		req.Categories,
		req.Price,
		req.Cost,
		req.Active,
		req.Images,
		req.Attributes,
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

	err := s.productSvc.UpdateProduct(
		c.Request.Context(),
		id,
		req.Name,
		req.Description,
		req.SKU,
		req.Categories,
		req.Price,
		req.Cost,
		req.Active,
		req.Images,
		req.Attributes,
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
