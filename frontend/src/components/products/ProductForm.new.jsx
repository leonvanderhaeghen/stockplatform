import React from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import {
  Box,
  Button,
  TextField,
  Grid,
  Paper,
  Typography,
  Divider,
  FormControlLabel,
  Switch,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  FormHelperText,
} from '@mui/material';

const validationSchema = Yup.object({
  name: Yup.string().required('Name is required'),
  sku: Yup.string().required('SKU is required'),
  description: Yup.string(),
  price: Yup.number().min(0, 'Price must be positive').required('Price is required'),
  cost: Yup.number().min(0, 'Cost must be positive').required('Cost is required'),
  category: Yup.string().required('Category is required'),
  inStock: Yup.boolean().default(true),
  stock: Yup.number().min(0, 'Stock cannot be negative').required('Stock is required'),
});

const ProductForm = ({
  initialValues = {},
  onSubmit,
  loading = false,
  isEdit = false,
  categories = [],
}) => {
  const formik = useFormik({
    initialValues: {
      name: '',
      sku: '',
      description: '',
      price: 0,
      cost: 0,
      category: '',
      inStock: true,
      stock: 0,
      ...initialValues,
    },
    validationSchema,
    onSubmit: (values) => {
      onSubmit(values);
    },
    enableReinitialize: true,
  });

  return (
    <form onSubmit={formik.handleSubmit}>
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          {isEdit ? 'Edit Product' : 'Create Product'}
        </Typography>
        <Divider sx={{ mb: 3 }} />
        
        <Grid container spacing={3}>
          <Grid item xs={12} md={8}>
            <TextField
              fullWidth
              id="name"
              name="name"
              label="Product Name"
              value={formik.values.name}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.name && Boolean(formik.errors.name)}
              helperText={formik.touched.name && formik.errors.name}
              margin="normal"
              disabled={loading}
            />

            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  id="sku"
                  name="sku"
                  label="SKU"
                  value={formik.values.sku}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  error={formik.touched.sku && Boolean(formik.errors.sku)}
                  helperText={formik.touched.sku && formik.errors.sku}
                  margin="normal"
                  disabled={loading}
                />
              </Grid>
              <Grid item xs={12} sm={6}>
                <FormControl fullWidth margin="normal" error={formik.touched.category && Boolean(formik.errors.category)}>
                  <InputLabel id="category-label">Category</InputLabel>
                  <Select
                    labelId="category-label"
                    id="category"
                    name="category"
                    value={formik.values.category}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    label="Category"
                    disabled={loading}
                  >
                    {categories.map((category) => (
                      <MenuItem key={category.id} value={category.id}>
                        {category.name}
                      </MenuItem>
                    ))}
                  </Select>
                  {formik.touched.category && formik.errors.category && (
                    <FormHelperText>{formik.errors.category}</FormHelperText>
                  )}
                </FormControl>
              </Grid>
            </Grid>

            <TextField
              fullWidth
              id="description"
              name="description"
              label="Description"
              multiline
              rows={4}
              value={formik.values.description}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.description && Boolean(formik.errors.description)}
              helperText={formik.touched.description && formik.errors.description}
              margin="normal"
              disabled={loading}
            />

            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  id="price"
                  name="price"
                  label="Price"
                  type="number"
                  value={formik.values.price}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  error={formik.touched.price && Boolean(formik.errors.price)}
                  helperText={formik.touched.price && formik.errors.price}
                  margin="normal"
                  disabled={loading}
                  InputProps={{
                    startAdornment: <span style={{ marginRight: 8 }}>$</span>,
                  }}
                />
              </Grid>
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  id="cost"
                  name="cost"
                  label="Cost"
                  type="number"
                  value={formik.values.cost}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  error={formik.touched.cost && Boolean(formik.errors.cost)}
                  helperText={formik.touched.cost && formik.errors.cost}
                  margin="normal"
                  disabled={loading}
                  InputProps={{
                    startAdornment: <span style={{ marginRight: 8 }}>$</span>,
                  }}
                />
              </Grid>
            </Grid>

            <Grid container spacing={2} marginTop={2}>
              <Grid item xs={12} sm={6}>
                <TextField
                  fullWidth
                  id="stock"
                  name="stock"
                  label="Stock Quantity"
                  type="number"
                  value={formik.values.stock}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  error={formik.touched.stock && Boolean(formik.errors.stock)}
                  helperText={formik.touched.stock && formik.errors.stock}
                  disabled={loading}
                />
              </Grid>
              <Grid item xs={12} sm={6}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={formik.values.inStock}
                      onChange={(e) => formik.setFieldValue('inStock', e.target.checked)}
                      name="inStock"
                      color="primary"
                      disabled={loading}
                    />
                  }
                  label={formik.values.inStock ? 'In Stock' : 'Out of Stock'}
                  sx={{ mt: 2 }}
                />
              </Grid>
            </Grid>
          </Grid>
        </Grid>

        <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 2, mt: 4 }}>
          <Button
            variant="outlined"
            onClick={() => window.history.back()}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            disabled={loading}
          >
            {loading ? 'Saving...' : isEdit ? 'Update Product' : 'Create Product'}
          </Button>
        </Box>
      </Paper>
    </form>
  );
};

export default ProductForm;
