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
  getLowStockItems: async (threshold = 10) => {
    const { data } = await api.get(`${INVENTORY_BASE}/low-stock`, { params: { threshold } });
    return data;
  },
};

export default inventoryService;
