import React, { useState, useCallback } from 'react';
import {
  Container,
  Typography,
  Box,
  Grid,
  Card,
  CardContent,
  CardActions,
  Button,
  IconButton,
  TextField,
  InputAdornment,
  Chip,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  CircularProgress,
  Fab,
  Pagination,
  FormControl,
  InputLabel,
  Select,
  Drawer,
  Divider,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Stepper,
  Step,
  StepLabel,
  Collapse
} from '@mui/material';
import {
  ShoppingCart as ShoppingCartIcon,
  Search,
  FilterList,
  Add,
  Visibility,
  Edit,
  Cancel,
  LocalShipping,
  CheckCircle,
  AccessTime,
  ExpandMore,
  ExpandLess,
  Receipt,
  Refresh,
  MoreVert,
  Person,
  Store
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import orderService from '../../services/orderService';
import { useSnackbar } from 'notistack';
import { format } from 'date-fns';

const OrdersPage = ({ userView = false }) => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();
  
  const [searchQuery, setSearchQuery] = useState('');
  const [filters, setFilters] = useState({});
  const [page, setPage] = useState(1);
  const [filterDrawerOpen, setFilterDrawerOpen] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [orderDetailOpen, setOrderDetailOpen] = useState(false);
  const [expandedOrders, setExpandedOrders] = useState(new Set());

  const isStaff = user?.role === 'STAFF' || user?.role === 'ADMIN';
  const itemsPerPage = 10;

  // Query for orders
  const { data: ordersData, isLoading: ordersLoading, error: ordersError } = useQuery({
    queryKey: ['orders', searchQuery, filters, page, userView],
    queryFn: () => {
      if (userView) {
        return orderService.getUserOrders({
          search: searchQuery,
          ...filters,
          page,
          limit: itemsPerPage
        });
      } else {
        return orderService.listOrders({
          search: searchQuery,
          ...filters,
          page,
          limit: itemsPerPage
        });
      }
    },
    staleTime: 30000,
  });

  // Update order status mutation
  const updateStatusMutation = useMutation({
    mutationFn: ({ orderId, status }) => orderService.updateOrderStatus(orderId, status),
    onSuccess: () => {
      queryClient.invalidateQueries(['orders']);
      enqueueSnackbar('Order status updated successfully', { variant: 'success' });
      setOrderDetailOpen(false);
    },
    onError: (error) => {
      enqueueSnackbar(`Failed to update order status: ${error.message}`, { variant: 'error' });
    },
  });

  const handleSearch = useCallback((query) => {
    setSearchQuery(query);
    setPage(1);
  }, []);

  const handleFiltersChange = useCallback((newFilters) => {
    setFilters(newFilters);
    setPage(1);
  }, []);

  const handleOrderView = (order) => {
    setSelectedOrder(order);
    setOrderDetailOpen(true);
  };

  const handleToggleExpand = (orderId) => {
    setExpandedOrders(prev => {
      const newSet = new Set(prev);
      if (newSet.has(orderId)) {
        newSet.delete(orderId);
      } else {
        newSet.add(orderId);
      }
      return newSet;
    });
  };

  const handlePageChange = (event, newPage) => {
    setPage(newPage);
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'COMPLETED': return 'success';
      case 'PENDING': return 'warning';
      case 'PROCESSING': return 'info';
      case 'SHIPPED': return 'primary';
      case 'CANCELLED': return 'error';
      case 'DELIVERED': return 'success';
      default: return 'default';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'COMPLETED': return <CheckCircle />;
      case 'PENDING': return <AccessTime />;
      case 'PROCESSING': return <LocalShipping />;
      case 'SHIPPED': return <LocalShipping />;
      case 'CANCELLED': return <Cancel />;
      case 'DELIVERED': return <CheckCircle />;
      default: return <AccessTime />;
    }
  };

  const orders = ordersData?.data || [];
  const totalPages = Math.ceil((ordersData?.total || 0) / itemsPerPage);

  const pageTitle = userView ? 'My Orders' : 'Orders Management';
  const pageDescription = userView 
    ? 'View and track your order history' 
    : 'Manage all customer orders and fulfillment';

  if (ordersLoading) {
    return (
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (ordersError) {
    return (
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Alert severity="error" sx={{ mb: 2 }}>
          Failed to load orders. Please try again.
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Typography variant="h4" component="h1" gutterBottom>
            <ShoppingCartIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
            {pageTitle}
          </Typography>
          <Typography variant="body1" color="text.secondary">
            {pageDescription} • {ordersData?.total || 0} orders found
          </Typography>
        </Box>
        {isStaff && !userView && (
          <Fab
            color="primary"
            aria-label="create order"
            onClick={() => navigate('/orders/new')}
          >
            <Add />
          </Fab>
        )}
      </Box>

      {/* Search and Controls */}
      <Box sx={{ mb: 3 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              placeholder="Search orders by ID, customer, or email..."
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
              <Button
                startIcon={<FilterList />}
                onClick={() => setFilterDrawerOpen(true)}
                variant="outlined"
              >
                Filters
              </Button>
              <Button
                startIcon={<Refresh />}
                onClick={() => queryClient.invalidateQueries(['orders'])}
                variant="outlined"
              >
                Refresh
              </Button>
            </Box>
          </Grid>
        </Grid>
      </Box>

      {/* Active Filters */}
      {Object.keys(filters).length > 0 && (
        <Box sx={{ mb: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          {Object.entries(filters).map(([key, value]) => {
            if (!value) return null;
            return (
              <Chip
                key={key}
                label={`${key}: ${value}`}
                onDelete={() => {
                  const newFilters = { ...filters };
                  delete newFilters[key];
                  handleFiltersChange(newFilters);
                }}
                size="small"
                color="primary"
              />
            );
          })}
        </Box>
      )}

      {/* Orders List */}
      {orders.length === 0 ? (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No orders found
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            {searchQuery || Object.keys(filters).length > 0
              ? 'Try adjusting your search or filters'
              : userView ? 'You haven\'t placed any orders yet' : 'No orders have been placed yet'}
          </Typography>
          {isStaff && !userView && (
            <Button
              variant="contained"
              startIcon={<Add />}
              onClick={() => navigate('/orders/new')}
            >
              Create Order
            </Button>
          )}
        </Box>
      ) : (
        <>
          {orders.map((order) => (
            <Card key={order.id} sx={{ mb: 2 }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                  <Box>
                    <Typography variant="h6" gutterBottom>
                      Order #{order.id?.slice(-8) || 'N/A'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      {format(new Date(order.createdAt), 'MMM d, yyyy HH:mm')}
                    </Typography>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                      {getStatusIcon(order.status)}
                      <Chip 
                        label={order.status} 
                        size="small" 
                        color={getStatusColor(order.status)}
                      />
                    </Box>
                    {order.customerInfo && (
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Person sx={{ fontSize: 16 }} />
                        <Typography variant="body2" color="text.secondary">
                          {order.customerInfo.name || order.customerInfo.email}
                        </Typography>
                      </Box>
                    )}
                    {order.storeInfo && (
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Store sx={{ fontSize: 16 }} />
                        <Typography variant="body2" color="text.secondary">
                          {order.storeInfo.name}
                        </Typography>
                      </Box>
                    )}
                  </Box>
                  <Box sx={{ textAlign: 'right' }}>
                    <Typography variant="h6" color="primary">
                      ${(order.total || 0).toFixed(2)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {order.items?.length || 0} items
                    </Typography>
                    {order.source && (
                      <Chip 
                        label={order.source} 
                        size="small" 
                        variant="outlined"
                        sx={{ mt: 1 }}
                      />
                    )}
                  </Box>
                </Box>

                {/* Order Items - Collapsible */}
                <Collapse in={expandedOrders.has(order.id)}>
                  <Divider sx={{ mb: 2 }} />
                  <Typography variant="subtitle2" gutterBottom>
                    Order Items
                  </Typography>
                  {order.items?.map((item, index) => (
                    <Box key={index} sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                      <Typography variant="body2">
                        {item.productName} (x{item.quantity})
                      </Typography>
                      <Typography variant="body2">
                        ${(item.unitPrice * item.quantity).toFixed(2)}
                      </Typography>
                    </Box>
                  ))}
                  <Divider sx={{ my: 1 }} />
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="subtitle2">Total:</Typography>
                    <Typography variant="subtitle2">${(order.total || 0).toFixed(2)}</Typography>
                  </Box>
                </Collapse>
              </CardContent>
              
              <CardActions sx={{ justifyContent: 'space-between' }}>
                <Box sx={{ display: 'flex', gap: 1 }}>
                  <Button size="small" startIcon={<Visibility />} onClick={() => handleOrderView(order)}>
                    View Details
                  </Button>
                  <Button 
                    size="small" 
                    startIcon={expandedOrders.has(order.id) ? <ExpandLess /> : <ExpandMore />}
                    onClick={() => handleToggleExpand(order.id)}
                  >
                    {expandedOrders.has(order.id) ? 'Less' : 'More'}
                  </Button>
                </Box>
                
                {isStaff && (
                  <Box sx={{ display: 'flex', gap: 1 }}>
                    {order.status === 'PENDING' && (
                      <Button size="small" startIcon={<Edit />} onClick={() => navigate(`/orders/${order.id}/edit`)}>
                        Edit
                      </Button>
                    )}
                    <IconButton size="small">
                      <MoreVert />
                    </IconButton>
                  </Box>
                )}
              </CardActions>
            </Card>
          ))}

          {/* Pagination */}
          {totalPages > 1 && (
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
              <Pagination
                count={totalPages}
                page={page}
                onChange={handlePageChange}
                color="primary"
                size="large"
              />
            </Box>
          )}
        </>
      )}

      {/* Filter Drawer */}
      <Drawer anchor="right" open={filterDrawerOpen} onClose={() => setFilterDrawerOpen(false)}>
        <Box sx={{ width: 300, p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Filters
          </Typography>
          <Divider sx={{ mb: 3 }} />
          
          <FormControl fullWidth sx={{ mb: 3 }}>
            <InputLabel>Status</InputLabel>
            <Select
              value={filters.status || ''}
              label="Status"
              onChange={(e) => handleFiltersChange({ ...filters, status: e.target.value })}
            >
              <MenuItem value="">All Status</MenuItem>
              <MenuItem value="PENDING">Pending</MenuItem>
              <MenuItem value="PROCESSING">Processing</MenuItem>
              <MenuItem value="SHIPPED">Shipped</MenuItem>
              <MenuItem value="DELIVERED">Delivered</MenuItem>
              <MenuItem value="COMPLETED">Completed</MenuItem>
              <MenuItem value="CANCELLED">Cancelled</MenuItem>
            </Select>
          </FormControl>

          <FormControl fullWidth sx={{ mb: 3 }}>
            <InputLabel>Source</InputLabel>
            <Select
              value={filters.source || ''}
              label="Source"
              onChange={(e) => handleFiltersChange({ ...filters, source: e.target.value })}
            >
              <MenuItem value="">All Sources</MenuItem>
              <MenuItem value="ONLINE">Online</MenuItem>
              <MenuItem value="POS">POS</MenuItem>
              <MenuItem value="PHONE">Phone</MenuItem>
            </Select>
          </FormControl>

          <Box sx={{ mt: 3, display: 'flex', gap: 1 }}>
            <Button 
              fullWidth 
              variant="outlined" 
              onClick={() => handleFiltersChange({})}
            >
              Clear Filters
            </Button>
            <Button fullWidth variant="contained" onClick={() => setFilterDrawerOpen(false)}>
              Apply
            </Button>
          </Box>
        </Box>
      </Drawer>

      {/* Order Detail Dialog */}
      <Dialog open={orderDetailOpen} onClose={() => setOrderDetailOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="h5">Order #{selectedOrder?.id?.slice(-8) || 'N/A'}</Typography>
            <Chip 
              label={selectedOrder?.status} 
              color={selectedOrder?.status === 'COMPLETED' ? 'success' : selectedOrder?.status === 'CANCELLED' ? 'error' : 'primary'}
            />
          </Box>
        </DialogTitle>
        
        <DialogContent>
          {selectedOrder && (
            <Grid container spacing={3}>
              {/* Order Progress */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom>Order Progress</Typography>
                <Stepper activeStep={['PENDING', 'PROCESSING', 'SHIPPED', 'DELIVERED'].indexOf(selectedOrder.status)} alternativeLabel>
                  <Step><StepLabel>Order Placed</StepLabel></Step>
                  <Step><StepLabel>Processing</StepLabel></Step>
                  <Step><StepLabel>Shipped</StepLabel></Step>
                  <Step><StepLabel>Delivered</StepLabel></Step>
                </Stepper>
              </Grid>

              {/* Order Items */}
              <Grid item xs={12}>
                <Typography variant="h6" gutterBottom>Order Items</Typography>
                <TableContainer component={Paper}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Product</TableCell>
                        <TableCell align="center">Quantity</TableCell>
                        <TableCell align="right">Unit Price</TableCell>
                        <TableCell align="right">Total</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {selectedOrder.items?.map((item, index) => (
                        <TableRow key={index}>
                          <TableCell>{item.productName}</TableCell>
                          <TableCell align="center">{item.quantity}</TableCell>
                          <TableCell align="right">${item.unitPrice?.toFixed(2)}</TableCell>
                          <TableCell align="right">${(item.unitPrice * item.quantity).toFixed(2)}</TableCell>
                        </TableRow>
                      ))}
                      <TableRow>
                        <TableCell colSpan={3}><strong>Total</strong></TableCell>
                        <TableCell align="right"><strong>${(selectedOrder.total || 0).toFixed(2)}</strong></TableCell>
                      </TableRow>
                    </TableBody>
                  </Table>
                </TableContainer>
              </Grid>
            </Grid>
          )}
        </DialogContent>
        
        <DialogActions>
          <Button onClick={() => setOrderDetailOpen(false)}>Close</Button>
          <Button variant="contained" startIcon={<Receipt />}>
            View Receipt
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default OrdersPage;
