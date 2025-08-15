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

const inventoryService = {
  // List inventory with filtering and pagination
  listInventory: async (params = {}) => {
    const response = await api.get('/inventory', { params });
    // Normalize response: backend returns { data: [...], success: true }
    // Frontend expects { items: [...] }
    const payload = response.data;
    return {
      items: payload.data || [],
      success: payload.success
    };
  },

  // Backward-compatible alias
  getInventory: async (params = {}) => {
    return inventoryService.listInventory(params);
  },

  // Get inventory item by ID
  getInventoryItemById: async (itemId) => {
    const response = await api.get(`/inventory/${itemId}`);
    // Normalize response for single item
    const payload = response.data;
    return payload.data || payload;
  },

  // Get inventory items by product ID
  getInventoryItemsByProductId: async (productId, params = {}) => {
    const response = await api.get(`/inventory/product/${productId}`, { params });
    // Normalize response
    const payload = response.data;
    return {
      items: payload.data || [],
      success: payload.success
    };
  },

  // Get inventory item by SKU
  getInventoryItemBySKU: async (sku) => {
    const response = await api.get(`/inventory/sku/${sku}`);
    return response.data;
  },

  // Create new inventory item
  createInventoryItem: async (itemData) => {
    const response = await api.post('/inventory', itemData);
    return response.data;
  },

  // Update inventory item
  updateInventoryItem: async (itemId, itemData) => {
    const response = await api.put(`/inventory/${itemId}`, itemData);
    return response.data;
  },

  // Delete inventory item
  deleteInventoryItem: async (itemId) => {
    const response = await api.delete(`/inventory/${itemId}`);
    return response.data;
  },

  // Stock Management
  // Add stock to inventory item
  addStock: async (itemId, stockOrQty, maybeReason) => {
    const stockData =
      typeof stockOrQty === 'number'
        ? { quantity: stockOrQty, reason: maybeReason }
        : stockOrQty;
    const response = await api.post(`/inventory/${itemId}/stock/add`, stockData);
    return response.data;
  },

  // Remove stock from inventory item
  removeStock: async (itemId, stockOrQty, maybeReason) => {
    const stockData =
      typeof stockOrQty === 'number'
        ? { quantity: stockOrQty, reason: maybeReason }
        : stockOrQty;
    const response = await api.post(`/inventory/${itemId}/stock/remove`, stockData);
    return response.data;
  },

  // Adjust stock levels
  adjustStock: async (itemId, adjustmentData) => {
    const response = await api.post(`/inventory/${itemId}/stock/adjust`, adjustmentData);
    return response.data;
  },

  // Bulk stock operations
  bulkStockUpdate: async (updates) => {
    const response = await api.post('/inventory/stock/bulk-update', { updates });
    return response.data;
  },

  // Get stock movements/history
  getStockMovements: async (itemId, params = {}) => {
    const response = await api.get(`/inventory/${itemId}/movements`, { params });
    return response.data;
  },

  // Inventory Reservations (for POS and order processing)
  // Create inventory reservation
  createReservation: async (reservationData) => {
    const response = await api.post('/inventory/reservations', reservationData);
    return response.data;
  },

  // Get reservations
  getReservations: async (params = {}) => {
    const response = await api.get('/inventory/reservations', { params });
    return response.data;
  },

  // Backward-compatible alias
  getInventoryReservations: async (params = {}) => {
    return inventoryService.getReservations(params);
  },

  // Get reservation by ID
  getReservationById: async (reservationId) => {
    const response = await api.get(`/inventory/reservations/${reservationId}`);
    return response.data;
  },

  // Update reservation
  updateReservation: async (reservationId, reservationData) => {
    const response = await api.put(`/inventory/reservations/${reservationId}`, reservationData);
    return response.data;
  },

  // Cancel reservation
  cancelReservation: async (reservationId) => {
    const response = await api.delete(`/inventory/reservations/${reservationId}`);
    return response.data;
  },

  // Confirm/fulfill reservation
  fulfillReservation: async (reservationId, fulfillmentData) => {
    const response = await api.post(`/inventory/reservations/${reservationId}/fulfill`, fulfillmentData);
    return response.data;
  },

  // POS-specific inventory operations (using consolidated endpoints)
  // Check inventory availability for POS
  checkPOSAvailability: async (items, storeId) => {
    const response = await api.post('/inventory/availability', {
      items,
      storeId,
      source: 'POS'
    });
    return response.data;
  },

  // Reserve inventory for POS transaction
  reserveForPOS: async (items, orderId, storeId) => {
    const response = await api.post('/inventory/reservations', {
      items,
      orderId,
      storeId,
      source: 'POS',
      type: 'POS_TRANSACTION'
    });
    return response.data;
  },

  // Complete POS inventory deduction
  completePOSDeduction: async (reservationId, deductionData) => {
    const response = await api.post(`/inventory/reservations/${reservationId}/fulfill`, {
      ...deductionData,
      source: 'POS'
    });
    return response.data;
  },

  // Direct POS stock deduction (for quick transactions)
  directPOSDeduction: async (items, storeId, reason = 'POS Sale') => {
    const response = await api.post('/inventory/stock/deduct', {
      items,
      storeId,
      source: 'POS',
      reason
    });
    return response.data;
  },

  // Low Stock Management
  // Get low stock items
  getLowStockItems: async (params = {}) => {
    const response = await api.get('/inventory/low-stock', { params });
    return response.data;
  },

  // Update low stock thresholds
  updateLowStockThreshold: async (itemId, threshold) => {
    const response = await api.put(`/inventory/${itemId}/threshold`, { threshold });
    return response.data;
  },

  // Get stock alerts
  getStockAlerts: async (params = {}) => {
    const response = await api.get('/inventory/alerts', { params });
    return response.data;
  },

  // Mark alert as resolved
  resolveStockAlert: async (alertId) => {
    const response = await api.put(`/inventory/alerts/${alertId}/resolve`);
    return response.data;
  },

  // Inventory Analytics
  // Get inventory analytics
  getInventoryAnalytics: async (params = {}) => {
    const response = await api.get('/inventory/analytics', { params });
    return response.data;
  },

  // Convenience alias used by Dashboard
  getInventoryStats: async (params = {}) => {
    return inventoryService.getInventoryAnalytics(params);
  },

  // Get inventory turnover rates
  getInventoryTurnover: async (params = {}) => {
    const response = await api.get('/inventory/analytics/turnover', { params });
    return response.data;
  },

  // Get inventory valuation
  getInventoryValuation: async (params = {}) => {
    const response = await api.get('/inventory/analytics/valuation', { params });
    return response.data;
  },

  // Location-based Inventory
  // Get inventory by location/store
  getInventoryByLocation: async (locationId, params = {}) => {
    const response = await api.get(`/inventory/location/${locationId}`, { params });
    return response.data;
  },

  // Transfer stock between locations
  transferStock: async (transferData) => {
    const response = await api.post('/inventory/transfers', transferData);
    return response.data;
  },

  // Get transfer history
  getTransferHistory: async (params = {}) => {
    const response = await api.get('/inventory/transfers', { params });
    return response.data;
  },

  // Get transfer by ID
  getTransferById: async (transferId) => {
    const response = await api.get(`/inventory/transfers/${transferId}`);
    return response.data;
  },

  // Update transfer status
  updateTransferStatus: async (transferId, status, notes) => {
    const response = await api.put(`/inventory/transfers/${transferId}/status`, {
      status,
      notes
    });
    return response.data;
  },

  // Inventory Audits
  // Create inventory audit
  createAudit: async (auditData) => {
    const response = await api.post('/inventory/audits', auditData);
    return response.data;
  },

  // Get audits
  getAudits: async (params = {}) => {
    const response = await api.get('/inventory/audits', { params });
    return response.data;
  },

  // Get audit by ID
  getAuditById: async (auditId) => {
    const response = await api.get(`/inventory/audits/${auditId}`);
    return response.data;
  },

  // Update audit
  updateAudit: async (auditId, auditData) => {
    const response = await api.put(`/inventory/audits/${auditId}`, auditData);
    return response.data;
  },

  // Complete audit
  completeAudit: async (auditId, results) => {
    const response = await api.post(`/inventory/audits/${auditId}/complete`, results);
    return response.data;
  },

  // Search inventory
  searchInventory: async (query, params = {}) => {
    const response = await api.get('/inventory/search', {
      params: { q: query, ...params }
    });
    return response.data;
  },

  // Export inventory data
  exportInventory: async (format = 'csv', params = {}) => {
    const response = await api.get('/inventory/export', {
      params: { format, ...params },
      responseType: 'blob'
    });
    return response.data;
  },
};

export default inventoryService;
