package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// POSOrderRequest represents a POS order request
type POSOrderRequest struct {
	Items       []OrderItemRequest `json:"items" binding:"required,min=1"`
	LocationID  string             `json:"locationId" binding:"required"`
	StaffID     string             `json:"staffId" binding:"required"`
	CustomerID  string             `json:"customerId"`
	PaymentType string             `json:"paymentType" binding:"required"`
	PaymentData map[string]string  `json:"paymentData"`
	Notes       string             `json:"notes"`
}

// POSQuickTransactionRequest represents a quick POS transaction that creates and pays for an order in one call
type POSQuickTransactionRequest struct {
	Items       []OrderItemRequest `json:"items" binding:"required,min=1"`
	LocationID  string             `json:"locationId" binding:"required"`
	StaffID     string             `json:"staffId" binding:"required"`
	CustomerID  string             `json:"customerId"`
	Amount      float64            `json:"amount" binding:"required,gt=0"`
	PaymentType string             `json:"paymentType" binding:"required"`
	PaymentData map[string]string  `json:"paymentData"`
	Notes       string             `json:"notes"`
}

// POSInventoryCheckRequest represents a request to check inventory availability
type POSInventoryCheckRequest struct {
	LocationID string                 `json:"locationId" binding:"required"`
	Items      []POSInventoryCheckItem `json:"items" binding:"required,min=1"`
}

// POSInventoryCheckItem represents an item to check in inventory
type POSInventoryCheckItem struct {
	ProductID string `json:"productId"`
	SKU       string `json:"sku"`
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
}

// POSReservationRequest represents a reservation request for in-store pickup
type POSReservationRequest struct {
	OrderID    string                `json:"orderId" binding:"required"`
	LocationID string                `json:"locationId" binding:"required"`
	StaffID    string                `json:"staffId" binding:"required"`
	Items      []POSReservationItem  `json:"items" binding:"required,min=1"`
}

// POSReservationItem represents an item to be reserved
type POSReservationItem struct {
	ProductID string `json:"productId"`
	SKU       string `json:"sku"`
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
}

// POSPickupCompletionRequest represents completing an order pickup
type POSPickupCompletionRequest struct {
	OrderID    string `json:"orderId" binding:"required"`
	LocationID string `json:"locationId" binding:"required"`
	StaffID    string `json:"staffId" binding:"required"`
	Notes      string `json:"notes"`
}

// POSDirectSaleRequest represents a direct POS sale without order creation
type POSDirectSaleRequest struct {
	LocationID string                `json:"locationId" binding:"required"`
	StaffID    string                `json:"staffId" binding:"required"`
	Items      []POSInventoryCheckItem `json:"items" binding:"required,min=1"`
	Notes      string                `json:"notes"`
}

// createPOSOrder creates a new order from a POS terminal
func (s *Server) createPOSOrder(c *gin.Context) {
	// Validate staff credentials and auth
	staffID, exists := c.Get("userID")
	staffIDStr, ok := staffID.(string)
	if !ok || !exists {
		respondWithError(c, http.StatusUnauthorized, "Staff ID required")
		return
	}

	// Get role to ensure staff/admin permissions
	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok || (roleStr != "staff" && roleStr != "admin") {
		respondWithError(c, http.StatusForbidden, "Staff or admin role required")
		return
	}

	var req POSOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert request items to service format
	items := make([]map[string]interface{}, len(req.Items))
	for i, item := range req.Items {
		items[i] = map[string]interface{}{
			"product_id": item.ProductID,
			"sku":        item.SKU,
			"quantity":   item.Quantity,
			"price":      item.Price,
		}
	}

	// Use staff ID from authenticated user if not provided in request
	if req.StaffID == "" {
		req.StaffID = staffIDStr
	}

	// Create and process POS order
	orderID, err := s.orderSvc.CreatePOSOrder(
		c.Request.Context(),
		staffIDStr,
		items,
		req.CustomerID,
		req.LocationID,
		"", // Empty payment method
		map[string]string{}, // Empty payment details
		req.Notes,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Create POS order")
		return
	}

	s.logger.Info("Created POS order",
		zap.String("order_id", orderID.(string)),
		zap.String("location_id", req.LocationID),
		zap.String("staff_id", req.StaffID),
	)

	respondWithSuccess(c, http.StatusCreated, gin.H{
		"message":  "POS order created successfully",
		"order_id": orderID,
	})
}

// processQuickPOSTransaction processes a quick POS transaction
func (s *Server) processQuickPOSTransaction(c *gin.Context) {
	// Validate staff credentials and auth
	staffID, exists := c.Get("userID")
	staffIDStr, ok := staffID.(string)
	if !ok || !exists {
		respondWithError(c, http.StatusUnauthorized, "Staff ID required")
		return
	}

	// Get role to ensure staff/admin permissions
	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok || (roleStr != "staff" && roleStr != "admin") {
		respondWithError(c, http.StatusForbidden, "Staff or admin role required")
		return
	}

	var req POSQuickTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert request items to service format
	items := make([]map[string]interface{}, len(req.Items))
	for i, item := range req.Items {
		items[i] = map[string]interface{}{
			"product_id": item.ProductID,
			"sku":        item.SKU,
			"quantity":   item.Quantity,
			"price":      item.Price,
		}
	}

	// Use staff ID from authenticated user if not provided
	if req.StaffID == "" {
		req.StaffID = staffIDStr
	}

	// Process quick POS transaction
	result, err := s.orderSvc.ProcessQuickPOSTransaction(
		c.Request.Context(),
		staffIDStr,
		req.LocationID,
		items,
		map[string]interface{}{
			"customerID": req.CustomerID,
			"totalAmount": req.Amount,
			"paymentMethod": req.PaymentType,
			"paymentDetails": req.PaymentData,
			"notes": req.Notes,
		},
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Process quick POS transaction")
		return
	}

	s.logger.Info("Processed quick POS transaction",
		zap.String("order_id", result.(string)),
		zap.String("location_id", req.LocationID),
		zap.String("staff_id", req.StaffID),
		zap.Float64("amount", req.Amount),
	)

	respondWithSuccess(c, http.StatusCreated, gin.H{
		"message":  "POS transaction processed successfully",
		"order_id": result,
	})
}

// checkPOSInventory checks inventory availability for POS transactions
func (s *Server) checkPOSInventory(c *gin.Context) {
	// Validate staff credentials and auth
	_, exists := c.Get("userID")
	if !exists {
		respondWithError(c, http.StatusUnauthorized, "Staff authentication required")
		return
	}

	// Get role to ensure staff/admin permissions
	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok || (roleStr != "staff" && roleStr != "admin") {
		respondWithError(c, http.StatusForbidden, "Staff or admin role required")
		return
	}

	var req POSInventoryCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert request items to domain format
	items := make([]map[string]interface{}, len(req.Items))
	for i, item := range req.Items {
		checkItem := map[string]interface{}{
			"product_id": item.ProductID,
			"sku":        item.SKU,
			"quantity":   item.Quantity,
		}
		items[i] = checkItem
	}

	// Check inventory
	results, err := s.inventorySvc.PerformPOSInventoryCheck(
		c.Request.Context(),
		req.LocationID,
		items,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Check POS inventory")
		return
	}

	respondWithSuccess(c, http.StatusOK, results)
}

// reserveForPOSTransaction reserves inventory for a POS transaction
func (s *Server) reserveForPOSTransaction(c *gin.Context) {
	// Validate staff credentials and auth
	_, exists := c.Get("userID")
	if !exists {
		respondWithError(c, http.StatusUnauthorized, "Staff authentication required")
		return
	}

	// Get role to ensure staff/admin permissions
	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok || (roleStr != "staff" && roleStr != "admin") {
		respondWithError(c, http.StatusForbidden, "Staff or admin role required")
		return
	}

	var req POSReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert request items to domain format
	items := make([]map[string]interface{}, len(req.Items))
	for i, item := range req.Items {
		reserveItem := map[string]interface{}{
			"product_id": item.ProductID,
			"sku":        item.SKU,
			"quantity":   item.Quantity,
		}
		items[i] = reserveItem
	}

	// Reserve inventory items for POS transaction
	_, err := s.inventorySvc.ReserveForPOSTransaction(
		c.Request.Context(),
		req.LocationID,
		req.OrderID, // Use OrderID as customer ID since CustomerID doesn't exist
		items,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Reserve for POS transaction")
		return
	}

	s.logger.Info("Reserved inventory for POS transaction",
		zap.String("order_id", req.OrderID),
		zap.String("location_id", req.LocationID),
		zap.String("staff_id", req.StaffID),
	)

	respondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Inventory reserved successfully",
	})
}

// completePickup completes an in-store pickup
func (s *Server) completePickup(c *gin.Context) {
	// Validate staff credentials and auth
	_, exists := c.Get("userID")
	if !exists {
		respondWithError(c, http.StatusUnauthorized, "Staff authentication required")
		return
	}

	// Get role to ensure staff/admin permissions
	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok || (roleStr != "staff" && roleStr != "admin") {
		respondWithError(c, http.StatusForbidden, "Staff or admin role required")
		return
	}

	var req POSPickupCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Complete pickup
	_, err := s.inventorySvc.CompletePickup(
		c.Request.Context(),
		req.OrderID,
		req.StaffID,
		req.Notes,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Complete pickup")
		return
	}

	s.logger.Info("Completed pickup",
		zap.String("order_id", req.OrderID),
		zap.String("location_id", req.LocationID),
		zap.String("staff_id", req.StaffID),
	)

	// Update the order status to reflect pickup completion
	err = s.orderSvc.UpdateOrderStatus(
		c.Request.Context(),
		req.OrderID,
		"completed",
		"Order picked up in store by customer",
	)
	if err != nil {
		s.logger.Warn("Failed to update order status after pickup completion",
			zap.String("order_id", req.OrderID),
			zap.Error(err),
		)
		// Continue since the pickup was still completed
	}

	respondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Pickup completed successfully",
	})
}

// deductForDirectPOSSale processes a direct POS sale without creating an order
func (s *Server) deductForDirectPOSSale(c *gin.Context) {
	// Validate staff credentials and auth
	staffID, exists := c.Get("userID")
	staffIDStr, ok := staffID.(string)
	if !ok || !exists {
		respondWithError(c, http.StatusUnauthorized, "Staff ID required")
		return
	}

	// Get role to ensure staff/admin permissions
	role, _ := c.Get("role")
	roleStr, ok := role.(string)
	if !ok || (roleStr != "staff" && roleStr != "admin") {
		respondWithError(c, http.StatusForbidden, "Staff or admin role required")
		return
	}

	var req POSDirectSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Convert request items to domain format
	items := make([]map[string]interface{}, len(req.Items))
	for i, item := range req.Items {
		saleItem := map[string]interface{}{
			"product_id": item.ProductID,
			"sku":        item.SKU,
			"quantity":   item.Quantity,
		}
		items[i] = saleItem
	}

	// Use staff ID from authenticated user if not provided in request
	if req.StaffID == "" {
		req.StaffID = staffIDStr
	}

	// Directly deduct inventory for a POS sale
	_, err := s.inventorySvc.DeductForDirectPOSTransaction(
		c.Request.Context(),
		req.LocationID,
		req.StaffID,
		items,
		req.Notes,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Process direct POS sale")
		return
	}

	s.logger.Info("Processed direct POS sale",
		zap.String("location_id", req.LocationID),
		zap.String("staff_id", req.StaffID),
		zap.Int("item_count", len(req.Items)),
	)

	respondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Direct POS sale processed successfully",
	})
}
