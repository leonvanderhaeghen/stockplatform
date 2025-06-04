import React from 'react';
import { Box, Typography, Paper, Button } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';

const InventoryPage = () => {
  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          Inventory Management
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
        >
          Add Inventory
        </Button>
      </Box>
      
      <Paper sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="h6" color="textSecondary" sx={{ mb: 2 }}>
          Inventory Management Coming Soon
        </Typography>
        <Typography variant="body1" color="textSecondary">
          This section is under development. You'll be able to manage your inventory here.
        </Typography>
      </Paper>
    </Box>
  );
};

export default InventoryPage;
