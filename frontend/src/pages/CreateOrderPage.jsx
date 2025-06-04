import React from 'react';
import { Box, Typography, Paper, Stepper, Step, StepLabel } from '@mui/material';

const steps = ['Customer Details', 'Order Items', 'Review & Place Order'];

const CreateOrderPage = () => {
  const [activeStep] = React.useState(0);

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Create New Order
      </Typography>
      
      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <Stepper activeStep={activeStep} alternativeLabel>
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
            </Step>
          ))}
        </Stepper>
      </Paper>
      
      <Paper elevation={3} sx={{ p: 3 }}>
        <Box sx={{ minHeight: '300px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <Typography>Order creation form will be implemented here</Typography>
        </Box>
      </Paper>
    </Box>
  );
};

export default CreateOrderPage;
