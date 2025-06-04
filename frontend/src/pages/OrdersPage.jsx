import React from 'react';
import { Box, Typography, Paper, Button } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';

const OrdersPage = () => {
  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          Order Management
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
        >
          Create Order
        </Button>
      </Box>
      
      <Paper sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="h6" color="textSecondary" sx={{ mb: 2 }}>
          Order Management Coming Soon
        </Typography>
        <Typography variant="body1" color="textSecondary">
          This section is under development. You'll be able to manage your orders here.
        </Typography>
      </Paper>
    </Box>
  );
};

export default OrdersPage;
