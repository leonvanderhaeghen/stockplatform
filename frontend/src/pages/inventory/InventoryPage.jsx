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
  Inventory,
  Search,
  Add,
  Remove,
  Edit,
  Warning,
  Refresh,
  GetApp,
  Analytics,
  ShoppingCart,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSnackbar } from 'notistack';
import { format } from 'date-fns';
import inventoryService from '../../services/inventoryService';
import productService from '../../services/productService';

const InventoryPage = () => {
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();

  // State for UI
  const [currentTab, setCurrentTab] = useState(0);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [lowStockOnly, setLowStockOnly] = useState(false);
  const [selectedItem, setSelectedItem] = useState(null);
  const [showAdjustmentDialog, setShowAdjustmentDialog] = useState(false);
  const [adjustmentType, setAdjustmentType] = useState('add');
  const [adjustmentQuantity, setAdjustmentQuantity] = useState('');
  const [adjustmentReason, setAdjustmentReason] = useState('');
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [newItemData, setNewItemData] = useState({ productId: '', sku: '', initialQuantity: '' });

  // Fetch inventory data
  const {
    data: inventoryData = { items: [], total: 0 },
    isLoading: inventoryLoading,
    error: inventoryError
  } = useQuery({
    queryKey: ['inventory', page + 1, rowsPerPage, searchTerm, statusFilter, lowStockOnly],
    queryFn: () => inventoryService.getInventory({
      page: page + 1,
      limit: rowsPerPage,
      search: searchTerm || undefined,
      status: statusFilter !== 'all' ? statusFilter : undefined,
      lowStock: lowStockOnly || undefined
    }),
  });

  // Fetch low stock items for alerts
  const { data: lowStockItems = [] } = useQuery({
    queryKey: ['low-stock-items'],
    queryFn: () => inventoryService.getLowStockItems(),
  });

  // Fetch reservations
  const { data: reservations = [] } = useQuery({
    queryKey: ['inventory-reservations'],
    queryFn: () => inventoryService.getInventoryReservations(),
    enabled: currentTab === 1,
  });

  // Fetch products for create dialog
  useQuery({
    queryKey: ['products', 'inventory-create'],
    queryFn: () => productService.getProducts({ page: 1, limit: 100, status: 'active' }),
    select: (data) => data.products || [],
    enabled: showCreateDialog,
  });

  // Stock adjustment mutation
  const adjustStockMutation = useMutation({
    mutationFn: async ({ itemId, type, quantity, reason }) => {
      if (type === 'add') {
        return await inventoryService.addStock(itemId, quantity, reason);
      } else {
        return await inventoryService.removeStock(itemId, quantity, reason);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory']);
      queryClient.invalidateQueries(['low-stock-items']);
      setShowAdjustmentDialog(false);
      setSelectedItem(null);
      setAdjustmentQuantity('');
      setAdjustmentReason('');
      enqueueSnackbar('Stock adjustment completed successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to adjust stock', { variant: 'error' });
    },
  });

  // Create inventory item mutation
  const createItemMutation = useMutation({
    mutationFn: async (itemData) => {
      return await productService.createProduct({
        name: itemData.name,
        sku: itemData.sku,
        description: itemData.description || '',
        price: 0,
        categoryId: 1,
        supplierId: 1,
        stockQuantity: parseInt(itemData.initialQuantity),
        status: 'active'
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory']);
      setShowCreateDialog(false);
      setNewItemData({ productId: '', sku: '', initialQuantity: '' });
      enqueueSnackbar('Inventory item created successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to create inventory item', { variant: 'error' });
    },
  });

  // Delete inventory item mutation
  const deleteItemMutation = useMutation({
    mutationFn: async (itemId) => {
      return await productService.deleteProduct(itemId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['inventory']);
      enqueueSnackbar('Inventory item removed successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.message || 'Failed to remove inventory item', { variant: 'error' });
    },
  });

  const handleStockAdjustment = () => {
    if (!selectedItem || !adjustmentQuantity || !adjustmentReason) {
      enqueueSnackbar('Please fill in all required fields', { variant: 'warning' });
      return;
    }

    adjustStockMutation.mutate({
      itemId: selectedItem.id,
      type: adjustmentType,
      quantity: parseInt(adjustmentQuantity),
      reason: adjustmentReason
    });
  };

  const handleCreateItem = () => {
    if (!newItemData.name || !newItemData.sku || !newItemData.initialQuantity) {
      enqueueSnackbar('Please fill in all required fields', { variant: 'warning' });
      return;
    }

    createItemMutation.mutate(newItemData);
  };

  const getStockStatusChip = (item) => {
    const quantity = item.stockQuantity || 0;
    const minThreshold = item.minThreshold || 10;
    
    if (quantity === 0) {
      return <Chip label="Out of Stock" color="error" size="small" />;
    } else if (quantity <= minThreshold) {
      return <Chip label="Low Stock" color="warning" size="small" />;
    } else {
      return <Chip label="In Stock" color="success" size="small" />;
    }
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
            <Inventory sx={{ fontSize: 40, mr: 2 }} />
            <Box>
              <Typography variant="h4" gutterBottom>
                Inventory Management
              </Typography>
              <Typography variant="subtitle1">
                Track and manage inventory levels, stock movements, and reservations
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
              Create Item
            </Button>
            <Button
              variant="contained"
              color="secondary"
              startIcon={<Refresh />}
              onClick={() => queryClient.invalidateQueries(['inventory'])}
            >
              Refresh
            </Button>
          </Box>
        </Box>
      </Paper>

      {/* Low Stock Alert */}
      {lowStockItems.length > 0 && (
        <Alert severity="warning" sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <Typography>
              ⚠️ {lowStockItems.length} item(s) are running low on stock
            </Typography>
            <Button
              size="small"
              onClick={() => setLowStockOnly(true)}
              startIcon={<Warning />}
            >
              View Low Stock
            </Button>
          </Box>
        </Alert>
      )}

      <Card>
        <CardContent>
          {/* Tabs */}
          <Tabs value={currentTab} onChange={(e, newValue) => setCurrentTab(newValue)}>
            <Tab label="Inventory Items" icon={<Inventory />} />
            <Tab label="Reservations" icon={<ShoppingCart />} />
            <Tab label="Analytics" icon={<Analytics />} />
          </Tabs>

          {/* Inventory Items Tab */}
          <TabPanel value={currentTab} index={0}>
            {/* Filters and Search */}
            <Grid container spacing={2} sx={{ mb: 3 }}>
              <Grid item xs={12} md={4}>
                <TextField
                  fullWidth
                  placeholder="Search by name, SKU, or description..."
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
              <Grid item xs={12} md={3}>
                <FormControl fullWidth>
                  <InputLabel>Status Filter</InputLabel>
                  <Select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value)}
                    label="Status Filter"
                  >
                    <MenuItem value="all">All Items</MenuItem>
                    <MenuItem value="in_stock">In Stock</MenuItem>
                    <MenuItem value="low_stock">Low Stock</MenuItem>
                    <MenuItem value="out_of_stock">Out of Stock</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={3}>
                <Button
                  fullWidth
                  variant={lowStockOnly ? 'contained' : 'outlined'}
                  startIcon={<Warning />}
                  onClick={() => setLowStockOnly(!lowStockOnly)}
                  sx={{ height: '56px' }}
                >
                  {lowStockOnly ? 'Show All' : 'Low Stock Only'}
                </Button>
              </Grid>
              <Grid item xs={12} md={2}>
                <Button
                  fullWidth
                  variant="outlined"
                  startIcon={<GetApp />}
                  sx={{ height: '56px' }}
                >
                  Export
                </Button>
              </Grid>
            </Grid>

            {/* Inventory Table */}
            {inventoryError ? (
              <Alert severity="error" sx={{ mb: 2 }}>
                Failed to load inventory: {inventoryError.message}
              </Alert>
            ) : (
              <>
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Product</TableCell>
                        <TableCell>SKU</TableCell>
                        <TableCell align="center">Stock Qty</TableCell>
                        <TableCell align="center">Reserved</TableCell>
                        <TableCell align="center">Available</TableCell>
                        <TableCell align="center">Status</TableCell>
                        <TableCell align="center">Last Updated</TableCell>
                        <TableCell align="center">Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {inventoryLoading ? (
                        // Loading skeletons
                        [...Array(rowsPerPage)].map((_, index) => (
                          <TableRow key={index}>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="rectangular" width={80} height={24} /></TableCell>
                            <TableCell><Skeleton variant="text" /></TableCell>
                            <TableCell><Skeleton variant="rectangular" width={100} height={36} /></TableCell>
                          </TableRow>
                        ))
                      ) : inventoryData.items.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={8} align="center">
                            <Typography color="text.secondary" sx={{ py: 4 }}>
                              No inventory items found. {searchTerm && 'Try adjusting your search terms.'}
                            </Typography>
                          </TableCell>
                        </TableRow>
                      ) : (
                        inventoryData.items.map((item) => {
                          const reserved = item.reservedQuantity || 0;
                          const available = (item.stockQuantity || 0) - reserved;
                          return (
                            <TableRow key={item.id}>
                              <TableCell>
                                <Box>
                                  <Typography variant="body1" fontWeight="medium">
                                    {item.name}
                                  </Typography>
                                  {item.description && (
                                    <Typography variant="caption" color="text.secondary">
                                      {item.description}
                                    </Typography>
                                  )}
                                </Box>
                              </TableCell>
                              <TableCell>
                                <Typography variant="body2" fontFamily="monospace">
                                  {item.sku}
                                </Typography>
                              </TableCell>
                              <TableCell align="center">
                                <Typography variant="body2" fontWeight="medium">
                                  {item.stockQuantity || 0}
                                </Typography>
                              </TableCell>
                              <TableCell align="center">
                                <Typography variant="body2" color={reserved > 0 ? 'warning.main' : 'text.secondary'}>
                                  {reserved}
                                </Typography>
                              </TableCell>
                              <TableCell align="center">
                                <Typography variant="body2" fontWeight="medium" color={available > 0 ? 'success.main' : 'error.main'}>
                                  {available}
                                </Typography>
                              </TableCell>
                              <TableCell align="center">
                                {getStockStatusChip(item)}
                              </TableCell>
                              <TableCell align="center">
                                <Typography variant="caption" color="text.secondary">
                                  {item.updatedAt ? format(new Date(item.updatedAt), 'MMM dd, yyyy') : 'N/A'}
                                </Typography>
                              </TableCell>
                              <TableCell align="center">
                                <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center' }}>
                                  <IconButton
                                    size="small"
                                    onClick={() => {
                                      setSelectedItem(item);
                                      setAdjustmentType('add');
                                      setShowAdjustmentDialog(true);
                                    }}
                                    title="Add Stock"
                                  >
                                    <Add color="success" />
                                  </IconButton>
                                  <IconButton
                                    size="small"
                                    onClick={() => {
                                      setSelectedItem(item);
                                      setAdjustmentType('remove');
                                      setShowAdjustmentDialog(true);
                                    }}
                                    title="Remove Stock"
                                  >
                                    <Remove color="warning" />
                                  </IconButton>
                                  <IconButton
                                    size="small"
                                    onClick={() => deleteItemMutation.mutate(item.id)}
                                    title="Remove Item"
                                    color="error"
                                  >
                                    <Edit />
                                  </IconButton>
                                </Box>
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
                  count={inventoryData.total || 0}
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

          {/* Reservations Tab */}
          <TabPanel value={currentTab} index={1}>
            <Typography variant="h6" gutterBottom>
              Active Inventory Reservations
            </Typography>
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Order ID</TableCell>
                    <TableCell>Product</TableCell>
                    <TableCell align="center">Quantity</TableCell>
                    <TableCell align="center">Reserved Date</TableCell>
                    <TableCell align="center">Expires</TableCell>
                    <TableCell align="center">Status</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {reservations.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} align="center">
                        <Typography color="text.secondary" sx={{ py: 4 }}>
                          No active reservations found.
                        </Typography>
                      </TableCell>
                    </TableRow>
                  ) : (
                    reservations.map((reservation) => (
                      <TableRow key={reservation.id}>
                        <TableCell>
                          <Typography variant="body2" fontFamily="monospace">
                            #{reservation.orderId}
                          </Typography>
                        </TableCell>
                        <TableCell>{reservation.productName}</TableCell>
                        <TableCell align="center">{reservation.quantity}</TableCell>
                        <TableCell align="center">
                          {format(new Date(reservation.createdAt), 'MMM dd, yyyy HH:mm')}
                        </TableCell>
                        <TableCell align="center">
                          {reservation.expiresAt ? format(new Date(reservation.expiresAt), 'MMM dd, yyyy HH:mm') : 'No expiry'}
                        </TableCell>
                        <TableCell align="center">
                          <Chip
                            label={reservation.status}
                            color={reservation.status === 'active' ? 'primary' : 'default'}
                            size="small"
                          />
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          </TabPanel>

          {/* Analytics Tab */}
          <TabPanel value={currentTab} index={2}>
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Card variant="outlined">
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Stock Summary
                    </Typography>
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Total Items:</Typography>
                        <Typography fontWeight="bold">{inventoryData.total || 0}</Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Low Stock Items:</Typography>
                        <Typography fontWeight="bold" color="warning.main">
                          {lowStockItems.length}
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Out of Stock:</Typography>
                        <Typography fontWeight="bold" color="error.main">
                          {inventoryData.items?.filter(item => (item.stockQuantity || 0) === 0).length || 0}
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography>Active Reservations:</Typography>
                        <Typography fontWeight="bold" color="info.main">
                          {reservations.length}
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
                      Recent Activity
                    </Typography>
                    <Typography color="text.secondary">
                      Stock movement analytics and reporting will be available here.
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </TabPanel>
        </CardContent>
      </Card>

      {/* Stock Adjustment Dialog */}
      <Dialog open={showAdjustmentDialog} onClose={() => setShowAdjustmentDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {adjustmentType === 'add' ? 'Add Stock' : 'Remove Stock'} - {selectedItem?.name}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2, display: 'flex', flexDirection: 'column', gap: 3 }}>
            <TextField
              label="Quantity"
              type="number"
              value={adjustmentQuantity}
              onChange={(e) => setAdjustmentQuantity(e.target.value)}
              fullWidth
              inputProps={{ min: 1 }}
            />
            <TextField
              label="Reason"
              value={adjustmentReason}
              onChange={(e) => setAdjustmentReason(e.target.value)}
              fullWidth
              multiline
              rows={3}
              placeholder="Enter reason for stock adjustment..."
            />
            {selectedItem && (
              <Alert severity="info">
                Current stock: {selectedItem.stockQuantity || 0} units
              </Alert>
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowAdjustmentDialog(false)}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handleStockAdjustment}
            disabled={adjustStockMutation.isLoading || !adjustmentQuantity || !adjustmentReason}
            startIcon={adjustStockMutation.isLoading ? <CircularProgress size={20} /> : adjustmentType === 'add' ? <Add /> : <Remove />}
          >
            {adjustStockMutation.isLoading ? 'Processing...' : (adjustmentType === 'add' ? 'Add Stock' : 'Remove Stock')}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Create Item Dialog */}
      <Dialog open={showCreateDialog} onClose={() => setShowCreateDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Create New Inventory Item</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2, display: 'flex', flexDirection: 'column', gap: 3 }}>
            <TextField
              label="Product Name"
              value={newItemData.name}
              onChange={(e) => setNewItemData({ ...newItemData, name: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="SKU"
              value={newItemData.sku}
              onChange={(e) => setNewItemData({ ...newItemData, sku: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Description"
              value={newItemData.description}
              onChange={(e) => setNewItemData({ ...newItemData, description: e.target.value })}
              fullWidth
              multiline
              rows={2}
            />
            <TextField
              label="Initial Quantity"
              type="number"
              value={newItemData.initialQuantity}
              onChange={(e) => setNewItemData({ ...newItemData, initialQuantity: e.target.value })}
              fullWidth
              required
              inputProps={{ min: 0 }}
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowCreateDialog(false)}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handleCreateItem}
            disabled={createItemMutation.isLoading || !newItemData.name || !newItemData.sku || !newItemData.initialQuantity}
            startIcon={createItemMutation.isLoading ? <CircularProgress size={20} /> : <Add />}
          >
            {createItemMutation.isLoading ? 'Creating...' : 'Create Item'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default InventoryPage;
