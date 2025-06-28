// Format currency with fallback for undefined/null values
export const formatCurrency = (value) => {
  if (value == null) return '$0.00';
  const numValue = typeof value === 'string' 
    ? parseFloat(value.replace(/[^0-9.-]+/g,"")) 
    : value;
  return `$${Number(numValue || 0).toFixed(2)}`;
};

// Format date with fallback for undefined/null values
export const formatDate = (date) => {
  if (!date) return 'N/A';
  try {
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    if (isNaN(dateObj.getTime())) return 'Invalid Date';
    return dateObj.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (error) {
    return 'Invalid Date';
  }
};
