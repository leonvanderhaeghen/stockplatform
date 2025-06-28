import React, { useContext, useState } from 'react';
import { useNavigate, Navigate, Link as RouterLink } from 'react-router-dom';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import {
  Box,
  Button,
  TextField,
  Typography,
  Paper,
  Container,
  Link,
  Alert,
  InputAdornment,
  IconButton,
  CircularProgress,
  MenuItem,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import { 
  Visibility, 
  VisibilityOff, 
  LockOutlined, 
  EmailOutlined,
  PersonOutlined,
} from '@mui/icons-material';
import { AuthContext } from '../App';
import { authService } from '../services';

const validationSchema = Yup.object({
  firstName: Yup.string().required('First name is required'),
  lastName: Yup.string().required('Last name is required'),
  email: Yup.string().email('Invalid email address').required('Email is required'),
  password: Yup.string()
    .min(8, 'Password must be at least 8 characters')
    .matches(/[a-z]/, 'Password must contain at least one lowercase letter')
    .matches(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .matches(/[0-9]/, 'Password must contain at least one number')
    .required('Password is required'),
  confirmPassword: Yup.string()
    .oneOf([Yup.ref('password'), null], 'Passwords must match')
    .required('Please confirm your password'),
  role: Yup.string().required('Role is required'),
});

const RegisterPage = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const { isAuthenticated } = useContext(AuthContext);
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const handleSubmit = async (values, { setSubmitting, setFieldError }) => {
    setError('');
    setSuccess('');
    
    try {
      const { confirmPassword, ...registrationData } = values;
      const result = await authService.register(registrationData);
      
      if (result?.token && result?.user) {
        setSuccess('Registration successful! Redirecting to login...');
        setTimeout(() => {
          navigate('/login', { 
            state: { 
              message: 'Registration successful! Please log in with your credentials.',
              email: values.email 
            } 
          });
        }, 2000);
      } else {
        setError('Registration failed. Please try again.');
      }
    } catch (err) {
      console.error('Registration error:', err);
      const errorMessage = err.response?.data?.message || 'Registration failed. Please try again.';
      
      // Set form field errors if available
      if (err.response?.data?.errors) {
        Object.entries(err.response.data.errors).forEach(([field, messages]) => {
          setFieldError(field, Array.isArray(messages) ? messages[0] : messages);
        });
      } else {
        setError(errorMessage);
      }
    } finally {
      setSubmitting(false);
    }
  };

  const formik = useFormik({
    initialValues: {
      firstName: '',
      lastName: '',
      email: '',
      password: '',
      confirmPassword: '',
      role: 'CUSTOMER',
    },
    validationSchema,
    onSubmit: handleSubmit,
  });

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const handleClickShowConfirmPassword = () => {
    setShowConfirmPassword(!showConfirmPassword);
  };

  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return (
    <Container component="main" maxWidth="sm">
      <Box
        sx={{
          marginTop: 4,
          marginBottom: 4,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Paper
          elevation={3}
          sx={{
            padding: isMobile ? 3 : 4,
            width: '100%',
            borderRadius: 2,
          }}
        >
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
            }}
          >
            <Typography component="h1" variant="h4" sx={{ mb: 1, fontWeight: 'bold' }}>
              Create Account
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 3, textAlign: 'center' }}>
              Join Stock Platform to manage your inventory and orders
            </Typography>

            {error && (
              <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
                {error}
              </Alert>
            )}

            {success && (
              <Alert severity="success" sx={{ width: '100%', mb: 2 }}>
                {success}
              </Alert>
            )}

            <Box component="form" onSubmit={formik.handleSubmit} sx={{ width: '100%' }}>
              <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  id="firstName"
                  label="First Name"
                  name="firstName"
                  autoComplete="given-name"
                  autoFocus
                  value={formik.values.firstName}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  error={formik.touched.firstName && Boolean(formik.errors.firstName)}
                  helperText={formik.touched.firstName && formik.errors.firstName}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <PersonOutlined color="action" />
                      </InputAdornment>
                    ),
                  }}
                />
                
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  id="lastName"
                  label="Last Name"
                  name="lastName"
                  autoComplete="family-name"
                  value={formik.values.lastName}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  error={formik.touched.lastName && Boolean(formik.errors.lastName)}
                  helperText={formik.touched.lastName && formik.errors.lastName}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <PersonOutlined color="action" />
                      </InputAdornment>
                    ),
                  }}
                />
              </Box>

              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                value={formik.values.email}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.touched.email && Boolean(formik.errors.email)}
                helperText={formik.touched.email && formik.errors.email}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <EmailOutlined color="action" />
                    </InputAdornment>
                  ),
                }}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                select
                id="role"
                label="Role"
                name="role"
                value={formik.values.role}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.touched.role && Boolean(formik.errors.role)}
                helperText={formik.touched.role && formik.errors.role}
              >
                <MenuItem value="CUSTOMER">Customer</MenuItem>
                <MenuItem value="USER">User</MenuItem>
                <MenuItem value="MANAGER">Manager</MenuItem>
              </TextField>

              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type={showPassword ? 'text' : 'password'}
                id="password"
                autoComplete="new-password"
                value={formik.values.password}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.touched.password && Boolean(formik.errors.password)}
                helperText={formik.touched.password && formik.errors.password}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <LockOutlined color="action" />
                    </InputAdornment>
                  ),
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        aria-label="toggle password visibility"
                        onClick={handleClickShowPassword}
                        edge="end"
                      >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                name="confirmPassword"
                label="Confirm Password"
                type={showConfirmPassword ? 'text' : 'password'}
                id="confirmPassword"
                autoComplete="new-password"
                value={formik.values.confirmPassword}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.touched.confirmPassword && Boolean(formik.errors.confirmPassword)}
                helperText={formik.touched.confirmPassword && formik.errors.confirmPassword}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <LockOutlined color="action" />
                    </InputAdornment>
                  ),
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        aria-label="toggle confirm password visibility"
                        onClick={handleClickShowConfirmPassword}
                        edge="end"
                      >
                        {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />

              <Button
                type="submit"
                fullWidth
                variant="contained"
                color="primary"
                disabled={formik.isSubmitting}
                sx={{
                  mt: 3,
                  mb: 2,
                  py: 1.5,
                  fontSize: '1rem',
                  fontWeight: 600,
                  borderRadius: 1,
                  textTransform: 'none',
                  '&:hover': {
                    boxShadow: 2,
                  },
                }}
              >
                {formik.isSubmitting ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  'Create Account'
                )}
              </Button>

              <Box sx={{ mt: 2, textAlign: 'center', width: '100%' }}>
                <Typography variant="body2" component="span" color="text.secondary">
                  Already have an account?{' '}
                  <Link 
                    component={RouterLink} 
                    to="/login" 
                    variant="body2"
                    sx={{
                      fontWeight: 600,
                      color: 'primary.main',
                      textDecoration: 'none',
                      '&:hover': {
                        textDecoration: 'underline',
                      },
                    }}
                  >
                    Sign in
                  </Link>
                </Typography>
              </Box>
            </Box>
          </Box>
        </Paper>
        
        <Box sx={{ mt: 3, textAlign: 'center' }}>
          <Typography variant="body2" color="text.secondary">
            © {new Date().getFullYear()} Stock Platform. All rights reserved.
          </Typography>
        </Box>
      </Box>
    </Container>
  );
};

export default RegisterPage;
