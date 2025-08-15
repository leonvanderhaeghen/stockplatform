import React, { useState } from 'react';
import {
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  Box,
  IconButton,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Button,
  Chip,
  Alert,
  CircularProgress,
  Paper,
  Divider
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  TrendingUp,
  Inventory,
  ShoppingCart,
  People,
  Store,
  Warning,
  Refresh,
  Add,
  ListAlt,
  Assessment,
  Settings,
  Notifications,
  TrendingDown,
  CheckCircle,
  AccessTime,
  LocalShipping
} from '@mui/icons-material';
import { useQuery } from '@tanstack/react-query';
import { useAuth } from '../../hooks/useAuth';
import { useNavigate } from 'react-router-dom';
import adminService from '../../services/adminService';
import orderService from '../../services/orderService';
import inventoryService from '../../services/inventoryService';
import userService from '../../services/userService';
import { format } from 'date-fns';

const StatCard = ({ title, value, icon, color, trend, alert }) => {
  const IconComponent = icon;
  return (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Box>
            <Typography color="text.secondary" gutterBottom variant="h6">
              {title}
            </Typography>
            <Typography variant="h4" component="div">
              {value}
            </Typography>
            {trend !== undefined && (
              <Box display="flex" alignItems="center" mt={1}>
                {trend > 0 ? (
                  <TrendingUp sx={{ fontSize: 16, color: 'success.main', mr: 0.5 }} />
                ) : trend < 0 ? (
                  <TrendingDown sx={{ fontSize: 16, color: 'error.main', mr: 0.5 }} />
                ) : null}
                <Typography variant="caption" color={trend > 0 ? 'success.main' : trend < 0 ? 'error.main' : 'text.secondary'}>
                  {trend > 0 ? '+' : ''}{trend}%
                </Typography>
              </Box>
            )}
          </Box>
          <IconComponent sx={{ 
            fontSize: 40, 
            color: alert ? color : color,
            opacity: alert ? 1 : 0.8
          }} />
        </Box>
        {alert && value > 0 && (
          <Alert severity="warning" sx={{ mt: 1, py: 0 }}>
            <Typography variant="caption">Requires attention</Typography>
          </Alert>
        )}
      </CardContent>
    </Card>
  );
};

const DashboardPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [refreshKey, setRefreshKey] = useState(0);

  // Query for dashboard analytics
  const { data: analytics, isLoading: analyticsLoading, error: analyticsError } = useQuery({
    queryKey: ['dashboard-analytics', user?.role, refreshKey],
    queryFn: async () => {
      if (user?.role === 'ADMIN') {
        return await adminService.getDashboardAnalytics();
      } else {
        // For STAFF and CUSTOMER, get basic stats
        return {
          orders: await orderService.getOrderAnalytics({ period: 'today' }),
          inventory: await inventoryService.getInventoryStats()
        };
      }
    },
    enabled: !!user,
    staleTime: 30000, // 30 seconds
  });

  // Query for recent orders
  const { data: recentOrders, isLoading: ordersLoading } = useQuery({
    queryKey: ['recent-orders', user?.role, refreshKey],
    queryFn: async () => {
      if (user?.role === 'CUSTOMER') {
        return await orderService.getUserOrders({ limit: 5 });
      } else {
        return await orderService.getOrders({ limit: 5, sortBy: 'createdAt', sortOrder: 'desc' });
      }
    },
    enabled: !!user,
    staleTime: 30000,
  });

  // Query for low stock items (ADMIN/STAFF only)
  const { data: lowStockItems, isLoading: stockLoading } = useQuery({
    queryKey: ['low-stock-items', refreshKey],
    queryFn: () => inventoryService.getLowStockItems({ limit: 5 }),
    enabled: !!user && (user.role === 'ADMIN' || user.role === 'STAFF'),
    staleTime: 60000, // 1 minute
  });

  // Query for user profile (CUSTOMER only)
  const { data: userProfile } = useQuery({
    queryKey: ['user-profile'],
    queryFn: () => userService.getProfile(),
    enabled: !!user && user.role === 'CUSTOMER',
    staleTime: 300000, // 5 minutes
  });

  const handleRefresh = () => {
    setRefreshKey(prev => prev + 1);
  };

  const getStatsForRole = () => {
    if (!analytics) return [];

    if (user?.role === 'ADMIN') {
      return [
        { 
          title: 'Total Revenue', 
          value: `$${(analytics.revenue?.total || 0).toLocaleString()}`, 
          icon: TrendingUp, 
          color: '#4caf50',
          trend: analytics.revenue?.trend || 0
        },
        { 
          title: 'Total Orders', 
          value: (analytics.orders?.total || 0).toLocaleString(), 
          icon: ShoppingCart, 
          color: '#2196f3',
          trend: analytics.orders?.trend || 0
        },
        { 
          title: 'Active Customers', 
          value: (analytics.customers?.active || 0).toLocaleString(), 
          icon: People, 
          color: '#ff9800',
          trend: analytics.customers?.trend || 0
        },
        { 
          title: 'Total Products', 
          value: (analytics.products?.total || 0).toLocaleString(), 
          icon: Inventory, 
          color: '#9c27b0',
          trend: analytics.products?.trend || 0
        },
        { 
          title: 'Low Stock Items', 
          value: lowStockItems?.length || 0, 
          icon: Warning, 
          color: '#f44336',
          alert: true
        },
        { 
          title: 'Active Stores', 
          value: (analytics.stores?.active || 0).toLocaleString(), 
          icon: Store, 
          color: '#00bcd4'
        }
      ];
    } else if (user?.role === 'STAFF') {
      return [
        { 
          title: "Today's Orders", 
          value: analytics.orders?.today || 0, 
          icon: ShoppingCart, 
          color: '#2196f3'
        },
        { 
          title: "Today's Revenue", 
          value: `$${(analytics.orders?.todayRevenue || 0).toLocaleString()}`, 
          icon: TrendingUp, 
          color: '#4caf50'
        },
        { 
          title: 'Pending Orders', 
          value: analytics.orders?.pending || 0, 
          icon: AccessTime, 
          color: '#ff9800'
        },
        { 
          title: 'Low Stock Items', 
          value: lowStockItems?.length || 0, 
          icon: Warning, 
          color: '#f44336',
          alert: true
        }
      ];
    } else {
      return [
        { 
          title: 'Recent Orders', 
          value: recentOrders?.data?.length || 0, 
          icon: ShoppingCart, 
          color: '#2196f3'
        },
        { 
          title: 'Total Spent', 
          value: `$${(userProfile?.totalSpent || 0).toLocaleString()}`, 
          icon: TrendingUp, 
          color: '#4caf50'
        },
        { 
          title: 'Loyalty Points', 
          value: userProfile?.loyaltyPoints || 0, 
          icon: People, 
          color: '#ff9800'
        },
        { 
          title: 'Active Orders', 
          value: recentOrders?.data?.filter(order => 
            ['PENDING', 'PROCESSING', 'SHIPPED'].includes(order.status)
          ).length || 0, 
          icon: LocalShipping, 
          color: '#9c27b0'
        }
      ];
    }
  };

  const getQuickActions = () => {
    if (user?.role === 'ADMIN') {
      return [
        { label: 'Create Product', icon: Add, action: () => navigate('/products/new'), color: 'primary' },
        { label: 'Manage Users', icon: People, action: () => navigate('/admin'), color: 'secondary' },
        { label: 'View Reports', icon: Assessment, action: () => navigate('/admin'), color: 'info' },
        { label: 'System Settings', icon: Settings, action: () => navigate('/admin'), color: 'warning' }
      ];
    } else if (user?.role === 'STAFF') {
      return [
        { label: 'New Order', icon: Add, action: () => navigate('/orders/new'), color: 'primary' },
        { label: 'POS Terminal', icon: Store, action: () => navigate('/pos'), color: 'secondary' },
        { label: 'Inventory Check', icon: Inventory, action: () => navigate('/inventory'), color: 'info' },
        { label: 'View Orders', icon: ListAlt, action: () => navigate('/orders'), color: 'success' }
      ];
    } else {
      return [
        { label: 'Browse Products', icon: Inventory, action: () => navigate('/products'), color: 'primary' },
        { label: 'View Orders', icon: ShoppingCart, action: () => navigate('/orders'), color: 'secondary' },
        { label: 'Update Profile', icon: People, action: () => navigate('/profile'), color: 'info' }
      ];
    }
  };

  const getOrderStatusColor = (status) => {
    switch (status) {
      case 'COMPLETED': return 'success';
      case 'PENDING': return 'warning';
      case 'PROCESSING': return 'info';
      case 'SHIPPED': return 'primary';
      case 'CANCELLED': return 'error';
      default: return 'default';
    }
  };

  const getWelcomeMessage = () => {
    const hour = new Date().getHours();
    let greeting = 'Good morning';
    if (hour >= 12 && hour < 17) greeting = 'Good afternoon';
    else if (hour >= 17) greeting = 'Good evening';
    
    return `${greeting}, ${user?.firstName || user?.username || 'User'}!`;
  };

  const statsCards = getStatsForRole();
  const quickActions = getQuickActions();

  if (analyticsLoading) {
    return (
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (analyticsError) {
    return (
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Alert severity="error" sx={{ mb: 2 }}>
          Failed to load dashboard data. Please try refreshing the page.
        </Alert>
      </Container>
    );
  }

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      {/* Welcome Section */}
      <Paper sx={{ p: 3, mb: 3, bgcolor: 'primary.main', color: 'primary.contrastText' }}>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box display="flex" alignItems="center">
            <DashboardIcon sx={{ fontSize: 40, mr: 2 }} />
            <Box>
              <Typography variant="h4" gutterBottom>
                {getWelcomeMessage()}
              </Typography>
              <Typography variant="subtitle1">
                Welcome to your StockPlatform dashboard for {format(new Date(), 'MMMM d, yyyy')}.
              </Typography>
            </Box>
          </Box>
          <IconButton color="inherit" onClick={handleRefresh} disabled={analyticsLoading}>
            <Refresh />
          </IconButton>
        </Box>
      </Paper>

      {/* Stats Cards */}
      <Grid container spacing={3} mb={4}>
        {statsCards.map((stat, index) => (
          <Grid item xs={12} sm={6} md={4} lg={2} key={index}>
            <StatCard
              title={stat.title}
              value={stat.value}
              icon={stat.icon}
              color={stat.color}
              trend={stat.trend}
              alert={stat.alert}
            />
          </Grid>
        ))}
      </Grid>

      {/* Main Content Grid */}
      <Grid container spacing={3}>
        {/* Recent Activity */}
        <Grid item xs={12} md={8}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                <Typography variant="h6">
                  Recent {user?.role === 'CUSTOMER' ? 'Orders' : 'Activity'}
                </Typography>
                <Chip 
                  icon={<Notifications />} 
                  label={`${recentOrders?.data?.length || 0} items`} 
                  size="small" 
                  color="primary" 
                />
              </Box>
              
              {ordersLoading ? (
                <Box display="flex" justifyContent="center" py={3}>
                  <CircularProgress size={24} />
                </Box>
              ) : recentOrders?.data?.length > 0 ? (
                <List>
                  {recentOrders.data.slice(0, 5).map((order, index) => (
                    <React.Fragment key={order.id}>
                      <ListItem 
                        button 
                        onClick={() => navigate(`/orders/${order.id}`)}
                        sx={{ px: 0 }}
                      >
                        <ListItemIcon>
                          <ShoppingCart color="primary" />
                        </ListItemIcon>
                        <ListItemText
                          primary={`Order #${order.id?.slice(-8) || 'N/A'}`}
                          secondary={`${format(new Date(order.createdAt), 'MMM d, yyyy HH:mm')} • $${(order.total || 0).toFixed(2)}`}
                        />
                        <Chip 
                          label={order.status} 
                          size="small" 
                          color={getOrderStatusColor(order.status)}
                        />
                      </ListItem>
                      {index < recentOrders.data.length - 1 && <Divider />}
                    </React.Fragment>
                  ))}
                </List>
              ) : (
                <Typography variant="body2" color="text.secondary" textAlign="center" py={3}>
                  No recent orders found.
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>

        {/* Quick Actions & Alerts */}
        <Grid item xs={12} md={4}>
          <Grid container spacing={2}>
            {/* Quick Actions */}
            <Grid item xs={12}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    Quick Actions
                  </Typography>
                  <Grid container spacing={1}>
                    {quickActions.map((action, index) => {
                      const IconComponent = action.icon;
                      return (
                        <Grid item xs={6} key={index}>
                          <Button
                            fullWidth
                            variant="outlined"
                            color={action.color}
                            startIcon={<IconComponent />}
                            onClick={action.action}
                            sx={{ 
                              justifyContent: 'flex-start',
                              textAlign: 'left',
                              py: 1
                            }}
                          >
                            {action.label}
                          </Button>
                        </Grid>
                      );
                    })}
                  </Grid>
                </CardContent>
              </Card>
            </Grid>

            {/* Low Stock Alert (ADMIN/STAFF only) */}
            {(user?.role === 'ADMIN' || user?.role === 'STAFF') && (
              <Grid item xs={12}>
                <Card>
                  <CardContent>
                    <Box display="flex" alignItems="center" mb={2}>
                      <Warning sx={{ mr: 1, color: 'warning.main' }} />
                      <Typography variant="h6">
                        Low Stock Alert
                      </Typography>
                    </Box>
                    
                    {stockLoading ? (
                      <Box display="flex" justifyContent="center" py={2}>
                        <CircularProgress size={24} />
                      </Box>
                    ) : lowStockItems?.length > 0 ? (
                      <List dense>
                        {lowStockItems.slice(0, 3).map((item, index) => (
                          <ListItem key={item.id} sx={{ px: 0 }}>
                            <ListItemIcon>
                              <Warning color="warning" />
                            </ListItemIcon>
                            <ListItemText
                              primary={item.productName}
                              secondary={`Stock: ${item.quantity}`}
                            />
                          </ListItem>
                        ))}
                        {lowStockItems.length > 3 && (
                          <ListItem button onClick={() => navigate('/inventory')} sx={{ px: 0 }}>
                            <ListItemText
                              primary={`View ${lowStockItems.length - 3} more items`}
                              sx={{ textAlign: 'center', color: 'primary.main' }}
                            />
                          </ListItem>
                        )}
                      </List>
                    ) : (
                      <Box display="flex" alignItems="center" py={2}>
                        <CheckCircle sx={{ mr: 1, color: 'success.main' }} />
                        <Typography variant="body2" color="success.main">
                          All items are well stocked
                        </Typography>
                      </Box>
                    )}
                  </CardContent>
                </Card>
              </Grid>
            )}
          </Grid>
        </Grid>
      </Grid>
    </Container>
  );
};

export default DashboardPage;
