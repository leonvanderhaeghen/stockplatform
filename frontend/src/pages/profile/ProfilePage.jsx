import React, { useState } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Paper,
  Tabs,
  Tab,
  Grid,
  TextField,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
  Chip,
  Divider,
  Alert,
  FormControlLabel,
  Switch,
  CircularProgress,
} from '@mui/material';
import {
  Person,
  Edit,
  Delete,
  Add,
  Home,
  Security,
  Settings,
  Save,
  Visibility,
  VisibilityOff,
  LocationOn,
  Phone,
  Email,
  Badge,
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useSnackbar } from 'notistack';
import { useAuth } from '../../hooks/useAuth';
import userService from '../../services/userService';

const ProfilePage = () => {
  const { user } = useAuth();
  const { enqueueSnackbar } = useSnackbar();
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState(0);
  const [editProfile, setEditProfile] = useState(false);
  const [changePasswordDialog, setChangePasswordDialog] = useState(false);
  const [addressDialog, setAddressDialog] = useState({ open: false, address: null, mode: 'create' });
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  
  const [profileData, setProfileData] = useState({
    firstName: user?.firstName || user?.first_name || '',
    lastName: user?.lastName || user?.last_name || '',
    phone: user?.phone || '',
    email: user?.email || '',
  });
  
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  
  const [newAddress, setNewAddress] = useState({
    name: '',
    street: '',
    city: '',
    state: '',
    postalCode: '',
    country: '',
    phone: '',
    isDefault: false,
  });

  // User Profile Query
  const { data: userProfile, isLoading: profileLoading } = useQuery({
    queryKey: ['user-profile'],
    queryFn: () => userService.getCurrentUser(),
  });

  // User Addresses Query
  const { data: addresses, isLoading: addressesLoading } = useQuery({
    queryKey: ['user-addresses'],
    queryFn: () => userService.getAddresses(),
    enabled: activeTab === 1,
  });

  // Update Profile Mutation
  const updateProfileMutation = useMutation({
    mutationFn: async (profileData) => userService.updateProfile(profileData),
    onSuccess: (updatedData) => {
      queryClient.invalidateQueries({ queryKey: ['user-profile'] });
      // Update local state with the updated data
      if (updatedData) {
        setProfileData({
          firstName: updatedData.firstName || updatedData.first_name || profileData.firstName,
          lastName: updatedData.lastName || updatedData.last_name || profileData.lastName,
          phone: updatedData.phone || profileData.phone,
          email: updatedData.email || profileData.email,
        });
      }
      setEditProfile(false);
      enqueueSnackbar('Profile updated successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.response?.data?.message || 'Failed to update profile', { variant: 'error' });
    },
  });

  // Change Password Mutation
  const changePasswordMutation = useMutation({
    mutationFn: async (passwordData) => userService.changePassword(passwordData.currentPassword, passwordData.newPassword),
    onSuccess: () => {
      setChangePasswordDialog(false);
      setPasswordData({ currentPassword: '', newPassword: '', confirmPassword: '' });
      enqueueSnackbar('Password changed successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.response?.data?.message || 'Failed to change password', { variant: 'error' });
    },
  });

  // Create Address Mutation
  const createAddressMutation = useMutation({
    mutationFn: async (addressData) => userService.createAddress(addressData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-addresses'] });
      setAddressDialog({ open: false, address: null, mode: 'create' });
      setNewAddress({ name: '', street: '', city: '', state: '', postalCode: '', country: '', phone: '', isDefault: false });
      enqueueSnackbar('Address created successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.response?.data?.message || 'Failed to create address', { variant: 'error' });
    },
  });

  // Delete Address Mutation
  const deleteAddressMutation = useMutation({
    mutationFn: async (addressId) => userService.deleteAddress(addressId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-addresses'] });
      enqueueSnackbar('Address deleted successfully', { variant: 'success' });
    },
    onError: (error) => {
      enqueueSnackbar(error.response?.data?.message || 'Failed to delete address', { variant: 'error' });
    },
  });

  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };

  const handleEditProfile = () => {
    setProfileData({
      firstName: userProfile?.data?.first_name || userProfile?.firstName || user?.firstName || '',
      lastName: userProfile?.data?.last_name || userProfile?.lastName || user?.lastName || '',
      phone: userProfile?.data?.phone || userProfile?.phone || user?.phone || '',
      email: userProfile?.data?.email || userProfile?.email || user?.email || '',
    });
    setEditProfile(true);
  };

  const handleSaveProfile = () => {
    if (!profileData.firstName || !profileData.lastName || !profileData.email) {
      enqueueSnackbar('Please fill in all required fields', { variant: 'warning' });
      return;
    }
    updateProfileMutation.mutate(profileData);
  };

  const handleChangePassword = () => {
    if (!passwordData.currentPassword || !passwordData.newPassword || !passwordData.confirmPassword) {
      enqueueSnackbar('Please fill in all password fields', { variant: 'warning' });
      return;
    }
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      enqueueSnackbar('New passwords do not match', { variant: 'warning' });
      return;
    }
    if (passwordData.newPassword.length < 6) {
      enqueueSnackbar('New password must be at least 6 characters long', { variant: 'warning' });
      return;
    }
    changePasswordMutation.mutate({
      currentPassword: passwordData.currentPassword,
      newPassword: passwordData.newPassword,
    });
  };

  const handleCreateAddress = () => {
    setAddressDialog({ open: true, address: null, mode: 'create' });
    setNewAddress({ name: '', street: '', city: '', state: '', postalCode: '', country: '', phone: '', isDefault: false });
  };

  const handleSubmitAddress = () => {
    if (!newAddress.name || !newAddress.street || !newAddress.city || !newAddress.postalCode || !newAddress.country) {
      enqueueSnackbar('Please fill in all required address fields (name, street, city, postal code, country)', { variant: 'warning' });
      return;
    }
    createAddressMutation.mutate(newAddress);
  };

  const handleDeleteAddress = (addressId) => {
    if (window.confirm('Are you sure you want to delete this address?')) {
      deleteAddressMutation.mutate(addressId);
    }
  };

  const renderProfileTab = () => (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h6">Personal Information</Typography>
        {!editProfile && (
          <Button variant="outlined" startIcon={<Edit />} onClick={handleEditProfile}>
            Edit Profile
          </Button>
        )}
      </Box>

      {profileLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                {editProfile ? (
                  <Grid container spacing={2}>
                    <Grid item xs={12} sm={6}>
                      <TextField
                        fullWidth
                        label="First Name"
                        value={profileData.firstName}
                        onChange={(e) => setProfileData({ ...profileData, firstName: e.target.value })}
                        required
                      />
                    </Grid>
                    <Grid item xs={12} sm={6}>
                      <TextField
                        fullWidth
                        label="Last Name"
                        value={profileData.lastName}
                        onChange={(e) => setProfileData({ ...profileData, lastName: e.target.value })}
                        required
                      />
                    </Grid>
                    <Grid item xs={12}>
                      <TextField
                        fullWidth
                        label="Email"
                        type="email"
                        value={profileData.email}
                        onChange={(e) => setProfileData({ ...profileData, email: e.target.value })}
                        required
                        disabled // Usually email changes require special verification
                      />
                    </Grid>
                    <Grid item xs={12}>
                      <TextField
                        fullWidth
                        label="Phone"
                        value={profileData.phone}
                        onChange={(e) => setProfileData({ ...profileData, phone: e.target.value })}
                      />
                    </Grid>
                    <Grid item xs={12}>
                      <Box sx={{ display: 'flex', gap: 2 }}>
                        <Button
                          variant="contained"
                          startIcon={<Save />}
                          onClick={handleSaveProfile}
                          disabled={updateProfileMutation.isPending}
                        >
                          {updateProfileMutation.isPending ? 'Saving...' : 'Save Changes'}
                        </Button>
                        <Button
                          variant="outlined"
                          onClick={() => setEditProfile(false)}
                        >
                          Cancel
                        </Button>
                      </Box>
                    </Grid>
                  </Grid>
                ) : (
                  <Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                      <Person sx={{ mr: 1, color: 'text.secondary' }} />
                      <Typography variant="h6">
                        {userProfile?.data?.first_name || userProfile?.firstName || user?.firstName || user?.first_name} {userProfile?.data?.last_name || userProfile?.lastName || user?.lastName || user?.last_name}
                      </Typography>
                    </Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                      <Email sx={{ mr: 1, color: 'text.secondary' }} />
                      <Typography color="text.secondary">
                        {userProfile?.data?.email || userProfile?.email || user?.email}
                      </Typography>
                    </Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                      <Phone sx={{ mr: 1, color: 'text.secondary' }} />
                      <Typography color="text.secondary">
                        {userProfile?.data?.phone || userProfile?.phone || user?.phone || 'No phone number'}
                      </Typography>
                    </Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                      <Badge sx={{ mr: 1, color: 'text.secondary' }} />
                      <Chip
                        label={user?.role}
                        color={user?.role === 'ADMIN' ? 'error' : user?.role === 'STAFF' ? 'primary' : 'success'}
                        size="small"
                      />
                    </Box>
                  </Box>
                )}
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} md={6}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Security
                </Typography>
                <Button
                  variant="outlined"
                  startIcon={<Security />}
                  onClick={() => setChangePasswordDialog(true)}
                  fullWidth
                  sx={{ mb: 2 }}
                >
                  Change Password
                </Button>
                <Alert severity="info">
                  Keep your account secure by using a strong password and updating it regularly.
                </Alert>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}
    </Box>
  );

  const renderAddressesTab = () => (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h6">Saved Addresses</Typography>
        <Button variant="contained" startIcon={<Add />} onClick={handleCreateAddress}>
          Add Address
        </Button>
      </Box>

      {addressesLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Card>
          {(addresses?.data?.length > 0 || addresses?.addresses?.length > 0) ? (
            <List>
              {(addresses.data || addresses.addresses || []).map((address, index) => (
                <React.Fragment key={address.id}>
                  <ListItem>
                    <LocationOn sx={{ mr: 2, color: 'text.secondary' }} />
                    <ListItemText
                      primary={
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Typography variant="subtitle1">
                            {address.name}: {address.street}, {address.city}
                          </Typography>
                          {(address.isDefault || address.is_default) && (
                            <Chip label="Default" color="primary" size="small" />
                          )}
                        </Box>
                      }
                      secondary={`${address.city}, ${address.state} ${address.postalCode || address.postal_code}, ${address.country}`}
                    />
                    <ListItemSecondaryAction>
                      <IconButton
                        edge="end"
                        onClick={() => handleDeleteAddress(address.id)}
                        color="error"
                      >
                        <Delete />
                      </IconButton>
                    </ListItemSecondaryAction>
                  </ListItem>
                  {index < (addresses.data || addresses.addresses || []).length - 1 && <Divider />}
                </React.Fragment>
              ))}
            </List>
          ) : (
            <CardContent>
              <Typography color="text.secondary" align="center">
                No addresses saved yet. Add your first address to get started.
              </Typography>
            </CardContent>
          )}
        </Card>
      )}
    </Box>
  );

  const renderPreferencesTab = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        Account Preferences
      </Typography>
      
      <Card>
        <CardContent>
          <Typography variant="subtitle1" gutterBottom>
            Notifications
          </Typography>
          <FormControlLabel
            control={<Switch defaultChecked />}
            label="Email notifications for orders"
            sx={{ mb: 1, display: 'block' }}
          />
          <FormControlLabel
            control={<Switch defaultChecked />}
            label="SMS notifications for deliveries"
            sx={{ mb: 1, display: 'block' }}
          />
          <FormControlLabel
            control={<Switch />}
            label="Marketing emails"
            sx={{ mb: 3, display: 'block' }}
          />
          
          <Divider sx={{ my: 2 }} />
          
          <Typography variant="subtitle1" gutterBottom color="error">
            Danger Zone
          </Typography>
          <Button variant="outlined" color="error">
            Deactivate Account
          </Button>
        </CardContent>
      </Card>
    </Box>
  );

  return (
    <Box>
      <Paper sx={{ p: 3, mb: 3, bgcolor: 'primary.main', color: 'primary.contrastText' }}>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Person sx={{ fontSize: 40, mr: 2 }} />
          <Box>
            <Typography variant="h4" gutterBottom>
              Profile
            </Typography>
            <Typography variant="subtitle1">
              Manage your account settings and preferences
            </Typography>
          </Box>
        </Box>
      </Paper>

      <Card>
        <Tabs 
          value={activeTab} 
          onChange={handleTabChange}
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab icon={<Person />} label="Profile" />
          <Tab icon={<Home />} label="Addresses" />
          <Tab icon={<Settings />} label="Preferences" />
        </Tabs>
        
        <Box sx={{ p: 3 }}>
          {activeTab === 0 && renderProfileTab()}
          {activeTab === 1 && renderAddressesTab()}
          {activeTab === 2 && renderPreferencesTab()}
        </Box>
      </Card>

      {/* Change Password Dialog */}
      <Dialog 
        open={changePasswordDialog} 
        onClose={() => setChangePasswordDialog(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Security sx={{ mr: 1 }} />
            Change Password
          </Box>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Current Password"
                type={showCurrentPassword ? 'text' : 'password'}
                value={passwordData.currentPassword}
                onChange={(e) => setPasswordData({ ...passwordData, currentPassword: e.target.value })}
                InputProps={{
                  endAdornment: (
                    <IconButton
                      onClick={() => setShowCurrentPassword(!showCurrentPassword)}
                      edge="end"
                    >
                      {showCurrentPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  ),
                }}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="New Password"
                type={showNewPassword ? 'text' : 'password'}
                value={passwordData.newPassword}
                onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
                InputProps={{
                  endAdornment: (
                    <IconButton
                      onClick={() => setShowNewPassword(!showNewPassword)}
                      edge="end"
                    >
                      {showNewPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  ),
                }}
                required
                helperText="Must be at least 6 characters long"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Confirm New Password"
                type={showConfirmPassword ? 'text' : 'password'}
                value={passwordData.confirmPassword}
                onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
                InputProps={{
                  endAdornment: (
                    <IconButton
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      edge="end"
                    >
                      {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  ),
                }}
                required
                error={passwordData.confirmPassword && passwordData.newPassword !== passwordData.confirmPassword}
                helperText={passwordData.confirmPassword && passwordData.newPassword !== passwordData.confirmPassword ? 'Passwords do not match' : ''}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setChangePasswordDialog(false)}>
            Cancel
          </Button>
          <Button 
            variant="contained" 
            onClick={handleChangePassword}
            disabled={changePasswordMutation.isPending}
          >
            {changePasswordMutation.isPending ? 'Changing...' : 'Change Password'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Address Dialog */}
      <Dialog 
        open={addressDialog.open} 
        onClose={() => setAddressDialog({ open: false, address: null, mode: 'create' })}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Add sx={{ mr: 1 }} />
            Add New Address
          </Box>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={2} sx={{ mt: 1 }}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Address Name"
                value={newAddress.name}
                onChange={(e) => setNewAddress({ ...newAddress, name: e.target.value })}
                placeholder="e.g., Home, Work, Billing"
                helperText="Give this address a name for easy identification"
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Street Address"
                value={newAddress.street}
                onChange={(e) => setNewAddress({ ...newAddress, street: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="City"
                value={newAddress.city}
                onChange={(e) => setNewAddress({ ...newAddress, city: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="State/Province"
                value={newAddress.state}
                onChange={(e) => setNewAddress({ ...newAddress, state: e.target.value })}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Postal Code"
                value={newAddress.postalCode}
                onChange={(e) => setNewAddress({ ...newAddress, postalCode: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Country"
                value={newAddress.country}
                onChange={(e) => setNewAddress({ ...newAddress, country: e.target.value })}
                required
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Phone Number"
                value={newAddress.phone}
                onChange={(e) => setNewAddress({ ...newAddress, phone: e.target.value })}
                placeholder="Optional"
              />
            </Grid>
            <Grid item xs={12}>
              <FormControlLabel
                control={
                  <Switch
                    checked={newAddress.isDefault}
                    onChange={(e) => setNewAddress({ ...newAddress, isDefault: e.target.checked })}
                  />
                }
                label="Set as default address"
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button 
            onClick={() => setAddressDialog({ open: false, address: null, mode: 'create' })}
          >
            Cancel
          </Button>
          <Button 
            variant="contained" 
            onClick={handleSubmitAddress}
            disabled={createAddressMutation.isPending}
          >
            {createAddressMutation.isPending ? 'Adding...' : 'Add Address'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default ProfilePage;
