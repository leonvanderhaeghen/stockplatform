import api from './api';

// Base path for inventory endpoints (relative to the API base URL which already includes /v1)
const INVENTORY_BASE = '/inventory';

const inventoryService = {
  // Get all inventory items with optional filters
  getInventoryItems: async (params = {}) => {
    const { data } = await api.get(INVENTORY_BASE, { params });
    return data;
  },

  // Get a single inventory item by ID
  getInventoryItem: async (id) => {
    const { data } = await api.get(`${INVENTORY_BASE}/${id}`);
    return data;
  },

  // Create a new inventory item
  createInventoryItem: async (itemData) => {
    const { data } = await api.post(INVENTORY_BASE, itemData);
    return data;
  },

  // Update an existing inventory item
  updateInventoryItem: async (id, itemData) => {
    const { data } = await api.put(`${INVENTORY_BASE}/${id}`, itemData);
    return data;
  },

  // Update inventory quantity (add/subtract)
  updateInventoryQuantity: async (id, change, reason = '') => {
    const { data } = await api.patch(`${INVENTORY_BASE}/${id}/quantity`, {
      change,
      reason,
    });
    return data;
  },

  // Delete an inventory item
  deleteInventoryItem: async (id) => {
    await api.delete(`${INVENTORY_BASE}/${id}`);
    return id;
  },

  // Get inventory history for an item
  getInventoryHistory: async (id) => {
    const { data } = await api.get(`${INVENTORY_BASE}/${id}/history`);
    return data;
  },

  // Get low stock items
  getLowStockItems: async (threshold = 10, location = '') => {
    const params = { threshold };
    if (location) {
      params.location = location;
    }
    const { data } = await api.get(`${INVENTORY_BASE}/low-stock`, { params });
    return data;
  },

  /**
   * Reserve inventory for an order
   * @param {Object} reservationData - Reservation data
   * @param {Array} reservationData.items - Array of items to reserve
   * @param {string} reservationData.items[].productId - Product ID
   * @param {number} reservationData.items[].quantity - Quantity to reserve
   * @param {string} reservationData.orderId - Order ID for the reservation
   * @param {string} [reservationData.notes] - Optional reservation notes
   * @returns {Promise<Object>} Reservation result
   */
  reserveInventory: async (reservationData) => {
    try {
      const response = await api.post(`${INVENTORY_BASE}/reserve`, reservationData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error reserving inventory:', error);
      throw error;
    }
  },

  /**
   * Release reserved inventory
   * @param {Object} releaseData - Release data
   * @param {Array} releaseData.items - Array of items to release
   * @param {string} releaseData.items[].productId - Product ID
   * @param {number} releaseData.items[].quantity - Quantity to release
   * @param {string} releaseData.orderId - Order ID for the release
   * @param {string} [releaseData.reason] - Reason for release
   * @param {string} [releaseData.notes] - Optional release notes
   * @returns {Promise<Object>} Release result
   */
  releaseInventory: async (releaseData) => {
    try {
      const response = await api.post(`${INVENTORY_BASE}/release`, releaseData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error releasing inventory:', error);
      throw error;
    }
  },

  /**
   * Get inventory reservations
   * @param {Object} params - Query parameters
   * @param {string} [params.orderId] - Filter by order ID
   * @param {string} [params.productId] - Filter by product ID
   * @param {string} [params.status] - Filter by reservation status
   * @param {number} [params.page] - Page number for pagination
   * @param {number} [params.limit] - Number of items per page
   * @returns {Promise<Object>} Paginated list of reservations
   */
  getReservations: async (params = {}) => {
    try {
      const response = await api.get(`${INVENTORY_BASE}/reservations`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching reservations:', error);
      throw error;
    }
  },

  /**
   * Get inventory movements/transactions
   * @param {Object} params - Query parameters
   * @param {string} [params.productId] - Filter by product ID
   * @param {string} [params.type] - Filter by movement type (IN, OUT, ADJUSTMENT, RESERVE, RELEASE)
   * @param {string} [params.startDate] - Filter movements from this date
   * @param {string} [params.endDate] - Filter movements until this date
   * @param {number} [params.page] - Page number for pagination
   * @param {number} [params.limit] - Number of items per page
   * @returns {Promise<Object>} Paginated list of inventory movements
   */
  getInventoryMovements: async (params = {}) => {
    try {
      const response = await api.get(`${INVENTORY_BASE}/movements`, { params });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching inventory movements:', error);
      throw error;
    }
  },

  /**
   * Perform bulk inventory update
   * @param {Object} bulkData - Bulk update data
   * @param {Array} bulkData.updates - Array of inventory updates
   * @param {string} bulkData.updates[].productId - Product ID
   * @param {number} bulkData.updates[].quantity - New quantity
   * @param {string} [bulkData.updates[].reason] - Reason for update
   * @param {string} [bulkData.notes] - Optional bulk update notes
   * @returns {Promise<Object>} Bulk update result
   */
  bulkUpdateInventory: async (bulkData) => {
    try {
      const response = await api.post(`${INVENTORY_BASE}/bulk-update`, bulkData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error performing bulk inventory update:', error);
      throw error;
    }
  },
};

export default inventoryService;
