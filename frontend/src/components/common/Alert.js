import React, { useState, useEffect, useCallback } from 'react';
import PropTypes from 'prop-types';

const Alert = ({ 
  children, 
  type = 'info', 
  title,
  dismissible = false,
  onDismiss,
  autoDismiss = false,
  autoDismissDelay = 5000,
  className = '',
  ...props 
}) => {
  const [isVisible, setIsVisible] = useState(true);

  const handleDismiss = useCallback(() => {
    setIsVisible(false);
    if (onDismiss) {
      onDismiss();
    }
  }, [onDismiss]);

  useEffect(() => {
    if (autoDismiss && isVisible) {
      const timer = setTimeout(() => {
        handleDismiss();
      }, autoDismissDelay);

      return () => clearTimeout(timer);
    }
  }, [autoDismiss, autoDismissDelay, isVisible, handleDismiss]);

  if (!isVisible) {
    return null;
  }

  const alertClass = `alert alert-${type} ${className}`.trim();
  
  const getIcon = () => {
    switch (type) {
      case 'success':
        return '✅';
      case 'warning':
        return '⚠️';
      case 'error':
        return '❌';
      case 'info':
      default:
        return 'ℹ️';
    }
  };

  return (
    <div 
      className={alertClass} 
      role="alert"
      aria-live="polite"
      {...props}
    >
      <div className="alert-icon">
        <span role="img" aria-hidden="true">
          {getIcon()}
        </span>
      </div>
      
      <div className="alert-content">
        {title && <h4 className="alert-title">{title}</h4>}
        <div className="alert-message">{children}</div>
      </div>
      
      {dismissible && (
        <button
          type="button"
          className="alert-close"
          onClick={handleDismiss}
          aria-label="Dismiss alert"
        >
          <span aria-hidden="true">×</span>
        </button>
      )}
    </div>
  );
};

Alert.propTypes = {
  children: PropTypes.node.isRequired,
  type: PropTypes.oneOf(['info', 'success', 'warning', 'error']),
  title: PropTypes.string,
  dismissible: PropTypes.bool,
  onDismiss: PropTypes.func,
  autoDismiss: PropTypes.bool,
  autoDismissDelay: PropTypes.number,
  className: PropTypes.string,
};

export default Alert;
