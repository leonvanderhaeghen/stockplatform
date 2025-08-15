import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { CssBaseline, Box } from '@mui/material';
import { SnackbarProvider } from 'notistack';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';

// Core Layout & Navigation
import AppLayout from './components/layout/AppLayout';
import AuthGuard from './components/auth/AuthGuard';
import RoleGuard from './components/auth/RoleGuard';

// Authentication Context
import { AuthProvider } from './hooks/useAuth';

// Authentication Pages
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';

// Customer Pages
import DashboardPage from './pages/dashboard/DashboardPage';
import ProductsPage from './pages/products/ProductsPage';
import OrdersPage from './pages/orders/OrdersPage';
import ProfilePage from './pages/profile/ProfilePage';

// Staff/Admin Pages
import InventoryPage from './pages/inventory/InventoryPage';
import SuppliersPage from './pages/suppliers/SuppliersPage';
import StoresPage from './pages/stores/StoresPage';
import AdminPage from './pages/admin/AdminPage';
import POSPage from './pages/pos/POSPage';
import CategoriesPage from './pages/categories/CategoriesPage';

// Error Pages
import NotFoundPage from './pages/error/NotFoundPage';
import UnauthorizedPage from './pages/error/UnauthorizedPage';

// Create React Query client with optimized configuration
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5 minutes
      gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
    },
    mutations: {
      retry: 1,
    },
  },
});

// Create Material-UI theme with StockPlatform branding
const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#1976d2',
      light: '#42a5f5',
      dark: '#1565c0',
    },
    secondary: {
      main: '#dc004e',
      light: '#ff5983',
      dark: '#9a0036',
    },
    background: {
      default: '#f5f5f5',
      paper: '#ffffff',
    },
    success: {
      main: '#2e7d32',
    },
    warning: {
      main: '#ed6c02',
    },
    error: {
      main: '#d32f2f',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 600,
    },
    h5: {
      fontWeight: 600,
    },
    h6: {
      fontWeight: 600,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 8,
        },
      },
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <SnackbarProvider 
            maxSnack={3} 
            anchorOrigin={{
              vertical: 'top',
              horizontal: 'right',
            }}
            autoHideDuration={5000}
          >
            <AuthProvider>
              <Router>
                <Box sx={{ display: 'flex', minHeight: '100vh' }}>
                  <Routes>
                  {/* Public Routes */}
                  <Route path="/login" element={<LoginPage />} />
                  <Route path="/register" element={<RegisterPage />} />
                  <Route path="/unauthorized" element={<UnauthorizedPage />} />
                  
                  {/* Protected Routes */}
                  <Route 
                    path="/*" 
                    element={
                      <AuthGuard>
                        <AppLayout>
                          <Routes>
                            {/* Dashboard - All authenticated users */}
                            <Route path="/" element={<Navigate to="/dashboard" replace />} />
                            <Route path="/dashboard" element={<DashboardPage />} />
                            
                            {/* Customer Routes */}
                            <Route path="/products" element={<ProductsPage />} />
                            <Route path="/my-orders" element={<OrdersPage userView={true} />} />
                            <Route path="/profile" element={<ProfilePage />} />
                            
                            {/* Staff/Admin Routes */}
                            <Route 
                              path="/inventory" 
                              element={
                                <RoleGuard allowedRoles={['STAFF', 'ADMIN']}>
                                  <InventoryPage />
                                </RoleGuard>
                              } 
                            />
                            <Route 
                              path="/orders" 
                              element={
                                <RoleGuard allowedRoles={['STAFF', 'ADMIN']}>
                                  <OrdersPage userView={false} />
                                </RoleGuard>
                              } 
                            />
                            <Route 
                              path="/pos" 
                              element={
                                <RoleGuard allowedRoles={['STAFF', 'ADMIN']}>
                                  <POSPage />
                                </RoleGuard>
                              } 
                            />
                            <Route 
                              path="/suppliers" 
                              element={
                                <RoleGuard allowedRoles={['STAFF', 'ADMIN']}>
                                  <SuppliersPage />
                                </RoleGuard>
                              } 
                            />
                            <Route 
                              path="/stores" 
                              element={
                                <RoleGuard allowedRoles={['STAFF', 'ADMIN']}>
                                  <StoresPage />
                                </RoleGuard>
                              } 
                            />
                            <Route 
                              path="/categories" 
                              element={
                                <RoleGuard allowedRoles={['STAFF', 'ADMIN']}>
                                  <CategoriesPage />
                                </RoleGuard>
                              } 
                            />
                            
                            {/* Admin Only Routes */}
                            <Route 
                              path="/admin" 
                              element={
                                <RoleGuard allowedRoles={['ADMIN']}>
                                  <AdminPage />
                                </RoleGuard>
                              } 
                            />
                            
                            {/* 404 Page */}
                            <Route path="*" element={<NotFoundPage />} />
                          </Routes>
                        </AppLayout>
                      </AuthGuard>
                    } 
                  />
                </Routes>
              </Box>
            </Router>
          </AuthProvider>
        </SnackbarProvider>
        </LocalizationProvider>
        
        {/* React Query Devtools - only in development */}
        {process.env.NODE_ENV === 'development' && (
          <ReactQueryDevtools initialIsOpen={false} />
        )}
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
