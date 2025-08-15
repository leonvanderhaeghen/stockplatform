import React, { useState } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Paper,
  Grid,
  TextField,
  Button,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Chip,
  Alert,
  Tabs,
  Tab,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  InputAdornment,
  CircularProgress,
  Skeleton,
} from '@mui/material';
import {
  Store,
  Search,
  Add,
  Edit,
  Delete,
  Sync,
  CheckCircle,
  Error,
  Warning,
  Refresh,
  CloudSync,
  Assessment,
  Business,
  Extension,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSnackbar } from 'notistack';
import { format } from 'date-fns';
import supplierService from '../../services/supplierService';

const SuppliersPage = () => {
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();

  // State for UI
  const [currentTab, setCurrentTab] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [selectedSupplier, setSelectedSupplier] = useState(null);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showEditDialog, setShowEditDialog] = useState(false);
  const [supplierFormData, setSupplierFormData] = useState({
    name: '', email: '', phone: '', address: '', contactPerson: '',
    website: '', notes: '', adapterName: ''
  });

  // Fetch suppliers
  const {
    data: suppliersData = { suppliers: [], total: 0 },
    isLoading: suppliersLoading,
    error: suppliersError
  } = useQuery({
    queryKey: ['suppliers', page + 1, rowsPerPage, searchTerm, statusFilter],
    queryFn: () => supplierService.getSuppliers({
      page: page + 1,
      limit: rowsPerPage,
      search: searchTerm || undefined,
      status: statusFilter !== 'all' ? statusFilter : undefined
    }),
  });

  // Fetch adapters
  const { data: adapters = [] } = useQuery({
    queryKey: ['supplier-adapters'],
    queryFn: () => supplierService.getAdapters(),
    enabled: currentTab === 1 || showCreateDialog || showEditDialog,
  });

  // Create supplier mutation
  const createSupplierMutation = useMutation({
    mutationFn: async (supplierData) => {
      return await supplierService.createSupplier(supplierData);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['suppliers']);
      setShowCreateDialog(false);
      resetForm();
      enqueueSnackbar('Supplier created successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to create supplier', { variant: 'error' });
    },
  });

  // Update supplier mutation
  const updateSupplierMutation = useMutation({
    mutationFn: async ({ id, ...data }) => {
      return await supplierService.updateSupplier(id, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['suppliers']);
      setShowEditDialog(false);
      setSelectedSupplier(null);
      resetForm();
      enqueueSnackbar('Supplier updated successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to update supplier', { variant: 'error' });
    },
  });

  // Delete supplier mutation
  const deleteSupplierMutation = useMutation({
    mutationFn: async (supplierId) => {
      return await supplierService.deleteSupplier(supplierId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['suppliers']);
      enqueueSnackbar('Supplier deleted successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to delete supplier', { variant: 'error' });
    },
  });

  const resetForm = () => {
    setSupplierFormData({
      name: '', email: '', phone: '', address: '', contactPerson: '',
      website: '', notes: '', adapterName: ''
    });
  };

  const handleCreateSupplier = () => {
    if (!supplierFormData.name || !supplierFormData.email) {
      enqueueSnackbar('Please fill in required fields (name and email)', { variant: 'warning' });
      return;
    }
    createSupplierMutation.mutate(supplierFormData);
  };

  const handleUpdateSupplier = () => {
    if (!selectedSupplier || !supplierFormData.name || !supplierFormData.email) {
      enqueueSnackbar('Please fill in required fields (name and email)', { variant: 'warning' });
      return;
    }
    updateSupplierMutation.mutate({ id: selectedSupplier.id, ...supplierFormData });
  };

  const handleEditSupplier = (supplier) => {
    setSelectedSupplier(supplier);
    setSupplierFormData({
      name: supplier.name || '',
      email: supplier.email || '',
      phone: supplier.phone || '',
      address: supplier.address || '',
      contactPerson: supplier.contactPerson || '',
      website: supplier.website || '',
      notes: supplier.notes || '',
      adapterName: supplier.adapterName || ''
    });
    setShowEditDialog(true);
  };

  const getStatusChip = (supplier) => {
    const isActive = supplier.status === 'active';
    return (
      <Chip
        label={isActive ? 'Active' : 'Inactive'}
        color={isActive ? 'success' : 'default'}
        size="small"
      />
    );
  };

  const getAdapterStatusIcon = (supplier) => {
    if (!supplier.adapterName) return <Warning color="warning" />;
    if (supplier.adapterStatus === 'connected') return <CheckCircle color="success" />;
    return <Error color="error" />;
  };

  const TabPanel = ({ children, value, index }) => {
    return value === index ? <Box sx={{ mt: 3 }}>{children}</Box> : null;
  };

  return (
    <Box>
      {/* Header */}
      <Paper sx={{ p: 3, mb: 3, bgcolor: 'primary.main', color: 'primary.contrastText' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Store sx={{ fontSize: 40, mr: 2 }} />
            <Box>
              <Typography variant="h4" gutterBottom>
                Supplier Management
              </Typography>
              <Typography variant="subtitle1">
                Manage suppliers, adapters, and product synchronization
              </Typography>
            </Box>
          </Box>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              variant="contained"
              color="secondary"
              startIcon={<Add />}
              onClick={() => setShowCreateDialog(true)}
            >
              Add Supplier
            </Button>
            <Button
              variant="contained"
              color="secondary"
              startIcon={<Refresh />}
              onClick={() => queryClient.invalidateQueries(['suppliers'])}
            >
              Refresh
            </Button>
          </Box>
        </Box>
      </Paper>

      <Card>
        <CardContent>
          {/* Tabs */}
          <Tabs value={currentTab} onChange={(e, newValue) => setCurrentTab(newValue)}>
            <Tab label="Suppliers" icon={<Business />} />
            <Tab label="Adapters" icon={<Extension />} />
            <Tab label="Analytics" icon={<Assessment />} />
          </Tabs>

          {/* Suppliers Tab */}
          <TabPanel value={currentTab} index={0}>
            {/* Filters and Search */}
            <Grid container spacing={2} sx={{ mb: 3 }}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  placeholder="Search suppliers by name, email, or contact..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <Search />
                      </InputAdornment>
                    ),
                  }}
                />
              </Grid>
              <Grid item xs={12} md={4}>
                <FormControl fullWidth>
                  <InputLabel>Status Filter</InputLabel>
                  <Select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value)}
                    label="Status Filter"
                  >
                    <MenuItem value="all">All Suppliers</MenuItem>
                    <MenuItem value="active">Active</MenuItem>
                    <MenuItem value="inactive">Inactive</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={2}>
                <Button
                  fullWidth
                  variant="outlined"
                  startIcon={<CloudSync />}
                  sx={{ height: '56px' }}
                >
                  Sync All
                </Button>
              </Grid>
            </Grid>

            {/* Suppliers Table */}
            {suppliersError ? (
              <Alert severity="error" sx={{ mb: 2 }}>
                Failed to load suppliers: {suppliersError.message}
              </Alert>
            ) : (
              <>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Supplier</TableCell>
                        <TableCell>Contact</TableCell>
                        <TableCell align="center">Adapter</TableCell>
                        <TableCell align="center">Status</TableCell>
                        <TableCell align="center">Last Sync</TableCell>
                        <TableCell align="center">Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {suppliersLoading ? (
                        // Loading skeletons
                        [...Array(rowsPerPage)].map((_, index) => (
                          <TableRow key={index}>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="rectangular" width={40} height={24} /></TableCell>
                            <TableCell><Skeleton variant="rectangular" width={60} height={24} /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="rectangular" width={120} height={36} /></TableCell>
                          </TableRow>
                        ))
                      ) : suppliersData.suppliers.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={6} align="center">
                            <Typography color="text.secondary" sx={{ py: 4 }}>
                              No suppliers found. {searchTerm && 'Try adjusting your search terms.'}
                            </Typography>
                          </TableCell>
                        </TableRow>
                      ) : (
                        suppliersData.suppliers.map((supplier) => (
                          <TableRow key={supplier.id}>
                            <TableCell>
                              <Box>
                                <Typography variant="body1" fontWeight="medium">
                                  {supplier.name}
                                </Typography>
                                {supplier.website && (
                                  <Typography variant="caption" color="primary">
                                    {supplier.website}
                                  </Typography>
                                )}
                              </Box>
                            </TableCell>
                            <TableCell>
                              <Box>
                                <Typography variant="body2">
                                  {supplier.contactPerson || 'N/A'}
                                </Typography>
                                <Typography variant="caption" color="text.secondary">
                                  {supplier.email}
                                </Typography>
                                {supplier.phone && (
                                  <Typography variant="caption" color="text.secondary" display="block">
                                    {supplier.phone}
                                  </Typography>
                                )}
                              </Box>
                            </TableCell>
                            <TableCell align="center">
                              <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                                {getAdapterStatusIcon(supplier)}
                                {supplier.adapterName && (
                                  <Typography variant="caption" sx={{ ml: 1 }}>
                                    {supplier.adapterName}
                                  </Typography>
                                )}
                              </Box>
                            </TableCell>
                            <TableCell align="center">
                              {getStatusChip(supplier)}
                            </TableCell>
                            <TableCell align="center">
                              <Typography variant="caption" color="text.secondary">
                                {supplier.lastSyncAt ? format(new Date(supplier.lastSyncAt), 'MMM dd, yyyy') : 'Never'}
                              </Typography>
                            </TableCell>
                            <TableCell align="center">
                              <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center' }}>
                                <IconButton
                                  size="small"
                                  onClick={() => handleEditSupplier(supplier)}
                                  title="Edit Supplier"
                                >
                                  <Edit />
                                </IconButton>
                                <IconButton
                                  size="small"
                                  title="Sync Products"
                                  disabled={!supplier.adapterName}
                                >
                                  <Sync />
                                </IconButton>
                                <IconButton
                                  size="small"
                                  onClick={() => deleteSupplierMutation.mutate(supplier.id)}
                                  title="Delete Supplier"
                                  color="error"
                                >
                                  <Delete />
                                </IconButton>
                              </Box>
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
                  count={suppliersData.total || 0}
                  rowsPerPage={rowsPerPage}
                  page={page}
                  onPageChange={(event, newPage) => setPage(newPage)}
                  onRowsPerPageChange={(event) => {
                    setRowsPerPage(parseInt(event.target.value, 10));
                    setPage(0);
                  }}
                />
              </>
            )}
          </TabPanel>

          {/* Adapters Tab */}
          <TabPanel value={currentTab} index={1}>
            <Typography variant="h6" gutterBottom>
              Available Supplier Adapters
            </Typography>
            <Grid container spacing={2}>
              {adapters.map((adapter) => (
                <Grid item xs={12} md={6} lg={4} key={adapter.name}>
                  <Card variant="outlined">
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        {adapter.displayName || adapter.name}
                      </Typography>
                      <Typography variant="body2" color="text.secondary" gutterBottom>
                        {adapter.description}
                      </Typography>
                      <Box sx={{ mt: 2 }}>
                        <Typography variant="caption" display="block">
                          Version: {adapter.version}
                        </Typography>
                        <Typography variant="caption" display="block">
                          Type: {adapter.type}
                        </Typography>
                      </Box>
                      <Button
                        variant="outlined"
                        size="small"
                        sx={{ mt: 2 }}
                      >
                        Test Connection
                      </Button>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
            {adapters.length === 0 && (
              <Typography color="text.secondary" sx={{ textAlign: 'center', py: 4 }}>
                No adapters available.
              </Typography>
            )}
          </TabPanel>

          {/* Analytics Tab */}
          <TabPanel value={currentTab} index={2}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Supplier Summary
                    </Typography>
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Total Suppliers:</Typography>
                        <Typography fontWeight="bold">{suppliersData.total || 0}</Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Active Suppliers:</Typography>
                        <Typography fontWeight="bold" color="success.main">
                          {suppliersData.suppliers?.filter(s => s.status === 'active').length || 0}
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>With Adapters:</Typography>
                        <Typography fontWeight="bold" color="info.main">
                          {suppliersData.suppliers?.filter(s => s.adapterName).length || 0}
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Connected:</Typography>
                        <Typography fontWeight="bold" color="primary.main">
                          {suppliersData.suppliers?.filter(s => s.adapterStatus === 'connected').length || 0}
                        </Typography>
                      </Box>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Sync Activity
                    </Typography>
                    <Typography color="text.secondary">
                      Synchronization analytics and recent activity will be displayed here.
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>
        </CardContent>
      </Card>

      {/* Create/Edit Supplier Dialog */}
      <Dialog
        open={showCreateDialog || showEditDialog}
        onClose={() => {
          setShowCreateDialog(false);
          setShowEditDialog(false);
          resetForm();
          setSelectedSupplier(null);
        }}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          {showCreateDialog ? 'Create New Supplier' : `Edit Supplier - ${selectedSupplier?.name}`}
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} md={6}>
              <TextField
                label="Supplier Name"
                value={supplierFormData.name}
                onChange={(e) => setSupplierFormData({ ...supplierFormData, name: e.target.value })}
                fullWidth
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Email"
                type="email"
                value={supplierFormData.email}
                onChange={(e) => setSupplierFormData({ ...supplierFormData, email: e.target.value })}
                fullWidth
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Phone"
                value={supplierFormData.phone}
                onChange={(e) => setSupplierFormData({ ...supplierFormData, phone: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Contact Person"
                value={supplierFormData.contactPerson}
                onChange={(e) => setSupplierFormData({ ...supplierFormData, contactPerson: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Address"
                value={supplierFormData.address}
                onChange={(e) => setSupplierFormData({ ...supplierFormData, address: e.target.value })}
                fullWidth
                multiline
                rows={2}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Website"
                value={supplierFormData.website}
                onChange={(e) => setSupplierFormData({ ...supplierFormData, website: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Adapter</InputLabel>
                <Select
                  value={supplierFormData.adapterName}
                  onChange={(e) => setSupplierFormData({ ...supplierFormData, adapterName: e.target.value })}
                  label="Adapter"
                >
                  <MenuItem value="">No Adapter</MenuItem>
                  {adapters.map((adapter) => (
                    <MenuItem key={adapter.name} value={adapter.name}>
                      {adapter.displayName || adapter.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Notes"
                value={supplierFormData.notes}
                onChange={(e) => setSupplierFormData({ ...supplierFormData, notes: e.target.value })}
                fullWidth
                multiline
                rows={3}
                placeholder="Additional notes about this supplier..."
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setShowCreateDialog(false);
              setShowEditDialog(false);
              resetForm();
              setSelectedSupplier(null);
            }}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            onClick={showCreateDialog ? handleCreateSupplier : handleUpdateSupplier}
            disabled={createSupplierMutation.isLoading || updateSupplierMutation.isLoading}
            startIcon={(createSupplierMutation.isLoading || updateSupplierMutation.isLoading) ? <CircularProgress size={20} /> : <Add />}
          >
            {createSupplierMutation.isLoading || updateSupplierMutation.isLoading
              ? 'Processing...'
              : showCreateDialog ? 'Create Supplier' : 'Update Supplier'
            }
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default SuppliersPage;
