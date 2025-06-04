/**
 * Authentication utility functions
 */

/**
 * Get the authentication token from localStorage
 * @returns {string|null} The authentication token or null if not found
 */
export const getAuthToken = () => {
  return localStorage.getItem('token');
};

/**
 * Check if the user is authenticated
 * @returns {boolean} True if the user is authenticated, false otherwise
 */
export const isAuthenticated = () => {
  return !!getAuthToken();
};

/**
 * Get the current user from localStorage
 * @returns {Object|null} The current user object or null if not found
 */
export const getCurrentUser = () => {
  const user = localStorage.getItem('user');
  return user ? JSON.parse(user) : null;
};

/**
 * Check if the current user has the required role
 * @param {string|string[]} requiredRoles - The required role or array of roles
 * @returns {boolean} True if the user has the required role, false otherwise
 */
export const hasRole = (requiredRoles) => {
  const user = getCurrentUser();
  console.log('Current user:', user);
  console.log('Required roles:', requiredRoles);
  
  if (!user || !user.role) {
    console.log('No user or role found');
    return false;
  }
  
  // If no roles required, allow access
  if (!requiredRoles || requiredRoles.length === 0) {
    console.log('No roles required, allowing access');
    return true;
  }
  
  // Convert single role to array for consistent handling
  const rolesToCheck = Array.isArray(requiredRoles) ? requiredRoles : [requiredRoles];
  
  // Check if user has any of the required roles
  const hasRole = rolesToCheck.some(role => {
    // Case-insensitive comparison
    const hasRole = user.role.toUpperCase() === role.toUpperCase();
    if (hasRole) {
      console.log(`User has required role: ${role}`);
    }
    return hasRole;
  });
  
  if (!hasRole) {
    console.log(`User role ${user.role} not in required roles: ${rolesToCheck.join(', ')}`);
  }
  
  return hasRole;
};

/**
 * Clear authentication data from localStorage
 */
export const clearAuthData = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
};

/**
 * Store authentication data in localStorage
 * @param {string} token - The authentication token
 * @param {Object} user - The user data to store
 */
export const storeAuthData = (token, user) => {
  localStorage.setItem('token', token);
  localStorage.setItem('user', JSON.stringify(user));
};
