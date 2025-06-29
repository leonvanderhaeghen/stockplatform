import React from 'react';
import { Alert, Box, Button, Typography } from '@mui/material';
import { Refresh as RefreshIcon } from '@mui/icons-material';

class InventoryErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true, error };
  }

  componentDidCatch(error, errorInfo) {
    console.error('Inventory Error Boundary caught an error:', error, errorInfo);
  }

  handleRetry = () => {
    this.setState({ hasError: false, error: null });
    if (this.props.onRetry) {
      this.props.onRetry();
    }
  };

  render() {
    if (this.state.hasError) {
      return (
        <Box sx={{ p: 3 }}>
          <Alert severity="error" sx={{ mb: 2 }}>
            <Typography variant="h6" gutterBottom>
              Inventory Data Error
            </Typography>
            <Typography variant="body2" sx={{ mb: 2 }}>
              {this.state.error?.message || 'An error occurred while loading inventory data.'}
            </Typography>
            <Button
              variant="outlined"
              startIcon={<RefreshIcon />}
              onClick={this.handleRetry}
              size="small"
            >
              Retry
            </Button>
          </Alert>
        </Box>
      );
    }

    return this.props.children;
  }
}

export default InventoryErrorBoundary;
