import { Link } from 'react-router-dom';
import { Button, Container, Typography, Box } from '@mui/material';
import { Lock as LockIcon } from '@mui/icons-material';

const UnauthorizedPage = () => {
  return (
    <Container component="main" maxWidth="md">
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          textAlign: 'center',
        }}
      >
        <LockIcon color="error" sx={{ fontSize: 80, mb: 2 }} />
        <Typography component="h1" variant="h3" gutterBottom>
          Access Denied
        </Typography>
        <Typography variant="h6" color="textSecondary" paragraph>
          You don't have permission to access this page.
        </Typography>
        <Typography variant="body1" color="textSecondary" paragraph>
          Please contact your administrator if you believe this is an error.
        </Typography>
        <Button
          variant="contained"
          color="primary"
          component={Link}
          to="/"
          sx={{ mt: 3 }}
        >
          Return to Home
        </Button>
      </Box>
    </Container>
  );
};

export default UnauthorizedPage;
