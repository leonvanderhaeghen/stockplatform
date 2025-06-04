import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import { Box, Typography, Paper, Button, Container, CircularProgress } from '@mui/material';
import { ArrowBack as ArrowBackIcon, Save as SaveIcon } from '@mui/icons-material';
import ProductForm from '../components/products/ProductForm';
import productService from '../services/productService';

// Mock categories - in a real app, this would come from an API
const mockCategories = [
  { id: '1', name: 'Electronics' },
  { id: '2', name: 'Clothing' },
  { id: '3', name: 'Books' },
  { id: '4', name: 'Home & Kitchen' },
  { id: '5', name: 'Sports & Outdoors' },
];

const ProductFormPage = ({ isEdit = false }) => {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [initialValues, setInitialValues] = useState(null);

  // Fetch product data if in edit mode
  const { data: product, isLoading: isLoadingProduct } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productService.getProduct(id),
    enabled: isEdit && !!id,
    onError: (error) => {
      toast.error(error.response?.data?.message || 'Failed to load product');
      navigate('/products');
    },
  });

  // Set initial values when product data is loaded
  useEffect(() => {
    if (isEdit && product) {
      setInitialValues({
        name: product.name || '',
        sku: product.sku || '',
        description: product.description || '',
        price: product.price || 0,
        cost: product.cost || 0,
        category: product.category?.id || '',
        inStock: product.inStock !== undefined ? product.inStock : true,
        stock: product.stock || 0,
        images: product.images || [],
      });
    } else if (!isEdit) {
      // Set default values for new product
      setInitialValues({
        name: '',
        sku: '',
        description: '',
        price: 0,
        cost: 0,
        category: '',
        inStock: true,
        stock: 0,
        images: [],
      });
    }
  }, [isEdit, product]);

  // Create or update product mutation
  const mutation = useMutation({
    mutationFn: (values) => {
      const productData = {
        ...values,
        // Convert string numbers to numbers
        price: parseFloat(values.price),
        cost: parseFloat(values.cost),
        stock: parseInt(values.stock, 10),
        inStock: values.inStock || values.stock > 0,
      };

      return isEdit
        ? productService.updateProduct(id, productData)
        : productService.createProduct(productData);
    },
    onSuccess: () => {
      const action = isEdit ? 'updated' : 'created';
      toast.success(`Product ${action} successfully`);
      queryClient.invalidateQueries({ queryKey: ['products'] });
      navigate('/products');
    },
    onError: (error) => {
      toast.error(error.response?.data?.message || `Failed to ${isEdit ? 'update' : 'create'} product`);
    },
  });

  const handleSubmit = (values) => {
    mutation.mutate(values);
  };

  if (isLoadingProduct || initialValues === null) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box sx={{ mb: 4 }}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate(-1)}
          sx={{ mb: 2 }}
        >
          Back
        </Button>
        
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h4" component="h1">
            {isEdit ? 'Edit Product' : 'Add New Product'}
          </Typography>
          <Box>
            <Button
              variant="contained"
              color="primary"
              startIcon={<SaveIcon />}
              form="product-form"
              type="submit"
              disabled={mutation.isLoading}
              sx={{ ml: 1 }}
            >
              {mutation.isLoading ? 'Saving...' : 'Save Product'}
            </Button>
          </Box>
        </Box>
        
        <Paper elevation={2} sx={{ p: { xs: 2, md: 4 } }}>
          <ProductForm
            initialValues={initialValues}
            onSubmit={handleSubmit}
            loading={mutation.isLoading}
            isEdit={isEdit}
            categories={mockCategories}
          />
        </Paper>
      </Box>
    </Container>
  );
};

export default ProductFormPage;
