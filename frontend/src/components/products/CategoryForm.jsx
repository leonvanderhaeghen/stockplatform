import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormControlLabel,
  Checkbox,
  Box,
  CircularProgress,
} from '@mui/material';
import { useSnackbar } from 'notistack';
import categoryService from '../../services/categoryService';

const CategoryForm = ({ open, handleClose, category, onSuccess }) => {
  const { enqueueSnackbar } = useSnackbar();
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState([]);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    parentId: '',
    isActive: true,
  });

  useEffect(() => {
    if (category) {
      setFormData({
        name: category.name || '',
        description: category.description || '',
        parentId: category.parentId || '',
        isActive: category.isActive !== undefined ? category.isActive : true,
      });
    } else {
      setFormData({
        name: '',
        description: '',
        parentId: '',
        isActive: true,
      });
    }
  }, [category, open]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        console.log('Fetching categories...');
        const response = await categoryService.getCategories();
        console.log('Categories API response:', response);
        
        // Handle both direct array response and response with data property
        let categoriesData = [];
        if (Array.isArray(response)) {
          categoriesData = response;
        } else if (response && response.data) {
          categoriesData = Array.isArray(response.data) ? response.data : [response.data];
        } else if (response && response.docs) {
          categoriesData = Array.isArray(response.docs) ? response.docs : [response.docs];
        }
        
        console.log('Processed categories data:', categoriesData);
        
        // Ensure we have an array before filtering
        if (!Array.isArray(categoriesData)) {
          console.error('Categories data is not an array:', categoriesData);
          setCategories([]);
          return;
        }
        
        // Filter out the current category if editing to prevent circular references
        const filteredCategories = category && category._id 
          ? categoriesData.filter(cat => cat && cat._id && cat._id !== category._id)
          : categoriesData.filter(cat => cat && cat._id); // Ensure each category has an _id
          
        console.log('Filtered categories:', filteredCategories);
        setCategories(filteredCategories);
      } catch (error) {
        console.error('Error fetching categories:', error);
        enqueueSnackbar('Failed to load categories', { variant: 'error' });
        setCategories([]);
      }
    };

    if (open) {
      fetchCategories();
    }
  }, [open, category, enqueueSnackbar]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      enqueueSnackbar('Category name is required', { variant: 'error' });
      return;
    }

    setLoading(true);
    try {
      if (category && category._id) {
        await categoryService.updateCategory(category._id, formData);
        enqueueSnackbar('Category updated successfully', { variant: 'success' });
      } else {
        await categoryService.createCategory(formData);
        enqueueSnackbar('Category created successfully', { variant: 'success' });
      }
      onSuccess();
      handleClose();
    } catch (error) {
      console.error('Error saving category:', error);
      enqueueSnackbar(
        error.response?.data?.message || 'Failed to save category', 
        { variant: 'error' }
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {category ? 'Edit Category' : 'Add New Category'}
      </DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Box mb={2}>
            <TextField
              fullWidth
              label="Category Name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              margin="normal"
              required
              disabled={loading}
            />
          </Box>
          
          <Box mb={2}>
            <TextField
              fullWidth
              label="Description"
              name="description"
              value={formData.description}
              onChange={handleChange}
              margin="normal"
              multiline
              rows={3}
              disabled={loading}
            />
          </Box>
          
          <Box mb={2}>
            <FormControl fullWidth margin="normal">
              <InputLabel id="parent-category-label">Parent Category (Optional)</InputLabel>
              <Select
                labelId="parent-category-label"
                name="parentId"
                value={formData.parentId || ''}
                onChange={handleChange}
                label="Parent Category (Optional)"
                disabled={loading}
              >
                <MenuItem value="">
                  <em>None</em>
                </MenuItem>
                {Array.isArray(categories) && categories.length > 0 ? (
                  categories.map((cat) => (
                    cat && cat._id ? (
                      <MenuItem key={cat._id} value={cat._id}>
                        {cat.name || 'Unnamed Category'}
                      </MenuItem>
                    ) : null
                  ))
                ) : (
                  <MenuItem disabled>No categories available</MenuItem>
                )}
              </Select>
            </FormControl>
          </Box>
          
          <Box>
            <FormControlLabel
              control={
                <Checkbox
                  checked={formData.isActive}
                  onChange={handleChange}
                  name="isActive"
                  color="primary"
                  disabled={loading}
                />
              }
              label="Active"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} disabled={loading}>
            Cancel
          </Button>
          <Button 
            type="submit" 
            color="primary" 
            variant="contained" 
            disabled={loading}
            startIcon={loading ? <CircularProgress size={20} /> : null}
          >
            {category ? 'Update' : 'Create'} Category
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default CategoryForm;
