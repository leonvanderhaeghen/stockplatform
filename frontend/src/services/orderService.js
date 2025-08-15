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

const orderService = {
  // Customer: Get user's orders
  getUserOrders: async (params = {}) => {
    const response = await api.get('/orders/me', { params });
    return response.data;
  },

  // Customer: Get specific user order
  getUserOrder: async (orderId) => {
    const response = await api.get(`/orders/me/${orderId}`);
    return response.data;
  },

  // Staff/Admin: List all orders
  listOrders: async (params = {}) => {
    const response = await api.get('/orders', { params });
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

  // Get order by ID
  getOrderById: async (orderId) => {
    const response = await api.get(`/orders/${orderId}`);
    return response.data;
  },

  // Create new order
  createOrder: async (orderData) => {
    const response = await api.post('/orders', orderData);
    return response.data;
  },

  // Update order status
  updateOrderStatus: async (orderId, status, notes) => {
    const response = await api.put(`/orders/${orderId}/status`, {
      status,
      notes
    });
    return response.data;
  },

  // Cancel order
  cancelOrder: async (orderId, reason) => {
    const response = await api.put(`/orders/${orderId}/cancel`, {
      reason
    });
    return response.data;
  },

  // Add payment to order
  addOrderPayment: async (orderId, paymentData) => {
    const response = await api.post(`/orders/${orderId}/payments`, paymentData);
    return response.data;
  },

  // Get order payments
  getOrderPayments: async (orderId) => {
    const response = await api.get(`/orders/${orderId}/payments`);
    return response.data;
  },

  // Add tracking information
  addOrderTracking: async (orderId, trackingData) => {
    const response = await api.post(`/orders/${orderId}/tracking`, trackingData);
    return response.data;
  },

  // Get order tracking
  getOrderTracking: async (orderId) => {
    const response = await api.get(`/orders/${orderId}/tracking`);
    return response.data;
  },

  // Update order items
  updateOrderItems: async (orderId, items) => {
    const response = await api.put(`/orders/${orderId}/items`, { items });
    return response.data;
  },

  // Add item to order
  addOrderItem: async (orderId, itemData) => {
    const response = await api.post(`/orders/${orderId}/items`, itemData);
    return response.data;
  },

  // Remove item from order
  removeOrderItem: async (orderId, itemId) => {
    const response = await api.delete(`/orders/${orderId}/items/${itemId}`);
    return response.data;
  },

  // Update order shipping address
  updateShippingAddress: async (orderId, addressData) => {
    const response = await api.put(`/orders/${orderId}/shipping-address`, addressData);
    return response.data;
  },

  // Update order billing address
  updateBillingAddress: async (orderId, addressData) => {
    const response = await api.put(`/orders/${orderId}/billing-address`, addressData);
    return response.data;
  },

  // POS Orders (using consolidated endpoints with source parameter)
  // Create POS order
  createPOSOrder: async (orderData) => {
    const response = await api.post('/orders', {
      ...orderData,
      source: 'POS'
    });
    return response.data;
  },

  // Process quick POS transaction
  processQuickPOSTransaction: async (transactionData) => {
    const response = await api.post('/orders', {
      ...transactionData,
      source: 'POS',
      type: 'QUICK_SALE'
    });
    return response.data;
  },

  // Get POS orders for a store
  getPOSOrders: async (storeId, params = {}) => {
    const response = await api.get('/orders', {
      params: {
        ...params,
        source: 'POS',
        storeId
      }
    });
    return response.data;
  },

  // Process POS payment
  processPOSPayment: async (orderId, paymentData) => {
    const response = await api.post(`/orders/${orderId}/payments`, {
      ...paymentData,
      source: 'POS'
    });
    return response.data;
  },

  // Complete POS order
  completePOSOrder: async (orderId, completionData) => {
    const response = await api.put(`/orders/${orderId}/status`, {
      status: 'COMPLETED',
      source: 'POS',
      ...completionData
    });
    return response.data;
  },

  // Order Analytics and Reporting
  // Get order analytics
  getOrderAnalytics: async (params = {}) => {
    const response = await api.get('/orders/analytics', { params });
    return response.data;
  },

  // Get sales analytics
  getSalesAnalytics: async (params = {}) => {
    const response = await api.get('/orders/analytics/sales', { params });
    return response.data;
  },

  // Get order trends
  getOrderTrends: async (params = {}) => {
    const response = await api.get('/orders/analytics/trends', { params });
    return response.data;
  },

  // Get revenue analytics
  getRevenueAnalytics: async (params = {}) => {
    const response = await api.get('/orders/analytics/revenue', { params });
    return response.data;
  },

  // Order Search and Filtering
  // Search orders
  searchOrders: async (query, params = {}) => {
    const response = await api.get('/orders/search', {
      params: { q: query, ...params }
    });
    return response.data;
  },

  // Get orders by status
  getOrdersByStatus: async (status, params = {}) => {
    const response = await api.get('/orders', {
      params: { status, ...params }
    });
    return response.data;
  },

  // Get orders by date range
  getOrdersByDateRange: async (startDate, endDate, params = {}) => {
    const response = await api.get('/orders', {
      params: {
        startDate: startDate.toISOString(),
        endDate: endDate.toISOString(),
        ...params
      }
    });
    return response.data;
  },

  // Get orders by customer
  getOrdersByCustomer: async (customerId, params = {}) => {
    const response = await api.get('/orders', {
      params: { customerId, ...params }
    });
    return response.data;
  },

  // Bulk Operations
  // Bulk update order status
  bulkUpdateOrderStatus: async (orderIds, status, notes) => {
    const response = await api.put('/orders/bulk/status', {
      orderIds,
      status,
      notes
    });
    return response.data;
  },

  // Bulk cancel orders
  bulkCancelOrders: async (orderIds, reason) => {
    const response = await api.put('/orders/bulk/cancel', {
      orderIds,
      reason
    });
    return response.data;
  },

  // Export orders
  exportOrders: async (format = 'csv', params = {}) => {
    const response = await api.get('/orders/export', {
      params: { format, ...params },
      responseType: 'blob'
    });
    return response.data;
  },

  // Returns and Refunds
  // Create return request
  createReturn: async (orderId, returnData) => {
    const response = await api.post(`/orders/${orderId}/returns`, returnData);
    return response.data;
  },

  // Get order returns
  getOrderReturns: async (orderId) => {
    const response = await api.get(`/orders/${orderId}/returns`);
    return response.data;
  },

  // Update return status
  updateReturnStatus: async (orderId, returnId, status, notes) => {
    const response = await api.put(`/orders/${orderId}/returns/${returnId}/status`, {
      status,
      notes
    });
    return response.data;
  },

  // Process refund
  processRefund: async (orderId, returnId, refundData) => {
    const response = await api.post(`/orders/${orderId}/returns/${returnId}/refund`, refundData);
    return response.data;
  },

  // Order History and Timeline
  // Get order history
  getOrderHistory: async (orderId) => {
    const response = await api.get(`/orders/${orderId}/history`);
    return response.data;
  },

  // Add order note
  addOrderNote: async (orderId, note, isPublic = false) => {
    const response = await api.post(`/orders/${orderId}/notes`, {
      note,
      isPublic
    });
    return response.data;
  },

  // Get order notes
  getOrderNotes: async (orderId) => {
    const response = await api.get(`/orders/${orderId}/notes`);
    return response.data;
  },

  // Update order note
  updateOrderNote: async (orderId, noteId, note, isPublic) => {
    const response = await api.put(`/orders/${orderId}/notes/${noteId}`, {
      note,
      isPublic
    });
    return response.data;
  },

  // Delete order note
  deleteOrderNote: async (orderId, noteId) => {
    const response = await api.delete(`/orders/${orderId}/notes/${noteId}`);
    return response.data;
  },

  // Order Notifications
  // Send order notification
  sendOrderNotification: async (orderId, notificationData) => {
    const response = await api.post(`/orders/${orderId}/notifications`, notificationData);
    return response.data;
  },

  // Get order notifications
  getOrderNotifications: async (orderId) => {
    const response = await api.get(`/orders/${orderId}/notifications`);
    return response.data;
  },

  // Order Validation
  // Validate order data
  validateOrder: async (orderData) => {
    const response = await api.post('/orders/validate', orderData);
    return response.data;
  },

  // Calculate order totals
  calculateOrderTotals: async (orderData) => {
    const response = await api.post('/orders/calculate', orderData);
    return response.data;
  },
};

export default orderService;
