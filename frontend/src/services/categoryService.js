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

/**
 * Category service for managing product categories
 * Handles API communication with the backend category endpoints
 */

const categoryService = {
  /**
   * Get all categories
   * @returns {Promise} List of categories
   */
  getCategories: async () => {
    const response = await api.get('/products/categories');
    const payload = response.data;
    // Normalize to an array of categories regardless of envelope shape
    // API returns { data: [...], success: true }
    return payload.data || [];
  },

  /**
   * Get category by ID
   * @param {string} id - Category ID
   * @returns {Promise} Category details
   */
  getCategoryById: async (id) => {
    const response = await api.get(`/products/categories/${id}`);
    return response.data;
  },

  /**
   * Create a new category
   * @param {Object} categoryData - Category data
   * @param {string} categoryData.name - Category name
   * @param {string} categoryData.description - Category description
   * @param {string} [categoryData.parent_id] - Parent category ID (optional)
   * @param {boolean} [categoryData.is_active] - Whether category is active (default: true)
   * @returns {Promise} Created category
   */
  createCategory: async (categoryData) => {
    const response = await api.post('/products/categories', categoryData);
    return response.data;
  },

  /**
   * Update an existing category
   * @param {string} id - Category ID
   * @param {Object} categoryData - Updated category data
   * @returns {Promise} Updated category
   */
  updateCategory: async (id, categoryData) => {
    const response = await api.put(`/products/categories/${id}`, categoryData);
    return response.data;
  },

  /**
   * Delete a category
   * @param {string} id - Category ID
   * @returns {Promise} Deletion confirmation
   */
  deleteCategory: async (id) => {
    const response = await api.delete(`/products/categories/${id}`);
    return response.data;
  },

  /**
   * Get categories in tree/hierarchical format
   * @returns {Promise} Hierarchical category structure
   */
  getCategoryTree: async () => {
    const categories = await categoryService.getCategories();
    return buildCategoryTree(categories);
  }
};

/**
 * Build hierarchical category tree from flat category list
 * @param {Array} categories - Flat list of categories
 * @returns {Array} Hierarchical category tree
 */
const buildCategoryTree = (categories) => {
  const categoryMap = {};
  const rootCategories = [];

  // Create a map of categories by ID
  categories.forEach(category => {
    categoryMap[category.id] = { ...category, children: [] };
  });

  // Build the tree structure
  categories.forEach(category => {
    if (category.parent_id && categoryMap[category.parent_id]) {
      // Add as child to parent
      categoryMap[category.parent_id].children.push(categoryMap[category.id]);
    } else {
      // Add as root category
      rootCategories.push(categoryMap[category.id]);
    }
  });

  return rootCategories;
};

/**
 * Get breadcrumb path for a category
 * @param {string} categoryId - Category ID
 * @param {Array} categories - All categories
 * @returns {Array} Breadcrumb path from root to category
 */
export const getCategoryBreadcrumb = (categoryId, categories) => {
  const categoryMap = {};
  categories.forEach(cat => categoryMap[cat.id] = cat);

  const breadcrumb = [];
  let currentCategory = categoryMap[categoryId];

  while (currentCategory) {
    breadcrumb.unshift(currentCategory);
    currentCategory = currentCategory.parent_id ? categoryMap[currentCategory.parent_id] : null;
  }

  return breadcrumb;
};

/**
 * Get all descendants of a category
 * @param {string} categoryId - Parent category ID
 * @param {Array} categories - All categories
 * @returns {Array} All descendant categories
 */
export const getCategoryDescendants = (categoryId, categories) => {
  const descendants = [];
  const directChildren = categories.filter(cat => cat.parent_id === categoryId);

  directChildren.forEach(child => {
    descendants.push(child);
    descendants.push(...getCategoryDescendants(child.id, categories));
  });

  return descendants;
};

export default categoryService;
