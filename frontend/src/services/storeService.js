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

const storeService = {
  // Store Management
  // List all stores
  listStores: async (params = {}) => {
    const response = await api.get('/stores', { params });
    // Normalize response: backend returns { data: [...], success: true }
    // Frontend expects { stores: [...], total: number }
    const payload = response.data;
    return {
      stores: payload.data || [],
      total: payload.total || (payload.data ? payload.data.length : 0),
      success: payload.success
    };
  },

  // Backward-compatible alias used by pages
  getStores: async (params = {}) => {
    return storeService.listStores(params);
  },

  // Get store by ID
  getStore: async (storeId) => {
    const response = await api.get(`/stores/${storeId}`);
    // Normalize response for single item
    const payload = response.data;
    return payload.data || payload;
  },

  // Create new store
  createStore: async (storeData) => {
    const response = await api.post('/stores', storeData);
    return response.data;
  },

  // Update store
  updateStore: async (storeId, storeData) => {
    const response = await api.put(`/stores/${storeId}`, storeData);
    return response.data;
  },

  // Delete store
  deleteStore: async (storeId) => {
    const response = await api.delete(`/stores/${storeId}`);
    return response.data;
  },

  // Activate store
  activateStore: async (storeId) => {
    const response = await api.put(`/stores/${storeId}/activate`);
    return response.data;
  },

  // Deactivate store
  deactivateStore: async (storeId) => {
    const response = await api.put(`/stores/${storeId}/deactivate`);
    return response.data;
  },

  // Search stores
  searchStores: async (query, params = {}) => {
    const response = await api.get('/stores/search', {
      params: { q: query, ...params }
    });
    return response.data;
  },

  // Get nearby stores
  getNearbyStores: async (latitude, longitude, radius = 10, params = {}) => {
    const response = await api.get('/stores/nearby', {
      params: { lat: latitude, lng: longitude, radius, ...params }
    });
    return response.data;
  },

  // Store Inventory
  // Get store inventory
  getStoreInventory: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/inventory`, { params });
    return response.data;
  },

  // Get store inventory for specific product
  getStoreProductInventory: async (storeId, productId) => {
    const response = await api.get(`/stores/${storeId}/inventory/product/${productId}`);
    return response.data;
  },

  // Transfer inventory between stores
  transferInventory: async (fromStoreId, toStoreId, transferData) => {
    const response = await api.post(`/stores/${fromStoreId}/inventory/transfer/${toStoreId}`, transferData);
    return response.data;
  },

  // Get transfer history for store
  getStoreTransferHistory: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/inventory/transfers`, { params });
    return response.data;
  },

  // Update store inventory levels
  updateStoreInventory: async (storeId, inventoryUpdates) => {
    const response = await api.put(`/stores/${storeId}/inventory`, { updates: inventoryUpdates });
    return response.data;
  },

  // Store Staff Management
  // Get store staff
  getStoreStaff: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/staff`, { params });
    return response.data;
  },

  // Add staff to store
  addStoreStaff: async (storeId, userId, role = 'STAFF') => {
    const response = await api.post(`/stores/${storeId}/staff`, { userId, role });
    return response.data;
  },

  // Remove staff from store
  removeStoreStaff: async (storeId, userId) => {
    const response = await api.delete(`/stores/${storeId}/staff/${userId}`);
    return response.data;
  },

  // Update staff role in store
  updateStoreStaffRole: async (storeId, userId, role) => {
    const response = await api.put(`/stores/${storeId}/staff/${userId}`, { role });
    return response.data;
  },

  // Get staff schedule
  getStoreSchedule: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/schedule`, { params });
    return response.data;
  },

  // Update staff schedule
  updateStoreSchedule: async (storeId, scheduleData) => {
    const response = await api.put(`/stores/${storeId}/schedule`, scheduleData);
    return response.data;
  },

  // Store Orders and Sales
  // Get store orders
  getStoreOrders: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/orders`, { params });
    return response.data;
  },

  // Get store sales analytics
  getStoreSales: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/sales`, { params });
    return response.data;
  },

  // Get store daily sales summary
  getStoreDailySales: async (storeId, date) => {
    const response = await api.get(`/stores/${storeId}/sales/daily`, {
      params: { date: date.toISOString().split('T')[0] }
    });
    return response.data;
  },

  // Get store performance metrics
  getStorePerformance: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/performance`, { params });
    return response.data;
  },

  // Store Configuration
  // Get store configuration
  getStoreConfiguration: async (storeId) => {
    const response = await api.get(`/stores/${storeId}/configuration`);
    return response.data;
  },

  // Update store configuration
  updateStoreConfiguration: async (storeId, configData) => {
    const response = await api.put(`/stores/${storeId}/configuration`, configData);
    return response.data;
  },

  // Get store operating hours
  getStoreHours: async (storeId) => {
    const response = await api.get(`/stores/${storeId}/hours`);
    return response.data;
  },

  // Update store operating hours
  updateStoreHours: async (storeId, hoursData) => {
    const response = await api.put(`/stores/${storeId}/hours`, hoursData);
    return response.data;
  },

  // Store Equipment and POS
  // Get store POS terminals
  getStorePOSTerminals: async (storeId) => {
    const response = await api.get(`/stores/${storeId}/pos-terminals`);
    return response.data;
  },

  // Register POS terminal
  registerPOSTerminal: async (storeId, terminalData) => {
    const response = await api.post(`/stores/${storeId}/pos-terminals`, terminalData);
    return response.data;
  },

  // Update POS terminal
  updatePOSTerminal: async (storeId, terminalId, terminalData) => {
    const response = await api.put(`/stores/${storeId}/pos-terminals/${terminalId}`, terminalData);
    return response.data;
  },

  // Deactivate POS terminal
  deactivatePOSTerminal: async (storeId, terminalId) => {
    const response = await api.delete(`/stores/${storeId}/pos-terminals/${terminalId}`);
    return response.data;
  },

  // Get terminal status
  getPOSTerminalStatus: async (storeId, terminalId) => {
    const response = await api.get(`/stores/${storeId}/pos-terminals/${terminalId}/status`);
    return response.data;
  },

  // Store Analytics and Reporting
  // Get store analytics dashboard
  getStoreAnalytics: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/analytics`, { params });
    return response.data;
  },

  // Get store revenue analytics
  getStoreRevenue: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/analytics/revenue`, { params });
    return response.data;
  },

  // Get store customer analytics
  getStoreCustomerAnalytics: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/analytics/customers`, { params });
    return response.data;
  },

  // Get store product performance
  getStoreProductPerformance: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/analytics/products`, { params });
    return response.data;
  },

  // Store Reports
  // Generate store report
  generateStoreReport: async (storeId, reportType, params = {}) => {
    const response = await api.post(`/stores/${storeId}/reports/${reportType}`, params);
    return response.data;
  },

  // Get store report
  getStoreReport: async (storeId, reportId) => {
    const response = await api.get(`/stores/${storeId}/reports/${reportId}`);
    return response.data;
  },

  // Download store report
  downloadStoreReport: async (storeId, reportId, format = 'pdf') => {
    const response = await api.get(`/stores/${storeId}/reports/${reportId}/download`, {
      params: { format },
      responseType: 'blob'
    });
    return response.data;
  },

  // Store Customers
  // Get store customers
  getStoreCustomers: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/customers`, { params });
    return response.data;
  },

  // Get customer visit history for store
  getCustomerVisitHistory: async (storeId, customerId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/customers/${customerId}/visits`, { params });
    return response.data;
  },

  // Store Promotions and Campaigns
  // Get store promotions
  getStorePromotions: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/promotions`, { params });
    return response.data;
  },

  // Create store promotion
  createStorePromotion: async (storeId, promotionData) => {
    const response = await api.post(`/stores/${storeId}/promotions`, promotionData);
    return response.data;
  },

  // Update store promotion
  updateStorePromotion: async (storeId, promotionId, promotionData) => {
    const response = await api.put(`/stores/${storeId}/promotions/${promotionId}`, promotionData);
    return response.data;
  },

  // Delete store promotion
  deleteStorePromotion: async (storeId, promotionId) => {
    const response = await api.delete(`/stores/${storeId}/promotions/${promotionId}`);
    return response.data;
  },

  // Bulk Operations
  // Bulk update stores
  bulkUpdateStores: async (updates) => {
    const response = await api.put('/stores/bulk', { stores: updates });
    return response.data;
  },

  // Export stores data
  exportStores: async (format = 'csv', params = {}) => {
    const response = await api.get('/stores/export', {
      params: { format, ...params },
      responseType: 'blob'
    });
    return response.data;
  },

  // Store Maintenance and Support
  // Report store issue
  reportStoreIssue: async (storeId, issueData) => {
    const response = await api.post(`/stores/${storeId}/issues`, issueData);
    return response.data;
  },

  // Get store issues
  getStoreIssues: async (storeId, params = {}) => {
    const response = await api.get(`/stores/${storeId}/issues`, { params });
    return response.data;
  },

  // Update issue status
  updateIssueStatus: async (storeId, issueId, status, notes) => {
    const response = await api.put(`/stores/${storeId}/issues/${issueId}`, {
      status,
      notes
    });
    return response.data;
  },

  // Format store data for autocomplete
  formatStoresForAutocomplete: (stores) => {
    return stores.map(store => ({
      id: store.id,
      label: `${store.name} - ${store.address?.city || 'Unknown City'}`,
      name: store.name,
      address: store.address,
      ...store
    }));
  },
};

export default storeService;
