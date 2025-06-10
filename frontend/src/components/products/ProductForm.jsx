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
  FormControlLabel,
  Switch,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  FormHelperText,
  InputAdornment,
  IconButton,
  Chip,
  Tooltip,
  Collapse,
  OutlinedInput,
  Checkbox,
  ListItemText,
} from '@mui/material';
import { 
  Delete as DeleteIcon, 
  Add as AddIcon, 
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  Image as ImageIcon,
  Videocam as VideoIcon,
  Link as LinkIcon,
} from '@mui/icons-material';

// Currency options for the currency selector
const CURRENCY_OPTIONS = [
  { value: 'USD', label: 'USD - US Dollar' },
  { value: 'EUR', label: 'EUR - Euro' },
  { value: 'GBP', label: 'GBP - British Pound' },
  { value: 'JPY', label: 'JPY - Japanese Yen' },
  { value: 'CNY', label: 'CNY - Chinese Yuan' },
];

// Helper component for array inputs (images, videos, etc.)
const ArrayInput = ({
  values = [],
  onAdd,
  onRemove,
  label,
  placeholder,
  helperText,
  error,
  touched,
  icon: Icon,
}) => {
  const [inputValue, setInputValue] = React.useState('');

  const handleAdd = () => {
    if (inputValue.trim()) {
      onAdd([...values, inputValue.trim()]);
      setInputValue('');
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAdd();
    }
  };

  return (
    <FormControl fullWidth margin="normal" error={error && touched}>
      <InputLabel shrink>{label}</InputLabel>
      <Box display="flex" alignItems="center" gap={1}>
        <TextField
          fullWidth
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder={placeholder}
          variant="outlined"
          size="small"
          InputProps={{
            startAdornment: Icon && (
              <InputAdornment position="start">
                <Icon color="action" />
              </InputAdornment>
            ),
            endAdornment: (
              <InputAdornment position="end">
                <Tooltip title={`Add ${label.toLowerCase()}`}>
                  <IconButton onClick={handleAdd} edge="end" size="small">
                    <AddIcon />
                  </IconButton>
                </Tooltip>
              </InputAdornment>
            ),
          }}
        />
      </Box>
      <Box mt={1}>
        {values.map((item, index) => (
          <Chip
            key={index}
            label={item}
            onDelete={() => onRemove(index)}
            style={{ marginRight: 8, marginBottom: 8 }}
            deleteIcon={<DeleteIcon />}
            variant="outlined"
          />
        ))}
      </Box>
      {error && touched && <FormHelperText>{error}</FormHelperText>}
      {helperText && !error && <FormHelperText>{helperText}</FormHelperText>}
    </FormControl>
  );
};

// Helper component for metadata key-value pairs
const MetadataInput = ({ values = {}, onChange, error, touched }) => {
  const [metadata, setMetadata] = React.useState(
    Object.entries(values).map(([key, value]) => ({ key, value }))
  );

  const handleAdd = () => {
    setMetadata([...metadata, { key: '', value: '' }]);
  };

  const handleRemove = (index) => {
    const newMetadata = [...metadata];
    newMetadata.splice(index, 1);
    setMetadata(newMetadata);
    updateParent(newMetadata);
  };

  const handleChange = (index, field, value) => {
    const newMetadata = [...metadata];
    newMetadata[index] = { ...newMetadata[index], [field]: value };
    setMetadata(newMetadata);
    updateParent(newMetadata);
  };

  const updateParent = (items) => {
    const metadataObj = {};
    items.forEach(({ key, value }) => {
      if (key.trim()) {
        metadataObj[key] = value;
      }
    });
    onChange(metadataObj);
  };

  return (
    <FormControl fullWidth margin="normal" error={error && touched}>
      <InputLabel shrink>Metadata (Key-Value Pairs)</InputLabel>
      <Box>
        {metadata.map((item, index) => (
          <Box key={index} display="flex" gap={1} mb={1}>
            <TextField
              label="Key"
              value={item.key}
              onChange={(e) => handleChange(index, 'key', e.target.value)}
              size="small"
              style={{ flex: 1 }}
            />
            <TextField
              label="Value"
              value={item.value}
              onChange={(e) => handleChange(index, 'value', e.target.value)}
              size="small"
              style={{ flex: 2 }}
            />
            <IconButton 
              onClick={() => handleRemove(index)}
              color="error"
              size="small"
              style={{ alignSelf: 'center' }}
            >
              <DeleteIcon />
            </IconButton>
          </Box>
        ))}
        <Button
          onClick={handleAdd}
          startIcon={<AddIcon />}
          size="small"
          color="primary"
        >
          Add Metadata
        </Button>
      </Box>
      {error && touched && <FormHelperText>{error}</FormHelperText>}
    </FormControl>
  );
};

// Validation schema for the product form
const validationSchema = Yup.object().shape({
  name: Yup.string().required('Product name is required'),
  sku: Yup.string().required('SKU is required'),
  barcode: Yup.string(),
  description: Yup.string(),
  costPrice: Yup.string()
    .required('Cost price is required')
    .test('is-decimal', 'Invalid price format', (value) => {
      if (!value) return false;
      return /^\d+(\.\d{1,2})?$/.test(value);
    }),
  sellingPrice: Yup.string()
    .required('Selling price is required')
    .test('is-decimal', 'Invalid price format', (value) => {
      if (!value) return false;
      return /^\d+(\.\d{1,2})?$/.test(value);
    }),
  currency: Yup.string().required('Currency is required'),
  inStock: Yup.boolean().default(true),
  stockQty: Yup.number()
    .min(0, 'Stock quantity cannot be negative')
    .required('Stock quantity is required'),
  lowStockAt: Yup.number()
    .min(0, 'Low stock threshold cannot be negative')
    .test('less-than-stock', 'Low stock threshold must be less than stock quantity', function(value) {
      if (value === undefined || value === '') return true;
      const stockQty = this.parent.stockQty;
      return value < stockQty;
    }),
  isActive: Yup.boolean().default(true),
  categoryIds: Yup.array().of(Yup.string()),
  supplierId: Yup.string().required('Supplier is required'),
  imageUrls: Yup.array().of(Yup.string().url('Must be a valid URL')),
  videoUrls: Yup.array().of(Yup.string().url('Must be a valid URL')),
  metadata: Yup.object().default({}),
});

const ProductForm = ({
  initialValues = {
    name: '',
    sku: '',
    barcode: '',
    description: '',
    costPrice: '',
    sellingPrice: '',
    currency: 'USD',
    inStock: true,
    stockQty: 0,
    lowStockAt: 0,
    isActive: true,
    categoryIds: [],
    supplierId: '',
    imageUrls: [],
    videoUrls: [],
    metadata: {},
  },
  onSubmit,
  loading = false,
  isEdit = false,
  categories = [],
  suppliers = [],
}) => {
  const [expandedSections, setExpandedSections] = React.useState({
    basic: true,
    pricing: true,
    inventory: true,
    media: false,
    metadata: false,
  });

  const toggleSection = (section) => {
    setExpandedSections(prev => ({
      ...prev,
      [section]: !prev[section]
    }));
  };

  const formik = useFormik({
    initialValues: {
      name: '',
      sku: '',
      barcode: '',
      description: '',
      costPrice: '0.00',
      sellingPrice: '0.00',
      currency: 'USD',
      categoryIds: [],
      supplierId: '',
      inStock: true,
      stockQty: 0,
      lowStockAt: 0,
      isActive: true,
      imageUrls: [],
      videoUrls: [],
      metadata: {},
      ...initialValues,
    },
    validationSchema,
    onSubmit: async (values, { setSubmitting }) => {
      try {
        const formattedValues = {
          ...values,
          costPrice: parseFloat(values.costPrice).toFixed(2),
          sellingPrice: parseFloat(values.sellingPrice).toFixed(2),
        };
        await onSubmit(formattedValues);
      } catch (error) {
        console.error('Error submitting form:', error);
      } finally {
        setSubmitting(false);
      }
    },
  });

  const renderSectionHeader = (title, section) => (
    <Box 
      display="flex" 
      alignItems="center" 
      justifyContent="space-between" 
      sx={{ 
        cursor: 'pointer', 
        p: 1, 
        '&:hover': { bgcolor: 'action.hover' } 
      }}
      onClick={() => toggleSection(section)}
    >
      <Typography variant="subtitle1" fontWeight="medium">
        {title}
      </Typography>
      {expandedSections[section] ? <ExpandLessIcon /> : <ExpandMoreIcon />}
    </Box>
  );

  return (
    <form onSubmit={formik.handleSubmit}>
      <Paper sx={{ p: 3 }}>
        <Grid container spacing={3}>
          {/* Basic Information Section */}
          <Grid item xs={12}>
            {renderSectionHeader('Basic Information', 'basic')}
            <Collapse in={expandedSections.basic}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
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
                </Grid>
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
                <Grid item xs={12}>
                  <TextField
                    fullWidth
                    id="description"
                    name="description"
                    label="Description"
                    multiline
                    rows={3}
                    value={formik.values.description}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    error={formik.touched.description && Boolean(formik.errors.description)}
                    helperText={formik.touched.description && formik.errors.description}
                    margin="normal"
                    disabled={loading}
                  />
                </Grid>
              </Grid>
            </Collapse>
          </Grid>

          {/* Pricing Section */}
          <Grid item xs={12}>
            {renderSectionHeader('Pricing', 'pricing')}
            <Collapse in={expandedSections.pricing}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    id="costPrice"
                    name="costPrice"
                    label="Cost Price"
                    value={formik.values.costPrice}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    error={formik.touched.costPrice && Boolean(formik.errors.costPrice)}
                    helperText={formik.touched.costPrice && formik.errors.costPrice}
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
                    id="sellingPrice"
                    name="sellingPrice"
                    label="Selling Price"
                    value={formik.values.sellingPrice}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    error={formik.touched.sellingPrice && Boolean(formik.errors.sellingPrice)}
                    helperText={formik.touched.sellingPrice && formik.errors.sellingPrice}
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
                    id="currency"
                    name="currency"
                    label="Currency"
                    value={formik.values.currency}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    error={formik.touched.currency && Boolean(formik.errors.currency)}
                    helperText={formik.touched.currency && formik.errors.currency}
                    margin="normal"
                    disabled={loading}
                    select
                  >
                    {CURRENCY_OPTIONS.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {option.label}
                      </MenuItem>
                    ))}
                  </TextField>
                </Grid>
              </Grid>
            </Collapse>
          </Grid>

          {/* Inventory Section */}
          <Grid item xs={12}>
            {renderSectionHeader('Inventory', 'inventory')}
            <Collapse in={expandedSections.inventory}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    id="stockQty"
                    name="stockQty"
                    label="Stock Quantity"
                    type="number"
                    value={formik.values.stockQty}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    error={formik.touched.stockQty && Boolean(formik.errors.stockQty)}
                    helperText={formik.touched.stockQty && formik.errors.stockQty}
                    margin="normal"
                    disabled={loading}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    fullWidth
                    id="lowStockAt"
                    name="lowStockAt"
                    label="Low Stock Threshold"
                    type="number"
                    value={formik.values.lowStockAt}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    error={formik.touched.lowStockAt && Boolean(formik.errors.lowStockAt)}
                    helperText={formik.touched.lowStockAt && formik.errors.lowStockAt}
                    margin="normal"
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
                <Grid item xs={12} sm={6}>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={formik.values.isActive}
                        onChange={(e) => formik.setFieldValue('isActive', e.target.checked)}
                        name="isActive"
                        color="primary"
                        disabled={loading}
                      />
                    }
                    label={formik.values.isActive ? 'Active' : 'Inactive'}
                    sx={{ mt: 2 }}
                  />
                </Grid>
              </Grid>
            </Collapse>
          </Grid>

          {/* Categories Section */}
          <Grid item xs={12}>
            {renderSectionHeader('Categories', 'categories')}
            <Collapse in={expandedSections.categories}>
              <FormControl fullWidth margin="normal" error={formik.touched.categoryIds && Boolean(formik.errors.categoryIds)}>
                <InputLabel id="category-select-label">Categories</InputLabel>
                <Select
                  labelId="category-select-label"
                  id="categoryIds"
                  name="categoryIds"
                  multiple
                  value={formik.values.categoryIds}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  input={<OutlinedInput label="Categories" />}
                  renderValue={(selected) => (
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                      {selected.map((value) => {
                        const category = categories.find(c => c.id === value);
                        return (
                          <Chip 
                            key={value} 
                            label={category ? category.name : value} 
                            size="small" 
                          />
                        );
                      })}
                    </Box>
                  )}
                  disabled={loading}
                >
                  {categories.map((category) => (
                    <MenuItem key={category.id} value={category.id}>
                      <Checkbox checked={formik.values.categoryIds.indexOf(category.id) > -1} />
                      <ListItemText primary={category.name} />
                    </MenuItem>
                  ))}
                </Select>
                {formik.touched.categoryIds && formik.errors.categoryIds && (
                  <FormHelperText>{formik.errors.categoryIds}</FormHelperText>
                )}
              </FormControl>
            </Collapse>
          </Grid>

          {/* Supplier Section */}
          <Grid item xs={12}>
            {renderSectionHeader('Supplier', 'supplier')}
            <Collapse in={expandedSections.supplier}>
              <FormControl fullWidth margin="normal" error={formik.touched.supplierId && Boolean(formik.errors.supplierId)}>
                <InputLabel id="supplier-select-label">Supplier</InputLabel>
                <Select
                  labelId="supplier-select-label"
                  id="supplierId"
                  name="supplierId"
                  value={formik.values.supplierId}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  label="Supplier"
                  disabled={loading}
                >
                  {suppliers.map((supplier) => (
                    <MenuItem key={supplier.id} value={supplier.id}>
                      {supplier.name}
                    </MenuItem>
                  ))}
                </Select>
                {formik.touched.supplierId && formik.errors.supplierId && (
                  <FormHelperText>{formik.errors.supplierId}</FormHelperText>
                )}
              </FormControl>
            </Collapse>
          </Grid>

          {/* Media Section */}
          <Grid item xs={12}>
            {renderSectionHeader('Media', 'media')}
            <Collapse in={expandedSections.media}>
              <ArrayInput
                values={formik.values.imageUrls}
                onAdd={() => {
                  const newImageUrls = [...formik.values.imageUrls, ''];
                  formik.setFieldValue('imageUrls', newImageUrls);
                }}
                onRemove={(index) => {
                  const newImageUrls = [...formik.values.imageUrls];
                  newImageUrls.splice(index, 1);
                  formik.setFieldValue('imageUrls', newImageUrls);
                }}
                label="Image URLs"
                placeholder="https://example.com/image.jpg"
                helperText="Add URLs for product images"
                error={formik.touched.imageUrls && formik.errors.imageUrls}
                touched={formik.touched.imageUrls}
                icon={ImageIcon}
              />

              <ArrayInput
                values={formik.values.videoUrls}
                onAdd={() => {
                  const newVideoUrls = [...formik.values.videoUrls, ''];
                  formik.setFieldValue('videoUrls', newVideoUrls);
                }}
                onRemove={(index) => {
                  const newVideoUrls = [...formik.values.videoUrls];
                  newVideoUrls.splice(index, 1);
                  formik.setFieldValue('videoUrls', newVideoUrls);
                }}
                label="Video URLs"
                placeholder="https://example.com/video.mp4"
                helperText="Add URLs for product videos"
                error={formik.touched.videoUrls && formik.errors.videoUrls}
                touched={formik.touched.videoUrls}
                icon={VideoIcon}
              />
            </Collapse>
          </Grid>

          {/* Metadata Section */}
          <Grid item xs={12}>
            {renderSectionHeader('Metadata', 'metadata')}
            <Collapse in={expandedSections.metadata}>
              <MetadataInput
                values={formik.values.metadata || {}}
                onChange={(metadata) => formik.setFieldValue('metadata', metadata)}
                error={formik.touched.metadata && formik.errors.metadata}
                touched={formik.touched.metadata}
              />
            </Collapse>
          </Grid>
        </Grid>

        <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 2, mt: 4, p: 2 }}>
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
