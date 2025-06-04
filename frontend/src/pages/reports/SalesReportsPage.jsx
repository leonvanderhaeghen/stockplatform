import React from 'react';
import { Box, Typography, Paper, Button, Grid, MenuItem, TextField } from '@mui/material';
import { DateRange as DateRangeIcon, Download as DownloadIcon } from '@mui/icons-material';

const SalesReportsPage = () => {
  const [reportType, setReportType] = React.useState('daily');
  const [dateRange, setDateRange] = React.useState('this_week');

  const reportTypes = [
    { value: 'daily', label: 'Daily Sales' },
    { value: 'weekly', label: 'Weekly Sales' },
    { value: 'monthly', label: 'Monthly Sales' },
    { value: 'yearly', label: 'Yearly Sales' },
    { value: 'custom', label: 'Custom Range' },
  ];

  const dateRanges = [
    { value: 'today', label: 'Today' },
    { value: 'yesterday', label: 'Yesterday' },
    { value: 'this_week', label: 'This Week' },
    { value: 'last_week', label: 'Last Week' },
    { value: 'this_month', label: 'This Month' },
    { value: 'last_month', label: 'Last Month' },
    { value: 'this_year', label: 'This Year' },
    { value: 'last_year', label: 'Last Year' },
  ];

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Sales Reports
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
              label="Date Range"
              value={dateRange}
              onChange={(e) => setDateRange(e.target.value)}
              size="small"
              InputProps={{
                startAdornment: <DateRangeIcon color="action" sx={{ mr: 1 }} />,
              }}
            >
              {dateRanges.map((option) => (
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
        <Box display="flex" alignItems="center" justifyContent="center" height="100%">
          <Typography>Sales report visualization will be displayed here</Typography>
        </Box>
      </Paper>
    </Box>
  );
};

export default SalesReportsPage;
