package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	"github.com/leonvanderhaeghen/stockplatform/services/gatewaySvc/internal/services"
)

// SupplierHandler handles HTTP requests for supplier operations
type SupplierHandler struct {
	svc    services.SupplierService
	logger *zap.Logger
}

// NewSupplierHandler creates a new supplier handler
func NewSupplierHandler(svc services.SupplierService, logger *zap.Logger) *SupplierHandler {
	return &SupplierHandler{
		svc:    svc,
		logger: logger.Named("supplier_handler"),
	}
}

// CreateSupplierRequest represents the request body for creating a supplier
type CreateSupplierRequest struct {
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email" binding:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	TaxID         string `json:"tax_id"`
	Website       string `json:"website"`
	Currency      string `json:"currency"`
	LeadTimeDays  int32  `json:"lead_time_days"`
	PaymentTerms  string `json:"payment_terms"`
}

// CreateSupplier creates a new supplier
// @Summary Create a new supplier
// @Description Create a new supplier with the provided details
// @Tags suppliers
// @Accept json
// @Produce json
// @Param request body CreateSupplierRequest true "Supplier details"
// @Success 201 {object} supplierv1.Supplier
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers [post]
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	supplier, err := h.svc.CreateSupplier(c.Request.Context(), &supplierv1.CreateSupplierRequest{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		City:          req.City,
		State:         req.State,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
		TaxId:         req.TaxID,
		Website:       req.Website,
		Currency:      req.Currency,
		LeadTimeDays:  req.LeadTimeDays,
		PaymentTerms:  req.PaymentTerms,
	})

	if err != nil {
		h.logger.Error("Failed to create supplier", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier"})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// GetSupplier retrieves a supplier by ID
// @Summary Get a supplier by ID
// @Description Get a supplier by its ID
// @Tags suppliers
// @Produce json
// @Param id path string true "Supplier ID"
// @Success 200 {object} supplierv1.Supplier
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id} [get]
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supplier ID is required"})
		return
	}

	supplier, err := h.svc.GetSupplier(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
			return
		}
		h.logger.Error("Failed to get supplier", zap.Error(err), zap.String("supplier_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get supplier"})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// UpdateSupplierRequest represents the request body for updating a supplier
type UpdateSupplierRequest struct {
	Name          string `json:"name"`
	ContactPerson string `json:"contact_person"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	TaxID         string `json:"tax_id"`
	Website       string `json:"website"`
	Currency      string `json:"currency"`
	LeadTimeDays  int32  `json:"lead_time_days"`
	PaymentTerms  string `json:"payment_terms"`
}

// UpdateSupplier updates an existing supplier
// @Summary Update a supplier
// @Description Update an existing supplier with the provided details
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Param request body UpdateSupplierRequest true "Supplier details"
// @Success 200 {object} supplierv1.Supplier
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id} [put]
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supplier ID is required"})
		return
	}

	var req UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	supplier, err := h.svc.UpdateSupplier(c.Request.Context(), &supplierv1.UpdateSupplierRequest{
		Id:            id,
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		City:          req.City,
		State:         req.State,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
		TaxId:         req.TaxID,
		Website:       req.Website,
		Currency:      req.Currency,
		LeadTimeDays:  req.LeadTimeDays,
		PaymentTerms:  req.PaymentTerms,
	})

	if err != nil {
		if status.Code(err) == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
			return
		}
		h.logger.Error("Failed to update supplier", zap.Error(err), zap.String("supplier_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier"})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// DeleteSupplier deletes a supplier by ID
// @Summary Delete a supplier
// @Description Delete a supplier by its ID
// @Tags suppliers
// @Param id path string true "Supplier ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id} [delete]
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supplier ID is required"})
		return
	}

	err := h.svc.DeleteSupplier(c.Request.Context(), id)
	if err != nil {
		if status.Code(err) == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
			return
		}
		h.logger.Error("Failed to delete supplier", zap.Error(err), zap.String("supplier_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListSuppliersResponse represents the response for listing suppliers
type ListSuppliersResponse struct {
	Suppliers []*supplierv1.Supplier `json:"suppliers"`
	Total     int32                  `json:"total"`
	Page      int32                  `json:"page"`
	PageSize  int32                  `json:"page_size"`
}

// ListSuppliers lists suppliers with pagination and search
// @Summary List suppliers
// @Description List suppliers with pagination and search
// @Tags suppliers
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Items per page (default: 10, max: 100)"
// @Param search query string false "Search query"
// @Success 200 {object} ListSuppliersResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers [get]
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	// Parse query parameters with defaults
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)
	search := c.Query("search")

	// Validate page and page size
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	suppliers, total, err := h.svc.ListSuppliers(c.Request.Context(), int32(page), int32(pageSize), search)
	if err != nil {
		h.logger.Error("Failed to list suppliers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list suppliers"})
		return
	}

	c.JSON(http.StatusOK, ListSuppliersResponse{
		Suppliers: suppliers,
		Total:     total,
		Page:      int32(page),
		PageSize:  int32(pageSize),
	})
}

// RegisterRoutes registers the supplier handler routes
func (h *SupplierHandler) RegisterRoutes(router *gin.RouterGroup) {
	suppliersGroup := router.Group("/suppliers")
	{
		suppliersGroup.POST("", h.CreateSupplier)
		suppliersGroup.GET(":id", h.GetSupplier)
		suppliersGroup.PUT(":id", h.UpdateSupplier)
		suppliersGroup.DELETE(":id", h.DeleteSupplier)
		suppliersGroup.GET("", h.ListSuppliers)
		
		// Adapter routes
		suppliersGroup.GET("/adapters", h.ListAdapters)
		suppliersGroup.GET("/adapters/:name/capabilities", h.GetAdapterCapabilities)
		suppliersGroup.POST("/adapters/:name/test-connection", h.TestAdapterConnection)
		
		// Sync routes
		suppliersGroup.POST(":id/sync/products", h.SyncProducts)
		suppliersGroup.POST(":id/sync/inventory", h.SyncInventory)
	}
}

// AdapterCapabilitiesResponse represents the response for adapter capabilities
type AdapterCapabilitiesResponse struct {
	Capabilities map[string]bool `json:"capabilities"`
}

// TestConnectionRequest represents the request to test a connection to a supplier
type TestConnectionRequest struct {
	Config map[string]string `json:"config" binding:"required"`
}

// SyncOptionsRequest represents synchronization options
type SyncOptionsRequest struct {
	FullSync       bool   `json:"full_sync"`
	BatchSize      int32  `json:"batch_size,omitempty"`
	Since          string `json:"since,omitempty"`
	IncludeInactive bool  `json:"include_inactive,omitempty"`
}

// SyncResponse represents the response for sync operations
type SyncResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message,omitempty"`
}

// ListAdapters lists all available supplier adapters
// @Summary List supplier adapters
// @Description List all available supplier adapters
// @Tags suppliers
// @Produce json
// @Success 200 {array} supplierv1.SupplierAdapter
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/adapters [get]
func (h *SupplierHandler) ListAdapters(c *gin.Context) {
	adapters, err := h.svc.ListAdapters(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list adapters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve adapters"})
		return
	}

	c.JSON(http.StatusOK, adapters)
}

// GetAdapterCapabilities retrieves the capabilities of a specific adapter
// @Summary Get adapter capabilities
// @Description Get the capabilities of a specific supplier adapter
// @Tags suppliers
// @Produce json
// @Param name path string true "Adapter Name"
// @Success 200 {object} AdapterCapabilitiesResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/adapters/{name}/capabilities [get]
func (h *SupplierHandler) GetAdapterCapabilities(c *gin.Context) {
	adapterName := c.Param("name")
	if adapterName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Adapter name is required"})
		return
	}

	capabilities, err := h.svc.GetAdapterCapabilities(c.Request.Context(), adapterName)
	if err != nil {
		if status.Code(err) == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Adapter not found"})
			return
		}
		h.logger.Error("Failed to get adapter capabilities", zap.Error(err), zap.String("adapter_name", adapterName))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve adapter capabilities"})
		return
	}

	c.JSON(http.StatusOK, capabilities)
}

// TestAdapterConnection tests the connection to a supplier system using a specified adapter
// @Summary Test adapter connection
// @Description Test the connection to a supplier system using a specified adapter
// @Tags suppliers
// @Accept json
// @Produce json
// @Param name path string true "Adapter Name"
// @Param request body TestConnectionRequest true "Connection configuration"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/adapters/{name}/test-connection [post]
func (h *SupplierHandler) TestAdapterConnection(c *gin.Context) {
	adapterName := c.Param("name")
	if adapterName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Adapter name is required"})
		return
	}

	var req TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.svc.TestAdapterConnection(c.Request.Context(), adapterName, req.Config)
	if err != nil {
		if status.Code(err) == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Adapter not found"})
			return
		}
		h.logger.Error("Failed to test connection", zap.Error(err), zap.String("adapter_name", adapterName))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection successful"})
}

// SyncProducts initiates a product synchronization job for a supplier
// @Summary Sync supplier products
// @Description Synchronize products from a supplier using their configured adapter
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Param options body SyncOptionsRequest false "Synchronization options"
// @Success 202 {object} SyncResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id}/sync/products [post]
func (h *SupplierHandler) SyncProducts(c *gin.Context) {
	supplierID := c.Param("id")
	if supplierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supplier ID is required"})
		return
	}

	var req SyncOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no body is provided, use default options
		if err.Error() != "EOF" {
			h.logger.Error("Failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	}

	// Convert relevant fields to protobuf SyncOptions
	var syncOpts supplierv1.SyncOptions
	syncOpts.FullSync = req.FullSync
	syncOpts.BatchSize = req.BatchSize
	syncOpts.IncludeInactive = req.IncludeInactive
	
	jobID, err := h.svc.SyncProducts(c.Request.Context(), supplierID, &syncOpts)
	if err != nil {
		if status.Code(err) == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
			return
		}
		h.logger.Error("Failed to sync products", zap.Error(err), zap.String("supplier_id", supplierID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate product synchronization"})
		return
	}

	c.JSON(http.StatusAccepted, SyncResponse{
		JobID:   jobID,
		Message: "Product synchronization job started",
	})
}

// SyncInventory initiates an inventory synchronization job for a supplier
// @Summary Sync supplier inventory
// @Description Synchronize inventory from a supplier using their configured adapter
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID"
// @Param options body SyncOptionsRequest false "Synchronization options"
// @Success 202 {object} SyncResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/suppliers/{id}/sync/inventory [post]
func (h *SupplierHandler) SyncInventory(c *gin.Context) {
	supplierID := c.Param("id")
	if supplierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Supplier ID is required"})
		return
	}

	var req SyncOptionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no body is provided, use default options
		if err.Error() != "EOF" {
			h.logger.Error("Failed to bind request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	}

	// Convert relevant fields to protobuf SyncOptions
	var syncOpts supplierv1.SyncOptions
	syncOpts.FullSync = req.FullSync
	syncOpts.BatchSize = req.BatchSize
	syncOpts.IncludeInactive = req.IncludeInactive
	
	jobID, err := h.svc.SyncInventory(c.Request.Context(), supplierID, &syncOpts)
	if err != nil {
		if status.Code(err) == 404 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
			return
		}
		h.logger.Error("Failed to sync inventory", zap.Error(err), zap.String("supplier_id", supplierID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate inventory synchronization"})
		return
	}

	c.JSON(http.StatusAccepted, SyncResponse{
		JobID:   jobID,
		Message: "Inventory synchronization job started",
	})
}
