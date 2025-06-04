import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import authService from '../../services/authService';
import { 
  getAuthToken, 
  storeAuthData as storeAuthDataUtil, 
  clearAuthData as clearAuthDataUtil,
  getCurrentUser as getCurrentUserUtil,
  isAuthenticated as isAuthenticatedUtil,
  hasRole as hasRoleUtil
} from '../auth';

// Map numeric role IDs to role strings
const mapRoleIdToString = (roleId) => {
  console.log('Mapping role ID to string:', roleId);
  
  // If already a string, ensure it's uppercase
  if (typeof roleId === 'string') {
    return roleId.toUpperCase();
  }
  
  const roleMap = {
    0: 'ADMIN',
    1: 'CUSTOMER',
    2: 'MANAGER',
    3: 'USER',
    // Add any other mappings as needed
    admin: 'ADMIN',
    customer: 'CUSTOMER',
    manager: 'MANAGER',
    user: 'USER'
  };
  
  const role = roleMap[roleId] || 'USER';
  console.log('Mapped role:', role);
  return role;
};

// Helper function to normalize user data
const normalizeUserData = (userData) => {
  if (!userData) return null;
  
  // Create a copy to avoid mutating the original
  const normalized = { ...userData };
  
  // Ensure role is properly mapped
  if (normalized.role !== undefined) {
    normalized.role = mapRoleIdToString(normalized.role);
  }
  
  return normalized;
};

const useAuth = () => {
  const [isAuth, setIsAuth] = useState(isAuthenticatedUtil());
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState(getCurrentUserUtil());
  const [authError, setAuthError] = useState(null);
  const navigate = useNavigate();

  // Store auth data in localStorage
  const storeAuthData = useCallback((token, userData) => {
    // Normalize user data before storing
    const normalizedUser = normalizeUserData(userData);
    console.log('Storing auth data:', { token, user: normalizedUser });
    
    storeAuthDataUtil(token, normalizedUser);
    setUser(normalizedUser);
    setIsAuth(true);
  }, []);

  // Clear auth data from localStorage
  const clearAuthData = useCallback(() => {
    clearAuthDataUtil();
    setUser(null);
    setIsAuth(false);
  }, []);

  // Validate token with backend
  const validateToken = useCallback(async () => {
    try {
      const userData = await authService.getCurrentUser();
      console.log('User data from validateToken:', userData);
      
      // Normalize the user data
      const normalizedUser = normalizeUserData(userData);
      
      storeAuthData(getAuthToken(), normalizedUser);
      return true;
    } catch (error) {
      console.error('Token validation failed:', error);
      clearAuthData();
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [clearAuthData, storeAuthData]);

  // Initialize auth state from localStorage
  useEffect(() => {
    let isMounted = true;
    
    const initializeAuth = async () => {
      const token = getAuthToken();
      const currentUser = getCurrentUserUtil();
      
      if (token && currentUser) {
        try {
          // Set initial state optimistically
          if (isMounted) {
            setUser(currentUser);
            setIsAuth(true);
          }
          
          // Validate token with backend in the background
          const isValid = await validateToken();
          
          if (!isValid && isMounted) {
            clearAuthData();
          }
        } catch (error) {
          console.error('Failed to initialize auth:', error);
          if (isMounted) {
            clearAuthData();
          }
        }
      }
      
      if (isMounted) {
        setIsLoading(false);
      }
    };

    initializeAuth();
    
    return () => {
      isMounted = false;
    };
  }, [validateToken, clearAuthData]);

  // Login function - navigation is handled by the component, not here
  const login = useCallback(async (email, password) => {
    console.log('useAuth: login called with email:', email);
    setAuthError(null);
    setIsLoading(true);
    
    try {
      console.log('useAuth: Calling authService.login...');
      const response = await authService.login(email, password);
      console.log('useAuth: authService.login response:', response);
      
      if (!response) {
        throw new Error('No response from authentication service');
      }
      
      const { token, user: userData } = response;
      
      if (!token) {
        throw new Error('No authentication token received');
      }
      
      if (!userData) {
        throw new Error('No user data received');
      }
      
      console.log('useAuth: Storing auth data...');
      
      // Map the numeric role to a string role if needed
      const mappedUserData = {
        ...userData,
        role: mapRoleIdToString(userData.role)
      };
      
      // Store the token and user data
      storeAuthData(token, mappedUserData);
      console.log('useAuth: Auth data stored, navigating to dashboard...');
      
      // Return success without navigating (navigation will be handled by the form)
      return { 
        success: true, 
        user: userData,
        token: token
      };
      
    } catch (error) {
      console.error('useAuth: Login error:', error);
      
      let errorMessage = 'Login failed. Please try again.';
      
      if (error.response) {
        // The request was made and the server responded with a status code
        // that falls out of the range of 2xx
        console.error('Response data:', error.response.data);
        console.error('Response status:', error.response.status);
        console.error('Response headers:', error.response.headers);
        
        errorMessage = error.response.data?.message || 
                     error.response.data?.error?.message || 
                     error.response.statusText || 
                     'Authentication failed';
                     
      } else if (error.request) {
        // The request was made but no response was received
        console.error('No response received:', error.request);
        errorMessage = 'No response from server. Please check your connection.';
      } else {
        // Something happened in setting up the request that triggered an Error
        console.error('Error setting up request:', error.message);
        errorMessage = error.message || 'An unexpected error occurred';
      }
      
      setAuthError(errorMessage);
      return { 
        success: false, 
        error: errorMessage,
        details: error.response?.data
      };
      
    } finally {
      console.log('useAuth: Login attempt completed');
      setIsLoading(false);
    }
  }, [storeAuthData]);

  // Logout function
  const logout = useCallback(async () => {
    try {
      await authService.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      clearAuthData();
      navigate('/login');
    }
  }, [navigate, clearAuthData]);

  // Check if user has specific role
  const hasRole = useCallback((roles) => {
    if (!user?.role) return false;
    return hasRoleUtil(roles, user.role);
  }, [user]);

  return {
    isAuthenticated: isAuth,
    isLoading,
    user,
    error: authError,
    login,
    logout,
    storeAuthData,
    clearAuthData,
    hasRole,
  };
};

export default useAuth;
