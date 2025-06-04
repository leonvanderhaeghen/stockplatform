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
    await api.delete(`${ORDERS_BASE}/${id}`);
    return id;
  },

  // Get orders by user ID
  getOrdersByUser: async (userId) => {
    const { data } = await api.get(`${ORDERS_BASE}/user/${userId}`);
    return data;
  },
};

export default orderService;
