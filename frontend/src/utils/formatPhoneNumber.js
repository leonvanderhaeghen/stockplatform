import { parsePhoneNumber as parseLibPhoneNumber, isValidPhoneNumber } from 'libphonenumber-js';

/**
 * Format a phone number for display
 * @param {string} phoneNumber - The phone number to format
 * @param {boolean} partial - Whether to allow partial formatting (while typing)
 * @returns {string} Formatted phone number
 */
const formatPhoneNumber = (phoneNumber, partial = false) => {
  if (!phoneNumber) return '';
  
  // Remove all non-digit characters
  const cleaned = ('' + phoneNumber).replace(/\D/g, '');
  
  // For empty or very short numbers, return as is
  if (cleaned.length === 0) return '';
  if (cleaned.length < 3) return `(${cleaned}`;
  
  try {
    // Try to parse as US number
    const phoneNumberObj = parseLibPhoneNumber(cleaned, 'US');
    
    if (phoneNumberObj && (partial || isValidPhoneNumber(phoneNumber, 'US'))) {
      // Format as US number: (123) 456-7890
      const nationalNumber = phoneNumberObj.formatNational();
      return nationalNumber;
    }
  } catch (e) {
    // If parsing fails, fall back to basic formatting
  }
  
  // Fallback formatting for non-US or invalid numbers
  const match = cleaned.match(/^(\d{1,3})(\d{0,3})(\d{0,4})$/);
  if (match) {
    const formatted = match[1] + (match[2] ? `-${match[2]}` : '') + (match[3] ? `-${match[3]}` : '');
    return formatted;
  }
  
  return cleaned;
};

/**
 * Parse a phone number string into a PhoneNumber object
 * @param {string} phoneNumber - The phone number to parse
 * @returns {import('libphonenumber-js').PhoneNumber} Parsed phone number
 */
const parsePhoneNumber = (phoneNumber) => {
  if (!phoneNumber) return null;
  
  try {
    // Remove all non-digit characters except leading +
    const cleaned = phoneNumber.replace(/[^\d+]/g, '');
    
    // If the number starts with +, assume it's in international format
    if (cleaned.startsWith('+')) {
      return parseLibPhoneNumber(cleaned);
    }
    
    // Otherwise, assume it's a US number
    return parseLibPhoneNumber(cleaned, 'US');
  } catch (e) {
    return null;
  }
};

const phoneUtils = {
  formatPhoneNumber,
  parsePhoneNumber,
};

export { formatPhoneNumber, parsePhoneNumber };
export default phoneUtils;
