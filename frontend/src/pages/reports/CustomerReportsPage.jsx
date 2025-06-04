import React from 'react';
import { Box, Typography, Paper, Button, Grid, MenuItem, TextField } from '@mui/material';
import { People as PeopleIcon, Download as DownloadIcon } from '@mui/icons-material';

const CustomerReportsPage = () => {
  const [reportType, setReportType] = React.useState('customer_list');
  const [customerGroup, setCustomerGroup] = React.useState('all');

  const reportTypes = [
    { value: 'customer_list', label: 'Customer List' },
    { value: 'customer_segments', label: 'Customer Segments' },
    { value: 'purchase_history', label: 'Purchase History' },
    { value: 'loyalty', label: 'Loyalty Program' },
    { value: 'customer_activity', label: 'Customer Activity' },
  ];

  const customerGroups = [
    { value: 'all', label: 'All Customers' },
    { value: 'repeat', label: 'Repeat Customers' },
    { value: 'new', label: 'New Customers' },
    { value: 'inactive', label: 'Inactive Customers' },
    { value: 'vip', label: 'VIP Customers' },
  ];

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Customer Reports
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
              label="Customer Group"
              value={customerGroup}
              onChange={(e) => setCustomerGroup(e.target.value)}
              size="small"
            >
              {customerGroups.map((option) => (
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
          <PeopleIcon color="action" sx={{ fontSize: 64, mb: 2, opacity: 0.5 }} />
          <Typography variant="h6" color="textSecondary" gutterBottom>
            {reportTypes.find(r => r.value === reportType)?.label}
          </Typography>
          <Typography variant="body1" color="textSecondary">
            {customerGroups.find(g => g.value === customerGroup)?.label} • {reportType === 'customer_list' ? 'List View' : 'Summary View'}
          </Typography>
        </Box>
      </Paper>
    </Box>
  );
};

export default CustomerReportsPage;
