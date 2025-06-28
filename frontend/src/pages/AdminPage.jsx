import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  Grid,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Chip,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  CircularProgress,
  Tabs,
  Tab,
  FormControl,
  InputLabel,
  Select,
  Divider,
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  People as PeopleIcon,
  ShoppingCart as OrdersIcon,
  Analytics as AnalyticsIcon,
  HealthAndSafety as HealthIcon,
  MoreVert as MoreVertIcon,
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminService } from '../services';
import { formatCurrency, formatDate } from '../utils/formatters';

const AdminPage = () => {
  const queryClient = useQueryClient();
  const [tabValue, setTabValue] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedUser, setSelectedUser] = useState(null);
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [anchorEl, setAnchorEl] = useState(null);
  const [userDialogOpen, setUserDialogOpen] = useState(false);
  const [orderDialogOpen, setOrderDialogOpen] = useState(false);
  const [createUserDialogOpen, setCreateUserDialogOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [userForm, setUserForm] = useState({
    firstName: '',
    lastName: '',
    email: '',
    role: 'USER',
    status: 'ACTIVE',
  });

  // Fetch admin analytics
  const { data: analytics, isLoading: analyticsLoading } = useQuery({
    queryKey: ['admin-analytics'],
    queryFn: () => adminService.getAnalytics(),
  });

  // Fetch system health
  const { data: systemHealth, isLoading: healthLoading } = useQuery({
    queryKey: ['admin-system-health'],
    queryFn: () => adminService.getSystemHealth(),
    refetchInterval: 30000, // Refresh every 30 seconds
  });

  // Fetch users for admin management
  const { data: usersData, isLoading: usersLoading, refetch: refetchUsers } = useQuery({
    queryKey: ['admin-users', page, rowsPerPage],
    queryFn: () => adminService.getUsers({
      page: page + 1,
      limit: rowsPerPage,
    }),
    enabled: tabValue === 1,
    keepPreviousData: true,
  });

  // Fetch orders for admin management
  const { data: ordersData, isLoading: ordersLoading, refetch: refetchOrders } = useQuery({
    queryKey: ['admin-orders', page, rowsPerPage],
    queryFn: () => adminService.getOrders({
      page: page + 1,
      limit: rowsPerPage,
    }),
    enabled: tabValue === 2,
    keepPreviousData: true,
  });

  // Create user mutation
  const createUserMutation = useMutation({
    mutationFn: (userData) => adminService.createUser(userData),
    onSuccess: () => {
      setSuccess('User created successfully');
      queryClient.invalidateQueries(['admin-users']);
      setCreateUserDialogOpen(false);
      resetUserForm();
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to create user');
    },
  });

  // Update user mutation
  const updateUserMutation = useMutation({
    mutationFn: ({ userId, userData }) => adminService.updateUser(userId, userData),
    onSuccess: () => {
      setSuccess('User updated successfully');
      queryClient.invalidateQueries(['admin-users']);
      setUserDialogOpen(false);
      setSelectedUser(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to update user');
    },
  });

  // Delete user mutation
  const deleteUserMutation = useMutation({
    mutationFn: (userId) => adminService.deleteUser(userId),
    onSuccess: () => {
      setSuccess('User deleted successfully');
      queryClient.invalidateQueries(['admin-users']);
      setSelectedUser(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to delete user');
    },
  });

  // Update order status mutation
  const updateOrderStatusMutation = useMutation({
    mutationFn: ({ orderId, status }) => adminService.updateOrderStatus(orderId, { status }),
    onSuccess: () => {
      setSuccess('Order status updated successfully');
      queryClient.invalidateQueries(['admin-orders']);
      setOrderDialogOpen(false);
      setSelectedOrder(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to update order status');
    },
  });

  const users = usersData?.data || [];
  const orders = ordersData?.data || [];
  const totalUsers = usersData?.total || 0;
  const totalOrders = ordersData?.total || 0;

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
    setPage(0);
  };

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleMenuClick = (event, item, type) => {
    setAnchorEl(event.currentTarget);
    if (type === 'user') {
      setSelectedUser(item);
    } else if (type === 'order') {
      setSelectedOrder(item);
    }
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedUser(null);
    setSelectedOrder(null);
  };

  const handleCreateUser = () => {
    if (userForm.firstName && userForm.lastName && userForm.email) {
      createUserMutation.mutate(userForm);
    }
  };

  const handleUpdateUser = () => {
    if (selectedUser) {
      updateUserMutation.mutate({
        userId: selectedUser.id,
        userData: userForm,
      });
    }
  };

  const handleDeleteUser = () => {
    if (selectedUser) {
      deleteUserMutation.mutate(selectedUser.id);
    }
  };

  const resetUserForm = () => {
    setUserForm({
      firstName: '',
      lastName: '',
      email: '',
      role: 'USER',
      status: 'ACTIVE',
    });
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'active': return 'success';
      case 'inactive': return 'error';
      case 'pending': return 'warning';
      case 'blocked': return 'error';
      default: return 'default';
    }
  };

  const getHealthColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'healthy': return 'success';
      case 'warning': return 'warning';
      case 'error': return 'error';
      default: return 'default';
    }
  };

  useEffect(() => {
    if (error || success) {
      const timer = setTimeout(() => {
        setError('');
        setSuccess('');
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [error, success]);

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1" sx={{ fontWeight: 'bold' }}>
          Admin Dashboard
        </Typography>
        <Button
          variant="outlined"
          startIcon={<RefreshIcon />}
          onClick={() => {
            queryClient.invalidateQueries(['admin-analytics']);
            queryClient.invalidateQueries(['admin-system-health']);
            if (tabValue === 1) refetchUsers();
            if (tabValue === 2) refetchOrders();
          }}
        >
          Refresh
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      {success && (
        <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccess('')}>
          {success}
        </Alert>
      )}

      <Tabs value={tabValue} onChange={handleTabChange} sx={{ mb: 3 }}>
        <Tab icon={<DashboardIcon />} label="Overview" />
        <Tab icon={<PeopleIcon />} label="User Management" />
        <Tab icon={<OrdersIcon />} label="Order Management" />
        <Tab icon={<AnalyticsIcon />} label="Analytics" />
        <Tab icon={<HealthIcon />} label="System Health" />
      </Tabs>

      {/* Overview Tab */}
      {tabValue === 0 && (
        <Grid container spacing={3}>
          {analytics && (
            <>
              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" color="text.secondary">Total Users</Typography>
                    <Typography variant="h4" color="primary">
                      {analytics.totalUsers || 0}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" color="text.secondary">Total Orders</Typography>
                    <Typography variant="h4" color="primary">
                      {analytics.totalOrders || 0}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" color="text.secondary">Revenue</Typography>
                    <Typography variant="h4" color="primary">
                      {formatCurrency(analytics.totalRevenue || 0)}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" color="text.secondary">Active Products</Typography>
                    <Typography variant="h4" color="primary">
                      {analytics.activeProducts || 0}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </>
          )}

          {systemHealth && (
            <Grid item xs={12}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>System Health</Typography>
                  <Grid container spacing={2}>
                    {Object.entries(systemHealth).map(([service, status]) => (
                      <Grid item xs={12} sm={6} md={4} key={service}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Chip
                            label={status?.status || 'Unknown'}
                            color={getHealthColor(status?.status)}
                            size="small"
                          />
                          <Typography variant="body2">
                            {service.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase())}
                          </Typography>
                        </Box>
                      </Grid>
                    ))}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>
          )}
        </Grid>
      )}

      {/* User Management Tab */}
      {tabValue === 1 && (
        <Box>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6">User Management</Typography>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => setCreateUserDialogOpen(true)}
            >
              Create User
            </Button>
          </Box>

          <Paper sx={{ width: '100%', overflow: 'hidden' }}>
            <TableContainer>
              <Table stickyHeader>
                <TableHead>
                  <TableRow>
                    <TableCell>Name</TableCell>
                    <TableCell>Email</TableCell>
                    <TableCell>Role</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Created</TableCell>
                    <TableCell align="center">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {usersLoading ? (
                    <TableRow>
                      <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                        <CircularProgress />
                      </TableCell>
                    </TableRow>
                  ) : users.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                        <Typography variant="body1" color="text.secondary">
                          No users found
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ) : (
                    users.map((user) => (
                      <TableRow key={user.id} hover>
                        <TableCell>
                          <Typography variant="body2" fontWeight="medium">
                            {user.firstName} {user.lastName}
                          </Typography>
                        </TableCell>
                        <TableCell>{user.email}</TableCell>
                        <TableCell>
                          <Chip
                            label={user.role}
                            color="primary"
                            size="small"
                            variant="outlined"
                          />
                        </TableCell>
                        <TableCell>
                          <Chip
                            label={user.status || 'Active'}
                            color={getStatusColor(user.status)}
                            size="small"
                            variant="outlined"
                          />
                        </TableCell>
                        <TableCell>{formatDate(user.createdAt)}</TableCell>
                        <TableCell align="center">
                          <IconButton
                            size="small"
                            onClick={(e) => handleMenuClick(e, user, 'user')}
                          >
                            <MoreVertIcon />
                          </IconButton>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TableContainer>
            
            <TablePagination
              rowsPerPageOptions={[5, 10, 25, 50]}
              component="div"
              count={totalUsers}
              rowsPerPage={rowsPerPage}
              page={page}
              onPageChange={handleChangePage}
              onRowsPerPageChange={handleChangeRowsPerPage}
            />
          </Paper>
        </Box>
      )}

      {/* Order Management Tab */}
      {tabValue === 2 && (
        <Box>
          <Typography variant="h6" sx={{ mb: 2 }}>Order Management</Typography>

          <Paper sx={{ width: '100%', overflow: 'hidden' }}>
            <TableContainer>
              <Table stickyHeader>
                <TableHead>
                  <TableRow>
                    <TableCell>Order ID</TableCell>
                    <TableCell>Customer</TableCell>
                    <TableCell>Total</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Date</TableCell>
                    <TableCell align="center">Actions</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {ordersLoading ? (
                    <TableRow>
                      <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                        <CircularProgress />
                      </TableCell>
                    </TableRow>
                  ) : orders.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                        <Typography variant="body1" color="text.secondary">
                          No orders found
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ) : (
                    orders.map((order) => (
                      <TableRow key={order.id} hover>
                        <TableCell>
                          <Typography variant="body2" fontWeight="medium">
                            {order.id}
                          </Typography>
                        </TableCell>
                        <TableCell>{order.customerName || order.customerId}</TableCell>
                        <TableCell>{formatCurrency(order.totalAmount)}</TableCell>
                        <TableCell>
                          <Chip
                            label={order.status}
                            color={getStatusColor(order.status)}
                            size="small"
                            variant="outlined"
                          />
                        </TableCell>
                        <TableCell>{formatDate(order.createdAt)}</TableCell>
                        <TableCell align="center">
                          <IconButton
                            size="small"
                            onClick={(e) => handleMenuClick(e, order, 'order')}
                          >
                            <MoreVertIcon />
                          </IconButton>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TableContainer>
            
            <TablePagination
              rowsPerPageOptions={[5, 10, 25, 50]}
              component="div"
              count={totalOrders}
              rowsPerPage={rowsPerPage}
              page={page}
              onPageChange={handleChangePage}
              onRowsPerPageChange={handleChangeRowsPerPage}
            />
          </Paper>
        </Box>
      )}

      {/* Analytics Tab */}
      {tabValue === 3 && (
        <Box>
          <Typography variant="h6" sx={{ mb: 2 }}>Analytics Dashboard</Typography>
          {analyticsLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : analytics ? (
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>User Statistics</Typography>
                    <Typography variant="body2">Total Users: {analytics.totalUsers}</Typography>
                    <Typography variant="body2">New Users (30d): {analytics.newUsers30d || 0}</Typography>
                    <Typography variant="body2">Active Users: {analytics.activeUsers || 0}</Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} md={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>Sales Statistics</Typography>
                    <Typography variant="body2">Total Revenue: {formatCurrency(analytics.totalRevenue)}</Typography>
                    <Typography variant="body2">Monthly Revenue: {formatCurrency(analytics.monthlyRevenue || 0)}</Typography>
                    <Typography variant="body2">Average Order: {formatCurrency(analytics.averageOrderValue || 0)}</Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          ) : (
            <Typography variant="body2" color="text.secondary">
              No analytics data available
            </Typography>
          )}
        </Box>
      )}

      {/* System Health Tab */}
      {tabValue === 4 && (
        <Box>
          <Typography variant="h6" sx={{ mb: 2 }}>System Health Monitor</Typography>
          {healthLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : systemHealth ? (
            <Grid container spacing={2}>
              {Object.entries(systemHealth).map(([service, status]) => (
                <Grid item xs={12} sm={6} md={4} key={service}>
                  <Card>
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        {service.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase())}
                      </Typography>
                      <Chip
                        label={status?.status || 'Unknown'}
                        color={getHealthColor(status?.status)}
                        sx={{ mb: 1 }}
                      />
                      {status?.message && (
                        <Typography variant="body2" color="text.secondary">
                          {status.message}
                        </Typography>
                      )}
                      {status?.lastCheck && (
                        <Typography variant="caption" color="text.secondary">
                          Last check: {formatDate(status.lastCheck)}
                        </Typography>
                      )}
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          ) : (
            <Typography variant="body2" color="text.secondary">
              No health data available
            </Typography>
          )}
        </Box>
      )}

      {/* User Action Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl) && Boolean(selectedUser)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={() => { setUserDialogOpen(true); handleMenuClose(); }}>
          <EditIcon sx={{ mr: 1 }} fontSize="small" />
          Edit User
        </MenuItem>
        <MenuItem onClick={() => { handleDeleteUser(); handleMenuClose(); }}>
          <DeleteIcon sx={{ mr: 1 }} fontSize="small" />
          Delete User
        </MenuItem>
      </Menu>

      {/* Order Action Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl) && Boolean(selectedOrder)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={() => { setOrderDialogOpen(true); handleMenuClose(); }}>
          <EditIcon sx={{ mr: 1 }} fontSize="small" />
          Update Status
        </MenuItem>
      </Menu>

      {/* Create User Dialog */}
      <Dialog open={createUserDialogOpen} onClose={() => setCreateUserDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create New User</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="First Name"
                value={userForm.firstName}
                onChange={(e) => setUserForm({ ...userForm, firstName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Last Name"
                value={userForm.lastName}
                onChange={(e) => setUserForm({ ...userForm, lastName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Email"
                type="email"
                value={userForm.email}
                onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Role</InputLabel>
                <Select
                  value={userForm.role}
                  label="Role"
                  onChange={(e) => setUserForm({ ...userForm, role: e.target.value })}
                >
                  <MenuItem value="USER">User</MenuItem>
                  <MenuItem value="MANAGER">Manager</MenuItem>
                  <MenuItem value="ADMIN">Admin</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Status</InputLabel>
                <Select
                  value={userForm.status}
                  label="Status"
                  onChange={(e) => setUserForm({ ...userForm, status: e.target.value })}
                >
                  <MenuItem value="ACTIVE">Active</MenuItem>
                  <MenuItem value="INACTIVE">Inactive</MenuItem>
                  <MenuItem value="PENDING">Pending</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateUserDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleCreateUser} 
            variant="contained"
            disabled={!userForm.firstName || !userForm.lastName || !userForm.email || createUserMutation.isLoading}
          >
            {createUserMutation.isLoading ? <CircularProgress size={20} /> : 'Create User'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default AdminPage;
