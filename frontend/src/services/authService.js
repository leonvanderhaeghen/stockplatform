import api from './api';

// Helper function to map role IDs to role strings
const mapRoleIdToString = (roleId) => {
  const roleMap = {
    0: 'ADMIN',
    1: 'CUSTOMER',
    2: 'MANAGER',
    3: 'USER'
  };
  return roleMap[roleId] || 'USER';
};

// Base paths for API endpoints - relative to the API base URL (which already includes /v1)
const AUTH_BASE = '/auth';
const USERS_BASE = '/users';

// Helper function to log API requests
const logRequest = (method, url, data = {}) => {
  console.log(`[authService] ${method.toUpperCase()} ${url}`, {
    ...(Object.keys(data).length > 0 && { data })
  });
};

// Helper function to log API responses
const logResponse = (method, url, response) => {
  console.log(`[authService] ${method.toUpperCase()} ${url} response:`, {
    status: response?.status,
    data: response?.data,
  });
  return response;
};

// Helper function to log API errors
const logError = (method, url, error) => {
  console.error(`[authService] ${method.toUpperCase()} ${url} error:`, {
    message: error.message,
    response: error.response ? {
      status: error.response.status,
      data: error.response.data,
      headers: error.response.headers,
    } : 'No response',
    request: error.request ? 'Request made but no response received' : 'No request was made',
    config: {
      url: error.config?.url,
      method: error.config?.method,
      headers: error.config?.headers,
      data: error.config?.data,
    },
  });
  throw error; // Re-throw to allow error handling up the chain
};

const authService = {
  // Login user
  login: async (email, password) => {
    const url = `${AUTH_BASE}/login`;
    try {
      logRequest('POST', url, { email, password: '***' });
      const response = await api.post(url, { email, password });
      logResponse('POST', url, response);
      
      // The backend wraps the response in a 'data' field
      const responseData = response.data.data || response.data;
      const { token, user } = responseData;
      
      if (!token || !user) {
        throw new Error('Invalid response format from server');
      }
      
      // Map the role if it's a number
      if (typeof user.role === 'number') {
        user.role = mapRoleIdToString(user.role);
      } else if (typeof user.role === 'string') {
        // Ensure the role is uppercase to match our constants
        user.role = user.role.toUpperCase();
      }
      
      return {
        token,
        user
      };
    } catch (error) {
      logError('POST', url, error);
    }
  },

  // Logout user
  logout: async () => {
    const url = `${AUTH_BASE}/logout`;
    try {
      logRequest('POST', url);
      // The backend doesn't require a body for logout, just the token in the Authorization header
      const response = await api.post(url, {});
      logResponse('POST', url, response);
      return response.data;
    } catch (error) {
      logError('POST', url, error);
    }
  },

  // Refresh access token
  refreshToken: async (refreshToken) => {
    const url = `${AUTH_BASE}/refresh`;
    try {
      logRequest('POST', url, { refreshToken: '***' });
      const response = await api.post(url, { refreshToken });
      logResponse('POST', url, response);
      
      // The backend wraps the response in a 'data' field
      const responseData = response.data.data || response.data;
      
      // If there's a user object in the response, ensure the role is properly formatted
      if (responseData.user) {
        if (typeof responseData.user.role === 'number') {
          responseData.user.role = mapRoleIdToString(responseData.user.role);
        } else if (typeof responseData.user.role === 'string') {
          responseData.user.role = responseData.user.role.toUpperCase();
        }
      }
      
      return responseData;
    } catch (error) {
      logError('POST', url, error);
    }
  },

  // Get current user
  getCurrentUser: async () => {
    const url = `${USERS_BASE}/me`;
    try {
      logRequest('GET', url);
      const response = await api.get(url);
      logResponse('GET', url, response);
      
      // The backend wraps the response in a 'data' field
      const userData = response.data.data || response.data;
      
      // Ensure we have a valid user object
      if (!userData || typeof userData !== 'object' || !userData.id) {
        throw new Error('Invalid user data received from server');
      }
      
      // Map the role if it's a number or ensure it's uppercase if it's a string
      if (typeof userData.role === 'number') {
        userData.role = mapRoleIdToString(userData.role);
      } else if (typeof userData.role === 'string') {
        userData.role = userData.role.toUpperCase();
      }
      
      return userData;
    } catch (error) {
      logError('GET', url, error);
      throw error; // Re-throw to allow error handling up the chain
    }
  },

  // Request password reset
  requestPasswordReset: async (email) => {
    const url = `${AUTH_BASE}/forgot-password`;
    try {
      logRequest('POST', url, { email });
      const response = await api.post(url, { email });
      logResponse('POST', url, response);
      // The backend wraps the response in a 'data' field
      return response.data.data || response.data;
    } catch (error) {
      logError('POST', url, error);
    }
  },

  // Reset password with token
  resetPassword: async (token, password, passwordConfirm) => {
    const url = `${AUTH_BASE}/reset-password`;
    const data = {
      token,
      password: '***', // Don't log actual password
      passwordConfirm: '***',
    };
    
    try {
      logRequest('POST', url, { ...data, passwordLength: password.length });
      const response = await api.post(url, { token, password, passwordConfirm });
      logResponse('POST', url, response);
      // The backend wraps the response in a 'data' field
      return response.data.data || response.data;
    } catch (error) {
      logError('POST', url, error);
    }
  },
};

export default authService;
