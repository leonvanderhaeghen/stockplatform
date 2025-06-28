import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  Button,
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
  CircularProgress,
  Grid,
  Card,
  CardContent,
} from '@mui/material';
import {
  Add as AddIcon,
  MoreVert as MoreVertIcon,
  Visibility as ViewIcon,
  Edit as EditIcon,
  Cancel as CancelIcon,
  Payment as PaymentIcon,
  LocalShipping as TrackingIcon,
  History as HistoryIcon,
  Refresh as RefreshIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { orderService } from '../services';
import { formatCurrency, formatDate } from '../utils/formatters';

const OrdersPage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedOrder, setSelectedOrder] = useState(null);
  const [anchorEl, setAnchorEl] = useState(null);
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false);
  const [paymentDialogOpen, setPaymentDialogOpen] = useState(false);
  const [trackingDialogOpen, setTrackingDialogOpen] = useState(false);
  const [historyDialogOpen, setHistoryDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [cancelReason, setCancelReason] = useState('');
  const [paymentData, setPaymentData] = useState({
    paymentMethod: '',
    paymentStatus: '',
    transactionId: '',
    notes: '',
  });
  const [trackingData, setTrackingData] = useState({
    trackingCode: '',
    carrier: '',
    estimatedDelivery: '',
    notes: '',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Fetch orders
  const { data: ordersData, isLoading, refetch } = useQuery({
    queryKey: ['orders', page, rowsPerPage],
    queryFn: () => orderService.getOrders({
      page: page + 1,
      limit: rowsPerPage,
    }),
    keepPreviousData: true,
  });

  // Fetch order history
  const { data: orderHistory, isLoading: historyLoading } = useQuery({
    queryKey: ['order-history', selectedOrder?.id],
    queryFn: () => orderService.getOrderHistory(selectedOrder.id),
    enabled: !!selectedOrder && historyDialogOpen,
  });

  // Cancel order mutation
  const cancelOrderMutation = useMutation({
    mutationFn: ({ orderId, cancelData }) => orderService.cancelOrder(orderId, cancelData),
    onSuccess: () => {
      setSuccess('Order cancelled successfully');
      queryClient.invalidateQueries(['orders']);
      setCancelDialogOpen(false);
      setCancelReason('');
      setSelectedOrder(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to cancel order');
    },
  });

  // Update payment mutation
  const updatePaymentMutation = useMutation({
    mutationFn: ({ orderId, paymentData }) => orderService.updateOrderPayment(orderId, paymentData),
    onSuccess: () => {
      setSuccess('Payment information updated successfully');
      queryClient.invalidateQueries(['orders']);
      setPaymentDialogOpen(false);
      setPaymentData({ paymentMethod: '', paymentStatus: '', transactionId: '', notes: '' });
      setSelectedOrder(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to update payment');
    },
  });

  // Add tracking mutation
  const addTrackingMutation = useMutation({
    mutationFn: ({ orderId, trackingData }) => orderService.addOrderTracking(orderId, trackingData),
    onSuccess: () => {
      setSuccess('Tracking information added successfully');
      queryClient.invalidateQueries(['orders']);
      setTrackingDialogOpen(false);
      setTrackingData({ trackingCode: '', carrier: '', estimatedDelivery: '', notes: '' });
      setSelectedOrder(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to add tracking');
    },
  });

  // Delete order mutation
  const deleteOrderMutation = useMutation({
    mutationFn: (orderId) => orderService.deleteOrder(orderId),
    onSuccess: () => {
      setSuccess('Order deleted successfully');
      queryClient.invalidateQueries(['orders']);
      setDeleteDialogOpen(false);
      setSelectedOrder(null);
    },
    onError: (error) => {
      setError(error.response?.data?.message || 'Failed to delete order');
    },
  });

  const orders = ordersData?.data || [];
  const totalCount = ordersData?.total || 0;

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleMenuClick = (event, order) => {
    setAnchorEl(event.currentTarget);
    setSelectedOrder(order);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedOrder(null);
  };

  const handleCancelOrder = () => {
    if (cancelReason.trim()) {
      cancelOrderMutation.mutate({
        orderId: selectedOrder.id,
        cancelData: { reason: cancelReason, notes: `Cancelled by user` },
      });
    }
  };

  const handleUpdatePayment = () => {
    if (paymentData.paymentMethod && paymentData.paymentStatus) {
      updatePaymentMutation.mutate({
        orderId: selectedOrder.id,
        paymentData,
      });
    }
  };

  const handleAddTracking = () => {
    if (trackingData.trackingCode) {
      addTrackingMutation.mutate({
        orderId: selectedOrder.id,
        trackingData,
      });
    }
  };

  const handleDeleteOrder = () => {
    deleteOrderMutation.mutate(selectedOrder.id);
  };

  const getStatusColor = (status) => {
    switch (status?.toLowerCase()) {
      case 'pending': return 'warning';
      case 'confirmed': return 'info';
      case 'processing': return 'primary';
      case 'shipped': return 'secondary';
      case 'delivered': return 'success';
      case 'cancelled': return 'error';
      default: return 'default';
    }
  };

  const canCancelOrder = (order) => {
    return ['pending', 'confirmed'].includes(order.status?.toLowerCase());
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
          Order Management
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<RefreshIcon />}
            onClick={() => refetch()}
            disabled={isLoading}
          >
            Refresh
          </Button>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => navigate('/orders/new')}
          >
            Create Order
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

      <Paper sx={{ width: '100%', overflow: 'hidden' }}>
        <TableContainer>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Order ID</TableCell>
                <TableCell>Customer</TableCell>
                <TableCell>Date</TableCell>
                <TableCell>Status</TableCell>
                <TableCell align="right">Total</TableCell>
                <TableCell>Payment</TableCell>
                <TableCell align="center">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 4 }}>
                    <CircularProgress />
                  </TableCell>
                </TableRow>
              ) : orders.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center" sx={{ py: 4 }}>
                    <Typography variant="body1" color="text.secondary">
                      No orders found
                    </Typography>
                  </TableCell>
                </TableRow>
              ) : (
                orders.map((order) => (
                  <TableRow key={order.id} hover>
                    <TableCell>
                      <Typography variant="body2" fontWeight="medium">
                        {order.id}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {order.customerName || order.customerId}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2">
                        {formatDate(order.createdAt)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={order.status}
                        color={getStatusColor(order.status)}
                        size="small"
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell align="right">
                      <Typography variant="body2" fontWeight="medium">
                        {formatCurrency(order.totalAmount)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={order.paymentStatus || 'Pending'}
                        color={order.paymentStatus === 'PAID' ? 'success' : 'warning'}
                        size="small"
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell align="center">
                      <IconButton
                        size="small"
                        onClick={(e) => handleMenuClick(e, order)}
                      >
                        <MoreVertIcon />
                      </IconButton>
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
          count={totalCount}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </Paper>

      {/* Action Menu */}
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        <MenuItem onClick={() => { navigate(`/orders/${selectedOrder?.id}`); handleMenuClose(); }}>
          <ViewIcon sx={{ mr: 1 }} fontSize="small" />
          View Details
        </MenuItem>
        <MenuItem onClick={() => { navigate(`/orders/${selectedOrder?.id}/edit`); handleMenuClose(); }}>
          <EditIcon sx={{ mr: 1 }} fontSize="small" />
          Edit Order
        </MenuItem>
        <MenuItem onClick={() => { setPaymentDialogOpen(true); handleMenuClose(); }}>
          <PaymentIcon sx={{ mr: 1 }} fontSize="small" />
          Update Payment
        </MenuItem>
        <MenuItem onClick={() => { setTrackingDialogOpen(true); handleMenuClose(); }}>
          <TrackingIcon sx={{ mr: 1 }} fontSize="small" />
          Add Tracking
        </MenuItem>
        <MenuItem onClick={() => { setHistoryDialogOpen(true); handleMenuClose(); }}>
          <HistoryIcon sx={{ mr: 1 }} fontSize="small" />
          View History
        </MenuItem>
        {selectedOrder && canCancelOrder(selectedOrder) && (
          <MenuItem 
            onClick={() => { setCancelDialogOpen(true); handleMenuClose(); }}
            sx={{ color: 'error.main' }}
          >
            <CancelIcon sx={{ mr: 1 }} fontSize="small" />
            Cancel Order
          </MenuItem>
        )}
        <MenuItem 
          onClick={() => { setDeleteDialogOpen(true); handleMenuClose(); }}
          sx={{ color: 'error.main' }}
        >
          <DeleteIcon sx={{ mr: 1 }} fontSize="small" />
          Delete Order
        </MenuItem>
      </Menu>

      {/* Cancel Order Dialog */}
      <Dialog open={cancelDialogOpen} onClose={() => setCancelDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Cancel Order</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Are you sure you want to cancel order {selectedOrder?.id}?
          </Typography>
          <TextField
            autoFocus
            margin="dense"
            label="Cancellation Reason"
            fullWidth
            multiline
            rows={3}
            value={cancelReason}
            onChange={(e) => setCancelReason(e.target.value)}
            placeholder="Please provide a reason for cancellation..."
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCancelDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleCancelOrder} 
            color="error" 
            variant="contained"
            disabled={!cancelReason.trim() || cancelOrderMutation.isLoading}
          >
            {cancelOrderMutation.isLoading ? <CircularProgress size={20} /> : 'Cancel Order'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Update Payment Dialog */}
      <Dialog open={paymentDialogOpen} onClose={() => setPaymentDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Update Payment Information</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12} sm={6}>
              <TextField
                select
                fullWidth
                label="Payment Method"
                value={paymentData.paymentMethod}
                onChange={(e) => setPaymentData({ ...paymentData, paymentMethod: e.target.value })}
              >
                <MenuItem value="CASH">Cash</MenuItem>
                <MenuItem value="CARD">Card</MenuItem>
                <MenuItem value="DIGITAL">Digital</MenuItem>
                <MenuItem value="BANK_TRANSFER">Bank Transfer</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                select
                fullWidth
                label="Payment Status"
                value={paymentData.paymentStatus}
                onChange={(e) => setPaymentData({ ...paymentData, paymentStatus: e.target.value })}
              >
                <MenuItem value="PENDING">Pending</MenuItem>
                <MenuItem value="PAID">Paid</MenuItem>
                <MenuItem value="FAILED">Failed</MenuItem>
                <MenuItem value="REFUNDED">Refunded</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Transaction ID"
                value={paymentData.transactionId}
                onChange={(e) => setPaymentData({ ...paymentData, transactionId: e.target.value })}
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
            onClick={handleUpdatePayment} 
            variant="contained"
            disabled={!paymentData.paymentMethod || !paymentData.paymentStatus || updatePaymentMutation.isLoading}
          >
            {updatePaymentMutation.isLoading ? <CircularProgress size={20} /> : 'Update Payment'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Tracking Dialog */}
      <Dialog open={trackingDialogOpen} onClose={() => setTrackingDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Add Tracking Information</DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Tracking Code"
                value={trackingData.trackingCode}
                onChange={(e) => setTrackingData({ ...trackingData, trackingCode: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Carrier"
                value={trackingData.carrier}
                onChange={(e) => setTrackingData({ ...trackingData, carrier: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                type="date"
                label="Estimated Delivery"
                value={trackingData.estimatedDelivery}
                onChange={(e) => setTrackingData({ ...trackingData, estimatedDelivery: e.target.value })}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                multiline
                rows={2}
                label="Notes"
                value={trackingData.notes}
                onChange={(e) => setTrackingData({ ...trackingData, notes: e.target.value })}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setTrackingDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleAddTracking} 
            variant="contained"
            disabled={!trackingData.trackingCode || addTrackingMutation.isLoading}
          >
            {addTrackingMutation.isLoading ? <CircularProgress size={20} /> : 'Add Tracking'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Order History Dialog */}
      <Dialog open={historyDialogOpen} onClose={() => setHistoryDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>Order History - {selectedOrder?.id}</DialogTitle>
        <DialogContent>
          {historyLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
              <CircularProgress />
            </Box>
          ) : orderHistory && orderHistory.length > 0 ? (
            <Box sx={{ mt: 2 }}>
              {orderHistory.map((event, index) => (
                <Card key={index} sx={{ mb: 2 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {event.action || event.type}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {event.description || event.details}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {formatDate(event.timestamp || event.createdAt)}
                    </Typography>
                  </CardContent>
                </Card>
              ))}
            </Box>
          ) : (
            <Typography variant="body2" color="text.secondary" sx={{ py: 4, textAlign: 'center' }}>
              No history available for this order
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setHistoryDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>

      {/* Delete Order Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Delete Order</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Are you sure you want to delete order {selectedOrder?.id}?
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleDeleteOrder} 
            color="error" 
            variant="contained"
            disabled={deleteOrderMutation.isLoading}
          >
            {deleteOrderMutation.isLoading ? <CircularProgress size={20} /> : 'Delete Order'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default OrdersPage;
