import api from './api';

// Base path for product endpoints (relative to the API base URL which already includes /v1)
const PRODUCTS_BASE = '/products';

const productService = {
  // Get all products with optional filters
  getProducts: async (params = {}) => {
    const { data } = await api.get(PRODUCTS_BASE, { params });
    return data;
  },

  // Get a single product by ID
  getProduct: async (id) => {
    const { data } = await api.get(`${PRODUCTS_BASE}/${id}`);
    return data;
  },

  // Create a new product
  createProduct: async (productData) => {
    const { data } = await api.post(PRODUCTS_BASE, productData);
    return data;
  },

  // Update an existing product
  updateProduct: async (id, productData) => {
    const { data } = await api.put(`${PRODUCTS_BASE}/${id}`, productData);
    return data;
  },

  // Delete a product
  deleteProduct: async (id) => {
    await api.delete(`${PRODUCTS_BASE}/${id}`);
    return id;
  },

    // Search products
  searchProducts: async (query) => {
    const { data } = await api.get(`${PRODUCTS_BASE}/search`, { params: { q: query } });
    return data;
  },
  
  // Get all product categories
  getCategories: async () => {
    const { data } = await api.get(`${PRODUCTS_BASE}/categories`);
    return data;
  },
};

export default productService;
