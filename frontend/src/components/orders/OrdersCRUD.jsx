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
  Typography,
  Chip,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from '@mui/material';
import { Add as AddIcon, Edit as EditIcon, Delete as DeleteIcon, Visibility as ViewIcon } from '@mui/icons-material';
import { orderService, userService, productService } from '../../services';

const OrdersCRUD = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [error, setError] = useState('');
  const [users, setUsers] = useState([]);
  const [products, setProducts] = useState([]);

  const statusColors = {
    pending: 'warning',
    processing: 'info',
    shipped: 'primary',
    delivered: 'success',
    cancelled: 'error',
  };

  const fetchData = async () => {
    try {
      setLoading(true);
      const [ordersData, usersData, productsData] = await Promise.all([
        orderService.getOrders(),
        userService.getUsers(),
        productService.getProducts(),
      ]);
      
      setOrders(ordersData.data || []);
      setUsers(usersData.data || []);
      setProducts(productsData.data || []);
    } catch (err) {
      setError('Failed to fetch data');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleOpen = (order = null) => {
    setSelectedOrder(order);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setSelectedOrder(null);
  };

  const handleStatusChange = async (orderId, status) => {
    try {
      await orderService.updateOrderStatus(orderId, status);
      await fetchData();
    } catch (err) {
      setError('Failed to update order status');
      console.error(err);
    }
  };

  const handleSubmit = async (orderData) => {
    try {
      if (selectedOrder) {
        await orderService.updateOrder(selectedOrder.id, orderData);
      } else {
        await orderService.createOrder(orderData);
      }
      await fetchData();
      handleClose();
    } catch (err) {
      setError(`Failed to ${selectedOrder ? 'update' : 'create'} order`);
      console.error(err);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this order?')) {
      try {
        await orderService.deleteOrder(id);
        await fetchData();
      } catch (err) {
        setError('Failed to delete order');
        console.error(err);
      }
    }
  };

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Typography color="error">{error}</Typography>;

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" mb={2}>
        <Typography variant="h4">Orders</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => handleOpen()}
        >
          Create Order
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Order #</TableCell>
              <TableCell>Customer</TableCell>
              <TableCell>Products</TableCell>
              <TableCell>Total</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Date</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {orders.map((order) => (
              <TableRow key={order.id}>
                <TableCell>#{order.orderNumber}</TableCell>
                <TableCell>
                  {users.find(u => u.id === order.userId)?.name || 'Unknown User'}
                </TableCell>
                <TableCell>
                  {order.items?.slice(0, 2).map((item, idx) => (
                    <div key={idx}>
                      {products.find(p => p.id === item.productId)?.name || 'Unknown Product'}
                      {item.quantity > 1 && ` × ${item.quantity}`}
                    </div>
                  ))}
                  {order.items?.length > 2 && `+${order.items.length - 2} more`}
                </TableCell>
                <TableCell>${order.total?.toFixed(2)}</TableCell>
                <TableCell>
                  <FormControl size="small" variant="outlined" fullWidth>
                    <Select
                      value={order.status}
                      onChange={(e) => handleStatusChange(order.id, e.target.value)}
                      sx={{ minWidth: 120 }}
                    >
                      {Object.keys(statusColors).map((status) => (
                        <MenuItem key={status} value={status}>
                          {status.charAt(0).toUpperCase() + status.slice(1)}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </TableCell>
                <TableCell>{new Date(order.createdAt).toLocaleDateString()}</TableCell>
                <TableCell>
                  <IconButton onClick={() => handleOpen(order)}>
                    <EditIcon />
                  </IconButton>
                  <IconButton onClick={() => handleDelete(order.id)}>
                    <DeleteIcon color="error" />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
};

export default OrdersCRUD;
