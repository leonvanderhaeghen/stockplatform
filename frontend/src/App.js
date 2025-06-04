import React, { createContext, useContext, useMemo } from 'react';
import { 
  Routes,
  Route,
  Navigate,
  useLocation
} from 'react-router-dom';
import { ThemeProvider, StyledEngineProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import CircularProgress from '@mui/material/CircularProgress';
import Typography from '@mui/material/Typography';
import { SnackbarProvider } from 'notistack';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { ThemeProvider as EmotionThemeProvider } from '@emotion/react';

// Layout
import MainLayout from './components/layout/MainLayout';

// Pages
import DashboardPage from './pages/DashboardPage';
import LoginPage from './pages/LoginPage';
import UnauthorizedPage from './pages/UnauthorizedPage';
import NotFoundPage from './pages/NotFoundPage';

// Product Pages
import ProductsPage from './components/products/ProductsCRUD';
import CategoriesPage from './pages/CategoriesPage';

// Inventory Pages
import InventoryLevelsPage from './pages/inventory/InventoryLevelsPage';
import StockTransfersPage from './pages/inventory/StockTransfersPage';
import StockAdjustmentsPage from './pages/inventory/StockAdjustmentsPage';

// Order Pages
import OrdersPage from './components/orders/OrdersCRUD';
import CreateOrderPage from './pages/CreateOrderPage';
import ReturnsPage from './pages/ReturnsPage';

// User & Customer Pages
import UsersPage from './components/users/UsersCRUD';
import CustomersPage from './pages/CustomersPage';

// Report Pages
import ReportsPage from './pages/ReportsPage';
import SalesReportsPage from './pages/reports/SalesReportsPage';
import InventoryReportsPage from './pages/reports/InventoryReportsPage';
import CustomerReportsPage from './pages/reports/CustomerReportsPage';

// Components
import ProtectedRoute from './components/common/ProtectedRoute';

// Theme
import theme from './theme';

// Hooks
import useAuth from './utils/hooks/useAuth';

// Create auth context
export const AuthContext = createContext();

const AuthProvider = ({ children }) => {
  const auth = useAuth();
  
  // Memoize the context value to prevent unnecessary re-renders
  const contextValue = useMemo(
    () => {
      console.log('Auth context updated:', {
        isAuthenticated: auth.isAuthenticated,
        isLoading: auth.isLoading,
        user: auth.user,
        hasError: !!auth.error
      });
      
      if (auth.user) {
        logUserInfo(auth.user);
      }
      
      return {
        isAuthenticated: auth.isAuthenticated,
        isLoading: auth.isLoading,
        user: auth.user,
        login: auth.login,
        logout: auth.logout,
        error: auth.error,
      };
    },
    [auth.isAuthenticated, auth.isLoading, auth.user, auth.login, auth.logout, auth.error]
  );

  if (auth.isLoading) {
    return (
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        flexDirection: 'column',
        gap: 2
      }}>
        <CircularProgress />
        <Typography>Loading application...</Typography>
      </Box>
    );
  }

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
};

// Public route component (for login/register pages)
const PublicRoute = ({ children }) => {
  const { isAuthenticated, isLoading } = useContext(AuthContext);
  const location = useLocation();

  // Show loading state while checking auth status
  if (isLoading) {
    return (
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh' 
      }}>
        <CircularProgress />
      </Box>
    );
  }

  // If authenticated, redirect to home or intended page
  if (isAuthenticated) {
    const from = (location.state?.from?.pathname || '/').startsWith('/login') 
      ? '/' 
      : location.state?.from?.pathname || '/';
    return <Navigate to={from} state={{ from: location }} replace />;
  }

  // If not authenticated, render the public route
  return children;
};

// Wrap the app with all the providers
const AppWrapper = () => {
  return (
    <StyledEngineProvider injectFirst>
      <ThemeProvider theme={theme}>
        <EmotionThemeProvider theme={theme}>
          <CssBaseline />
          <LocalizationProvider dateAdapter={AdapterDateFns}>
            <SnackbarProvider
              maxSnack={3}
              anchorOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              autoHideDuration={3000}
            >
              <AuthProvider>
                <AppRoutes />
              </AuthProvider>
            </SnackbarProvider>
          </LocalizationProvider>
        </EmotionThemeProvider>
      </ThemeProvider>
    </StyledEngineProvider>
  );
};

// Define role-based access control
export const ROLES = {
  ADMIN: 'ADMIN',
  MANAGER: 'MANAGER',
  USER: 'USER',
  CUSTOMER: 'CUSTOMER',
};

// Debug function to log user info
const logUserInfo = (user) => {
  console.log('Current user info:', {
    id: user?.id,
    email: user?.email,
    role: user?.role,
    isAuthenticated: !!user
  });
};

// Define the main routes
const AppRoutes = () => {
  // Get the current location for navigation
  const location = useLocation();
  
  // Only use the location state if it's not the login page to prevent loops
  const getLocationState = () => {
    if (location.pathname === '/login') {
      return { from: { pathname: '/' } };
    }
    return { from: location };
  };
  
  return (
    <Routes>
      {/* Public routes */}
      <Route 
        path="/login" 
        element={
          <PublicRoute>
            <LoginPage />
          </PublicRoute>
        } 
      />
      
      {/* Protected routes */}
      <Route element={
        <ProtectedRoute 
          allowedRoles={[]} // Allow all authenticated users
          redirectTo="/login"
          state={getLocationState()}
        />
      }>
        <Route element={<MainLayout />}>
          <Route index element={<DashboardPage />} />
          <Route path="dashboard" element={<DashboardPage />} />
          
          {/* Admin only routes */}
          <Route element={
            <ProtectedRoute 
              allowedRoles={[ROLES.ADMIN]} 
              redirectTo="/unauthorized"
            />
          }>
            <Route path="users" element={<UsersPage />} />
          </Route>
          
          {/* Manager and Admin routes */}
          <Route element={
            <ProtectedRoute 
              allowedRoles={[ROLES.ADMIN, ROLES.MANAGER]} 
              redirectTo="/unauthorized"
            />
          }>
            {/* Product Routes */}
            <Route path="products">
              <Route index element={<ProductsPage />} />
              <Route path="categories" element={<CategoriesPage />} />
            </Route>

            {/* Inventory Routes */}
            <Route path="inventory">
              <Route index element={<InventoryLevelsPage />} />
              <Route path="levels" element={<InventoryLevelsPage />} />
              <Route path="transfers" element={<StockTransfersPage />} />
              <Route path="adjustments" element={<StockAdjustmentsPage />} />
            </Route>

            {/* Customer Routes */}
            <Route path="customers" element={<CustomersPage />} />

            {/* Report Routes */}
            <Route path="reports">
              <Route index element={<ReportsPage />} />
              <Route path="sales" element={<SalesReportsPage />} />
              <Route path="inventory" element={<InventoryReportsPage />} />
              <Route path="customers" element={<CustomerReportsPage />} />
            </Route>
          </Route>
          
          {/* All authenticated users */}
          <Route path="orders">
            <Route index element={<OrdersPage />} />
            <Route path="new" element={<CreateOrderPage />} />
          </Route>
          <Route path="returns" element={<ReturnsPage />} />
        </Route>
      </Route>
      
      {/* Error pages */}
      <Route path="/unauthorized" element={<UnauthorizedPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
};

// Main App component
export default AppWrapper;
