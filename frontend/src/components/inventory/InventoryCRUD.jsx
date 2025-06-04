import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  InputAdornment,
  Divider,
} from '@mui/material';
import { 
  Add as AddIcon, 
  Edit as EditIcon, 
  Delete as DeleteIcon,
  AddCircle as AddCircleIcon,
  RemoveCircle as RemoveCircleIcon,
} from '@mui/icons-material';
import { inventoryService, productService } from '../../services';

const InventoryForm = ({ initialValues, onSubmit, onCancel, loading = false }) => {
  const [formData, setFormData] = useState({
    productId: '',
    quantity: 0,
    location: '',
    reorderThreshold: 10,
    ...initialValues,
  });
  const [products, setProducts] = useState([]);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        const { data } = await productService.getProducts();
        setProducts(data || []);
      } catch (err) {
        console.error('Failed to fetch products', err);
      }
    };
    fetchProducts();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleQuantityChange = (amount) => {
    setFormData(prev => ({
      ...prev,
      quantity: Math.max(0, (parseInt(prev.quantity) || 0) + amount),
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit({
      ...formData,
      quantity: parseInt(formData.quantity) || 0,
      reorderThreshold: parseInt(formData.reorderThreshold) || 10,
    });
  };

  return (
    <form onSubmit={handleSubmit}>
      <DialogContent>
        <FormControl fullWidth margin="normal">
          <InputLabel>Product</InputLabel>
          <Select
            name="productId"
            value={formData.productId}
            label="Product"
            onChange={handleChange}
            required
            disabled={!!initialValues?.id}
          >
            {products.map((product) => (
              <MenuItem key={product.id} value={product.id}>
                {product.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        
        <TextField
          margin="normal"
          required
          fullWidth
          label="Location"
          name="location"
          value={formData.location}
          onChange={handleChange}
        />

        <Box sx={{ display: 'flex', alignItems: 'center', mt: 2, mb: 2 }}>
          <Typography variant="subtitle1" sx={{ mr: 2, minWidth: 120 }}>
            Quantity:
          </Typography>
          <IconButton 
            onClick={() => handleQuantityChange(-1)}
            color="error"
            disabled={formData.quantity <= 0}
          >
            <RemoveCircleIcon />
          </IconButton>
          <TextField
            name="quantity"
            type="number"
            value={formData.quantity}
            onChange={handleChange}
            inputProps={{ min: 0, style: { textAlign: 'center' } }}
            sx={{ width: 80, mx: 1 }}
          />
          <IconButton 
            onClick={() => handleQuantityChange(1)}
            color="primary"
          >
            <AddCircleIcon />
          </IconButton>
        </Box>

        <TextField
          margin="normal"
          required
          fullWidth
          label="Reorder Threshold"
          name="reorderThreshold"
          type="number"
          value={formData.reorderThreshold}
          onChange={handleChange}
          inputProps={{ min: 1 }}
          helperText="Alert when quantity falls below this number"
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onCancel} disabled={loading}>
          Cancel
        </Button>
        <Button type="submit" color="primary" variant="contained" disabled={loading}>
          {loading ? 'Saving...' : 'Save'}
        </Button>
      </DialogActions>
    </form>
  );
};

const InventoryCRUD = () => {
  const [inventory, setInventory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState(null);
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState(false);
  const [products, setProducts] = useState({});

  const fetchData = async () => {
    try {
      setLoading(true);
      const [inventoryData, productsData] = await Promise.all([
        inventoryService.getInventoryItems(),
        productService.getProducts(),
      ]);
      
      setInventory(inventoryData.data || []);
      
      // Create a products map for quick lookup
      const productsMap = {};
      (productsData.data || []).forEach(product => {
        productsMap[product.id] = product;
      });
      setProducts(productsMap);
    } catch (err) {
      setError('Failed to fetch inventory data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleOpen = (item = null) => {
    setSelectedItem(item);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setSelectedItem(null);
  };

  const handleSubmit = async (itemData) => {
    try {
      setActionLoading(true);
      if (selectedItem) {
        await inventoryService.updateInventoryItem(selectedItem.id, itemData);
      } else {
        await inventoryService.createInventoryItem(itemData);
      }
      await fetchData();
      handleClose();
    } catch (err) {
      setError(`Failed to ${selectedItem ? 'update' : 'create'} inventory item`);
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this inventory item?')) {
      try {
        setActionLoading(true);
        await inventoryService.deleteInventoryItem(id);
        await fetchData();
      } catch (err) {
        setError('Failed to delete inventory item');
        console.error(err);
      } finally {
        setActionLoading(false);
      }
    }
  };

  const handleAdjustQuantity = async (id, change) => {
    try {
      setActionLoading(true);
      await inventoryService.updateInventoryQuantity(id, change, 'Manual adjustment');
      await fetchData();
    } catch (err) {
      setError('Failed to update quantity');
      console.error(err);
    } finally {
      setActionLoading(false);
    }
  };

  const getStatusColor = (quantity, threshold) => {
    if (quantity === 0) return 'error';
    if (quantity <= threshold) return 'warning';
    return 'success';
  };

  const getStatusText = (quantity, threshold) => {
    if (quantity === 0) return 'Out of Stock';
    if (quantity <= threshold) return 'Low Stock';
    return 'In Stock';
  };

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Typography color="error">{error}</Typography>;

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" mb={2}>
        <Typography variant="h4">Inventory</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => handleOpen()}
          disabled={actionLoading}
        >
          Add Item
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Product</TableCell>
              <TableCell>Location</TableCell>
              <TableCell align="center">Quantity</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Reorder At</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {inventory.map((item) => {
              const product = products[item.productId] || {};
              const statusColor = getStatusColor(item.quantity, item.reorderThreshold);
              const statusText = getStatusText(item.quantity, item.reorderThreshold);
              
              return (
                <TableRow key={item.id}>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <Box>
                        <Typography variant="body1">{product.name || 'Unknown Product'}</Typography>
                        <Typography variant="body2" color="textSecondary">
                          SKU: {product.sku || 'N/A'}
                        </Typography>
                      </Box>
                    </Box>
                  </TableCell>
                  <TableCell>{item.location}</TableCell>
                  <TableCell align="center">
                    <Box display="flex" alignItems="center" justifyContent="center">
                      <IconButton 
                        size="small" 
                        onClick={() => handleAdjustQuantity(item.id, -1)}
                        disabled={actionLoading}
                      >
                        <RemoveCircleIcon color={item.quantity > 0 ? 'error' : 'disabled'} />
                      </IconButton>
                      <Typography sx={{ mx: 1, minWidth: 30, textAlign: 'center' }}>
                        {item.quantity}
                      </Typography>
                      <IconButton 
                        size="small" 
                        onClick={() => handleAdjustQuantity(item.id, 1)}
                        disabled={actionLoading}
                      >
                        <AddCircleIcon color="primary" />
                      </IconButton>
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Chip 
                      label={statusText}
                      color={statusColor}
                      size="small"
                      variant={statusColor === 'success' ? 'outlined' : 'filled'}
                    />
                  </TableCell>
                  <TableCell>{item.reorderThreshold}</TableCell>
                  <TableCell>
                    <IconButton 
                      onClick={() => handleOpen(item)}
                      disabled={actionLoading}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton 
                      onClick={() => handleDelete(item.id)}
                      disabled={actionLoading}
                    >
                      <DeleteIcon color="error" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>
          {selectedItem ? 'Edit Inventory Item' : 'Add New Inventory Item'}
        </DialogTitle>
        <InventoryForm
          initialValues={selectedItem || {}}
          onSubmit={handleSubmit}
          onCancel={handleClose}
          loading={actionLoading}
        />
      </Dialog>
    </Box>
  );
};

export default InventoryCRUD;
