import api from './api';

// Base path for POS endpoints (relative to the API base URL which already includes /v1)
const POS_BASE = '/pos';

const posService = {
  /**
   * Create a quick POS sale
   * @param {Object} saleData - POS sale data
   * @param {Array} saleData.items - Array of sale items
   * @param {string} saleData.items[].productId - Product ID
   * @param {number} saleData.items[].quantity - Quantity sold
   * @param {string} saleData.items[].price - Price per unit
   * @param {string} saleData.paymentMethod - Payment method (CASH, CARD, DIGITAL)
   * @param {string} saleData.totalAmount - Total sale amount
   * @param {string} saleData.storeId - Store ID where sale occurred
   * @param {string} [saleData.customerId] - Optional customer ID
   * @param {Object} [saleData.metadata] - Additional metadata
   * @returns {Promise<Object>} Created POS sale details
   */
  createQuickSale: async (saleData) => {
    try {
      const response = await api.post(`${POS_BASE}/sale`, saleData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error creating POS sale:', error);
      throw error;
    }
  },

  /**
   * Get POS sales with optional filters
   * @param {Object} params - Query parameters for filtering and pagination
   * @param {number} params.page - Page number for pagination
   * @param {number} params.limit - Number of items per page
   * @param {string} params.storeId - Filter by store ID
   * @param {string} params.startDate - Filter sales from this date
   * @param {string} params.endDate - Filter sales until this date
   * @param {string} params.paymentMethod - Filter by payment method
   * @returns {Promise<Object>} Paginated list of POS sales
   */
  getSales: async (params = {}) => {
    try {
      const response = await api.get(`${POS_BASE}/sales`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching POS sales:', error);
      throw error;
    }
  },

  /**
   * Get a specific POS sale by ID
   * @param {string} id - Sale ID
   * @returns {Promise<Object>} POS sale details
   */
  getSale: async (id) => {
    try {
      const response = await api.get(`${POS_BASE}/sales/${id}`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching POS sale:', error);
      throw error;
    }
  },

  /**
   * Update a POS sale
   * @param {string} id - Sale ID
   * @param {Object} saleData - Updated sale data
   * @returns {Promise<Object>} Updated POS sale details
   */
  updateSale: async (id, saleData) => {
    try {
      const response = await api.put(`${POS_BASE}/sales/${id}`, saleData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating POS sale:', error);
      throw error;
    }
  },

  /**
   * Delete a POS sale
   * @param {string} id - Sale ID
   * @returns {Promise<Object>} Deletion confirmation
   */
  deleteSale: async (id) => {
    try {
      const response = await api.delete(`${POS_BASE}/sales/${id}`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error deleting POS sale:', error);
      throw error;
    }
  },

  /**
   * Complete pickup for an order
   * @param {string} orderId - Order ID to complete pickup for
   * @param {Object} pickupData - Pickup completion data
   * @param {string} pickupData.staffId - Staff member completing pickup
   * @param {string} [pickupData.notes] - Optional pickup notes
   * @param {Array} [pickupData.items] - Optional specific items picked up
   * @returns {Promise<Object>} Pickup completion result
   */
  completePickup: async (orderId, pickupData) => {
    try {
      const response = await api.post(`${POS_BASE}/pickup/${orderId}/complete`, pickupData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error completing pickup:', error);
      throw error;
    }
  },

  /**
   * Deduct inventory for direct POS sale
   * @param {Object} deductionData - Inventory deduction data
   * @param {Array} deductionData.items - Array of items to deduct
   * @param {string} deductionData.items[].productId - Product ID
   * @param {number} deductionData.items[].quantity - Quantity to deduct
   * @param {string} deductionData.storeId - Store ID where deduction occurs
   * @param {string} deductionData.reason - Reason for deduction
   * @param {string} [deductionData.staffId] - Staff member performing deduction
   * @returns {Promise<Object>} Inventory deduction result
   */
  deductInventory: async (deductionData) => {
    try {
      const response = await api.post(`${POS_BASE}/inventory/deduct`, deductionData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error deducting inventory:', error);
      throw error;
    }
  },

  /**
   * Get POS sales statistics
   * @param {Object} params - Query parameters for statistics
   * @param {string} params.storeId - Store ID for statistics
   * @param {string} params.startDate - Start date for statistics
   * @param {string} params.endDate - End date for statistics
   * @param {string} [params.groupBy] - Group statistics by (day, week, month)
   * @returns {Promise<Object>} POS sales statistics
   */
  getSalesStatistics: async (params = {}) => {
    try {
      const response = await api.get(`${POS_BASE}/statistics`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching POS statistics:', error);
      throw error;
    }
  },

  /**
   * Get daily sales summary
   * @param {Object} params - Query parameters
   * @param {string} params.storeId - Store ID
   * @param {string} [params.date] - Specific date (defaults to today)
   * @returns {Promise<Object>} Daily sales summary
   */
  getDailySummary: async (params = {}) => {
    try {
      const response = await api.get(`${POS_BASE}/summary/daily`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching daily summary:', error);
      throw error;
    }
  },
};

export default posService;
