import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';

const FormInput = forwardRef(({ 
  label,
  name,
  type = 'text',
  placeholder,
  value,
  onChange,
  onBlur,
  error,
  required = false,
  disabled = false,
  readOnly = false,
  autoComplete,
  autoFocus = false,
  maxLength,
  minLength,
  pattern,
  size = 'medium',
  variant = 'default',
  icon,
  iconPosition = 'left',
  className = '',
  ...props 
}, ref) => {
  const inputId = `input-${name}`;
  const errorId = `${inputId}-error`;
  
  const inputClass = `form-input form-input-${size} form-input-${variant} ${error ? 'error' : ''} ${className}`.trim();
  const containerClass = `form-group ${icon ? 'form-group-with-icon' : ''} ${iconPosition === 'right' ? 'form-group-icon-right' : ''}`;
  
  const renderIcon = () => {
    if (!icon) return null;
    
    const iconElement = typeof icon === 'string' ? (
      <span role="img" aria-hidden="true" className="form-input-icon">
        {icon}
      </span>
    ) : (
      <span className="form-input-icon">{icon}</span>
    );
    
    return iconElement;
  };

  return (
    <div className={containerClass}>
      {label && (
        <label htmlFor={inputId} className="form-label">
          {label}
          {required && <span className="required-asterisk" aria-label="required">*</span>}
        </label>
      )}
      
      <div className="form-input-wrapper">
        {icon && iconPosition === 'left' && renderIcon()}
        
        <input
          ref={ref}
          id={inputId}
          name={name}
          type={type}
          value={value}
          onChange={onChange}
          onBlur={onBlur}
          placeholder={placeholder}
          required={required}
          disabled={disabled}
          readOnly={readOnly}
          autoComplete={autoComplete}
          autoFocus={autoFocus}
          maxLength={maxLength}
          minLength={minLength}
          pattern={pattern}
          className={inputClass}
          aria-describedby={error ? errorId : undefined}
          aria-invalid={error ? 'true' : 'false'}
          {...props}
        />
        
        {icon && iconPosition === 'right' && renderIcon()}
      </div>
      
      {error && (
        <div id={errorId} className="error-text" role="alert">
          <span role="img" aria-hidden="true">⚠️</span>
          {error}
        </div>
      )}
    </div>
  );
});

FormInput.propTypes = {
  label: PropTypes.string,
  name: PropTypes.string.isRequired,
  type: PropTypes.string,
  placeholder: PropTypes.string,
  value: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
  onChange: PropTypes.func,
  onBlur: PropTypes.func,
  error: PropTypes.string,
  required: PropTypes.bool,
  disabled: PropTypes.bool,
  readOnly: PropTypes.bool,
  autoComplete: PropTypes.string,
  autoFocus: PropTypes.bool,
  maxLength: PropTypes.number,
  minLength: PropTypes.number,
  pattern: PropTypes.string,
  size: PropTypes.oneOf(['small', 'medium', 'large']),
  variant: PropTypes.oneOf(['default', 'outlined', 'filled']),
  icon: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
  iconPosition: PropTypes.oneOf(['left', 'right']),
  className: PropTypes.string,
};

FormInput.displayName = 'FormInput';

export default FormInput;
