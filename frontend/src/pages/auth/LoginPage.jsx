import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useFormik } from 'formik';
import * as yup from 'yup';
import {
  Container,
  Paper,
  Box,
  TextField,
  Button,
  Typography,
  Alert,
  CircularProgress,
  Divider,
} from '@mui/material';
import { useSnackbar } from 'notistack';
import { useAuth } from '../../hooks/useAuth';

const validationSchema = yup.object({
  email: yup
    .string('Enter your email')
    .email('Enter a valid email')
    .required('Email is required'),
  password: yup
    .string('Enter your password')
    .min(6, 'Password should be of minimum 6 characters length')
    .required('Password is required'),
});

const LoginPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const { enqueueSnackbar } = useSnackbar();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const from = location.state?.from?.pathname || '/dashboard';

  const formik = useFormik({
    initialValues: {
      email: '',
      password: '',
    },
    validationSchema: validationSchema,
    onSubmit: async (values) => {
      setIsLoading(true);
      setError('');

      try {
        const result = await login(values.email, values.password);
        
        if (result.success) {
          enqueueSnackbar('Login successful!', { variant: 'success' });
          navigate(from, { replace: true });
        } else {
          setError(result.error || 'Login failed');
        }
      } catch (err) {
        setError('An unexpected error occurred');
        console.error('Login error:', err);
      } finally {
        setIsLoading(false);
      }
    },
  });

  return (
    <Container component="main" maxWidth="sm">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Paper
          elevation={3}
          sx={{
            padding: 4,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            width: '100%',
          }}
        >
          <Typography component="h1" variant="h4" gutterBottom>
            StockPlatform
          </Typography>
          <Typography variant="h6" color="text.secondary" gutterBottom>
            Sign in to your account
          </Typography>

          <Box component="form" onSubmit={formik.handleSubmit} sx={{ mt: 1, width: '100%' }}>
            {error && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {error}
              </Alert>
            )}

            <TextField
              margin="normal"
              required
              fullWidth
              id="email"
              label="Email Address"
              name="email"
              autoComplete="email"
              autoFocus
              value={formik.values.email}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.email && Boolean(formik.errors.email)}
              helperText={formik.touched.email && formik.errors.email}
            />
            
            <TextField
              margin="normal"
              required
              fullWidth
              name="password"
              label="Password"
              type="password"
              id="password"
              autoComplete="current-password"
              value={formik.values.password}
              onChange={formik.handleChange}
              onBlur={formik.handleBlur}
              error={formik.touched.password && Boolean(formik.errors.password)}
              helperText={formik.touched.password && formik.errors.password}
            />

            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
              disabled={isLoading}
              startIcon={isLoading ? <CircularProgress size={20} /> : null}
            >
              {isLoading ? 'Signing In...' : 'Sign In'}
            </Button>

            <Divider sx={{ my: 2 }}>
              <Typography variant="body2" color="text.secondary">
                OR
              </Typography>
            </Divider>

            <Button
              component={Link}
              to="/register"
              fullWidth
              variant="outlined"
              sx={{ mb: 2 }}
            >
              Create New Account
            </Button>

            <Box sx={{ textAlign: 'center', mt: 2 }}>
              <Link to="/forgot-password" style={{ textDecoration: 'none' }}>
                <Typography variant="body2" color="primary">
                  Forgot your password?
                </Typography>
              </Link>
            </Box>
          </Box>
        </Paper>

        {/* Demo credentials info */}
        <Paper sx={{ mt: 3, p: 2, width: '100%', bgcolor: 'grey.50' }}>
          <Typography variant="subtitle2" gutterBottom>
            Demo Credentials:
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Admin: admin@stockplatform.com / password123
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Staff: staff@stockplatform.com / password123
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Customer: customer@stockplatform.com / password123
          </Typography>
        </Paper>
      </Box>
    </Container>
  );
};

export default LoginPage;
