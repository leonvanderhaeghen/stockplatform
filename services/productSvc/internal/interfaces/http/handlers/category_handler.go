package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/application"
	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/domain"
)

type CategoryHandler struct {
	categoryService *application.CategoryService
	logger         *zap.Logger
}

func NewCategoryHandler(categoryService *application.CategoryService, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:         logger.Named("http_category_handler"),
	}
}

func (h *CategoryHandler) RegisterRoutes(router *mux.Router) {
	h.logger.Info("Registering category routes...")

	routes := []struct {
		path    string
		handler http.HandlerFunc
		methods []string
		queries []string
	}{
		{"/products/categories", h.CreateCategory, []string{"POST"}, nil},
		{"/products/categories/{id}", h.GetCategory, []string{"GET"}, nil},
		{"/products/categories/{id}", h.UpdateCategory, []string{"PUT"}, nil},
		{"/products/categories/{id}", h.DeleteCategory, []string{"DELETE"}, nil},
		{"/products/categories", h.ListCategories, []string{"GET"}, []string{"parent_id", "{parent_id}"}},
		{"/products/categories", h.ListCategories, []string{"GET"}, nil},
	}

	for _, route := range routes {
		r := router.HandleFunc(route.path, route.handler).Methods(route.methods...)
		if route.queries != nil && len(route.queries) == 2 {
			r.Queries(route.queries[0], route.queries[1])
		}
		h.logger.Info("Registered route",
			zap.String("path", route.path),
			zap.Strings("methods", route.methods),
			zap.Strings("queries", route.queries),
		)
	}

	h.logger.Info("Category routes registered successfully")
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var category domain.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	createdCategory, err := h.categoryService.CreateCategory(ctx, &category)
	if err != nil {
		h.logger.Error("Failed to create category", zap.Error(err))
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCategory)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	category, err := h.categoryService.GetCategory(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		h.logger.Error("Failed to get category", zap.Error(err))
		http.Error(w, "Failed to get category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert string ID to primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.logger.Error("Invalid category ID format", zap.Error(err))
		http.Error(w, "Invalid category ID format", http.StatusBadRequest)
		return
	}

	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	category.ID = objID
	if err := h.categoryService.UpdateCategory(ctx, &category); err != nil {
		h.logger.Error("Failed to update category", zap.Error(err))
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.categoryService.DeleteCategory(ctx, id); err != nil {
		h.logger.Error("Failed to delete category", zap.Error(err))
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parentID := r.URL.Query().Get("parent_id")

	// If parent_id is empty, it will return root categories
	categories, err := h.categoryService.ListCategories(ctx, parentID, 0)
	if err != nil {
		h.logger.Error("Failed to list categories", zap.Error(err))
		http.Error(w, "Failed to list categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
