import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  Chip,
  CircularProgress,
  Grid,
  FormControl,
  InputLabel,
  Select,
  InputAdornment,
} from '@mui/material';
import {
  Add as AddIcon,
  MoreVert as MoreVertIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  Person as PersonIcon,
  Email as EmailIcon,
  Phone as PhoneIcon,
  LocationOn as LocationIcon,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { userService } from '../services';
import { formatDate } from '../utils/formatters';

const CustomersPage = () => {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCustomer, setSelectedCustomer] = useState(null);
  const [anchorEl, setAnchorEl] = useState(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [customerForm, setCustomerForm] = useState({
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
    status: 'ACTIVE',
    address: {
      street: '',
      city: '',
      state: '',
      zipCode: '',
      country: '',
    },
  });

  // Fetch customers with search and pagination
  const { data: customersData, isLoading } = useQuery({
    queryKey: ['customers', page, rowsPerPage, searchTerm],
    queryFn: () => userService.getUsers({
      page: page + 1,
      limit: rowsPerPage,
      search: searchTerm,
      role: 'USER', // Only get customers, not admin users
    }),
    keepPreviousData: true,
  });

  // Create customer mutation
  const createCustomerMutation = useMutation({
    mutationFn: (customerData) => userService.createUser(customerData),
    onSuccess: () => {
      setSuccess('Customer created successfully');
      queryClient.invalidateQueries(['customers']);
      setCreateDialogOpen(false);
      resetForm();
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to create customer');
    },
  });

  // Update customer mutation
  const updateCustomerMutation = useMutation({
    mutationFn: ({ customerId, customerData }) => userService.updateUser(customerId, customerData),
    onSuccess: () => {
      setSuccess('Customer updated successfully');
      queryClient.invalidateQueries(['customers']);
      setEditDialogOpen(false);
      setSelectedCustomer(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to update customer');
    },
  });

  // Delete customer mutation
  const deleteCustomerMutation = useMutation({
    mutationFn: (customerId) => userService.deleteUser(customerId),
    onSuccess: () => {
      setSuccess('Customer deleted successfully');
      queryClient.invalidateQueries(['customers']);
      setDeleteDialogOpen(false);
      setSelectedCustomer(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to delete customer');
    },
  });

  const customers = customersData?.data || [];
  const totalCount = customersData?.total || 0;

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleMenuClick = (event, customer) => {
    setAnchorEl(event.currentTarget);
    setSelectedCustomer(customer);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleCreateCustomer = () => {
    setCreateDialogOpen(true);
    resetForm();
  };

  const handleEditCustomer = () => {
    if (selectedCustomer) {
      setCustomerForm({
        firstName: selectedCustomer.firstName || '',
        lastName: selectedCustomer.lastName || '',
        email: selectedCustomer.email || '',
        phone: selectedCustomer.phone || '',
        status: selectedCustomer.status || 'ACTIVE',
        address: {
          street: selectedCustomer.address?.street || '',
          city: selectedCustomer.address?.city || '',
          state: selectedCustomer.address?.state || '',
          zipCode: selectedCustomer.address?.zipCode || '',
          country: selectedCustomer.address?.country || '',
        },
      });
      setEditDialogOpen(true);
    }
    handleMenuClose();
  };

  const handleDeleteCustomer = () => {
    setDeleteDialogOpen(true);
    handleMenuClose();
  };

  const handleCreateSubmit = () => {
    const customerData = {
      ...customerForm,
      role: 'USER',
    };
    createCustomerMutation.mutate(customerData);
  };

  const handleUpdateSubmit = () => {
    if (selectedCustomer) {
      updateCustomerMutation.mutate({
        customerId: selectedCustomer.id,
        customerData: customerForm,
      });
    }
  };

  const handleDeleteConfirm = () => {
    if (selectedCustomer) {
      deleteCustomerMutation.mutate(selectedCustomer.id);
    }
  };

  const resetForm = () => {
    setCustomerForm({
      firstName: '',
      lastName: '',
      email: '',
      phone: '',
      status: 'ACTIVE',
      address: {
        street: '',
        city: '',
        state: '',
        zipCode: '',
        country: '',
      },
    });
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'active': return 'success';
      case 'inactive': return 'error';
      case 'pending': return 'warning';
      default: return 'default';
    }
  };

  // Clear alerts after 5 seconds
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
    <Box>
      {/* Header */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">
          <PersonIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
          Customers
        </Typography>
        <Box display="flex" gap={2} alignItems="center">
          <TextField
            size="small"
            placeholder="Search customers..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
          />
          <Button 
            variant="contained" 
            color="primary"
            startIcon={<AddIcon />}
            onClick={handleCreateCustomer}
          >
            Add Customer
          </Button>
        </Box>
      </Box>

      {/* Alerts */}
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

      {/* Customers Table */}
      <Paper elevation={3}>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Phone</TableCell>
                <TableCell>Location</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Joined</TableCell>
                <TableCell align="center">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : customers.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    <Typography variant="body2" color="text.secondary">
                      No customers found
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                customers.map((customer) => (
                  <TableRow key={customer.id} hover>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <PersonIcon color="action" sx={{ mr: 1 }} />
                        <Box>
                          <Typography variant="body2" fontWeight="medium">
                            {customer.firstName} {customer.lastName}
                          </Typography>
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <EmailIcon color="action" sx={{ mr: 1, fontSize: 16 }} />
                        {customer.email}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <PhoneIcon color="action" sx={{ mr: 1, fontSize: 16 }} />
                        {customer.phone || 'N/A'}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <LocationIcon color="action" sx={{ mr: 1, fontSize: 16 }} />
                        {customer.address?.city && customer.address?.state 
                          ? `${customer.address.city}, ${customer.address.state}`
                          : 'N/A'
                        }
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip 
                        label={customer.status || 'Active'} 
                        color={getStatusColor(customer.status)}
                        size="small"
                      />
                    </TableCell>
                    <TableCell>
                      {formatDate(customer.createdAt)}
                    </TableCell>
                    <TableCell align="center">
                      <IconButton
                        size="small"
                        onClick={(e) => handleMenuClick(e, customer)}
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
          component="div"
          count={totalCount}
          page={page}
          onPageChange={handleChangePage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={handleChangeRowsPerPage}
          rowsPerPageOptions={[5, 10, 25, 50]}
        />
      </Paper>

      {/* Action Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleEditCustomer}>
          <EditIcon sx={{ mr: 1 }} fontSize="small" />
          Edit Customer
        </MenuItem>
        <MenuItem onClick={handleDeleteCustomer} sx={{ color: 'error.main' }}>
          <DeleteIcon sx={{ mr: 1 }} fontSize="small" />
          Delete Customer
        </MenuItem>
      </Menu>

      {/* Create Customer Dialog */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Create New Customer</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="First Name"
                value={customerForm.firstName}
                onChange={(e) => setCustomerForm({ ...customerForm, firstName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Last Name"
                value={customerForm.lastName}
                onChange={(e) => setCustomerForm({ ...customerForm, lastName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Email"
                type="email"
                value={customerForm.email}
                onChange={(e) => setCustomerForm({ ...customerForm, email: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Phone"
                value={customerForm.phone}
                onChange={(e) => setCustomerForm({ ...customerForm, phone: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Status</InputLabel>
                <Select
                  value={customerForm.status}
                  label="Status"
                  onChange={(e) => setCustomerForm({ ...customerForm, status: e.target.value })}
                >
                  <MenuItem value="ACTIVE">Active</MenuItem>
                  <MenuItem value="INACTIVE">Inactive</MenuItem>
                  <MenuItem value="PENDING">Pending</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>Address</Typography>
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Street Address"
                value={customerForm.address.street}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, street: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="City"
                value={customerForm.address.city}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, city: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="State"
                value={customerForm.address.state}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, state: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="ZIP Code"
                value={customerForm.address.zipCode}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, zipCode: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Country"
                value={customerForm.address.country}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, country: e.target.value }
                })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleCreateSubmit} 
            variant="contained"
            disabled={!customerForm.firstName || !customerForm.lastName || !customerForm.email || createCustomerMutation.isLoading}
          >
            {createCustomerMutation.isLoading ? <CircularProgress size={20} /> : 'Create Customer'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Customer Dialog */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Edit Customer</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="First Name"
                value={customerForm.firstName}
                onChange={(e) => setCustomerForm({ ...customerForm, firstName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Last Name"
                value={customerForm.lastName}
                onChange={(e) => setCustomerForm({ ...customerForm, lastName: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Email"
                type="email"
                value={customerForm.email}
                onChange={(e) => setCustomerForm({ ...customerForm, email: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Phone"
                value={customerForm.phone}
                onChange={(e) => setCustomerForm({ ...customerForm, phone: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Status</InputLabel>
                <Select
                  value={customerForm.status}
                  label="Status"
                  onChange={(e) => setCustomerForm({ ...customerForm, status: e.target.value })}
                >
                  <MenuItem value="ACTIVE">Active</MenuItem>
                  <MenuItem value="INACTIVE">Inactive</MenuItem>
                  <MenuItem value="PENDING">Pending</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <Typography variant="h6" sx={{ mt: 2, mb: 1 }}>Address</Typography>
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Street Address"
                value={customerForm.address.street}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, street: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="City"
                value={customerForm.address.city}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, city: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="State"
                value={customerForm.address.state}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, state: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="ZIP Code"
                value={customerForm.address.zipCode}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, zipCode: e.target.value }
                })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Country"
                value={customerForm.address.country}
                onChange={(e) => setCustomerForm({ 
                  ...customerForm, 
                  address: { ...customerForm.address, country: e.target.value }
                })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleUpdateSubmit} 
            variant="contained"
            disabled={!customerForm.firstName || !customerForm.lastName || !customerForm.email || updateCustomerMutation.isLoading}
          >
            {updateCustomerMutation.isLoading ? <CircularProgress size={20} /> : 'Update Customer'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Customer Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Delete Customer</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary">
            Are you sure you want to delete customer "{selectedCustomer?.firstName} {selectedCustomer?.lastName}"? 
            This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleDeleteConfirm} 
            color="error" 
            variant="contained"
            disabled={deleteCustomerMutation.isLoading}
          >
            {deleteCustomerMutation.isLoading ? <CircularProgress size={20} /> : 'Delete Customer'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default CustomersPage;
