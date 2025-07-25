import api from './api';

const storeService = {
  // Helper to normalize list responses which may come as array or {items: []}
  _extractList: (data) => {
    if (Array.isArray(data)) return data;
    if (data && Array.isArray(data.items)) return data.items;
    if (data && Array.isArray(data.data)) return data.data;
    return [];
  },

  // Get all stores with pagination
  getStores: async (params = {}) => {
    const { limit = 50, offset = 0 } = params;
    const response = await api.get('/stores', {
      params: { limit, offset }
    });
    return storeService._extractList(response.data);
  },

  // Get a specific store by ID
  getStore: async (id) => {
    const response = await api.get(`/stores/${id}`);
    return response.data;
  },

  // Get stores for autocomplete (simplified response)
  getStoresForAutocomplete: async (searchTerm = '') => {
    const response = await api.get('/stores', {
      params: { 
        limit: 20, 
        offset: 0,
        search: searchTerm 
      }
    });

    const list = storeService._extractList(response.data);

    return list.map(store => ({
      id: store.id,
      name: store.name,
      address: store.address,
      label: `${store.name} - ${store.address}`,
      value: store.id
    }));
  }
  ,

  // Create a new store
  createStore: async (storeData) => {
    const response = await api.post('/stores', storeData);
    return response.data;
  }
};

export default storeService;
