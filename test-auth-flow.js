const { chromium } = require('playwright');

// Configuration
const FRONTEND_URL = 'http://localhost:3001';
const TEST_CREDENTIALS = {
  email: 'admin@example.com',
  password: 'admin123'
};

// Helper function to log with timestamp
const log = (message) => {
  const timestamp = new Date().toISOString();
  console.log(`[${timestamp}] ${message}`);
};

(async () => {
  // Launch the browser
  const browser = await chromium.launch({ 
    headless: false, 
    slowMo: 100,
    args: ['--disable-web-security'] // Disable CORS for testing
  });
  
  const context = await browser.newContext({
    viewport: { width: 1280, height: 800 },
    ignoreHTTPSErrors: true
  });
  
  const page = await context.newPage();

  try {
    // Test 1: Access login page
    log('Test 1: Accessing login page...');
    await page.goto(`${FRONTEND_URL}/login`, { waitUntil: 'networkidle' });
    
    // Wait for the login form to be visible
    await page.waitForSelector('input[name="email"]', { state: 'visible', timeout: 10000 });
    log('✅ Login page loaded successfully');
    
    // Test 2: Fill and submit login form
    log('Test 2: Submitting login form...');
    await page.fill('input[name="email"]', TEST_CREDENTIALS.email);
    await page.fill('input[name="password"]', TEST_CREDENTIALS.password);
    
    // Click the login button and wait for navigation
    await Promise.all([
      page.waitForNavigation({ waitUntil: 'networkidle' }),
      page.click('button[type="submit"]')
    ]);
    
    // Test 3: Verify successful login
    const currentUrl = page.url();
    log(`Current URL after login: ${currentUrl}`);
    
    if (currentUrl.includes('/dashboard') || currentUrl.endsWith('/') || currentUrl === FRONTEND_URL) {
      log('✅ Login successful!');
      
      // Test 4: Check for authentication token
      const token = await page.evaluate(() => localStorage.getItem('token'));
      if (token) {
        log('✅ Authentication token found in localStorage');
      } else {
        log('❌ No authentication token found in localStorage');
      }
      
      // Test 5: Access protected route
      log('Test 5: Accessing protected route...');
      await page.goto(`${FRONTEND_URL}/dashboard`, { waitUntil: 'networkidle' });
      
      // Check for dashboard content
      const dashboardTitle = await page.$('h1, h2, h3, [role="heading"]');
      if (dashboardTitle) {
        const titleText = await dashboardTitle.textContent();
        log(`✅ Successfully accessed protected route. Page title: ${titleText}`);
      } else {
        log('❌ Failed to verify protected route content');
      }
      
      // Test 6: Check for user session
      const userData = await page.evaluate(() => {
        try {
          return JSON.parse(localStorage.getItem('user') || 'null');
        } catch (e) {
          return null;
        }
      });
      
      if (userData && userData.email) {
        log(`✅ User session found for: ${userData.email}`);
        log(`✅ User role: ${userData.role || 'No role specified'}`);
      } else {
        log('❌ No valid user session found');
      }
      
    } else {
      log('❌ Login failed or unexpected redirect');
      
      // Check for error messages
      const errorMessage = await page.$('.MuiAlert-message, .error-message, [role="alert"]');
      if (errorMessage) {
        const errorText = await errorMessage.textContent();
        log(`Error message: ${errorText}`);
      }
      
      // Check for form validation errors
      const formErrors = await page.$$eval('.MuiFormHelperText-root.Mui-error', 
        errors => errors.map(e => e.textContent)
      );
      
      if (formErrors.length > 0) {
        log('Form validation errors:');
        formErrors.forEach((error, i) => log(`  ${i + 1}. ${error}`));
      }
    }
    
    // Take a screenshot for debugging
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const screenshotPath = `test-result-${timestamp}.png`;
    await page.screenshot({ path: screenshotPath, fullPage: true });
    log(`Screenshot saved as ${screenshotPath}`);
    
  } catch (error) {
    log(`❌ Test failed with error: ${error.message}`);
    
    // Take a screenshot on error
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const errorScreenshotPath = `test-error-${timestamp}.png`;
    await page.screenshot({ path: errorScreenshotPath, fullPage: true });
    log(`Error screenshot saved as ${errorScreenshotPath}`);
    
  } finally {
    // Close the browser
    await browser.close();
    log('Test completed. Browser closed.');
  }
})();
