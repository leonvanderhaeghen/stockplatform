import React from 'react';
import { Chip, Tooltip, Box, Typography } from '@mui/material';
import { 
  CheckCircle as InStockIcon,
  Warning as LowStockIcon,
  Error as OutOfStockIcon,
  HelpOutline as UnknownIcon 
} from '@mui/icons-material';

const InventoryStatus = ({ 
  quantity, 
  lowStockThreshold = 10, 
  showQuantity = false, 
  size = 'small',
  variant = 'filled' 
}) => {
  // Handle cases where inventory data might not be available
  if (quantity === null || quantity === undefined) {
    return (
      <Tooltip title="Inventory data not available">
        <Chip
          icon={<UnknownIcon />}
          label="Unknown"
          color="default"
          size={size}
          variant={variant}
        />
      </Tooltip>
    );
  }

  const qty = Number(quantity) || 0;
  const isOutOfStock = qty === 0;
  const isLowStock = qty > 0 && qty <= lowStockThreshold;

  let color, icon, label, tooltipText;

  if (isOutOfStock) {
    color = 'error';
    icon = <OutOfStockIcon />;
    label = showQuantity ? `Out of Stock (${qty})` : 'Out of Stock';
    tooltipText = 'This item is currently out of stock';
  } else if (isLowStock) {
    color = 'warning';
    icon = <LowStockIcon />;
    label = showQuantity ? `Low Stock (${qty})` : 'Low Stock';
    tooltipText = `Low stock warning: Only ${qty} units remaining`;
  } else {
    color = 'success';
    icon = <InStockIcon />;
    label = showQuantity ? `In Stock (${qty})` : 'In Stock';
    tooltipText = `${qty} units available in stock`;
  }

  return (
    <Tooltip title={tooltipText}>
      <Chip
        icon={icon}
        label={label}
        color={color}
        size={size}
        variant={variant}
      />
    </Tooltip>
  );
};

// Simplified version that just shows quantity with color coding
export const InventoryQuantity = ({ quantity, lowStockThreshold = 10 }) => {
  const qty = Number(quantity) || 0;
  const isOutOfStock = qty === 0;
  const isLowStock = qty > 0 && qty <= lowStockThreshold;

  let color;
  if (isOutOfStock) {
    color = 'error';
  } else if (isLowStock) {
    color = 'warning';
  } else {
    color = 'success';
  }

  return (
    <Chip
      label={qty}
      color={color}
      size="small"
      variant="filled"
    />
  );
};

// Component for displaying detailed inventory information
export const InventoryDetails = ({ inventoryData, loading = false }) => {
  if (loading) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Typography variant="body2" color="text.secondary">
          Loading inventory...
        </Typography>
      </Box>
    );
  }

  if (!inventoryData) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Typography variant="body2" color="text.secondary">
          No inventory data
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
      <InventoryStatus 
        quantity={inventoryData.quantity} 
        lowStockThreshold={inventoryData.low_stock_threshold}
        showQuantity={true}
      />
      {inventoryData.reserved_quantity > 0 && (
        <Chip
          label={`${inventoryData.reserved_quantity} Reserved`}
          color="info"
          size="small"
          variant="outlined"
        />
      )}
    </Box>
  );
};

export default InventoryStatus;
