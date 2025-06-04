import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Grid,
  Paper,
  Typography,
  Button,
  Divider,
  useTheme,
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
import { formatCurrency } from '../utils/formatters';

// Mock data for recent orders
const recentOrders = [
  { id: 1, orderId: 'ORD-001', customer: 'John Doe', amount: 125.99, status: 'Completed', date: '2023-06-15' },
  { id: 2, orderId: 'ORD-002', customer: 'Jane Smith', amount: 89.50, status: 'Processing', date: '2023-06-14' },
  { id: 3, orderId: 'ORD-003', customer: 'Acme Corp', amount: 245.75, status: 'Shipped', date: '2023-06-14' },
  { id: 4, orderId: 'ORD-004', customer: 'Tech Solutions', amount: 179.99, status: 'Pending', date: '2023-06-13' },
  { id: 5, orderId: 'ORD-005', customer: 'Retail Plus', amount: 320.00, status: 'Completed', date: '2023-06-12' },
];

const StatCard = ({ title, value, icon: Icon, color, onClick }) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

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
          boxShadow: theme.shadows[8],
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
            {value}
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
  const theme = useTheme();

  // Mock products data to use when API is not available
  const mockProducts = [
    { id: 1, name: 'Sample Product 1', category: { name: 'Electronics' }, stock: 42, price: 99.99 },
    { id: 2, name: 'Sample Product 2', category: { name: 'Clothing' }, stock: 15, price: 29.99 },
    { id: 3, name: 'Sample Product 3', category: { name: 'Home' }, stock: 7, price: 149.99 },
    { id: 4, name: 'Sample Product 4', category: { name: 'Electronics' }, stock: 23, price: 199.99 },
    { id: 5, name: 'Sample Product 5', category: { name: 'Books' }, stock: 31, price: 14.99 },
  ];

  // Fetch products data with fallback to mock data
  const { data: products = [], isLoading } = useQuery({
    queryKey: ['dashboard-products'],
    queryFn: async () => {
      try {
        const data = await productService.getProducts({ limit: 5 });
        return data;
      } catch (error) {
        console.warn('Using mock data due to API error:', error);
        return mockProducts;
      }
    },
    initialData: mockProducts, // Use mock data initially
  });

  // Calculate statistics
  const totalProducts = 142;
  const totalOrders = 89;
  const revenue = 12543.75;
  const growth = 12.5;

  const handleAddProduct = () => {
    navigate('/products/new');
  };

  const columns = [
    { 
      field: 'name', 
      headerName: 'Product', 
      flex: 1,
      minWidth: 200,
    },
    { 
      field: 'category', 
      headerName: 'Category', 
      width: 150,
      renderCell: (params) => (
        <Typography variant="body2" color="text.secondary">
          {params.row.category?.name || 'N/A'}
        </Typography>
      ),
    },
    { 
      field: 'stock', 
      headerName: 'Stock', 
      width: 120,
      align: 'right',
      headerAlign: 'right',
    },
    { 
      field: 'price', 
      headerName: 'Price', 
      width: 120,
      align: 'right',
      headerAlign: 'right',
      valueFormatter: (params) => formatCurrency(params.value),
    },
  ];

  const orderColumns = [
    { 
      field: 'orderId', 
      headerName: 'Order ID', 
      width: 130,
    },
    { 
      field: 'customer', 
      headerName: 'Customer', 
      flex: 1,
      minWidth: 150,
    },
    { 
      field: 'date', 
      headerName: 'Date', 
      width: 120,
    },
    { 
      field: 'amount', 
      headerName: 'Amount', 
      width: 120,
      align: 'right',
      headerAlign: 'right',
      valueFormatter: (params) => formatCurrency(params.value),
    },
    { 
      field: 'status', 
      headerName: 'Status', 
      width: 130,
      renderCell: (params) => (
        <Box
          sx={{
            backgroundColor: 
              params.value === 'Completed' ? 'success.light' :
              params.value === 'Processing' ? 'info.light' :
              params.value === 'Shipped' ? 'primary.light' : 'warning.light',
            color: 
              params.value === 'Completed' ? 'success.dark' :
              params.value === 'Processing' ? 'info.dark' :
              params.value === 'Shipped' ? 'primary.dark' : 'warning.dark',
            px: 1.5,
            py: 0.5,
            borderRadius: 4,
            fontSize: '0.75rem',
            fontWeight: 500,
            textTransform: 'capitalize',
          }}
        >
          {params.value}
        </Box>
      ),
    },
  ];

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
            value={totalProducts.toLocaleString()}
            icon={InventoryIcon}
            color="primary"
            onClick={() => navigate('/products')}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Total Orders"
            value={totalOrders.toLocaleString()}
            icon={OrdersIcon}
            color="secondary"
            onClick={() => navigate('/orders')}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Revenue"
            value={formatCurrency(revenue)}
            icon={TrendingIcon}
            color="success"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Growth"
            value={`${growth}%`}
            icon={CategoryIcon}
            color="info"
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
            columns={columns}
            pageSize={5}
            rowsPerPageOptions={[5]}
            disableSelectionOnClick
            loading={isLoading}
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
            rows={recentOrders}
            columns={orderColumns}
            pageSize={5}
            rowsPerPageOptions={[5]}
            disableSelectionOnClick
          />
        </Box>
      </Paper>
    </Box>
  );
};

export default DashboardPage;
