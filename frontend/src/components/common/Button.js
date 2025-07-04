import React from 'react';
import PropTypes from 'prop-types';
import './Button.css';

const Button = ({ 
  children, 
  onClick, 
  type = 'button', 
  variant = 'primary',
  size = 'medium',
  disabled = false,
  loading = false,
  className = '',
  icon,
  iconPosition = 'left',
  fullWidth = false,
  ...props 
}) => {
  const baseClass = 'btn';
  const variantClass = `btn--${variant}`;
  const sizeClass = `btn--${size}`;
  const disabledClass = disabled || loading ? 'btn--disabled' : '';
  const fullWidthClass = fullWidth ? 'btn--full-width' : '';
  const iconClass = icon ? `btn--with-icon btn--icon-${iconPosition}` : '';
  const loadingClass = loading ? 'btn--loading' : '';
  
  const buttonClass = `${baseClass} ${variantClass} ${sizeClass} ${disabledClass} ${fullWidthClass} ${iconClass} ${loadingClass} ${className}`.trim();
  
  const renderIcon = () => {
    if (!icon) return null;
    
    const iconElement = typeof icon === 'string' ? (
      <span role="img" aria-hidden="true" className="btn__icon">
        {icon}
      </span>
    ) : (
      <span className="btn__icon">{icon}</span>
    );
    
    return iconElement;
  };

  const renderSpinner = () => (
    <svg className="btn__spinner" viewBox="0 0 24 24" fill="none">
      <circle
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeDasharray="32"
        strokeDashoffset="32"
      >
        <animate
          attributeName="stroke-dasharray"
          dur="2s"
          values="0 32;16 16;0 32;0 32"
          repeatCount="indefinite"
        />
        <animate
          attributeName="stroke-dashoffset"
          dur="2s"
          values="0;-16;-32;-32"
          repeatCount="indefinite"
        />
      </circle>
    </svg>
  );
  
  return (
    <button
      type={type}
      onClick={onClick}
      className={buttonClass}
      disabled={disabled || loading}
      aria-disabled={disabled || loading}
      {...props}
    >
      {loading && renderSpinner()}
      {!loading && iconPosition === 'left' && renderIcon()}
      <span className="btn__content">{children}</span>
      {!loading && iconPosition === 'right' && renderIcon()}
    </button>
  );
};

Button.propTypes = {
  children: PropTypes.node.isRequired,
  onClick: PropTypes.func,
  type: PropTypes.oneOf(['button', 'submit', 'reset']),
  variant: PropTypes.oneOf(['primary', 'secondary', 'success', 'danger', 'warning', 'ghost', 'outline']),
  size: PropTypes.oneOf(['small', 'medium', 'large']),
  disabled: PropTypes.bool,
  loading: PropTypes.bool,
  className: PropTypes.string,
  icon: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
  iconPosition: PropTypes.oneOf(['left', 'right']),
  fullWidth: PropTypes.bool,
};

export default Button;
