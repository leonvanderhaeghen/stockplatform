import React from 'react';
import { Box, Typography, Paper, Button, Grid, MenuItem, TextField } from '@mui/material';
import { Assessment as AssessmentIcon, Download as DownloadIcon } from '@mui/icons-material';

const InventoryReportsPage = () => {
  const [reportType, setReportType] = React.useState('stock_levels');
  const [location, setLocation] = React.useState('all');

  const reportTypes = [
    { value: 'stock_levels', label: 'Stock Levels' },
    { value: 'low_stock', label: 'Low Stock' },
    { value: 'out_of_stock', label: 'Out of Stock' },
    { value: 'stock_movements', label: 'Stock Movements' },
    { value: 'inventory_valuation', label: 'Inventory Valuation' },
  ];

  const locations = [
    { value: 'all', label: 'All Locations' },
    { value: 'warehouse', label: 'Main Warehouse' },
    { value: 'store', label: 'Retail Store' },
    { value: 'online', label: 'Online Store' },
  ];

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Inventory Reports
      </Typography>
      
      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={5}>
            <TextField
              select
              fullWidth
              label="Report Type"
              value={reportType}
              onChange={(e) => setReportType(e.target.value)}
              size="small"
            >
              {reportTypes.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </TextField>
          </Grid>
          <Grid item xs={12} md={5}>
            <TextField
              select
              fullWidth
              label="Location"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              size="small"
            >
              {locations.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </TextField>
          </Grid>
          <Grid item xs={12} md={2} display="flex" alignItems="flex-end">
            <Button
              fullWidth
              variant="contained"
              color="primary"
              startIcon={<DownloadIcon />}
            >
              Generate
            </Button>
          </Grid>
        </Grid>
      </Paper>
      
      <Paper elevation={3} sx={{ p: 3, minHeight: '400px' }}>
        <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height="100%" textAlign="center">
          <AssessmentIcon color="action" sx={{ fontSize: 64, mb: 2, opacity: 0.5 }} />
          <Typography variant="h6" color="textSecondary" gutterBottom>
            {reportTypes.find(r => r.value === reportType)?.label} Report
          </Typography>
          <Typography variant="body1" color="textSecondary">
            Select a report type and location to generate the report
          </Typography>
        </Box>
      </Paper>
    </Box>
  );
};

export default InventoryReportsPage;
