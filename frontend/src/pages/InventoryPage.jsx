import React, { useState, useEffect } from 'react';
import Autocomplete from '@mui/material/Autocomplete';
import {
  Box,
  Typography,
  Button,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  CircularProgress,
  Tabs,
  Tab,
  Card,
  CardContent,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  Divider,
  FormControl,
  InputLabel,
  Select,
  Alert,
  AlertTitle,
  Grid,
  List,
  ListItem,
  ListItemText,
  Menu
} from '@mui/material';
import {
  Inventory as InventoryIcon,
  Add as AddIcon,
  Lock as ReserveIcon,
  LockOpen as ReleaseIcon,
  Analytics as AnalyticsIcon,
  Warning as WarningIcon,
  CheckCircle as CheckIcon,
  MoreVert as MoreVertIcon,
  Refresh as RefreshIcon,
  History as HistoryIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { inventoryService, productService, storeService } from '../services';
import { formatCurrency, formatDate } from '../utils/formatters';
import {
  useInventoryItems,
  useInventoryReservations,
  useLowStockItems,
  useInventoryHistory,
  useUpdateInventoryQuantity,
  useReserveInventory,
  useReleaseInventory,
  useCreateInventoryItem,
  useDeleteInventoryItem,
} from '../utils/hooks/useInventory';
import InventoryStatus from '../components/inventory/InventoryStatus';
import InventoryErrorBoundary from '../components/inventory/InventoryErrorBoundary';

const InventoryPage = () => {
  const queryClient = useQueryClient();
  const [tabValue, setTabValue] = useState(0);
  // Removed unused selectedProduct state
  const [anchorEl, setAnchorEl] = useState(null);
  const [reservationDialogOpen, setReservationDialogOpen] = useState(false);
  const [releaseDialogOpen, setReleaseDialogOpen] = useState(false);
  const [adjustmentDialogOpen, setAdjustmentDialogOpen] = useState(false);
  const [historyDialogOpen, setHistoryDialogOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [products, setProducts] = useState([]);
  const [totalProducts, setTotalProducts] = useState(0);
  
  // Location filtering state
  const [selectedLocation, setSelectedLocation] = useState('');
  const [locationOptions, setLocationOptions] = useState([]);
  const [locationInputValue, setLocationInputValue] = useState('');
  const [lowStockThreshold, setLowStockThreshold] = useState(10);
  const [reservationData, setReservationData] = useState({
    quantity: 1,
    orderId: '',
    notes: '',
  });
  const [releaseData, setReleaseData] = useState({
    reservationId: '',
    quantity: 1,
    reason: '',
  });
  const [adjustmentData, setAdjustmentData] = useState({
    quantity: 0,
    type: 'ADD',
    reason: '',
    notes: '',
  });
  const [createItemDialogOpen, setCreateItemDialogOpen] = useState(false);
  const [removeItemDialogOpen, setRemoveItemDialogOpen] = useState(false);
  const [selectedInventory, setSelectedInventory] = useState(null);
  
  // Product lookup and store management state
  const [productLookupLoading, setProductLookupLoading] = useState(false);
  const [productOptions, setProductOptions] = useState([]);
  const [productLookupValue, setProductLookupValue] = useState('');
  const [storeOptions, setStoreOptions] = useState([]);
  const [storeDialogOpen, setStoreDialogOpen] = useState(false);
  const [newStore, setNewStore] = useState({ name: '', address: '', phone: '', email: '' });

  // Fetch products and inventory data using new hooks
  const { data: inventoryItems = [], isLoading: inventoryLoading } = useInventoryItems();
  
  // Fetch products for analytics
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const data = await productService.getProducts();
        setProducts(data.products || []);
        setTotalProducts(data.total || 0);
      } catch (error) {
        console.error('Failed to fetch products:', error);
        setError('Failed to load product data');
      }
    };
    
    if (tabValue === 3) { // Only fetch when analytics tab is active
      fetchProducts();
    }
  }, [tabValue]);

  // Fetch inventory reservations using new hook
  const { data: reservationsData, isLoading: reservationsLoading } = useInventoryReservations(
    tabValue === 1 ? {} : undefined
  );

  // Fetch low stock alerts using new hook with location filtering
  const { data: lowStockData, isLoading: lowStockLoading } = useLowStockItems(
    tabValue === 2 ? lowStockThreshold : undefined,
    tabValue === 2 ? selectedLocation : undefined
  );

  // Fetch inventory history for selected product using new hook
  const { data: historyData, isLoading: historyLoading } = useInventoryHistory(
    historyDialogOpen && selectedInventory?.id ? selectedInventory.id : null
  );

  // Use new inventory hooks for mutations
  const updateQuantityMutation = useUpdateInventoryQuantity({
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory']);
      setSuccess('Inventory quantity updated successfully');
    },
    onError: (error) => {
      setError(error.message || 'Failed to update inventory quantity');
    }
  });
  
  const reserveInventoryMutation = useReserveInventory({
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory', 'reservations']);
      setReservationDialogOpen(false);
      setSuccess('Inventory reserved successfully');
    },
    onError: (error) => {
      setError(error.message || 'Failed to reserve inventory');
    }
  });
  
  const releaseInventoryMutation = useReleaseInventory({
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory', 'reservations']);
      setReleaseDialogOpen(false);
      setSuccess('Inventory released successfully');
    },
    onError: (error) => {
      setError(error.message || 'Failed to release inventory');
    }
  });
  
  const createItemMutation = useCreateInventoryItem({
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory']);
      setCreateItemDialogOpen(false);
      setSuccess('Inventory item created successfully');
    },
    onError: (error) => {
      setError(error.message || 'Failed to create inventory item');
    }
  });
  
  const removeItemMutation = useDeleteInventoryItem({
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory']);
      setRemoveItemDialogOpen(false);
      setSuccess('Inventory item removed successfully');
    },
    onError: (error) => {
      setError(error.message || 'Failed to remove inventory item');
    }
  });

  // Fetch location autocomplete options when low-stock tab is active
  useEffect(() => {
    const fetchLocationOptions = async () => {
      if (tabValue === 2 && locationInputValue.trim()) {
        try {
          const stores = await storeService.getStoresForAutocomplete(locationInputValue);
          setLocationOptions(stores);
        } catch (error) {
          console.error('Failed to fetch location options:', error);
          setLocationOptions([]);
        }
      } else {
        setLocationOptions([]);
      }
    };

    fetchLocationOptions();
  }, [tabValue, locationInputValue]);

  // Fetch store options for inventory creation
  useEffect(() => {
    const fetchStoreOptions = async () => {
      if (createItemDialogOpen) {
        try {
          const stores = await storeService.getStores();
          setStoreOptions(stores.map(store => ({
            value: store.id,
            label: `${store.name} - ${store.address}`,
            store: store
          })));
        } catch (error) {
          console.error('Failed to fetch store options:', error);
          setStoreOptions([]);
        }
      }
    };

    fetchStoreOptions();
  }, [createItemDialogOpen]);

  // Store creation function
  const handleCreateStore = async () => {
    try {
      const createdStore = await storeService.createStore(newStore);
      setStoreOptions(prev => [...prev, {
        value: createdStore.id,
        label: `${createdStore.name} - ${createdStore.address}`,
        store: createdStore
      }]);
      setSelectedInventory({
        ...selectedInventory,
        storeId: createdStore.id
      });
      setStoreDialogOpen(false);
      setNewStore({ name: '', address: '', phone: '', email: '' });
      setSuccess(`Store created: ${createdStore.name}`);
    } catch (error) {
      console.error('Store creation failed:', error);
      setError('Failed to create store. Please try again.');
    }
  };

  // Check stock mutation (keeping original for compatibility)
  const checkStockMutation = useMutation({
    mutationFn: (productId) => inventoryService.checkStock(productId),
    onSuccess: (data) => {
      setSuccess(`Stock check: ${data.available} units available`);
      setError('');
    },
    onError: (error) => {
      setError(`Error checking stock: ${error.message}`);
    },
  });

  const reservations = reservationsData?.data || [];
  const lowStockItems = lowStockData?.data || [];

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleMenuOpen = (event, product) => {
    setSelectedProduct(product);
    setSelectedInventory(product);
    setAnchorEl(event.currentTarget);
  };

  const handleCheckStock = () => {
    if (selectedInventory) {
      checkStockMutation.mutate(selectedInventory.id);
    }
  };

  const handleReserveInventory = () => {
    if (selectedInventory && reservationData.quantity > 0) {
      reserveInventoryMutation.mutate({
        inventoryId: selectedInventory.id,
        quantity: reservationData.quantity,
        orderId: reservationData.orderId,
        notes: reservationData.notes
      });
      setReservationDialogOpen(false);
      resetReservationData();
    }
  };

  const handleReleaseInventory = () => {
    if (selectedInventory && releaseData.quantity > 0) {
      releaseInventoryMutation.mutate({
        reservationId: releaseData.reservationId,
        quantity: releaseData.quantity,
        reason: releaseData.reason
      });
      setReleaseDialogOpen(false);
      resetReleaseData();
    }
  };

  const handleAdjustInventory = () => {
    if (selectedInventory && adjustmentData.quantity !== 0) {
      updateQuantityMutation.mutate({
        id: selectedInventory.id,
        change: adjustmentData.type === 'ADD' ? adjustmentData.quantity : -adjustmentData.quantity,
        reason: adjustmentData.reason || 'Manual adjustment'
      });
      setAdjustmentDialogOpen(false);
      resetAdjustmentData();
    }
  };

  const resetReservationData = () => {
    setReservationData({
      quantity: 1,
      orderId: '',
      notes: '',
    });
  };

  const resetReleaseData = () => {
    setReleaseData({
      reservationId: '',
      quantity: 1,
      reason: '',
    });
  };

  const resetAdjustmentData = () => {
    setAdjustmentData({
      quantity: 0,
      type: 'ADD',
      reason: '',
      notes: '',
    });
  };

  const getStockStatusColor = (stock, lowStockThreshold = 10) => {
    if (stock === 0) return 'error';
    if (stock <= lowStockThreshold) return 'warning';
    return 'success';
  };

  const getStockStatusText = (stock, lowStockThreshold = 10) => {
    if (stock === 0) return 'Out of Stock';
    if (stock <= lowStockThreshold) return 'Low Stock';
    return 'In Stock';
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
          Inventory Management
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Button 
            variant="outlined" 
            startIcon={<AddIcon />}
            onClick={() => {
              setSelectedInventory({ name: '', sku: '', quantity: 0 });
              setCreateItemDialogOpen(true);
            }}
            sx={{ mr: 1 }}
          >
            Create Item
          </Button>
          <Button 
            variant="outlined" 
            startIcon={<RefreshIcon />}
            onClick={() => {
              queryClient.invalidateQueries(['inventory-reservations']);
              queryClient.invalidateQueries(['inventory-low-stock']);
            }}
          >
            Refresh
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

      <Tabs value={tabValue} onChange={handleTabChange} sx={{ mb: 3 }}>
        <Tab icon={<InventoryIcon />} label="Stock Overview" />
        <Tab icon={<ReserveIcon />} label="Reservations" />
        <Tab icon={<WarningIcon />} label="Low Stock Alerts" />
        <Tab icon={<AnalyticsIcon />} label="Analytics" />
      </Tabs>

      {/* Stock Overview Tab */}
      {tabValue === 0 && (
  <Paper sx={{ width: '100%', overflow: 'hidden' }}>
    <TableContainer>
      <Table stickyHeader>
        <TableHead>
          <TableRow>
            <TableCell>Product</TableCell>
            <TableCell>SKU</TableCell>
            <TableCell align="right">Current Stock</TableCell>
            <TableCell align="right">Reserved</TableCell>
            <TableCell align="right">Available</TableCell>
            <TableCell>Status</TableCell>
            <TableCell align="right">Value</TableCell>
            <TableCell align="center">Actions</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {inventoryLoading ? (
            <TableRow>
              <TableCell colSpan={8} align="center" sx={{ py: 4 }}>
                <CircularProgress />
              </TableCell>
            </TableRow>
          ) : inventoryItems.length === 0 ? (
            <TableRow>
              <TableCell colSpan={8} align="center" sx={{ py: 4 }}>
                <Typography variant="body1" color="text.secondary">
                  No inventory items found
                </Typography>
              </TableCell>
            </TableRow>
          ) : (
            (Array.isArray(inventoryItems) ? inventoryItems : [])
              .filter(item => item && typeof item === 'object' && item.id)
              .map((item) => {
                if (!item || typeof item !== 'object' || !item.id) return null;
                
                // Safely extract and convert values with fallbacks
                const currentStock = item.quantity != null ? Number(item.quantity) : 0;
                const reserved = item.reserved != null ? Number(item.reserved) : 0;
                const available = item.available != null && Number.isFinite(Number(item.available))
                  ? Number(item.available)
                  : Math.max(0, currentStock - reserved);
                const sellingPrice = item.selling_price != null ? Number(item.selling_price) : 0;
                const stockValue = currentStock * sellingPrice;
                const productName = item.name || 'Unnamed Product';
                const sku = item.sku || 'N/A';
              return (
                <TableRow key={item.id} hover>
                  <TableCell>
                    <Typography variant="body2" fontWeight="medium">
                      {productName}
                    </Typography>
                  </TableCell>
                  <TableCell>{sku}</TableCell>
                  <TableCell align="right">{currentStock}</TableCell>
                  <TableCell align="right">{reserved}</TableCell>
                  <TableCell align="right">{available}</TableCell>
                  <TableCell>
                    <Chip
                      label={getStockStatusText(available)}
                      color={getStockStatusColor(available)}
                      size="small"
                      variant="outlined"
                    />
                  </TableCell>
                  <TableCell align="right">
                    {formatCurrency(stockValue)}
                  </TableCell>
                  <TableCell align="center">
                    <IconButton
                      size="small"
                      onClick={(e) => {
                        setSelectedInventory(item);
                        handleMenuOpen(e, item);
                      }}
                      aria-label="actions"
                    >
                      <MoreVertIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              );
            })
          )}
        </TableBody>
      </Table>
    </TableContainer>
    {/* Pagination can be added here if inventoryItems supports it */}
  </Paper>
)}

      {/* Reservations Tab */}
      {tabValue === 1 && (
        <Box>
          <Typography variant="h6" sx={{ mb: 2 }}>Active Reservations</Typography>
          
          {reservationsLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : reservations.length === 0 ? (
            <Paper sx={{ p: 3, textAlign: 'center' }}>
              <Typography variant="body1" color="text.secondary">
                No active reservations
              </Typography>
            </Paper>
          ) : (
            <Grid container spacing={2}>
              {reservations.map((reservation) => (
                <Grid item xs={12} md={6} lg={4} key={reservation.id}>
                  <Card>
                    <CardContent>
                      <Typography variant="h6" gutterBottom>
                        {reservation.productName}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Quantity: {reservation.quantity}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Order ID: {reservation.orderId}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Reserved: {formatDate(reservation.createdAt)}
                      </Typography>
                      {reservation.notes && (
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                          Notes: {reservation.notes}
                        </Typography>
                      )}
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          )}
        </Box>
      )}

      {/* Low Stock Alerts Tab */}
      {tabValue === 2 && (
        <Box>
          <Typography variant="h6" sx={{ mb: 2 }}>Low Stock Alerts</Typography>
          
          {/* Location and Threshold Filters */}
          <Paper sx={{ p: 2, mb: 3 }}>
            <Grid container spacing={2} alignItems="center">
              <Grid item xs={12} sm={6} md={4}>
                <TextField
                  fullWidth
                  label="Store Location"
                  placeholder="Search store locations..."
                  value={locationInputValue}
                  onChange={(e) => setLocationInputValue(e.target.value)}
                  select={locationOptions.length > 0}
                  SelectProps={{
                    native: false,
                  }}
                  helperText="Filter by store location"
                >
                  <MenuItem value="">
                    <em>All Locations</em>
                  </MenuItem>
                  {locationOptions.map((option) => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </TextField>
              </Grid>
              
              <Grid item xs={12} sm={6} md={3}>
                <TextField
                  fullWidth
                  type="number"
                  label="Stock Threshold"
                  value={lowStockThreshold}
                  onChange={(e) => setLowStockThreshold(Math.max(0, parseInt(e.target.value) || 0))}
                  inputProps={{ min: 0, max: 1000 }}
                  helperText="Items below this quantity"
                />
              </Grid>
              
              <Grid item xs={12} sm={12} md={5}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Button
                    variant="outlined"
                    onClick={() => {
                      setSelectedLocation(locationInputValue);
                      // Trigger refetch of low stock items with new filters
                      queryClient.invalidateQueries(['inventory-low-stock']);
                    }}
                    disabled={lowStockLoading}
                  >
                    Apply Filters
                  </Button>
                  <Button
                    variant="text"
                    onClick={() => {
                      setLocationInputValue('');
                      setSelectedLocation('');
                      setLowStockThreshold(10);
                      queryClient.invalidateQueries(['inventory-low-stock']);
                    }}
                    disabled={lowStockLoading}
                  >
                    Clear
                  </Button>
                  {(selectedLocation || lowStockThreshold !== 10) && (
                    <Chip
                      label={`Filtered: ${selectedLocation || 'All'} | Threshold: ${lowStockThreshold}`}
                      onDelete={() => {
                        setLocationInputValue('');
                        setSelectedLocation('');
                        setLowStockThreshold(10);
                        queryClient.invalidateQueries(['inventory-low-stock']);
                      }}
                      size="small"
                      color="primary"
                      variant="outlined"
                    />
                  )}
                </Box>
              </Grid>
            </Grid>
          </Paper>
          
          {lowStockLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : lowStockItems.length === 0 ? (
            <Paper sx={{ p: 3, textAlign: 'center' }}>
              <CheckIcon color="success" sx={{ fontSize: 48, mb: 2 }} />
              <Typography variant="h6" color="success.main">
                All products are well stocked!
              </Typography>
            </Paper>
          ) : (
            <List>
              {lowStockItems.map((item, index) => (
                <React.Fragment key={item.id}>
                  <ListItem>
                    <ListItemText
                      primary={item.name}
                      secondary={`SKU: ${item.sku} | Current Stock: ${item.stockQuantity}`}
                    />
                    <InventoryErrorBoundary>
                      <InventoryStatus
                        quantity={inventoryItems.find(inv => inv.product_id === item.id)?.quantity}
                        lowStockThreshold={inventoryItems.find(inv => inv.product_id === item.id)?.low_stock_threshold}
                        showQuantity={true}
                      />
                    </InventoryErrorBoundary>
                  </ListItem>
                  {index < lowStockItems.length - 1 && <Divider />}
                </React.Fragment>
              ))}
            </List>
          )}
        </Box>
      )}

      {/* Analytics Tab */}
      {tabValue === 3 && (
        <Box>
          <Typography variant="h6" sx={{ mb: 2 }}>Inventory Analytics</Typography>
          
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography variant="h6" color="text.secondary">Total Products</Typography>
                  <Typography variant="h4" color="primary">
                    {totalProducts}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography variant="h6" color="text.secondary">Low Stock Items</Typography>
                  <Typography variant="h4" color="warning.main">
                    {lowStockItems.length}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography variant="h6" color="text.secondary">Active Reservations</Typography>
                  <Typography variant="h4" color="info.main">
                    {reservations.length}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card>
                <CardContent>
                  <Typography variant="h6" color="text.secondary">Total Inventory Value</Typography>
                  <Typography variant="h4" color="success.main">
                    {formatCurrency(
                      products.reduce((total, product) => 
                        total + (product.stockQuantity || 0) * (product.price || 0), 0
                      )
                    )}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        </Box>
      )}

      {/* Product Action Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={handleCheckStock}>
          <CheckIcon sx={{ mr: 1 }} fontSize="small" />
          Check Stock
        </MenuItem>
        <MenuItem onClick={() => { setReservationDialogOpen(true); handleMenuClose(); }}>
          <ReserveIcon sx={{ mr: 1 }} fontSize="small" />
          Reserve Inventory
        </MenuItem>
        <MenuItem onClick={() => { setReleaseDialogOpen(true); handleMenuClose(); }}>
          <ReleaseIcon sx={{ mr: 1 }} fontSize="small" />
          Release Inventory
        </MenuItem>
        <MenuItem onClick={() => { setAdjustmentDialogOpen(true); handleMenuClose(); }}>
          <AddIcon sx={{ mr: 1 }} fontSize="small" />
          Adjust Inventory
        </MenuItem>
        <MenuItem onClick={() => { setHistoryDialogOpen(true); handleMenuClose(); }}>
          <HistoryIcon sx={{ mr: 1 }} fontSize="small" />
          View History
        </MenuItem>
        <MenuItem 
          onClick={() => { 
            setRemoveItemDialogOpen(true); 
            handleMenuClose(); 
          }}
          sx={{ color: 'error.main' }}
        >
          <DeleteIcon sx={{ mr: 1 }} fontSize="small" />
          Remove Item
        </MenuItem>
      </Menu>

      {/* Reserve Inventory Dialog */}
      <Dialog open={reservationDialogOpen} onClose={() => setReservationDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Reserve Inventory</DialogTitle>
        <DialogContent>
          {selectedProduct && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="h6">{selectedProduct.name}</Typography>
              <Typography variant="body2" color="text.secondary">
                Available: {(selectedProduct.stockQuantity || 0) - (selectedProduct.reservedQuantity || 0)}
              </Typography>
            </Box>
          )}
          
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Quantity"
                value={reservationData.quantity}
                onChange={(e) => setReservationData({ ...reservationData, quantity: parseInt(e.target.value) || 0 })}
                inputProps={{ min: 1 }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Order ID"
                value={reservationData.orderId}
                onChange={(e) => setReservationData({ ...reservationData, orderId: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={2}
                label="Notes"
                value={reservationData.notes}
                onChange={(e) => setReservationData({ ...reservationData, notes: e.target.value })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setReservationDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleReserveInventory} 
            variant="contained"
            disabled={reservationData.quantity <= 0 || reserveInventoryMutation.isLoading}
          >
            {reserveInventoryMutation.isLoading ? <CircularProgress size={20} /> : 'Reserve'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Release Inventory Dialog */}
      <Dialog open={releaseDialogOpen} onClose={() => setReleaseDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Release Inventory</DialogTitle>
        <DialogContent>
          {selectedInventory && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="h6">{selectedInventory.name}</Typography>
              <Typography variant="body2" color="text.secondary">
                Current Stock: {selectedInventory.stockQuantity || 0}
              </Typography>
            </Box>
          )}
          
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Reservation ID"
                value={releaseData.reservationId}
                onChange={(e) => setReleaseData({ ...releaseData, reservationId: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Quantity"
                value={releaseData.quantity}
                onChange={(e) => setReleaseData({ ...releaseData, quantity: parseInt(e.target.value) || 0 })}
                inputProps={{ min: 1 }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Reason"
                value={releaseData.reason}
                onChange={(e) => setReleaseData({ ...releaseData, reason: e.target.value })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setReleaseDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleReleaseInventory} 
            variant="contained"
            disabled={releaseData.quantity <= 0 || releaseInventoryMutation.isLoading}
          >
            {releaseInventoryMutation.isLoading ? <CircularProgress size={20} /> : 'Release'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Adjust Inventory Dialog */}
      <Dialog open={adjustmentDialogOpen} onClose={() => setAdjustmentDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Adjust Inventory</DialogTitle>
        <DialogContent>
          {selectedInventory && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="h6">{selectedInventory.name}</Typography>
              <Typography variant="body2" color="text.secondary">
                Current Stock: {selectedInventory.stockQuantity || 0}
              </Typography>
            </Box>
          )}
          
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth>
                <InputLabel>Type</InputLabel>
                <Select
                  value={adjustmentData.type}
                  label="Type"
                  onChange={(e) => setAdjustmentData({ ...adjustmentData, type: e.target.value })}
                >
                  <MenuItem value="ADD">Add Stock</MenuItem>
                  <MenuItem value="REMOVE">Remove Stock</MenuItem>
                  <MenuItem value="SET">Set Stock</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Quantity"
                value={adjustmentData.quantity}
                onChange={(e) => setAdjustmentData({ ...adjustmentData, quantity: parseInt(e.target.value) || 0 })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Reason"
                value={adjustmentData.reason}
                onChange={(e) => setAdjustmentData({ ...adjustmentData, reason: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={2}
                label="Notes"
                value={adjustmentData.notes}
                onChange={(e) => setAdjustmentData({ ...adjustmentData, notes: e.target.value })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAdjustmentDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleAdjustInventory} 
            variant="contained"
            disabled={adjustmentData.quantity === 0 || !adjustmentData.reason || updateQuantityMutation.isLoading}
          >
            {updateQuantityMutation.isLoading ? <CircularProgress size={20} /> : 'Adjust'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Inventory History Dialog */}
      <Dialog open={historyDialogOpen} onClose={() => setHistoryDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Inventory History</DialogTitle>
        <DialogContent>
          {selectedInventory && (
            <Typography variant="h6" sx={{ mb: 2 }}>
              {selectedInventory.name}
            </Typography>
          )}
          
          {historyLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : historyData?.length === 0 ? (
            <Typography variant="body2" color="text.secondary" sx={{ py: 4, textAlign: 'center' }}>
              No history available
            </Typography>
          ) : (
            <List>
              {(historyData || []).map((entry, index) => (
                <React.Fragment key={entry.id || index}>
                  <ListItem>
                    <ListItemText
                      primary={`${entry.type}: ${entry.quantity} units`}
                      secondary={
                        <Box>
                          <Typography variant="body2">
                            {entry.reason} - {formatDate(entry.createdAt)}
                          </Typography>
                          {entry.notes && (
                            <Typography variant="body2" color="text.secondary">
                              {entry.notes}
                            </Typography>
                          )}
                        </Box>
                      }
                    />
                  </ListItem>
                  {index < (historyData?.length || 0) - 1 && <Divider />}
                </React.Fragment>
              ))}
            </List>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setHistoryDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Create Item Dialog */}
      <Dialog open={createItemDialogOpen} onClose={() => setCreateItemDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <AddIcon />
            Create New Inventory Item
          </Box>
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              Create a new inventory item to track stock levels and manage product availability.
            </Typography>
            
            <Grid container spacing={3}>
              {/* Product Lookup */}
              <Grid item xs={12}>
                <Typography variant="h6" sx={{ mb: 2, color: 'primary.main' }}>
                  Product Lookup
                </Typography>
              </Grid>
              
              <Grid item xs={12}>
  <Autocomplete
    fullWidth
    options={productOptions}
    loading={productLookupLoading}
    value={selectedInventory?.productOption || null}
    inputValue={productLookupValue}
    onInputChange={async (event, value, reason) => {
      setProductLookupValue(value);
      if (reason === 'input' && value.trim().length > 1) {
        setProductLookupLoading(true);
        try {
          const results = await productService.getProducts({ q: value });
          setProductOptions(results);
        } catch (e) {
          setProductOptions([]);
        } finally {
          setProductLookupLoading(false);
        }
      }
    }}
    onChange={(event, newValue) => {
      if (newValue) {
        setSelectedInventory({
          ...selectedInventory,
          productOption: newValue,
          name: newValue.name,
          sku: newValue.sku,
          productId: newValue.id // invisible to user, sent to backend
        });
      }
    }}
    getOptionLabel={(option) => option.name ? `${option.name} – ${option.sku}` : ''}
    renderInput={(params) => (
      <TextField
        {...params}
        label="Product Lookup (SKU or Barcode)"
        placeholder="Type to search by SKU or barcode"
        InputProps={{
          ...params.InputProps,
          endAdornment: (
            <>
              {productLookupLoading ? <CircularProgress size={20} /> : null}
              {params.InputProps.endAdornment}
            </>
          ),
        }}
        helperText="Search and select a product to link inventory item."
      />
    )}
  />
</Grid>

              {/* Basic Information */}
              <Grid item xs={12}>
                <Typography variant="h6" sx={{ mb: 2, mt: 2, color: 'primary.main' }}>
                  Basic Information
                </Typography>
              </Grid>
              
              <Grid item xs={12} sm={8}>
                <TextField
                  fullWidth
                  label="Product Name"
                  placeholder="Enter product name"
                  value={selectedInventory?.name || ''}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, name: e.target.value })}
                  required
                  error={!selectedInventory?.name}
                  helperText={!selectedInventory?.name ? "Product name is required" : ""}
                />
              </Grid>
              
              <Grid item xs={12} sm={4}>
                <TextField
                  fullWidth
                  label="SKU"
                  placeholder="e.g., PROD-001"
                  value={selectedInventory?.sku || ''}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, sku: e.target.value.toUpperCase() })}
                  required
                  error={!selectedInventory?.sku}
                  helperText={!selectedInventory?.sku ? "SKU is required" : "Unique product identifier"}
                />
              </Grid>
              
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  multiline
                  rows={2}
                  label="Description"
                  placeholder="Optional product description"
                  value={selectedInventory?.description || ''}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, description: e.target.value })}
                  helperText="Brief description of the product"
                />
              </Grid>

              {/* Stock Information */}
              <Grid item xs={12}>
                <Typography variant="h6" sx={{ mb: 2, mt: 2, color: 'primary.main' }}>
                  Stock Information
                </Typography>
              </Grid>
              
              <Grid item xs={12} sm={4}>
                <TextField
                  fullWidth
                  type="number"
                  label="Initial Quantity"
                  value={selectedInventory?.quantity || 0}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, quantity: Math.max(0, parseInt(e.target.value) || 0) })}
                  inputProps={{ min: 0 }}
                  helperText="Starting stock quantity"
                />
              </Grid>
              
              <Grid item xs={12} sm={4}>
                <TextField
                  fullWidth
                  type="number"
                  label="Minimum Stock Level"
                  value={selectedInventory?.minStockLevel || 0}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, minStockLevel: Math.max(0, parseInt(e.target.value) || 0) })}
                  inputProps={{ min: 0 }}
                  helperText="Low stock alert threshold"
                />
              </Grid>
              
              <Grid item xs={12} sm={4}>
                <TextField
                  fullWidth
                  type="number"
                  label="Maximum Stock Level"
                  value={selectedInventory?.maxStockLevel || 0}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, maxStockLevel: Math.max(0, parseInt(e.target.value) || 0) })}
                  inputProps={{ min: 0 }}
                  helperText="Maximum stock capacity"
                />
              </Grid>

              {/* Location & Category */}
              <Grid item xs={12}>
                <Typography variant="h6" sx={{ mb: 2, mt: 2, color: 'primary.main' }}>
                  Organization
                </Typography>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <FormControl fullWidth>
                  <InputLabel>Store Location</InputLabel>
                  <Select
                    value={selectedInventory?.storeId || ''}
                    onChange={(e) => setSelectedInventory({ ...selectedInventory, storeId: e.target.value })}
                    label="Store Location"
                  >
                    <MenuItem value="">
                      <em>Select a store</em>
                    </MenuItem>
                    {storeOptions.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {option.label}
                      </MenuItem>
                    ))}
                    <Divider />
                    <MenuItem onClick={() => setStoreDialogOpen(true)}>
                      <AddIcon sx={{ mr: 1 }} fontSize="small" />
                      Create New Store
                    </MenuItem>
                  </Select>
                  <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, ml: 1.5 }}>
                    Select the physical store location
                  </Typography>
                </FormControl>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  label="Storage Location"
                  placeholder="e.g., Warehouse A, Shelf B-3"
                  value={selectedInventory?.location || ''}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, location: e.target.value })}
                  helperText="Physical storage location within store"
                />
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  label="Category"
                  placeholder="e.g., Electronics, Clothing"
                  value={selectedInventory?.category || ''}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, category: e.target.value })}
                  helperText="Product category"
                />
              </Grid>
              
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  multiline
                  rows={2}
                  label="Notes"
                  placeholder="Additional notes or special instructions"
                  value={selectedInventory?.notes || ''}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, notes: e.target.value })}
                  helperText="Optional notes for this inventory item"
                />
              </Grid>
            </Grid>
            
            {/* Validation Summary */}
            {(!selectedInventory?.name || !selectedInventory?.sku) && (
              <Alert severity="warning" sx={{ mt: 3 }}>
                <AlertTitle>Required Fields Missing</AlertTitle>
                Please fill in all required fields (Product Name and SKU) before creating the item.
              </Alert>
            )}
          </Box>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 3 }}>
          <Button 
            onClick={() => {
              setCreateItemDialogOpen(false);
              setSelectedInventory(null);
            }}
            size="large"
          >
            Cancel
          </Button>
          <Button 
            onClick={() => {
              if (selectedInventory?.productId && selectedInventory?.name && selectedInventory?.sku) {
                createItemMutation.mutate({
                  productId: selectedInventory.productId,
                  sku: selectedInventory.sku, // Add SKU to the request payload
                  quantity: selectedInventory.quantity || 0,
                  minStockLevel: selectedInventory.minStockLevel || 0,
                  maxStockLevel: selectedInventory.maxStockLevel || 0,
                  location: selectedInventory.location || '',
                  storeId: selectedInventory.storeId || '',
                  notes: selectedInventory.notes || '',
                });
              }
            }} 
            variant="contained"
            size="large"
            disabled={!selectedInventory?.name || !selectedInventory?.sku || createItemMutation.isLoading}
            startIcon={createItemMutation.isLoading ? <CircularProgress size={20} /> : <AddIcon />}
          >
            {createItemMutation.isLoading ? 'Creating...' : 'Create Item'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Remove Item Dialog */}
      <Dialog open={removeItemDialogOpen} onClose={() => setRemoveItemDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Remove Item</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <Typography variant="body2" color="text.secondary">
                Are you sure you want to remove this item?
              </Typography>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRemoveItemDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={() => {
              removeItemMutation.mutate(selectedInventory?.id);
            }} 
            variant="contained"
            color="error"
          >
            Remove
          </Button>
        </DialogActions>
      </Dialog>

      {/* Create Store Dialog */}
      <Dialog open={storeDialogOpen} onClose={() => setStoreDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <AddIcon />
            Create New Store
          </Box>
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
              Create a new physical store location for inventory management.
            </Typography>
            
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Store Name"
                  placeholder="e.g., Main Warehouse, Downtown Store"
                  value={newStore.name}
                  onChange={(e) => setNewStore({ ...newStore, name: e.target.value })}
                  required
                  error={!newStore.name}
                  helperText={!newStore.name ? "Store name is required" : ""}
                />
              </Grid>
              
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Address"
                  placeholder="Full store address"
                  value={newStore.address}
                  onChange={(e) => setNewStore({ ...newStore, address: e.target.value })}
                  required
                  error={!newStore.address}
                  helperText={!newStore.address ? "Address is required" : ""}
                />
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  label="Phone"
                  placeholder="Store phone number"
                  value={newStore.phone}
                  onChange={(e) => setNewStore({ ...newStore, phone: e.target.value })}
                  helperText="Optional contact number"
                />
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  label="Email"
                  type="email"
                  placeholder="Store email address"
                  value={newStore.email}
                  onChange={(e) => setNewStore({ ...newStore, email: e.target.value })}
                  helperText="Optional contact email"
                />
              </Grid>
            </Grid>
            
            {/* Validation Summary */}
            {(!newStore.name || !newStore.address) && (
              <Alert severity="warning" sx={{ mt: 2 }}>
                <AlertTitle>Required Fields Missing</AlertTitle>
                Please fill in the store name and address before creating the store.
              </Alert>
            )}
          </Box>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 3 }}>
          <Button 
            onClick={() => {
              setStoreDialogOpen(false);
              setNewStore({ name: '', address: '', phone: '', email: '' });
            }}
            size="large"
          >
            Cancel
          </Button>
          <Button 
            onClick={handleCreateStore}
            variant="contained"
            size="large"
            disabled={!newStore.name || !newStore.address}
            startIcon={<AddIcon />}
          >
            Create Store
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default InventoryPage;
