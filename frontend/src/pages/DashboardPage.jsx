import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Grid,
  Paper,
  Typography,
  Button,
  useMediaQuery,
} from '@mui/material';
import {
  Inventory as InventoryIcon,
  ShoppingCart as OrdersIcon,
  Category as CategoryIcon,
  TrendingUp as TrendingIcon,
  Add as AddIcon,
} from '@mui/icons-material';
import { DataGrid } from '@mui/x-data-grid';
import { useQuery } from '@tanstack/react-query';
import productService from '../services/productService';
import orderService from '../services/orderService';
import categoryService from '../services/categoryService';
import { formatCurrency } from '../utils/formatters';

const StatCard = ({ title, value, icon: Icon, color, onClick, loading = false }) => {
  const isMobile = useMediaQuery('sm');

  return (
    <Paper 
      sx={{ 
        p: 3, 
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        transition: 'transform 0.2s, box-shadow 0.2s',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.2)',
        },
        cursor: 'pointer',
      }}
      onClick={onClick}
    >
      <Box display="flex" justifyContent="space-between" alignItems="center">
        <Box>
          <Typography variant="subtitle2" color="text.secondary" gutterBottom>
            {title}
          </Typography>
          <Typography variant="h4" component="div" sx={{ fontWeight: 'bold' }}>
            {loading ? '...' : value}
          </Typography>
        </Box>
        <Box
          sx={{
            backgroundColor: `${color}.light`,
            color: `${color}.dark`,
            borderRadius: '50%',
            width: 56,
            height: 56,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Icon fontSize="large" />
        </Box>
      </Box>
      {!isMobile && (
        <Box sx={{ mt: 'auto', pt: 1 }}>
          <Typography variant="caption" color="text.secondary">
            View all
          </Typography>
        </Box>
      )}
    </Paper>
  );
};

const DashboardPage = () => {
  const navigate = useNavigate();

  // Fetch dashboard statistics
  const { data: productsData, isLoading: productsLoading } = useQuery({
    queryKey: ['dashboard-products'],
    queryFn: async () => {
      try {
        return await productService.getProducts({ limit: 5 });
      } catch (error) {
        console.error('Error fetching products:', error);
        return [];
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });

  const { data: ordersData, isLoading: ordersLoading } = useQuery({
    queryKey: ['dashboard-orders'],
    queryFn: async () => {
      try {
        return await orderService.getOrders({ limit: 5 });
      } catch (error) {
        console.error('Error fetching orders:', error);
        return [];
      }
    },
    staleTime: 5 * 60 * 1000,
  });

  const { data: categoriesData, isLoading: categoriesLoading } = useQuery({
    queryKey: ['dashboard-categories'],
    queryFn: async () => {
      try {
        return await categoryService.getCategories({ limit: 100 });
      } catch (error) {
        console.error('Error fetching categories:', error);
        return [];
      }
    },
    staleTime: 5 * 60 * 1000,
  });

  const { data: statsData, isLoading: statsLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: async () => {
      try {
        // Fetch basic statistics
        const [products, orders] = await Promise.all([
          productService.getProducts({ limit: 1000 }),
          orderService.getOrders({ limit: 1000 })
        ]);

        // Calculate revenue from orders
        const revenue = Array.isArray(orders) 
          ? orders.reduce((sum, order) => sum + (order.total || order.amount || 0), 0)
          : 0;

        // Calculate growth (mock for now - would need historical data)
        const growth = 12.5;

        return {
          totalProducts: Array.isArray(products) ? products.length : 0,
          totalOrders: Array.isArray(orders) ? orders.length : 0,
          revenue,
          growth
        };
      } catch (error) {
        console.error('Error fetching stats:', error);
        return {
          totalProducts: 0,
          totalOrders: 0,
          revenue: 0,
          growth: 0
        };
      }
    },
    staleTime: 5 * 60 * 1000,
  });

  const handleAddProduct = () => {
    navigate('/products/new');
  };

  const productColumns = [
    { 
      field: 'name', 
      headerName: 'Product Name', 
      flex: 1,
      minWidth: 200,
    },
    { 
      field: 'category', 
      headerName: 'Category', 
      flex: 1,
      valueGetter: (params) => params.row.category?.name || 'N/A',
    },
    { 
      field: 'stock', 
      headerName: 'Stock', 
      width: 100,
      type: 'number',
    },
    { 
      field: 'price', 
      headerName: 'Price', 
      width: 120,
      valueFormatter: (params) => formatCurrency(params.value || 0),
    },
  ];

  const orderColumns = [
    { 
      field: 'orderId', 
      headerName: 'Order ID', 
      width: 120,
      valueGetter: (params) => params.row.orderId || params.row._id || params.row.id,
    },
    { 
      field: 'customer', 
      headerName: 'Customer', 
      flex: 1,
      valueGetter: (params) => params.row.customer?.name || params.row.customerName || 'N/A',
    },
    { 
      field: 'amount', 
      headerName: 'Amount', 
      width: 120,
      valueFormatter: (params) => formatCurrency(params.row.total || params.row.amount || 0),
    },
    { 
      field: 'status', 
      headerName: 'Status', 
      width: 130,
      renderCell: (params) => (
        <Box
          sx={{
            backgroundColor: 
              params.value === 'Completed' || params.value === 'completed' ? 'success.light' :
              params.value === 'Processing' || params.value === 'processing' ? 'info.light' :
              params.value === 'Shipped' || params.value === 'shipped' ? 'primary.light' : 'warning.light',
            color: 
              params.value === 'Completed' || params.value === 'completed' ? 'success.dark' :
              params.value === 'Processing' || params.value === 'processing' ? 'info.dark' :
              params.value === 'Shipped' || params.value === 'shipped' ? 'primary.dark' : 'warning.dark',
            px: 1.5,
            py: 0.5,
            borderRadius: 4,
            fontSize: '0.75rem',
            fontWeight: 500,
            textTransform: 'capitalize',
          }}
        >
          {params.value || 'pending'}
        </Box>
      ),
    },
  ];

  const stats = statsData || { totalProducts: 0, totalOrders: 0, revenue: 0, growth: 0 };
  const products = Array.isArray(productsData) ? productsData : [];
  const orders = Array.isArray(ordersData) ? ordersData : [];

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
          Dashboard
        </Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={handleAddProduct}
        >
          Add Product
        </Button>
      </Box>

      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Total Products"
            value={stats.totalProducts.toLocaleString()}
            icon={InventoryIcon}
            color="primary"
            onClick={() => navigate('/products')}
            loading={statsLoading}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Total Orders"
            value={stats.totalOrders.toLocaleString()}
            icon={OrdersIcon}
            color="secondary"
            onClick={() => navigate('/orders')}
            loading={statsLoading}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Revenue"
            value={formatCurrency(stats.revenue)}
            icon={TrendingIcon}
            color="success"
            loading={statsLoading}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Categories"
            value={Array.isArray(categoriesData) ? categoriesData.length : 0}
            icon={CategoryIcon}
            color="info"
            onClick={() => navigate('/categories')}
            loading={categoriesLoading}
          />
        </Grid>
      </Grid>

      {/* Recent Products */}
      <Paper sx={{ p: 3, mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" component="h2" sx={{ fontWeight: 'bold' }}>
            Recent Products
          </Typography>
          <Button
            variant="text"
            color="primary"
            size="small"
            onClick={() => navigate('/products')}
          >
            View All
          </Button>
        </Box>
        <Box sx={{ height: 400, width: '100%' }}>
          <DataGrid
            rows={products}
            columns={productColumns}
            pageSize={5}
            rowsPerPageOptions={[5]}
            disableSelectionOnClick
            loading={productsLoading}
            getRowId={(row) => row._id || row.id}
          />
        </Box>
      </Paper>

      {/* Recent Orders */}
      <Paper sx={{ p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" component="h2" sx={{ fontWeight: 'bold' }}>
            Recent Orders
          </Typography>
          <Button
            variant="text"
            color="primary"
            size="small"
            onClick={() => navigate('/orders')}
          >
            View All
          </Button>
        </Box>
        <Box sx={{ height: 400, width: '100%' }}>
          <DataGrid
            rows={orders}
            columns={orderColumns}
            pageSize={5}
            rowsPerPageOptions={[5]}
            disableSelectionOnClick
            loading={ordersLoading}
            getRowId={(row) => row._id || row.id}
          />
        </Box>
      </Paper>
    </Box>
  );
};

export default DashboardPage;
