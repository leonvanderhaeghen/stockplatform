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

const supplierService = {
  // Supplier Management
  // List all suppliers
  listSuppliers: async (params = {}) => {
    const response = await api.get('/suppliers', { params });
    // Handle new standardized API response: { data: { suppliers: [], total_count: number }, success: boolean }
    if (response.data && typeof response.data === 'object' && response.data.data) {
      return {
        suppliers: Array.isArray(response.data.data.suppliers) ? response.data.data.suppliers : [],
        total: response.data.data.total_count || response.data.data.total || 0,
        success: response.data.success || true
      };
    }
    // Legacy fallback: handle old supplier API response format
    if (response.data && typeof response.data === 'object' && response.data.suppliers) {
      return {
        suppliers: Array.isArray(response.data.suppliers) ? response.data.suppliers : [],
        total: response.data.total_count || response.data.total || 0,
        success: true
      };
    }
    // Final fallback: wrap any other response in expected structure
    const suppliers = Array.isArray(response.data) ? response.data : [];
    return {
      suppliers: suppliers,
      total: suppliers.length,
      success: true
    };
  },

  // Backward-compatible alias used by pages (returns object with suppliers property)
  getSuppliers: async (params = {}) => {
    return supplierService.listSuppliers(params);
  },

  // New method for autocomplete components that need just the array
  getSuppliersArray: async (params = {}) => {
    const result = await supplierService.listSuppliers(params);
    return result.suppliers || [];
  },

  // Get supplier by ID
  getSupplier: async (supplierId) => {
    const response = await api.get(`/suppliers/${supplierId}`);
    return response.data;
  },

  // Create new supplier
  createSupplier: async (supplierData) => {
    const response = await api.post('/suppliers', supplierData);
    return response.data;
  },

  // Update supplier
  updateSupplier: async (supplierId, supplierData) => {
    const response = await api.put(`/suppliers/${supplierId}`, supplierData);
    return response.data;
  },

  // Delete supplier
  deleteSupplier: async (supplierId) => {
    const response = await api.delete(`/suppliers/${supplierId}`);
    return response.data;
  },

  // Activate supplier
  activateSupplier: async (supplierId) => {
    const response = await api.put(`/suppliers/${supplierId}/activate`);
    return response.data;
  },

  // Deactivate supplier
  deactivateSupplier: async (supplierId) => {
    const response = await api.put(`/suppliers/${supplierId}/deactivate`);
    return response.data;
  },

  // Search suppliers
  searchSuppliers: async (query, params = {}) => {
    const response = await api.get('/suppliers/search', {
      params: { q: query, ...params }
    });
    return response.data;
  },

  // Supplier Adapters
  // List all available adapters
  listAdapters: async () => {
    const response = await api.get('/suppliers/adapters');
    return response.data;
  },

  // Backward-compatible alias used by pages
  getAdapters: async () => {
    return supplierService.listAdapters();
  },

  // Get adapter capabilities
  getAdapterCapabilities: async (adapterType) => {
    const response = await api.get(`/suppliers/adapters/${adapterType}/capabilities`);
    return response.data;
  },

  // Test adapter connection
  testAdapterConnection: async (supplierId, connectionData) => {
    const response = await api.post(`/suppliers/${supplierId}/test-connection`, connectionData);
    return response.data;
  },

  // Configure supplier adapter
  configureAdapter: async (supplierId, adapterConfig) => {
    const response = await api.put(`/suppliers/${supplierId}/adapter`, adapterConfig);
    return response.data;
  },

  // Get supplier adapter configuration
  getAdapterConfiguration: async (supplierId) => {
    const response = await api.get(`/suppliers/${supplierId}/adapter`);
    return response.data;
  },

  // Product Synchronization
  // Sync products from supplier
  syncProducts: async (supplierId, syncOptions = {}) => {
    const response = await api.post(`/suppliers/${supplierId}/sync/products`, syncOptions);
    return response.data;
  },

  // Sync inventory from supplier
  syncInventory: async (supplierId, syncOptions = {}) => {
    const response = await api.post(`/suppliers/${supplierId}/sync/inventory`, syncOptions);
    return response.data;
  },

  // Get sync history
  getSyncHistory: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/sync/history`, { params });
    return response.data;
  },

  // Get sync status
  getSyncStatus: async (supplierId, syncId) => {
    const response = await api.get(`/suppliers/${supplierId}/sync/${syncId}/status`);
    return response.data;
  },

  // Cancel sync operation
  cancelSync: async (supplierId, syncId) => {
    const response = await api.delete(`/suppliers/${supplierId}/sync/${syncId}`);
    return response.data;
  },

  // Schedule automatic sync
  scheduleSync: async (supplierId, scheduleData) => {
    const response = await api.post(`/suppliers/${supplierId}/sync/schedule`, scheduleData);
    return response.data;
  },

  // Get sync schedule
  getSyncSchedule: async (supplierId) => {
    const response = await api.get(`/suppliers/${supplierId}/sync/schedule`);
    return response.data;
  },

  // Update sync schedule
  updateSyncSchedule: async (supplierId, scheduleData) => {
    const response = await api.put(`/suppliers/${supplierId}/sync/schedule`, scheduleData);
    return response.data;
  },

  // Delete sync schedule
  deleteSyncSchedule: async (supplierId) => {
    const response = await api.delete(`/suppliers/${supplierId}/sync/schedule`);
    return response.data;
  },

  // Supplier Products
  // Get products from supplier
  getSupplierProducts: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/products`, { params });
    return response.data;
  },

  // Map supplier product to internal product
  mapSupplierProduct: async (supplierId, supplierProductId, internalProductId) => {
    const response = await api.post(`/suppliers/${supplierId}/products/${supplierProductId}/map`, {
      internalProductId
    });
    return response.data;
  },

  // Unmap supplier product
  unmapSupplierProduct: async (supplierId, supplierProductId) => {
    const response = await api.delete(`/suppliers/${supplierId}/products/${supplierProductId}/map`);
    return response.data;
  },

  // Get product mapping
  getProductMapping: async (supplierId, supplierProductId) => {
    const response = await api.get(`/suppliers/${supplierId}/products/${supplierProductId}/map`);
    return response.data;
  },

  // Bulk map products
  bulkMapProducts: async (supplierId, mappings) => {
    const response = await api.post(`/suppliers/${supplierId}/products/bulk-map`, {
      mappings
    });
    return response.data;
  },

  // Purchase Orders
  // Create purchase order
  createPurchaseOrder: async (supplierId, orderData) => {
    const response = await api.post(`/suppliers/${supplierId}/purchase-orders`, orderData);
    return response.data;
  },

  // Get purchase orders
  getPurchaseOrders: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/purchase-orders`, { params });
    return response.data;
  },

  // Get purchase order by ID
  getPurchaseOrder: async (supplierId, orderId) => {
    const response = await api.get(`/suppliers/${supplierId}/purchase-orders/${orderId}`);
    return response.data;
  },

  // Update purchase order
  updatePurchaseOrder: async (supplierId, orderId, orderData) => {
    const response = await api.put(`/suppliers/${supplierId}/purchase-orders/${orderId}`, orderData);
    return response.data;
  },

  // Cancel purchase order
  cancelPurchaseOrder: async (supplierId, orderId, reason) => {
    const response = await api.put(`/suppliers/${supplierId}/purchase-orders/${orderId}/cancel`, {
      reason
    });
    return response.data;
  },

  // Submit purchase order to supplier
  submitPurchaseOrder: async (supplierId, orderId) => {
    const response = await api.post(`/suppliers/${supplierId}/purchase-orders/${orderId}/submit`);
    return response.data;
  },

  // Receive purchase order
  receivePurchaseOrder: async (supplierId, orderId, receiptData) => {
    const response = await api.post(`/suppliers/${supplierId}/purchase-orders/${orderId}/receive`, receiptData);
    return response.data;
  },

  // Supplier Analytics
  // Get supplier performance analytics
  getSupplierAnalytics: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/analytics`, { params });
    return response.data;
  },

  // Get supplier delivery performance
  getDeliveryPerformance: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/analytics/delivery`, { params });
    return response.data;
  },

  // Get supplier quality metrics
  getQualityMetrics: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/analytics/quality`, { params });
    return response.data;
  },

  // Get supplier cost analysis
  getCostAnalysis: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/analytics/cost`, { params });
    return response.data;
  },

  // Supplier Contacts
  // Get supplier contacts
  getSupplierContacts: async (supplierId) => {
    const response = await api.get(`/suppliers/${supplierId}/contacts`);
    return response.data;
  },

  // Add supplier contact
  addSupplierContact: async (supplierId, contactData) => {
    const response = await api.post(`/suppliers/${supplierId}/contacts`, contactData);
    return response.data;
  },

  // Update supplier contact
  updateSupplierContact: async (supplierId, contactId, contactData) => {
    const response = await api.put(`/suppliers/${supplierId}/contacts/${contactId}`, contactData);
    return response.data;
  },

  // Delete supplier contact
  deleteSupplierContact: async (supplierId, contactId) => {
    const response = await api.delete(`/suppliers/${supplierId}/contacts/${contactId}`);
    return response.data;
  },

  // Set primary contact
  setPrimaryContact: async (supplierId, contactId) => {
    const response = await api.put(`/suppliers/${supplierId}/contacts/${contactId}/primary`);
    return response.data;
  },

  // Supplier Documents
  // Upload supplier document
  uploadDocument: async (supplierId, file, documentType, description) => {
    const formData = new FormData();
    formData.append('document', file);
    formData.append('type', documentType);
    formData.append('description', description);
    
    const response = await api.post(`/suppliers/${supplierId}/documents`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Get supplier documents
  getSupplierDocuments: async (supplierId, params = {}) => {
    const response = await api.get(`/suppliers/${supplierId}/documents`, { params });
    return response.data;
  },

  // Download supplier document
  downloadDocument: async (supplierId, documentId) => {
    const response = await api.get(`/suppliers/${supplierId}/documents/${documentId}/download`, {
      responseType: 'blob'
    });
    return response.data;
  },

  // Delete supplier document
  deleteDocument: async (supplierId, documentId) => {
    const response = await api.delete(`/suppliers/${supplierId}/documents/${documentId}`);
    return response.data;
  },

  // Bulk Operations
  // Bulk update suppliers
  bulkUpdateSuppliers: async (updates) => {
    const response = await api.put('/suppliers/bulk', { suppliers: updates });
    return response.data;
  },

  // Bulk delete suppliers
  bulkDeleteSuppliers: async (supplierIds) => {
    const response = await api.delete('/suppliers/bulk', { data: { supplierIds } });
    return response.data;
  },

  // Export suppliers
  exportSuppliers: async (format = 'csv', params = {}) => {
    const response = await api.get('/suppliers/export', {
      params: { format, ...params },
      responseType: 'blob'
    });
    return response.data;
  },

  // Import suppliers
  importSuppliers: async (file, options = {}) => {
    const formData = new FormData();
    formData.append('file', file);
    Object.keys(options).forEach(key => {
      formData.append(key, options[key]);
    });
    
    const response = await api.post('/suppliers/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Get import status
  getImportStatus: async (importId) => {
    const response = await api.get(`/suppliers/import/${importId}/status`);
    return response.data;
  },
};

export default supplierService;
