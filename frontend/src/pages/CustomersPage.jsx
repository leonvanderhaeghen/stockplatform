import React from 'react';
import { Box, Typography, Paper, Button, TextField, InputAdornment } from '@mui/material';
import { Search as SearchIcon, Add as AddIcon } from '@mui/icons-material';

const CustomersPage = () => {
  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">Customers</Typography>
        <Box display="flex" gap={2}>
          <TextField
            size="small"
            placeholder="Search customers..."
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
          />
          <Button 
            variant="contained" 
            color="primary"
            startIcon={<AddIcon />}
          >
            Add Customer
          </Button>
        </Box>
      </Box>
      
      <Paper elevation={3} sx={{ p: 3 }}>
        <Typography>Customer list will be displayed here</Typography>
      </Paper>
    </Box>
  );
};

export default CustomersPage;
