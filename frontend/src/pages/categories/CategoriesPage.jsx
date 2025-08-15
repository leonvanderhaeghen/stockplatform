import React, { useState, useCallback } from 'react';
import {
  Container,
  Typography,
  Box,
  Card,
  CardContent,
  Button,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
  Alert,
  CircularProgress,
  Chip,
  Grid,
  Paper,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Collapse,
  Breadcrumbs,
  Link,
  Divider,
  Fab,
  Tooltip
} from '@mui/material';
import {
  Add,
  Edit,
  Delete,
  ExpandMore,
  ExpandLess,
  Category,
  Visibility,
  VisibilityOff,
  FolderOpen,
  Folder,
  NavigateNext
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSnackbar } from 'notistack';
import categoryService, { getCategoryBreadcrumb } from '../../services/categoryService';
import { useAuth } from '../../hooks/useAuth';

const CategoriesPage = () => {
  const { user } = useAuth();
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();

  // UI State
  const [openDialog, setOpenDialog] = useState(false);
  const [editingCategory, setEditingCategory] = useState(null);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  const [categoryToDelete, setCategoryToDelete] = useState(null);
  const [expandedCategories, setExpandedCategories] = useState(new Set());
  const [selectedCategory, setSelectedCategory] = useState(null);

  // Form State
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    parent_id: '',
    is_active: true
  });

  // Fetch categories
  const {
    data: categories = [],
    isLoading,
    error
  } = useQuery({
    queryKey: ['categories'],
    queryFn: categoryService.getCategories
  });

  // Fetch category tree
  const {
    data: categoryTree = [],
    isLoading: isTreeLoading
  } = useQuery({
    queryKey: ['categoryTree'],
    queryFn: categoryService.getCategoryTree,
    enabled: categories.length > 0
  });

  // Create category mutation
  const createMutation = useMutation({
    mutationFn: categoryService.createCategory,
    onSuccess: () => {
      queryClient.invalidateQueries(['categories']);
      queryClient.invalidateQueries(['categoryTree']);
      handleCloseDialog();
      enqueueSnackbar('Category created successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(
        error.response?.data?.message || 'Failed to create category',
        { variant: 'error' }
      );
    }
  });

  // Update category mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, data }) => categoryService.updateCategory(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries(['categories']);
      queryClient.invalidateQueries(['categoryTree']);
      handleCloseDialog();
      enqueueSnackbar('Category updated successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(
        error.response?.data?.message || 'Failed to update category',
        { variant: 'error' }
      );
    }
  });

  // Delete category mutation
  const deleteMutation = useMutation({
    mutationFn: categoryService.deleteCategory,
    onSuccess: () => {
      queryClient.invalidateQueries(['categories']);
      queryClient.invalidateQueries(['categoryTree']);
      setDeleteConfirmOpen(false);
      setCategoryToDelete(null);
      enqueueSnackbar('Category deleted successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(
        error.response?.data?.message || 'Failed to delete category',
        { variant: 'error' }
      );
    }
  });

  // Event Handlers
  const handleOpenCreateDialog = useCallback(() => {
    setEditingCategory(null);
    setFormData({
      name: '',
      description: '',
      parent_id: '',
      is_active: true
    });
    setOpenDialog(true);
  }, []);

  const handleOpenEditDialog = useCallback((category) => {
    setEditingCategory(category);
    setFormData({
      name: category.name || '',
      description: category.description || '',
      parent_id: category.parent_id || '',
      is_active: category.is_active !== false
    });
    setOpenDialog(true);
  }, []);

  const handleCloseDialog = useCallback(() => {
    setOpenDialog(false);
    setEditingCategory(null);
    setFormData({
      name: '',
      description: '',
      parent_id: '',
      is_active: true
    });
  }, []);

  const handleInputChange = useCallback((field, value) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  }, []);

  const handleSubmit = useCallback(() => {
    if (!formData.name.trim()) {
      enqueueSnackbar('Category name is required', { variant: 'error' });
      return;
    }

    const submitData = {
      name: formData.name.trim(),
      description: formData.description.trim(),
      is_active: formData.is_active
    };

    // Only include parent_id if it's not empty
    if (formData.parent_id) {
      submitData.parent_id = formData.parent_id;
    }

    if (editingCategory) {
      updateMutation.mutate({ id: editingCategory.id, data: submitData });
    } else {
      createMutation.mutate(submitData);
    }
  }, [formData, editingCategory, createMutation, updateMutation, enqueueSnackbar]);

  const handleDeleteClick = useCallback((category) => {
    setCategoryToDelete(category);
    setDeleteConfirmOpen(true);
  }, []);

  const handleDeleteConfirm = useCallback(() => {
    if (categoryToDelete) {
      deleteMutation.mutate(categoryToDelete.id);
    }
  }, [categoryToDelete, deleteMutation]);

  const toggleExpanded = useCallback((categoryId) => {
    setExpandedCategories(prev => {
      const newSet = new Set(prev);
      if (newSet.has(categoryId)) {
        newSet.delete(categoryId);
      } else {
        newSet.add(categoryId);
      }
      return newSet;
    });
  }, []);

  // Recursive Category Tree Component
  const CategoryTreeItem = ({ category, level = 0 }) => {
    const hasChildren = category.children && category.children.length > 0;
    const isExpanded = expandedCategories.has(category.id);
    const breadcrumb = getCategoryBreadcrumb(category.id, categories);

    return (
      <Box key={category.id}>
        <ListItem
          sx={{
            pl: 2 + level * 2,
            backgroundColor: selectedCategory?.id === category.id ? 'action.selected' : 'transparent',
            '&:hover': { backgroundColor: 'action.hover' }
          }}
          onClick={() => setSelectedCategory(category)}
        >
          <Box sx={{ display: 'flex', alignItems: 'center', flex: 1 }}>
            {hasChildren ? (
              <IconButton
                size="small"
                onClick={(e) => {
                  e.stopPropagation();
                  toggleExpanded(category.id);
                }}
              >
                {isExpanded ? <ExpandLess /> : <ExpandMore />}
              </IconButton>
            ) : (
              <Box sx={{ width: 40 }} />
            )}
            
            <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
              {hasChildren ? <FolderOpen /> : <Folder />}
            </Box>

            <ListItemText
              primary={category.name}
              secondary={category.description}
              sx={{ flex: 1 }}
            />

            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Chip
                label={category.is_active ? 'Active' : 'Inactive'}
                color={category.is_active ? 'success' : 'default'}
                size="small"
              />
              
              <Tooltip title="Edit Category">
                <IconButton
                  size="small"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleOpenEditDialog(category);
                  }}
                >
                  <Edit />
                </IconButton>
              </Tooltip>

              <Tooltip title="Delete Category">
                <IconButton
                  size="small"
                  color="error"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDeleteClick(category);
                  }}
                >
                  <Delete />
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
        </ListItem>

        {hasChildren && isExpanded && (
          <Collapse in={isExpanded}>
            {category.children.map(child => (
              <CategoryTreeItem
                key={child.id}
                category={child}
                level={level + 1}
              />
            ))}
          </Collapse>
        )}
      </Box>
    );
  };

  // Get available parent categories (excluding current category when editing)
  const getAvailableParents = useCallback(() => {
    if (!categories) return [];
    
    return categories.filter(cat => {
      if (editingCategory && cat.id === editingCategory.id) {
        return false; // Can't be parent of itself
      }
      return cat.is_active;
    });
  }, [categories, editingCategory]);

  if (isLoading) {
    return (
      <Container>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (error) {
    return (
      <Container>
        <Alert severity="error">
          Failed to load categories: {error.message}
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg">
      <Box py={3}>
        {/* Header */}
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Typography variant="h4" component="h1">
            Categories Management
          </Typography>
          <Button
            variant="contained"
            startIcon={<Add />}
            onClick={handleOpenCreateDialog}
          >
            Create Category
          </Button>
        </Box>

        {/* Main Content */}
        <Grid container spacing={3}>
          {/* Category Tree */}
          <Grid item xs={12} md={8}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Category Hierarchy
                </Typography>
                {isTreeLoading ? (
                  <Box display="flex" justifyContent="center" p={3}>
                    <CircularProgress />
                  </Box>
                ) : categoryTree.length === 0 ? (
                  <Box textAlign="center" p={3}>
                    <Category sx={{ fontSize: 48, color: 'text.secondary', mb: 2 }} />
                    <Typography variant="body1" color="text.secondary">
                      No categories found. Create your first category to get started.
                    </Typography>
                  </Box>
                ) : (
                  <List>
                    {categoryTree.map(category => (
                      <CategoryTreeItem key={category.id} category={category} />
                    ))}
                  </List>
                )}
              </CardContent>
            </Card>
          </Grid>

          {/* Category Details */}
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Category Details
                </Typography>
                
                {selectedCategory ? (
                  <Box>
                    <Typography variant="h6" gutterBottom>
                      {selectedCategory.name}
                    </Typography>
                    
                    {selectedCategory.description && (
                      <Typography variant="body2" color="text.secondary" paragraph>
                        {selectedCategory.description}
                      </Typography>
                    )}

                    <Box mb={2}>
                      <Chip
                        label={selectedCategory.is_active ? 'Active' : 'Inactive'}
                        color={selectedCategory.is_active ? 'success' : 'default'}
                        size="small"
                      />
                    </Box>

                    {/* Breadcrumb */}
                    {selectedCategory.parent_id && (
                      <Box mb={2}>
                        <Typography variant="subtitle2" gutterBottom>
                          Path:
                        </Typography>
                        <Breadcrumbs separator={<NavigateNext fontSize="small" />}>
                          {getCategoryBreadcrumb(selectedCategory.id, categories).map((cat, index, arr) => (
                            <Typography
                              key={cat.id}
                              color={index === arr.length - 1 ? 'text.primary' : 'text.secondary'}
                              variant="body2"
                            >
                              {cat.name}
                            </Typography>
                          ))}
                        </Breadcrumbs>
                      </Box>
                    )}

                    <Box display="flex" gap={1} mt={2}>
                      <Button
                        variant="outlined"
                        size="small"
                        startIcon={<Edit />}
                        onClick={() => handleOpenEditDialog(selectedCategory)}
                      >
                        Edit
                      </Button>
                      <Button
                        variant="outlined"
                        color="error"
                        size="small"
                        startIcon={<Delete />}
                        onClick={() => handleDeleteClick(selectedCategory)}
                      >
                        Delete
                      </Button>
                    </Box>
                  </Box>
                ) : (
                  <Typography variant="body2" color="text.secondary">
                    Select a category from the tree to view its details.
                  </Typography>
                )}
              </CardContent>
            </Card>
          </Grid>
        </Grid>

        {/* Create/Edit Dialog */}
        <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="sm" fullWidth>
          <DialogTitle>
            {editingCategory ? 'Edit Category' : 'Create New Category'}
          </DialogTitle>
          <DialogContent>
            <Box display="flex" flexDirection="column" gap={2} mt={1}>
              <TextField
                label="Category Name"
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                required
                fullWidth
              />

              <TextField
                label="Description"
                value={formData.description}
                onChange={(e) => handleInputChange('description', e.target.value)}
                multiline
                rows={3}
                fullWidth
              />

              <FormControl fullWidth>
                <InputLabel>Parent Category</InputLabel>
                <Select
                  value={formData.parent_id}
                  onChange={(e) => handleInputChange('parent_id', e.target.value)}
                  label="Parent Category"
                >
                  <MenuItem value="">
                    <em>None (Root Category)</em>
                  </MenuItem>
                  {getAvailableParents().map(category => (
                    <MenuItem key={category.id} value={category.id}>
                      {category.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>

              <FormControlLabel
                control={
                  <Switch
                    checked={formData.is_active}
                    onChange={(e) => handleInputChange('is_active', e.target.checked)}
                  />
                }
                label="Active"
              />
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseDialog}>Cancel</Button>
            <Button
              onClick={handleSubmit}
              variant="contained"
              disabled={createMutation.isPending || updateMutation.isPending}
            >
              {createMutation.isPending || updateMutation.isPending ? (
                <CircularProgress size={20} />
              ) : (
                editingCategory ? 'Update' : 'Create'
              )}
            </Button>
          </DialogActions>
        </Dialog>

        {/* Delete Confirmation Dialog */}
        <Dialog open={deleteConfirmOpen} onClose={() => setDeleteConfirmOpen(false)}>
          <DialogTitle>Confirm Delete</DialogTitle>
          <DialogContent>
            <Typography>
              Are you sure you want to delete the category "{categoryToDelete?.name}"?
              This action cannot be undone.
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDeleteConfirmOpen(false)}>Cancel</Button>
            <Button
              onClick={handleDeleteConfirm}
              color="error"
              variant="contained"
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? <CircularProgress size={20} /> : 'Delete'}
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </Container>
  );
};

export default CategoriesPage;
