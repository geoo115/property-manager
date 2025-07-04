import React, { forwardRef } from 'react';
import PropTypes from 'prop-types';

const FormSelect = forwardRef(({ 
  label,
  name,
  value,
  onChange,
  onBlur,
  options = [],
  placeholder,
  error,
  required = false,
  disabled = false,
  multiple = false,
  size = 'medium',
  variant = 'default',
  icon,
  className = '',
  ...props 
}, ref) => {
  const selectId = `select-${name}`;
  const errorId = `${selectId}-error`;
  
  const selectClass = `form-select form-select-${size} form-select-${variant} ${error ? 'error' : ''} ${className}`.trim();
  const containerClass = `form-group ${icon ? 'form-group-with-icon' : ''}`;
  
  const renderIcon = () => {
    if (!icon) return null;
    
    const iconElement = typeof icon === 'string' ? (
      <span role="img" aria-hidden="true" className="form-select-icon">
        {icon}
      </span>
    ) : (
      <span className="form-select-icon">{icon}</span>
    );
    
    return iconElement;
  };

  return (
    <div className={containerClass}>
      {label && (
        <label htmlFor={selectId} className="form-label">
          {label}
          {required && <span className="required-asterisk" aria-label="required">*</span>}
        </label>
      )}
      
      <div className="form-select-wrapper">
        {icon && renderIcon()}
        
        <select
          ref={ref}
          id={selectId}
          name={name}
          value={value}
          onChange={onChange}
          onBlur={onBlur}
          required={required}
          disabled={disabled}
          multiple={multiple}
          className={selectClass}
          aria-describedby={error ? errorId : undefined}
          aria-invalid={error ? 'true' : 'false'}
          {...props}
        >
          {placeholder && (
            <option value="" disabled>
              {placeholder}
            </option>
          )}
          
          {options.map((option) => (
            <option
              key={option.value}
              value={option.value}
              disabled={option.disabled}
            >
              {option.label}
            </option>
          ))}
        </select>
        
        <div className="form-select-arrow" aria-hidden="true">
          <span>▼</span>
        </div>
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

FormSelect.propTypes = {
  label: PropTypes.string,
  name: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([PropTypes.string, PropTypes.number, PropTypes.array]),
  onChange: PropTypes.func,
  onBlur: PropTypes.func,
  options: PropTypes.arrayOf(
    PropTypes.shape({
      value: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
      label: PropTypes.string.isRequired,
      disabled: PropTypes.bool,
    })
  ),
  placeholder: PropTypes.string,
  error: PropTypes.string,
  required: PropTypes.bool,
  disabled: PropTypes.bool,
  multiple: PropTypes.bool,
  size: PropTypes.oneOf(['small', 'medium', 'large']),
  variant: PropTypes.oneOf(['default', 'outlined', 'filled']),
  icon: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
  className: PropTypes.string,
};

FormSelect.displayName = 'FormSelect';

export default FormSelect;
