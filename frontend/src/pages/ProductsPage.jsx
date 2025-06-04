import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Chip, IconButton, Tooltip } from '@mui/material';
import { Add as AddIcon, Edit as EditIcon, Delete as DeleteIcon, Visibility as ViewIcon } from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import DataTable from '../components/common/DataTable';
import productService from '../services/productService';
import ConfirmDialog from '../components/common/ConfirmDialog';

const ProductsPage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);

  // Fetch products
  const { data: products = [], isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: () => productService.getProducts(),
  });

  // Delete product mutation
  const deleteMutation = useMutation({
    mutationFn: (id) => productService.deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] });
      toast.success('Product deleted successfully');
      setDeleteDialogOpen(false);
    },
    onError: (error) => {
      toast.error(error.response?.data?.message || 'Failed to delete product');
    },
  });

  const columns = [
    { 
      field: 'name', 
      headerName: 'Name', 
      flex: 1,
      minWidth: 200,
    },
    { 
      field: 'sku', 
      headerName: 'SKU', 
      width: 150,
    },
    { 
      field: 'price', 
      headerName: 'Price', 
      width: 120,
      valueFormatter: (params) => `$${params.value.toFixed(2)}`,
    },
    { 
      field: 'category', 
      headerName: 'Category', 
      width: 150,
      renderCell: (params) => (
        <Chip 
          label={params.value} 
          size="small" 
          color="primary"
          variant="outlined"
        />
      ),
    },
    { 
      field: 'inStock', 
      headerName: 'Status', 
      width: 120,
      renderCell: (params) => (
        <Chip 
          label={params.value ? 'In Stock' : 'Out of Stock'} 
          color={params.value ? 'success' : 'error'} 
          size="small"
        />
      ),
    },
  ];

  const handleAddProduct = () => {
    navigate('/products/new');
  };

  const handleEditProduct = (product) => {
    navigate(`/products/edit/${product.id}`);
  };

  const handleViewProduct = (product) => {
    navigate(`/products/${product.id}`);
  };

  const handleDeleteClick = (product) => {
    setSelectedProduct(product);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    if (selectedProduct) {
      deleteMutation.mutate(selectedProduct.id);
    }
  };

  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <DataTable
        rows={products}
        columns={columns}
        loading={isLoading || deleteMutation.isPending}
        title="Products"
        onAdd={handleAddProduct}
        onEdit={handleEditProduct}
        onView={handleViewProduct}
        onDelete={handleDeleteClick}
        pageSize={10}
        pageSizeOptions={[5, 10, 25]}
        getRowId={(row) => row.id}
      />

      <ConfirmDialog
        open={deleteDialogOpen}
        title="Delete Product"
        content={`Are you sure you want to delete "${selectedProduct?.name}"?`}
        onClose={() => setDeleteDialogOpen(false)}
        onConfirm={handleDeleteConfirm}
        confirmText="Delete"
        confirmColor="error"
        loading={deleteMutation.isPending}
      />
    </Box>
  );
};

export default ProductsPage;
