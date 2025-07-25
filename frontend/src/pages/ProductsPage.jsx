import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Chip } from '@mui/material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'react-toastify';
import DataTable from '../components/common/DataTable';
import productService from '../services/productService';
import { useInventoryItems } from '../utils/hooks/useInventory';
import ConfirmDialog from '../components/common/ConfirmDialog';
import InventoryStatus from '../components/inventory/InventoryStatus';
import InventoryErrorBoundary from '../components/inventory/InventoryErrorBoundary';

const ProductsPage = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);

  // Fetch products with proper data extraction
  const { data: products, isLoading: productsLoading } = useQuery({
    queryKey: ['products'],
    queryFn: () => productService.getProducts(),
    select: (data) => {
      // Ensure we always return an array, even if data is missing or malformed
      if (!data) return [];
      if (Array.isArray(data)) return data;
      if (Array.isArray(data.products)) return data.products;
      if (Array.isArray(data.data?.products)) return data.data.products;
      return [];
    },
  });

  // Fetch inventory items to join with products using the new hook
  const { data: inventoryItems = [], isLoading: inventoryLoading } = useInventoryItems();

  // Create a map of product ID to inventory data for efficient lookup
  const inventoryMap = useMemo(() => {
    const map = new Map();
    inventoryItems.forEach(item => {
      if (item.product_id) {
        map.set(item.product_id, item);
      }
    });
    return map;
  }, [inventoryItems]);

  // Combine products with their inventory data
  const productsWithInventory = useMemo(() => {
    return products.map(product => {
      const inventory = inventoryMap.get(product.id);
      // Ensure we have a numeric price field for the table and formatting
      const sellingPrice = product.selling_price ?? product.sellingPrice ?? product.price;
      let priceNumber = 0;
      if (typeof sellingPrice === 'number') {
        priceNumber = sellingPrice;
      } else if (typeof sellingPrice === 'string') {
        priceNumber = parseFloat(sellingPrice);
      }
      return {
        ...product,
        price: priceNumber,
        inventory: inventory || null,
        stockQty: inventory?.quantity || 0,
        inStock: inventory ? inventory.quantity > 0 : false,
      };
    });
  }, [products, inventoryMap]);

  const isLoading = productsLoading || inventoryLoading;

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
      valueFormatter: (params) => {
        const val = params.value;
        if (val === undefined || val === null || Number.isNaN(val)) {
          return '—';
        }
        return `$${Number(val).toFixed(2)}`;
      },
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
      field: 'stockQty', 
      headerName: 'Stock', 
      width: 100,
      renderCell: (params) => {
        const inventoryItem = inventoryMap.get(params.id);
        return (
          <InventoryErrorBoundary>
            <InventoryStatus 
              quantity={inventoryItem?.quantity}
              lowStockThreshold={inventoryItem?.low_stock_threshold}
            />
          </InventoryErrorBoundary>
        );
      },
    },
    { 
      field: 'inStock', 
      headerName: 'Status', 
      width: 120,
      renderCell: (params) => {
        const inventoryItem = inventoryMap.get(params.id);
        return (
          <InventoryErrorBoundary>
            <Chip
              label={inventoryItem?.quantity || 0}
              color={inventoryItem && inventoryItem.quantity > 0 ? 'success' : 'default'}
              size="small"
              variant="outlined"
            />
          </InventoryErrorBoundary>
        );
      },
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
        rows={productsWithInventory}
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
