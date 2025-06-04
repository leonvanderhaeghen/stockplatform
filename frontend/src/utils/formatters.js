// Format currency with fallback for undefined/null values
export const formatCurrency = (value) => {
  if (value == null) return '$0.00';
  const numValue = typeof value === 'string' 
    ? parseFloat(value.replace(/[^0-9.-]+/g,"")) 
    : value;
  return `$${Number(numValue || 0).toFixed(2)}`;
};
