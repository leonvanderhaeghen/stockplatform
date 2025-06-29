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
  AlertTitle,
  CircularProgress,
  Tabs,
  Tab,
  FormControl,
  InputLabel,
  Select,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Divider,
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
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
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
  useDeleteInventoryItem
} from '../utils/hooks/useInventory';
import InventoryStatus from '../components/inventory/InventoryStatus';
import InventoryErrorBoundary from '../components/inventory/InventoryErrorBoundary';

const InventoryPage = () => {
  const queryClient = useQueryClient();
  const [tabValue, setTabValue] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [anchorEl, setAnchorEl] = useState(null);
  const [reservationDialogOpen, setReservationDialogOpen] = useState(false);
  const [releaseDialogOpen, setReleaseDialogOpen] = useState(false);
  const [adjustmentDialogOpen, setAdjustmentDialogOpen] = useState(false);
  const [historyDialogOpen, setHistoryDialogOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
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

  // Fetch products and inventory data using new hooks
  const { data: productsData, isLoading: productsLoading, refetch: refetchProducts } = useQuery({
    queryKey: ['products', page, rowsPerPage],
    queryFn: () => productService.getProducts({
      page: page + 1,
      limit: rowsPerPage,
    }),
    keepPreviousData: true,
  });

  // Fetch inventory items using the new hook
  const { data: inventoryItems = [] } = useInventoryItems();

  // Defensive extraction of products array
  const products = Array.isArray(productsData)
    ? productsData
    : Array.isArray(productsData?.products)
      ? productsData.products
      : Array.isArray(productsData?.items)
        ? productsData.items
        : Array.isArray(productsData?.data)
          ? productsData.data
          : [];

  // Fetch inventory reservations using new hook
  const { data: reservationsData, isLoading: reservationsLoading } = useInventoryReservations(
    tabValue === 1 ? {} : undefined
  );

  // Fetch low stock alerts using new hook
  const { data: lowStockData, isLoading: lowStockLoading } = useLowStockItems(
    tabValue === 2 ? 10 : undefined
  );

  // Fetch inventory history for selected product using new hook
  const { data: historyData, isLoading: historyLoading } = useInventoryHistory(
    historyDialogOpen && selectedProduct?.id ? selectedProduct.id : null
  );

  // Use new inventory hooks for mutations
  const updateQuantityMutation = useUpdateInventoryQuantity();
  const reserveInventoryMutation = useReserveInventory();
  const releaseInventoryMutation = useReleaseInventory();
  const createItemMutation = useCreateInventoryItem();
  const removeItemMutation = useDeleteInventoryItem();

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
  const totalProducts = productsData?.total || 0;

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

  const handleMenuClick = (event, product) => {
    setAnchorEl(event.currentTarget);
    setSelectedProduct(product);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedProduct(null);
  };

  const handleCheckStock = () => {
    if (selectedProduct) {
      checkStockMutation.mutate(selectedProduct.id);
    }
    handleMenuClose();
  };

  const handleReserveInventory = () => {
    if (selectedProduct && reservationData.quantity > 0) {
      reserveInventoryMutation.mutate({
        productId: selectedProduct.id,
        quantity: reservationData.quantity,
        orderId: reservationData.orderId,
        notes: reservationData.notes
      });
      setReservationDialogOpen(false);
      resetReservationData();
    }
  };

  const handleReleaseInventory = () => {
    if (selectedProduct && releaseData.quantity > 0) {
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
    if (selectedProduct && adjustmentData.quantity !== 0) {
      updateQuantityMutation.mutate({
        id: selectedProduct.id,
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
              refetchProducts();
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
                {productsLoading ? (
                  <TableRow>
                    <TableCell colSpan={8} align="center" sx={{ py: 4 }}>
                      <CircularProgress />
                    </TableCell>
                  </TableRow>
                ) : products.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} align="center" sx={{ py: 4 }}>
                      <Typography variant="body1" color="text.secondary">
                        No products found
                      </Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  products.map((product) => {
                    const currentStock = product.stockQuantity || 0;
                    const reserved = product.reservedQuantity || 0;
                    const available = currentStock - reserved;
                    const stockValue = currentStock * (product.price || 0);

                    return (
                      <TableRow key={product.id} hover>
                        <TableCell>
                          <Typography variant="body2" fontWeight="medium">
                            {product.name}
                          </Typography>
                        </TableCell>
                        <TableCell>{product.sku}</TableCell>
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
                            onClick={(e) => handleMenuClick(e, product)}
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
          
          <TablePagination
            rowsPerPageOptions={[5, 10, 25, 50]}
            component="div"
            count={totalProducts}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
          />
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
                    <ListItemSecondaryAction>
                      <InventoryErrorBoundary>
                        <InventoryStatus
                          quantity={inventoryItems.find(inv => inv.product_id === item.id)?.quantity}
                          lowStockThreshold={inventoryItems.find(inv => inv.product_id === item.id)?.low_stock_threshold}
                          showQuantity={true}
                        />
                      </InventoryErrorBoundary>
                    </ListItemSecondaryAction>
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
          {selectedProduct && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="h6">{selectedProduct.name}</Typography>
              <Typography variant="body2" color="text.secondary">
                Reserved: {selectedProduct.reservedQuantity || 0}
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
          {selectedProduct && (
            <Box sx={{ mb: 2 }}>
              <Typography variant="h6">{selectedProduct.name}</Typography>
              <Typography variant="body2" color="text.secondary">
                Current Stock: {selectedProduct.stockQuantity || 0}
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
          {selectedProduct && (
            <Typography variant="h6" sx={{ mb: 2 }}>
              {selectedProduct.name}
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
              {/* Basic Information */}
              <Grid item xs={12}>
                <Typography variant="h6" sx={{ mb: 2, color: 'primary.main' }}>
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
                <TextField
                  fullWidth
                  label="Storage Location"
                  placeholder="e.g., Warehouse A, Shelf B-3"
                  value={selectedInventory?.location || ''}
                  onChange={(e) => setSelectedInventory({ ...selectedInventory, location: e.target.value })}
                  helperText="Physical storage location"
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
              if (selectedInventory?.name && selectedInventory?.sku) {
                createItemMutation.mutate({
                  ...selectedInventory,
                  // Ensure numeric fields are properly formatted
                  quantity: selectedInventory.quantity || 0,
                  minStockLevel: selectedInventory.minStockLevel || 0,
                  maxStockLevel: selectedInventory.maxStockLevel || 0,
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
    </Box>
  );
};

export default InventoryPage;
