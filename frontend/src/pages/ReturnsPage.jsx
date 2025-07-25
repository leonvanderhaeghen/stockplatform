import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Paper,
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
  TextField,
  Snackbar,
  Alert,
  CircularProgress,
} from '@mui/material';
import RefreshIcon from '@mui/icons-material/Refresh';
import UndoIcon from '@mui/icons-material/Undo';
import orderService from '../services/orderService';

const ReturnsPage = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(false);
  const [dialog, setDialog] = useState({ open: false, id: '', reason: '' });
  const [snack, setSnack] = useState({ open: false, message: '', severity: 'success' });

  const load = async () => {
    try {
      setLoading(true);
      const response = await orderService.getMyOrders();
      // Handle both direct array response and nested data.orders response
      const ordersList = Array.isArray(response) 
        ? response 
        : (response.data?.orders || response.orders || []);
      setOrders(ordersList);
    } catch (e) {
      console.error(e);
      setSnack({ 
        open: true, 
        message: 'Failed to load orders', 
        severity: 'error' 
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openDialog = (id) => setDialog({ open: true, id, reason: '' });
  const closeDialog = () => setDialog({ open: false, id: '', reason: '' });

  const handleReturn = async () => {
    try {
      await orderService.cancelOrder(dialog.id, dialog.reason);
      setSnack({ open: true, message: 'Return requested', severity: 'success' });
      closeDialog();
      load();
    } catch (e) {
      console.error(e);
      setSnack({ open: true, message: 'Failed to request return', severity: 'error' });
    }
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Returns Management</Typography>
        <Button variant="outlined" startIcon={<RefreshIcon />} onClick={load} disabled={loading}>
          Refresh
        </Button>
      </Box>
      <Paper sx={{ p: 2 }}>
        {loading ? (
          <CircularProgress />
        ) : (
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>ID</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Items</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {orders.map((o) => (
                  <TableRow key={o.id}>
                    <TableCell>{o.id}</TableCell>
                    <TableCell>{o.status}</TableCell>
                    <TableCell>{o.items?.length}</TableCell>
                    <TableCell align="right">
                      <IconButton color="primary" onClick={() => openDialog(o.id)} disabled={o.status !== 'PLACED'}>
                        <UndoIcon />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Paper>

      <Dialog open={dialog.open} onClose={closeDialog}>
        <DialogTitle>Request Return</DialogTitle>
        <DialogContent>
          <TextField
            fullWidth
            label="Reason"
            multiline
            rows={3}
            value={dialog.reason}
            onChange={(e) => setDialog({ ...dialog, reason: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={closeDialog}>Cancel</Button>
          <Button variant="contained" onClick={handleReturn} disabled={!dialog.reason}>
            Submit
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snack.open}
        onClose={() => setSnack({ ...snack, open: false })}
        autoHideDuration={4000}
      >
        <Alert severity={snack.severity}>{snack.message}</Alert>
      </Snackbar>
    </Box>
  );
};

export default ReturnsPage;
