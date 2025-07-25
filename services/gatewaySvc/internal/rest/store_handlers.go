package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getStores returns a list of stores with pagination
func (s *Server) getStores(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid limit parameter")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid offset parameter")
		return
	}

	stores, err := s.storeSvc.ListStores(c.Request.Context(), limit, offset)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get stores")
		return
	}

	respondWithSuccess(c, http.StatusOK, stores)
}

// getStore returns a specific store by ID
func (s *Server) getStore(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respondWithError(c, http.StatusBadRequest, "Store ID is required")
		return
	}

	store, err := s.storeSvc.GetStore(c.Request.Context(), id)
	if err != nil {
		genericErrorHandler(c, err, s.logger, "Get store")
		return
	}

	respondWithSuccess(c, http.StatusOK, store)
}

// createStore handles creating a new store
func (s *Server) createStore(c *gin.Context) {
    var req struct {
        Name    string `json:"name" binding:"required"`
        Address string `json:"address" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        respondWithError(c, http.StatusBadRequest, "Invalid request payload")
        return
    }

    store, err := s.storeSvc.CreateStore(c.Request.Context(), req.Name, req.Address)
    if err != nil {
        genericErrorHandler(c, err, s.logger, "Create store")
        return
    }

    respondWithSuccess(c, http.StatusCreated, store)
}
