import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TablePagination,
  IconButton,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  Alert,
  Chip,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  CircularProgress,
  Grid,
} from '@mui/material';
import {
  MoreVert as MoreVertIcon,
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Sync as SyncIcon,
  Science as TestIcon,
  Visibility as ViewIcon,
  Cable as AdapterIcon,
  CheckCircle as CheckIcon,
  Error as ErrorIcon,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { supplierService } from '../services';
import { formatDate } from '../utils/formatters';

const SuppliersPage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedSupplier, setSelectedSupplier] = useState(null);
  const [anchorEl, setAnchorEl] = useState(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [adapterDialogOpen, setAdapterDialogOpen] = useState(false);
  const [testConnectionDialogOpen, setTestConnectionDialogOpen] = useState(false);
  const [createSupplierDialogOpen, setCreateSupplierDialogOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [supplierForm, setSupplierForm] = useState({
    name: '',
    contactEmail: '',
    contactPhone: '',
    address: '',
    website: '',
    description: '',
    status: 'ACTIVE',
  });
  const [testResults, setTestResults] = useState(null);

  // Fetch suppliers
  const { data: suppliersData, isLoading, refetch } = useQuery({
    queryKey: ['suppliers', page, rowsPerPage],
    queryFn: () => supplierService.getSuppliers({
      page: page + 1,
      limit: rowsPerPage,
    }),
    keepPreviousData: true,
  });

  // Fetch supplier adapters
  const { data: adaptersData, isLoading: adaptersLoading } = useQuery({
    queryKey: ['supplier-adapters'],
    queryFn: () => supplierService.getSupplierAdapters(),
    enabled: adapterDialogOpen,
  });

  // Create supplier mutation
  const createSupplierMutation = useMutation({
    mutationFn: (supplierData) => supplierService.createSupplier(supplierData),
    onSuccess: () => {
      setSuccess('Supplier created successfully');
      queryClient.invalidateQueries(['suppliers']);
      setCreateSupplierDialogOpen(false);
      resetSupplierForm();
      setError(''); // Clear any previous errors
    },
    onError: (error) => {
      console.error('Create supplier error:', error);
      setError(error.response?.data?.message || error.message || 'Failed to create supplier');
    },
  });

  // Delete supplier mutation
  const deleteSupplierMutation = useMutation({
    mutationFn: (supplierId) => supplierService.deleteSupplier(supplierId),
    onSuccess: () => {
      setSuccess('Supplier deleted successfully');
      queryClient.invalidateQueries(['suppliers']);
      setDeleteDialogOpen(false);
      setSelectedSupplier(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to delete supplier');
    },
  });

  // Test connection mutation
  const testConnectionMutation = useMutation({
    mutationFn: ({ supplierId, adapterId }) => 
      supplierService.testSupplierConnection(supplierId, adapterId),
    onSuccess: (data) => {
      setTestResults(data);
      setSuccess('Connection test completed');
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Connection test failed');
      setTestResults({ success: false, error: error.message });
    },
  });

  const suppliers = Array.isArray(suppliersData)
  ? suppliersData
  : Array.isArray(suppliersData?.items)
    ? suppliersData.items
    : Array.isArray(suppliersData?.data)
      ? suppliersData.data
      : Array.isArray(suppliersData?.suppliers)
        ? suppliersData.suppliers
        : [];

  const totalCount = suppliersData?.total || 0;
  const adapters = adaptersData || [];

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleMenuClick = (event, supplier) => {
    setAnchorEl(event.currentTarget);
    setSelectedSupplier(supplier);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedSupplier(null);
  };

  const handleCreateSupplier = () => {
    if (supplierForm.name && supplierForm.contactEmail) {
      createSupplierMutation.mutate(supplierForm);
    }
  };

  const handleDeleteSupplier = () => {
    if (selectedSupplier) {
      deleteSupplierMutation.mutate(selectedSupplier.id);
    }
  };

  const handleTestConnection = (adapterId) => {
    if (selectedSupplier && adapterId) {
      testConnectionMutation.mutate({
        supplierId: selectedSupplier.id,
        adapterId,
      });
    }
  };

  const resetSupplierForm = () => {
    setSupplierForm({
      name: '',
      contactEmail: '',
      contactPhone: '',
      address: '',
      website: '',
      description: '',
      status: 'ACTIVE',
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
          Supplier Management
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<SyncIcon />}
            onClick={() => refetch()}
            disabled={isLoading}
          >
            Refresh
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => setCreateSupplierDialogOpen(true)}
          >
            Add Supplier
          </Button>
        </Box>
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

      <Paper sx={{ width: '100%', overflow: 'hidden' }}>
        <TableContainer>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Contact Email</TableCell>
                <TableCell>Phone</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Created</TableCell>
                <TableCell align="center">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : suppliers.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center" sx={{ py: 4 }}>
                    <Typography variant="body1" color="text.secondary">
                      No suppliers found
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                suppliers.map((supplier) => (
                  <TableRow key={supplier.id} hover>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {supplier.name}
                      </Typography>
                      {supplier.website && (
                        <Typography variant="caption" color="text.secondary">
                          {supplier.website}
                        </Typography>
                      )}
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {supplier.contactEmail}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {supplier.contactPhone || '-'}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={supplier.status || 'Active'}
                        color={getStatusColor(supplier.status)}
                        size="small"
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {formatDate(supplier.createdAt)}
                      </Typography>
                    </TableCell>
                    <TableCell align="center">
                      <IconButton
                        size="small"
                        onClick={(e) => handleMenuClick(e, supplier)}
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
          count={totalCount}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </Paper>

      {/* Action Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={() => { navigate(`/suppliers/${selectedSupplier?.id}`); handleMenuClose(); }}>
          <ViewIcon sx={{ mr: 1 }} fontSize="small" />
          View Details
        </MenuItem>
        <MenuItem onClick={() => { navigate(`/suppliers/${selectedSupplier?.id}/edit`); handleMenuClose(); }}>
          <EditIcon sx={{ mr: 1 }} fontSize="small" />
          Edit Supplier
        </MenuItem>
        <MenuItem onClick={() => { setAdapterDialogOpen(true); handleMenuClose(); }}>
          <AdapterIcon sx={{ mr: 1 }} fontSize="small" />
          Manage Adapters
        </MenuItem>
        <MenuItem onClick={() => { setTestConnectionDialogOpen(true); handleMenuClose(); }}>
          <TestIcon sx={{ mr: 1 }} fontSize="small" />
          Test Connection
        </MenuItem>
        <MenuItem 
          onClick={() => { setDeleteDialogOpen(true); handleMenuClose(); }}
          sx={{ color: 'error.main' }}
        >
          <DeleteIcon sx={{ mr: 1 }} fontSize="small" />
          Delete Supplier
        </MenuItem>
      </Menu>

      {/* Create Supplier Dialog */}
      <Dialog open={createSupplierDialogOpen} onClose={() => setCreateSupplierDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Add New Supplier</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Supplier Name"
                value={supplierForm.name}
                onChange={(e) => setSupplierForm({ ...supplierForm, name: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Contact Email"
                type="email"
                value={supplierForm.contactEmail}
                onChange={(e) => setSupplierForm({ ...supplierForm, contactEmail: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Contact Phone"
                value={supplierForm.contactPhone}
                onChange={(e) => setSupplierForm({ ...supplierForm, contactPhone: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Website"
                value={supplierForm.website}
                onChange={(e) => setSupplierForm({ ...supplierForm, website: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Address"
                multiline
                rows={2}
                value={supplierForm.address}
                onChange={(e) => setSupplierForm({ ...supplierForm, address: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Description"
                multiline
                rows={3}
                value={supplierForm.description}
                onChange={(e) => setSupplierForm({ ...supplierForm, description: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Status</InputLabel>
                <Select
                  value={supplierForm.status}
                  label="Status"
                  onChange={(e) => setSupplierForm({ ...supplierForm, status: e.target.value })}
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
          <Button onClick={() => setCreateSupplierDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleCreateSupplier} 
            variant="contained"
            disabled={!supplierForm.name || !supplierForm.contactEmail || createSupplierMutation.isLoading}
          >
            {createSupplierMutation.isLoading ? <CircularProgress size={20} /> : 'Create Supplier'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Supplier Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Delete Supplier</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary">
            Are you sure you want to delete supplier "{selectedSupplier?.name}"? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleDeleteSupplier} 
            color="error" 
            variant="contained"
            disabled={deleteSupplierMutation.isLoading}
          >
            {deleteSupplierMutation.isLoading ? <CircularProgress size={20} /> : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Supplier Adapters Dialog */}
      <Dialog open={adapterDialogOpen} onClose={() => setAdapterDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Supplier Adapters - {selectedSupplier?.name}</DialogTitle>
        <DialogContent>
          {adaptersLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : adapters && adapters.length > 0 ? (
            <List>
              {adapters.map((adapter) => (
                <ListItem key={adapter.id} divider>
                  <ListItemText
                    primary={adapter.name || adapter.type}
                    secondary={
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          Type: {adapter.type} | Version: {adapter.version || 'N/A'}
                        </Typography>
                        {adapter.description && (
                          <Typography variant="body2" color="text.secondary">
                            {adapter.description}
                          </Typography>
                        )}
                      </Box>
                    }
                  />
                  <ListItemSecondaryAction>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      <Chip
                        label={adapter.status || 'Active'}
                        color={adapter.status === 'ACTIVE' ? 'success' : 'default'}
                        size="small"
                      />
                      <IconButton
                        size="small"
                        onClick={() => handleTestConnection(adapter.id)}
                        disabled={testConnectionMutation.isLoading}
                      >
                        <TestIcon />
                      </IconButton>
                    </Box>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>
          ) : (
            <Typography variant="body2" color="text.secondary" sx={{ py: 4, textAlign: 'center' }}>
              No adapters found for this supplier
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAdapterDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Test Connection Dialog */}
      <Dialog open={testConnectionDialogOpen} onClose={() => setTestConnectionDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Test Supplier Connection</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Test connection to supplier: {selectedSupplier?.name}
          </Typography>
          
          {testResults && (
            <Alert 
              severity={testResults.success ? 'success' : 'error'} 
              sx={{ mb: 2 }}
              icon={testResults.success ? <CheckIcon /> : <ErrorIcon />}
            >
              {testResults.success 
                ? 'Connection test successful!' 
                : `Connection test failed: ${testResults.error || 'Unknown error'}`
              }
            </Alert>
          )}

          {testConnectionMutation.isLoading && (
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, py: 2 }}>
              <CircularProgress size={20} />
              <Typography variant="body2">Testing connection...</Typography>
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setTestConnectionDialogOpen(false)}>Close</Button>
          <Button 
            onClick={() => handleTestConnection('default')} 
            variant="contained"
            disabled={testConnectionMutation.isLoading}
          >
            {testConnectionMutation.isLoading ? 'Testing...' : 'Test Connection'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default SuppliersPage;
