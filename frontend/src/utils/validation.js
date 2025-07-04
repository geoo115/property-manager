// Validation rules
export const validationRules = {
  required: (fieldName) => (value) => {
    if (!value || value.trim() === '') {
      return `${fieldName} is required`;
    }
    return '';
  },

  email: (value) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(value)) {
      return 'Please enter a valid email address';
    }
    return '';
  },

  minLength: (min) => (value) => {
    if (value && value.length < min) {
      return `Must be at least ${min} characters long`;
    }
    return '';
  },

  maxLength: (max) => (value) => {
    if (value && value.length > max) {
      return `Must be no more than ${max} characters long`;
    }
    return '';
  },

  password: (value) => {
    if (!value) return 'Password is required';
    if (value.length < 8) return 'Password must be at least 8 characters long';
    if (!/(?=.*[a-z])/.test(value)) return 'Password must contain at least one lowercase letter';
    if (!/(?=.*[A-Z])/.test(value)) return 'Password must contain at least one uppercase letter';
    if (!/(?=.*\d)/.test(value)) return 'Password must contain at least one number';
    return '';
  },

  confirmPassword: (password) => (value) => {
    if (value !== password) {
      return 'Passwords do not match';
    }
    return '';
  },

  phone: (value) => {
    const phoneRegex = /^[+]?[1-9][\d]{0,15}$/;
    if (value && !phoneRegex.test(value.replace(/\s/g, ''))) {
      return 'Please enter a valid phone number';
    }
    return '';
  },

  username: (value) => {
    if (!value) return 'Username is required';
    if (value.length < 3) return 'Username must be at least 3 characters long';
    if (!/^[a-zA-Z0-9_]+$/.test(value)) return 'Username can only contain letters, numbers, and underscores';
    return '';
  },
};

// Helper function to create form validation rules
export const createFormValidation = (fields) => {
  const rules = {};
  
  Object.keys(fields).forEach(fieldName => {
    const fieldRules = fields[fieldName];
    rules[fieldName] = fieldRules.map(rule => {
      if (typeof rule === 'function') {
        return rule;
      }
      if (typeof rule === 'string') {
        return validationRules[rule];
      }
      if (typeof rule === 'object' && rule.type) {
        return validationRules[rule.type](rule.value);
      }
      return () => '';
    });
  });
  
  return rules;
};
