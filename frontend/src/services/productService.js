import axios from 'axios';

const API_BASE_URL = '/api/v1';

// Create axios instance with default configuration
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

const productService = {
  // List products with filtering and pagination
  listProducts: async (params = {}) => {
    const response = await api.get('/products', { params });
    // Normalize response to ensure consistent structure
    if (response.data && typeof response.data === 'object' && response.data.data) {
      // Backend returns { data: [...], total: number, success: boolean }
      return {
        data: Array.isArray(response.data.data) ? response.data.data : [],
        total: response.data.total || 0,
        success: response.data.success || true
      };
    }
    // Fallback: assume direct array or wrap in expected structure
    return {
      data: Array.isArray(response.data) ? response.data : [],
      total: Array.isArray(response.data) ? response.data.length : 0,
      success: true
    };
  },

  // Get product by ID
  getProductById: async (productId) => {
    const response = await api.get(`/products/${productId}`);
    return response.data;
  },

  // Create new product (Staff/Admin only)
  createProduct: async (productData) => {
    const response = await api.post('/products', productData);
    return response.data;
  },

  // Update product (Staff/Admin only)
  updateProduct: async (productId, productData) => {
    const response = await api.put(`/products/${productId}`, productData);
    return response.data;
  },

  // Delete product (Staff/Admin only)
  deleteProduct: async (productId) => {
    const response = await api.delete(`/products/${productId}`);
    return response.data;
  },

  // Search products
  searchProducts: async (query, params = {}) => {
    const response = await api.get('/products/search', {
      params: { q: query, ...params },
    });
    return response.data;
  },

  // Get products by category
  getProductsByCategory: async (categoryId, params = {}) => {
    const response = await api.get(`/products/category/${categoryId}`, { params });
    return response.data;
  },

  // Upload product image
  uploadProductImage: async (productId, file) => {
    const formData = new FormData();
    formData.append('image', file);
    
    const response = await api.post(`/products/${productId}/image`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Delete product image
  deleteProductImage: async (productId, imageId) => {
    const response = await api.delete(`/products/${productId}/images/${imageId}`);
    return response.data;
  },

  // Get product reviews
  getProductReviews: async (productId, params = {}) => {
    const response = await api.get(`/products/${productId}/reviews`, { params });
    return response.data;
  },

  // Create product review
  createProductReview: async (productId, reviewData) => {
    const response = await api.post(`/products/${productId}/reviews`, reviewData);
    return response.data;
  },

  // Update product review
  updateProductReview: async (productId, reviewId, reviewData) => {
    const response = await api.put(`/products/${productId}/reviews/${reviewId}`, reviewData);
    return response.data;
  },

  // Delete product review
  deleteProductReview: async (productId, reviewId) => {
    const response = await api.delete(`/products/${productId}/reviews/${reviewId}`);
    return response.data;
  },

  // Get product variants
  getProductVariants: async (productId) => {
    const response = await api.get(`/products/${productId}/variants`);
    return response.data;
  },

  // Create product variant
  createProductVariant: async (productId, variantData) => {
    const response = await api.post(`/products/${productId}/variants`, variantData);
    return response.data;
  },

  // Update product variant
  updateProductVariant: async (productId, variantId, variantData) => {
    const response = await api.put(`/products/${productId}/variants/${variantId}`, variantData);
    return response.data;
  },

  // Delete product variant
  deleteProductVariant: async (productId, variantId) => {
    const response = await api.delete(`/products/${productId}/variants/${variantId}`);
    return response.data;
  },

  // Category Management
  // List all categories
  listCategories: async () => {
    const response = await api.get('/products/categories');
    return response.data;
  },

  // Get category by ID
  getCategoryById: async (categoryId) => {
    const response = await api.get(`/products/categories/${categoryId}`);
    return response.data;
  },

  // Create new category (Staff/Admin only)
  createCategory: async (categoryData) => {
    const response = await api.post('/products/categories', categoryData);
    return response.data;
  },

  // Update category (Staff/Admin only)
  updateCategory: async (categoryId, categoryData) => {
    const response = await api.put(`/products/categories/${categoryId}`, categoryData);
    return response.data;
  },

  // Delete category (Staff/Admin only)
  deleteCategory: async (categoryId) => {
    const response = await api.delete(`/products/categories/${categoryId}`);
    return response.data;
  },

  // Get category hierarchy
  getCategoryHierarchy: async () => {
    const response = await api.get('/products/categories/hierarchy');
    return response.data;
  },

  // Bulk operations
  bulkUpdateProducts: async (productUpdates) => {
    const response = await api.put('/products/bulk', { products: productUpdates });
    return response.data;
  },

  bulkDeleteProducts: async (productIds) => {
    const response = await api.delete('/products/bulk', { data: { productIds } });
    return response.data;
  },

  // Get product analytics
  getProductAnalytics: async (productId, params = {}) => {
    const response = await api.get(`/products/${productId}/analytics`, { params });
    return response.data;
  },

  // Get trending products
  getTrendingProducts: async (params = {}) => {
    const response = await api.get('/products/trending', { params });
    return response.data;
  },

  // Get featured products
  getFeaturedProducts: async (params = {}) => {
    const response = await api.get('/products/featured', { params });
    return response.data;
  },
};

export default productService;
