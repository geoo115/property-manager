import React from 'react';
import PropTypes from 'prop-types';
import useAuth from '../../hooks/useAuth';
import Button from './Button';

/**
 * RoleBasedActions component - Shows action buttons based on user permissions
 */
const RoleBasedActions = ({ 
  resource, 
  item, 
  onView, 
  onEdit, 
  onDelete, 
  onApprove, 
  onReject,
  customActions = [],
  className = '' 
}) => {
  const { hasUserPermission } = useAuth();

  const actions = [];

  // Standard CRUD actions
  if (hasUserPermission(resource, 'read') && onView) {
    actions.push({
      name: 'View',
      icon: '👁️',
      onClick: () => onView(item),
      variant: 'secondary',
      size: 'small',
    });
  }

  if (hasUserPermission(resource, 'update') && onEdit) {
    actions.push({
      name: 'Edit',
      icon: '✏️',
      onClick: () => onEdit(item),
      variant: 'primary',
      size: 'small',
    });
  }

  if (hasUserPermission(resource, 'delete') && onDelete) {
    actions.push({
      name: 'Delete',
      icon: '🗑️',
      onClick: () => onDelete(item),
      variant: 'danger',
      size: 'small',
    });
  }

  // Workflow actions
  if (hasUserPermission(resource, 'approve') && onApprove) {
    actions.push({
      name: 'Approve',
      icon: '✅',
      onClick: () => onApprove(item),
      variant: 'success',
      size: 'small',
    });
  }

  if (hasUserPermission(resource, 'reject') && onReject) {
    actions.push({
      name: 'Reject',
      icon: '❌',
      onClick: () => onReject(item),
      variant: 'danger',
      size: 'small',
    });
  }

  // Custom actions
  customActions.forEach(action => {
    if (hasUserPermission(resource, action.permission)) {
      actions.push(action);
    }
  });

  if (actions.length === 0) {
    return null;
  }

  return (
    <div className={`role-based-actions ${className}`}>
      {actions.map((action, index) => (
        <Button
          key={index}
          variant={action.variant}
          size={action.size}
          onClick={action.onClick}
          title={action.name}
          className="action-btn"
        >
          {action.icon && <span className="action-icon">{action.icon}</span>}
          {action.name}
        </Button>
      ))}
    </div>
  );
};

RoleBasedActions.propTypes = {
  resource: PropTypes.string.isRequired,
  item: PropTypes.object.isRequired,
  onView: PropTypes.func,
  onEdit: PropTypes.func,
  onDelete: PropTypes.func,
  onApprove: PropTypes.func,
  onReject: PropTypes.func,
  customActions: PropTypes.arrayOf(PropTypes.shape({
    name: PropTypes.string.isRequired,
    icon: PropTypes.string,
    onClick: PropTypes.func.isRequired,
    variant: PropTypes.string,
    size: PropTypes.string,
    permission: PropTypes.string.isRequired,
  })),
  className: PropTypes.string,
};

export default RoleBasedActions;
