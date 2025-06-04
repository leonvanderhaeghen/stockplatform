import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  Typography,
} from '@mui/material';
import { Add as AddIcon, Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { productService } from '../../services';
import ProductForm from './ProductForm';

const ProductsCRUD = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [error, setError] = useState('');
  const [categories, setCategories] = useState([]);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [productsResponse, categoriesResponse] = await Promise.all([
        productService.getProducts(),
        productService.getCategories()
      ]);
      
      // Handle products response - check if it has a data.products array
      const productsData = Array.isArray(productsResponse.data?.products) 
        ? productsResponse.data.products 
        : [];
      
      // Handle categories response - check if it has a data.categories array
      const categoriesData = Array.isArray(categoriesResponse.data?.categories) 
        ? categoriesResponse.data.categories 
        : [];
      
      setProducts(productsData);
      setCategories(categoriesData);
    } catch (err) {
      setError('Failed to fetch data');
      console.error('Error fetching data:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleOpen = (product = null) => {
    setSelectedProduct(product);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setSelectedProduct(null);
  };

  const handleSubmit = async (productData) => {
    try {
      if (selectedProduct) {
        await productService.updateProduct(selectedProduct.id, productData);
      } else {
        await productService.createProduct(productData);
      }
      await fetchData();
      handleClose();
    } catch (err) {
      setError(`Failed to ${selectedProduct ? 'update' : 'create'} product: ${err.message}`);
      console.error(err);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this product?')) {
      try {
        await productService.deleteProduct(id);
        await fetchData();
      } catch (err) {
        setError(`Failed to delete product: ${err.message}`);
        console.error(err);
      }
    }
  };

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Typography color="error">{error}</Typography>;

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" mb={2}>
        <Typography variant="h4">Products</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => handleOpen()}
        >
          Add Product
        </Button>
      </Box>

      <TableContainer component={Paper}>
        {loading ? (
          <Box p={4} textAlign="center">
            <Typography>Loading products...</Typography>
          </Box>
        ) : error ? (
          <Box p={4} textAlign="center" color="error.main">
            <Typography>{error}</Typography>
            <Button onClick={fetchData} color="primary" variant="outlined" sx={{ mt: 2 }}>
              Retry
            </Button>
          </Box>
        ) : products.length === 0 ? (
          <Box p={4} textAlign="center">
            <Typography>No products found</Typography>
            <Button 
              onClick={() => handleOpen()} 
              variant="contained" 
              color="primary" 
              startIcon={<AddIcon />}
              sx={{ mt: 2 }}
            >
              Add Your First Product
            </Button>
          </Box>
        ) : (
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Description</TableCell>
                <TableCell>Price</TableCell>
                <TableCell>Stock</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {products.map((product) => (
                <TableRow key={product.id}>
                  <TableCell>{product.name}</TableCell>
                  <TableCell>{product.description}</TableCell>
                  <TableCell>${product.price?.toFixed(2)}</TableCell>
                  <TableCell>{product.stock}</TableCell>
                  <TableCell>
                    <IconButton onClick={() => handleOpen(product)}>
                      <EditIcon />
                    </IconButton>
                    <IconButton onClick={() => handleDelete(product.id)}>
                      <DeleteIcon color="error" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </TableContainer>

      <Dialog open={open} onClose={handleClose} maxWidth="md" fullWidth>
        <DialogTitle>{selectedProduct ? 'Edit Product' : 'Add New Product'}</DialogTitle>
        <DialogContent>
          <ProductForm
            initialValues={selectedProduct || { 
              name: '', 
              description: '', 
              price: 0, 
              cost: 0,
              sku: '',
              category: '',
              inStock: true,
              stock: 0,
              images: []
            }}
            onSubmit={handleSubmit}
            onCancel={handleClose}
            loading={loading}
            isEdit={!!selectedProduct}
            categories={categories}
          />
        </DialogContent>
      </Dialog>
    </Box>
  );
};

export default ProductsCRUD;
