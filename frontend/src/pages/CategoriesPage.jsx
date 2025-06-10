import React, { useState, useEffect, useCallback } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Button, 
  IconButton, 
  Chip, 
  Tooltip,
  CircularProgress,
  TextField,
  InputAdornment,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from '@mui/material';
import { 
  Add as AddIcon, 
  Edit as EditIcon, 
  Delete as DeleteIcon, 
  Search as SearchIcon,
  Category as CategoryIcon,
} from '@mui/icons-material';
import { useSnackbar } from 'notistack';
import { format } from 'date-fns';
import { DataGrid } from '@mui/x-data-grid';
import CategoryForm from '../components/products/CategoryForm';
import categoryService from '../services/categoryService';

const CategoriesPage = () => {
  const { enqueueSnackbar } = useSnackbar();
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [openForm, setOpenForm] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [categoryToDelete, setCategoryToDelete] = useState(null);
  const [paginationModel, setPaginationModel] = useState({
    page: 0,
    pageSize: 10,
  });
  const [rowCount, setRowCount] = useState(0);

  const fetchCategories = useCallback(async () => {
    try {
      setLoading(true);
      const params = {
        q: searchQuery,
        page: paginationModel.page + 1,
        limit: paginationModel.pageSize,
      };
      
      const response = await categoryService.getCategories(params);
      
      // Handle both array and paginated response formats
      if (Array.isArray(response)) {
        setCategories(response);
        setRowCount(response.length);
      } else {
        setCategories(response.data || response.docs || []);
        setRowCount(response.total || response.docs?.length || 0);
      }
    } catch (error) {
      console.error('Error fetching categories:', error);
      enqueueSnackbar('Failed to load categories', { variant: 'error' });
      setCategories([]);
    } finally {
      setLoading(false);
    }
  }, [searchQuery, paginationModel.page, paginationModel.pageSize, enqueueSnackbar]);

  useEffect(() => {
    const timer = setTimeout(() => {
      fetchCategories();
    }, 300);

    return () => clearTimeout(timer);
  }, [searchQuery, paginationModel, fetchCategories]);

  const handleSearch = (e) => {
    setSearchQuery(e.target.value);
    setPaginationModel(prev => ({ ...prev, page: 0 }));
  };

  const handleAddCategory = () => {
    setSelectedCategory(null);
    setOpenForm(true);
  };

  const handleEditCategory = (category) => {
    setSelectedCategory(category);
    setOpenForm(true);
  };

  const handleDeleteClick = (category) => {
    setCategoryToDelete(category);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!categoryToDelete) return;
    
    try {
      await categoryService.deleteCategory(categoryToDelete._id);
      enqueueSnackbar('Category deleted successfully', { variant: 'success' });
      fetchCategories();
    } catch (error) {
      console.error('Error deleting category:', error);
      enqueueSnackbar(
        error.response?.data?.message || 'Failed to delete category', 
        { variant: 'error' }
      );
    } finally {
      setDeleteDialogOpen(false);
      setCategoryToDelete(null);
    }
  };

  const handleFormSuccess = () => {
    fetchCategories();
  };

  const columns = [
    { 
      field: 'name', 
      headerName: 'Name', 
      flex: 2,
      renderCell: (params) => (
        <Box display="flex" alignItems="center">
          <CategoryIcon color="action" sx={{ mr: 1 }} />
          <Typography variant="body1">{params.value}</Typography>
        </Box>
      )
    },
    { 
      field: 'description', 
      headerName: 'Description', 
      flex: 3,
      renderCell: (params) => (
        <Typography variant="body2" color="textSecondary" noWrap>
          {params.value || 'No description'}
        </Typography>
      )
    },
    { 
      field: 'parent', 
      headerName: 'Parent Category', 
      flex: 2,
      valueGetter: (params) => {
        if (!params.row.parentId) return '—';
        const parent = categories.find(cat => cat._id === params.row.parentId);
        return parent ? parent.name : 'Unknown';
      },
      renderCell: (params) => (
        <Chip 
          label={params.value} 
          size="small" 
          variant="outlined"
          sx={{ 
            backgroundColor: params.value !== '—' ? 'action.hover' : 'transparent',
            borderColor: 'divider',
            color: 'text.secondary'
          }}
        />
      )
    },
    { 
      field: 'productCount', 
      headerName: 'Products', 
      flex: 1,
      align: 'center',
      headerAlign: 'center',
      renderCell: (params) => (
        <Chip 
          label={params.value || 0} 
          size="small"
          variant="outlined"
          color="primary"
        />
      )
    },
    { 
      field: 'status', 
      headerName: 'Status', 
      flex: 1,
      renderCell: (params) => (
        <Chip 
          label={params.row.isActive ? 'Active' : 'Inactive'} 
          color={params.row.isActive ? 'success' : 'default'}
          size="small"
          variant="outlined"
        />
      )
    },
    { 
      field: 'createdAt', 
      headerName: 'Created', 
      flex: 1.5,
      valueFormatter: (params) => 
        params.value ? format(new Date(params.value), 'PPpp') : '—',
    },
    {
      field: 'actions',
      headerName: 'Actions',
      sortable: false,
      flex: 1,
      renderCell: (params) => (
        <Box>
          <Tooltip title="Edit">
            <IconButton 
              size="small" 
              onClick={(e) => {
                e.stopPropagation();
                handleEditCategory(params.row);
              }}
              color="primary"
            >
              <EditIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete">
            <IconButton 
              size="small" 
              onClick={(e) => {
                e.stopPropagation();
                handleDeleteClick(params.row);
              }}
              color="error"
            >
              <DeleteIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        </Box>
      ),
    },
  ];

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1">Product Categories</Typography>
        <Button 
          variant="contained" 
          color="primary" 
          startIcon={<AddIcon />}
          onClick={handleAddCategory}
        >
          Add Category
        </Button>
      </Box>
      
      <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Search categories..."
          value={searchQuery}
          onChange={handleSearch}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon color="action" />
              </InputAdornment>
            ),
            sx: { backgroundColor: 'background.paper' }
          }}
        />
      </Paper>
      
      <Paper elevation={3} sx={{ height: 600, width: '100%' }}>
        <DataGrid
          rows={categories}
          columns={columns}
          loading={loading}
          getRowId={(row) => row._id}
          pageSizeOptions={[5, 10, 25, 50]}
          paginationMode="server"
          paginationModel={paginationModel}
          onPaginationModelChange={setPaginationModel}
          rowCount={rowCount}
          disableRowSelectionOnClick
          sx={{
            '& .MuiDataGrid-cell:focus': {
              outline: 'none',
            },
            '& .MuiDataRow-root:hover': {
              backgroundColor: 'action.hover',
              cursor: 'pointer',
            },
          }}
        />
      </Paper>

      {/* Category Form Dialog */}
      <CategoryForm
        open={openForm}
        handleClose={() => setOpenForm(false)}
        category={selectedCategory}
        onSuccess={handleFormSuccess}
      />

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => setDeleteDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Delete Category</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete the category "{categoryToDelete?.name}"? 
            This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button 
            onClick={() => setDeleteDialogOpen(false)}
            color="primary"
          >
            Cancel
          </Button>
          <Button 
            onClick={handleDeleteConfirm}
            color="error"
            variant="contained"
            startIcon={loading ? <CircularProgress size={20} /> : null}
            disabled={loading}
          >
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default CategoriesPage;
