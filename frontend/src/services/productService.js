import api from './api';

// Base path for product endpoints (relative to the API base URL which already includes /v1)
const PRODUCTS_BASE = '/products';

const productService = {
  /**
   * Get all products with optional filters
   * @param {Object} params - Query parameters for filtering and pagination
   * @param {string} params.categoryId - Filter by category ID
   * @param {string} params.query - Search query string
   * @param {boolean} params.active - Filter by active status
   * @param {number} params.limit - Maximum number of results to return
   * @param {number} params.offset - Number of results to skip
   * @param {string} params.sortBy - Field to sort by
   * @param {boolean} params.ascending - Sort in ascending order
   * @returns {Promise<Array>} List of products
   */
  getProducts: async (params = {}) => {
    const { data } = await api.get(PRODUCTS_BASE, { params });
    return data;
  },

  /**
   * Get a single product by ID
   * @param {string} id - Product ID
   * @returns {Promise<Object>} Product details
   */
  getProduct: async (id) => {
    const { data } = await api.get(`${PRODUCTS_BASE}/${id}`);
    return data;
  },

  /**
   * Create a new product
   * @param {Object} productData - Product data
   * @param {string} productData.name - Product name
   * @param {string} productData.description - Product description
   * @param {string} productData.costPrice - Cost price as a string (e.g., "10.99")
   * @param {string} productData.sellingPrice - Selling price as a string (e.g., "19.99")
   * @param {string} productData.currency - Currency code (e.g., "USD")
   * @param {string} productData.sku - Stock Keeping Unit
   * @param {string} productData.barcode - Barcode
   * @param {Array<string>} productData.categoryIds - Array of category IDs
   * @param {string} productData.supplierId - Supplier ID
   * @param {boolean} productData.isActive - Whether the product is active
   * @param {boolean} productData.inStock - Whether the product is in stock
   * @param {number} productData.stockQty - Current stock quantity
   * @param {number} productData.lowStockAt - Threshold for low stock alert
   * @param {Array<string>} productData.imageUrls - Array of image URLs
   * @param {Array<string>} productData.videoUrls - Array of video URLs
   * @param {Object} productData.metadata - Additional metadata as key-value pairs
   * @returns {Promise<Object>} Created product details
   */
  createProduct: async (productData) => {
    const { data } = await api.post(PRODUCTS_BASE, productData);
    return data;
  },

  /**
   * Update an existing product
   * Note: This is not fully implemented in the backend API
   * @param {string} id - Product ID
   * @param {Object} productData - Updated product data
   * @returns {Promise<Object>} Updated product details
   */
  updateProduct: async (id, productData) => {
    const { data } = await api.put(`${PRODUCTS_BASE}/${id}`, productData);
    return data;
  },

  /**
   * Delete a product
   * Note: This is not implemented in the backend API
   * @param {string} id - Product ID
   * @returns {Promise<string>} Deleted product ID
   */
  deleteProduct: async (id) => {
    await api.delete(`${PRODUCTS_BASE}/${id}`);
    return id;
  },

  /**
   * Search products by query
   * @param {string} query - Search query string
   * @returns {Promise<Array>} List of matching products
   */
  searchProducts: async (query) => {
    const { data } = await api.get(`${PRODUCTS_BASE}/search`, { params: { q: query } });
    return data;
  },
  
  /**
   * Get all product categories
   * @returns {Promise<Array>} List of product categories
   */
  getCategories: async () => {
    const { data } = await api.get(`${PRODUCTS_BASE}/categories`);
    return data;
  },
};

export default productService;
