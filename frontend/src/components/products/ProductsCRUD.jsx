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
import { productService, supplierService } from '../../services';
import ProductForm from './ProductForm';

const ProductsCRUD = () => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [error, setError] = useState('');
  const [categories, setCategories] = useState([]);
  const [suppliers, setSuppliers] = useState([]);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [productsResponse, categoriesResponse, suppliersResponse] = await Promise.all([
        productService.getProducts(),
        productService.getCategories(),
        supplierService.getSuppliers()
      ]);
      
      // Handle products response - check if it has a data.products array
      const productsData = Array.isArray(productsResponse.data?.products) 
        ? productsResponse.data.products 
        : [];
      
      // Handle categories response - check if it has a data.categories array
      const categoriesData = Array.isArray(categoriesResponse.data?.categories) 
        ? categoriesResponse.data.categories 
        : [];

      // Handle suppliers response - check various formats
      let suppliersData = [];
      if (Array.isArray(suppliersResponse)) {
        suppliersData = suppliersResponse;
      } else if (Array.isArray(suppliersResponse?.suppliers)) {
        suppliersData = suppliersResponse.suppliers;
      } else if (Array.isArray(suppliersResponse?.data?.suppliers)) {
        suppliersData = suppliersResponse.data.suppliers;
      }
      
      setProducts(productsData);
      setCategories(categoriesData);
      setSuppliers(suppliersData);
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
      // Transform field names to match backend's expected snake_case format
      const transformedProductData = {};
      
      // Field mapping from camelCase to snake_case
      const fieldMappings = {
        name: 'name',
        description: 'description',
        costPrice: 'cost_price',
        sellingPrice: 'selling_price',
        sku: 'sku',
        barcode: 'barcode',
        categoryIds: 'category_ids',
        supplierId: 'supplier_id',
        isActive: 'is_active',
        inStock: 'in_stock',
        stockQty: 'stock_qty',
        lowStockAt: 'low_stock_at',
        imageUrls: 'image_urls',
        videoUrls: 'video_urls',
        metadata: 'metadata',
        currency: 'currency'
      };
      
      // Process all fields with proper snake_case naming for backend
      Object.entries(productData).forEach(([key, value]) => {
        // Skip field entry if there's no value
        if (value === undefined || value === null) return;
        
        // Convert to snake_case using mapping or fallback to camelCase to snake_case conversion
        const snakeCaseKey = fieldMappings[key] || key.replace(/([A-Z])/g, '_$1').toLowerCase();
        
        // For price fields, ensure they are properly formatted strings
        if (key === 'costPrice' || key === 'sellingPrice') {
          // Make sure price is a string with exactly 2 decimal places
          const priceValue = typeof value === 'number' ? value.toFixed(2) : 
                           (typeof value === 'string' ? parseFloat(value).toFixed(2) : '0.00');
          transformedProductData[snakeCaseKey] = priceValue;
        } else {
          transformedProductData[snakeCaseKey] = value;
        }
      });
      
      // Debug payload
      console.log('Product payload before submission:', transformedProductData);
      
      try {
        if (selectedProduct) {
          await productService.updateProduct(selectedProduct.id, transformedProductData);
        } else {
          await productService.createProduct(transformedProductData);
        }
        await fetchData();
        handleClose();
      } catch (error) {
        console.error('Error saving product:', error);
        
        // Extract meaningful error message
        let errorMessage = 'Failed to save product';
        
        if (error.response?.data?.error) {
          errorMessage = error.response.data.error;
        } else if (typeof error.message === 'string') {
          errorMessage = error.message;
        }
        
        // Check for supplier not found error
        if (errorMessage.includes('supplier not found') || errorMessage.toLowerCase().includes('supplier')) {
          setError(`Supplier error: The selected supplier ID (${transformedProductData.supplier_id}) could not be found. Please select a different supplier.`);
        } else {
          setError(`Error: ${errorMessage}`);
        }
      }
    } catch (err) {
      setError(`Failed to ${selectedProduct ? 'update' : 'create'} product: ${err.message}`);
      console.error('Product creation error:', err);
      console.log('Product data sent:', productData);
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
              costPrice: '0.00', 
              sellingPrice: '0.00',
              sku: '',
              barcode: '',
              category: '',
              inStock: true,
              stockQty: 0,
              images: []
            }}
            onSubmit={handleSubmit}
            onCancel={handleClose}
            loading={loading}
            isEdit={!!selectedProduct}
            categories={categories}
            suppliers={suppliers}
          />
        </DialogContent>
      </Dialog>
    </Box>
  );
};

export default ProductsCRUD;
