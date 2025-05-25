package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InventoryRequest represents the inventory request body
type InventoryRequest struct {
	ProductID  string  `json:"productId" binding:"required"`
	SKU        string  `json:"sku" binding:"required"`
	Quantity   int32   `json:"quantity" binding:"required,gte=0"`
	Location   string  `json:"location"`
	ReorderAt  int32   `json:"reorderAt" binding:"gte=0"`
	ReorderQty int32   `json:"reorderQty" binding:"gte=0"`
	Cost       float64 `json:"cost" binding:"gte=0"`
}

// StockAdjustRequest represents the stock adjustment request body
type StockAdjustRequest struct {
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
	Reason    string `json:"reason"`
	Reference string `json:"reference"`
}

// listInventory returns a list of inventory items
func (s *Server) listInventory(c *gin.Context) {
	location := c.Query("location")
	lowStock := c.Query("lowStock") == "true"
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

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

	items, err := s.inventorySvc.ListInventory(c.Request.Context(), location, lowStock, limit, offset)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "List inventory")
		return
	}

	respondWithSuccess(c, http.StatusOK, items)
}

// getInventoryItem returns an inventory item by ID
func (s *Server) getInventoryItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Inventory item ID is required")
		return
	}

	item, err := s.inventorySvc.GetInventoryItemByID(c.Request.Context(), id)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get inventory item")
		return
	}

	respondWithSuccess(c, http.StatusOK, item)
}

// getInventoryItemByProduct returns inventory items for a product
func (s *Server) getInventoryItemByProduct(c *gin.Context) {
	productID := c.Param("productId")
	if productID == "" {
		respondWithError(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	items, err := s.inventorySvc.GetInventoryItemsByProductID(c.Request.Context(), productID)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get inventory by product")
		return
	}

	respondWithSuccess(c, http.StatusOK, items)
}

// getInventoryItemBySKU returns an inventory item by SKU
func (s *Server) getInventoryItemBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		respondWithError(c, http.StatusBadRequest, "SKU is required")
		return
	}

	item, err := s.inventorySvc.GetInventoryItemBySKU(c.Request.Context(), sku)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get inventory by SKU")
		return
	}

	respondWithSuccess(c, http.StatusOK, item)
}

// createInventoryItem creates a new inventory item
func (s *Server) createInventoryItem(c *gin.Context) {
	var req InventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	item, err := s.inventorySvc.CreateInventoryItem(
		c.Request.Context(),
		req.ProductID,
		req.SKU,
		req.Quantity,
		req.Location,
		req.ReorderAt,
		req.ReorderQty,
		req.Cost,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Create inventory item")
		return
	}

	respondWithSuccess(c, http.StatusCreated, item)
}

// updateInventoryItem updates an existing inventory item
func (s *Server) updateInventoryItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Inventory item ID is required")
		return
	}

	var req InventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := s.inventorySvc.UpdateInventoryItem(
		c.Request.Context(),
		id,
		req.ProductID,
		req.SKU,
		req.Quantity,
		req.Location,
		req.ReorderAt,
		req.ReorderQty,
		req.Cost,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Update inventory item")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Inventory item updated successfully"})
}

// deleteInventoryItem deletes an inventory item
func (s *Server) deleteInventoryItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Inventory item ID is required")
		return
	}

	err := s.inventorySvc.DeleteInventoryItem(c.Request.Context(), id)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Delete inventory item")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Inventory item deleted successfully"})
}

// addStock adds stock to an inventory item
func (s *Server) addStock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Inventory item ID is required")
		return
	}

	var req StockAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	item, err := s.inventorySvc.AddStock(
		c.Request.Context(),
		id,
		req.Quantity,
		req.Reason,
		req.Reference,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Add stock")
		return
	}

	respondWithSuccess(c, http.StatusOK, item)
}

// removeStock removes stock from an inventory item
func (s *Server) removeStock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Inventory item ID is required")
		return
	}

	var req StockAdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	item, err := s.inventorySvc.RemoveStock(
		c.Request.Context(),
		id,
		req.Quantity,
		req.Reason,
		req.Reference,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Remove stock")
		return
	}

	respondWithSuccess(c, http.StatusOK, item)
}
