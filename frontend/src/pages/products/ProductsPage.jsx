import React, { useState, useCallback } from 'react';
import {
  Container,
  Typography,
  Box,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CardActions,
  Button,
  IconButton,
  TextField,
  InputAdornment,
  Chip,
  Menu,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
  CircularProgress,
  Fab,
  Pagination,
  FormControl,
  InputLabel,
  Select,
  Drawer,
  Divider,
  Rating,
  Skeleton,
  Switch,
  FormControlLabel,
  Autocomplete
} from '@mui/material';
import {
  Inventory as InventoryIcon,
  Search,
  FilterList,
  Add,
  Edit,
  Delete,
  MoreVert,
  Visibility,
  ShoppingCart,
  Close,
  Category,
  GridView,
  ViewList
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import productService from '../../services/productService';
import categoryService from '../../services/categoryService';
import supplierService from '../../services/supplierService';
import { useSnackbar } from 'notistack';
import { format } from 'date-fns';

const ProductCard = ({ product, onEdit, onDelete, onView, isStaff }) => {
  const [anchorEl, setAnchorEl] = useState(null);

  const handleMenuClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleEdit = () => {
    onEdit(product);
    handleMenuClose();
  };

  const handleDelete = () => {
    onDelete(product);
    handleMenuClose();
  };

  const handleView = () => {
    onView(product);
    handleMenuClose();
  };

  return (
    <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <CardMedia
        component="img"
        height="200"
        image={product.imageUrl || '/placeholder-product.jpg'}
        alt={product.name}
        sx={{ objectFit: 'cover' }}
      />
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography gutterBottom variant="h6" component="div" noWrap>
          {product.name}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
          SKU: {product.sku}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
          {product.description}
        </Typography>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
          <Typography variant="h6" color="primary">
            ${product.price}
          </Typography>
          <Chip 
            label={product.category} 
            size="small" 
            color="primary" 
            variant="outlined"
          />
        </Box>
        {product.rating && (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Rating value={product.rating} precision={0.1} size="small" readOnly />
            <Typography variant="caption" color="text.secondary">
              ({product.reviewCount || 0})
            </Typography>
          </Box>
        )}
        <Box sx={{ display: 'flex', gap: 1, mt: 1, flexWrap: 'wrap' }}>
          <Chip 
            label={product.status === 'ACTIVE' ? 'Active' : 'Inactive'} 
            size="small" 
            color={product.status === 'ACTIVE' ? 'success' : 'default'}
          />
          {product.stockQuantity !== undefined && (
            <Chip 
              label={`Stock: ${product.stockQuantity}`} 
              size="small" 
              color={product.stockQuantity > 10 ? 'success' : product.stockQuantity > 0 ? 'warning' : 'error'}
            />
          )}
        </Box>
      </CardContent>
      <CardActions sx={{ justifyContent: 'space-between' }}>
        <Button size="small" startIcon={<Visibility />} onClick={handleView}>
          View
        </Button>
        {isStaff && (
          <>
            <Button size="small" startIcon={<ShoppingCart />}>
              Add to Cart
            </Button>
            <IconButton size="small" onClick={handleMenuClick}>
              <MoreVert />
            </IconButton>
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleMenuClose}
            >
              <MenuItem onClick={handleEdit}>
                <Edit sx={{ mr: 1 }} /> Edit
              </MenuItem>
              <MenuItem onClick={handleDelete} sx={{ color: 'error.main' }}>
                <Delete sx={{ mr: 1 }} /> Delete
              </MenuItem>
            </Menu>
          </>
        )}
      </CardActions>
    </Card>
  );
};

const ProductDetailDialog = ({ open, onClose, product }) => {
  if (!product) return null;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h5">{product.name}</Typography>
        <IconButton onClick={onClose}>
          <Close />
        </IconButton>
      </DialogTitle>
      <DialogContent>
        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Box
              component="img"
              src={product.imageUrl || '/placeholder-product.jpg'}
              alt={product.name}
              sx={{ width: '100%', height: 300, objectFit: 'cover', borderRadius: 1 }}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Typography variant="h4" color="primary" gutterBottom>
              ${product.price}
            </Typography>
            <Typography variant="body1" paragraph>
              {product.description}
            </Typography>
            <Box sx={{ mb: 2 }}>
              <Typography variant="subtitle2" gutterBottom>Product Details</Typography>
              <Typography variant="body2">SKU: {product.sku}</Typography>
              <Typography variant="body2">Category: {product.category}</Typography>
              <Typography variant="body2">Status: {product.status}</Typography>
              {product.stockQuantity !== undefined && (
                <Typography variant="body2">Stock: {product.stockQuantity}</Typography>
              )}
              {product.weight && (
                <Typography variant="body2">Weight: {product.weight} lbs</Typography>
              )}
              {product.dimensions && (
                <Typography variant="body2">Dimensions: {product.dimensions}</Typography>
              )}
            </Box>
            {product.rating && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                <Rating value={product.rating} precision={0.1} readOnly />
                <Typography variant="body2" color="text.secondary">
                  {product.rating} ({product.reviewCount || 0} reviews)
                </Typography>
              </Box>
            )}
            <Typography variant="caption" color="text.secondary">
              Created: {format(new Date(product.createdAt || Date.now()), 'MMM d, yyyy')}
            </Typography>
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
        <Button variant="contained" startIcon={<ShoppingCart />}>
          Add to Cart
        </Button>
      </DialogActions>
    </Dialog>
  );
};

const FilterDrawer = ({ open, onClose, filters, onFiltersChange, categories }) => {
  return (
    <Drawer anchor="right" open={open} onClose={onClose}>
      <Box sx={{ width: 300, p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Filters
        </Typography>
        <Divider sx={{ mb: 3 }} />
        
        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel>Category</InputLabel>
          <Select
            value={filters.category || ''}
            label="Category"
            onChange={(e) => onFiltersChange({ ...filters, category: e.target.value })}
          >
            <MenuItem value="">All Categories</MenuItem>
            {categories.map((category) => (
              <MenuItem key={category.id} value={category.name}>
                {category.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel>Status</InputLabel>
          <Select
            value={filters.status || ''}
            label="Status"
            onChange={(e) => onFiltersChange({ ...filters, status: e.target.value })}
          >
            <MenuItem value="">All Status</MenuItem>
            <MenuItem value="ACTIVE">Active</MenuItem>
            <MenuItem value="INACTIVE">Inactive</MenuItem>
          </Select>
        </FormControl>

        <Typography variant="subtitle2" gutterBottom>
          Price Range
        </Typography>
        <Grid container spacing={1} sx={{ mb: 3 }}>
          <Grid item xs={6}>
            <TextField
              fullWidth
              label="Min Price"
              type="number"
              value={filters.minPrice || ''}
              onChange={(e) => onFiltersChange({ ...filters, minPrice: e.target.value })}
            />
          </Grid>
          <Grid item xs={6}>
            <TextField
              fullWidth
              label="Max Price"
              type="number"
              value={filters.maxPrice || ''}
              onChange={(e) => onFiltersChange({ ...filters, maxPrice: e.target.value })}
            />
          </Grid>
        </Grid>

        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel>Sort By</InputLabel>
          <Select
            value={filters.sortBy || 'name'}
            label="Sort By"
            onChange={(e) => onFiltersChange({ ...filters, sortBy: e.target.value })}
          >
            <MenuItem value="name">Name</MenuItem>
            <MenuItem value="price">Price</MenuItem>
            <MenuItem value="createdAt">Date Created</MenuItem>
            <MenuItem value="rating">Rating</MenuItem>
          </Select>
        </FormControl>

        <FormControl fullWidth>
          <InputLabel>Sort Order</InputLabel>
          <Select
            value={filters.sortOrder || 'asc'}
            label="Sort Order"
            onChange={(e) => onFiltersChange({ ...filters, sortOrder: e.target.value })}
          >
            <MenuItem value="asc">Ascending</MenuItem>
            <MenuItem value="desc">Descending</MenuItem>
          </Select>
        </FormControl>

        <Box sx={{ mt: 3, display: 'flex', gap: 1 }}>
          <Button 
            fullWidth 
            variant="outlined" 
            onClick={() => onFiltersChange({})}
          >
            Clear Filters
          </Button>
          <Button fullWidth variant="contained" onClick={onClose}>
            Apply
          </Button>
        </Box>
      </Box>
    </Drawer>
  );
};

const ProductsPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();
  
  const [searchQuery, setSearchQuery] = useState('');
  const [filters, setFilters] = useState({});
  const [page, setPage] = useState(1);
  const [viewMode, setViewMode] = useState('grid'); // 'grid' or 'list'
  const [filterDrawerOpen, setFilterDrawerOpen] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [productDetailOpen, setProductDetailOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [productToDelete, setProductToDelete] = useState(null);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);
  const [productFormData, setProductFormData] = useState({
    name: '',
    description: '',
    sku: '',
    category_ids: [],
    supplier_id: '',
    cost_price: '',
    selling_price: '',
    currency: 'USD',
    is_active: true
  });

  const isStaff = user?.role === 'STAFF' || user?.role === 'ADMIN';
  const itemsPerPage = 12;

  // Query for products
  const { data: productsData, isLoading: productsLoading, error: productsError } = useQuery({
    queryKey: ['products', searchQuery, filters, page],
    queryFn: () => productService.listProducts({
      search: searchQuery,
      ...filters,
      page,
      limit: itemsPerPage
    }),
    staleTime: 30000,
  });

  // Fetch categories for filters
  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryService.getCategories(),
    staleTime: 300000, // 5 minutes
  });

  // Fetch suppliers for product creation
  const { data: suppliers = [] } = useQuery({
    queryKey: ['suppliers'],
    queryFn: () => supplierService.getSuppliersArray(),
    staleTime: 300000, // 5 minutes
  });

  // Create product mutation
  const createProductMutation = useMutation({
    mutationFn: (productData) => productService.createProduct(productData),
    onSuccess: () => {
      queryClient.invalidateQueries(['products']);
      enqueueSnackbar('Product created successfully', { variant: 'success' });
      setCreateDialogOpen(false);
      resetProductForm();
    },
    onError: (error) => {
      enqueueSnackbar(`Failed to create product: ${error.message}`, { variant: 'error' });
    },
  });

  // Update product mutation
  const updateProductMutation = useMutation({
    mutationFn: ({ id, data }) => productService.updateProduct(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries(['products']);
      enqueueSnackbar('Product updated successfully', { variant: 'success' });
      setCreateDialogOpen(false);
      resetProductForm();
    },
    onError: (error) => {
      enqueueSnackbar(`Failed to update product: ${error.message}`, { variant: 'error' });
    },
  });

  // Delete product mutation
  const deleteProductMutation = useMutation({
    mutationFn: (productId) => productService.deleteProduct(productId),
    onSuccess: () => {
      queryClient.invalidateQueries(['products']);
      enqueueSnackbar('Product deleted successfully', { variant: 'success' });
      setDeleteDialogOpen(false);
      setProductToDelete(null);
    },
    onError: (error) => {
      enqueueSnackbar(`Failed to delete product: ${error.message}`, { variant: 'error' });
    },
  });

  const handleSearch = useCallback((query) => {
    setSearchQuery(query);
    setPage(1);
  }, []);

  const handleFiltersChange = useCallback((newFilters) => {
    setFilters(newFilters);
    setPage(1);
  }, []);

  const handleProductView = (product) => {
    setSelectedProduct(product);
    setProductDetailOpen(true);
  };

  const handleProductEdit = (product) => {
    setEditingProduct(product);
    setProductFormData({
      name: product.name || '',
      description: product.description || '',
      sku: product.sku || '',
      category_ids: product.category_ids || [],
      supplier_id: product.supplier_id || '',
      cost_price: product.cost_price || '',
      selling_price: product.selling_price || '',
      currency: product.currency || 'USD',
      is_active: product.is_active !== undefined ? product.is_active : true
    });
    setCreateDialogOpen(true);
  };

  const handleCreateProduct = () => {
    setEditingProduct(null);
    resetProductForm();
    setCreateDialogOpen(true);
  };

  const resetProductForm = () => {
    setProductFormData({
      name: '',
      description: '',
      sku: '',
      category_ids: [],
      supplier_id: '',
      cost_price: '',
      selling_price: '',
      currency: 'USD',
      is_active: true
    });
  };

  const handleProductFormSubmit = useCallback(() => {
    if (editingProduct) {
      updateProductMutation.mutate({ id: editingProduct.id, data: productFormData });
    } else {
      createProductMutation.mutate(productFormData);
    }
  }, [productFormData, editingProduct, createProductMutation, updateProductMutation]);

  const handleProductDelete = (product) => {
    setProductToDelete(product);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (productToDelete) {
      deleteProductMutation.mutate(productToDelete.id);
    }
  };

  const handlePageChange = (event, newPage) => {
    setPage(newPage);
  };

  const products = productsData?.data || [];
  const totalPages = Math.ceil((productsData?.total || 0) / itemsPerPage);

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      {/* Header */}
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box>
          <Typography variant="h4" component="h1" gutterBottom>
            <InventoryIcon sx={{ mr: 1, verticalAlign: 'middle' }} />
            Products
          </Typography>
          <Typography variant="body1" color="text.secondary">
            {productsData?.total || 0} products available
          </Typography>
        </Box>
        {isStaff && (
          <Fab
            color="primary"
            aria-label="add product"
            onClick={handleCreateProduct}
          >
            <Add />
          </Fab>
        )}
      </Box>

      {/* Search and Controls */}
      <Box sx={{ mb: 3 }}>
        <Grid container spacing={2} alignItems="center">
          <Grid item xs={12} md={6}>
            <TextField
              fullWidth
              placeholder="Search products..."
              value={searchQuery}
              onChange={(e) => handleSearch(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search />
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
              <Button
                startIcon={<Category />}
                onClick={() => navigate('/categories')}
                variant="outlined"
              >
                Categories
              </Button>
              <Button
                startIcon={<FilterList />}
                onClick={() => setFilterDrawerOpen(true)}
                variant="outlined"
              >
                Filters
              </Button>
              <IconButton
                onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
                color="primary"
              >
                {viewMode === 'grid' ? <ViewList /> : <GridView />}
              </IconButton>
            </Box>
          </Grid>
        </Grid>
      </Box>

      {/* Active Filters */}
      {Object.keys(filters).length > 0 && (
        <Box sx={{ mb: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          {Object.entries(filters).map(([key, value]) => {
            if (!value) return null;
            return (
              <Chip
                key={key}
                label={`${key}: ${value}`}
                onDelete={() => {
                  const newFilters = { ...filters };
                  delete newFilters[key];
                  handleFiltersChange(newFilters);
                }}
                size="small"
                color="primary"
              />
            );
          })}
        </Box>
      )}

      {/* Error State */}
      {productsError && (
        <Alert severity="error" sx={{ mb: 3 }}>
          Failed to load products. Please try again.
        </Alert>
      )}

      {/* Loading State */}
      {productsLoading ? (
        <Grid container spacing={3}>
          {Array.from({ length: 8 }).map((_, index) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={index}>
              <Card>
                <Skeleton variant="rectangular" height={200} />
                <CardContent>
                  <Skeleton variant="text" height={32} />
                  <Skeleton variant="text" height={20} />
                  <Skeleton variant="text" height={20} width="60%" />
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      ) : products.length === 0 ? (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No products found
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            {searchQuery || Object.keys(filters).length > 0
              ? 'Try adjusting your search or filters'
              : 'Start by adding your first product'}
          </Typography>
          {isStaff && (
            <Button
              variant="contained"
              startIcon={<Add />}
              onClick={handleCreateProduct}
            >
              Add Product
            </Button>
          )}
        </Box>
      ) : (
        /* Products Grid */
        <>
          <Grid container spacing={3}>
            {products.map((product) => (
              <Grid 
                item 
                xs={12} 
                sm={viewMode === 'grid' ? 6 : 12} 
                md={viewMode === 'grid' ? 4 : 12} 
                lg={viewMode === 'grid' ? 3 : 12} 
                key={product.id}
              >
                <ProductCard
                  product={product}
                  onView={handleProductView}
                  onEdit={handleProductEdit}
                  onDelete={handleProductDelete}
                  isStaff={isStaff}
                />
              </Grid>
            ))}
          </Grid>

          {/* Pagination */}
          {totalPages > 1 && (
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
              <Pagination
                count={totalPages}
                page={page}
                onChange={handlePageChange}
                color="primary"
                size="large"
              />
            </Box>
          )}
        </>
      )}

      {/* Filter Drawer */}
      <FilterDrawer
        open={filterDrawerOpen}
        onClose={() => setFilterDrawerOpen(false)}
        filters={filters}
        onFiltersChange={handleFiltersChange}
        categories={categories || []}
      />

      {/* Product Detail Dialog */}
      <ProductDetailDialog
        open={productDetailOpen}
        onClose={() => setProductDetailOpen(false)}
        product={selectedProduct}
      />

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Delete Product</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete "{productToDelete?.name}"? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={confirmDelete} 
            color="error" 
            disabled={deleteProductMutation.isLoading}
          >
            {deleteProductMutation.isLoading ? <CircularProgress size={20} /> : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Product Creation/Edit Dialog */}
      <Dialog open={createDialogOpen} onClose={() => setCreateDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingProduct ? 'Edit Product' : 'Create New Product'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 2 }}>
            <TextField
              label="Product Name"
              value={productFormData.name}
              onChange={(e) => setProductFormData(prev => ({ ...prev, name: e.target.value }))}
              fullWidth
              required
            />
            <TextField
              label="Description"
              value={productFormData.description}
              onChange={(e) => setProductFormData(prev => ({ ...prev, description: e.target.value }))}
              fullWidth
              multiline
              rows={3}
            />
            <TextField
              label="SKU"
              value={productFormData.sku}
              onChange={(e) => setProductFormData(prev => ({ ...prev, sku: e.target.value }))}
              fullWidth
              required
            />
            <Autocomplete
              multiple
              options={categories}
              getOptionLabel={(option) => option.name || ''}
              value={categories.filter(cat => productFormData.category_ids.includes(cat.id))}
              onChange={(event, newValue) => {
                setProductFormData(prev => ({
                  ...prev,
                  category_ids: newValue.map(cat => cat.id)
                }));
              }}
              renderTags={(value, getTagProps) =>
                value.map((option, index) => (
                  <Chip
                    variant="outlined"
                    label={option.name}
                    {...getTagProps({ index })}
                    key={option.id}
                  />
                ))
              }
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Categories"
                  placeholder="Select categories..."
                />
              )}
            />
            <Autocomplete
              options={suppliers}
              getOptionLabel={(option) => option.name || ''}
              value={suppliers.find(supplier => supplier.id === productFormData.supplier_id) || null}
              onChange={(event, newValue) => {
                setProductFormData(prev => ({
                  ...prev,
                  supplier_id: newValue ? newValue.id : ''
                }));
              }}
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Supplier"
                  placeholder="Select supplier..."
                  required
                />
              )}
            />
            <Box sx={{ display: 'flex', gap: 2 }}>
              <TextField
                label="Cost Price"
                value={productFormData.cost_price}
                onChange={(e) => setProductFormData(prev => ({ ...prev, cost_price: e.target.value }))}
                type="number"
                inputProps={{ step: '0.01', min: '0' }}
                fullWidth
                required
              />
              <TextField
                label="Selling Price"
                value={productFormData.selling_price}
                onChange={(e) => setProductFormData(prev => ({ ...prev, selling_price: e.target.value }))}
                type="number"
                inputProps={{ step: '0.01', min: '0' }}
                fullWidth
                required
              />
            </Box>
            <FormControl fullWidth>
              <InputLabel>Currency</InputLabel>
              <Select
                value={productFormData.currency}
                label="Currency"
                onChange={(e) => setProductFormData(prev => ({ ...prev, currency: e.target.value }))}
              >
                <MenuItem value="USD">USD</MenuItem>
                <MenuItem value="EUR">EUR</MenuItem>
                <MenuItem value="GBP">GBP</MenuItem>
              </Select>
            </FormControl>
            <FormControlLabel
              control={
                <Switch
                  checked={productFormData.is_active}
                  onChange={(e) => setProductFormData(prev => ({ ...prev, is_active: e.target.checked }))}
                />
              }
              label="Active"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setCreateDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleProductFormSubmit}
            variant="contained"
            disabled={createProductMutation.isLoading || updateProductMutation.isLoading}
          >
            {createProductMutation.isLoading || updateProductMutation.isLoading ? (
              <CircularProgress size={20} />
            ) : (
              editingProduct ? 'Update Product' : 'Create Product'
            )}
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default ProductsPage;
