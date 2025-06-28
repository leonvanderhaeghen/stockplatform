import axios from 'axios';
import { getAuthToken } from '../utils/auth';

// Base URL for API requests - using a relative path that will be handled by the nginx proxy
const API_BASE_URL = '/api/v1';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: false, // Disable credentials to avoid CORS preflight
  timeout: 30000, // 30 second timeout
});

// Request interceptor to add auth token to requests
api.interceptors.request.use(
  (config) => {
    const token = getAuthToken() || localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    console.log('API Request:', config.method?.toUpperCase(), config.url);
    return config;
  },
  (error) => {
    console.error('Request interceptor error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor to handle common errors
api.interceptors.response.use(
  (response) => {
    console.log('API Response:', response.config.url, response.status);
    return response;
  },
  (error) => {
    const { response } = error;
    
    console.error('API Error:', error.config?.url, response?.status, error.message);
    
    // Handle common error cases
    if (response) {
      // Handle specific HTTP status codes
      switch (response.status) {
        case 401: // Unauthorized
          console.warn('Unauthorized access - clearing auth data');
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          // Don't redirect automatically to avoid infinite loops
          break;
        case 403: // Forbidden
          console.error('Access forbidden');
          break;
        case 404: // Not Found
          console.error('Resource not found');
          break;
        case 500: // Internal Server Error
          console.error('Server error occurred');
          break;
        default:
          console.error('An error occurred');
      }
    } else if (error.request) {
      // The request was made but no response was received
      console.error('No response from server. Please check your connection.');
    } else {
      // Something happened in setting up the request
      console.error('Request error:', error.message);
    }
    
    return Promise.reject(error);
  }
);

export default api;
