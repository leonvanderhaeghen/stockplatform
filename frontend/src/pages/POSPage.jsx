import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  Grid,
  Card,
  CardContent,
  TextField,
  Alert,
  CircularProgress,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  IconButton,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import {
  ShoppingCart as CartIcon,
  Add as AddIcon,
  Remove as RemoveIcon,
  Delete as DeleteIcon,
  Payment as PaymentIcon,
  Receipt as ReceiptIcon,
  Inventory as InventoryIcon,
  Analytics as AnalyticsIcon,
  Today as TodayIcon,
  Clear as ClearIcon,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { posService, productService, inventoryService } from '../services';
import { formatCurrency, formatDate } from '../utils/formatters';

const POSPage = () => {
  const queryClient = useQueryClient();
  const [cart, setCart] = useState([]);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [productSearchTerm, setProductSearchTerm] = useState('');
  const [paymentDialogOpen, setPaymentDialogOpen] = useState(false);
  const [pickupDialogOpen, setPickupDialogOpen] = useState(false);
  const [deductionDialogOpen, setDeductionDialogOpen] = useState(false);
  const [statsDialogOpen, setStatsDialogOpen] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [paymentData, setPaymentData] = useState({
    paymentMethod: 'CASH',
    customerName: '',
    customerEmail: '',
    notes: '',
  });
  const [pickupData, setPickupData] = useState({
    orderId: '',
    customerInfo: '',
    notes: '',
  });
  const [deductionData, setDeductionData] = useState({
    productId: '',
    quantity: 1,
    reason: '',
    notes: '',
  });

  // Fetch products for POS
  const { data: productsData, isLoading: productsLoading } = useQuery({
    queryKey: ['products', productSearchTerm],
    queryFn: () => productService.getProducts({
      search: productSearchTerm,
      limit: 20,
    }),
    keepPreviousData: true,
  });

  // Defensive extraction of products array for POS product grid
  const products = Array.isArray(productsData)
    ? productsData
    : Array.isArray(productsData?.products)
      ? productsData.products
      : Array.isArray(productsData?.items)
        ? productsData.items
        : Array.isArray(productsData?.data)
          ? productsData.data
          : [];


  // Fetch POS sales statistics
  const { data: salesStats, isLoading: statsLoading } = useQuery({
    queryKey: ['pos-sales-stats'],
    queryFn: () => posService.getSalesStatistics(),
    enabled: statsDialogOpen,
  });

  // Fetch daily summary
  const { data: dailySummary, isLoading: summaryLoading } = useQuery({
    queryKey: ['pos-daily-summary'],
    queryFn: () => posService.getDailySummary(),
  });

  // Quick sale mutation
  const quickSaleMutation = useMutation({
    mutationFn: (saleData) => posService.createQuickSale(saleData),
    onSuccess: (data) => {
      setSuccess(`Sale completed successfully! Transaction ID: ${data.transactionId}`);
      setCart([]);
      setPaymentDialogOpen(false);
      resetPaymentData();
      queryClient.invalidateQueries(['pos-daily-summary']);
      queryClient.invalidateQueries(['pos-sales-stats']);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to complete sale');
    },
  });

  // Pickup completion mutation
  const pickupCompletionMutation = useMutation({
    mutationFn: (pickupData) => posService.completePickup(pickupData),
    onSuccess: () => {
      setSuccess('Pickup completed successfully');
      setPickupDialogOpen(false);
      resetPickupData();
      queryClient.invalidateQueries(['pos-daily-summary']);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to complete pickup');
    },
  });

  // Inventory deduction mutation
  const inventoryDeductionMutation = useMutation({
    mutationFn: (deductionData) => posService.deductInventory(deductionData),
    onSuccess: () => {
      setSuccess('Inventory deduction completed successfully');
      setDeductionDialogOpen(false);
      resetDeductionData();
      queryClient.invalidateQueries(['products']);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to deduct inventory');
    },
  });

  const cartTotal = cart.reduce((total, item) => total + (item.price * item.quantity), 0);
  const cartItemCount = cart.reduce((total, item) => total + item.quantity, 0);

  const addToCart = (product) => {
    const existingItem = cart.find(item => item.id === product.id);
    if (existingItem) {
      setCart(cart.map(item =>
        item.id === product.id
          ? { ...item, quantity: item.quantity + 1 }
          : item
      ));
    } else {
      setCart([...cart, { ...product, quantity: 1 }]);
    }
  };

  const updateCartItemQuantity = (productId, newQuantity) => {
    if (newQuantity <= 0) {
      removeFromCart(productId);
    } else {
      setCart(cart.map(item =>
        item.id === productId
          ? { ...item, quantity: newQuantity }
          : item
      ));
    }
  };

  const removeFromCart = (productId) => {
    setCart(cart.filter(item => item.id !== productId));
  };

  const clearCart = () => {
    setCart([]);
  };

  const handleQuickSale = () => {
    if (cart.length === 0) {
      setError('Cart is empty');
      return;
    }

    const saleData = {
      items: cart.map(item => ({
        productId: item.id,
        quantity: item.quantity,
        unitPrice: item.price,
      })),
      totalAmount: cartTotal,
      paymentMethod: paymentData.paymentMethod,
      customerName: paymentData.customerName,
      customerEmail: paymentData.customerEmail,
      notes: paymentData.notes,
    };

    quickSaleMutation.mutate(saleData);
  };

  const handlePickupCompletion = () => {
    if (!pickupData.orderId) {
      setError('Order ID is required');
      return;
    }

    pickupCompletionMutation.mutate(pickupData);
  };

  const handleInventoryDeduction = () => {
    if (!deductionData.productId || deductionData.quantity <= 0) {
      setError('Product and quantity are required');
      return;
    }

    inventoryDeductionMutation.mutate(deductionData);
  };

  const resetPaymentData = () => {
    setPaymentData({
      paymentMethod: 'CASH',
      customerName: '',
      customerEmail: '',
      notes: '',
    });
  };

  const resetPickupData = () => {
    setPickupData({
      orderId: '',
      customerInfo: '',
      notes: '',
    });
  };

  const resetDeductionData = () => {
    setDeductionData({
      productId: '',
      quantity: 1,
      reason: '',
      notes: '',
    });
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
          Point of Sale
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<AnalyticsIcon />}
            onClick={() => setStatsDialogOpen(true)}
          >
            View Stats
          </Button>
          <Button
            variant="outlined"
            startIcon={<InventoryIcon />}
            onClick={() => setDeductionDialogOpen(true)}
          >
            Deduct Inventory
          </Button>
          <Button
            variant="outlined"
            startIcon={<ReceiptIcon />}
            onClick={() => setPickupDialogOpen(true)}
          >
            Complete Pickup
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

      {/* Daily Summary Card */}
      {dailySummary && (
        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>
              <TodayIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
              Today's Summary
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={3}>
                <Typography variant="body2" color="text.secondary">Total Sales</Typography>
                <Typography variant="h6">{formatCurrency(dailySummary.totalSales || 0)}</Typography>
              </Grid>
              <Grid item xs={12} sm={3}>
                <Typography variant="body2" color="text.secondary">Transactions</Typography>
                <Typography variant="h6">{dailySummary.transactionCount || 0}</Typography>
              </Grid>
              <Grid item xs={12} sm={3}>
                <Typography variant="body2" color="text.secondary">Items Sold</Typography>
                <Typography variant="h6">{dailySummary.itemsSold || 0}</Typography>
              </Grid>
              <Grid item xs={12} sm={3}>
                <Typography variant="body2" color="text.secondary">Avg. Transaction</Typography>
                <Typography variant="h6">{formatCurrency(dailySummary.averageTransaction || 0)}</Typography>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      )}

      <Grid container spacing={3}>
        {/* Product Selection */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Product Selection
            </Typography>
            <TextField
              fullWidth
              label="Search Products"
              value={productSearchTerm}
              onChange={(e) => setProductSearchTerm(e.target.value)}
              sx={{ mb: 2 }}
            />
            
            {productsLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
                <CircularProgress />
              </Box>
            ) : (
              <Grid container spacing={2}>
                {products.map((product) => (
                  <Grid item xs={12} sm={6} md={4} key={product.id}>
                    <Card 
                      sx={{ 
                        cursor: 'pointer',
                        '&:hover': { boxShadow: 2 },
                        height: '100%',
                      }}
                      onClick={() => addToCart(product)}
                    >
                      <CardContent>
                        <Typography variant="h6" component="h3" noWrap>
                          {product.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" noWrap>
                          SKU: {product.sku}
                        </Typography>
                        <Typography variant="h6" color="primary" sx={{ mt: 1 }}>
                          {formatCurrency(product.price)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Stock: {product.stockQuantity || 0}
                        </Typography>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            )}
          </Paper>
        </Grid>

        {/* Shopping Cart */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2, position: 'sticky', top: 20 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6">
                <CartIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
                Cart ({cartItemCount})
              </Typography>
              {cart.length > 0 && (
                <IconButton size="small" onClick={clearCart} color="error">
                  <ClearIcon />
                </IconButton>
              )}
            </Box>

            {cart.length === 0 ? (
              <Typography variant="body2" color="text.secondary" sx={{ textAlign: 'center', py: 4 }}>
                Cart is empty
              </Typography>
            ) : (
              <>
                <List dense>
                  {cart.map((item) => (
                    <ListItem key={item.id} divider>
                      <ListItemText
                        primary={item.name}
                        secondary={`${formatCurrency(item.price)} each`}
                      />
                      <ListItemSecondaryAction>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <IconButton
                            size="small"
                            onClick={() => updateCartItemQuantity(item.id, item.quantity - 1)}
                          >
                            <RemoveIcon />
                          </IconButton>
                          <Typography variant="body2" sx={{ minWidth: 20, textAlign: 'center' }}>
                            {item.quantity}
                          </Typography>
                          <IconButton
                            size="small"
                            onClick={() => updateCartItemQuantity(item.id, item.quantity + 1)}
                          >
                            <AddIcon />
                          </IconButton>
                          <IconButton
                            size="small"
                            onClick={() => removeFromCart(item.id)}
                            color="error"
                          >
                            <DeleteIcon />
                          </IconButton>
                        </Box>
                      </ListItemSecondaryAction>
                    </ListItem>
                  ))}
                </List>

                <Divider sx={{ my: 2 }} />
                
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                  <Typography variant="h6">Total:</Typography>
                  <Typography variant="h6" color="primary">
                    {formatCurrency(cartTotal)}
                  </Typography>
                </Box>

                <Button
                  fullWidth
                  variant="contained"
                  size="large"
                  startIcon={<PaymentIcon />}
                  onClick={() => setPaymentDialogOpen(true)}
                  disabled={quickSaleMutation.isLoading}
                >
                  {quickSaleMutation.isLoading ? <CircularProgress size={20} /> : 'Process Payment'}
                </Button>
              </>
            )}
          </Paper>
        </Grid>
      </Grid>

      {/* Payment Dialog */}
      <Dialog open={paymentDialogOpen} onClose={() => setPaymentDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Process Payment</DialogTitle>
        <DialogContent>
          <Typography variant="h6" sx={{ mb: 2 }}>
            Total: {formatCurrency(cartTotal)}
          </Typography>
          
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Payment Method</InputLabel>
                <Select
                  value={paymentData.paymentMethod}
                  label="Payment Method"
                  onChange={(e) => setPaymentData({ ...paymentData, paymentMethod: e.target.value })}
                >
                  <MenuItem value="CASH">Cash</MenuItem>
                  <MenuItem value="CARD">Card</MenuItem>
                  <MenuItem value="DIGITAL">Digital</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Customer Name"
                value={paymentData.customerName}
                onChange={(e) => setPaymentData({ ...paymentData, customerName: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Customer Email"
                type="email"
                value={paymentData.customerEmail}
                onChange={(e) => setPaymentData({ ...paymentData, customerEmail: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={2}
                label="Notes"
                value={paymentData.notes}
                onChange={(e) => setPaymentData({ ...paymentData, notes: e.target.value })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setPaymentDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleQuickSale} 
            variant="contained"
            disabled={quickSaleMutation.isLoading}
          >
            {quickSaleMutation.isLoading ? <CircularProgress size={20} /> : 'Complete Sale'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Pickup Completion Dialog */}
      <Dialog open={pickupDialogOpen} onClose={() => setPickupDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Complete Pickup</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Order ID"
                value={pickupData.orderId}
                onChange={(e) => setPickupData({ ...pickupData, orderId: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Customer Information"
                value={pickupData.customerInfo}
                onChange={(e) => setPickupData({ ...pickupData, customerInfo: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={2}
                label="Notes"
                value={pickupData.notes}
                onChange={(e) => setPickupData({ ...pickupData, notes: e.target.value })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setPickupDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handlePickupCompletion} 
            variant="contained"
            disabled={!pickupData.orderId || pickupCompletionMutation.isLoading}
          >
            {pickupCompletionMutation.isLoading ? <CircularProgress size={20} /> : 'Complete Pickup'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Inventory Deduction Dialog */}
      <Dialog open={deductionDialogOpen} onClose={() => setDeductionDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Deduct Inventory</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Product</InputLabel>
                <Select
                  value={deductionData.productId}
                  label="Product"
                  onChange={(e) => setDeductionData({ ...deductionData, productId: e.target.value })}
                >
                  {products.map((product) => (
                    <MenuItem key={product.id} value={product.id}>
                      {product.name} (Stock: {product.stockQuantity})
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="number"
                label="Quantity"
                value={deductionData.quantity}
                onChange={(e) => setDeductionData({ ...deductionData, quantity: parseInt(e.target.value) || 0 })}
                inputProps={{ min: 1 }}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Reason"
                value={deductionData.reason}
                onChange={(e) => setDeductionData({ ...deductionData, reason: e.target.value })}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={2}
                label="Notes"
                value={deductionData.notes}
                onChange={(e) => setDeductionData({ ...deductionData, notes: e.target.value })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeductionDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleInventoryDeduction} 
            variant="contained"
            disabled={!deductionData.productId || deductionData.quantity <= 0 || inventoryDeductionMutation.isLoading}
          >
            {inventoryDeductionMutation.isLoading ? <CircularProgress size={20} /> : 'Deduct Inventory'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Sales Statistics Dialog */}
      <Dialog open={statsDialogOpen} onClose={() => setStatsDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Sales Statistics</DialogTitle>
        <DialogContent>
          {statsLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : salesStats ? (
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">Weekly Sales</Typography>
                    <Typography variant="h4" color="primary">
                      {formatCurrency(salesStats.weeklySales || 0)}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">Monthly Sales</Typography>
                    <Typography variant="h4" color="primary">
                      {formatCurrency(salesStats.monthlySales || 0)}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">Total Transactions</Typography>
                    <Typography variant="h4" color="primary">
                      {salesStats.totalTransactions || 0}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6">Top Product</Typography>
                    <Typography variant="h6" color="primary">
                      {salesStats.topProduct || 'N/A'}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          ) : (
            <Typography variant="body2" color="text.secondary" sx={{ py: 4, textAlign: 'center' }}>
              No statistics available
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setStatsDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default POSPage;
