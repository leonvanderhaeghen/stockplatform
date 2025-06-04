import api from './api';

// Base path for user endpoints (relative to the API base URL which already includes /v1)
const USERS_BASE = '/users';

const userService = {
  // Get all users (admin only)
  getUsers: async () => {
    try {
      const response = await api.get(USERS_BASE);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  },

  // Get current user profile
  getProfile: async () => {
    try {
      const response = await api.get(`${USERS_BASE}/me`);
      // The backend wraps the response in a 'data' field
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching user profile:', error);
      throw error;
    }
  },

  // Update current user profile
  updateProfile: async (userData) => {
    try {
      const response = await api.put(`${USERS_BASE}/me`, userData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating user profile:', error);
      throw error;
    }
  },

  // Change current user's password
  updatePassword: async (currentPassword, newPassword) => {
    try {
      const response = await api.put(`${USERS_BASE}/me/password`, {
        currentPassword,
        newPassword,
      });
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating password:', error);
      throw error;
    }
  },

  // Get user addresses
  getAddresses: async () => {
    try {
      const response = await api.get(`${USERS_BASE}/me/addresses`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error fetching addresses:', error);
      throw error;
    }
  },

  // Create a new address
  createAddress: async (addressData) => {
    try {
      const response = await api.post(`${USERS_BASE}/me/addresses`, addressData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error creating address:', error);
      throw error;
    }
  },

  // Update an address
  updateAddress: async (addressId, addressData) => {
    try {
      const response = await api.put(`${USERS_BASE}/me/addresses/${addressId}`, addressData);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error updating address:', error);
      throw error;
    }
  },

  // Delete an address
  deleteAddress: async (addressId) => {
    try {
      await api.delete(`${USERS_BASE}/me/addresses/${addressId}`);
      return addressId;
    } catch (error) {
      console.error('Error deleting address:', error);
      throw error;
    }
  },

  // Set default address
  setDefaultAddress: async (addressId) => {
    try {
      const response = await api.put(`${USERS_BASE}/me/addresses/${addressId}/default`);
      return response.data.data || response.data;
    } catch (error) {
      console.error('Error setting default address:', error);
      throw error;
    }
  },

  // Note: The backend doesn't support updating user roles directly through the API
  // Role updates should be done through the admin interface or backend directly
};

export default userService;
