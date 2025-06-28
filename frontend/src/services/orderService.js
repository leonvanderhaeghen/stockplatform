import api from './api';

// Base path for order endpoints (relative to the API base URL which already includes /v1)
const ORDERS_BASE = '/orders';

const orderService = {
  // Get all orders with optional filters
  getOrders: async (params = {}) => {
    const { data } = await api.get(ORDERS_BASE, { params });
    return data;
  },

  // Get a single order by ID
  getOrder: async (id) => {
    const { data } = await api.get(`${ORDERS_BASE}/${id}`);
    return data;
  },

  // Create a new order
  createOrder: async (orderData) => {
    const { data } = await api.post(ORDERS_BASE, orderData);
    return data;
  },

  // Update an existing order
  updateOrder: async (id, orderData) => {
    const { data } = await api.put(`${ORDERS_BASE}/${id}`, orderData);
    return data;
  },

  // Update order status
  updateOrderStatus: async (id, status) => {
    const { data } = await api.patch(`${ORDERS_BASE}/${id}/status`, { status });
    return data;
  },

  // Delete an order
  deleteOrder: async (id) => {
    const { data } = await api.delete(`${ORDERS_BASE}/${id}`);
    return data;
  },

  // Get orders by user ID
  getOrdersByUser: async (userId) => {
    const { data } = await api.get(`${ORDERS_BASE}/user/${userId}`);
    return data;
  },

  /**
   * Cancel an order
   * @param {string} id - Order ID
   * @param {Object} cancelData - Cancellation data
   * @param {string} cancelData.reason - Reason for cancellation
   * @param {string} [cancelData.notes] - Optional cancellation notes
   * @returns {Promise<Object>} Cancelled order details
   */
  cancelOrder: async (id, cancelData) => {
    try {
      const response = await api.post(`${ORDERS_BASE}/${id}/cancel`, cancelData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error cancelling order:', error);
      throw error;
    }
  },

  /**
   * Update order payment information
   * @param {string} id - Order ID
   * @param {Object} paymentData - Payment data
   * @param {string} paymentData.paymentMethod - Payment method
   * @param {string} paymentData.paymentStatus - Payment status
   * @param {string} [paymentData.transactionId] - Transaction ID
   * @param {string} [paymentData.notes] - Optional payment notes
   * @returns {Promise<Object>} Updated order details
   */
  updateOrderPayment: async (id, paymentData) => {
    try {
      const response = await api.put(`${ORDERS_BASE}/${id}/payment`, paymentData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating order payment:', error);
      throw error;
    }
  },

  /**
   * Add tracking information to an order
   * @param {string} id - Order ID
   * @param {Object} trackingData - Tracking data
   * @param {string} trackingData.trackingCode - Tracking code
   * @param {string} [trackingData.carrier] - Carrier name
   * @param {string} [trackingData.estimatedDelivery] - Estimated delivery date
   * @param {string} [trackingData.notes] - Optional tracking notes
   * @returns {Promise<Object>} Updated order details
   */
  addOrderTracking: async (id, trackingData) => {
    try {
      const response = await api.put(`${ORDERS_BASE}/${id}/tracking`, trackingData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error adding order tracking:', error);
      throw error;
    }
  },

  /**
   * Get order tracking information
   * @param {string} id - Order ID
   * @returns {Promise<Object>} Order tracking details
   */
  getOrderTracking: async (id) => {
    try {
      const response = await api.get(`${ORDERS_BASE}/${id}/tracking`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching order tracking:', error);
      throw error;
    }
  },

  /**
   * Get order history/timeline
   * @param {string} id - Order ID
   * @returns {Promise<Array>} Order history timeline
   */
  getOrderHistory: async (id) => {
    try {
      const response = await api.get(`${ORDERS_BASE}/${id}/history`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching order history:', error);
      throw error;
    }
  },
};

export default orderService;
