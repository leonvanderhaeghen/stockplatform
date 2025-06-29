import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  Box,
  Typography,
  Paper,
  Grid,
  Button,
  Divider,
  Chip,
  Card,
  CardMedia,
  IconButton,
} from '@mui/material';
import {
  Edit as EditIcon,
  ArrowBack as ArrowBackIcon,
  Delete as DeleteIcon,
  Inventory as InventoryIcon,
  AttachMoney as PriceIcon,
  Category as CategoryIcon,
  Code as SkuIcon,
} from '@mui/icons-material';
import { format } from 'date-fns';
import productService from '../services/productService';
import inventoryService from '../services/inventoryService';
import ConfirmDialog from '../components/common/ConfirmDialog';

const ProductDetailPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [deleteDialogOpen, setDeleteDialogOpen] = React.useState(false);

  const { data: product, isLoading: productLoading, error } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productService.getProduct(id),
    enabled: !!id,
  });

  // Fetch inventory data for this product
  const { data: inventoryData, isLoading: inventoryLoading } = useQuery({
    queryKey: ['inventory-item', id],
    queryFn: () => inventoryService.getInventoryItem(id),
    enabled: !!id && !!product,
    retry: false, // Don't retry if inventory item doesn't exist
  });

  const isLoading = productLoading || inventoryLoading;

  const handleEdit = () => {
    navigate(`/products/edit/${id}`);
  };

  const handleDelete = () => {
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = () => {
    // Implement delete logic
    setDeleteDialogOpen(false);
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading product: {error.message}</div>;
  if (!product) return <div>Product not found</div>;

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <IconButton onClick={() => navigate(-1)}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h4" component="h1">
            {product.name}
          </Typography>
          <Chip
            label={inventoryData && inventoryData.quantity > 0 ? 'In Stock' : 'Out of Stock'}
            color={inventoryData && inventoryData.quantity > 0 ? 'success' : 'error'}
            size="small"
          />
        </Box>
        <Box>
          <Button
            variant="outlined"
            startIcon={<EditIcon />}
            onClick={handleEdit}
            sx={{ mr: 1 }}
          >
            Edit
          </Button>
          <Button
            variant="outlined"
            color="error"
            startIcon={<DeleteIcon />}
            onClick={handleDelete}
          >
            Delete
          </Button>
        </Box>
      </Box>

      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              Product Information
            </Typography>
            <Divider sx={{ mb: 3 }} />
            
            <Grid container spacing={3}>
              <Grid item xs={12} sm={6}>
                <DetailItem 
                  icon={<SkuIcon color="action" />} 
                  label="SKU" 
                  value={product.sku} 
                />
                <DetailItem 
                  icon={<CategoryIcon color="action" />} 
                  label="Category IDs" 
                  value={product.category_ids?.join(', ') || 'N/A'} 
                />
                <DetailItem 
                  icon={<PriceIcon color="action" />} 
                  label="Price" 
                  value={`$${parseFloat(product.selling_price).toFixed(2)}`} 
                />
                <DetailItem 
                  icon={<InventoryIcon color="action" />} 
                  label="Stock Quantity" 
                  value={inventoryData?.quantity || 0} 
                />
                <DetailItem 
                  label="Low Stock Threshold" 
                  value={inventoryData?.low_stock_threshold || 'N/A'} 
                />
              </Grid>
              <Grid item xs={12} sm={6}>
                <DetailItem 
                  label="Currency" 
                  value={product.currency || 'N/A'} 
                />
                <DetailItem 
                  label="Cost Price" 
                  value={`$${parseFloat(product.cost_price).toFixed(2)}`} 
                />
                <DetailItem 
                  label="Created" 
                  value={format(new Date(product.created_at.seconds * 1000), 'PPpp')} 
                />
                <DetailItem 
                  label="Last Updated" 
                  value={format(new Date(product.updated_at.seconds * 1000), 'PPpp')} 
                />
              </Grid>
            </Grid>

            <Box sx={{ mt: 3 }}>
              <Typography variant="subtitle1" gutterBottom>
                Description
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {product.description || 'No description available.'}
              </Typography>
            </Box>
          </Paper>
        </Grid>

        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              Product Images
            </Typography>
            <Divider sx={{ mb: 3 }} />
            
            {product.images?.length > 0 ? (
              <Grid container spacing={2}>
                {product.images.map((image, index) => (
                  <Grid item xs={6} key={index}>
                    <Card>
                      <CardMedia
                        component="img"
                        height="140"
                        image={image}
                        alt={`Product ${index + 1}`}
                      />
                    </Card>
                  </Grid>
                ))}
              </Grid>
            ) : (
              <Box
                sx={{
                  height: 200,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  bgcolor: 'action.hover',
                  borderRadius: 1,
                }}
              >
                <Typography color="text.secondary">No images available</Typography>
              </Box>
            )}
          </Paper>
        </Grid>
      </Grid>

      <ConfirmDialog
        open={deleteDialogOpen}
        title="Delete Product"
        content={`Are you sure you want to delete "${product.name}"? This action cannot be undone.`}
        onClose={() => setDeleteDialogOpen(false)}
        onConfirm={handleDeleteConfirm}
        confirmText="Delete"
        confirmColor="error"
      />
    </Box>
  );
};

const DetailItem = ({ icon, label, value }) => (
  <Box sx={{ display: 'flex', mb: 1.5 }}>
    <Box sx={{ display: 'flex', alignItems: 'center', mr: 1.5, color: 'text.secondary' }}>
      {icon}
    </Box>
    <Box>
      <Typography variant="caption" color="text.secondary">
        {label}
      </Typography>
      <Typography variant="body1">
        {value}
      </Typography>
    </Box>
  </Box>
);

export default ProductDetailPage;
