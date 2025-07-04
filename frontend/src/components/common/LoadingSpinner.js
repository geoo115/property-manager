import React from 'react';
import PropTypes from 'prop-types';

const LoadingSpinner = ({ 
  size = 'medium', 
  variant = 'primary',
  text = 'Loading...',
  className = '',
  ...props 
}) => {
  const spinnerClass = `loading-spinner loading-spinner-${size} loading-spinner-${variant} ${className}`.trim();
  
  return (
    <div 
      className={spinnerClass}
      role="status"
      aria-label={text}
      {...props}
    >
      <div className="spinner" aria-hidden="true" />
      {text && (
        <div className="spinner-text" aria-live="polite">
          {text}
        </div>
      )}
    </div>
  );
};

LoadingSpinner.propTypes = {
  size: PropTypes.oneOf(['small', 'medium', 'large']),
  variant: PropTypes.oneOf(['primary', 'secondary', 'success', 'danger', 'warning']),
  text: PropTypes.string,
  className: PropTypes.string,
};

export default LoadingSpinner;
