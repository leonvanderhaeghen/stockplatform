import React, { useState } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Tabs, 
  Tab, 
  Grid, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  Button
} from '@mui/material';
import { DateRange as DateRangeIcon, Download as DownloadIcon } from '@mui/icons-material';

const reportTypes = [
  'Sales Report',
  'Inventory Report',
  'Customer Report',
  'Revenue Report',
  'Product Performance'
];

const timeRanges = [
  'Today',
  'This Week',
  'This Month',
  'This Quarter',
  'This Year',
  'Custom Range'
];

const ReportsPage = () => {
  const [activeTab, setActiveTab] = useState(0);
  const [reportType, setReportType] = useState('');
  const [timeRange, setTimeRange] = useState('This Month');

  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Reports
      </Typography>
      
      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Tabs 
          value={activeTab} 
          onChange={handleTabChange}
          indicatorColor="primary"
          textColor="primary"
          variant="scrollable"
          scrollButtons="auto"
        >
          <Tab label="Sales" />
          <Tab label="Inventory" />
          <Tab label="Customers" />
          <Tab label="Revenue" />
          <Tab label="Products" />
        </Tabs>
        
        <Box mt={3}>
          <Grid container spacing={3}>
            <Grid item xs={12} md={5}>
              <FormControl fullWidth size="small">
                <InputLabel>Report Type</InputLabel>
                <Select
                  value={reportType}
                  label="Report Type"
                  onChange={(e) => setReportType(e.target.value)}
                >
                  {reportTypes.map((type) => (
                    <MenuItem key={type} value={type}>{type}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={5}>
              <FormControl fullWidth size="small">
                <InputLabel>Time Range</InputLabel>
                <Select
                  value={timeRange}
                  label="Time Range"
                  onChange={(e) => setTimeRange(e.target.value)}
                  startAdornment={<DateRangeIcon />}
                >
                  {timeRanges.map((range) => (
                    <MenuItem key={range} value={range}>{range}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={2} display="flex" alignItems="flex-end">
              <Button 
                fullWidth 
                variant="contained" 
                color="primary"
                startIcon={<DownloadIcon />}
                disabled={!reportType}
              >
                Generate
              </Button>
            </Grid>
          </Grid>
        </Box>
      </Paper>
      
      <Paper elevation={3} sx={{ p: 3, minHeight: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Typography>Report visualization will be displayed here</Typography>
      </Paper>
    </Box>
  );
};

export default ReportsPage;
