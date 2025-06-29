import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import { Box, Typography, Paper, Button, Container, CircularProgress } from '@mui/material';
import { ArrowBack as ArrowBackIcon, Save as SaveIcon } from '@mui/icons-material';
import ProductForm from '../components/products/ProductForm';
import productService from '../services/productService';


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

  // Fetch categories from API
  const { data: categoriesData, isLoading: isLoadingCategories, error: categoriesError } = useQuery({
    queryKey: ['categories'],
    queryFn: async () => {
      const res = await productService.getCategories();
      // Defensive: API may return { data: { categories: [...] } } or { categories: [...] } or [...]
      if (Array.isArray(res)) return res;
      if (Array.isArray(res?.data?.categories)) return res.data.categories;
      if (Array.isArray(res?.categories)) return res.categories;
      return [];
    },
  });

  // Fetch suppliers from API
  const { data: suppliersData, isLoading: isLoadingSuppliers, error: suppliersError } = useQuery({
    queryKey: ['suppliers'],
    queryFn: async () => {
      const { supplierService } = await import('../services');
      const res = await supplierService.getSuppliers();
      // Defensive: API may return { suppliers: [...] } or { data: { suppliers: [...] } } or [...]
      if (Array.isArray(res)) return res;
      if (Array.isArray(res?.suppliers)) return res.suppliers;
      if (Array.isArray(res?.data?.suppliers)) return res.data.suppliers;
      return [];
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

  if (isLoadingProduct || initialValues === null || isLoadingCategories || isLoadingSuppliers) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (categoriesError) {
    return <Box color="error.main">Failed to load categories: {categoriesError.message}</Box>;
  }
  if (suppliersError) {
    return <Box color="error.main">Failed to load suppliers: {suppliersError.message}</Box>;
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
            categories={categoriesData}
            suppliers={suppliersData}
          />
        </Paper>
      </Box>
    </Container>
  );
};

export default ProductFormPage;
