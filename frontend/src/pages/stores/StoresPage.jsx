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
  Switch,
  FormControlLabel,
  Divider,
} from '@mui/material';
import {
  Store,
  Search,
  Add,
  Edit,
  Delete,
  LocationOn,
  Phone,
  Email,
  Schedule,
  Settings,
  Refresh,
  Assessment,
  Storefront,
  Map,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSnackbar } from 'notistack';
import storeService from '../../services/storeService';

const StoresPage = () => {
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();

  // State for UI
  const [currentTab, setCurrentTab] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [selectedStore, setSelectedStore] = useState(null);
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showEditDialog, setShowEditDialog] = useState(false);
  const [showSettingsDialog, setShowSettingsDialog] = useState(false);
  const [storeFormData, setStoreFormData] = useState({
    name: '', address: '', city: '', state: '', country: '', zipCode: '',
    phone: '', email: '', managerName: '', description: '',
    operatingHours: { open: '09:00', close: '18:00' },
    settings: { allowPOS: true, allowPickup: true, allowDelivery: false }
  });

  // Fetch stores
  const {
    data: storesData = { stores: [], total: 0 },
    isLoading: storesLoading,
    error: storesError
  } = useQuery({
    queryKey: ['stores', page + 1, rowsPerPage, searchTerm, statusFilter],
    queryFn: () => storeService.getStores({
      page: page + 1,
      limit: rowsPerPage,
      search: searchTerm || undefined,
      status: statusFilter !== 'all' ? statusFilter : undefined
    }),
  });

  // Create store mutation
  const createStoreMutation = useMutation({
    mutationFn: async (storeData) => {
      return await storeService.createStore(storeData);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['stores']);
      setShowCreateDialog(false);
      resetForm();
      enqueueSnackbar('Store created successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to create store', { variant: 'error' });
    },
  });

  // Update store mutation
  const updateStoreMutation = useMutation({
    mutationFn: async ({ id, ...data }) => {
      return await storeService.updateStore(id, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['stores']);
      setShowEditDialog(false);
      setShowSettingsDialog(false);
      setSelectedStore(null);
      resetForm();
      enqueueSnackbar('Store updated successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to update store', { variant: 'error' });
    },
  });

  // Delete store mutation
  const deleteStoreMutation = useMutation({
    mutationFn: async (storeId) => {
      return await storeService.deleteStore(storeId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['stores']);
      enqueueSnackbar('Store deleted successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to delete store', { variant: 'error' });
    },
  });

  const resetForm = () => {
    setStoreFormData({
      name: '', address: '', city: '', state: '', country: '', zipCode: '',
      phone: '', email: '', managerName: '', description: '',
      operatingHours: { open: '09:00', close: '18:00' },
      settings: { allowPOS: true, allowPickup: true, allowDelivery: false }
    });
  };

  const handleCreateStore = () => {
    if (!storeFormData.name || !storeFormData.address || !storeFormData.city) {
      enqueueSnackbar('Please fill in required fields (name, address, city)', { variant: 'warning' });
      return;
    }
    createStoreMutation.mutate(storeFormData);
  };

  const handleUpdateStore = () => {
    if (!selectedStore || !storeFormData.name || !storeFormData.address || !storeFormData.city) {
      enqueueSnackbar('Please fill in required fields (name, address, city)', { variant: 'warning' });
      return;
    }
    updateStoreMutation.mutate({ id: selectedStore.id, ...storeFormData });
  };

  const handleEditStore = (store) => {
    setSelectedStore(store);
    setStoreFormData({
      name: store.name || '',
      address: store.address || '',
      city: store.city || '',
      state: store.state || '',
      country: store.country || '',
      zipCode: store.zipCode || '',
      phone: store.phone || '',
      email: store.email || '',
      managerName: store.managerName || '',
      description: store.description || '',
      operatingHours: store.operatingHours || { open: '09:00', close: '18:00' },
      settings: store.settings || { allowPOS: true, allowPickup: true, allowDelivery: false }
    });
    setShowEditDialog(true);
  };

  const handleStoreSettings = (store) => {
    setSelectedStore(store);
    setStoreFormData({
      ...storeFormData,
      settings: store.settings || { allowPOS: true, allowPickup: true, allowDelivery: false }
    });
    setShowSettingsDialog(true);
  };

  const getStatusChip = (store) => {
    const isActive = store.status === 'active';
    return (
      <Chip
        label={isActive ? 'Active' : 'Inactive'}
        color={isActive ? 'success' : 'default'}
        size="small"
      />
    );
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
                Store Management
              </Typography>
              <Typography variant="subtitle1">
                Manage store locations, inventory, and operations
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
              Add Store
            </Button>
            <Button
              variant="contained"
              color="secondary"
              startIcon={<Refresh />}
              onClick={() => queryClient.invalidateQueries(['stores'])}
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
            <Tab label="Stores" icon={<Storefront />} />
            <Tab label="Analytics" icon={<Assessment />} />
            <Tab label="Locations" icon={<Map />} />
          </Tabs>

          {/* Stores Tab */}
          <TabPanel value={currentTab} index={0}>
            {/* Filters and Search */}
            <Grid container spacing={2} sx={{ mb: 3 }}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  placeholder="Search stores by name, city, or manager..."
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
                    <MenuItem value="all">All Stores</MenuItem>
                    <MenuItem value="active">Active</MenuItem>
                    <MenuItem value="inactive">Inactive</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
            </Grid>

            {/* Stores Table */}
            {storesError ? (
              <Alert severity="error" sx={{ mb: 2 }}>
                Failed to load stores: {storesError.message}
              </Alert>
            ) : (
              <>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Store</TableCell>
                        <TableCell>Location</TableCell>
                        <TableCell>Contact</TableCell>
                        <TableCell align="center">Hours</TableCell>
                        <TableCell align="center">Status</TableCell>
                        <TableCell align="center">Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {storesLoading ? (
                        // Loading skeletons
                        [...Array(rowsPerPage)].map((_, index) => (
                          <TableRow key={index}>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="rectangular" width={60} height={24} /></TableCell>
                            <TableCell><Skeleton variant="rectangular" width={120} height={36} /></TableCell>
                          </TableRow>
                        ))
                      ) : storesData.stores.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={6} align="center">
                            <Typography color="text.secondary" sx={{ py: 4 }}>
                              No stores found. {searchTerm && 'Try adjusting your search terms.'}
                            </Typography>
                          </TableCell>
                        </TableRow>
                      ) : (
                        storesData.stores.map((store) => (
                          <TableRow key={store.id}>
                            <TableCell>
                              <Box>
                                <Typography variant="body1" fontWeight="medium">
                                  {store.name}
                                </Typography>
                                <Typography variant="caption" color="text.secondary">
                                  Manager: {store.managerName || 'N/A'}
                                </Typography>
                              </Box>
                            </TableCell>
                            <TableCell>
                              <Box>
                                <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center' }}>
                                  <LocationOn sx={{ fontSize: 16, mr: 0.5 }} />
                                  {store.address?.street || 'No address'}
                                </Typography>
                                <Typography variant="caption" color="text.secondary">
                                  {store.address?.city && `${store.address.city}, `}
                                  {store.address?.state && `${store.address.state} `}
                                  {store.address?.zip_code}
                                </Typography>
                                {store.address?.country && (
                                  <Typography variant="caption" color="text.secondary" display="block">
                                    {store.address.country}
                                  </Typography>
                                )}
                              </Box>
                            </TableCell>
                            <TableCell>
                              <Box>
                                {store.phone && (
                                  <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center' }}>
                                    <Phone sx={{ fontSize: 16, mr: 0.5 }} />
                                    {store.phone}
                                  </Typography>
                                )}
                                {store.email && (
                                  <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center' }}>
                                    <Email sx={{ fontSize: 16, mr: 0.5 }} />
                                    {store.email}
                                  </Typography>
                                )}
                              </Box>
                            </TableCell>
                            <TableCell align="center">
                              <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                                <Schedule sx={{ fontSize: 16, mr: 0.5 }} />
                                {store.operatingHours ? 
                                  `${store.operatingHours.open} - ${store.operatingHours.close}` : 
                                  'Not set'
                                }
                              </Typography>
                            </TableCell>
                            <TableCell align="center">
                              {getStatusChip(store)}
                            </TableCell>
                            <TableCell align="center">
                              <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center' }}>
                                <IconButton
                                  size="small"
                                  onClick={() => handleEditStore(store)}
                                  title="Edit Store"
                                >
                                  <Edit />
                                </IconButton>
                                <IconButton
                                  size="small"
                                  onClick={() => handleStoreSettings(store)}
                                  title="Store Settings"
                                >
                                  <Settings />
                                </IconButton>
                                <IconButton
                                  size="small"
                                  onClick={() => deleteStoreMutation.mutate(store.id)}
                                  title="Delete Store"
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
                  count={storesData.total || 0}
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

          {/* Analytics Tab */}
          <TabPanel value={currentTab} index={1}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Store Summary
                    </Typography>
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Total Stores:</Typography>
                        <Typography fontWeight="bold">{storesData.total || 0}</Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Active Stores:</Typography>
                        <Typography fontWeight="bold" color="success.main">
                          {storesData.stores?.filter(s => s.status === 'active').length || 0}
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>POS Enabled:</Typography>
                        <Typography fontWeight="bold" color="primary.main">
                          {storesData.stores?.filter(s => s.settings?.allowPOS).length || 0}
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Pickup Enabled:</Typography>
                        <Typography fontWeight="bold" color="info.main">
                          {storesData.stores?.filter(s => s.settings?.allowPickup).length || 0}
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
                      Store Performance
                    </Typography>
                    <Typography color="text.secondary">
                      Store performance metrics and analytics will be displayed here.
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>

          {/* Locations Tab */}
          <TabPanel value={currentTab} index={2}>
            <Typography variant="h6" gutterBottom>
              Store Locations Map
            </Typography>
            <Card variant="outlined" sx={{ height: 400, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <Box sx={{ textAlign: 'center' }}>
                <Map sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
                <Typography color="text.secondary">
                  Interactive map showing all store locations will be displayed here.
                </Typography>
              </Box>
            </Card>
          </TabPanel>
        </CardContent>
      </Card>

      {/* Create/Edit Store Dialog */}
      <Dialog
        open={showCreateDialog || showEditDialog}
        onClose={() => {
          setShowCreateDialog(false);
          setShowEditDialog(false);
          resetForm();
          setSelectedStore(null);
        }}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          {showCreateDialog ? 'Create New Store' : `Edit Store - ${selectedStore?.name}`}
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} md={6}>
              <TextField
                label="Store Name"
                value={storeFormData.name}
                onChange={(e) => setStoreFormData({ ...storeFormData, name: e.target.value })}
                fullWidth
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Manager Name"
                value={storeFormData.managerName}
                onChange={(e) => setStoreFormData({ ...storeFormData, managerName: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Address"
                value={storeFormData.address}
                onChange={(e) => setStoreFormData({ ...storeFormData, address: e.target.value })}
                fullWidth
                required
              />
            </Grid>
            <Grid item xs={12} md={4}>
              <TextField
                label="City"
                value={storeFormData.city}
                onChange={(e) => setStoreFormData({ ...storeFormData, city: e.target.value })}
                fullWidth
                required
              />
            </Grid>
            <Grid item xs={12} md={4}>
              <TextField
                label="State/Province"
                value={storeFormData.state}
                onChange={(e) => setStoreFormData({ ...storeFormData, state: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={4}>
              <TextField
                label="ZIP/Postal Code"
                value={storeFormData.zipCode}
                onChange={(e) => setStoreFormData({ ...storeFormData, zipCode: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Country"
                value={storeFormData.country}
                onChange={(e) => setStoreFormData({ ...storeFormData, country: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Phone"
                value={storeFormData.phone}
                onChange={(e) => setStoreFormData({ ...storeFormData, phone: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                label="Email"
                type="email"
                value={storeFormData.email}
                onChange={(e) => setStoreFormData({ ...storeFormData, email: e.target.value })}
                fullWidth
              />
            </Grid>
            <Grid item xs={12} md={3}>
              <TextField
                label="Opening Time"
                type="time"
                value={storeFormData.operatingHours.open}
                onChange={(e) => setStoreFormData({ 
                  ...storeFormData, 
                  operatingHours: { ...storeFormData.operatingHours, open: e.target.value }
                })}
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} md={3}>
              <TextField
                label="Closing Time"
                type="time"
                value={storeFormData.operatingHours.close}
                onChange={(e) => setStoreFormData({ 
                  ...storeFormData, 
                  operatingHours: { ...storeFormData.operatingHours, close: e.target.value }
                })}
                fullWidth
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                label="Description"
                value={storeFormData.description}
                onChange={(e) => setStoreFormData({ ...storeFormData, description: e.target.value })}
                fullWidth
                multiline
                rows={3}
                placeholder="Additional notes about this store..."
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
              setSelectedStore(null);
            }}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            onClick={showCreateDialog ? handleCreateStore : handleUpdateStore}
            disabled={createStoreMutation.isLoading || updateStoreMutation.isLoading}
            startIcon={(createStoreMutation.isLoading || updateStoreMutation.isLoading) ? <CircularProgress size={20} /> : <Add />}
          >
            {createStoreMutation.isLoading || updateStoreMutation.isLoading
              ? 'Processing...'
              : showCreateDialog ? 'Create Store' : 'Update Store'
            }
          </Button>
        </DialogActions>
      </Dialog>

      {/* Store Settings Dialog */}
      <Dialog open={showSettingsDialog} onClose={() => setShowSettingsDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          Store Settings - {selectedStore?.name}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography variant="h6" gutterBottom>
              Operation Settings
            </Typography>
            <Divider sx={{ mb: 2 }} />
            
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              <FormControlLabel
                control={
                  <Switch
                    checked={storeFormData.settings.allowPOS}
                    onChange={(e) => setStoreFormData({
                      ...storeFormData,
                      settings: { ...storeFormData.settings, allowPOS: e.target.checked }
                    })}
                  />
                }
                label="Enable POS Operations"
              />
              <FormControlLabel
                control={
                  <Switch
                    checked={storeFormData.settings.allowPickup}
                    onChange={(e) => setStoreFormData({
                      ...storeFormData,
                      settings: { ...storeFormData.settings, allowPickup: e.target.checked }
                    })}
                  />
                }
                label="Enable Customer Pickup"
              />
              <FormControlLabel
                control={
                  <Switch
                    checked={storeFormData.settings.allowDelivery}
                    onChange={(e) => setStoreFormData({
                      ...storeFormData,
                      settings: { ...storeFormData.settings, allowDelivery: e.target.checked }
                    })}
                  />
                }
                label="Enable Delivery Service"
              />
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowSettingsDialog(false)}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handleUpdateStore}
            disabled={updateStoreMutation.isLoading}
            startIcon={updateStoreMutation.isLoading ? <CircularProgress size={20} /> : <Settings />}
          >
            {updateStoreMutation.isLoading ? 'Saving...' : 'Save Settings'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default StoresPage;
