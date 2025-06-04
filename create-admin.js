const axios = require('axios');

// Configuration
const API_BASE_URL = 'http://localhost:8080/api/v1';
const ADMIN_CREDENTIALS = {
  email: 'admin@admin.com',
  password: 'Admin@123',
  firstName: 'Admin',
  lastName: 'User',
  role: 'ADMIN'
};

// Helper function to log with timestamp
const log = (message) => {
  const timestamp = new Date().toISOString();
  console.log(`[${timestamp}] ${message}`);
};

// Register admin user
async function createAdminUser() {
  try {
    log('Attempting to create admin user...');
    
    // First, check if user already exists by trying to log in
    try {
      log('Checking if admin user already exists...');
      const loginResponse = await axios.post(`${API_BASE_URL}/auth/login`, {
        email: ADMIN_CREDENTIALS.email,
        password: ADMIN_CREDENTIALS.password
      });
      
      log('Admin user already exists. Login successful!');
      log('Token:', loginResponse.data.token);
      return;
    } catch (loginError) {
      // If login fails, it means user doesn't exist or password is wrong
      if (loginError.response?.status === 401) {
        log('Admin user does not exist or password is incorrect. Proceeding with registration...');
      } else {
        throw loginError;
      }
    }
    
    // Register the admin user
    log('Registering new admin user...');
    const registerResponse = await axios.post(`${API_BASE_URL}/auth/register`, {
      email: ADMIN_CREDENTIALS.email,
      password: ADMIN_CREDENTIALS.password,
      firstName: ADMIN_CREDENTIALS.firstName,
      lastName: ADMIN_CREDENTIALS.lastName,
      role: ADMIN_CREDENTIALS.role
    });
    
    log('Admin user created successfully!');
    log('Registration response:', registerResponse.data);
    
    // Log in to get the token
    log('Logging in with new admin credentials...');
    const loginResponse = await axios.post(`${API_BASE_URL}/auth/login`, {
      email: ADMIN_CREDENTIALS.email,
      password: ADMIN_CREDENTIALS.password
    });
    
    log('Login successful!');
    log('Token:', loginResponse.data.token);
    
  } catch (error) {
    console.error('Error creating admin user:');
    
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      console.error('Response data:', error.response.data);
      console.error('Response status:', error.response.status);
      console.error('Response headers:', error.response.headers);
    } else if (error.request) {
      // The request was made but no response was received
      console.error('No response received:', error.request);
    } else {
      // Something happened in setting up the request that triggered an Error
      console.error('Error setting up request:', error.message);
    }
    
    process.exit(1);
  }
}

// Run the script
createAdminUser().catch(console.error);
