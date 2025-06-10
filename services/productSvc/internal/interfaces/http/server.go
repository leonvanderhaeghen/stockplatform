package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/application"
	categoryHandlers "github.com/leonvanderhaeghen/stockplatform/services/productSvc/internal/interfaces/http/handlers"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(port string, productSvc *application.ProductService, categorySvc *application.CategoryService, logger *zap.Logger) *Server {
	// Create a new router
	router := mux.NewRouter()

	// Add CORS middleware
	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Create a subrouter for API v1
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Register routes with the API subrouter
	handlers := NewHandlers(productSvc, categorySvc, logger)
	handlers.RegisterRoutes(apiRouter)

	// Add a test route for debugging
	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Log all registered routes
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		logger.Info("Registered route",
			zap.String("path", t),
			zap.Strings("methods", m),
		)
		return nil
	})

	// Add CORS and logging middleware
	handler := cors(router)
	handler = loggingMiddleware(logger)(handler)

	return &Server{
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

func loggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response wrapper to capture the status code
			wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Process the request
			next.ServeHTTP(wrapper, r)


			// Log the request details
		duration := time.Since(start)
		logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.Int("status", wrapper.status),
			zap.Duration("duration", duration),
			zap.String("remote_addr", r.RemoteAddr),
		)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

type Handlers struct {
	productService  *application.ProductService
	categoryService *application.CategoryService
	logger         *zap.Logger
}

func NewHandlers(productService *application.ProductService, categoryService *application.CategoryService, logger *zap.Logger) *Handlers {
	return &Handlers{
		productService:  productService,
		categoryService: categoryService,
		logger:         logger.Named("http_handlers"),
	}
}

func (h *Handlers) RegisterRoutes(router *mux.Router) {
	h.logger.Info("Registering routes...")

	// Register category routes
	categoryHandler := categoryHandlers.NewCategoryHandler(h.categoryService, h.logger)
	categoryHandler.RegisterRoutes(router)

	h.logger.Info("Routes registered successfully")
}
