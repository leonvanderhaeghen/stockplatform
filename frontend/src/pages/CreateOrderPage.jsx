import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  MenuItem,
  Grid,
  Snackbar,
  Alert,
  CircularProgress,
} from '@mui/material';
import productService from '../services/productService';
import orderService from '../services/orderService';

const CreateOrderPage = () => {
  const [products, setProducts] = useState([]);
  const [loadingProducts, setLoadingProducts] = useState(false);
  const [saving, setSaving] = useState(false);
  const [snack, setSnack] = useState({ open: false, message: '', severity: 'success' });

  const [form, setForm] = useState({
    productId: '',
    quantity: 1,
    addressId: '',
    paymentType: 'CARD',
    shippingType: 'STANDARD',
  });

  // Fetch product list once
  useEffect(() => {
    const load = async () => {
      try {
        setLoadingProducts(true);
        const list = await productService.getProducts({ limit: 100 });
        setProducts(list);
      } catch (e) {
        console.error(e);
      } finally {
        setLoadingProducts(false);
      }
    };
    load();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async () => {
    if (!form.productId || !form.addressId) return;
    setSaving(true);
    try {
      const prod = products.find((p) => p.id === form.productId);
      const body = {
        items: [
          {
            productId: form.productId,
            sku: prod?.sku || '',
            quantity: Number(form.quantity),
            price: prod?.price || 0,
          },
        ],
        addressId: form.addressId,
        paymentType: form.paymentType,
        shippingType: form.shippingType,
      };
      await orderService.createOrder(body);
      setSnack({ open: true, message: 'Order placed successfully', severity: 'success' });
      setForm({ ...form, quantity: 1 });
    } catch (err) {
      console.error(err);
      setSnack({ open: true, message: 'Failed to place order', severity: 'error' });
    } finally {
      setSaving(false);
    }
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Create Order
      </Typography>
      <Paper sx={{ p: 3 }}>
        {loadingProducts ? (
          <CircularProgress />
        ) : (
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <TextField
                select
                fullWidth
                label="Product"
                name="productId"
                value={form.productId}
                onChange={handleChange}
              >
                {products.map((p) => (
                  <MenuItem key={p.id} value={p.id}>
                    {p.name}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12} md={3}>
              <TextField
                type="number"
                label="Quantity"
                name="quantity"
                value={form.quantity}
                onChange={handleChange}
                fullWidth
                inputProps={{ min: 1 }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Address ID"
                name="addressId"
                value={form.addressId}
                onChange={handleChange}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                select
                fullWidth
                label="Payment Type"
                name="paymentType"
                value={form.paymentType}
                onChange={handleChange}
              >
                <MenuItem value="CARD">Card</MenuItem>
                <MenuItem value="CASH">Cash</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                select
                fullWidth
                label="Shipping"
                name="shippingType"
                value={form.shippingType}
                onChange={handleChange}
              >
                <MenuItem value="STANDARD">Standard</MenuItem>
                <MenuItem value="EXPRESS">Express</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12}>
              <Button variant="contained" onClick={handleSubmit} disabled={saving}>
                {saving ? 'Saving…' : 'Place Order'}
              </Button>
            </Grid>
          </Grid>
        )}
      </Paper>
      <Snackbar
        open={snack.open}
        onClose={() => setSnack({ ...snack, open: false })}
        autoHideDuration={4000}
      >
        <Alert severity={snack.severity} onClose={() => setSnack({ ...snack, open: false })}>
          {snack.message}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default CreateOrderPage;