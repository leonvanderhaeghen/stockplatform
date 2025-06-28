import api from './api';

// Base path for supplier endpoints (relative to the API base URL which already includes /v1)
const SUPPLIERS_BASE = '/suppliers';

const supplierService = {
  /**
   * Get all suppliers with optional filters
   * @param {Object} params - Query parameters for filtering and pagination
   * @param {number} params.page - Page number for pagination
   * @param {number} params.limit - Number of items per page
   * @param {string} params.name - Filter by supplier name
   * @param {string} params.email - Filter by supplier email
   * @param {boolean} params.active - Filter by active status
   * @returns {Promise<Object>} Paginated list of suppliers
   */
  getSuppliers: async (params = {}) => {
    try {
      const response = await api.get(SUPPLIERS_BASE, { params });
      return response.data;
    } catch (error) {
      console.error('Error fetching suppliers:', error);
      throw error;
    }
  },

  /**
   * Get a single supplier by ID
   * @param {string} id - Supplier ID
   * @returns {Promise<Object>} Supplier details
   */
  getSupplier: async (id) => {
    try {
      const response = await api.get(`${SUPPLIERS_BASE}/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching supplier:', error);
      throw error;
    }
  },

  /**
   * Create a new supplier
   * @param {Object} supplierData - Supplier data
   * @param {string} supplierData.name - Supplier name
   * @param {string} supplierData.contactPerson - Contact person name
   * @param {string} supplierData.email - Supplier email
   * @param {string} supplierData.phone - Supplier phone
   * @param {string} supplierData.address - Supplier address
   * @param {string} supplierData.city - Supplier city
   * @param {string} supplierData.country - Supplier country
   * @param {string} supplierData.postalCode - Supplier postal code
   * @param {boolean} supplierData.isActive - Whether the supplier is active
   * @param {Object} supplierData.metadata - Additional metadata
   * @returns {Promise<Object>} Created supplier details
   */
  createSupplier: async (supplierData) => {
    try {
      const response = await api.post(SUPPLIERS_BASE, supplierData);
      return response.data;
    } catch (error) {
      console.error('Error creating supplier:', error);
      throw error;
    }
  },

  /**
   * Update an existing supplier
   * @param {string} id - Supplier ID
   * @param {Object} supplierData - Updated supplier data
   * @returns {Promise<Object>} Updated supplier details
   */
  updateSupplier: async (id, supplierData) => {
    try {
      const response = await api.put(`${SUPPLIERS_BASE}/${id}`, supplierData);
      return response.data;
    } catch (error) {
      console.error('Error updating supplier:', error);
      throw error;
    }
  },

  /**
   * Delete a supplier
   * @param {string} id - Supplier ID
   * @returns {Promise<Object>} Deletion confirmation
   */
  deleteSupplier: async (id) => {
    try {
      const response = await api.delete(`${SUPPLIERS_BASE}/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting supplier:', error);
      throw error;
    }
  },

  /**
   * Sync supplier data (generic method - deprecated, use specific sync methods)
   * @param {string} id - Supplier ID
   * @returns {Promise<Object>} Sync result
   */
  syncSupplier: async (id) => {
    try {
      const response = await api.post(`${SUPPLIERS_BASE}/${id}/sync`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error syncing supplier:', error);
      throw error;
    }
  },

  /**
   * Sync supplier products
   * @param {string} id - Supplier ID
   * @param {Object} options - Sync options
   * @returns {Promise<Object>} Sync result
   */
  syncSupplierProducts: async (id, options = {}) => {
    try {
      const response = await api.post(`${SUPPLIERS_BASE}/${id}/sync/products`, options);
      return response.data;
    } catch (error) {
      console.error('Error syncing supplier products:', error);
      throw error;
    }
  },

  /**
   * Sync supplier inventory
   * @param {string} id - Supplier ID
   * @param {Object} options - Sync options
   * @returns {Promise<Object>} Sync result
   */
  syncSupplierInventory: async (id, options = {}) => {
    try {
      const response = await api.post(`${SUPPLIERS_BASE}/${id}/sync/inventory`, options);
      return response.data;
    } catch (error) {
      console.error('Error syncing supplier inventory:', error);
      throw error;
    }
  },

  // Supplier Adapter Operations
  /**
   * Get all available supplier adapters
   * @returns {Promise<Array>} List of supplier adapters
   */
  getSupplierAdapters: async () => {
    try {
      const response = await api.get(`${SUPPLIERS_BASE}/adapters`);
      return response.data;
    } catch (error) {
      console.error('Error fetching supplier adapters:', error);
      throw error;
    }
  },

  /**
   * Get capabilities of a specific supplier adapter
   * @param {string} adapterName - Name of the adapter
   * @returns {Promise<Object>} Adapter capabilities
   */
  getAdapterCapabilities: async (adapterName) => {
    try {
      const response = await api.get(`${SUPPLIERS_BASE}/adapters/${adapterName}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching adapter capabilities:', error);
      throw error;
    }
  },

  /**
   * Test connection to a supplier adapter
   * @param {string} adapterName - Name of the adapter
   * @param {Object} connectionConfig - Connection configuration
   * @returns {Promise<Object>} Connection test result
   */
  testAdapterConnection: async (adapterName, connectionConfig = {}) => {
    try {
      const response = await api.post(`${SUPPLIERS_BASE}/adapters/${adapterName}/test-connection`, connectionConfig);
      return response.data;
    } catch (error) {
      console.error('Error testing adapter connection:', error);
      throw error;
    }
  },
};

export default supplierService;
