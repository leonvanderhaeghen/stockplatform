import api from './api';

// Base paths for admin endpoints (relative to the API base URL which already includes /v1)
const ADMIN_BASE = '/admin';

const adminService = {
  // Admin User Management
  /**
   * Get all users (Admin only)
   * @param {Object} params - Query parameters for filtering and pagination
   * @param {number} params.page - Page number for pagination
   * @param {number} params.limit - Number of items per page
   * @param {string} params.role - Filter by user role
   * @param {string} params.email - Filter by email
   * @param {boolean} params.active - Filter by active status
   * @returns {Promise<Object>} Paginated list of users
   */
  getAllUsers: async (params = {}) => {
    try {
      const response = await api.get(`${ADMIN_BASE}/users`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching all users:', error);
      throw error;
    }
  },

  /**
   * Get a specific user by ID (Admin only)
   * @param {string} userId - User ID
   * @returns {Promise<Object>} User details
   */
  getUserById: async (userId) => {
    try {
      const response = await api.get(`${ADMIN_BASE}/users/${userId}`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching user by ID:', error);
      throw error;
    }
  },

  /**
   * Create a new user (Admin only)
   * @param {Object} userData - User data
   * @param {string} userData.email - User email
   * @param {string} userData.password - User password
   * @param {string} userData.firstName - User first name
   * @param {string} userData.lastName - User last name
   * @param {string} userData.role - User role (ADMIN, CUSTOMER, MANAGER, USER)
   * @param {boolean} userData.isActive - Whether the user is active
   * @returns {Promise<Object>} Created user details
   */
  createUser: async (userData) => {
    try {
      const response = await api.post(`${ADMIN_BASE}/users`, userData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error creating user:', error);
      throw error;
    }
  },

  /**
   * Update a user (Admin only)
   * @param {string} userId - User ID
   * @param {Object} userData - Updated user data
   * @returns {Promise<Object>} Updated user details
   */
  updateUser: async (userId, userData) => {
    try {
      const response = await api.put(`${ADMIN_BASE}/users/${userId}`, userData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating user:', error);
      throw error;
    }
  },

  /**
   * Delete a user (Admin only)
   * @param {string} userId - User ID
   * @returns {Promise<Object>} Deletion confirmation
   */
  deleteUser: async (userId) => {
    try {
      const response = await api.delete(`${ADMIN_BASE}/users/${userId}`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error deleting user:', error);
      throw error;
    }
  },

  /**
   * Update user role (Admin only)
   * @param {string} userId - User ID
   * @param {string} role - New role (ADMIN, CUSTOMER, MANAGER, USER)
   * @returns {Promise<Object>} Updated user details
   */
  updateUserRole: async (userId, role) => {
    try {
      const response = await api.put(`${ADMIN_BASE}/users/${userId}/role`, { role });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating user role:', error);
      throw error;
    }
  },

  /**
   * Activate/deactivate a user (Admin only)
   * @param {string} userId - User ID
   * @param {boolean} isActive - Active status
   * @returns {Promise<Object>} Updated user details
   */
  updateUserStatus: async (userId, isActive) => {
    try {
      const response = await api.put(`${ADMIN_BASE}/users/${userId}/status`, { isActive });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating user status:', error);
      throw error;
    }
  },

  // Admin Order Management
  /**
   * Get all orders (Admin only)
   * @param {Object} params - Query parameters for filtering and pagination
   * @param {number} params.page - Page number for pagination
   * @param {number} params.limit - Number of items per page
   * @param {string} params.status - Filter by order status
   * @param {string} params.userId - Filter by user ID
   * @param {string} params.startDate - Filter orders from this date
   * @param {string} params.endDate - Filter orders until this date
   * @returns {Promise<Object>} Paginated list of orders
   */
  getAllOrders: async (params = {}) => {
    try {
      const response = await api.get(`${ADMIN_BASE}/orders`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching all orders:', error);
      throw error;
    }
  },

  /**
   * Get a specific order by ID (Admin only)
   * @param {string} orderId - Order ID
   * @returns {Promise<Object>} Order details
   */
  getOrderById: async (orderId) => {
    try {
      const response = await api.get(`${ADMIN_BASE}/orders/${orderId}`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching order by ID:', error);
      throw error;
    }
  },

  /**
   * Update an order (Admin only)
   * @param {string} orderId - Order ID
   * @param {Object} orderData - Updated order data
   * @returns {Promise<Object>} Updated order details
   */
  updateOrder: async (orderId, orderData) => {
    try {
      const response = await api.put(`${ADMIN_BASE}/orders/${orderId}`, orderData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating order:', error);
      throw error;
    }
  },

  /**
   * Update order status (Admin only)
   * @param {string} orderId - Order ID
   * @param {Object} statusData - Status update data
   * @param {string} statusData.status - New status
   * @param {string} [statusData.trackingNumber] - Tracking number (required for SHIPPED status)
   * @param {string} [statusData.notes] - Optional notes
   * @returns {Promise<Object>} Updated order details
   */
  updateOrderStatus: async (orderId, statusData) => {
    try {
      const response = await api.put(`${ADMIN_BASE}/orders/${orderId}/status`, statusData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating order status:', error);
      throw error;
    }
  },

  /**
   * Add tracking code to order (Admin only)
   * @param {string} orderId - Order ID
   * @param {Object} trackingData - Tracking data
   * @param {string} trackingData.trackingCode - Tracking code
   * @param {string} [trackingData.carrier] - Carrier name
   * @param {string} [trackingData.notes] - Optional notes
   * @returns {Promise<Object>} Updated order details
   */
  addOrderTracking: async (orderId, trackingData) => {
    try {
      const response = await api.put(`${ADMIN_BASE}/orders/${orderId}/tracking`, trackingData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error adding order tracking:', error);
      throw error;
    }
  },

  // Admin Analytics and Reports
  /**
   * Get user statistics (Admin only)
   * @param {Object} params - Query parameters
   * @param {string} [params.startDate] - Start date for statistics
   * @param {string} [params.endDate] - End date for statistics
   * @returns {Promise<Object>} User statistics
   */
  getUserStatistics: async (params = {}) => {
    try {
      const response = await api.get(`${ADMIN_BASE}/statistics/users`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching user statistics:', error);
      throw error;
    }
  },

  /**
   * Get order statistics (Admin only)
   * @param {Object} params - Query parameters
   * @param {string} [params.startDate] - Start date for statistics
   * @param {string} [params.endDate] - End date for statistics
   * @param {string} [params.groupBy] - Group by (day, week, month)
   * @returns {Promise<Object>} Order statistics
   */
  getOrderStatistics: async (params = {}) => {
    try {
      const response = await api.get(`${ADMIN_BASE}/statistics/orders`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching order statistics:', error);
      throw error;
    }
  },

  /**
   * Get sales statistics (Admin only)
   * @param {Object} params - Query parameters
   * @param {string} [params.startDate] - Start date for statistics
   * @param {string} [params.endDate] - End date for statistics
   * @param {string} [params.groupBy] - Group by (day, week, month)
   * @returns {Promise<Object>} Sales statistics
   */
  getSalesStatistics: async (params = {}) => {
    try {
      const response = await api.get(`${ADMIN_BASE}/statistics/sales`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching sales statistics:', error);
      throw error;
    }
  },

  /**
   * Get system health status (Admin only)
   * @returns {Promise<Object>} System health status
   */
  getSystemHealth: async () => {
    try {
      const response = await api.get(`${ADMIN_BASE}/system/health`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching system health:', error);
      throw error;
    }
  },
};

export default adminService;
