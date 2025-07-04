import React from 'react';
import './FormField.css';

const FormField = ({ 
  type = 'text',
  label,
  name,
  value,
  onChange,
  onBlur,
  placeholder,
  required = false,
  disabled = false,
  error,
  hint,
  options = [],
  rows = 3,
  multiple = false,
  accept,
  className = '',
  ...props
}) => {
  const baseId = `field-${name}`;
  
  const renderInput = () => {
    const commonProps = {
      id: baseId,
      name,
      value: value || '',
      onChange,
      onBlur,
      placeholder,
      required,
      disabled,
      className: `form-input ${error ? 'form-input--error' : ''}`,
      'aria-describedby': error ? `${baseId}-error` : hint ? `${baseId}-hint` : undefined,
      ...props
    };

    switch (type) {
      case 'textarea':
        return <textarea {...commonProps} rows={rows} />;
      
      case 'select':
        return (
          <select {...commonProps} multiple={multiple}>
            {placeholder && !multiple && (
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
        );
      
      case 'checkbox':
        return (
          <div className="checkbox-wrapper">
            <input
              {...commonProps}
              type="checkbox"
              checked={value || false}
              className="form-checkbox"
            />
            <span className="checkbox-indicator"></span>
          </div>
        );
      
      case 'radio':
        return (
          <div className="radio-group">
            {options.map((option) => (
              <label key={option.value} className="radio-label">
                <input
                  type="radio"
                  name={name}
                  value={option.value}
                  checked={value === option.value}
                  onChange={onChange}
                  disabled={disabled}
                  className="form-radio"
                />
                <span className="radio-indicator"></span>
                <span className="radio-text">{option.label}</span>
              </label>
            ))}
          </div>
        );
      
      case 'file':
        return (
          <div className="file-input-wrapper">
            <input
              {...commonProps}
              type="file"
              accept={accept}
              multiple={multiple}
              className="form-file-input"
            />
            <div className="file-input-display">
              <span className="file-input-icon">📁</span>
              <span className="file-input-text">
                {value ? 'File selected' : 'Choose file...'}
              </span>
            </div>
          </div>
        );
      
      default:
        return <input {...commonProps} type={type} />;
    }
  };

  return (
    <div className={`form-field ${className} ${error ? 'form-field--error' : ''}`}>
      {label && type !== 'checkbox' && (
        <label htmlFor={baseId} className="form-label">
          {label}
          {required && <span className="form-required">*</span>}
        </label>
      )}
      
      {type === 'checkbox' ? (
        <label htmlFor={baseId} className="form-label form-label--checkbox">
          {renderInput()}
          <span className="checkbox-label-text">
            {label}
            {required && <span className="form-required">*</span>}
          </span>
        </label>
      ) : (
        renderInput()
      )}
      
      {hint && !error && (
        <div id={`${baseId}-hint`} className="form-hint">
          {hint}
        </div>
      )}
      
      {error && (
        <div id={`${baseId}-error`} className="form-error" role="alert">
          {error}
        </div>
      )}
    </div>
  );
};

export default FormField;
