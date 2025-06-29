import api from './api';

const storeService = {
  // Get all stores with pagination
  getStores: async (params = {}) => {
    const { limit = 50, offset = 0 } = params;
    const response = await api.get('/stores', {
      params: { limit, offset }
    });
    return response.data;
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
    
    // Transform the response for autocomplete usage
    if (response.data && Array.isArray(response.data)) {
      return response.data.map(store => ({
        id: store.id,
        name: store.name,
        address: store.address,
        label: `${store.name} - ${store.address}`,
        value: store.id
      }));
    }
    
    return [];
  }
};

export default storeService;
