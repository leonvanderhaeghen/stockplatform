import React, { useState, useEffect } from 'react';
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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Divider,
  InputAdornment,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  CircularProgress,
  ListItem,
  ListItemText,
  Autocomplete,
} from '@mui/material';
import {
  PointOfSale,
  Search,
  Add,
  Remove,
  Delete,
  ShoppingCart,
  Payment,
  Print,
  Person,
  Clear,
  Refresh,
  Receipt,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSnackbar } from 'notistack';
import { format } from 'date-fns';
import posService from '../../services/posService';
import productService from '../../services/productService';
import { useAuth } from '../../hooks/useAuth';

const POSPage = () => {
  const { user } = useAuth();
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();

  // POS State
  const [cart, setCart] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [quantity, setQuantity] = useState(1);
  const [customerInfo, setCustomerInfo] = useState({ name: '', email: '', phone: '' });
  const [paymentMethod, setPaymentMethod] = useState('cash');
  const [amountReceived, setAmountReceived] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  const [showPaymentDialog, setShowPaymentDialog] = useState(false);
  const [showReceiptDialog, setShowReceiptDialog] = useState(false);
  const [lastTransaction, setLastTransaction] = useState(null);
  const [productSuggestions, setProductSuggestions] = useState([]);

  // Fetch products for autocomplete
  const { data: products = [] } = useQuery({
    queryKey: ['products', 'pos'],
    queryFn: () => productService.getProducts({ page: 1, limit: 100, status: 'active' }),
    select: (data) => data.products || [],
  });

  // Calculate totals
  const subtotal = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  const tax = subtotal * 0.1; // 10% tax
  const total = subtotal + tax;
  const change = amountReceived ? Math.max(0, parseFloat(amountReceived) - total) : 0;

  // Search for products
  useEffect(() => {
    if (searchTerm.length > 2) {
      const filtered = products.filter(product => 
        product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        product.sku?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        product.description?.toLowerCase().includes(searchTerm.toLowerCase())
      );
      setProductSuggestions(filtered.slice(0, 10));
    } else {
      setProductSuggestions([]);
    }
  }, [searchTerm, products]);

  // Add item to cart
  const addToCart = async (product, qty = quantity) => {
    try {
      // Check inventory availability
      const availability = await posService.checkInventoryAvailability({
        productId: product.id,
        quantity: qty
      });

      if (!availability.available) {
        enqueueSnackbar(`Insufficient inventory. Available: ${availability.availableQuantity}`, { variant: 'warning' });
        return;
      }

      const existingItem = cart.find(item => item.id === product.id);
      if (existingItem) {
        const newQuantity = existingItem.quantity + qty;
        if (newQuantity <= availability.availableQuantity) {
          setCart(cart.map(item => 
            item.id === product.id 
              ? { ...item, quantity: newQuantity }
              : item
          ));
        } else {
          enqueueSnackbar(`Cannot add more. Max available: ${availability.availableQuantity}`, { variant: 'warning' });
        }
      } else {
        setCart([...cart, {
          ...product,
          quantity: qty,
          availableQuantity: availability.availableQuantity
        }]);
      }

      setSelectedProduct(null);
      setQuantity(1);
      setSearchTerm('');
      setProductSuggestions([]);
      enqueueSnackbar(`Added ${product.name} to cart`, { variant: 'success' });
    } catch (error) {
      enqueueSnackbar('Failed to check inventory', { variant: 'error' });
    }
  };

  // Update item quantity in cart
  const updateCartQuantity = (productId, newQuantity) => {
    if (newQuantity <= 0) {
      removeFromCart(productId);
      return;
    }

    const item = cart.find(item => item.id === productId);
    if (item && newQuantity <= item.availableQuantity) {
      setCart(cart.map(item => 
        item.id === productId 
          ? { ...item, quantity: newQuantity }
          : item
      ));
    } else {
      enqueueSnackbar(`Max available: ${item?.availableQuantity}`, { variant: 'warning' });
    }
  };

  // Remove item from cart
  const removeFromCart = (productId) => {
    setCart(cart.filter(item => item.id !== productId));
    enqueueSnackbar('Item removed from cart', { variant: 'info' });
  };

  // Clear cart
  const clearCart = () => {
    setCart([]);
    setCustomerInfo({ name: '', email: '', phone: '' });
    setAmountReceived('');
    enqueueSnackbar('Cart cleared', { variant: 'info' });
  };

  // Process transaction
  const processTransactionMutation = useMutation({
    mutationFn: async (transactionData) => {
      return await posService.processTransaction(transactionData);
    },
    onSuccess: (data) => {
      setLastTransaction(data);
      setShowReceiptDialog(true);
      clearCart();
      setIsProcessing(false);
      setShowPaymentDialog(false);
      enqueueSnackbar('Transaction processed successfully!', { variant: 'success' });
      queryClient.invalidateQueries(['inventory']);
    },
    onError: (error) => {
      setIsProcessing(false);
      enqueueSnackbar(error.message || 'Transaction failed', { variant: 'error' });
    },
  });

  // Handle payment processing
  const handlePayment = () => {
    if (cart.length === 0) {
      enqueueSnackbar('Cart is empty', { variant: 'warning' });
      return;
    }

    if (paymentMethod === 'cash' && (!amountReceived || parseFloat(amountReceived) < total)) {
      enqueueSnackbar('Insufficient payment amount', { variant: 'warning' });
      return;
    }

    setIsProcessing(true);

    const transactionData = {
      items: cart.map(item => ({
        productId: item.id,
        sku: item.sku,
        name: item.name,
        price: item.price,
        quantity: item.quantity
      })),
      customerInfo: customerInfo.name ? customerInfo : null,
      paymentMethod,
      amountReceived: paymentMethod === 'cash' ? parseFloat(amountReceived) : total,
      subtotal,
      tax,
      total,
      change: paymentMethod === 'cash' ? change : 0,
      cashierId: user.id,
      cashierName: user.name
    };

    processTransactionMutation.mutate(transactionData);
  };

  // Print receipt
  const printReceipt = () => {
    if (!lastTransaction) return;

    const receiptContent = `
      STOCKPLATFORM POS
      ==================
      Date: ${format(new Date(lastTransaction.timestamp), 'yyyy-MM-dd HH:mm:ss')}
      Transaction ID: ${lastTransaction.id}
      Cashier: ${lastTransaction.cashierName}
      
      ${lastTransaction.customerInfo ? `Customer: ${lastTransaction.customerInfo.name}\n` : ''}
      
      ITEMS:
      ${lastTransaction.items.map(item => 
        `${item.name} x${item.quantity} @ $${item.price.toFixed(2)} = $${(item.price * item.quantity).toFixed(2)}`
      ).join('\n')}
      
      ==================
      Subtotal: $${lastTransaction.subtotal.toFixed(2)}
      Tax (10%): $${lastTransaction.tax.toFixed(2)}
      TOTAL: $${lastTransaction.total.toFixed(2)}
      
      ${lastTransaction.paymentMethod === 'cash' ? 
        `Cash Received: $${lastTransaction.amountReceived.toFixed(2)}\nChange: $${lastTransaction.change.toFixed(2)}` :
        `Payment Method: ${lastTransaction.paymentMethod.toUpperCase()}`
      }
      
      Thank you for your business!
    `;

    const printWindow = window.open('', '_blank');
    printWindow.document.write(`<pre>${receiptContent}</pre>`);
    printWindow.document.close();
    printWindow.print();
  };

  return (
    <Box>
      {/* Header */}
      <Paper sx={{ p: 3, mb: 3, bgcolor: 'primary.main', color: 'primary.contrastText' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <PointOfSale sx={{ fontSize: 40, mr: 2 }} />
            <Box>
              <Typography variant="h4" gutterBottom>
                Point of Sale
              </Typography>
              <Typography variant="subtitle1">
                Process in-store transactions • Cashier: {user?.name}
              </Typography>
            </Box>
          </Box>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              variant="contained"
              color="secondary"
              startIcon={<Clear />}
              onClick={clearCart}
              disabled={cart.length === 0}
            >
              Clear Cart
            </Button>
            <Button
              variant="contained"
              color="secondary"
              startIcon={<Refresh />}
              onClick={() => queryClient.invalidateQueries(['products'])}
            >
              Refresh
            </Button>
          </Box>
        </Box>
      </Paper>

      <Grid container spacing={3}>
        {/* Product Search */}
        <Grid item xs={12} md={8}>
          <Card sx={{ mb: 3 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Product Search
              </Typography>
              <Box sx={{ position: 'relative' }}>
                <Autocomplete
                  options={productSuggestions}
                  getOptionLabel={(option) => `${option.name} (${option.sku}) - $${option.price}`}
                  value={selectedProduct}
                  onChange={(event, newValue) => {
                    setSelectedProduct(newValue);
                    if (newValue) {
                      setSearchTerm(newValue.name);
                    }
                  }}
                  inputValue={searchTerm}
                  onInputChange={(event, newInputValue) => {
                    setSearchTerm(newInputValue);
                  }}
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      fullWidth
                      placeholder="Search by name, SKU, or scan barcode..."
                      InputProps={{
                        ...params.InputProps,
                        startAdornment: (
                          <InputAdornment position="start">
                            <Search />
                          </InputAdornment>
                        ),
                      }}
                    />
                  )}
                  renderOption={(props, option) => (
                    <ListItem {...props} key={option.id}>
                      <ListItemText
                        primary={`${option.name} - $${option.price}`}
                        secondary={`SKU: ${option.sku} | Stock: ${option.stockQuantity || 'N/A'}`}
                      />
                    </ListItem>
                  )}
                />
              </Box>
              
              {selectedProduct && (
                <Box sx={{ mt: 2, display: 'flex', gap: 2, alignItems: 'center' }}>
                  <TextField
                    type="number"
                    label="Quantity"
                    value={quantity}
                    onChange={(e) => setQuantity(parseInt(e.target.value) || 1)}
                    sx={{ width: 120 }}
                    inputProps={{ min: 1 }}
                  />
                  <Button
                    variant="contained"
                    startIcon={<Add />}
                    onClick={() => addToCart(selectedProduct, quantity)}
                  >
                    Add to Cart
                  </Button>
                </Box>
              )}
            </CardContent>
          </Card>

          {/* Shopping Cart */}
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
                <ShoppingCart sx={{ mr: 1 }} />
                Shopping Cart ({cart.length} items)
              </Typography>
              
              {cart.length === 0 ? (
                <Typography color="text.secondary" sx={{ textAlign: 'center', py: 4 }}>
                  Cart is empty. Add products to start a transaction.
                </Typography>
              ) : (
                <TableContainer>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Product</TableCell>
                        <TableCell align="center">Price</TableCell>
                        <TableCell align="center">Quantity</TableCell>
                        <TableCell align="right">Total</TableCell>
                        <TableCell align="center">Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {cart.map((item) => (
                        <TableRow key={item.id}>
                          <TableCell>
                            <Box>
                              <Typography variant="body1">{item.name}</Typography>
                              <Typography variant="caption" color="text.secondary">
                                SKU: {item.sku}
                              </Typography>
                            </Box>
                          </TableCell>
                          <TableCell align="center">${item.price.toFixed(2)}</TableCell>
                          <TableCell align="center">
                            <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                              <IconButton
                                size="small"
                                onClick={() => updateCartQuantity(item.id, item.quantity - 1)}
                              >
                                <Remove />
                              </IconButton>
                              <Typography sx={{ mx: 2, minWidth: 30, textAlign: 'center' }}>
                                {item.quantity}
                              </Typography>
                              <IconButton
                                size="small"
                                onClick={() => updateCartQuantity(item.id, item.quantity + 1)}
                              >
                                <Add />
                              </IconButton>
                            </Box>
                          </TableCell>
                          <TableCell align="right">
                            ${(item.price * item.quantity).toFixed(2)}
                          </TableCell>
                          <TableCell align="center">
                            <IconButton
                              color="error"
                              onClick={() => removeFromCart(item.id)}
                            >
                              <Delete />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Transaction Summary & Customer Info */}
        <Grid item xs={12} md={4}>
          {/* Customer Information */}
          <Card sx={{ mb: 3 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom sx={{ display: 'flex', alignItems: 'center' }}>
                <Person sx={{ mr: 1 }} />
                Customer Information
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                <TextField
                  label="Name (Optional)"
                  value={customerInfo.name}
                  onChange={(e) => setCustomerInfo({ ...customerInfo, name: e.target.value })}
                  fullWidth
                />
                <TextField
                  label="Email (Optional)"
                  type="email"
                  value={customerInfo.email}
                  onChange={(e) => setCustomerInfo({ ...customerInfo, email: e.target.value })}
                  fullWidth
                />
                <TextField
                  label="Phone (Optional)"
                  value={customerInfo.phone}
                  onChange={(e) => setCustomerInfo({ ...customerInfo, phone: e.target.value })}
                  fullWidth
                />
              </Box>
            </CardContent>
          </Card>

          {/* Transaction Summary */}
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Transaction Summary
              </Typography>
              
              <Box sx={{ mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                  <Typography>Subtotal:</Typography>
                  <Typography>${subtotal.toFixed(2)}</Typography>
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                  <Typography>Tax (10%):</Typography>
                  <Typography>${tax.toFixed(2)}</Typography>
                </Box>
                <Divider sx={{ my: 1 }} />
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                  <Typography variant="h6">Total:</Typography>
                  <Typography variant="h6" color="primary">
                    ${total.toFixed(2)}
                  </Typography>
                </Box>
              </Box>

              <Button
                variant="contained"
                fullWidth
                size="large"
                startIcon={<Payment />}
                onClick={() => setShowPaymentDialog(true)}
                disabled={cart.length === 0}
                sx={{ mb: 2 }}
              >
                Process Payment
              </Button>
              
              <Typography variant="caption" color="text.secondary" sx={{ textAlign: 'center', display: 'block' }}>
                {cart.length} item(s) in cart
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Payment Dialog */}
      <Dialog open={showPaymentDialog} onClose={() => setShowPaymentDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Process Payment</DialogTitle>
        <DialogContent>
          <Box sx={{ mb: 3 }}>
            <Typography variant="h6" color="primary">
              Total Amount: ${total.toFixed(2)}
            </Typography>
          </Box>
          
          <FormControl fullWidth sx={{ mb: 3 }}>
            <InputLabel>Payment Method</InputLabel>
            <Select
              value={paymentMethod}
              onChange={(e) => setPaymentMethod(e.target.value)}
              label="Payment Method"
            >
              <MenuItem value="cash">Cash</MenuItem>
              <MenuItem value="card">Credit/Debit Card</MenuItem>
              <MenuItem value="mobile">Mobile Payment</MenuItem>
            </Select>
          </FormControl>

          {paymentMethod === 'cash' && (
            <>
              <TextField
                label="Amount Received"
                type="number"
                value={amountReceived}
                onChange={(e) => setAmountReceived(e.target.value)}
                fullWidth
                sx={{ mb: 2 }}
                inputProps={{ step: 0.01, min: 0 }}
                InputProps={{
                  startAdornment: <InputAdornment position="start">$</InputAdornment>,
                }}
              />
              {amountReceived && (
                <Alert severity={change >= 0 ? 'success' : 'error'} sx={{ mb: 2 }}>
                  Change to give: ${change.toFixed(2)}
                </Alert>
              )}
            </>
          )}

          {paymentMethod !== 'cash' && (
            <Alert severity="info" sx={{ mb: 2 }}>
              {paymentMethod === 'card' ? 'Insert or swipe card to process payment' : 'Present mobile device for payment'}
            </Alert>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowPaymentDialog(false)}>Cancel</Button>
          <Button
            variant="contained"
            onClick={handlePayment}
            disabled={isProcessing || (paymentMethod === 'cash' && (!amountReceived || parseFloat(amountReceived) < total))}
            startIcon={isProcessing ? <CircularProgress size={20} /> : <Payment />}
          >
            {isProcessing ? 'Processing...' : 'Complete Transaction'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Receipt Dialog */}
      <Dialog open={showReceiptDialog} onClose={() => setShowReceiptDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ display: 'flex', alignItems: 'center' }}>
          <Receipt sx={{ mr: 1 }} />
          Transaction Complete
        </DialogTitle>
        <DialogContent>
          {lastTransaction && (
            <Box>
              <Alert severity="success" sx={{ mb: 3 }}>
                Transaction processed successfully!
              </Alert>
              
              <Typography variant="h6" gutterBottom>
                Transaction ID: {lastTransaction.id}
              </Typography>
              
              <Typography variant="body2" color="text.secondary" gutterBottom>
                {format(new Date(lastTransaction.timestamp), 'yyyy-MM-dd HH:mm:ss')}
              </Typography>
              
              <Box sx={{ mt: 2, mb: 2 }}>
                <Typography variant="subtitle1" gutterBottom>Items:</Typography>
                {lastTransaction.items?.map((item, index) => (
                  <Box key={index} sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography>{item.name} x{item.quantity}</Typography>
                    <Typography>${(item.price * item.quantity).toFixed(2)}</Typography>
                  </Box>
                ))}
              </Box>
              
              <Divider sx={{ my: 2 }} />
              
              <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                <Typography variant="h6">Total:</Typography>
                <Typography variant="h6" color="primary">
                  ${lastTransaction.total?.toFixed(2)}
                </Typography>
              </Box>
              
              {lastTransaction.paymentMethod === 'cash' && (
                <Box sx={{ mt: 2 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography>Cash Received:</Typography>
                    <Typography>${lastTransaction.amountReceived?.toFixed(2)}</Typography>
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography>Change:</Typography>
                    <Typography>${lastTransaction.change?.toFixed(2)}</Typography>
                  </Box>
                </Box>
              )}
            </Box>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowReceiptDialog(false)}>Close</Button>
          <Button
            variant="contained"
            startIcon={<Print />}
            onClick={printReceipt}
          >
            Print Receipt
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default POSPage;
