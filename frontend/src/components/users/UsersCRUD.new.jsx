import React, { useState, useEffect } from 'react';
import {
  Container,
  Typography,
  Alert,
  Button,
  CircularProgress,
  Box,
  Card,
  CardHeader,
  CardContent,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
  Chip,
  Snackbar,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Menu,
  MenuItem,
  ListItemIcon as MuiListItemIcon,
} from '@mui/material';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Add as AddIcon,
  MoreVert as MoreVertIcon,
  Star as StarIcon,
  Home as HomeIcon,
  Work as WorkIcon,
  LocationOn as LocationIcon,
  Phone as PhoneIcon,
  Email as EmailIcon,
  Person as PersonIcon,
} from '@mui/icons-material';
import userService from '../../services/userService';
import { formatPhoneNumber } from '../../utils/formatPhoneNumber';
import ProfileForm from './ProfileForm';
import AddressForm from './AddressForm';

// Styled components
const StyledListItemIcon = ({ children }) => (
  <MuiListItemIcon sx={{ minWidth: 36, color: 'inherit' }}>
    {children}
  </MuiListItemIcon>
);

// Export ListItemIcon for consistency
export const ListItemIcon = StyledListItemIcon;

// Main UsersCRUD component for managing user profile and addresses
const UsersCRUD = () => {
  // User data state
  const [user, setUser] = useState(null);
  const [addresses, setAddresses] = useState([]);
  
  // UI state
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  
  // Dialog states
  const [editProfileOpen, setEditProfileOpen] = useState(false);
  const [addressDialogOpen, setAddressDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  
  // Selected items
  const [selectedAddress, setSelectedAddress] = useState(null);
  
  // Loading states
  const [profileUpdating, setProfileUpdating] = useState(false);
  const [addressUpdating, setAddressUpdating] = useState(false);
  const [addressDeleting, setAddressDeleting] = useState(false);
  const [settingDefault, setSettingDefault] = useState(false);

  // Fetch user data and addresses
  const fetchUserData = async () => {
    try {
      setLoading(true);
      setError('');
      
      const [userRes, addressesRes] = await Promise.all([
        userService.getProfile(),
        userService.getAddresses()
      ]);
      
      setUser(userRes.data);
      setAddresses(addressesRes.data || []);
    } catch (err) {
      setError('Failed to load user data. Please try again.');
      console.error('Error fetching user data:', err);
    } finally {
      setLoading(false);
    }
  };

  // Load data on component mount
  useEffect(() => {
    fetchUserData();
  }, []);

  // Profile handlers
  const handleOpenEditProfile = () => setEditProfileOpen(true);
  const handleCloseEditProfile = () => setEditProfileOpen(false);
  
  const handleUpdateProfile = async (profileData) => {
    try {
      setProfileUpdating(true);
      setError('');
      
      const { data } = await userService.updateProfile(profileData);
      setUser(data);
      setSuccess('Profile updated successfully');
      setSnackbarOpen(true);
      handleCloseEditProfile();
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to update profile');
      console.error('Error updating profile:', err);
    } finally {
      setProfileUpdating(false);
    }
  };

  // Address dialog handlers
  const handleOpenAddressDialog = (address = null) => {
    setSelectedAddress(address);
    setAddressDialogOpen(true);
  };

  const handleCloseAddressDialog = () => {
    setAddressDialogOpen(false);
    setSelectedAddress(null);
  };

  // Address CRUD operations
  const handleSaveAddress = async (addressData) => {
    try {
      setAddressUpdating(true);
      setError('');
      
      if (selectedAddress?.id) {
        // Update existing address
        const { data } = await userService.updateAddress(selectedAddress.id, addressData);
        setAddresses(prev => prev.map(addr => 
          addr.id === selectedAddress.id ? data : addr
        ));
        setSuccess('Address updated successfully');
      } else {
        // Add new address
        const { data } = await userService.createAddress(addressData);
        setAddresses(prev => [...prev, data]);
        setSuccess('Address added successfully');
      }
      
      setSnackbarOpen(true);
      handleCloseAddressDialog();
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to save address');
      console.error('Error saving address:', err);
    } finally {
      setAddressUpdating(false);
    }
  };

  // Delete address
  const handleDeleteAddress = async () => {
    if (!selectedAddress?.id) return;
    
    try {
      setAddressDeleting(true);
      setError('');
      
      await userService.deleteAddress(selectedAddress.id);
      setAddresses(prev => prev.filter(addr => addr.id !== selectedAddress.id));
      setSuccess('Address deleted successfully');
      setSnackbarOpen(true);
      setDeleteDialogOpen(false);
      setSelectedAddress(null);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to delete address');
      console.error('Error deleting address:', err);
    } finally {
      setAddressDeleting(false);
    }
  };

  // Set default address
  const handleSetDefaultAddress = async (addressId) => {
    try {
      setSettingDefault(true);
      setError('');
      
      await userService.setDefaultAddress(addressId);
      
      // Update addresses list with new default
      setAddresses(prev => 
        prev.map(addr => ({
          ...addr,
          isDefault: addr.id === addressId
        }))
      );
      
      setSuccess('Default address updated successfully');
      setSnackbarOpen(true);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to set default address');
      console.error('Error setting default address:', err);
    } finally {
      setSettingDefault(false);
    }
  };

  // Snackbar handler
  const handleCloseSnackbar = () => {
    setSnackbarOpen(false);
  };

  // Loading state
  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  // Error state
  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  // No user data
  if (!user) {
    return (
      <Alert severity="info">
        No user data available. Please try again later.
      </Alert>
    );
  }

  return (
    <Container maxWidth="lg">
      <Typography variant="h4" gutterBottom>
        My Profile
      </Typography>
      
      {/* Success Snackbar */}
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleCloseSnackbar}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert 
          onClose={handleCloseSnackbar} 
          severity="success" 
          sx={{ width: '100%' }}
        >
          {success}
        </Alert>
      </Snackbar>
      
      {/* Profile Section */}
      <Card sx={{ mb: 4 }}>
        <CardHeader
          title="Profile Information"
          action={
            <Button
              startIcon={<EditIcon />}
              onClick={handleOpenEditProfile}
              disabled={profileUpdating}
              variant="outlined"
            >
              Edit Profile
            </Button>
          }
        />
        <CardContent>
          <Box display="flex" alignItems="center" mb={2}>
            <PersonIcon color="action" sx={{ mr: 2 }} />
            <Typography variant="body1">
              {user.firstName} {user.lastName}
            </Typography>
          </Box>
          <Box display="flex" alignItems="center" mb={2}>
            <EmailIcon color="action" sx={{ mr: 2 }} />
            <Typography variant="body1">
              {user.email}
            </Typography>
          </Box>
          {user.phone && (
            <Box display="flex" alignItems="center">
              <PhoneIcon color="action" sx={{ mr: 2 }} />
              <Typography variant="body1">
                {formatPhoneNumber(user.phone)}
              </Typography>
            </Box>
          )}
        </CardContent>
      </Card>
      
      {/* Addresses Section */}
      <Card>
        <CardHeader
          title="Shipping Addresses"
          subheader={`${addresses.length} saved address${addresses.length !== 1 ? 'es' : ''}`}
          action={
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => handleOpenAddressDialog()}
              disabled={addressUpdating}
            >
              Add Address
            </Button>
          }
        />
        <CardContent>
          {addresses.length === 0 ? (
            <Box textAlign="center" py={4}>
              <LocationIcon color="disabled" sx={{ fontSize: 48, mb: 2 }} />
              <Typography variant="h6" color="textSecondary" gutterBottom>
                No saved addresses
              </Typography>
              <Typography color="textSecondary" paragraph>
                Add your first address to get started
              </Typography>
              <Button
                variant="outlined"
                startIcon={<AddIcon />}
                onClick={() => handleOpenAddressDialog()}
              >
                Add Address
              </Button>
            </Box>
          ) : (
            <List>
              {addresses.map((address) => (
                <React.Fragment key={address.id}>
                  <ListItem 
                    alignItems="flex-start"
                    sx={{
                      border: '1px solid',
                      borderColor: 'divider',
                      borderRadius: 1,
                      mb: 2,
                      position: 'relative',
                      bgcolor: address.isDefault ? 'action.hover' : 'background.paper',
                    }}
                  >
                    {address.isDefault && (
                      <Chip
                        icon={<StarIcon />}
                        label="Default"
                        color="primary"
                        size="small"
                        sx={{
                          position: 'absolute',
                          top: 8,
                          right: 8,
                        }}
                      />
                    )}
                    <ListItemIcon sx={{ minWidth: 40, mt: 1 }}>
                      {address.name?.toLowerCase().includes('work') ? (
                        <WorkIcon color="primary" />
                      ) : address.name?.toLowerCase().includes('home') ? (
                        <HomeIcon color="primary" />
                      ) : (
                        <LocationIcon color="primary" />
                      )}
                    </ListItemIcon>
                    <ListItemText
                      primary={
                        <Box display="flex" alignItems="center">
                          <Typography variant="subtitle1" component="span" fontWeight="bold">
                            {address.name}
                          </Typography>
                        </Box>
                      }
                      secondary={
                        <Box component="span" sx={{ display: 'block', mt: 0.5 }}>
                          <Typography variant="body2" color="text.primary" component="span" display="block">
                            {address.street}
                          </Typography>
                          <Typography variant="body2" color="text.secondary" component="span" display="block">
                            {[address.city, address.state, address.postalCode].filter(Boolean).join(', ')}
                          </Typography>
                          <Typography variant="body2" color="text.secondary" component="span" display="block">
                            {address.country}
                          </Typography>
                          <Typography variant="body2" color="text.secondary" component="span" display="flex" alignItems="center" mt={0.5}>
                            <PhoneIcon fontSize="small" sx={{ mr: 0.5, fontSize: '1rem' }} />
                            {formatPhoneNumber(address.phone) || 'N/A'}
                          </Typography>
                        </Box>
                      }
                    />
                    <ListItemSecondaryAction>
                      <IconButton 
                        edge="end" 
                        aria-label="more"
                        onClick={(e) => {
                          setSelectedAddress(address);
                        }}
                      >
                        <MoreVertIcon />
                      </IconButton>
                    </ListItemSecondaryAction>
                  </ListItem>
                </React.Fragment>
              ))}
            </List>
          )}
        </CardContent>
      </Card>
      
      {/* Edit Profile Dialog */}
      <Dialog 
        open={editProfileOpen} 
        onClose={handleCloseEditProfile}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Edit Profile</DialogTitle>
        <ProfileForm 
          user={user}
          onUpdate={handleUpdateProfile}
          onCancel={handleCloseEditProfile}
          loading={profileUpdating}
        />
      </Dialog>
      
      {/* Add/Edit Address Dialog */}
      <Dialog 
        open={addressDialogOpen} 
        onClose={handleCloseAddressDialog}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          {selectedAddress ? 'Edit Address' : 'Add New Address'}
        </DialogTitle>
        <AddressForm 
          address={selectedAddress || {}}
          onSubmit={handleSaveAddress}
          onCancel={handleCloseAddressDialog}
          loading={addressUpdating}
        />
      </Dialog>
      
      {/* Delete Confirmation Dialog */}
      <Dialog 
        open={deleteDialogOpen} 
        onClose={() => setDeleteDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Delete Address</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete this address? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button 
            onClick={() => setDeleteDialogOpen(false)}
            disabled={addressDeleting}
          >
            Cancel
          </Button>
          <Button 
            onClick={handleDeleteAddress}
            color="error"
            variant="contained"
            disabled={addressDeleting}
            startIcon={addressDeleting ? <CircularProgress size={20} /> : <DeleteIcon />}
          >
            {addressDeleting ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* Address Actions Menu */}
      <Menu
        anchorEl={document.querySelector('[aria-label="more"]')}
        open={!!selectedAddress}
        onClose={() => setSelectedAddress(null)}
      >
        <MenuItem 
          onClick={() => {
            handleOpenAddressDialog(selectedAddress);
            setSelectedAddress(null);
          }}
        >
          <ListItemIcon>
            <EditIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText>Edit</ListItemText>
        </MenuItem>
        {!selectedAddress?.isDefault && (
          <MenuItem 
            onClick={() => {
              handleSetDefaultAddress(selectedAddress.id);
              setSelectedAddress(null);
            }}
            disabled={settingDefault}
          >
            <ListItemIcon>
              {settingDefault ? (
                <CircularProgress size={20} />
              ) : (
                <StarIcon fontSize="small" />
              )}
            </ListItemIcon>
            <ListItemText>
              {settingDefault ? 'Setting as default...' : 'Set as default'}
            </ListItemText>
          </MenuItem>
        )}
        <MenuItem 
          onClick={() => {
            setDeleteDialogOpen(true);
          }}
          disabled={addressDeleting}
        >
          <ListItemIcon>
            <DeleteIcon fontSize="small" color="error" />
          </ListItemIcon>
          <ListItemText primaryTypographyProps={{ color: 'error.main' }}>
            Delete
          </ListItemText>
        </MenuItem>
      </Menu>
    </Container>
  );
};

export default UsersCRUD;
