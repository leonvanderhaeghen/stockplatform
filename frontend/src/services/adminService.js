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

const adminService = {
  // System Dashboard and Analytics
  // Get system overview
  getSystemOverview: async () => {
    const response = await api.get('/admin/dashboard');
    return response.data;
  },

  // Get dashboard analytics (aggregated data for admin dashboard)
  getDashboardAnalytics: async () => {
    const response = await api.get('/admin/analytics/dashboard');
    return response.data;
  },

  // Get system health status
  getSystemHealth: async () => {
    const response = await api.get('/admin/health');
    return response.data;
  },

  // Get system metrics
  getSystemMetrics: async (params = {}) => {
    const response = await api.get('/admin/metrics', { params });
    return response.data;
  },

  // Get system performance data
  getSystemPerformance: async (params = {}) => {
    const response = await api.get('/admin/performance', { params });
    return response.data;
  },

  // User Management (Admin)
  // Get all users with advanced filtering
  getAllUsers: async (params = {}) => {
    const response = await api.get('/admin/users', { params });
    return response.data;
  },

  // Create new user (Admin)
  createUser: async (userData) => {
    const response = await api.post('/admin/users', userData);
    return response.data;
  },

  // Update user (Admin)
  updateUser: async (userId, userData) => {
    const response = await api.put(`/admin/users/${userId}`, userData);
    return response.data;
  },

  // Delete user (Admin)
  deleteUser: async (userId) => {
    const response = await api.delete(`/admin/users/${userId}`);
    return response.data;
  },

  // Bulk user operations
  bulkUpdateUsers: async (userUpdates) => {
    const response = await api.put('/admin/users/bulk', { users: userUpdates });
    return response.data;
  },

  // Reset user password
  resetUserPassword: async (userId, newPassword) => {
    const response = await api.put(`/admin/users/${userId}/password-reset`, {
      newPassword
    });
    return response.data;
  },

  // Force password change
  forcePasswordChange: async (userId) => {
    const response = await api.put(`/admin/users/${userId}/force-password-change`);
    return response.data;
  },

  // Get user activity logs
  getUserActivityLogs: async (userId, params = {}) => {
    const response = await api.get(`/admin/users/${userId}/activity`, { params });
    return response.data;
  },

  // Get user login history
  getUserLoginHistory: async (userId, params = {}) => {
    const response = await api.get(`/admin/users/${userId}/login-history`, { params });
    return response.data;
  },

  // Order Management (Admin)
  // Get all orders with admin privileges
  getAllOrders: async (params = {}) => {
    const response = await api.get('/admin/orders', { params });
    return response.data;
  },

  // Admin order operations
  adminUpdateOrder: async (orderId, orderData) => {
    const response = await api.put(`/admin/orders/${orderId}`, orderData);
    return response.data;
  },

  // Force cancel order
  forceCancelOrder: async (orderId, reason) => {
    const response = await api.put(`/admin/orders/${orderId}/force-cancel`, {
      reason
    });
    return response.data;
  },

  // Override order status
  overrideOrderStatus: async (orderId, status, reason) => {
    const response = await api.put(`/admin/orders/${orderId}/override-status`, {
      status,
      reason
    });
    return response.data;
  },

  // Get problematic orders
  getProblematicOrders: async (params = {}) => {
    const response = await api.get('/admin/orders/problematic', { params });
    return response.data;
  },

  // Get order disputes
  getOrderDisputes: async (params = {}) => {
    const response = await api.get('/admin/orders/disputes', { params });
    return response.data;
  },

  // Resolve order dispute
  resolveOrderDispute: async (disputeId, resolution) => {
    const response = await api.put(`/admin/orders/disputes/${disputeId}/resolve`, {
      resolution
    });
    return response.data;
  },

  // Inventory Management (Admin)
  // Get global inventory overview
  getGlobalInventoryOverview: async (params = {}) => {
    const response = await api.get('/admin/inventory/overview', { params });
    return response.data;
  },

  // Get inventory discrepancies
  getInventoryDiscrepancies: async (params = {}) => {
    const response = await api.get('/admin/inventory/discrepancies', { params });
    return response.data;
  },

  // Resolve inventory discrepancy
  resolveInventoryDiscrepancy: async (discrepancyId, resolution) => {
    const response = await api.put(`/admin/inventory/discrepancies/${discrepancyId}/resolve`, {
      resolution
    });
    return response.data;
  },

  // Global stock adjustment
  globalStockAdjustment: async (adjustmentData) => {
    const response = await api.post('/admin/inventory/global-adjustment', adjustmentData);
    return response.data;
  },

  // Get inventory audit trails
  getInventoryAuditTrails: async (params = {}) => {
    const response = await api.get('/admin/inventory/audit-trails', { params });
    return response.data;
  },

  // Product Management (Admin)
  // Admin product operations
  adminGetAllProducts: async (params = {}) => {
    const response = await api.get('/admin/products', { params });
    return response.data;
  },

  // Bulk product operations
  bulkProductOperations: async (operation, productIds, data = {}) => {
    const response = await api.post('/admin/products/bulk', {
      operation,
      productIds,
      data
    });
    return response.data;
  },

  // Approve pending products
  approveProducts: async (productIds) => {
    const response = await api.put('/admin/products/approve', { productIds });
    return response.data;
  },

  // Reject pending products
  rejectProducts: async (productIds, reason) => {
    const response = await api.put('/admin/products/reject', { productIds, reason });
    return response.data;
  },

  // Get product quality issues
  getProductQualityIssues: async (params = {}) => {
    const response = await api.get('/admin/products/quality-issues', { params });
    return response.data;
  },

  // System Configuration
  // Get system configuration
  getSystemConfiguration: async () => {
    const response = await api.get('/admin/configuration');
    return response.data;
  },

  // Update system configuration
  updateSystemConfiguration: async (configData) => {
    const response = await api.put('/admin/configuration', configData);
    return response.data;
  },

  // Get feature flags
  getFeatureFlags: async () => {
    const response = await api.get('/admin/feature-flags');
    return response.data;
  },

  // Update feature flags
  updateFeatureFlags: async (flags) => {
    const response = await api.put('/admin/feature-flags', flags);
    return response.data;
  },

  // System maintenance mode
  enableMaintenanceMode: async (message, duration) => {
    const response = await api.post('/admin/maintenance/enable', {
      message,
      duration
    });
    return response.data;
  },

  // Disable maintenance mode
  disableMaintenanceMode: async () => {
    const response = await api.post('/admin/maintenance/disable');
    return response.data;
  },

  // Audit Logs and Security
  // Get system audit logs
  getAuditLogs: async (params = {}) => {
    const response = await api.get('/admin/audit-logs', { params });
    return response.data;
  },

  // Get security events
  getSecurityEvents: async (params = {}) => {
    const response = await api.get('/admin/security/events', { params });
    return response.data;
  },

  // Get failed login attempts
  getFailedLoginAttempts: async (params = {}) => {
    const response = await api.get('/admin/security/failed-logins', { params });
    return response.data;
  },

  // Block IP address
  blockIP: async (ipAddress, reason, duration) => {
    const response = await api.post('/admin/security/block-ip', {
      ipAddress,
      reason,
      duration
    });
    return response.data;
  },

  // Unblock IP address
  unblockIP: async (ipAddress) => {
    const response = await api.delete(`/admin/security/block-ip/${ipAddress}`);
    return response.data;
  },

  // Get blocked IPs
  getBlockedIPs: async (params = {}) => {
    const response = await api.get('/admin/security/blocked-ips', { params });
    return response.data;
  },

  // Analytics and Reporting
  // Generate system report
  generateSystemReport: async (reportType, params = {}) => {
    const response = await api.post('/admin/reports/generate', {
      type: reportType,
      params
    });
    return response.data;
  },

  // Get available reports
  getAvailableReports: async () => {
    const response = await api.get('/admin/reports/available');
    return response.data;
  },

  // Get generated reports
  getGeneratedReports: async (params = {}) => {
    const response = await api.get('/admin/reports', { params });
    return response.data;
  },

  // Download report
  downloadReport: async (reportId, format = 'pdf') => {
    const response = await api.get(`/admin/reports/${reportId}/download`, {
      params: { format },
      responseType: 'blob'
    });
    return response.data;
  },

  // Get business analytics
  getBusinessAnalytics: async (params = {}) => {
    const response = await api.get('/admin/analytics/business', { params });
    return response.data;
  },

  // Get financial analytics
  getFinancialAnalytics: async (params = {}) => {
    const response = await api.get('/admin/analytics/financial', { params });
    return response.data;
  },

  // Get customer analytics
  getCustomerAnalytics: async (params = {}) => {
    const response = await api.get('/admin/analytics/customers', { params });
    return response.data;
  },

  // Notifications and Alerts
  // Get system alerts
  getSystemAlerts: async (params = {}) => {
    const response = await api.get('/admin/alerts', { params });
    return response.data;
  },

  // Create system alert
  createSystemAlert: async (alertData) => {
    const response = await api.post('/admin/alerts', alertData);
    return response.data;
  },

  // Update alert status
  updateAlertStatus: async (alertId, status) => {
    const response = await api.put(`/admin/alerts/${alertId}/status`, { status });
    return response.data;
  },

  // Delete alert
  deleteAlert: async (alertId) => {
    const response = await api.delete(`/admin/alerts/${alertId}`);
    return response.data;
  },

  // Send system notification
  sendSystemNotification: async (notificationData) => {
    const response = await api.post('/admin/notifications/send', notificationData);
    return response.data;
  },

  // Backup and Recovery
  // Create system backup
  createBackup: async (backupData) => {
    const response = await api.post('/admin/backup/create', backupData);
    return response.data;
  },

  // Get backup history
  getBackupHistory: async (params = {}) => {
    const response = await api.get('/admin/backup/history', { params });
    return response.data;
  },

  // Restore from backup
  restoreFromBackup: async (backupId) => {
    const response = await api.post(`/admin/backup/${backupId}/restore`);
    return response.data;
  },

  // Delete backup
  deleteBackup: async (backupId) => {
    const response = await api.delete(`/admin/backup/${backupId}`);
    return response.data;
  },

  // API Management
  // Get API usage statistics
  getAPIUsage: async (params = {}) => {
    const response = await api.get('/admin/api/usage', { params });
    return response.data;
  },

  // Get API rate limits
  getAPIRateLimits: async () => {
    const response = await api.get('/admin/api/rate-limits');
    return response.data;
  },

  // Update API rate limits
  updateAPIRateLimits: async (rateLimits) => {
    const response = await api.put('/admin/api/rate-limits', rateLimits);
    return response.data;
  },

  // Get API keys
  getAPIKeys: async (params = {}) => {
    const response = await api.get('/admin/api/keys', { params });
    return response.data;
  },

  // Create API key
  createAPIKey: async (keyData) => {
    const response = await api.post('/admin/api/keys', keyData);
    return response.data;
  },

  // Revoke API key
  revokeAPIKey: async (keyId) => {
    const response = await api.delete(`/admin/api/keys/${keyId}`);
    return response.data;
  },

  // Integration Management
  // Get system integrations
  getSystemIntegrations: async () => {
    const response = await api.get('/admin/integrations');
    return response.data;
  },

  // Configure integration
  configureIntegration: async (integrationType, configData) => {
    const response = await api.put(`/admin/integrations/${integrationType}`, configData);
    return response.data;
  },

  // Test integration
  testIntegration: async (integrationType) => {
    const response = await api.post(`/admin/integrations/${integrationType}/test`);
    return response.data;
  },

  // Enable/disable integration
  toggleIntegration: async (integrationType, enabled) => {
    const response = await api.put(`/admin/integrations/${integrationType}/toggle`, {
      enabled
    });
    return response.data;
  },

  // System Utilities
  // Clear system cache
  clearSystemCache: async () => {
    const response = await api.post('/admin/utilities/clear-cache');
    return response.data;
  },

  // Reindex search
  reindexSearch: async () => {
    const response = await api.post('/admin/utilities/reindex-search');
    return response.data;
  },

  // Clean up old data
  cleanupOldData: async (cleanupOptions) => {
    const response = await api.post('/admin/utilities/cleanup', cleanupOptions);
    return response.data;
  },

  // Optimize database
  optimizeDatabase: async () => {
    const response = await api.post('/admin/utilities/optimize-database');
    return response.data;
  },

  // Export system data
  exportSystemData: async (exportOptions) => {
    const response = await api.post('/admin/export', exportOptions, {
      responseType: 'blob'
    });
    return response.data;
  },
};

export default adminService;
