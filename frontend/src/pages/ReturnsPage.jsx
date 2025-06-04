import React from 'react';
import { Box, Typography, Paper, Button } from '@mui/material';
import { Refresh as RefreshIcon } from '@mui/icons-material';

const ReturnsPage = () => {
  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">Returns Management</Typography>
        <Button 
          variant="outlined" 
          startIcon={<RefreshIcon />}
        >
          Refresh
        </Button>
      </Box>
      
      <Paper elevation={3} sx={{ p: 3 }}>
        <Typography>Returns management interface will be implemented here</Typography>
      </Paper>
    </Box>
  );
};

export default ReturnsPage;
