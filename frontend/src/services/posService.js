import axios from 'axios';
import orderService from './orderService';
import inventoryService from './inventoryService';

const API_BASE_URL = '/api/v1';

// Create axios instance with default configuration
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);



const posService = {
  // POS Transaction Management (using consolidated order endpoints)
  // Create POS transaction
  createTransaction: async (transactionData, storeId) => {
    return await orderService.createPOSOrder({
      ...transactionData,
      storeId,
      source: 'POS',
      type: 'POS_TRANSACTION'
    });
  },

  // Process quick sale (consolidated endpoint)
  processQuickSale: async (items, paymentData, storeId, customerInfo = null) => {
    const transactionData = {
      items,
      payments: [paymentData],
      storeId,
      customerInfo,
      source: 'POS',
      type: 'QUICK_SALE',
      status: 'COMPLETED'
    };

    return await orderService.processQuickPOSTransaction(transactionData);
  },

  // Get POS transactions for current session
  getSessionTransactions: async (storeId, sessionId = null) => {
    const params = {
      source: 'POS',
      storeId
    };

    if (sessionId) {
      params.sessionId = sessionId;
    } else {
      // Get today's transactions
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      params.startDate = today.toISOString();
    }

    return await orderService.getPOSOrders(storeId, params);
  },

  // Get transaction by ID
  getTransaction: async (transactionId) => {
    return await orderService.getOrderById(transactionId);
  },

  // Void/cancel POS transaction
  voidTransaction: async (transactionId, reason) => {
    return await orderService.cancelOrder(transactionId, reason);
  },

  // Process refund for POS transaction
  processRefund: async (transactionId, refundData) => {
    return await orderService.createReturn(transactionId, {
      ...refundData,
      type: 'REFUND',
      source: 'POS'
    });
  },

  // Inventory Management for POS (using consolidated endpoints)
  // Check product availability
  checkProductAvailability: async (items, storeId) => {
    return await inventoryService.checkPOSAvailability(items, storeId);
  },

  // Search products for POS
  searchProducts: async (query, storeId, params = {}) => {
    return await inventoryService.searchInventory(query, {
      ...params,
      storeId,
      availableOnly: true
    });
  },

  // Get product by SKU/barcode
  getProductBySKU: async (sku, storeId) => {
    const product = await inventoryService.getInventoryItemBySKU(sku);
    if (product && storeId) {
      // Check availability at specific store
      const availability = await inventoryService.getInventoryByLocation(storeId, {
        productId: product.productId
      });
      return {
        ...product,
        storeAvailability: availability
      };
    }
    return product;
  },

  // Reserve inventory for POS transaction
  reserveInventory: async (items, orderId, storeId) => {
    return await inventoryService.reserveForPOS(items, orderId, storeId);
  },

  // Complete inventory deduction after payment
  completeInventoryDeduction: async (reservationId, deductionData) => {
    return await inventoryService.completePOSDeduction(reservationId, deductionData);
  },

  // Direct inventory deduction for quick sales
  directInventoryDeduction: async (items, storeId, reason = 'POS Sale') => {
    return await inventoryService.directPOSDeduction(items, storeId, reason);
  },

  // Payment Processing
  // Process payment
  processPayment: async (orderId, paymentData) => {
    return await orderService.processPOSPayment(orderId, {
      ...paymentData,
      source: 'POS'
    });
  },

  // Process multiple payments for split payment
  processMultiplePayments: async (orderId, payments) => {
    const results = [];
    for (const payment of payments) {
      try {
        const result = await this.processPayment(orderId, payment);
        results.push(result);
      } catch (error) {
        // If any payment fails, we should handle it appropriately
        throw new Error(`Payment failed: ${error.message}`);
      }
    }
    return results;
  },

  // Validate payment amount
  validatePayment: async (orderTotal, paymentAmount, paymentMethod) => {
    // Basic validation - can be enhanced with payment gateway integration
    if (paymentAmount < orderTotal) {
      throw new Error('Payment amount is less than order total');
    }

    const change = paymentAmount - orderTotal;
    return {
      valid: true,
      change,
      paymentMethod
    };
  },

  // POS Session Management
  // Start POS session
  startSession: async (terminalId, storeId, staffId) => {
    const sessionData = {
      terminalId,
      storeId,
      staffId,
      startTime: new Date().toISOString(),
      status: 'ACTIVE'
    };

    // This would typically create a session record
    // For now, we'll store it locally and return session info
    const sessionId = `pos_session_${Date.now()}`;
    localStorage.setItem('posSession', JSON.stringify({
      ...sessionData,
      sessionId
    }));

    return {
      sessionId,
      ...sessionData
    };
  },

  // End POS session
  endSession: async (sessionId) => {
    const session = this.getCurrentSession();
    if (!session || session.sessionId !== sessionId) {
      throw new Error('Invalid session');
    }

    // Get session summary
    const summary = await this.getSessionSummary(sessionId);

    // Update session status
    const updatedSession = {
      ...session,
      endTime: new Date().toISOString(),
      status: 'CLOSED',
      summary
    };

    localStorage.removeItem('posSession');
    return updatedSession;
  },

  // Get current session
  getCurrentSession: () => {
    const sessionData = localStorage.getItem('posSession');
    return sessionData ? JSON.parse(sessionData) : null;
  },

  // Get session summary
  getSessionSummary: async (sessionId) => {
    const session = this.getCurrentSession();
    if (!session) {
      throw new Error('No active session');
    }

    const transactions = await this.getSessionTransactions(session.storeId, sessionId);
    
    const summary = {
      totalTransactions: transactions.length,
      totalSales: transactions.reduce((sum, tx) => sum + (tx.total || 0), 0),
      totalItems: transactions.reduce((sum, tx) => sum + (tx.items?.length || 0), 0),
      paymentMethods: {},
      sessionStart: session.startTime,
      sessionEnd: new Date().toISOString()
    };

    // Calculate payment method breakdown
    transactions.forEach(tx => {
      if (tx.payments) {
        tx.payments.forEach(payment => {
          const method = payment.method || 'UNKNOWN';
          summary.paymentMethods[method] = (summary.paymentMethods[method] || 0) + payment.amount;
        });
      }
    });

    return summary;
  },

  // Customer Management for POS
  // Search customers
  searchCustomers: async (query) => {
    // This would integrate with customer service
    // For now, return a simple search
    return await api.get('/customers/search', {
      params: { q: query, limit: 10 }
    }).then(response => response.data);
  },

  // Create quick customer
  createQuickCustomer: async (customerData) => {
    return await api.post('/customers', {
      ...customerData,
      source: 'POS'
    }).then(response => response.data);
  },

  // Get customer by phone/email
  getCustomerByContact: async (contact) => {
    return await api.get('/customers/search', {
      params: { contact }
    }).then(response => response.data);
  },

  // Discounts and Promotions
  // Apply discount to transaction
  applyDiscount: async (transactionId, discountData) => {
    return await api.post(`/orders/${transactionId}/discounts`, discountData);
  },

  // Remove discount from transaction
  removeDiscount: async (transactionId, discountId) => {
    return await api.delete(`/orders/${transactionId}/discounts/${discountId}`);
  },

  // Get available promotions
  getAvailablePromotions: async (storeId) => {
    return await api.get('/promotions', {
      params: { storeId, active: true }
    }).then(response => response.data);
  },

  // Apply promotion code
  applyPromotionCode: async (transactionId, promoCode) => {
    return await api.post(`/orders/${transactionId}/promotions`, { code: promoCode });
  },

  // POS Analytics and Reporting
  // Get daily sales summary
  getDailySummary: async (storeId, date = new Date()) => {
    const dateStr = date.toISOString().split('T')[0];
    return await orderService.getOrderAnalytics({
      storeId,
      source: 'POS',
      startDate: dateStr,
      endDate: dateStr,
      groupBy: 'day'
    });
  },

  // Get hourly sales breakdown
  getHourlySales: async (storeId, date = new Date()) => {
    const dateStr = date.toISOString().split('T')[0];
    return await orderService.getOrderAnalytics({
      storeId,
      source: 'POS',
      startDate: dateStr,
      endDate: dateStr,
      groupBy: 'hour'
    });
  },

  // Get top selling products
  getTopSellingProducts: async (storeId, params = {}) => {
    return await orderService.getOrderAnalytics({
      storeId,
      source: 'POS',
      groupBy: 'product',
      ...params
    });
  },

  // Get staff performance
  getStaffPerformance: async (storeId, staffId, params = {}) => {
    return await orderService.getOrderAnalytics({
      storeId,
      staffId,
      source: 'POS',
      ...params
    });
  },

  // Receipt and Printing
  // Generate receipt data
  generateReceipt: async (transactionId) => {
    const transaction = await this.getTransaction(transactionId);
    const store = transaction.storeInfo || {};
    
    return {
      transactionId,
      timestamp: transaction.createdAt,
      store: {
        name: store.name,
        address: store.address,
        phone: store.phone,
        taxId: store.taxId
      },
      items: transaction.items?.map(item => ({
        name: item.productName,
        sku: item.sku,
        quantity: item.quantity,
        unitPrice: item.unitPrice,
        total: item.total
      })) || [],
      subtotal: transaction.subtotal,
      tax: transaction.tax,
      total: transaction.total,
      payments: transaction.payments,
      staff: transaction.staffInfo?.name,
      customer: transaction.customerInfo
    };
  },

  // Print receipt (would integrate with printer)
  printReceipt: async (receiptData) => {
    // This would integrate with a receipt printer
    // For now, we'll just return the formatted receipt
    console.log('Printing receipt:', receiptData);
    return { printed: true, receiptData };
  },

  // Email receipt
  emailReceipt: async (transactionId, emailAddress) => {
    const receipt = await this.generateReceipt(transactionId);
    return await api.post(`/orders/${transactionId}/email-receipt`, {
      email: emailAddress,
      receipt
    });
  },

  // POS Configuration
  // Get POS terminal configuration
  getTerminalConfig: async (terminalId) => {
    return await api.get(`/pos/terminals/${terminalId}/config`);
  },

  // Update terminal configuration
  updateTerminalConfig: async (terminalId, config) => {
    return await api.put(`/pos/terminals/${terminalId}/config`, config);
  },

  // Utility Functions
  // Calculate order totals
  calculateTotals: (items, taxRate = 0, discounts = []) => {
    const subtotal = items.reduce((sum, item) => {
      return sum + (item.quantity * item.unitPrice);
    }, 0);

    const discountAmount = discounts.reduce((sum, discount) => {
      if (discount.type === 'PERCENTAGE') {
        return sum + (subtotal * discount.value / 100);
      } else {
        return sum + discount.value;
      }
    }, 0);

    const discountedSubtotal = Math.max(0, subtotal - discountAmount);
    const tax = discountedSubtotal * taxRate;
    const total = discountedSubtotal + tax;

    return {
      subtotal,
      discountAmount,
      discountedSubtotal,
      tax,
      total
    };
  },

  // Format currency
  formatCurrency: (amount, currency = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency
    }).format(amount);
  },

  // Generate barcode for transaction
  generateTransactionBarcode: (transactionId) => {
    // Simple barcode generation - would integrate with proper barcode library
    return `TXN${transactionId.toString().padStart(10, '0')}`;
  },

  // Validate POS operation permissions
  validatePermissions: (operation, userRole) => {
    const permissions = {
      'ADMIN': ['*'],
      'MANAGER': ['create_transaction', 'void_transaction', 'process_refund', 'apply_discount'],
      'STAFF': ['create_transaction', 'apply_discount'],
      'CASHIER': ['create_transaction']
    };

    const userPermissions = permissions[userRole] || [];
    return userPermissions.includes('*') || userPermissions.includes(operation);
  }
};

export default posService;
