import React from 'react';
import PropTypes from 'prop-types';

const Card = ({ 
  children, 
  className = '', 
  title, 
  subtitle,
  headerAction,
  footer,
  variant = 'default',
  elevation = 'default',
  ...props 
}) => {
  const cardClass = `card card-${variant} card-elevation-${elevation} ${className}`.trim();
  
  return (
    <div className={cardClass} {...props}>
      {(title || subtitle || headerAction) && (
        <div className="card-header">
          <div className="card-header-content">
            {title && <h3 className="card-title">{title}</h3>}
            {subtitle && <p className="card-subtitle">{subtitle}</p>}
          </div>
          {headerAction && (
            <div className="card-header-action">
              {headerAction}
            </div>
          )}
        </div>
      )}
      <div className="card-body">{children}</div>
      {footer && <div className="card-footer">{footer}</div>}
    </div>
  );
};

Card.propTypes = {
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
  title: PropTypes.string,
  subtitle: PropTypes.string,
  headerAction: PropTypes.node,
  footer: PropTypes.node,
  variant: PropTypes.oneOf(['default', 'outlined', 'elevated', 'flat']),
  elevation: PropTypes.oneOf(['none', 'sm', 'default', 'md', 'lg', 'xl']),
};

export default Card;
