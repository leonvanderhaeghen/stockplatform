import api from './api';

// Base path for category endpoints
const CATEGORIES_BASE = '/products/categories';

const categoryService = {
  /**
   * Get all categories
   * @param {Object} params - Query parameters for filtering and pagination
   * @param {string} params.query - Search query string
   * @param {boolean} params.active - Filter by active status
   * @param {number} params.limit - Maximum number of results to return
   * @param {number} params.offset - Number of results to skip
   * @returns {Promise<Array>} List of categories
   */
  getCategories: async (params = {}) => {
    try {
      console.log('Fetching categories from:', CATEGORIES_BASE, 'with params:', params);
      const response = await api.get(CATEGORIES_BASE, { params });
      console.log('Raw categories response:', response);
      
      // Return the data property if it exists, otherwise return the full response
      // This handles both { data: [...] } and direct array responses
      const data = response.data !== undefined ? response.data : response;
      console.log('Processed categories data:', data);
      
      return data || [];
    } catch (error) {
      console.error('Error fetching categories:', error);
      console.error('Error details:', error.response?.data || error.message);
      return [];
    }
  },

  /**
   * Get a single category by ID
   * @param {string} id - Category ID
   * @returns {Promise<Object>} Category details
   */
  getCategory: async (id) => {
    const { data } = await api.get(`${CATEGORIES_BASE}/${id}`);
    return data;
  },

  /**
   * Create a new category
   * @param {Object} categoryData - Category data
   * @param {string} categoryData.name - Category name
   * @param {string} [categoryData.description] - Category description
   * @param {string} [categoryData.parentId] - Parent category ID
   * @param {boolean} [categoryData.isActive=true] - Whether the category is active
   * @returns {Promise<Object>} Created category details
   */
  createCategory: async (categoryData) => {
    const { data } = await api.post(CATEGORIES_BASE, categoryData);
    return data;
  },

  /**
   * Update an existing category
   * @param {string} id - Category ID
   * @param {Object} categoryData - Updated category data
   * @returns {Promise<Object>} Updated category details
   */
  updateCategory: async (id, categoryData) => {
    const { data } = await api.put(`${CATEGORIES_BASE}/${id}`, categoryData);
    return data;
  },

  /**
   * Delete a category
   * @param {string} id - Category ID
   * @returns {Promise<Object>} Deletion status
   */
  deleteCategory: async (id) => {
    const { data } = await api.delete(`${CATEGORIES_BASE}/${id}`);
    return data;
  },

  /**
   * Search categories by query
   * @param {string} query - Search query string
   * @returns {Promise<Array>} List of matching categories
   */
  searchCategories: async (query) => {
    const { data } = await api.get(`${CATEGORIES_BASE}/search`, { params: { q: query } });
    return data;
  },
};

export default categoryService;
