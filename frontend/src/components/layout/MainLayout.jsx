import React, { useState, useMemo } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { 
  Box, 
  CssBaseline, 
  useTheme, 
  useMediaQuery, 
  Drawer, 
  Toolbar,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Divider,
  Typography,
  Collapse,
  IconButton
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  Inventory as InventoryIcon,
  ShoppingCart as OrdersIcon,
  Category as CategoryIcon,
  ChevronLeft as ChevronLeftIcon,
  ChevronRight as ChevronRightIcon,
  ExpandLess,
  ExpandMore,
  Settings as SettingsIcon,
  People as PeopleIcon,
  Group as GroupIcon,
  BarChart as ReportsIcon,
  Store as SuppliersIcon,
  PointOfSale as POSIcon,
  AdminPanelSettings as AdminIcon,
} from '@mui/icons-material';
import Header from './Header';
import { AuthContext } from '../../App';

const drawerWidth = 240;

const menuItems = [
  { text: 'Dashboard', icon: <DashboardIcon />, path: '/' },
  { 
    text: 'Products', 
    icon: <CategoryIcon />, 
    path: '/products',
    children: [
      { text: 'All Products', path: '/products' },
      { text: 'Categories', path: '/products/categories' },
    ]
  },
  { 
    text: 'Inventory', 
    icon: <InventoryIcon />, 
    path: '/inventory',
    children: [
      { text: 'Overview', path: '/inventory' },
      { text: 'Stock Levels', path: '/inventory/levels' },
      { text: 'Stock Transfers', path: '/inventory/transfers' },
      { text: 'Stock Adjustments', path: '/inventory/adjustments' },
    ]
  },
  { 
    text: 'Orders', 
    icon: <OrdersIcon />, 
    path: '/orders',
    children: [
      { text: 'All Orders', path: '/orders' },
      { text: 'Create Order', path: '/orders/new' },
      { text: 'Returns', path: '/returns' },
    ]
  },
  { 
    text: 'Suppliers', 
    icon: <SuppliersIcon />, 
    path: '/suppliers',
    managerOnly: true
  },
  { 
    text: 'Point of Sale', 
    icon: <POSIcon />, 
    path: '/pos',
    managerOnly: true
  },
  { 
    text: 'Customers', 
    icon: <PeopleIcon />, 
    path: '/customers'
  },
  { 
    text: 'Users', 
    icon: <GroupIcon />, 
    path: '/users',
    adminOnly: true
  },
  { 
    text: 'Admin Panel', 
    icon: <AdminIcon />, 
    path: '/admin',
    adminOnly: true
  },
  { 
    text: 'Reports', 
    icon: <ReportsIcon />, 
    path: '/reports',
    children: [
      { text: 'Sales Reports', path: '/reports/sales' },
      { text: 'Inventory Reports', path: '/reports/inventory' },
      { text: 'Customer Reports', path: '/reports/customers' },
    ]
  },
  { 
    text: 'Settings', 
    icon: <SettingsIcon />, 
    path: '/settings',
    adminOnly: true,
    children: [
      { text: 'General', path: '/settings/general' },
      { text: 'Users', path: '/settings/users' },
      { text: 'Roles & Permissions', path: '/settings/roles' },
    ]
  },
];

const Sidebar = ({ drawerWidth, mobileOpen, onClose, isMobile }) => {
  const theme = useTheme();
  const location = useLocation();
  const navigate = useNavigate();
  const { user } = React.useContext(AuthContext);
  const [expandedItems, setExpandedItems] = React.useState({});

  // Filter menu items based on user role
  const filteredMenuItems = useMemo(() => {
    return menuItems.filter(item => {
      if (item.adminOnly && user?.role !== 'ADMIN') {
        return false;
      }
      if (item.managerOnly && !['ADMIN', 'MANAGER'].includes(user?.role)) {
        return false;
      }
      return true;
    });
  }, [user?.role]);

  const handleItemClick = (item) => {
    if (item.children) {
      setExpandedItems(prev => ({
        ...prev,
        [item.text]: !prev[item.text]
      }));
    } else {
      navigate(item.path);
      if (isMobile) {
        onClose();
      }
    }
  };

  const isItemActive = (itemPath) => {
    return location.pathname === itemPath || 
           location.pathname.startsWith(`${itemPath}/`);
  };

  const drawer = (
    <div>
      <Toolbar sx={{ 
        display: 'flex', 
        alignItems: 'center',
        justifyContent: 'space-between',
        px: 2,
        minHeight: '64px !important',
        borderBottom: '1px solid',
        borderColor: 'divider',
      }}>
        <Typography 
          variant="h6" 
          noWrap 
          component="div"
          sx={{
            background: 'linear-gradient(45deg, #3f51b5 30%, #2196f3 90%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            backgroundClip: 'text',
            textFillColor: 'transparent',
            fontWeight: 700,
          }}
        >
          StockPro
        </Typography>
        {!isMobile && (
          <IconButton onClick={onClose} size="small">
            {theme.direction === 'ltr' ? <ChevronLeftIcon /> : <ChevronRightIcon />}
          </IconButton>
        )}
      </Toolbar>
      <Divider />
      <List>
        {filteredMenuItems.map((item) => (
          <React.Fragment key={item.text}>
            <ListItem 
              disablePadding 
              sx={{ 
                display: 'block',
                '&:hover': {
                  backgroundColor: 'action.hover',
                },
                bgcolor: isItemActive(item.path) ? 'action.selected' : 'transparent',
              }}
            >
              <ListItemButton
                onClick={() => handleItemClick(item)}
                sx={{
                  minHeight: 48,
                  justifyContent: 'initial',
                  px: 2.5,
                }}
              >
                <ListItemIcon
                  sx={{
                    minWidth: 0,
                    mr: 3,
                    justifyContent: 'center',
                    color: isItemActive(item.path) ? 'primary.main' : 'text.secondary',
                  }}
                >
                  {React.cloneElement(item.icon, {
                    sx: { fontSize: 20 },
                  })}
                </ListItemIcon>
                <ListItemText 
                  primary={item.text} 
                  primaryTypographyProps={{
                    fontWeight: isItemActive(item.path) ? 600 : 400,
                    color: isItemActive(item.path) ? 'text.primary' : 'text.secondary',
                  }}
                />
                {item.children && (
                  expandedItems[item.text] ? <ExpandLess /> : <ExpandMore />
                )}
              </ListItemButton>
            </ListItem>
            {item.children && (
              <Collapse in={expandedItems[item.text] || isItemActive(item.path)} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                  {item.children.map((child) => (
                    <ListItemButton
                      key={child.text}
                      onClick={() => handleItemClick(child)}
                      sx={{
                        pl: 8,
                        minHeight: 40,
                        bgcolor: isItemActive(child.path) ? 'action.selected' : 'transparent',
                        '&:hover': {
                          bgcolor: 'action.hover',
                        },
                      }}
                    >
                      <ListItemText 
                        primary={child.text} 
                        primaryTypographyProps={{
                          fontSize: '0.875rem',
                          color: isItemActive(child.path) ? 'primary.main' : 'text.secondary',
                          fontWeight: isItemActive(child.path) ? 600 : 400,
                        }}
                      />
                    </ListItemButton>
                  ))}
                </List>
              </Collapse>
            )}
          </React.Fragment>
        ))}
      </List>
    </div>
  );

  if (isMobile) {
    return (
      <Drawer
        variant="temporary"
        open={mobileOpen}
        onClose={onClose}
        ModalProps={{
          keepMounted: true, // Better open performance on mobile.
        }}
        sx={{
          display: { xs: 'block', sm: 'none' },
          '& .MuiDrawer-paper': { 
            boxSizing: 'border-box', 
            width: drawerWidth,
            borderRight: 'none',
            boxShadow: theme.shadows[8],
          },
        }}
      >
        {drawer}
      </Drawer>
    );
  }

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: drawerWidth,
          boxSizing: 'border-box',
          borderRight: 'none',
          backgroundColor: 'background.paper',
        },
        display: { xs: 'none', sm: 'block' },
      }}
      open
    >
      {drawer}
    </Drawer>
  );
};

const MainLayout = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const [mobileOpen, setMobileOpen] = useState(false);
  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  
  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <CssBaseline />
      <Header onMenuToggle={handleDrawerToggle} />
      <Sidebar 
        drawerWidth={drawerWidth}
        mobileOpen={mobileOpen}
        onClose={handleDrawerToggle}
        isMobile={isMobile}
      />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: { xs: 2, sm: 3 },
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          marginTop: '64px', // Height of the header
          backgroundColor: 'background.default',
          minHeight: 'calc(100vh - 64px)',
        }}
      >
        <Outlet />
      </Box>
    </Box>
  );
};

export default MainLayout;
