import api from './api';

const unwrap = (res) => (res && res.data ? res.data : res);
const BASE = '/orders';

const orderService = {
  /* ----------  customer  ---------- */
  getMyOrders: async (params = {}) => {
    const { data } = await api.get(`${BASE}/me`, { params });
    return unwrap(data);
  },

  getMyOrder: async (id) => {
    const { data } = await api.get(`${BASE}/me/${id}`);
    return unwrap(data);
  },

  createOrder: async (order) => {
    const { data } = await api.post(BASE, order);
    return unwrap(data);
  },

  /* ----------  admin / staff  ---------- */
  getOrders: async (params = {}) => {
    const { data } = await api.get(BASE, { params });
    return unwrap(data);
  },

  getOrder: async (id) => {
    const { data } = await api.get(`${BASE}/${id}`);
    return unwrap(data);
  },

  cancelOrder: async (id, reason = '') => {
    const { data } = await api.put(`${BASE}/${id}/cancel`, null, {
      params: { reason },
    });
    return unwrap(data);
  },
};

export default orderService;