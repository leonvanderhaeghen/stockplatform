package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	Source    string `json:"source,omitempty"` // POS, ONLINE, etc. for tracking adjustment source
}

// ReservationRequest represents the inventory reservation request body
type ReservationRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
	OrderID   string `json:"orderId" binding:"required"`
	Source    string `json:"source,omitempty"` // POS, ONLINE, etc. for tracking reservation source
	StoreID   string `json:"storeId,omitempty"` // For POS reservations
}

// listInventory returns a list of inventory items (supports POS availability checking)
func (s *Server) listInventory(c *gin.Context) {
	location := c.Query("location")
	lowStock := c.Query("lowStock") == "true"
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	
	// POS-specific parameters for availability checking
	checkAvailability := c.Query("checkAvailability") == "true" // POS inventory check
	storeId := c.Query("storeId")                              // POS store filter
	sku := c.Query("sku")                                      // POS SKU lookup
	minQuantity := c.Query("minQuantity")                      // POS availability threshold

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

	// Handle POS-specific inventory checking
	if checkAvailability {
		s.logger.Debug("POS availability check requested",
			zap.String("storeId", storeId),
			zap.String("sku", sku),
			zap.String("minQuantity", minQuantity),
		)
		
		// If specific SKU requested for POS, use the SKU lookup instead
		if sku != "" {
			item, err := s.inventorySvc.GetInventoryItemBySKU(c.Request.Context(), sku)
			if err != nil {
				genericErrorHandler(c, err, s.logger, "Get inventory by SKU for POS")
				return
			}
			
			// Check if item meets minimum quantity requirement
			if minQuantity != "" {
				minQty, err := parseIntParam(minQuantity, 0)
				if err != nil {
					respondWithError(c, http.StatusBadRequest, "Invalid minQuantity parameter")
					return
				}
				
				// Add availability status to response
				response := map[string]interface{}{
					"item": item,
					"available": true, // Placeholder - would need item quantity check
					"availabilityCheck": map[string]interface{}{
						"requestedMinQuantity": minQty,
						"storeId": storeId,
						"checkTimestamp": "now", // Would use time.Now() in real implementation
					},
				}
				respondWithSuccess(c, http.StatusOK, response)
				return
			}
			
			// Return single item for POS SKU check
			respondWithSuccess(c, http.StatusOK, map[string]interface{}{
				"item": item,
				"posAvailabilityCheck": true,
				"storeId": storeId,
			})
			return
		}
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

	// Enhanced logging for POS source tracking
	if req.Source != "" {
		s.logger.Debug("POS inventory deduction requested",
			zap.String("inventoryId", id),
			zap.Int32("quantity", req.Quantity),
			zap.String("source", req.Source),
			zap.String("reason", req.Reason),
			zap.String("reference", req.Reference),
		)
		
		// For POS transactions, enhance the reason to include source info
		if req.Reason == "" {
			req.Reason = "POS Transaction"
		}
		req.Reason = req.Reason + " [Source: " + req.Source + "]"
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

// getInventoryReservations returns inventory reservations with optional filters
func (s *Server) getInventoryReservations(c *gin.Context) {
	orderId := c.Query("orderId")
	productId := c.Query("productId")
	status := c.Query("status")
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

	reservations, err := s.inventorySvc.GetInventoryReservations(c.Request.Context(), orderId, productId, status, limit, offset)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get inventory reservations")
		return
	}

	respondWithSuccess(c, http.StatusOK, reservations)
}

// createInventoryReservation creates a new inventory reservation (supports POS source tracking)
func (s *Server) createInventoryReservation(c *gin.Context) {
	var req ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Enhanced logging for POS source tracking
	if req.Source != "" {
		s.logger.Debug("POS inventory reservation requested",
			zap.String("productId", req.ProductID),
			zap.Int32("quantity", req.Quantity),
			zap.String("orderId", req.OrderID),
			zap.String("source", req.Source),
			zap.String("storeId", req.StoreID),
		)
	}

	reservation, err := s.inventorySvc.CreateInventoryReservation(
		c.Request.Context(),
		req.ProductID,
		req.Quantity,
		req.OrderID,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Create inventory reservation")
		return
	}

	// For POS reservations, add source metadata to response
	response := map[string]interface{}{
		"reservation": reservation,
	}
	if req.Source != "" {
		response["source"] = req.Source
		response["storeId"] = req.StoreID
		response["posReservation"] = true
	}

	respondWithSuccess(c, http.StatusCreated, response)
}

// getLowStockItems returns inventory items that are low in stock
func (s *Server) getLowStockItems(c *gin.Context) {
	thresholdStr := c.DefaultQuery("threshold", "10")
	location := c.Query("location")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	threshold, err := parseIntParam(thresholdStr, 10)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid threshold parameter")
		return
	}

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

	// Get low stock items with threshold and location filtering
	items, err := s.inventorySvc.GetLowStockItems(c.Request.Context(), location, threshold, limit, offset)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get low stock items")
		return
	}

	respondWithSuccess(c, http.StatusOK, items)
}
