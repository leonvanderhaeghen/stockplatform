import { Navigate, Outlet, useLocation } from 'react-router-dom';
import { useContext, useMemo, useEffect, useState } from 'react';
import { hasRole } from '../../utils/auth';
import { AuthContext } from '../../App';
import { CircularProgress, Box } from '@mui/material';

/**
 * ProtectedRoute component for handling authentication and authorization
 * @param {string[]} [props.allowedRoles=[]] - Array of allowed roles
 * @param {string} [props.redirectTo] - Path to redirect to when user is not authorized/authenticated
 * @param {Object} [props.state] - State to pass to the redirect
 * @returns {JSX.Element} Rendered component
 */
const ProtectedRoute = ({
  allowedRoles = [],
  redirectTo = '/login',
  state,
  ...rest
}) => {
  const location = useLocation();
  const { isAuthenticated, isLoading, user } = useContext(AuthContext);
  const [initialCheckComplete, setInitialCheckComplete] = useState(false);
  
  // Use effect to set initial check complete after first render
  useEffect(() => {
    if (!isLoading) {
      setInitialCheckComplete(true);
    }
  }, [isLoading]);
  
  // Memoize the auth check to prevent unnecessary re-renders
  const authCheck = useMemo(() => {
    console.log('Auth check - isAuthenticated:', isAuthenticated);
    console.log('Auth check - allowedRoles:', allowedRoles);
    console.log('Auth check - current user:', user);
    
    // If still loading, don't make any decisions yet
    if (isLoading) {
      console.log('Auth check - still loading...');
      return { isAuth: false, isLoading: true };
    }
    
    // If no roles required, just check authentication
    if (allowedRoles.length === 0) {
      console.log('Auth check - no roles required, checking authentication only');
      return { 
        isAuth: isAuthenticated, 
        isLoading: false,
        hasRequiredRole: true
      };
    }
    
    // Check if user has any of the required roles
    const userHasRole = hasRole(allowedRoles);
    console.log('Auth check - user has required role?', userHasRole);
    
    return { 
      isAuth: isAuthenticated, 
      isLoading: false,
      hasRequiredRole: userHasRole
    };
  }, [isLoading, isAuthenticated, allowedRoles, user]);

  // Show loading state while checking auth status
  if (!initialCheckComplete || authCheck.isLoading) {
    return (
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        backgroundColor: 'background.default'
      }}>
        <CircularProgress />
      </Box>
    );
  }

  // If not authenticated, redirect to login with the current location
  if (!authCheck.isAuth) {
    console.log('Not authenticated, redirecting to login');
    // Only redirect if we're not already on the login page to prevent loops
    if (location.pathname !== redirectTo) {
      const redirectState = { 
        from: location.pathname === '/login' ? '/' : location,
        error: 'Please sign in to access this page.',
        ...state 
      };
      console.log('Redirecting to login with state:', redirectState);
      
      return (
        <Navigate 
          to={redirectTo} 
          state={redirectState} 
          replace 
        />
      );
    }
    return null; // Prevent rendering anything if we're already redirecting to login
  }

  // If no specific roles required, allow access
  if (allowedRoles.length === 0) {
    console.log('No roles required, allowing access');
    return <Outlet {...rest} />;
  }

  // If user has required role, render the child routes
  if (authCheck.hasRequiredRole) {
    console.log('User has required role, allowing access');
    return <Outlet {...rest} />;
  }
  
  console.log('User does not have required role');

  // If user doesn't have required role, redirect to unauthorized or specified route
  return (
    <Navigate 
      to="/unauthorized" 
      state={{ 
        from: location,
        error: 'You do not have permission to access this page.',
        ...state 
      }} 
      replace 
    />
  );

};

export default ProtectedRoute;
