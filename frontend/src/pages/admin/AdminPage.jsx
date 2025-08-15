import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Paper,
  Tabs,
  Tab,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  Chip,
  IconButton,
  TablePagination,
  Switch,
  FormControlLabel,
  Alert,
  CircularProgress,
  Divider,
} from '@mui/material';
import {
  AdminPanelSettings,
  People,
  Analytics,
  Settings,
  Add,
  Edit,
  Delete,
  CheckCircle,
  Cancel,
  Refresh,
  PersonAdd,
  Security,
  Dashboard as DashboardIcon,
  TrendingUp,
  Group,
  Business,
  Inventory,
  ShoppingCart,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSnackbar } from 'notistack';
import { useAuth } from '../../hooks/useAuth';
import userService from '../../services/userService';
import productService from '../../services/productService';
import orderService from '../../services/orderService';
import inventoryService from '../../services/inventoryService';
import supplierService from '../../services/supplierService';
import storeService from '../../services/storeService';

const AdminPage = () => {
  const { user } = useAuth();
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState(0);
  const [userDialog, setUserDialog] = useState({ open: false, user: null, mode: 'create' });
  const [userPage, setUserPage] = useState(0);
  const [userRowsPerPage, setUserRowsPerPage] = useState(10);
  const [newUser, setNewUser] = useState({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    phone: '',
    role: 'CUSTOMER',
    isActive: true,
  });

  // System Analytics Queries
  const { data: systemStats, isLoading: statsLoading } = useQuery({
    queryKey: ['admin-stats'],
    queryFn: async () => {
      const [users, products, orders, inventory, suppliers, stores] = await Promise.all([
        userService.listUsers({ page: 1, limit: 1 }),
        productService.listProducts({ page: 1, limit: 1 }),
        orderService.listOrders({ page: 1, limit: 1 }),
        inventoryService.listInventory({ page: 1, limit: 1 }),
        supplierService.listSuppliers({ page: 1, limit: 1 }),
        storeService.listStores({ page: 1, limit: 1 }),
      ]);
      return {
        totalUsers: users.total || users.users?.length || 0,
        totalProducts: products.total || products.products?.length || 0,
        totalOrders: orders.total || orders.orders?.length || 0,
        totalInventoryItems: inventory.total || inventory.items?.length || 0,
        totalSuppliers: suppliers.total || suppliers.suppliers?.length || 0,
        totalStores: stores.total || stores.stores?.length || 0,
      };
    },
    enabled: activeTab === 1, // Only load when analytics tab is active
  });

  // Users Query
  const { data: usersData, isLoading: usersLoading } = useQuery({
    queryKey: ['admin-users', userPage, userRowsPerPage],
    queryFn: () => userService.listUsers({ 
      page: userPage + 1, 
      limit: userRowsPerPage 
    }),
    enabled: activeTab === 0,
  });

  // User Management Mutations
  const createUserMutation = useMutation({
    mutationFn: async (userData) => userService.registerUser(userData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] });
      queryClient.invalidateQueries({ queryKey: ['admin-stats'] });
      setUserDialog({ open: false, user: null, mode: 'create' });
      setNewUser({ email: '', password: '', firstName: '', lastName: '', phone: '', role: 'CUSTOMER', isActive: true });
      enqueueSnackbar('User created successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.response?.data?.message || 'Failed to create user', { variant: 'error' });
    },
  });

  const toggleUserStatusMutation = useMutation({
    mutationFn: async ({ userId, activate }) => {
      if (activate) {
        return userService.activateUser(userId);
      } else {
        return userService.deactivateUser(userId);
      }
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] });
      enqueueSnackbar(
        `User ${variables.activate ? 'activated' : 'deactivated'} successfully`, 
        { variant: 'success' }
      );
    },
    onError: (error) => {
      enqueueSnackbar(error.response?.data?.message || 'Failed to update user status', { variant: 'error' });
    },
  });

  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };

  const handleCreateUser = () => {
    setUserDialog({ open: true, user: null, mode: 'create' });
    setNewUser({ email: '', password: '', firstName: '', lastName: '', phone: '', role: 'CUSTOMER', isActive: true });
  };

  const handleSubmitUser = () => {
    if (!newUser.email || !newUser.password || !newUser.firstName || !newUser.lastName) {
      enqueueSnackbar('Please fill in all required fields', { variant: 'warning' });
      return;
    }
    createUserMutation.mutate(newUser);
  };

  const handleToggleUserStatus = (userId, isActive) => {
    toggleUserStatusMutation.mutate({ userId, activate: !isActive });
  };

  const getRoleColor = (role) => {
    switch (role) {
      case 'ADMIN': return 'error';
      case 'STAFF': return 'primary';
      case 'CUSTOMER': return 'success';
      default: return 'default';
    }
  };

  const getStatusColor = (isActive) => {
    return isActive ? 'success' : 'default';
  };

  const StatCard = ({ title, value, icon, color = 'primary' }) => (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box>
            <Typography color="text.secondary" gutterBottom variant="h6">
              {title}
            </Typography>
            <Typography variant="h4" component="div" color={`${color}.main`}>
              {value}
            </Typography>
          </Box>
          <Box sx={{ color: `${color}.main` }}>
            {icon}
          </Box>
        </Box>
      </CardContent>
    </Card>
  );

  const renderUsersTab = () => (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h6">User Management</Typography>
        <Button
          variant="contained"
          startIcon={<PersonAdd />}
          onClick={handleCreateUser}
          sx={{ ml: 'auto' }}
        >
          Create User
        </Button>
      </Box>

      {usersLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Card>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  <TableCell>Email</TableCell>
                  <TableCell>Phone</TableCell>
                  <TableCell>Role</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {usersData?.users?.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell>
                      <Typography variant="body1">
                        {user.firstName} {user.lastName}
                      </Typography>
                    </TableCell>
                    <TableCell>{user.email}</TableCell>
                    <TableCell>{user.phone || 'N/A'}</TableCell>
                    <TableCell>
                      <Chip
                        label={user.role}
                        color={getRoleColor(user.role)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={user.isActive ? 'Active' : 'Inactive'}
                        color={getStatusColor(user.isActive)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      <IconButton
                        onClick={() => handleToggleUserStatus(user.id, user.isActive)}
                        color={user.isActive ? 'error' : 'success'}
                        size="small"
                      >
                        {user.isActive ? <Cancel /> : <CheckCircle />}
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
          <TablePagination
            rowsPerPageOptions={[5, 10, 25]}
            component="div"
            count={usersData?.total || 0}
            rowsPerPage={userRowsPerPage}
            page={userPage}
            onPageChange={(event, newPage) => setUserPage(newPage)}
            onRowsPerPageChange={(event) => {
              setUserRowsPerPage(parseInt(event.target.value, 10));
              setUserPage(0);
            }}
          />
        </Card>
      )}
    </Box>
  );

  const renderAnalyticsTab = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        System Analytics
      </Typography>
      
      {statsLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6} md={4}>
            <StatCard
              title="Total Users"
              value={systemStats?.totalUsers || 0}
              icon={<People sx={{ fontSize: 40 }} />}
              color="primary"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={4}>
            <StatCard
              title="Total Products"
              value={systemStats?.totalProducts || 0}
              icon={<Inventory sx={{ fontSize: 40 }} />}
              color="secondary"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={4}>
            <StatCard
              title="Total Orders"
              value={systemStats?.totalOrders || 0}
              icon={<ShoppingCart sx={{ fontSize: 40 }} />}
              color="success"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={4}>
            <StatCard
              title="Inventory Items"
              value={systemStats?.totalInventoryItems || 0}
              icon={<Business sx={{ fontSize: 40 }} />}
              color="info"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={4}>
            <StatCard
              title="Suppliers"
              value={systemStats?.totalSuppliers || 0}
              icon={<Group sx={{ fontSize: 40 }} />}
              color="warning"
            />
          </Grid>
          <Grid item xs={12} sm={6} md={4}>
            <StatCard
              title="Stores"
              value={systemStats?.totalStores || 0}
              icon={<Business sx={{ fontSize: 40 }} />}
              color="error"
            />
          </Grid>
        </Grid>
      )}
      
      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            System Health
          </Typography>
          <Alert severity="success" sx={{ mb: 2 }}>
            All services are running normally
          </Alert>
          <Typography variant="body2" color="text.secondary">
            Last updated: {new Date().toLocaleString()}
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );

  const renderSettingsTab = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        System Configuration
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Security Settings
              </Typography>
              <FormControlLabel
                control={<Switch defaultChecked />}
                label="Enable JWT Authentication"
                sx={{ mb: 2, display: 'block' }}
              />
              <FormControlLabel
                control={<Switch defaultChecked />}
                label="Enforce Password Policy"
                sx={{ mb: 2, display: 'block' }}
              />
              <FormControlLabel
                control={<Switch defaultChecked />}
                label="Enable Rate Limiting"
                sx={{ mb: 2, display: 'block' }}
              />
              <Button variant="outlined" startIcon={<Security />}>
                Update Security Settings
              </Button>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                System Maintenance
              </Typography>
              <Button 
                variant="outlined" 
                startIcon={<Refresh />} 
                sx={{ mb: 2, display: 'block', width: '100%' }}
              >
                Clear System Cache
              </Button>
              <Button 
                variant="outlined" 
                startIcon={<DashboardIcon />} 
                sx={{ mb: 2, display: 'block', width: '100%' }}
              >
                Generate System Report
              </Button>
              <Button 
                variant="outlined" 
                color="warning"
                sx={{ mb: 2, display: 'block', width: '100%' }}
              >
                Backup Database
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );

  return (
    <Box>
      <Paper sx={{ p: 3, mb: 3, bgcolor: 'primary.main', color: 'primary.contrastText' }}>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <AdminPanelSettings sx={{ fontSize: 40, mr: 2 }} />
          <Box>
            <Typography variant="h4" gutterBottom>
              Admin Panel
            </Typography>
            <Typography variant="subtitle1">
              System administration, user management, and analytics
            </Typography>
          </Box>
        </Box>
      </Paper>

      <Card>
        <Tabs 
          value={activeTab} 
          onChange={handleTabChange}
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab icon={<People />} label="Users" />
          <Tab icon={<Analytics />} label="Analytics" />
          <Tab icon={<Settings />} label="Settings" />
        </Tabs>
        
        <Box sx={{ p: 3 }}>
          {activeTab === 0 && renderUsersTab()}
          {activeTab === 1 && renderAnalyticsTab()}
          {activeTab === 2 && renderSettingsTab()}
        </Box>
      </Card>

      {/* Create User Dialog */}
      <Dialog 
        open={userDialog.open} 
        onClose={() => setUserDialog({ open: false, user: null, mode: 'create' })}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <PersonAdd sx={{ mr: 1 }} />
            Create New User
          </Box>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="First Name"
                value={newUser.firstName}
                onChange={(e) => setNewUser({ ...newUser, firstName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Last Name"
                value={newUser.lastName}
                onChange={(e) => setNewUser({ ...newUser, lastName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Email"
                type="email"
                value={newUser.email}
                onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Phone"
                value={newUser.phone}
                onChange={(e) => setNewUser({ ...newUser, phone: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Password"
                type="password"
                value={newUser.password}
                onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Role</InputLabel>
                <Select
                  value={newUser.role}
                  onChange={(e) => setNewUser({ ...newUser, role: e.target.value })}
                  label="Role"
                >
                  <MenuItem value="CUSTOMER">Customer</MenuItem>
                  <MenuItem value="STAFF">Staff</MenuItem>
                  <MenuItem value="ADMIN">Admin</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <FormControlLabel
                control={
                  <Switch
                    checked={newUser.isActive}
                    onChange={(e) => setNewUser({ ...newUser, isActive: e.target.checked })}
                  />
                }
                label="Active User"
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button 
            onClick={() => setUserDialog({ open: false, user: null, mode: 'create' })}
          >
            Cancel
          </Button>
          <Button 
            variant="contained" 
            onClick={handleSubmitUser}
            disabled={createUserMutation.isPending}
          >
            {createUserMutation.isPending ? 'Creating...' : 'Create User'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default AdminPage;
