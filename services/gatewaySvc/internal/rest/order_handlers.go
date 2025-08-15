package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// OrderItemRequest represents an order item in the order request
type OrderItemRequest struct {
	ProductID string  `json:"productId" binding:"required"`
	SKU       string  `json:"sku" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required,gt=0"`
	Price     float64 `json:"price" binding:"required,gt=0"`
}

// OrderRequest represents the order request body
type OrderRequest struct {
	Items        []OrderItemRequest `json:"items" binding:"required,min=1"`
	AddressID    string             `json:"addressId"`                    // Optional for POS orders
	PaymentType  string             `json:"paymentType" binding:"required"`
	PaymentData  map[string]string  `json:"paymentData"`
	ShippingType string             `json:"shippingType"`                 // Optional for POS orders
	Notes        string             `json:"notes"`
	Source       string             `json:"source"`                       // "POS", "ONLINE", "QUICK_POS" - determines order type
	StoreID      string             `json:"storeId"`                      // Required for POS orders
	CustomerInfo map[string]string  `json:"customerInfo"`                // For walk-in customers (POS)
}

// OrderStatusRequest represents the order status update request
type OrderStatusRequest struct {
	Status      string `json:"status" binding:"required"`
	Description string `json:"description"`
}

// OrderPaymentRequest represents adding a payment to an order
type OrderPaymentRequest struct {
	Amount      float64           `json:"amount" binding:"required,gt=0"`
	Type        string            `json:"type" binding:"required"`
	Reference   string            `json:"reference"`
	Status      string            `json:"status" binding:"required"`
	Date        time.Time         `json:"date"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// OrderTrackingRequest represents adding tracking info to an order
type OrderTrackingRequest struct {
	Carrier     string    `json:"carrier" binding:"required"`
	TrackingNum string    `json:"trackingNum" binding:"required"`
	ShipDate    time.Time `json:"shipDate"`
	EstDelivery time.Time `json:"estDelivery"`
	Notes       string    `json:"notes"`
}

// getUserOrders returns orders for the current user
func (s *Server) getUserOrders(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	status := c.Query("status")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
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

	orders, err := s.orderSvc.GetUserOrders(
		c.Request.Context(),
		userIDStr,
		status,
		startDate,
		endDate,
		limit,
		offset,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get user orders")
		return
	}

	respondWithSuccess(c, http.StatusOK, orders)
}

// getUserOrder returns a specific order for the current user
func (s *Server) getUserOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		respondWithError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := s.orderSvc.GetUserOrder(c.Request.Context(), orderID, userIDStr)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get user order")
		return
	}

	respondWithSuccess(c, http.StatusOK, order)
}

// createOrder creates a new order for the current user
func (s *Server) createOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		respondWithError(c, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Validate request based on source type
	if req.Source == "POS" || req.Source == "QUICK_POS" {
		// POS orders require storeId
		if req.StoreID == "" {
			respondWithError(c, http.StatusBadRequest, "Store ID is required for POS orders")
			return
		}
	} else {
		// Online orders require addressId and shippingType
		if req.AddressID == "" {
			respondWithError(c, http.StatusBadRequest, "Address ID is required for online orders")
			return
		}
		if req.ShippingType == "" {
			respondWithError(c, http.StatusBadRequest, "Shipping type is required for online orders")
			return
		}
	}

	// Convert request items to service items
	items := make([]map[string]interface{}, len(req.Items))
	for i, item := range req.Items {
		items[i] = map[string]interface{}{
			"productId": item.ProductID,
			"sku":       item.SKU,
			"quantity":  item.Quantity,
			"price":     item.Price,
		}
	}

	order, err := s.orderSvc.CreateOrder(
		c.Request.Context(),
		userIDStr,
		items,
		req.AddressID,
		req.PaymentType,
		req.PaymentData,
		req.ShippingType,
		req.Notes,
		req.Source,     // POS-specific: "POS", "ONLINE", "QUICK_POS"
		req.StoreID,    // POS-specific: required for POS orders
		req.CustomerInfo, // POS-specific: walk-in customer information
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Create order")
		return
	}

	respondWithSuccess(c, http.StatusCreated, order)
}

// listOrders returns all orders (admin/staff only)
func (s *Server) listOrders(c *gin.Context) {
	status := c.Query("status")
	userID := c.Query("userId")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
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

	orders, err := s.orderSvc.ListOrders(
		c.Request.Context(),
		status,
		userID,
		startDate,
		endDate,
		limit,
		offset,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "List orders")
		return
	}

	respondWithSuccess(c, http.StatusOK, orders)
}

// getOrder returns a specific order (admin/staff only)
func (s *Server) getOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		respondWithError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	order, err := s.orderSvc.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get order")
		return
	}

	respondWithSuccess(c, http.StatusOK, order)
}

// updateOrderStatus updates the status of an order (admin/staff only)
func (s *Server) updateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		respondWithError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	var req OrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := s.orderSvc.UpdateOrderStatus(
		c.Request.Context(),
		orderID,
		req.Status,
		req.Description,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Update order status")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// addOrderPayment adds a payment to an order (admin/staff only)
func (s *Server) addOrderPayment(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		respondWithError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	var req OrderPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := s.orderSvc.AddOrderPayment(
		c.Request.Context(),
		orderID,
		req.Amount,
		req.Type,
		req.Reference,
		req.Status,
		req.Date,
		req.Description,
		req.Metadata,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Add order payment")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Order payment added successfully"})
}

// addOrderTracking adds tracking information to an order (admin/staff only)
func (s *Server) addOrderTracking(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		respondWithError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	var req OrderTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := s.orderSvc.AddOrderTracking(
		c.Request.Context(),
		orderID,
		req.Carrier,
		req.TrackingNum,
		req.ShipDate,
		req.EstDelivery,
		req.Notes,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Add order tracking")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Order tracking added successfully"})
}

// cancelOrder cancels an order (admin/staff only)
func (s *Server) cancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		respondWithError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	reason := c.Query("reason")

	err := s.orderSvc.CancelOrder(
		c.Request.Context(),
		orderID,
		reason,
	)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Cancel order")
		return
	}

	respondWithSuccess(c, http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
