import { useState, useEffect, useCallback, createContext, useContext } from 'react';
import { jwtDecode } from 'jwt-decode';
import authService from '../services/authService';

// Create authentication context
const AuthContext = createContext(null);

// Custom hook to use the auth context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// Auth provider component
export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [token, setToken] = useState(localStorage.getItem('token'));

  // Check if token is expired
  const isTokenExpired = useCallback((token) => {
    if (!token) return true;
    
    console.log('Checking token expiry for token:', token);
    console.log('Token type:', typeof token);
    console.log('Token length:', token?.length);
    
    try {
      const decoded = jwtDecode(token);
      const now = Date.now() / 1000;
      return decoded.exp < now;
    } catch (error) {
      console.error('Error decoding token in isTokenExpired:', error);
      console.error('Problematic token:', token);
      return true;
    }
  }, []);

  // Extract user from token
  const getUserFromToken = useCallback((token) => {
    if (!token || isTokenExpired(token)) return null;
    
    try {
      const decoded = jwtDecode(token);
      return {
        id: decoded.sub || decoded.userId,
        email: decoded.email,
        role: decoded.role,
        firstName: decoded.firstName,
        lastName: decoded.lastName,
      };
    } catch (error) {
      console.error('Error decoding token:', error);
      return null;
    }
  }, [isTokenExpired]);

  // Login function
  const login = useCallback(async (email, password) => {
    setIsLoading(true);
    try {
      const response = await authService.login(email, password);
      console.log('Login response:', response);
      const { token: newToken, user: userData } = response.data;
      
      console.log('Extracted token:', newToken);
      console.log('Token length:', newToken?.length);
      console.log('Token parts:', newToken?.split('.').length);
      
      setToken(newToken);
      setUser(userData);
      localStorage.setItem('token', newToken);
      
      console.log('Stored in localStorage:', localStorage.getItem('token'));
      
      return { success: true, user: userData };
    } catch (error) {
      console.error('Login error:', error);
      return { 
        success: false, 
        error: error.response?.data?.message || 'Login failed' 
      };
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Register function
  const register = useCallback(async (userData) => {
    setIsLoading(true);
    try {
      const response = await authService.register(userData);
      return { success: true, data: response };
    } catch (error) {
      console.error('Registration error:', error);
      return { 
        success: false, 
        error: error.response?.data?.message || 'Registration failed' 
      };
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Logout function
  const logout = useCallback(() => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
  }, []);

  // Refresh token function - disabled since backend doesn't support it
  const refreshToken = useCallback(async () => {
    console.warn('Token refresh not supported by backend - logging out');
    logout();
    return false;
  }, [logout]);

  // Check authentication on mount and token change
  useEffect(() => {
    const initializeAuth = async () => {
      setIsLoading(true);
      
      if (!token) {
        setIsLoading(false);
        return;
      }

      if (isTokenExpired(token)) {
        // Try to refresh token
        const refreshed = await refreshToken();
        if (!refreshed) {
          logout();
        }
      } else {
        // Extract user from valid token
        const userData = getUserFromToken(token);
        setUser(userData);
      }
      
      setIsLoading(false);
    };

    initializeAuth();
  }, [token, isTokenExpired, getUserFromToken, refreshToken, logout]);

  // Set up automatic token refresh
  useEffect(() => {
    if (!token || isTokenExpired(token)) return;

    const decoded = jwtDecode(token);
    const expirationTime = decoded.exp * 1000;
    const currentTime = Date.now();
    const timeUntilExpiration = expirationTime - currentTime;
    
    // Refresh token 5 minutes before it expires
    const refreshTime = Math.max(timeUntilExpiration - (5 * 60 * 1000), 0);

    const timeoutId = setTimeout(() => {
      refreshToken();
    }, refreshTime);

    return () => clearTimeout(timeoutId);
  }, [token, isTokenExpired, refreshToken]);

  const value = {
    user,
    token,
    isLoading,
    isAuthenticated: !!user && !!token && !isTokenExpired(token),
    login,
    register,
    logout,
    refreshToken,
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default useAuth;
