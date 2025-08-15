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

const userService = {
  // Get current user profile
  getCurrentUser: async () => {
    const response = await api.get('/users/me');
    return response.data;
  },

  // Update user profile
  updateProfile: async (userData) => {
    const response = await api.put('/users/me', userData);
    return response.data;
  },

  // Change password
  changePassword: async (currentPassword, newPassword) => {
    const response = await api.put('/users/me/password', {
      currentPassword,
      newPassword,
    });
    return response.data;
  },

  // Get user addresses
  getAddresses: async () => {
    const response = await api.get('/users/me/addresses');
    return response.data;
  },

  // Create new address
  createAddress: async (addressData) => {
    const response = await api.post('/users/me/addresses', addressData);
    return response.data;
  },

  // Update address
  updateAddress: async (addressId, addressData) => {
    const response = await api.put(`/users/me/addresses/${addressId}`, addressData);
    return response.data;
  },

  // Delete address
  deleteAddress: async (addressId) => {
    const response = await api.delete(`/users/me/addresses/${addressId}`);
    return response.data;
  },

  // Set default address
  setDefaultAddress: async (addressId) => {
    const response = await api.put(`/users/me/addresses/${addressId}/default`);
    return response.data;
  },

  // Get default address
  getDefaultAddress: async () => {
    const response = await api.get('/users/me/addresses/default');
    return response.data;
  },

  // Admin: List all users
  listUsers: async (params = {}) => {
    const response = await api.get('/users', { params });
    return response.data;
  },

  // Admin: Get user by ID
  getUserById: async (userId) => {
    const response = await api.get(`/users/${userId}`);
    return response.data;
  },

  // Admin: Activate user
  activateUser: async (userId) => {
    const response = await api.put(`/users/${userId}/activate`);
    return response.data;
  },

  // Admin: Deactivate user
  deactivateUser: async (userId) => {
    const response = await api.put(`/users/${userId}/deactivate`);
    return response.data;
  },

  // Admin: Update user role
  updateUserRole: async (userId, role) => {
    const response = await api.put(`/users/${userId}/role`, { role });
    return response.data;
  },

  // Upload profile picture
  uploadProfilePicture: async (file) => {
    const formData = new FormData();
    formData.append('profilePicture', file);
    
    const response = await api.post('/users/me/profile-picture', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};

export default userService;
