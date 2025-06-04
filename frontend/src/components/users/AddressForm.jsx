import React, { useState } from 'react';
import {
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Grid,
  Typography,
  CircularProgress,
  InputAdornment,
  Alert,
  FormControlLabel,
  Checkbox,
  MenuItem,
  Box,
  Divider
} from '@mui/material';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { formatPhoneNumber, parsePhoneNumber } from '../../utils/formatPhoneNumber';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import HomeIcon from '@mui/icons-material/Home';
import WorkIcon from '@mui/icons-material/Work';
import LocalShippingIcon from '@mui/icons-material/LocalShipping';
import PhoneIcon from '@mui/icons-material/Phone';

// US States for the dropdown
const US_STATES = [
  'Alabama', 'Alaska', 'Arizona', 'Arkansas', 'California', 'Colorado',
  'Connecticut', 'Delaware', 'Florida', 'Georgia', 'Hawaii', 'Idaho',
  'Illinois', 'Indiana', 'Iowa', 'Kansas', 'Kentucky', 'Louisiana',
  'Maine', 'Maryland', 'Massachusetts', 'Michigan', 'Minnesota',
  'Mississippi', 'Missouri', 'Montana', 'Nebraska', 'Nevada',
  'New Hampshire', 'New Jersey', 'New Mexico', 'New York',
  'North Carolina', 'North Dakota', 'Ohio', 'Oklahoma', 'Oregon',
  'Pennsylvania', 'Rhode Island', 'South Carolina', 'South Dakota',
  'Tennessee', 'Texas', 'Utah', 'Vermont', 'Virginia', 'Washington',
  'West Virginia', 'Wisconsin', 'Wyoming'
];

// Address types for the name field
const ADDRESS_TYPES = [
  { value: 'Home', label: 'Home', icon: <HomeIcon /> },
  { value: 'Work', label: 'Work', icon: <WorkIcon /> },
  { value: 'Other', label: 'Other', icon: <LocalShippingIcon /> },
];

// Validation schema
const addressSchema = Yup.object().shape({
  name: Yup.string()
    .required('Address name is required')
    .max(50, 'Name is too long'),
  street: Yup.string()
    .required('Street address is required')
    .max(100, 'Address is too long'),
  street2: Yup.string()
    .max(100, 'Address line 2 is too long'),
  city: Yup.string()
    .required('City is required')
    .max(50, 'City name is too long'),
  state: Yup.string()
    .required('State is required'),
  postalCode: Yup.string()
    .required('ZIP/Postal code is required')
    .matches(/^[0-9]{5}(-[0-9]{4})?$/, 'Invalid ZIP/Postal code'),
  country: Yup.string()
    .required('Country is required'),
  phone: Yup.string()
    .test('phone', 'Invalid phone number', (value) => {
      if (!value) return true; // Optional field
      const phoneNumber = parsePhoneNumber(value);
      return phoneNumber.isValid();
    }),
  isDefault: Yup.boolean()
});

const AddressForm = ({ address, onSubmit, onCancel, loading }) => {
  const [error, setError] = useState('');
  
  const formik = useFormik({
    initialValues: {
      name: address?.name || 'Home',
      street: address?.street || '',
      street2: address?.street2 || '',
      city: address?.city || '',
      state: address?.state || '',
      postalCode: address?.postalCode || '',
      country: address?.country || 'United States',
      phone: address?.phone ? formatPhoneNumber(address.phone, false) : '',
      isDefault: address?.isDefault || false,
    },
    validationSchema: addressSchema,
    enableReinitialize: true,
    onSubmit: async (values) => {
      try {
        setError('');
        // Format phone number before submitting
        const formattedValues = {
          ...values,
          phone: values.phone ? parsePhoneNumber(values.phone).format('E.164') : ''
        };
        await onSubmit(formattedValues);
      } catch (err) {
        setError(err.response?.data?.message || 'Failed to save address');
      }
    },
  });

  // Format phone number as user types
  const handlePhoneChange = (e) => {
    const formatted = formatPhoneNumber(e.target.value, true);
    formik.setFieldValue('phone', formatted);
  };

  // Get icon for address type
  const getAddressTypeIcon = (type) => {
    const addressType = ADDRESS_TYPES.find(t => t.value === type);
    return addressType ? addressType.icon : <LocationOnIcon />;
  };

  return (
    <form onSubmit={formik.handleSubmit}>
      <DialogContent dividers>
        {error && (
          <Alert severity="error" sx={{ mb: 3 }}>
            {error}
          </Alert>
        )}
        
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <TextField
              select
              fullWidth
              id="name"
              name="name"
              label="Address Type"
              value={formik.values.name}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.name && Boolean(formik.errors.name)}
              helperText={formik.touched.name && formik.errors.name}
              disabled={loading}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    {getAddressTypeIcon(formik.values.name)}
                  </InputAdornment>
                ),
              }}
              SelectProps={{
                renderValue: (selected) => {
                  const option = ADDRESS_TYPES.find(opt => opt.value === selected);
                  return option ? option.label : selected;
                }
              }}
            >
              {ADDRESS_TYPES.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  <Box display="flex" alignItems="center">
                    {option.icon}
                    <Box ml={1}>{option.label}</Box>
                  </Box>
                </MenuItem>
              ))}
            </TextField>
          </Grid>
          
          <Grid item xs={12}>
            <TextField
              fullWidth
              id="street"
              name="street"
              label="Street Address"
              value={formik.values.street}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.street && Boolean(formik.errors.street)}
              helperText={formik.touched.street && formik.errors.street}
              disabled={loading}
              placeholder="123 Main St"
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <LocationOnIcon color="action" />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          
          <Grid item xs={12}>
            <TextField
              fullWidth
              id="street2"
              name="street2"
              label="Apartment, suite, etc. (Optional)"
              value={formik.values.street2}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.street2 && Boolean(formik.errors.street2)}
              helperText={formik.touched.street2 && formik.errors.street2}
              disabled={loading}
              placeholder="Apt 4B"
            />
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              id="city"
              name="city"
              label="City"
              value={formik.values.city}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.city && Boolean(formik.errors.city)}
              helperText={formik.touched.city && formik.errors.city}
              disabled={loading}
            />
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <TextField
              select
              fullWidth
              id="state"
              name="state"
              label="State/Province/Region"
              value={formik.values.state}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.state && Boolean(formik.errors.state)}
              helperText={formik.touched.state && formik.errors.state}
              disabled={loading}
              SelectProps={{
                native: true,
              }}
            >
              <option value="">Select a state</option>
              {US_STATES.map((state) => (
                <option key={state} value={state}>
                  {state}
                </option>
              ))}
            </TextField>
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              id="postalCode"
              name="postalCode"
              label="ZIP/Postal Code"
              value={formik.values.postalCode}
              onChange={(e) => {
                // Allow only numbers and format as user types
                const value = e.target.value.replace(/\D/g, '');
                let formatted = value;
                
                // Add hyphen for ZIP+4 format (12345-6789)
                if (value.length > 5) {
                  formatted = `${value.slice(0, 5)}-${value.slice(5, 9)}`;
                }
                
                formik.setFieldValue('postalCode', formatted);
              }}
              onBlur={formik.handleBlur}
              error={formik.touched.postalCode && Boolean(formik.errors.postalCode)}
              helperText={formik.touched.postalCode && formik.errors.postalCode}
              disabled={loading}
              inputProps={{
                maxLength: 10, // 12345-6789
              }}
            />
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              id="country"
              name="country"
              label="Country/Region"
              value={formik.values.country}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.country && Boolean(formik.errors.country)}
              helperText={formik.touched.country && formik.errors.country}
              disabled={loading || true} // Currently only US is supported
              SelectProps={{
                native: true,
              }}
            >
              <option value="United States">United States</option>
            </TextField>
          </Grid>
          
          <Grid item xs={12}>
            <Divider sx={{ my: 2 }} />
            <Typography variant="subtitle2" gutterBottom>
              Contact Information (Optional)
            </Typography>
          </Grid>
          
          <Grid item xs={12}>
            <TextField
              fullWidth
              id="phone"
              name="phone"
              label="Phone Number"
              value={formik.values.phone}
              onChange={handlePhoneChange}
              onBlur={formik.handleBlur}
              error={formik.touched.phone && Boolean(formik.errors.phone)}
              helperText={formik.touched.phone ? formik.errors.phone || 'For delivery questions only' : 'For delivery questions only'}
              disabled={loading}
              placeholder="(123) 456-7890"
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <PhoneIcon color="action" />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          
          <Grid item xs={12}>
            <FormControlLabel
              control={
                <Checkbox
                  checked={formik.values.isDefault}
                  onChange={(e) => formik.setFieldValue('isDefault', e.target.checked)}
                  name="isDefault"
                  color="primary"
                  disabled={loading}
                />
              }
              label="Set as default shipping address"
            />
          </Grid>
        </Grid>
      </DialogContent>
      
      <DialogActions sx={{ p: 2, borderTop: 1, borderColor: 'divider' }}>
        <Button 
          onClick={onCancel} 
          disabled={loading}
          variant="outlined"
          sx={{ minWidth: 100 }}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          variant="contained"
          color="primary"
          disabled={loading || !formik.isValid || !formik.dirty}
          sx={{ minWidth: 100 }}
          startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}
        >
          {loading ? 'Saving...' : 'Save Address'}
        </Button>
      </DialogActions>
    </form>
  );
};

export default AddressForm;
